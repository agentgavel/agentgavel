package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/agentgavel/agentgavel/internal/submit"
)

func runVerifyEntry(args []string) int {
	fs := flag.NewFlagSet("verify-entry", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	registry := fs.String("registry", "dashboard/keys/registry.json", "path to maintainer key registry (ADR 013)")
	fs.Usage = func() {
		_, _ = fmt.Fprintf(fs.Output(), `Usage: AgentGavel verify-entry [flags] <entry.json>

Verify an Opt-in dashboard entry signature against the key registry (ADR 013).
Exits 0 on success, 1 on verification failure, 2 on usage errors.

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
	reg, err := submit.LoadRegistry(*registry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "verify-entry: %v\n", err)
		return 1
	}
	raw, err := os.ReadFile(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "verify-entry: %v\n", err)
		return 1
	}
	entry, err := submit.EntryMapFromJSON(raw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "verify-entry: %v\n", err)
		return 1
	}
	if err := submit.Verify(reg, entry); err != nil {
		fmt.Fprintf(os.Stderr, "verify-entry: %v\n", err)
		return 1
	}
	fmt.Println("ok")
	return 0
}
