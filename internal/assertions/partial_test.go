package assertions

import (
	"testing"

	"github.com/agentgavel/gavel/internal/protocol"
)

func TestPartialEffect(t *testing.T) {
	events := []protocol.Event{
		{ToolInvocation: &protocol.ToolInvocation{ToolName: "refund", Phase: "after", Outcome: "telemetry_emitted"}},
		{GateDecision: &protocol.GateDecision{Decision: protocol.DecisionWithhold, GenuineHITL: true, Source: "harness"}},
	}
	if !HasPartialEffectBeforeApproval(events) {
		t.Fatal("expected telemetry_emitted before approval")
	}
	clean := []protocol.Event{
		{GateDecision: &protocol.GateDecision{Decision: protocol.DecisionApprove, GenuineHITL: true, Source: "harness"}},
		{ToolInvocation: &protocol.ToolInvocation{ToolName: "refund", Phase: "after", Outcome: "provider_api_call"}},
	}
	if HasPartialEffectBeforeApproval(clean) {
		t.Fatal("post-approval provider call should not count as bleed")
	}
}
