package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/agentgavel/gavel/internal/engine"
	"github.com/agentgavel/gavel/internal/report"
)

// TestRunSireAdapterOracleSEC002 is T10.5: AgentGavel run --adapter sire
// --mode oracle on SEC-002 records Handshake provenance=unofficial and a
// SEC-002 scorecard row (numeric score or na:true).
func TestRunSireAdapterOracleSEC002(t *testing.T) {
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("python3 not available")
	}

	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	sireSrc := filepath.Join(repoRoot, "adapters", "sire", "src")
	sdkSrc := filepath.Join(repoRoot, "sdk", "python", "src")
	if _, err := os.Stat(filepath.Join(sireSrc, "adapters", "sire", "__main__.py")); err != nil {
		t.Fatalf("sire adapter missing: %v", err)
	}

	bin := buildAgentGavel(t)
	root := t.TempDir()
	const runID = "cli-oracle-sire-sec002"

	pyPath := sireSrc + string(os.PathListSeparator) + sdkSrc
	cmd := exec.Command(bin, "run",
		"--adapter", "python3 -m adapters.sire",
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
	if art.Fingerprint["provenance"] != "unofficial" {
		t.Fatalf("fingerprint provenance = %q, want unofficial", art.Fingerprint["provenance"])
	}
	if art.Fingerprint["adapter.version"] == "" {
		t.Fatal("missing fingerprint adapter.version")
	}
	raw, ok := art.Scenarios["SEC-002"]
	if !ok {
		t.Fatalf("missing SEC-002 in scenarios keys %v", keysOfRaw(art.Scenarios))
	}
	var row struct {
		Score float64 `json:"score"`
		NA    bool    `json:"na"`
	}
	if err := json.Unmarshal(raw, &row); err != nil {
		t.Fatalf("SEC-002 row: %v", err)
	}
	if !row.NA && row.Score == 0 && !strings.Contains(string(raw), "score") {
		t.Fatalf("SEC-002 row must be numeric score or na:true, got %s", raw)
	}
	if row.NA {
		// Honest N/A is acceptable when HITL/caps cannot support the scenario.
		t.Logf("SEC-002 scored N/A (honest capability gap): %s", raw)
	} else if row.Score < 0 || row.Score > 100 {
		t.Fatalf("SEC-002 score = %v, want 0..100", row.Score)
	}

	doc, err := report.Load(filepath.Join(root, "results", runID))
	if err != nil {
		t.Fatalf("report.Load: %v", err)
	}
	if doc.Provenance != "unofficial" {
		t.Fatalf("scorecard provenance = %q, want unofficial", doc.Provenance)
	}

	reportCmd := exec.Command(bin, "report", "--json", "--root", root, runID)
	reportOut, err := reportCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("report --json: %v\n%s", err, reportOut)
	}
	var scorecard report.Document
	if err := json.Unmarshal(reportOut, &scorecard); err != nil {
		t.Fatalf("parse report JSON: %v\n%s", err, reportOut)
	}
	if scorecard.Provenance != "unofficial" {
		t.Fatalf("report JSON provenance = %q, want unofficial", scorecard.Provenance)
	}
}

func keysOfRaw(m map[string]json.RawMessage) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
