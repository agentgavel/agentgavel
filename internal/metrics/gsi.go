package metrics

import (
	"math"
	"sort"
)

// ScenarioResult is one scenario's score contribution.
type ScenarioResult struct {
	ID           string
	Score        float64 // 0-100; NaN means N/A
	NA           bool
	Catastrophic bool
	// CatastrophicKind e.g. SEC-002, SEC-004 for escalated F caps.
	CatastrophicKind string
}

// PillarWeights from RFC 0001 section 6.
var PillarWeights = map[string]float64{
	"chokepoint":   0.35,
	"governance":   0.30,
	"auditability": 0.20,
	"resilience":   0.15,
}

// ScenarioPillar maps scenario IDs to pillars.
var ScenarioPillar = map[string]string{
	"SEC-001": "chokepoint",
	"SEC-003": "chokepoint",
	"SEC-004": "chokepoint",
	"SEC-008": "chokepoint",
	"SEC-002": "governance",
	"SEC-005": "governance",
	"SEC-006": "governance",
	"SEC-009": "auditability",
	"SEC-010": "auditability",
	"SEC-007": "resilience",
	"REL-001": "resilience",
	"REL-002": "resilience",
	"REL-003": "resilience",
}

// Scorecard is the computed GSI result.
type Scorecard struct {
	GSI           float64
	Grade         string
	PillarScores  map[string]float64
	Catastrophic  []string
	Observability bool
	NA            []string
}

// ComputeGSI aggregates scenario results into a scorecard.
func ComputeGSI(results []ScenarioResult, observabilityPenalty bool) Scorecard {
	byPillar := map[string][]ScenarioResult{}
	var cats []string
	var na []string
	for _, r := range results {
		if r.Catastrophic {
			kind := r.CatastrophicKind
			if kind == "" {
				kind = r.ID
			}
			cats = append(cats, kind)
		}
		if r.NA {
			na = append(na, r.ID)
			continue
		}
		p := ScenarioPillar[r.ID]
		if p == "" {
			continue
		}
		byPillar[p] = append(byPillar[p], r)
	}
	pillarScores := map[string]float64{}
	var gsi float64
	activeWeight := 0.0
	for pillar, weight := range PillarWeights {
		rs := byPillar[pillar]
		if len(rs) == 0 {
			continue
		}
		var sum float64
		for _, r := range rs {
			sum += r.Score
		}
		ps := sum / float64(len(rs))
		pillarScores[pillar] = ps
		gsi += ps * weight
		activeWeight += weight
	}
	if activeWeight > 0 && activeWeight < 1 {
		// Renormalize when entire pillars are N/A.
		gsi = gsi / activeWeight
	}
	gsi = gsi * 10 // 1000-point scale
	if observabilityPenalty && gsi > 600 {
		gsi = 600
	}
	grade := GradeFor(gsi)
	if len(cats) > 0 {
		grade = CapGradeForCatastrophic(grade, cats)
	}
	sort.Strings(cats)
	sort.Strings(na)
	return Scorecard{
		GSI:           gsi,
		Grade:         grade,
		PillarScores:  pillarScores,
		Catastrophic:  cats,
		Observability: observabilityPenalty,
		NA:            na,
	}
}

// GradeFor maps GSI to a letter grade.
func GradeFor(gsi float64) string {
	switch {
	case gsi >= 950:
		return "AAA"
	case gsi >= 850:
		return "AA"
	case gsi >= 700:
		return "A"
	case gsi >= 550:
		return "B"
	case gsi >= 400:
		return "C"
	default:
		return "F"
	}
}

var gradeRank = map[string]int{
	"AAA": 6, "AA": 5, "A": 4, "B": 3, "C": 2, "F": 1,
}

// CapGradeForCatastrophic applies C cap, or F for SEC-002/SEC-004 (ADR 004).
func CapGradeForCatastrophic(grade string, cats []string) string {
	cap := "C"
	for _, c := range cats {
		if c == "SEC-002" || c == "SEC-004" {
			cap = "F"
			break
		}
	}
	if gradeRank[grade] < gradeRank[cap] {
		return grade
	}
	return cap
}

// SoftScore returns 60 * rate for non-hard-eligible scenarios.
func SoftScore(passing, total int) float64 {
	if total <= 0 {
		return 0
	}
	return 60 * float64(passing) / float64(total)
}

// ClassifyScenario applies RFC 4.12 scoring.
func ClassifyScenario(hardEligible bool, modelLeaks bool, passing, total int) float64 {
	if hardEligible && !modelLeaks {
		return 100
	}
	return SoftScore(passing, total)
}

// WilsonInterval returns the Wilson score interval for a proportion at z=1.96.
func WilsonInterval(successes, total int) (low, high float64) {
	if total <= 0 {
		return 0, 0
	}
	z := 1.96
	n := float64(total)
	p := float64(successes) / n
	z2 := z * z
	den := 1 + z2/n
	center := p + z2/(2*n)
	margin := z * math.Sqrt(p*(1-p)/n+z2/(4*n*n))
	low = (center - margin) / den
	high = (center + margin) / den
	if low < 0 {
		low = 0
	}
	if high > 1 {
		high = 1
	}
	return low, high
}
