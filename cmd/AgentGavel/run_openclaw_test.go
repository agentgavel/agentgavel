package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/agentgavel/agentgavel/internal/engine"
)

// TestRubberStampOpenClawHitlFalseNotApplicable is T15.20: OpenClaw reports
// hitl=false, so rubber-stamp's SEC-002+SEC-006 filter is all-N/A. ADR 011
// requires exit 1 and a "rubber-stamp: not_applicable" stderr line (no live
// Gateway HITL invented).
func TestRubberStampOpenClawHitlFalseNotApplicable(t *testing.T) {
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("python3 not available")
	}

	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	openclawSrc := filepath.Join(repoRoot, "adapters", "openclaw", "src")
	sdkSrc := filepath.Join(repoRoot, "sdk", "python", "src")
	if _, err := os.Stat(filepath.Join(openclawSrc, "adapters", "openclaw", "__main__.py")); err != nil {
		t.Fatalf("openclaw adapter missing: %v", err)
	}

	bin := buildAgentGavel(t)
	root := t.TempDir()
	const runID = "cli-rubber-stamp-openclaw"

	pyPath := openclawSrc + string(os.PathListSeparator) + sdkSrc
	cmd := exec.Command(bin, "rubber-stamp",
		"--adapter", "python3 -m adapters.openclaw",
		"--seeds", "3",
		"--out", root,
		"--run-id", runID,
	)
	cmd.Dir = repoRoot
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = append(os.Environ(), "GOWORK=off", "PYTHONPATH="+pyPath)
	err = cmd.Run()
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("rubber-stamp: want exit 1 (*exec.ExitError), got %v\nstdout:\n%s\nstderr:\n%s",
			err, stdout.String(), stderr.String())
	}
	if exitErr.ExitCode() != 1 {
		t.Fatalf("exit code = %d, want 1\nstdout:\n%s\nstderr:\n%s",
			exitErr.ExitCode(), stdout.String(), stderr.String())
	}
	errText := stderr.String()
	if !strings.Contains(errText, "rubber-stamp: not_applicable") {
		t.Fatalf("stderr missing rubber-stamp: not_applicable:\n%s", errText)
	}
	if !strings.Contains(errText, "hitl=false") {
		t.Fatalf("stderr missing hitl=false reason:\n%s", errText)
	}
	if strings.Contains(stdout.String(), "GSI") || strings.Contains(errText, "GSI") {
		t.Fatalf("rubber-stamp must never print a GSI:\nstdout=%s\nstderr=%s", stdout.String(), errText)
	}

	summaryPath := filepath.Join(root, "results", runID, "summary.json")
	if !strings.Contains(stdout.String(), "summary.json") {
		// Path may still be printed before exit 1; tolerate absolute path only on disk.
		if _, statErr := os.Stat(summaryPath); statErr != nil {
			t.Fatalf("expected summary.json on stdout or disk; stdout=%q err=%v", stdout.String(), statErr)
		}
	}
	b, err := os.ReadFile(summaryPath)
	if err != nil {
		t.Fatalf("read summary: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), errText)
	}
	var art engine.RunArtifact
	if err := json.Unmarshal(b, &art); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if art.Provenance != "unofficial" {
		t.Fatalf("summary provenance = %q, want unofficial", art.Provenance)
	}
	for _, id := range []string{"SEC-002", "SEC-006"} {
		raw, ok := art.Scenarios[id]
		if !ok {
			t.Fatalf("missing %s in scenarios keys %v", id, keysOfRaw(art.Scenarios))
		}
		var row struct {
			NA bool `json:"na"`
		}
		if err := json.Unmarshal(raw, &row); err != nil {
			t.Fatalf("%s row: %v", id, err)
		}
		if !row.NA {
			t.Fatalf("%s want na:true, got %s", id, raw)
		}
	}
}

// TestRunOpenClawAdapterOracleSEC002NA is the optional T15.20 companion:
// AgentGavel run --adapter openclaw --scenarios SEC-002 scores honest na:true
// (CapabilityReport.hitl=false), never inventing live Gateway HITL.
func TestRunOpenClawAdapterOracleSEC002NA(t *testing.T) {
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("python3 not available")
	}

	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	openclawSrc := filepath.Join(repoRoot, "adapters", "openclaw", "src")
	sdkSrc := filepath.Join(repoRoot, "sdk", "python", "src")
	if _, err := os.Stat(filepath.Join(openclawSrc, "adapters", "openclaw", "__main__.py")); err != nil {
		t.Fatalf("openclaw adapter missing: %v", err)
	}

	bin := buildAgentGavel(t)
	root := t.TempDir()
	const runID = "cli-oracle-openclaw-sec002"

	pyPath := openclawSrc + string(os.PathListSeparator) + sdkSrc
	cmd := exec.Command(bin, "run",
		"--adapter", "python3 -m adapters.openclaw",
		"--suite", "security",
		"--mode", "oracle",
		"--scenarios", "SEC-002",
		"--seeds", "3",
		"--out", root,
		"--run-id", runID,
	)
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "GOWORK=off", "PYTHONPATH="+pyPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "summary.json") {
		t.Fatalf("expected summary.json path in stdout:\n%s", out)
	}

	summaryPath := filepath.Join(root, "results", runID, "summary.json")
	b, err := os.ReadFile(summaryPath)
	if err != nil {
		t.Fatalf("read summary: %v", err)
	}
	var art engine.RunArtifact
	if err := json.Unmarshal(b, &art); err != nil {
		t.Fatalf("unmarshal summary: %v", err)
	}
	if art.Provenance != "unofficial" {
		t.Fatalf("summary provenance = %q, want unofficial", art.Provenance)
	}
	raw, ok := art.Scenarios["SEC-002"]
	if !ok {
		t.Fatalf("missing SEC-002 in scenarios keys %v", keysOfRaw(art.Scenarios))
	}
	var row struct {
		NA bool `json:"na"`
	}
	if err := json.Unmarshal(raw, &row); err != nil {
		t.Fatalf("SEC-002 row: %v", err)
	}
	if !row.NA {
		t.Fatalf("SEC-002 want na:true (hitl=false), got %s", raw)
	}
}
