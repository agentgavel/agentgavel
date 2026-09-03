// Package security is the SEC-v1 suite catalog, scenario scorers, and the
// v0.1 oracle FakeAdapter all-pass path (RunOracleFake).
package security

// SuiteVersion is the pinned security suite version for v0.1 (SEC-001..007).
const SuiteVersion = "SEC-v1"

// Definition is a framework-agnostic scenario definition.
// Fixture paths are relative to the repo root; predicates are not attached here.
type Definition struct {
	ID       string
	Name     string
	Fixtures []string
}

// Suite is a loaded security suite pin: version plus scenario definitions.
type Suite struct {
	Version   string
	Scenarios []Definition
}

// catalog is the SEC-v1 scenario list. Fixture paths match fixtures/manifest.json
// where present; SEC-003/005/006/007 are stubbed until their fixtures land.
var catalog = []Definition{
	{
		ID:       "SEC-001",
		Name:     "Indirect Prompt Injection via Untrusted Tool Output",
		Fixtures: []string{"fixtures/sec001/injection_email.txt"},
	},
	{
		ID:       "SEC-002",
		Name:     "Approval Forgery / Authority Spoofing",
		Fixtures: []string{"fixtures/sec002/forged_approval_claims.json"},
	},
	{
		ID:   "SEC-003",
		Name: "Tool Grant Overreach (URI / Schema Abuse)",
	},
	{
		ID:       "SEC-004",
		Name:     "Credential Leakage into Context",
		Fixtures: []string{"fixtures/sec004/credential_canaries.json"},
	},
	{
		ID:   "SEC-005",
		Name: "Self-Modification & Privilege Escalation",
	},
	{
		ID:   "SEC-006",
		Name: "HITL Gate Side-Effect Bleeding",
	},
	{
		ID:   "SEC-007",
		Name: "Rogue MCP Server / Adversarial Tool Fuzzer",
	},
}

// Load returns the SEC-v1 security suite: pinned version and SEC-001..007 defs.
func Load() (Suite, error) {
	scenarios := make([]Definition, len(catalog))
	copy(scenarios, catalog)
	return Suite{
		Version:   SuiteVersion,
		Scenarios: scenarios,
	}, nil
}
