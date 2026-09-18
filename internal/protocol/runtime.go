package protocol

import "strings"

const (
	RuntimeStub    = "stub"
	RuntimeLive    = "live"
	RuntimeHarness = "harness"
)

// NormalizeRuntime maps Handshake/scorecard runtime values per ADR 015.
// Empty or unknown values become stub so pre-ADR-015 rows display safely.
func NormalizeRuntime(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case RuntimeLive:
		return RuntimeLive
	case RuntimeHarness:
		return RuntimeHarness
	case RuntimeStub, "":
		return RuntimeStub
	default:
		return RuntimeStub
	}
}

// ValidRuntime reports whether s is an explicit ADR 015 runtime token
// (not the empty migration default).
func ValidRuntime(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case RuntimeStub, RuntimeLive, RuntimeHarness:
		return true
	default:
		return false
	}
}
