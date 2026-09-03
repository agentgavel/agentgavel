package engine

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/agentgavel/agentgavel/internal/protocol"
)

// pythonSDKDir resolves the sdk/python directory relative to this package,
// or skips the test if python3 or the SDK source tree isn't available.
func pythonSDKDir(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("python3 not available")
	}
	dir, err := filepath.Abs(filepath.Join("..", "..", "sdk", "python"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := exec.Command("python3", "-c", "import sys").CombinedOutput(); err != nil {
		t.Skipf("python3 unusable: %v", err)
	}
	return dir
}

// TestHandshakePythonFakeAdapterExample launches the fake adapter example
// shipped with the Python SDK and drives a Handshake over stdio JSON-RPC,
// proving the Go engine launcher interoperates with the Python transport.
func TestHandshakePythonFakeAdapterExample(t *testing.T) {
	sdkDir := pythonSDKDir(t)

	ctx := context.Background()
	sess, err := Launch(ctx, LaunchConfig{
		Command: "python3",
		Args:    []string{"-m", "agentgavel_adapter.examples.fake"},
		Dir:     sdkDir,
		Env:     []string{"PYTHONPATH=src"},
		Timeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sess.Close() }()

	rep, err := sess.Handshake(ctx, protocol.HandshakeRequest{EngineProtocolVersion: "1.0"})
	if err != nil {
		t.Fatal(err)
	}
	if rep.AdapterName != "fake" || rep.AdapterProtocolVersion != "1.0" {
		t.Fatalf("capability report %#v", rep)
	}

	if err := sess.StopSession(ctx, protocol.SessionID{ID: "sess-1"}); err != nil {
		t.Fatal(err)
	}

	select {
	case <-sess.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("python fake adapter did not exit after StopSession")
	}
}

// TestPythonFakeAdapter launches the Python FakeAdapter example and runs the
// same oracle-mode noop flow as TestIntegrationNoopScenario against the Go
// engine (UC-001, UC-006).
func TestPythonFakeAdapter(t *testing.T) {
	sdkDir := pythonSDKDir(t)

	ctx := context.Background()
	adapter, err := Launch(ctx, LaunchConfig{
		Command: "python3",
		Args:    []string{"-m", "agentgavel_adapter.examples.fake"},
		Dir:     sdkDir,
		Env:     []string{"PYTHONPATH=src"},
		Timeout: 10 * time.Second,
	})
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
	case <-time.After(10 * time.Second):
		t.Fatal("python fake adapter did not exit after StopSession")
	}
}
