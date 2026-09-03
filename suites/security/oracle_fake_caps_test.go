package security

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/agentgavel/gavel/internal/engine"
	"github.com/agentgavel/gavel/internal/protocol"
)

func TestRunOracleFakeRecordsProvenanceAndNAFromCapabilities(t *testing.T) {
	root := t.TempDir()
	const runID = "oracle-caps-na"
	caps := protocol.CapabilityReport{
		AdapterProtocolVersion: "1.0",
		AdapterName:            "sire",
		AdapterVersion:         "0.0.1",
		Provenance:             "unofficial",
		HITL:                   true,
		Ledger:                 false,
		Observability:          false,
		ContextMode:            "none",
		FrameworkVersion:       "unknown",
	}
	result, err := RunOracleFake(root, runID, OracleFakeOptions{
		RepoRoot:     repoRoot(t),
		Seeds:        2,
		Scenarios:    []string{"SEC-002", "SEC-004"},
		Capabilities: &caps,
	})
	if err != nil {
		t.Fatalf("RunOracleFake: %v", err)
	}
	if !result.AllPass {
		t.Fatalf("expected all-pass (SEC-002 scored, SEC-004 N/A), failures: %v", result.Failures)
	}

	b, err := os.ReadFile(result.Path)
	if err != nil {
		t.Fatal(err)
	}
	var art engine.RunArtifact
	if err := json.Unmarshal(b, &art); err != nil {
		t.Fatal(err)
	}
	if art.Provenance != "unofficial" {
		t.Fatalf("provenance = %q", art.Provenance)
	}
	if !art.ObservabilityPenalty {
		t.Fatal("expected observability_penalty from capabilities")
	}
	if art.Fingerprint["provenance"] != "unofficial" {
		t.Fatalf("fingerprint provenance = %q", art.Fingerprint["provenance"])
	}

	var sec002 struct {
		Score float64 `json:"score"`
		NA    bool    `json:"na"`
	}
	if err := json.Unmarshal(art.Scenarios["SEC-002"], &sec002); err != nil {
		t.Fatal(err)
	}
	if sec002.NA || sec002.Score != 100 {
		t.Fatalf("SEC-002 = %#v, want score 100", sec002)
	}

	var sec004 struct {
		NA bool `json:"na"`
	}
	if err := json.Unmarshal(art.Scenarios["SEC-004"], &sec004); err != nil {
		t.Fatal(err)
	}
	if !sec004.NA {
		t.Fatalf("SEC-004 want na:true, got %s", art.Scenarios["SEC-004"])
	}

	fpPath := filepath.Join(root, "results", runID, "fingerprint.json")
	if _, err := os.Stat(fpPath); err != nil {
		t.Fatal(err)
	}
}
