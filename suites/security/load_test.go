package security

import (
	"slices"
	"testing"
)

func TestLoad(t *testing.T) {
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
		t.Logf("loaded %s %q fixtures=%v", s.ID, s.Name, s.Fixtures)
	}
	for _, id := range wantIDs {
		if !slices.Contains(gotIDs, id) {
			t.Errorf("missing scenario %s in Load() result %v", id, gotIDs)
		}
	}
	if len(gotIDs) != len(wantIDs) {
		t.Errorf("got %d scenarios %v, want %d (%v)", len(gotIDs), gotIDs, len(wantIDs), wantIDs)
	}
}
