package engine

import "testing"

func TestFingerprintStableAcrossSeedOrder(t *testing.T) {
	a := Fingerprint{
		ScenarioVersion:  "v1",
		FrameworkVersion: "v2",
		ConfigHash:       "abc123",
		AdapterVersion:   "v3",
		SeedSet:          []int64{3, 1, 2},
	}
	b := Fingerprint{
		ScenarioVersion:  "v1",
		FrameworkVersion: "v2",
		ConfigHash:       "abc123",
		AdapterVersion:   "v3",
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
		SeedSet:          []int64{1, 2, 3},
	}
	changed := base
	changed.AdapterVersion = "v4"
	if base.Hash() == changed.Hash() {
		t.Fatal("expected different hashes when adapter version changes")
	}
}
