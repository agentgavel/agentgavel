// Command AgentGavel is the AgentGavel benchmarking harness CLI.
package main

import (
	"flag"
	"fmt"
	"os"
)

// version is overridden at link time via -ldflags "-X main.version=...".
var version = "0.0.0-dev"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		printHelp(os.Stderr)
		return 0
	}
	switch args[0] {
	case "version":
		fmt.Println(version)
		return 0
	case "oracle":
		return runOracle(args[1:])
	case "report":
		return runReport(args[1:])
	case "run":
		return runRun(args[1:])
	case "help", "-h", "--help":
		printHelp(os.Stdout)
		return 0
	default:
		fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
		fs.SetOutput(os.Stderr)
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", args[0])
		printHelp(os.Stderr)
		_ = fs
		return 2
	}
}

func printHelp(w *os.File) {
	_, _ = fmt.Fprintf(w, `AgentGavel is an adversarial governance and security benchmark harness.

Usage:
  AgentGavel <command>

Commands:
  version  Print the AgentGavel version
  oracle   Run the Compliance Oracle HTTP server (--listen host:port)
  report   Print GSI scorecard text (or --json) from results/<run-id>
  run      Run a suite against an adapter; write results/<run-id>/summary.json (--ci for CI exit codes)
  help     Show this help

See docs/RFC-0001.md for the specification.
`)
}
