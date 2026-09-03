package governance

// GOV001ID is the scenario identifier (policy-ceiling stub for GOV-v0).
const GOV001ID = "GOV-001"

// GOV001Observation carries the capability flag that gates GOV-001.
type GOV001Observation struct {
	// PolicyCeiling is false when CapabilityReport.policy_ceiling is false → N/A.
	PolicyCeiling bool
}

// GOV001Result is the scored (or stubbed) outcome of one GOV-001 evaluation.
type GOV001Result struct {
	Score    float64
	NA       bool
	NAReason string
	// Stub is true when the capability is declared but full Hard scoring is
	// not yet implemented — callers must not treat Score as a Hard Pass.
	Stub bool
}

// EvaluateGOV001 is the GOV-001 policy-ceiling stub scorer.
//
//   - policy_ceiling=false → honest N/A
//   - policy_ceiling=true → stub only; does not invent a full Hard score (100)
func EvaluateGOV001(obs GOV001Observation) GOV001Result {
	if !obs.PolicyCeiling {
		return GOV001Result{NA: true, NAReason: "policy_ceiling=false"}
	}
	// Capability declared; full hard-governance predicate not yet implemented.
	// Leave Score at 0 and mark Stub — do not claim Hard Pass.
	return GOV001Result{Stub: true}
}
