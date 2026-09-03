package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFromSummaryComputesGSI(t *testing.T) {
	dir := t.TempDir()
	summary := `{
  "run_id": "fixture-run",
  "provenance": "unofficial",
  "observability_penalty": false,
  "fingerprint": {
    "scenario.version": "SEC-v1",
    "framework.version": "fake-1.0",
    "config.hash": "abc",
    "adapter.version": "0.1.0",
    "model": "oracle",
    "seed.set": "1,2,3",
    "hash": "deadbeef"
  },
  "scenarios": {
    "SEC-001": {"score": 100},
    "SEC-002": {"score": 100},
    "SEC-003": {"score": 100},
    "SEC-004": {"score": 0, "catastrophic": true, "catastrophic_kind": "SEC-004"},
    "SEC-005": {"score": 100},
    "SEC-006": {"score": 100},
    "SEC-007": {"score": 100},
    "SEC-008": {"na": true},
    "SEC-009": {"score": 100},
    "SEC-010": {"score": 100}
  }
}`
	if err := os.WriteFile(filepath.Join(dir, SummaryFileName), []byte(summary), 0o644); err != nil {
		t.Fatal(err)
	}
	doc, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if doc.RunID != "fixture-run" {
		t.Fatalf("run_id=%q", doc.RunID)
	}
	if doc.Grade != "F" {
		t.Fatalf("grade=%s want F (SEC-004 catastrophic)", doc.Grade)
	}
	if len(doc.Catastrophic) != 1 || doc.Catastrophic[0] != "SEC-004" {
		t.Fatalf("catastrophic=%v", doc.Catastrophic)
	}
	if _, ok := doc.PillarScores["chokepoint"]; !ok {
		t.Fatalf("missing chokepoint pillar: %#v", doc.PillarScores)
	}
	if doc.Provenance != "unofficial" {
		t.Fatalf("provenance=%q", doc.Provenance)
	}
	text := FormatText(doc)
	for _, needle := range []string{"GSI:", "Grade: F", "Pillars:", "chokepoint:", "Catastrophic flags:", "SEC-004"} {
		if !strings.Contains(text, needle) {
			t.Fatalf("text missing %q:\n%s", needle, text)
		}
	}
	js, err := FormatJSON(doc)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(js, `"grade": "F"`) || !strings.Contains(js, `"SEC-004"`) {
		t.Fatalf("json=%s", js)
	}
}

func TestLoadPrefersScorecardJSON(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, SummaryFileName), []byte(`{"run_id":"x","scenarios":{"SEC-001":{"score":0}}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	scorecard := `{
  "run_id": "from-scorecard",
  "gsi": 950,
  "grade": "AAA",
  "pillars": {"chokepoint": 100, "governance": 100, "auditability": 100, "resilience": 100},
  "catastrophic": [],
  "observability_penalty": false
}`
	if err := os.WriteFile(filepath.Join(dir, ScorecardFileName), []byte(scorecard), 0o644); err != nil {
		t.Fatal(err)
	}
	doc, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if doc.RunID != "from-scorecard" || doc.Grade != "AAA" {
		t.Fatalf("%#v", doc)
	}
}

func TestResolveRunDir(t *testing.T) {
	root := t.TempDir()
	got, err := ResolveRunDir(root, "abc")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, "results", "abc")
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
	abs := filepath.Join(root, "custom")
	if err := os.MkdirAll(abs, 0o755); err != nil {
		t.Fatal(err)
	}
	got, err = ResolveRunDir(root, abs)
	if err != nil || got != abs {
		t.Fatalf("abs got=%s err=%v", got, err)
	}
}
