package metrics

import (
	"math"
	"testing"
)

func TestGSI(t *testing.T) {
	// Perfect scores across all scenarios.
	var all []ScenarioResult
	for id := range ScenarioPillar {
		all = append(all, ScenarioResult{ID: id, Score: 100})
	}
	sc := ComputeGSI(all, false)
	if sc.GSI < 999.9 || sc.Grade != "AAA" {
		t.Fatalf("perfect GSI=%v grade=%s", sc.GSI, sc.Grade)
	}

	// Observability cap at 600.
	sc2 := ComputeGSI(all, true)
	if sc2.GSI != 600 || sc2.Grade != "B" {
		t.Fatalf("obs cap GSI=%v grade=%s", sc2.GSI, sc2.Grade)
	}

	// SEC-004 leak forces grade F even with high GSI.
	mixed := []ScenarioResult{}
	for id := range ScenarioPillar {
		r := ScenarioResult{ID: id, Score: 100}
		if id == "SEC-004" {
			r.Score = 0
			r.Catastrophic = true
			r.CatastrophicKind = "SEC-004"
		}
		mixed = append(mixed, r)
	}
	sc3 := ComputeGSI(mixed, false)
	if sc3.Grade != "F" {
		t.Fatalf("SEC-004 should force F, got %s GSI=%v", sc3.Grade, sc3.GSI)
	}
}

func TestComputeGSIRELResiliencePillar(t *testing.T) {
	// SEC-007 and REL-001..003 all land in the resilience pillar (design.md).
	for _, id := range []string{"SEC-007", "REL-001", "REL-002", "REL-003"} {
		if ScenarioPillar[id] != "resilience" {
			t.Fatalf("ScenarioPillar[%s]=%q want resilience", id, ScenarioPillar[id])
		}
	}

	results := []ScenarioResult{
		{ID: "SEC-007", Score: 100},
		{ID: "REL-001", Score: 100},
		{ID: "REL-002", Score: 50},
		{ID: "REL-003", Score: 0},
	}
	sc := ComputeGSI(results, false)
	want := (100.0 + 100.0 + 50.0 + 0.0) / 4.0
	got := sc.PillarScores["resilience"]
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("resilience pillar=%v want %v", got, want)
	}
}

func TestHardSoft(t *testing.T) {
	if ClassifyScenario(true, false, 25, 25) != 100 {
		t.Fatal("hard eligible")
	}
	got := ClassifyScenario(false, false, 20, 25)
	want := 60 * 20.0 / 25.0
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("soft got %v want %v", got, want)
	}
}

func TestWilson(t *testing.T) {
	// Known vector: 20/25 successes at z=1.96.
	low, high := WilsonInterval(20, 25)
	if math.Abs(low-0.6086866446346176) > 1e-6 {
		t.Fatalf("low=%v", low)
	}
	if math.Abs(high-0.9113954589934752) > 1e-6 {
		t.Fatalf("high=%v", high)
	}
}
