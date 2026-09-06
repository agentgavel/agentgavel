package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/agentgavel/agentgavel/internal/publish"
	"github.com/agentgavel/agentgavel/internal/report"
	"github.com/agentgavel/agentgavel/internal/submit"
)

func runReport(args []string) int {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	root := fs.String("root", ".", "directory containing results/<run-id>/")
	asJSON := fs.Bool("json", false, "emit machine-readable scorecard JSON")
	doPublish := fs.Bool("publish", false, "write a dashboard entry and update index.json")
	doSign := fs.Bool("sign", false, "sign a dashboard entry with Ed25519 (ADR 013)")
	keyPath := fs.String("key", "", "path to base64 Ed25519 private key (for --sign)")
	keyID := fs.String("key-id", "", "registry key_id to embed (for --sign)")
	entryPath := fs.String("entry", "", "existing dashboard entry JSON to sign (optional; else build from run)")
	outPath := fs.String("out", "", "write signed entry to this path instead of stdout (sign-only)")
	dashboard := fs.String("dashboard", "dashboard", "dashboard root directory for --publish")
	framework := fs.String("framework", "", "framework display name for --publish/--sign")
	adapterName := fs.String("adapter-name", "", "adapter package/module name for --publish/--sign")
	tab := fs.String("tab", publish.TabUnratified, "leaderboard tab (unratified, or opt-in when signed per ADR 013)")
	fs.Usage = func() {
		_, _ = fmt.Fprintf(fs.Output(), `Usage: AgentGavel report [flags] <run-id|path>

Render a GSI scorecard from a completed run's results directory.
Looks for scorecard.json, otherwise computes GSI from summary.json.

With --publish, write <dashboard>/data/<run-id>.json and update index.json.
Default tab is unratified (sample=false). Opt-in publish requires --sign
with --key and --key-id (ADR 013); unsigned --tab opt-in is rejected.

With --sign alone, produce a signed entry JSON (ADR 013): either from
--entry PATH or from a run plus --framework/--adapter-name. Requires
--key and --key-id. Combine with --publish to write a signed Opt-in row.

Flags:
`)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return 2
	}

	// Sign-only (no publish): existing path.
	if *doSign && !*doPublish {
		return runReportSign(fs, *root, *keyPath, *keyID, *entryPath, *outPath, *framework, *adapterName)
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return 2
	}

	if *doPublish {
		if code := rejectPublishTab(*tab, *doSign); code != 0 {
			return code
		}
		if *doSign {
			if strings.TrimSpace(*keyPath) == "" || strings.TrimSpace(*keyID) == "" {
				fmt.Fprintf(os.Stderr, "report: --publish --sign requires --key and --key-id\n")
				return 2
			}
		}
	}

	dir, err := report.ResolveRunDir(*root, fs.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "report: %v\n", err)
		return 2
	}
	doc, err := report.Load(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "report: %v\n", err)
		return 1
	}

	if *doPublish {
		entry := publish.FromDocument(doc, *framework, *adapterName)
		pubTab := strings.TrimSpace(strings.ToLower(*tab))
		if pubTab == "" {
			pubTab = publish.TabUnratified
		}
		if pubTab == publish.TabOptIn {
			entry.Tab = publish.TabOptIn
			entry.Sample = false
			if err := signPublishEntry(&entry, *keyPath, *keyID); err != nil {
				fmt.Fprintf(os.Stderr, "report: publish: %v\n", err)
				return 1
			}
		}
		path, err := publish.Write(*dashboard, entry)
		if err != nil {
			fmt.Fprintf(os.Stderr, "report: publish: %v\n", err)
			return 1
		}
		enc, err := json.Marshal(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "report: %v\n", err)
			return 1
		}
		fmt.Println(string(enc))
		return 0
	}

	if *asJSON {
		out, err := report.FormatJSON(doc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "report: %v\n", err)
			return 1
		}
		fmt.Print(out)
		return 0
	}
	fmt.Print(report.FormatText(doc))
	return 0
}

func signPublishEntry(entry *publish.Entry, keyPath, keyID string) error {
	priv, err := submit.LoadPrivateKeyFile(keyPath)
	if err != nil {
		return err
	}
	raw, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal entry for sign: %w", err)
	}
	m, err := submit.EntryMapFromJSON(raw)
	if err != nil {
		return err
	}
	if err := submit.Sign(priv, keyID, m); err != nil {
		return err
	}
	kid, _ := m["key_id"].(string)
	sig, _ := m["signature"].(string)
	entry.KeyID = kid
	entry.Signature = sig
	return nil
}

func runReportSign(fs *flag.FlagSet, root, keyPath, keyID, entryPath, outPath, framework, adapterName string) int {
	if strings.TrimSpace(keyPath) == "" || strings.TrimSpace(keyID) == "" {
		fmt.Fprintf(os.Stderr, "report: --sign requires --key and --key-id\n")
		return 2
	}
	priv, err := submit.LoadPrivateKeyFile(keyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "report: sign: %v\n", err)
		return 1
	}

	var entry map[string]any
	if strings.TrimSpace(entryPath) != "" {
		raw, err := os.ReadFile(entryPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "report: sign: %v\n", err)
			return 1
		}
		entry, err = submit.EntryMapFromJSON(raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "report: sign: %v\n", err)
			return 1
		}
	} else {
		if fs.NArg() != 1 {
			fmt.Fprintf(os.Stderr, "report: --sign without --entry requires <run-id>\n")
			return 2
		}
		if strings.TrimSpace(framework) == "" || strings.TrimSpace(adapterName) == "" {
			fmt.Fprintf(os.Stderr, "report: --sign from a run requires --framework and --adapter-name\n")
			return 2
		}
		dir, err := report.ResolveRunDir(root, fs.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "report: %v\n", err)
			return 2
		}
		doc, err := report.Load(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "report: %v\n", err)
			return 1
		}
		pubEntry := publish.FromDocument(doc, framework, adapterName)
		pubEntry.Tab = publish.TabOptIn
		pubEntry.Sample = false
		raw, err := json.Marshal(pubEntry)
		if err != nil {
			fmt.Fprintf(os.Stderr, "report: sign: %v\n", err)
			return 1
		}
		entry, err = submit.EntryMapFromJSON(raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "report: sign: %v\n", err)
			return 1
		}
	}

	if err := submit.Sign(priv, keyID, entry); err != nil {
		fmt.Fprintf(os.Stderr, "report: sign: %v\n", err)
		return 1
	}
	out, err := submit.MarshalEntry(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "report: sign: %v\n", err)
		return 1
	}
	out = append(out, '\n')
	if strings.TrimSpace(outPath) != "" {
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil && filepath.Dir(outPath) != "." {
			fmt.Fprintf(os.Stderr, "report: sign: %v\n", err)
			return 1
		}
		if err := os.WriteFile(outPath, out, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "report: sign: %v\n", err)
			return 1
		}
		fmt.Println(outPath)
		return 0
	}
	fmt.Print(string(out))
	return 0
}

// rejectPublishTab enforces ADR 013: unsigned opt-in publish is rejected;
// signed opt-in (--sign) is allowed; unratified is always allowed.
func rejectPublishTab(tab string, signed bool) int {
	t := strings.TrimSpace(strings.ToLower(tab))
	if t == "" || t == publish.TabUnratified {
		return 0
	}
	if t == publish.TabOptIn {
		if signed {
			return 0
		}
		fmt.Fprintf(os.Stderr, "report: --tab opt-in requires --sign with --key and --key-id (ADR 013); use unratified for unsigned publishes\n")
		return 2
	}
	fmt.Fprintf(os.Stderr, "report: --tab %q invalid (want unratified or signed opt-in per ADR 013)\n", tab)
	return 2
}
