package engine

import (
	"context"
	"testing"
	"time"

	"github.com/agentgavel/agentgavel/internal/protocol"
)

// TestIntegrationNoopScenario launches the fake adapter, starts an oracle-mode
// session, submits a task, drains emitted events, scores the noop scenario, and
// stops cleanly (UC-001, UC-004).
func TestIntegrationNoopScenario(t *testing.T) {
	bin := buildFakeAdapter(t)
	ctx := context.Background()
	adapter, err := Launch(ctx, LaunchConfig{Command: bin, Timeout: 5 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = adapter.Close() }()

	if _, err := adapter.Handshake(ctx, protocol.HandshakeRequest{EngineProtocolVersion: "1.0"}); err != nil {
		t.Fatal(err)
	}

	const oracleURL = "http://127.0.0.1:9"
	sid, cfg, err := adapter.StartSessionForMode(ctx, ModeOracle, ModeEndpoints{OracleURL: oracleURL})
	if err != nil {
		t.Fatal(err)
	}
	if sid.ID == "" {
		t.Fatal("empty session id")
	}
	if cfg.ModelBaseURL != oracleURL || cfg.RunMode != ModeOracle {
		t.Fatalf("session config %#v", cfg)
	}

	if err := adapter.Client.SubmitTask(ctx, protocol.SubmitTaskRequest{
		SessionID: sid.ID,
		Task:      protocol.Task{ID: "noop-1", Prompt: "ping"},
	}); err != nil {
		t.Fatal(err)
	}

	events, err := DrainEvents(adapter.Client.Conn, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 {
		t.Fatalf("events=%d want 2", len(events))
	}
	if err := protocol.CheckToolInvocationOrder(events); err != nil {
		t.Fatal(err)
	}
	for i, ev := range events {
		if ev.SessionID != sid.ID {
			t.Fatalf("event[%d] session %q want %q", i, ev.SessionID, sid.ID)
		}
		if ev.Kind() != "tool_invocation" {
			t.Fatalf("event[%d] kind %q", i, ev.Kind())
		}
	}

	scenario := &noopScenario{}
	sess := &Session{
		ID:     sid.ID,
		Seed:   0,
		Mode:   ModeOracle,
		Client: adapter.Client,
		Config: cfg,
	}
	if err := scenario.Setup(ctx, sess); err != nil {
		t.Fatal(err)
	}
	if err := scenario.Probe(ctx, sess); err != nil {
		t.Fatal(err)
	}
	observed, err := scenario.Observe(ctx, sess)
	if err != nil {
		t.Fatal(err)
	}
	// noop Observe is a no-op; use drained adapter events for predicates.
	if observed == nil {
		observed = events
	}
	result, err := scenario.Predicate(ctx, sess, observed)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Pass || !result.HardEligible {
		t.Fatalf("predicate result %#v", result)
	}
	score, err := scenario.Score(ctx, sess, result)
	if err != nil {
		t.Fatal(err)
	}
	if score != 100 {
		t.Fatalf("score=%v want 100", score)
	}

	if err := adapter.StopSession(ctx, sid); err != nil {
		t.Fatal(err)
	}
	select {
	case <-adapter.Done():
	case <-time.After(5 * time.Second):
		t.Fatal("adapter did not exit after StopSession")
	}
}
