package mcpfuzz

import (
	"fmt"
	"io"
	"time"
)

// Named SEC-007 fuzz modes (RFC 0001 §4.7).
const (
	ModeToxicOutput     = "toxic-output"
	ModeSchemaFlooder   = "schema-flooder"
	ModeEarlyDisconnect = "early-disconnect"
	ModeToolRenamer     = "tool-renamer"
	ModeSlowloris       = "slowloris"
	ModeMasquerade      = "masquerade"
)

// DefaultSlowlorisDelay is used when starting slowloris via NewByName.
const DefaultSlowlorisDelay = 200 * time.Millisecond

// Modes lists registered fuzz mode names.
func Modes() []string {
	return []string{
		ModeToxicOutput,
		ModeSchemaFlooder,
		ModeEarlyDisconnect,
		ModeToolRenamer,
		ModeSlowloris,
		ModeMasquerade,
	}
}

// NewByName constructs a fuzz mode server by RFC name.
func NewByName(name string, in io.Reader, out io.Writer) (*Server, error) {
	switch name {
	case ModeToxicOutput:
		return NewToxicOutput(in, out), nil
	case ModeSchemaFlooder:
		return NewSchemaFlooder(in, out), nil
	case ModeEarlyDisconnect:
		return NewEarlyDisconnect(in, out), nil
	case ModeToolRenamer:
		return NewToolRenamer(in, out), nil
	case ModeSlowloris:
		return NewSlowloris(in, out, DefaultSlowlorisDelay), nil
	case ModeMasquerade:
		return NewMasquerade(in, out), nil
	default:
		return nil, fmt.Errorf("mcpfuzz: unknown mode %q", name)
	}
}
