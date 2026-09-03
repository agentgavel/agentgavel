package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/agentgavel/agentgavel/internal/engine"
)

// TestRunReliabilityFakeAdapterWritesRELRows is T14.17: AgentGavel run
// --suite reliability --adapter <fakeadapter> --mode oracle --ci writes
// summary.json containing REL-001, REL-002, REL-003 rows and exits 0. The
// FakeAdapter reports hitl=true/ledger=true, so all three score (no N/A).
func TestRunReliabilityFakeAdapterWritesRELRows(t *testing.T) {
	bin := buildAgentGavel(t)
	fake := buildFakeAdapterBin(t)
	root := t.TempDir()
	const runID = "cli-oracle-fake-rel"

	cmd := exec.Command(bin, "run",
		"--ci",
		"--adapter", fake,
		"--suite", "reliability",
		"--seeds", "25",
		"--mode", "oracle",
		"--out", root,
		"--run-id", runID,
	)
	cmd.Dir = filepath.Join("..", "..")
	cmd.Env = append(os.Environ(), "GOWORK=off")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}

	summary := filepath.Join(root, "results", runID, "summary.json")
	b, err := os.ReadFile(summary)
	if err != nil {
		t.Fatalf("read summary: %v\nstdout:\n%s", err, out)
	}
	var art engine.RunArtifact
	if err := json.Unmarshal(b, &art); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if art.RunID != runID {
		t.Errorf("run_id = %q, want %q", art.RunID, runID)
	}
	if art.Provenance != "unofficial" {
		t.Errorf("summary provenance = %q, want unofficial", art.Provenance)
	}
	if len(art.Scenarios) != 3 {
		t.Fatalf("scenarios = %d, want 3 (REL-v0)", len(art.Scenarios))
	}
	for _, id := range []string{"REL-001", "REL-002", "REL-003"} {
		raw, ok := art.Scenarios[id]
		if !ok {
			t.Errorf("missing %s", id)
			continue
		}
		var row struct {
			Score float64 `json:"score"`
			NA    bool    `json:"na"`
		}
		if err := json.Unmarshal(raw, &row); err != nil {
			t.Errorf("%s: %v", id, err)
			continue
		}
		if row.NA {
			t.Errorf("%s scored N/A, want a score (FakeAdapter reports hitl=true, ledger=true)", id)
			continue
		}
		if row.Score != 100 {
			t.Errorf("%s score = %v, want 100", id, row.Score)
		}
	}
}

// TestRunUnsupportedSuiteStillExits2 covers --suite values outside
// security|reliability: usage error, exit 2.
func TestRunUnsupportedSuiteStillExits2(t *testing.T) {
	bin := buildAgentGavel(t)
	cmd := exec.Command(bin, "run", "--adapter", "noop", "--suite", "bogus", "--mode", "oracle")
	out, err := cmd.CombinedOutput()
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected *exec.ExitError, got %v (out=%s)", err, out)
	}
	if exitErr.ExitCode() != 2 {
		t.Fatalf("exit code = %d, want 2\n%s", exitErr.ExitCode(), out)
	}
	if !strings.Contains(string(out), "--suite") {
		t.Fatalf("unexpected stderr: %s", out)
	}
}
