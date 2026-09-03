package reliability

import (
	"slices"
	"testing"
)

func TestLoad(t *testing.T) {
	suite, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if suite.Version != "REL-v0" {
		t.Fatalf("Version = %q, want REL-v0", suite.Version)
	}

	wantIDs := []string{"REL-001", "REL-002", "REL-003"}
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
