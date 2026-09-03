package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/agentgavel/agentgavel/internal/engine"
	"github.com/agentgavel/agentgavel/suites/security"
)

// rubberStampScenarios is the fixed SEC-002 + SEC-006 filter (ADR 011): a fast
// standalone check for approval forgery, never the full suite.
var rubberStampScenarios = []string{security.SEC002ID, security.SEC006ID}

func runRubberStamp(args []string) int {
	fs := flag.NewFlagSet("rubber-stamp", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	adapter := fs.String("adapter", "", "adapter command or path (required); may include args")
	seeds := fs.Int("seeds", 25, "number of deterministic seeds for the run fingerprint")
	outDir := fs.String("out", "", "results root directory (writes results/<run-id>/)")
	rootDir := fs.String("root", "", "alias for --out (directory containing results/)")
	runID := fs.String("run-id", "", "run id for results/<run-id>/ (default: generated)")
	fingerprintPath := fs.String("fingerprint", "", "reload seed-set (and other pins) from a prior fingerprint.json or summary.json")
	fs.Usage = func() {
		_, _ = fmt.Fprintf(fs.Output(), `Usage: AgentGavel rubber-stamp [flags]

Fast standalone check for SEC-002 (approval forgery) and SEC-006 (timeout
auto-approve) only — see docs/adr/011-rubber-stamp-verdict.md. It is not a
substitute for "run" and never prints a GSI or a grade.

stdout is the absolute path to summary.json, same as "run --ci". Exit codes:
0=pass, 1=fail, 2=catastrophic (catastrophic wins). When every selected
scenario is N/A (e.g. hitl=false), exits 1 and prints a
"rubber-stamp: not_applicable" line to stderr naming the missing capability.

Flags:
`)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *adapter == "" {
		fmt.Fprintln(os.Stderr, "rubber-stamp: --adapter is required")
		fs.Usage()
		return 2
	}
	resultsRoot := *outDir
	if resultsRoot == "" {
		resultsRoot = *rootDir
	}
	if resultsRoot == "" {
		resultsRoot = "."
	}
	if *seeds <= 0 {
		fmt.Fprintln(os.Stderr, "rubber-stamp: --seeds must be > 0")
		return 2
	}

	var pinned *engine.Fingerprint
	if strings.TrimSpace(*fingerprintPath) != "" {
		fp, err := engine.LoadFingerprintFile(*fingerprintPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "rubber-stamp: %v\n", err)
			return 2
		}
		pinned = &fp
	}

	id := *runID
	if id == "" {
		id = fmt.Sprintf("run-%d", time.Now().UTC().UnixNano())
	}

	repoRoot, err := security.FindRepoRoot(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "rubber-stamp: %v\n", err)
		return 1
	}

	cmd, cmdArgs := splitAdapterCommand(*adapter)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	caps, err := pingAdapter(ctx, cmd, cmdArgs, engine.ModeOracle)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rubber-stamp: adapter: %v\n", err)
		return 1
	}

	opts := security.OracleFakeOptions{
		RepoRoot:         repoRoot,
		Seeds:            *seeds,
		Scenarios:        rubberStampScenarios,
		AdapterVersion:   caps.AdapterVersion,
		Model:            "oracle",
		FrameworkVersion: version,
		Provenance:       caps.Provenance,
		Capabilities:     &caps,
	}
	if pinned != nil {
		opts.SeedSet = pinned.SeedSet
		if pinned.Model != "" {
			opts.Model = pinned.Model
		}
		if pinned.FrameworkVersion != "" {
			opts.FrameworkVersion = pinned.FrameworkVersion
		}
		if pinned.AdapterVersion != "" {
			opts.AdapterVersion = pinned.AdapterVersion
		}
	}

	result, err := security.RunOracleFake(resultsRoot, id, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rubber-stamp: %v\n", err)
		return 1
	}

	rel := result.Path
	if abs, err := filepath.Abs(result.Path); err == nil {
		rel = abs
	}
	// Machine-readable, same shape as "run --ci": summary.json path only.
	fmt.Println(rel)

	code := rubberStampExitCode(!result.AllPass, result.Catastrophic, len(result.NA), len(rubberStampScenarios))
	if code == 1 && len(result.NA) == len(rubberStampScenarios) {
		fmt.Fprintf(os.Stderr, "rubber-stamp: not_applicable (%s): no approval store to check\n", rubberStampNAReasons(result.NA))
	}
	return code
}

// rubberStampExitCode maps a completed rubber-stamp run to an exit code
// (ADR 011): pass->0, fail->1, catastrophic->2 (catastrophic wins), and when
// every selected scenario came back N/A it fails closed with 1 regardless of
// AllPass/Catastrophic — an all-N/A run has certified nothing.
func rubberStampExitCode(failed, catastrophic bool, naCount, total int) int {
	if total > 0 && naCount == total {
		return 1
	}
	return ciExitCode(failed, catastrophic)
}

// rubberStampNAReasons extracts the deduplicated "<reason>" half of each
// "<scenario id>: <reason>" entry in na, joined for the not_applicable line.
func rubberStampNAReasons(na []string) string {
	seen := make(map[string]bool, len(na))
	reasons := make([]string, 0, len(na))
	for _, entry := range na {
		reason := entry
		if _, r, ok := strings.Cut(entry, ": "); ok {
			reason = r
		}
		if !seen[reason] {
			seen[reason] = true
			reasons = append(reasons, reason)
		}
	}
	return strings.Join(reasons, ", ")
}
