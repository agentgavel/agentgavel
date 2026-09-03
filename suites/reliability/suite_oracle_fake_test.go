package reliability

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/agentgavel/agentgavel/internal/engine"
)

// TestSuiteOracleFakeREL is the T14.6 fixture E2E: score all-pass FakeAdapter
// golden observations through EvaluateREL001..003 and write
// results/<runID>/summary.json via RunOracleFakeREL — the REL-v0 counterpart
// of suites/security's TestSuiteOracleFake.
func TestSuiteOracleFakeREL(t *testing.T) {
	suite, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if suite.Version != SuiteVersion {
		t.Fatalf("Version = %q, want %s", suite.Version, SuiteVersion)
	}

	const runID = "suite-oracle-fake-rel"
	root := t.TempDir()
	result, err := RunOracleFakeREL(root, runID, OracleFakeOptions{
		Seeds: 25,
		Model: "oracle",
	})
	if err != nil {
		t.Fatalf("RunOracleFakeREL: %v", err)
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
	if art.Fingerprint["model"] != "oracle" {
		t.Errorf("fingerprint model = %q, want oracle", art.Fingerprint["model"])
	}
	if art.Fingerprint["scenario.version"] != SuiteVersion {
		t.Errorf("fingerprint scenario.version = %q, want %s", art.Fingerprint["scenario.version"], SuiteVersion)
	}

	wantIDs := []string{REL001ID, REL002ID, REL003ID}
	if len(art.Scenarios) != len(wantIDs) {
		t.Fatalf("scenarios = %d (%v), want %d", len(art.Scenarios), keysOfREL(art.Scenarios), len(wantIDs))
	}
	for _, id := range wantIDs {
		t.Run(id, func(t *testing.T) {
			raw, ok := art.Scenarios[id]
			if !ok {
				t.Fatalf("summary missing scenario row %s", id)
			}
			var row struct {
				Score float64 `json:"score"`
			}
			if err := json.Unmarshal(raw, &row); err != nil {
				t.Fatalf("scenario %s JSON: %v", id, err)
			}
			if row.Score != 100 {
				t.Fatalf("scenario %s score = %v, want 100", id, row.Score)
			}
			t.Logf("%s row: %s", id, raw)
		})
	}
}

func keysOfREL(m map[string]json.RawMessage) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
