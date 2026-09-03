package engine

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/agentgavel/gavel/internal/protocol"
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
		Env:     append([]string{}, "PYTHONPATH=src"),
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
