package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/agentgavel/agentgavel/internal/report"
)

func runReport(args []string) int {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	root := fs.String("root", ".", "directory containing results/<run-id>/")
	asJSON := fs.Bool("json", false, "emit machine-readable scorecard JSON")
	fs.Usage = func() {
		_, _ = fmt.Fprintf(fs.Output(), `Usage: AgentGavel report [flags] <run-id|path>

Render a GSI scorecard from a completed run's results directory.
Looks for scorecard.json, otherwise computes GSI from summary.json.

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
