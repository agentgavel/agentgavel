package protocol_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestAdapterProtoDefinesSevenRPCs(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	// This test file will live under internal/protocol once codec lands;
	// for T2.1 we keep the check next to proto via relative path from module root.
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))
	protoPath := filepath.Join(root, "proto", "adapter.proto")
	b, err := os.ReadFile(protoPath)
	if err != nil {
		// Fallback: cwd is module root during `go test ./internal/protocol`
		b, err = os.ReadFile(filepath.Join("..", "..", "proto", "adapter.proto"))
		if err != nil {
			b, err = os.ReadFile("proto/adapter.proto")
		}
	}
	if err != nil {
		t.Fatalf("read proto: %v", err)
	}
	body := string(b)
	for _, rpc := range []string{
		"rpc Handshake(",
		"rpc StartSession(",
		"rpc SubmitTask(",
		"rpc ResolveApproval(",
		"rpc Events(",
		"rpc ExportLedger(",
		"rpc StopSession(",
	} {
		if !strings.Contains(body, rpc) {
			t.Errorf("missing %s", rpc)
		}
	}
	for _, kind := range []string{
		"ToolInvocation tool_invocation",
		"GateDecision gate_decision",
		"ContextSnapshot context_snapshot",
		"ContextAttestation context_attestation",
		"LedgerAppend ledger_append",
		"SessionError session_error",
	} {
		if !strings.Contains(body, kind) {
			t.Errorf("missing event kind %s", kind)
		}
	}
}
