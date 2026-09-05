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

// TestRubberStampFakeAdapterWritesOnlySEC002SEC006 covers T14.7's core
// contract (ADR 011): rubber-stamp is a thin SEC-002 + SEC-006 filter, not
// the full suite, and it never prints a GSI or a grade.
func TestRubberStampFakeAdapterWritesOnlySEC002SEC006(t *testing.T) {
	bin := buildAgentGavel(t)
	fake := buildFakeAdapterBin(t)
	root := t.TempDir()
	const runID = "cli-rubber-stamp"

	cmd := exec.Command(bin, "rubber-stamp",
		"--adapter", fake,
		"--seeds", "25",
		"--out", root,
		"--run-id", runID,
	)
	// Working directory must be the module root so FindRepoRoot finds fixtures/.
	cmd.Dir = filepath.Join("..", "..")
	cmd.Env = append(os.Environ(), "GOWORK=off")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("rubber-stamp: %v\n%s", err, out)
	}
	text := strings.TrimSpace(string(out))
	if strings.Contains(text, "GSI") || strings.Contains(text, "grade") {
		t.Fatalf("rubber-stamp must never print a GSI or a grade:\n%s", text)
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
	if len(art.Scenarios) != 2 {
		t.Fatalf("scenarios = %d, want 2 (SEC-002, SEC-006 only)", len(art.Scenarios))
	}
	for _, id := range []string{"SEC-002", "SEC-006"} {
		if _, ok := art.Scenarios[id]; !ok {
			t.Errorf("missing %s in summary.json", id)
		}
	}
}

// TestRubberStampExitMap covers rubberStampExitCode (ADR 011): pass->0,
// fail->1, catastrophic->2 (catastrophic wins), and the all-N/A fail-closed
// case, without spinning up a live adapter.
func TestRubberStampExitMap(t *testing.T) {
	cases := []struct {
		name         string
		failed       bool
		catastrophic bool
		naCount      int
		total        int
		want         int
	}{
		{name: "pass", failed: false, catastrophic: false, naCount: 0, total: 2, want: 0},
		{name: "fail", failed: true, catastrophic: false, naCount: 0, total: 2, want: 1},
		{name: "catastrophic", failed: false, catastrophic: true, naCount: 0, total: 2, want: 2},
		{name: "fail_and_catastrophic_wins", failed: true, catastrophic: true, naCount: 0, total: 2, want: 2},
		{name: "one_na_still_scores_the_other", failed: false, catastrophic: false, naCount: 1, total: 2, want: 0},
		{name: "all_na_fails_closed", failed: false, catastrophic: false, naCount: 2, total: 2, want: 1},
		{name: "all_na_wins_over_pass", failed: false, catastrophic: false, naCount: 2, total: 2, want: 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rubberStampExitCode(tc.failed, tc.catastrophic, tc.naCount, tc.total)
			if got != tc.want {
				t.Fatalf("rubberStampExitCode(%v, %v, %d, %d) = %d, want %d",
					tc.failed, tc.catastrophic, tc.naCount, tc.total, got, tc.want)
			}
		})
	}
}

// TestRubberStampNAReasons covers the stderr message helper used when every
// selected scenario is N/A (ADR 011's "not_applicable (hitl=false)" line).
func TestRubberStampNAReasons(t *testing.T) {
	got := rubberStampNAReasons([]string{"SEC-002: hitl=false", "SEC-006: hitl=false"})
	if got != "hitl=false" {
		t.Fatalf("rubberStampNAReasons = %q, want deduplicated %q", got, "hitl=false")
	}
}
