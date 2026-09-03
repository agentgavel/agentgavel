package security

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/agentgavel/gavel/internal/engine"
	"github.com/agentgavel/gavel/internal/mcpfuzz"
)

// TestSuiteOracleFake is the T8.10 fixture E2E: Load the SEC-v1 catalog, score
// all-pass FakeAdapter golden observations through EvaluateSEC001..007 (oracle
// hard path), and write results/<runID>/summary.json via the engine artifact
// writer. Full adapter event streams are out of scope (T9.2).
func TestSuiteOracleFake(t *testing.T) {
	suite, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if suite.Version != "SEC-v1" {
		t.Fatalf("Version = %q, want SEC-v1", suite.Version)
	}

	wantIDs := []string{
		"SEC-001", "SEC-002", "SEC-003", "SEC-004",
		"SEC-005", "SEC-006", "SEC-007",
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

	creds, err := LoadSEC004Credentials(repoRoot(t))
	if err != nil {
		t.Fatalf("LoadSEC004Credentials: %v", err)
	}
	if len(creds) == 0 {
		t.Fatal("SEC-004 credentials fixture empty")
	}

	type scored struct {
		id           string
		score        float64
		catastrophic bool
	}
	scores := []scored{
		{
			id: SEC001ID,
			score: EvaluateSEC001(SEC001Observation{
				HardEligible: true,
				Events:       fakeAdapterHardRefuse(),
			}).Score,
		},
		func() scored {
			r := EvaluateSEC002(SEC002Observation{
				HardEligible: true,
				Events:       fakeAdapterSEC002HardStore(),
			})
			return scored{id: SEC002ID, score: r.Score, catastrophic: r.Catastrophic}
		}(),
		{
			id: SEC003ID,
			score: EvaluateSEC003(SEC003Observation{
				Events: fakeAdapterSEC003OutOfScopeRefused(),
			}).Score,
		},
		func() scored {
			r := EvaluateSEC004(SEC004Observation{
				Credentials: creds,
				Events:      fakeAdapterSEC004Clean(creds[0]),
			})
			return scored{id: SEC004ID, score: r.Score, catastrophic: r.Catastrophic}
		}(),
		func() scored {
			r := EvaluateSEC005(SEC005Observation{
				StructuralSeparation: true,
				Events:               fakeAdapterSEC005Structural(),
			})
			return scored{id: SEC005ID, score: r.Score, catastrophic: r.Catastrophic}
		}(),
		func() scored {
			r := EvaluateSEC006(SEC006Observation{
				Events: fakeAdapterSEC006CleanHold(),
			})
			return scored{id: SEC006ID, score: r.Score, catastrophic: r.Catastrophic}
		}(),
		func() scored {
			r := EvaluateSEC007(SEC007Observation{
				Runs: []SEC007RunObservation{
					{Mode: mcpfuzz.ModeToxicOutput},
					{Mode: mcpfuzz.ModeSchemaFlooder},
					{Mode: mcpfuzz.ModeEarlyDisconnect, Interrupted: true, ReceiptPresent: true},
					{Mode: mcpfuzz.ModeToolRenamer},
					{Mode: mcpfuzz.ModeSlowloris},
					{Mode: mcpfuzz.ModeMasquerade},
				},
			})
			return scored{id: SEC007ID, score: r.Score, catastrophic: r.Catastrophic}
		}(),
	}

	scenarios := make(map[string]json.RawMessage, len(scores))
	for _, s := range scores {
		if s.score != 100 {
			t.Errorf("%s score = %v, want 100 (all-pass oracle golden)", s.id, s.score)
		}
		if s.catastrophic {
			t.Errorf("%s unexpected Catastrophic on all-pass oracle golden", s.id)
		}
		raw, err := json.Marshal(map[string]any{
			"score":        s.score,
			"catastrophic": s.catastrophic,
		})
		if err != nil {
			t.Fatalf("marshal %s: %v", s.id, err)
		}
		scenarios[s.id] = raw
		t.Logf("%s → score=%.0f catastrophic=%v", s.id, s.score, s.catastrophic)
	}
	if t.Failed() {
		t.Fatal("oracle golden evaluations failed; skipping artifact write")
	}

	const runID = "suite-oracle-fake"
	root := t.TempDir()
	fp := engine.Fingerprint{
		Model:           "oracle",
		ScenarioVersion: suite.Version,
	}
	path, err := engine.WriteRunArtifact(root, runID, engine.RunArtifact{
		Fingerprint: fp.Fields(),
		Scenarios:   scenarios,
	})
	if err != nil {
		t.Fatalf("WriteRunArtifact: %v", err)
	}

	wantPath := filepath.Join(root, "results", runID, "summary.json")
	if path != wantPath {
		t.Fatalf("artifact path = %s, want %s", path, wantPath)
	}
	b, err := os.ReadFile(path)
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
	if len(art.Scenarios) != 7 {
		t.Fatalf("scenarios keys = %d (%v), want 7", len(art.Scenarios), keysOf(art.Scenarios))
	}
	for _, id := range wantIDs {
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
