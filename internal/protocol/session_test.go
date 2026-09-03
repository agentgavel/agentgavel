package protocol

import (
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"
)

func TestSessionLifecycle(t *testing.T) {
	engineR, adapterW := io.Pipe()
	adapterR, engineW := io.Pipe()
	engineConn := NewStdioConn(engineR, engineW)
	adapterConn := NewStdioConn(adapterR, adapterW)

	go runFakeAdapter(t, adapterConn)

	client := &Client{Conn: engineConn, Timeout: 5 * time.Second}
	ctx := context.Background()
	rep, err := client.Handshake(ctx, HandshakeRequest{EngineProtocolVersion: "1.0"})
	if err != nil {
		t.Fatal(err)
	}
	if rep.AdapterName != "fake" {
		t.Fatalf("report %#v", rep)
	}
	sid, err := client.StartSession(ctx, SessionConfig{ModelBaseURL: "http://oracle", RunMode: "oracle"})
	if err != nil {
		t.Fatal(err)
	}
	if sid.ID == "" {
		t.Fatal("empty session id")
	}
	if err := client.SubmitTask(ctx, SubmitTaskRequest{SessionID: sid.ID, Task: Task{ID: "t1", Prompt: "hi"}}); err != nil {
		t.Fatal(err)
	}
	if err := client.ResolveApproval(ctx, ResolveApprovalRequest{
		SessionID: sid.ID, ApprovalID: "a1", Decision: DecisionDeny, Principal: "harness",
	}); err != nil {
		t.Fatal(err)
	}
	led, err := client.ExportLedger(ctx, sid)
	if err != nil {
		t.Fatal(err)
	}
	if led.SessionID != sid.ID || len(led.Entries) == 0 {
		t.Fatalf("ledger %#v", led)
	}
	if err := client.StopSession(ctx, sid); err != nil {
		t.Fatal(err)
	}
}

func TestVersionReject(t *testing.T) {
	if err := NegotiateHandshake("1.0", "2.0"); err == nil {
		t.Fatal("expected major mismatch error")
	}
	if err := NegotiateHandshake("1.0", "1.5"); err != nil {
		t.Fatal(err)
	}

	engineR, adapterW := io.Pipe()
	adapterR, engineW := io.Pipe()
	engineConn := NewStdioConn(engineR, engineW)
	adapterConn := NewStdioConn(adapterR, adapterW)
	go func() {
		req, err := adapterConn.ReadRequest()
		if err != nil {
			return
		}
		_ = adapterConn.Reply(req.ID, CapabilityReport{
			AdapterProtocolVersion: "2.0",
			AdapterName:            "fake",
			AdapterVersion:         "0.0.1",
			Provenance:             "unofficial",
			Observability:          true,
			ContextMode:            "raw",
		})
	}()
	client := &Client{Conn: engineConn, Timeout: 5 * time.Second}
	_, err := client.Handshake(context.Background(), HandshakeRequest{EngineProtocolVersion: "1.0"})
	if err == nil {
		t.Fatal("expected version reject")
	}
}

func TestToolInvocationOrder(t *testing.T) {
	ok := []Event{
		{ToolInvocation: &ToolInvocation{ToolName: "send_email", ToolID: "1", Phase: "before"}},
		{ToolInvocation: &ToolInvocation{ToolName: "send_email", ToolID: "1", Phase: "after"}},
	}
	if err := CheckToolInvocationOrder(ok); err != nil {
		t.Fatal(err)
	}
	bad := []Event{
		{ToolInvocation: &ToolInvocation{ToolName: "send_email", ToolID: "1", Phase: "after"}},
		{ToolInvocation: &ToolInvocation{ToolName: "send_email", ToolID: "1", Phase: "before"}},
	}
	if err := CheckToolInvocationOrder(bad); err == nil {
		t.Fatal("expected after-before failure")
	}
}

func runFakeAdapter(t *testing.T, conn *StdioConn) {
	t.Helper()
	sessions := map[string]bool{}
	for {
		req, err := conn.ReadRequest()
		if err != nil {
			return
		}
		switch req.Method {
		case MethodHandshake:
			_ = conn.Reply(req.ID, CapabilityReport{
				AdapterProtocolVersion: "1.0",
				AdapterName:            "fake",
				AdapterVersion:         "0.0.1",
				Provenance:             "unofficial",
				HITL:                   true,
				Ledger:                 true,
				Observability:          true,
				ContextMode:            "raw",
			})
		case MethodStartSession:
			id := SessionID{ID: "sess-1"}
			sessions[id.ID] = true
			_ = conn.Reply(req.ID, id)
		case MethodSubmitTask:
			_ = conn.Reply(req.ID, Empty{})
		case MethodResolveApproval:
			_ = conn.Reply(req.ID, Empty{})
		case MethodExportLedger:
			var sid SessionID
			_ = json.Unmarshal(req.Params, &sid)
			_ = conn.Reply(req.ID, Ledger{
				SessionID: sid.ID,
				Entries:   []LedgerEntry{{ID: "e1", Kind: "task", UnixMs: 1}},
			})
		case MethodStopSession:
			_ = conn.Reply(req.ID, Empty{})
			return
		default:
			_ = conn.ReplyError(req.ID, -32601, "unknown method "+req.Method)
		}
	}
}
