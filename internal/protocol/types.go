// Package protocol implements the AgentGavel adapter wire types and codecs.
package protocol

import (
	"encoding/json"
	"fmt"
)

// Decision is an approval outcome.
type Decision string

const (
	DecisionUnspecified Decision = ""
	DecisionApprove     Decision = "approve"
	DecisionDeny        Decision = "deny"
	DecisionWithhold    Decision = "withhold"
)

// Empty is an empty response.
type Empty struct{}

// HandshakeRequest is sent by the engine at connection start.
type HandshakeRequest struct {
	EngineProtocolVersion string `json:"engine_protocol_version"`
	EngineVersion         string `json:"engine_version,omitempty"`
}

// CapabilityReport is returned by Handshake.
type CapabilityReport struct {
	AdapterProtocolVersion string `json:"adapter_protocol_version"`
	AdapterName            string `json:"adapter_name"`
	AdapterVersion         string `json:"adapter_version"`
	Provenance             string `json:"provenance"`
	HITL                   bool   `json:"hitl"`
	Tenancy                bool   `json:"tenancy"`
	Ledger                 bool   `json:"ledger"`
	Observability          bool   `json:"observability"`
	ContextMode            string `json:"context_mode"`
	FrameworkName          string `json:"framework_name,omitempty"`
	FrameworkVersion       string `json:"framework_version,omitempty"`
}

// SessionConfig configures a target session.
type SessionConfig struct {
	ModelBaseURL        string            `json:"model_base_url"`
	ModelName           string            `json:"model_name,omitempty"`
	MCPEndpoints        []MCPEndpoint     `json:"mcp_endpoints,omitempty"`
	TenantID            string            `json:"tenant_id,omitempty"`
	FrameworkConfigHash string            `json:"framework_config_hash,omitempty"`
	RunMode             string            `json:"run_mode,omitempty"`
	Extra               map[string]string `json:"extra,omitempty"`
}

// MCPEndpoint describes a fixture MCP server.
type MCPEndpoint struct {
	Name      string   `json:"name"`
	Transport string   `json:"transport,omitempty"`
	Command   string   `json:"command,omitempty"`
	Args      []string `json:"args,omitempty"`
	URL       string   `json:"url,omitempty"`
}

// SessionID identifies a session.
type SessionID struct {
	ID string `json:"id"`
}

// Task is a scenario task payload.
type Task struct {
	ID       string            `json:"id"`
	Prompt   string            `json:"prompt"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// SubmitTaskRequest wraps a task submission.
type SubmitTaskRequest struct {
	SessionID string `json:"session_id"`
	Task      Task   `json:"task"`
}

// ResolveApprovalRequest is a harness-driven approval decision.
type ResolveApprovalRequest struct {
	SessionID  string   `json:"session_id"`
	ApprovalID string   `json:"approval_id"`
	Decision   Decision `json:"decision"`
	Principal  string   `json:"principal,omitempty"`
}

// Event is a tagged union pushed by the adapter.
type Event struct {
	SessionID string `json:"session_id"`
	Seq       uint64 `json:"seq"`
	UnixMs    int64  `json:"unix_ms"`

	ToolInvocation     *ToolInvocation     `json:"tool_invocation,omitempty"`
	GateDecision       *GateDecision       `json:"gate_decision,omitempty"`
	ContextSnapshot    *ContextSnapshot    `json:"context_snapshot,omitempty"`
	ContextAttestation *ContextAttestation `json:"context_attestation,omitempty"`
	LedgerAppend       *LedgerAppend       `json:"ledger_append,omitempty"`
	SessionError       *SessionError       `json:"session_error,omitempty"`
}

// Kind returns the event kind name for the set payload.
func (e Event) Kind() string {
	switch {
	case e.ToolInvocation != nil:
		return "tool_invocation"
	case e.GateDecision != nil:
		return "gate_decision"
	case e.ContextSnapshot != nil:
		return "context_snapshot"
	case e.ContextAttestation != nil:
		return "context_attestation"
	case e.LedgerAppend != nil:
		return "ledger_append"
	case e.SessionError != nil:
		return "session_error"
	default:
		return ""
	}
}

// ToolInvocation records a tool call before or after dispatch.
type ToolInvocation struct {
	ToolName      string `json:"tool_name"`
	ToolID        string `json:"tool_id,omitempty"`
	Phase         string `json:"phase"` // before | after
	ArgumentsJSON string `json:"arguments_json,omitempty"`
	Outcome       string `json:"outcome,omitempty"`
	Error         string `json:"error,omitempty"`
	Refused       bool   `json:"refused,omitempty"`
}

// GateDecision records an approval gate outcome.
type GateDecision struct {
	ApprovalID  string   `json:"approval_id"`
	Source      string   `json:"source"` // store | tool_output | llm | harness
	Decision    Decision `json:"decision"`
	Principal   string   `json:"principal,omitempty"`
	GenuineHITL bool     `json:"genuine_hitl"`
}

// ContextSnapshot is a raw context fragment.
type ContextSnapshot struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ContextAttestation carries hashed n-grams (ADR 005).
type ContextAttestation struct {
	NgramHashes []string `json:"ngram_hashes"`
	Algorithm   string   `json:"algorithm"`
}

// LedgerAppend is a ledger entry notification.
type LedgerAppend struct {
	EntryJSON string `json:"entry_json"`
}

// SessionError is a session-level error event.
type SessionError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Ledger is the exported audit ledger.
type Ledger struct {
	SessionID string        `json:"session_id"`
	Entries   []LedgerEntry `json:"entries"`
}

// LedgerEntry is one ledger row.
type LedgerEntry struct {
	ID          string `json:"id"`
	UnixMs      int64  `json:"unix_ms"`
	Kind        string `json:"kind"`
	PayloadJSON string `json:"payload_json,omitempty"`
	PrevHash    string `json:"prev_hash,omitempty"`
	Hash        string `json:"hash,omitempty"`
}

// MarshalJSON encodes v with encoding/json.
func MarshalJSON(v any) ([]byte, error) {
	return json.Marshal(v)
}

// UnmarshalJSON decodes data into v.
func UnmarshalJSON(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// RoundTripJSON marshals then unmarshals into a fresh value of the same type.
func RoundTripJSON[T any](in T) (T, error) {
	var out T
	b, err := json.Marshal(in)
	if err != nil {
		return out, fmt.Errorf("marshal: %w", err)
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return out, fmt.Errorf("unmarshal: %w", err)
	}
	return out, nil
}
