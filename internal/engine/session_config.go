package engine

import (
	"context"
	"fmt"

	"github.com/agentgavel/agentgavel/internal/protocol"
)

// Evaluation modes (ADR 003 / RFC 0001 §4.12).
const (
	ModeOracle = "oracle"
	ModeModel  = "model"
)

// ModeEndpoints holds the model base URLs for each evaluation mode.
// Oracle mode points adapters at the Compliance Oracle listen URL so the
// framework cannot special-case oracle behavior (ADR 003).
type ModeEndpoints struct {
	OracleURL string // Compliance Oracle HTTP base URL (e.g. http://127.0.0.1:8080)
	ModelURL  string // Pinned real-model OpenAI/Anthropic-compatible base URL
}

// SessionConfigForMode builds a protocol.SessionConfig with ModelBaseURL and
// RunMode set from the evaluation mode. The Oracle is injected only via
// model_base_url — it is not part of the adapter contract.
func SessionConfigForMode(mode string, endpoints ModeEndpoints) (protocol.SessionConfig, error) {
	cfg := protocol.SessionConfig{RunMode: mode}
	switch mode {
	case ModeOracle:
		if endpoints.OracleURL == "" {
			return protocol.SessionConfig{}, fmt.Errorf("engine: oracle mode requires OracleURL")
		}
		cfg.ModelBaseURL = endpoints.OracleURL
	case ModeModel:
		if endpoints.ModelURL == "" {
			return protocol.SessionConfig{}, fmt.Errorf("engine: model mode requires ModelURL")
		}
		cfg.ModelBaseURL = endpoints.ModelURL
	default:
		return protocol.SessionConfig{}, fmt.Errorf("engine: unknown run mode %q (want %q or %q)", mode, ModeOracle, ModeModel)
	}
	return cfg, nil
}

// ApplyModeConfig sets sess.Config from sess.Mode and the given endpoints.
func ApplyModeConfig(sess *Session, endpoints ModeEndpoints) error {
	if sess == nil {
		return fmt.Errorf("engine: nil session")
	}
	cfg, err := SessionConfigForMode(sess.Mode, endpoints)
	if err != nil {
		return err
	}
	sess.Config = cfg
	return nil
}

// StartSessionForMode builds SessionConfig for mode and starts the session on
// the launched adapter. This is the launcher path that injects ModelBaseURL.
func (s *AdapterSession) StartSessionForMode(ctx context.Context, mode string, endpoints ModeEndpoints) (protocol.SessionID, protocol.SessionConfig, error) {
	if s == nil || s.Client == nil {
		return protocol.SessionID{}, protocol.SessionConfig{}, fmt.Errorf("engine: adapter session not launched")
	}
	cfg, err := SessionConfigForMode(mode, endpoints)
	if err != nil {
		return protocol.SessionID{}, protocol.SessionConfig{}, err
	}
	id, err := s.Client.StartSession(ctx, cfg)
	if err != nil {
		return protocol.SessionID{}, cfg, fmt.Errorf("engine: start session: %w", err)
	}
	return id, cfg, nil
}
