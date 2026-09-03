package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// RunArtifact is the on-disk summary for a harness run.
type RunArtifact struct {
	RunID                string                     `json:"run_id"`
	CreatedAt            time.Time                  `json:"created_at"`
	Provenance           string                     `json:"provenance,omitempty"`
	ObservabilityPenalty bool                       `json:"observability_penalty,omitempty"`
	Fingerprint          map[string]string          `json:"fingerprint"`
	Scenarios            map[string]json.RawMessage `json:"scenarios"`
}

// WriteRunArtifact writes results/<run-id>/summary.json.
func WriteRunArtifact(root, runID string, art RunArtifact) (string, error) {
	dir := filepath.Join(root, "results", runID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	art.RunID = runID
	if art.CreatedAt.IsZero() {
		art.CreatedAt = time.Now().UTC()
	}
	if art.Fingerprint == nil {
		art.Fingerprint = map[string]string{}
	}
	if art.Scenarios == nil {
		art.Scenarios = map[string]json.RawMessage{}
	}
	path := filepath.Join(dir, "summary.json")
	b, err := json.MarshalIndent(art, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return "", err
	}
	if _, err := WriteFingerprintFile(dir, art.Fingerprint); err != nil {
		return "", err
	}
	return path, nil
}
