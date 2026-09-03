package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/agentgavel/agentgavel/internal/publish"
	"github.com/agentgavel/agentgavel/internal/report"
)

func runReport(args []string) int {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	root := fs.String("root", ".", "directory containing results/<run-id>/")
	asJSON := fs.Bool("json", false, "emit machine-readable scorecard JSON")
	doPublish := fs.Bool("publish", false, "write an Unratified dashboard entry (ADR 006)")
	dashboard := fs.String("dashboard", "dashboard", "dashboard root directory for --publish")
	framework := fs.String("framework", "", "framework display name for --publish")
	adapterName := fs.String("adapter-name", "", "adapter package/module name for --publish")
	tab := fs.String("tab", publish.TabUnratified, "leaderboard tab (v0.3: unratified only; opt-in rejected per ADR 006)")
	fs.Usage = func() {
		_, _ = fmt.Fprintf(fs.Output(), `Usage: AgentGavel report [flags] <run-id|path>

Render a GSI scorecard from a completed run's results directory.
Looks for scorecard.json, otherwise computes GSI from summary.json.

With --publish, write <dashboard>/data/<run-id>.json (tab=unratified,
sample=false) and update index.json. Opt-in publish is rejected until
v1.0 signatures (ADR 006 addendum).

Flags:
`)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return 2
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
		// --publish implies JSON output of the written path.
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
