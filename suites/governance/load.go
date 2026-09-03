// Package governance is the GOV-v0 suite catalog and policy-ceiling stubs.
package governance

// SuiteVersion is the pinned governance suite version for v0.2 (GOV-001 stub).
const SuiteVersion = "GOV-v0"

// Definition is a framework-agnostic scenario definition.
// Fixture paths are relative to the repo root; predicates are not attached here.
type Definition struct {
	ID       string
	Name     string
	Fixtures []string
}

// Suite is a loaded governance suite pin: version plus scenario definitions.
type Suite struct {
	Version   string
	Scenarios []Definition
}

// catalog is the GOV-v0 scenario list. GOV-001 is a stub until full
// hard-governance scoring lands.
var catalog = []Definition{
	{
		ID:   "GOV-001",
		Name: "Policy Ceiling Enforcement",
	},
}

// Load returns the GOV-v0 governance suite: pinned version and GOV-001 def.
func Load() (Suite, error) {
	scenarios := make([]Definition, len(catalog))
	copy(scenarios, catalog)
	return Suite{
		Version:   SuiteVersion,
		Scenarios: scenarios,
	}, nil
}
