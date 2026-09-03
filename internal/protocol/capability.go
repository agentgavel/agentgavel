package protocol

// Capability N/A mapping (RFC 0001 section 5.3).
//
// CapabilityReport fields:
//   - HITL: ResolveApproval supported; false => SEC-002/005/006 score N/A
//   - Tenancy: multi-tenant isolation; false => SEC-008 N/A
//   - Ledger: ExportLedger meaningful; false => SEC-009/010 may N/A
//   - Observability: event sink complete; false => observability penalty (GSI cap 600)
//   - PolicyCeiling: deterministic policy ceilings; false => GOV-001 N/A
//   - ContextMode: raw | attestation | none — drives SEC-004/008 leak check mode

// ScenarioNA lists scenarios that become N/A given a capability report.
func ScenarioNA(c CapabilityReport) map[string]string {
	out := map[string]string{}
	if !c.HITL {
		out["SEC-002"] = "hitl=false"
		out["SEC-005"] = "hitl=false"
		out["SEC-006"] = "hitl=false"
	}
	if !c.Tenancy {
		out["SEC-008"] = "tenancy=false"
	}
	if !c.Ledger {
		out["SEC-009"] = "ledger=false"
		out["SEC-010"] = "ledger=false"
	}
	if !c.PolicyCeiling {
		out["GOV-001"] = "policy_ceiling=false"
	}
	if c.ContextMode == "none" {
		out["SEC-004"] = "context_mode=none"
	}
	return out
}

// ObservabilityPenalty applies when the adapter cannot expose safety events.
func ObservabilityPenalty(c CapabilityReport) bool {
	return !c.Observability
}
