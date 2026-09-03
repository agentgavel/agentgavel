package assertions

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/agentgavel/gavel/internal/protocol"
)

// FindToolInvocations returns tool_invocation events matching name (empty = all).
func FindToolInvocations(events []protocol.Event, name string) []protocol.ToolInvocation {
	var out []protocol.ToolInvocation
	for _, e := range events {
		if e.ToolInvocation == nil {
			continue
		}
		if name != "" && e.ToolInvocation.ToolName != name {
			continue
		}
		out = append(out, *e.ToolInvocation)
	}
	return out
}

// ToolFired reports whether a non-refused after-phase invocation occurred.
func ToolFired(events []protocol.Event, name string) bool {
	for _, inv := range FindToolInvocations(events, name) {
		if inv.Phase == "after" && !inv.Refused {
			return true
		}
		if inv.Phase == "" && !inv.Refused {
			return true
		}
	}
	return false
}

// ToolRecipient extracts a string field (e.g. "to") from invocation arguments JSON.
func ToolRecipient(inv protocol.ToolInvocation, field string) (string, error) {
	if inv.ArgumentsJSON == "" {
		return "", fmt.Errorf("empty arguments")
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(inv.ArgumentsJSON), &m); err != nil {
		return "", err
	}
	v, ok := m[field]
	if !ok {
		return "", fmt.Errorf("field %q missing", field)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("field %q not string", field)
	}
	return s, nil
}

// HasToolRecipient true if any matching tool invocation sends to recipient.
func HasToolRecipient(events []protocol.Event, toolName, field, recipient string) bool {
	for _, inv := range FindToolInvocations(events, toolName) {
		got, err := ToolRecipient(inv, field)
		if err != nil {
			continue
		}
		if strings.EqualFold(got, recipient) {
			return true
		}
	}
	return false
}
