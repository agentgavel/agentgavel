package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestFingerprintStableAcrossSeedOrder(t *testing.T) {
	a := Fingerprint{
		ScenarioVersion:  "v1",
		FrameworkVersion: "v2",
		ConfigHash:       "abc123",
		AdapterVersion:   "v3",
		Model:            "oracle",
		SeedSet:          []int64{3, 1, 2},
	}
	b := Fingerprint{
		ScenarioVersion:  "v1",
		FrameworkVersion: "v2",
		ConfigHash:       "abc123",
		AdapterVersion:   "v3",
		Model:            "oracle",
		SeedSet:          []int64{1, 2, 3},
	}
	if a.Hash() != b.Hash() {
		t.Fatalf("expected equal hashes for equivalent seed sets, got %s vs %s", a.Hash(), b.Hash())
	}
	if a.Hash() == "" {
		t.Fatal("expected non-empty hash")
	}
}

func TestFingerprintDiffersOnFieldChange(t *testing.T) {
	base := Fingerprint{
		ScenarioVersion:  "v1",
		FrameworkVersion: "v2",
		ConfigHash:       "abc123",
		AdapterVersion:   "v3",
		Model:            "oracle",
		SeedSet:          []int64{1, 2, 3},
	}
	changed := base
	changed.AdapterVersion = "v4"
	if base.Hash() == changed.Hash() {
		t.Fatal("expected different hashes when adapter version changes")
	}
}

func TestParseSeedSet(t *testing.T) {
	got, err := ParseSeedSet(" 3, 1,2 ")
	if err != nil {
		t.Fatal(err)
	}
	want := []int64{3, 1, 2}
	if len(got) != len(want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v want %v", got, want)
		}
	}
}

func TestLoadFingerprintFileRoundTrip(t *testing.T) {
	dir := t.TempDir()
	fp := Fingerprint{
		ScenarioVersion:  "SEC-v1",
		FrameworkVersion: "0.1.0",
		ConfigHash:       "cfg",
		AdapterVersion:   "0.0.1",
		Model:            "oracle",
		SeedSet:          []int64{7, 8, 9},
	}
	path, err := WriteFingerprintFile(dir, fp.Fields())
	if err != nil {
		t.Fatal(err)
	}
	got, err := LoadFingerprintFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Hash() != fp.Hash() {
		t.Fatalf("hash %s want %s", got.Hash(), fp.Hash())
	}
	if len(got.SeedSet) != 3 || got.SeedSet[0] != 7 {
		t.Fatalf("seed set %#v", got.SeedSet)
	}
}

func TestLoadFingerprintFileFromSummary(t *testing.T) {
	dir := t.TempDir()
	fp := Fingerprint{
		Model:   "oracle",
		SeedSet: []int64{0, 1, 2},
	}
	summary := filepath.Join(dir, "summary.json")
	body := `{"run_id":"r1","fingerprint":` + mustJSON(t, fp.Fields()) + `,"scenarios":{}}`
	if err := os.WriteFile(summary, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := LoadFingerprintFile(summary)
	if err != nil {
		t.Fatal(err)
	}
	if got.Fields()["seed.set"] != "0,1,2" {
		t.Fatalf("seed.set = %q", got.Fields()["seed.set"])
	}
}

func mustJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
