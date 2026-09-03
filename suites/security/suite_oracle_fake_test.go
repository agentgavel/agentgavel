package security

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/agentgavel/agentgavel/internal/engine"
)

// TestSuiteOracleFake is the T8.10 fixture E2E: Load the SEC-v2 catalog, score
// all-pass FakeAdapter golden observations through EvaluateSEC001..007 (oracle
// hard path), and write results/<runID>/summary.json via RunOracleFake (shared
// with the T9.2 CLI run path). SEC-008..010 are listed in Load but scored in
// SuiteOracleFakeV2 (T13.6).
func TestSuiteOracleFake(t *testing.T) {
	suite, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if suite.Version != "SEC-v2" {
		t.Fatalf("Version = %q, want SEC-v2", suite.Version)
	}

	wantIDs := []string{
		"SEC-001", "SEC-002", "SEC-003", "SEC-004",
		"SEC-005", "SEC-006", "SEC-007", "SEC-008",
		"SEC-009", "SEC-010",
	}
	gotIDs := make([]string, 0, len(suite.Scenarios))
	for _, s := range suite.Scenarios {
		gotIDs = append(gotIDs, s.ID)
	}
	for _, id := range wantIDs {
		if !slices.Contains(gotIDs, id) {
			t.Fatalf("missing scenario %s in Load() result %v", id, gotIDs)
		}
	}
	if len(gotIDs) != len(wantIDs) {
		t.Fatalf("got %d scenarios %v, want %d", len(gotIDs), gotIDs, len(wantIDs))
	}

	const runID = "suite-oracle-fake"
	root := t.TempDir()
	result, err := RunOracleFake(root, runID, OracleFakeOptions{
		RepoRoot: repoRoot(t),
		Seeds:    25,
		Model:    "oracle",
		// v0.1 all-pass path still scores SEC-001..007 only until T13.6.
		Scenarios: []string{
			"SEC-001", "SEC-002", "SEC-003", "SEC-004",
			"SEC-005", "SEC-006", "SEC-007",
		},
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
	if art.RunID != runID {
		t.Errorf("run_id = %q, want %q", art.RunID, runID)
	}
	if art.Fingerprint["model"] != "oracle" {
		t.Errorf("fingerprint model = %q, want oracle", art.Fingerprint["model"])
	}
	if art.Fingerprint["scenario.version"] != SuiteVersion {
		t.Errorf("fingerprint scenario.version = %q, want %s", art.Fingerprint["scenario.version"], SuiteVersion)
	}
	v1IDs := []string{
		"SEC-001", "SEC-002", "SEC-003", "SEC-004",
		"SEC-005", "SEC-006", "SEC-007",
	}
	if len(art.Scenarios) != 7 {
		t.Fatalf("scenarios keys = %d (%v), want 7", len(art.Scenarios), keysOf(art.Scenarios))
	}
	for _, id := range v1IDs {
		raw, ok := art.Scenarios[id]
		if !ok {
			t.Errorf("summary missing scenario key %s", id)
			continue
		}
		var body struct {
			Score float64 `json:"score"`
		}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Errorf("scenario %s JSON: %v", id, err)
			continue
		}
		if body.Score != 100 {
			t.Errorf("scenario %s score = %v, want 100", id, body.Score)
		}
	}
}

func keysOf(m map[string]json.RawMessage) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}
