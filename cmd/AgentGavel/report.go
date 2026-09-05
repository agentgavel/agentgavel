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
	doPublish := fs.Bool("publish", false, "write an Unratified dashboard entry (ADR 006)")
	doSign := fs.Bool("sign", false, "sign a dashboard entry with Ed25519 (ADR 013); writes JSON to stdout")
	keyPath := fs.String("key", "", "path to base64 Ed25519 private key (for --sign)")
	keyID := fs.String("key-id", "", "registry key_id to embed (for --sign)")
	entryPath := fs.String("entry", "", "existing dashboard entry JSON to sign (optional; else build from run)")
	outPath := fs.String("out", "", "write signed entry to this path instead of stdout")
	dashboard := fs.String("dashboard", "dashboard", "dashboard root directory for --publish")
	framework := fs.String("framework", "", "framework display name for --publish/--sign")
	adapterName := fs.String("adapter-name", "", "adapter package/module name for --publish/--sign")
	tab := fs.String("tab", publish.TabUnratified, "leaderboard tab (v0.3: unratified only; opt-in rejected per ADR 006)")
	fs.Usage = func() {
		_, _ = fmt.Fprintf(fs.Output(), `Usage: AgentGavel report [flags] <run-id|path>

Render a GSI scorecard from a completed run's results directory.
Looks for scorecard.json, otherwise computes GSI from summary.json.

With --publish, write <dashboard>/data/<run-id>.json (tab=unratified,
sample=false) and update index.json. Opt-in publish is rejected until
v1.0 signatures land in later waves (ADR 006 addendum).

With --sign, produce a signed entry JSON (ADR 013): either from --entry
PATH or from a run plus --framework/--adapter-name. Requires --key and
--key-id.

Flags:
`)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *doSign {
		return runReportSign(fs, *root, *keyPath, *keyID, *entryPath, *outPath, *framework, *adapterName)
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return 2
	}

	if *doPublish {
		if code := rejectPublishTab(*tab); code != 0 {
			return code
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

// rejectPublishTab enforces ADR 006: report --publish writes Unratified only.
func rejectPublishTab(tab string) int {
	t := strings.TrimSpace(strings.ToLower(tab))
	if t == "" || t == publish.TabUnratified {
		return 0
	}
	if t == publish.TabOptIn {
		fmt.Fprintf(os.Stderr, "report: --tab opt-in rejected until v1.0 signatures (ADR 006); use unratified\n")
		return 2
	}
	fmt.Fprintf(os.Stderr, "report: --tab %q invalid (want unratified; opt-in rejected per ADR 006)\n", tab)
	return 2
}
