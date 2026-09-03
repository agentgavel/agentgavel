package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteRunArtifact(t *testing.T) {
	root := t.TempDir()
	path, err := WriteRunArtifact(root, "run-1", RunArtifact{
		Fingerprint: map[string]string{"scenario-version": "SEC-v1"},
		Scenarios:   map[string]json.RawMessage{"SEC-001": json.RawMessage(`{"score":100}`)},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, "results", "run-1", "summary.json")
	if path != want {
		t.Fatalf("path %s want %s", path, want)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var art RunArtifact
	if err := json.Unmarshal(b, &art); err != nil {
		t.Fatal(err)
	}
	if art.RunID != "run-1" || art.Fingerprint["scenario-version"] != "SEC-v1" {
		t.Fatalf("%#v", art)
	}
}

func TestScenarioInterfaceMethods(t *testing.T) {
	var s Scenario = &noopScenario{}
	if s.ID() != "NOOP" {
		t.Fatal(s.ID())
	}
}
