// Package oracle is the Compliance Oracle HTTP service (RFC 0001 section 4.12).
//
// Probe directive binding (authoritative detail in README.md in this directory):
//
//  1. Header X-AgentGavel-Probe-Directive: JSON {"tool_name","arguments"}
//  2. Else a system message starting with AGENTGAVEL_PROBE_DIRECTIVE: + JSON
//
// A missing directive must yield 4xx; the Oracle must not invent actions.
package oracle

// Directive is the probe instruction the Oracle must obey on the next completion.
type Directive struct {
	ToolName  string         `json:"tool_name"`
	Arguments map[string]any `json:"arguments"`
}

const (
	// HeaderProbeDirective is the preferred binding for probe instructions.
	HeaderProbeDirective = "X-AgentGavel-Probe-Directive"
	// SystemPrefixProbeDirective prefixes a system message carrying the directive JSON.
	SystemPrefixProbeDirective = "AGENTGAVEL_PROBE_DIRECTIVE:"
)
