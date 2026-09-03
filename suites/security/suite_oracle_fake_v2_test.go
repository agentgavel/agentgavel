package security

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/agentgavel/agentgavel/internal/engine"
)

// TestSuiteOracleFakeV2 is the T13.6 E2E: score SEC-008..010 FakeAdapter
// all-pass goldens through RunOracleFake and assert summary rows.
func TestSuiteOracleFakeV2(t *testing.T) {
	const runID = "suite-oracle-fake-v2"
	root := t.TempDir()
	v2IDs := []string{SEC008ID, SEC009ID, SEC010ID}
	result, err := RunOracleFake(root, runID, OracleFakeOptions{
		RepoRoot:  repoRoot(t),
		Seeds:     25,
		Model:     "oracle",
		Scenarios: v2IDs,
	})
	if err != nil {
		t.Fatalf("RunOracleFake: %v", err)
	}
	if !result.AllPass {
		t.Fatalf("expected all-pass, failures: %v", result.Failures)
	}

	wantPath := filepath.Join(root, "results", runID, "summary.json")
	if result.Path != wantPath {
		t.Fatalf("artifact path = %s, want %s", result.Path, wantPath)
	}
	b, err := os.ReadFile(result.Path)
	if err != nil {
		t.Fatalf("read summary.json: %v", err)
	}
	var art engine.RunArtifact
	if err := json.Unmarshal(b, &art); err != nil {
		t.Fatalf("unmarshal summary.json: %v", err)
	}
	if art.Fingerprint["scenario.version"] != SuiteVersion {
		t.Errorf("fingerprint scenario.version = %q, want %s", art.Fingerprint["scenario.version"], SuiteVersion)
	}
	if len(art.Scenarios) != 3 {
		t.Fatalf("scenarios keys = %d (%v), want 3", len(art.Scenarios), keysOf(art.Scenarios))
	}
	for _, id := range v2IDs {
		raw, ok := art.Scenarios[id]
		if !ok {
			t.Errorf("summary missing scenario key %s", id)
			continue
		}
		var body struct {
			Score float64 `json:"score"`
			NA    bool    `json:"na"`
		}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Errorf("scenario %s JSON: %v", id, err)
			continue
		}
		if body.NA {
			t.Errorf("scenario %s unexpected na:true", id)
			continue
		}
		if body.Score != 100 {
			t.Errorf("scenario %s score = %v, want 100", id, body.Score)
		}
	}
}
