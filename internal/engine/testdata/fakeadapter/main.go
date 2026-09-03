// Command fakeadapter is a minimal JSON-RPC stdio adapter for engine launcher tests.
package main

import (
	"encoding/json"
	"os"

	"github.com/agentgavel/agentgavel/internal/protocol"
)

func main() {
	conn := protocol.NewStdioConn(os.Stdin, os.Stdout)
	for {
		req, err := conn.ReadRequest()
		if err != nil {
			return
		}
		switch req.Method {
		case protocol.MethodHandshake:
			_ = conn.Reply(req.ID, protocol.CapabilityReport{
				AdapterProtocolVersion: "1.0",
				AdapterName:            "fake",
				AdapterVersion:         "0.0.1",
				Provenance:             "unofficial",
				HITL:                   true,
				Tenancy:                true,
				Ledger:                 true,
				Observability:          true,
				ContextMode:            "raw",
			})
		case protocol.MethodStartSession:
			_ = conn.Reply(req.ID, protocol.SessionID{ID: "sess-1"})
		case protocol.MethodSubmitTask:
			var st protocol.SubmitTaskRequest
			_ = json.Unmarshal(req.Params, &st)
			// Reply first so Call returns; Event notifies stay buffered for DrainEvents.
			_ = conn.Reply(req.ID, protocol.Empty{})
			sid := st.SessionID
			if sid == "" {
				sid = "sess-1"
			}
			_ = conn.Notify(protocol.MethodEventNotify, protocol.Event{
				SessionID: sid,
				Seq:       1,
				UnixMs:    1,
				ToolInvocation: &protocol.ToolInvocation{
					ToolName: "noop",
					ToolID:   "1",
					Phase:    "before",
				},
			})
			_ = conn.Notify(protocol.MethodEventNotify, protocol.Event{
				SessionID: sid,
				Seq:       2,
				UnixMs:    2,
				ToolInvocation: &protocol.ToolInvocation{
					ToolName: "noop",
					ToolID:   "1",
					Phase:    "after",
					Outcome:  "ok",
				},
			})
		case protocol.MethodResolveApproval:
			_ = conn.Reply(req.ID, protocol.Empty{})
		case protocol.MethodExportLedger:
			var sid protocol.SessionID
			_ = json.Unmarshal(req.Params, &sid)
			_ = conn.Reply(req.ID, protocol.Ledger{
				SessionID: sid.ID,
				Entries:   []protocol.LedgerEntry{{ID: "e1", Kind: "task", UnixMs: 1}},
			})
		case protocol.MethodStopSession:
			_ = conn.Reply(req.ID, protocol.Empty{})
			return
		default:
			_ = conn.ReplyError(req.ID, -32601, "unknown method "+req.Method)
		}
	}
}
