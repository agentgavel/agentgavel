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
	"github.com/agentgavel/agentgavel/internal/protocol"
	"github.com/agentgavel/agentgavel/suites/security"
)

func runRun(args []string) int {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	adapter := fs.String("adapter", "", "adapter command or path (required); may include args")
	suite := fs.String("suite", "security", "suite to run (security)")
	seeds := fs.Int("seeds", 25, "number of deterministic seeds for the run fingerprint")
	mode := fs.String("mode", "oracle", "evaluation mode: oracle|model")
	outDir := fs.String("out", "", "results root directory (writes results/<run-id>/)")
	rootDir := fs.String("root", "", "alias for --out (directory containing results/)")
	scenarios := fs.String("scenarios", "", "optional comma-separated scenario IDs (default: all)")
	runID := fs.String("run-id", "", "run id for results/<run-id>/ (default: generated)")
	fingerprintPath := fs.String("fingerprint", "", "reload seed-set (and other pins) from a prior fingerprint.json or summary.json")
	fs.Usage = func() {
		_, _ = fmt.Fprintf(fs.Output(), `Usage: AgentGavel run [flags]

Run a suite against an adapter and write results/<run-id>/summary.json.

Flags:
`)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *adapter == "" {
		fmt.Fprintln(os.Stderr, "run: --adapter is required")
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

	switch *mode {
	case engine.ModeOracle, engine.ModeModel:
	default:
		fmt.Fprintf(os.Stderr, "run: unknown --mode %q (want oracle|model)\n", *mode)
		return 2
	}
	if *suite != "security" {
		fmt.Fprintf(os.Stderr, "run: unsupported --suite %q (want security)\n", *suite)
		return 2
	}
	if *mode == engine.ModeModel {
		fmt.Fprintln(os.Stderr, "run: --mode model is not implemented in v0.1 (use --mode oracle)")
		return 2
	}
	if *seeds <= 0 {
		fmt.Fprintln(os.Stderr, "run: --seeds must be > 0")
		return 2
	}

	var pinned *engine.Fingerprint
	if strings.TrimSpace(*fingerprintPath) != "" {
		fp, err := engine.LoadFingerprintFile(*fingerprintPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "run: %v\n", err)
			return 2
		}
		pinned = &fp
	}

	id := *runID
	if id == "" {
		id = fmt.Sprintf("run-%d", time.Now().UTC().UnixNano())
	}

	var scenarioList []string
	if strings.TrimSpace(*scenarios) != "" {
		for _, s := range strings.Split(*scenarios, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				scenarioList = append(scenarioList, s)
			}
		}
	}

	repoRoot, err := security.FindRepoRoot(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "run: %v\n", err)
		return 1
	}

	cmd, cmdArgs := splitAdapterCommand(*adapter)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	caps, err := pingAdapter(ctx, cmd, cmdArgs, *mode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "run: adapter: %v\n", err)
		return 1
	}

	opts := security.OracleFakeOptions{
		RepoRoot:         repoRoot,
		Seeds:            *seeds,
		Scenarios:        scenarioList,
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
		fmt.Fprintf(os.Stderr, "run: %v\n", err)
		return 1
	}

	rel := result.Path
	if abs, err := filepath.Abs(result.Path); err == nil {
		rel = abs
	}
	fmt.Printf("wrote %s\n", rel)
	if !result.AllPass {
		for _, f := range result.Failures {
			fmt.Fprintf(os.Stderr, "run: fail %s\n", f)
		}
		return 1
	}
	return 0
}

// splitAdapterCommand splits "--adapter" into exec path + args.
// A bare path with no spaces is treated as the executable.
func splitAdapterCommand(s string) (string, []string) {
	s = strings.TrimSpace(s)
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return "", nil
	}
	return fields[0], fields[1:]
}

// pingAdapter launches the adapter, Handshake + StartSession + StopSession.
// Returns the Handshake CapabilityReport (provenance, version, N/A drivers).
func pingAdapter(ctx context.Context, command string, args []string, mode string) (protocol.CapabilityReport, error) {
	sess, err := engine.Launch(ctx, engine.LaunchConfig{
		Command: command,
		Args:    args,
		Timeout: 10 * time.Second,
	})
	if err != nil {
		return protocol.CapabilityReport{}, err
	}
	defer func() { _ = sess.Close() }()

	rep, err := sess.Handshake(ctx, protocol.HandshakeRequest{EngineProtocolVersion: "1.0"})
	if err != nil {
		return protocol.CapabilityReport{}, fmt.Errorf("handshake: %w", err)
	}

	sid, _, err := sess.StartSessionForMode(ctx, mode, engine.ModeEndpoints{
		// v0.1 oracle all-pass path scores golden observations locally; the
		// adapter still receives a SessionConfig with ModelBaseURL set.
		OracleURL: "http://127.0.0.1:9",
	})
	if err != nil {
		return protocol.CapabilityReport{}, fmt.Errorf("start session: %w", err)
	}
	if err := sess.StopSession(ctx, sid); err != nil {
		return protocol.CapabilityReport{}, fmt.Errorf("stop session: %w", err)
	}
	return rep, nil
}
