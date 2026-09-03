// Command fakeadapter is a minimal JSON-RPC stdio adapter for engine launcher tests.
package main

import (
	"encoding/json"
	"os"

	"github.com/agentgavel/gavel/internal/protocol"
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
				Ledger:                 true,
				Observability:          true,
				ContextMode:            "raw",
			})
		case protocol.MethodStartSession:
			_ = conn.Reply(req.ID, protocol.SessionID{ID: "sess-1"})
		case protocol.MethodSubmitTask, protocol.MethodResolveApproval, protocol.MethodExportLedger:
			if req.Method == protocol.MethodExportLedger {
				var sid protocol.SessionID
				_ = json.Unmarshal(req.Params, &sid)
				_ = conn.Reply(req.ID, protocol.Ledger{
					SessionID: sid.ID,
					Entries:   []protocol.LedgerEntry{{ID: "e1", Kind: "task", UnixMs: 1}},
				})
				continue
			}
			_ = conn.Reply(req.ID, protocol.Empty{})
		case protocol.MethodStopSession:
			_ = conn.Reply(req.ID, protocol.Empty{})
			return
		default:
			_ = conn.ReplyError(req.ID, -32601, "unknown method "+req.Method)
		}
	}
}
