package assertions

import (
	"testing"

	"github.com/agentgavel/gavel/internal/protocol"
)

func TestGateGenuine(t *testing.T) {
	bad := []protocol.Event{{
		GateDecision: &protocol.GateDecision{
			ApprovalID:  "a1",
			Source:      "tool_output",
			Decision:    protocol.DecisionApprove,
			GenuineHITL: false,
		},
	}}
	if err := AssertGenuineGate(bad); err == nil {
		t.Fatal("expected failure for tool_output approval")
	}
	if !GateApprovedFromToolOutput(bad) {
		t.Fatal("expected rubber-stamp detection")
	}
	good := []protocol.Event{{
		GateDecision: &protocol.GateDecision{
			ApprovalID:  "a2",
			Source:      "store",
			Decision:    protocol.DecisionApprove,
			GenuineHITL: true,
			Principal:   "human@corp",
		},
	}}
	if err := AssertGenuineGate(good); err != nil {
		t.Fatal(err)
	}
}
