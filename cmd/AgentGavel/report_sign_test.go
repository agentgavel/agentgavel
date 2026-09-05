package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestReportSignAndVerifyEntry(t *testing.T) {
	bin := buildAgentGavel(t)
	repoRoot := findRepoRoot(t)
	key := filepath.Join(repoRoot, "internal", "submit", "testdata", "example-framework-test-1.priv.b64")
	registry := filepath.Join(repoRoot, "dashboard", "keys", "registry.json")

	entry := map[string]any{
		"run_id":          "cli-sign-test",
		"framework":       "Example Framework",
		"adapter":         "example",
		"adapter_version": "1.0.0",
		"provenance":      "ratified",
		"tab":             "opt-in",
		"sample":          false,
		"gsi":             900.0,
		"grade":           "AA",
		"pillars":         map[string]float64{"chokepoint": 90},
		"catastrophic":    []string{},
		"na":              []string{},
		"fingerprint":     map[string]string{"hash": "abcd"},
		"generated_at":    "2026-09-05T12:00:00Z",
	}
	raw, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	inPath := filepath.Join(t.TempDir(), "in.json")
	outPath := filepath.Join(t.TempDir(), "out.json")
	if err := os.WriteFile(inPath, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	sign := exec.Command(bin, "report", "--sign",
		"--key", key,
		"--key-id", "example-framework-test-1",
		"--entry", inPath,
		"--out", outPath,
	)
	sign.Dir = repoRoot
	if out, err := sign.CombinedOutput(); err != nil {
		t.Fatalf("report --sign: %v\n%s", err, out)
	}
	signed, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(signed, &m); err != nil {
		t.Fatal(err)
	}
	if m["key_id"] != "example-framework-test-1" {
		t.Fatalf("key_id = %v", m["key_id"])
	}
	if _, ok := m["signature"].(string); !ok || m["signature"] == "" {
		t.Fatalf("missing signature: %v", m["signature"])
	}

	verify := exec.Command(bin, "verify-entry", "--registry", registry, outPath)
	verify.Dir = repoRoot
	if out, err := verify.CombinedOutput(); err != nil {
		t.Fatalf("verify-entry: %v\n%s", err, out)
	}

	// Tamper
	m["gsi"] = 1.0
	tampered, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	badPath := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(badPath, tampered, 0o644); err != nil {
		t.Fatal(err)
	}
	bad := exec.Command(bin, "verify-entry", "--registry", registry, badPath)
	bad.Dir = repoRoot
	out, err := bad.CombinedOutput()
	if err == nil {
		t.Fatalf("expected verify-entry failure on tamper, got: %s", out)
	}
	if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 1 {
		t.Fatalf("want exit 1, got %v\n%s", err, out)
	}
}

func TestHelpListsSignAndVerifyEntry(t *testing.T) {
	bin := buildAgentGavel(t)
	cmd := exec.Command(bin, "help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("help: %v\n%s", err, out)
	}
	s := string(out)
	if !strings.Contains(s, "verify-entry") {
		t.Fatalf("help missing verify-entry:\n%s", s)
	}
	if !strings.Contains(s, "--sign") && !strings.Contains(s, "sign") {
		t.Fatalf("help missing sign mention:\n%s", s)
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}
