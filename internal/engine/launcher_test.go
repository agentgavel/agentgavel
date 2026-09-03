package engine

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/agentgavel/gavel/internal/protocol"
)

func buildFakeAdapter(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "fakeadapter")
	cmd := exec.Command("go", "build", "-o", bin, "./testdata/fakeadapter")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build fakeadapter: %v\n%s", err, out)
	}
	return bin
}

func TestLaunchFakeAdapter(t *testing.T) {
	bin := buildFakeAdapter(t)
	ctx := context.Background()
	sess, err := Launch(ctx, LaunchConfig{Command: bin, Timeout: 5 * time.Second})
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
	case <-time.After(5 * time.Second):
		t.Fatal("adapter process did not exit after StopSession")
	}
	ps := sess.ProcessState()
	if ps == nil {
		t.Fatal("expected ProcessState after exit")
	}
	if !ps.Exited() {
		t.Fatalf("process state not exited: %#v", ps)
	}
}

func TestLaunchFakeAdapter_ContextCancel(t *testing.T) {
	bin := buildFakeAdapter(t)
	ctx, cancel := context.WithCancel(context.Background())
	sess, err := Launch(ctx, LaunchConfig{Command: bin, Timeout: 5 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sess.Close() }()

	if _, err := sess.Handshake(context.Background(), protocol.HandshakeRequest{EngineProtocolVersion: "1.0"}); err != nil {
		t.Fatal(err)
	}

	cancel()

	select {
	case <-sess.Done():
	case <-time.After(5 * time.Second):
		t.Fatal("adapter process did not exit after context cancel")
	}
	ps := sess.ProcessState()
	if ps == nil {
		t.Fatal("expected ProcessState after context cancel")
	}
	// CommandContext kills with a signal; on Unix Exited() is false for signals.
	if ps.Success() {
		t.Fatal("expected non-zero / signal termination after context cancel")
	}
}
