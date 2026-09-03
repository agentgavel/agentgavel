package engine

import (
	"context"
	"fmt"

	"github.com/agentgavel/gavel/internal/mcpfuzz"
	"github.com/agentgavel/gavel/internal/protocol"
)

// FuzzModeHandle owns a running mcpfuzz mode listener and the SessionConfig
// slice entry that points adapters at it (RFC 0001 SEC-007 / UC-013).
type FuzzModeHandle struct {
	Mode     string
	Endpoint protocol.MCPEndpoint
	Config   protocol.SessionConfig

	running *mcpfuzz.RunningMode
}

// StartFuzzMode launches a named mcpfuzz mode on a local TCP listener and
// returns a SessionConfig with that dialable endpoint in MCPEndpoints.
// The caller must Close the handle when finished.
func StartFuzzMode(ctx context.Context, mode string) (*FuzzModeHandle, error) {
	if mode == "" {
		return nil, fmt.Errorf("engine: empty fuzz mode name")
	}
	running, err := mcpfuzz.StartMode(ctx, mode)
	if err != nil {
		return nil, fmt.Errorf("engine: start fuzz mode %q: %w", mode, err)
	}
	ep := protocol.MCPEndpoint{
		Name:      mode,
		Transport: "tcp",
		URL:       running.URL,
	}
	cfg := protocol.SessionConfig{
		MCPEndpoints: []protocol.MCPEndpoint{ep},
	}
	return &FuzzModeHandle{
		Mode:     mode,
		Endpoint: ep,
		Config:   cfg,
		running:  running,
	}, nil
}

// Addr returns the host:port suitable for net.Dial("tcp", ...).
func (h *FuzzModeHandle) Addr() string {
	if h == nil || h.running == nil {
		return ""
	}
	return h.running.Addr
}

// Close shuts down the fuzz mode listener.
func (h *FuzzModeHandle) Close() error {
	if h == nil || h.running == nil {
		return nil
	}
	return h.running.Close()
}

// InjectFuzzMode appends the fuzz endpoint into sess.Config.MCPEndpoints,
// starting the mode if needed. Prefer StartFuzzMode when you need the handle.
func InjectFuzzMode(ctx context.Context, sess *Session, mode string) (*FuzzModeHandle, error) {
	if sess == nil {
		return nil, fmt.Errorf("engine: nil session")
	}
	h, err := StartFuzzMode(ctx, mode)
	if err != nil {
		return nil, err
	}
	sess.Config.MCPEndpoints = append(sess.Config.MCPEndpoints, h.Endpoint)
	return h, nil
}
