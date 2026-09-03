// Package reliability is the REL-v0 suite catalog (ADR 010): demotion
// latches, replay/config drift, and ledger completeness.
package reliability

// SuiteVersion is the pinned reliability suite version for v0.3 (REL-001..003).
const SuiteVersion = "REL-v0"

// Definition is a framework-agnostic scenario definition.
// Fixture paths are relative to the repo root; predicates are not attached here.
type Definition struct {
	ID       string
	Name     string
	Fixtures []string
}

// Suite is a loaded reliability suite pin: version plus scenario definitions.
type Suite struct {
	Version   string
	Scenarios []Definition
}

// catalog is the REL-v0 scenario list, pinned by ADR 010. Fixtures land with
// their scoring predicates (demotion probe, drifted replay plan, ledger-gap
// corpus) and may be empty until those tasks complete.
var catalog = []Definition{
	{
		ID:   "REL-001",
		Name: "Demotion Latch",
	},
	{
		ID:   "REL-002",
		Name: "Replay / Config Drift",
	},
	{
		ID:   "REL-003",
		Name: "Ledger Completeness",
	},
}

// Load returns the REL-v0 reliability suite: pinned version and REL-001..003 defs.
func Load() (Suite, error) {
	scenarios := make([]Definition, len(catalog))
	copy(scenarios, catalog)
	return Suite{
		Version:   SuiteVersion,
		Scenarios: scenarios,
	}, nil
}
