package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/agentgavel/gavel/internal/engine"
)

func buildFakeAdapterBin(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "fakeadapter")
	// Build from module root so ./internal/engine/testdata/fakeadapter resolves.
	cmd := exec.Command("go", "build", "-o", bin, "./internal/engine/testdata/fakeadapter")
	cmd.Dir = filepath.Join("..", "..")
	cmd.Env = append(os.Environ(), "GOWORK=off")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build fakeadapter: %v\n%s", err, out)
	}
	return bin
}

func TestRunOracleFakeAllPassWritesSummary(t *testing.T) {
	bin := buildAgentGavel(t)
	fake := buildFakeAdapterBin(t)
	root := t.TempDir()
	const runID = "cli-oracle-fake"

	cmd := exec.Command(bin, "run",
		"--adapter", fake,
		"--suite", "security",
		"--seeds", "25",
		"--mode", "oracle",
		"--out", root,
		"--run-id", runID,
	)
	// Working directory must be the module root so FindRepoRoot finds fixtures/.
	cmd.Dir = filepath.Join("..", "..")
	cmd.Env = append(os.Environ(), "GOWORK=off")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	text := string(out)
	if !strings.Contains(text, "summary.json") {
		t.Fatalf("expected summary.json path in stdout:\n%s", text)
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
	if art.RunID != runID {
		t.Errorf("run_id = %q, want %q", art.RunID, runID)
	}
	if art.Fingerprint["model"] != "oracle" {
		t.Errorf("fingerprint model = %q, want oracle", art.Fingerprint["model"])
	}
	if art.Fingerprint["adapter.version"] != "0.0.1" {
		t.Errorf("fingerprint adapter.version = %q, want 0.0.1", art.Fingerprint["adapter.version"])
	}
	wantSeeds := "0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24"
	if art.Fingerprint["seed.set"] != wantSeeds {
		t.Errorf("seed.set = %q, want 0..24", art.Fingerprint["seed.set"])
	}
	fpPath := filepath.Join(root, "results", runID, "fingerprint.json")
	if _, err := os.Stat(fpPath); err != nil {
		t.Fatalf("fingerprint.json: %v", err)
	}
	if len(art.Scenarios) != 7 {
		t.Fatalf("scenarios = %d, want 7", len(art.Scenarios))
	}
	for _, id := range []string{"SEC-001", "SEC-002", "SEC-003", "SEC-004", "SEC-005", "SEC-006", "SEC-007"} {
		raw, ok := art.Scenarios[id]
		if !ok {
			t.Errorf("missing %s", id)
			continue
		}
		var body struct {
			Score float64 `json:"score"`
		}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Errorf("%s: %v", id, err)
			continue
		}
		if body.Score != 100 {
			t.Errorf("%s score = %v, want 100", id, body.Score)
		}
	}
}

func TestRunFingerprintReloadSameSeedSet(t *testing.T) {
	bin := buildAgentGavel(t)
	fake := buildFakeAdapterBin(t)
	root := t.TempDir()

	first := exec.Command(bin, "run",
		"--adapter", fake,
		"--suite", "security",
		"--seeds", "7",
		"--mode", "oracle",
		"--out", root,
		"--run-id", "first",
		"--scenarios", "SEC-001",
	)
	first.Dir = filepath.Join("..", "..")
	first.Env = append(os.Environ(), "GOWORK=off")
	if out, err := first.CombinedOutput(); err != nil {
		t.Fatalf("first run: %v\n%s", err, out)
	}

	fpPath := filepath.Join(root, "results", "first", "fingerprint.json")
	fb, err := os.ReadFile(fpPath)
	if err != nil {
		t.Fatalf("read fingerprint: %v", err)
	}
	var firstFields map[string]string
	if err := json.Unmarshal(fb, &firstFields); err != nil {
		t.Fatal(err)
	}
	wantSeed := firstFields["seed.set"]
	if wantSeed != "0,1,2,3,4,5,6" {
		t.Fatalf("first seed.set = %q, want 0..6", wantSeed)
	}
	wantHash := firstFields["hash"]
	if wantHash == "" {
		t.Fatal("empty fingerprint hash")
	}

	second := exec.Command(bin, "run",
		"--adapter", fake,
		"--suite", "security",
		"--seeds", "99", // ignored when --fingerprint pins seed-set
		"--mode", "oracle",
		"--out", root,
		"--run-id", "second",
		"--scenarios", "SEC-001",
		"--fingerprint", fpPath,
	)
	second.Dir = filepath.Join("..", "..")
	second.Env = append(os.Environ(), "GOWORK=off")
	if out, err := second.CombinedOutput(); err != nil {
		t.Fatalf("second run: %v\n%s", err, out)
	}

	sb, err := os.ReadFile(filepath.Join(root, "results", "second", "fingerprint.json"))
	if err != nil {
		t.Fatal(err)
	}
	var secondFields map[string]string
	if err := json.Unmarshal(sb, &secondFields); err != nil {
		t.Fatal(err)
	}
	if secondFields["seed.set"] != wantSeed {
		t.Fatalf("reloaded seed.set = %q, want %q", secondFields["seed.set"], wantSeed)
	}
	if secondFields["hash"] != wantHash {
		t.Fatalf("reloaded hash = %q, want %q (pins should reproduce fingerprint hash)", secondFields["hash"], wantHash)
	}
}

func TestRunMissingAdapter(t *testing.T) {
	bin := buildAgentGavel(t)
	cmd := exec.Command(bin, "run", "--suite", "security", "--mode", "oracle")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected error, got: %s", out)
	}
	if !strings.Contains(string(out), "--adapter") {
		t.Fatalf("unexpected stderr: %s", out)
	}
}
