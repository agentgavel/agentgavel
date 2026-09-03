package engine

import (
	"encoding/json"
	"fmt"

	"github.com/agentgavel/gavel/internal/protocol"
)

// DrainEvents reads n adapter→engine Event notifications from conn.
// Call this after RPCs that may have interleaved notify lines already
// consumed by Call; events emitted after an RPC response remain buffered
// on the pipe for this drain.
func DrainEvents(conn *protocol.StdioConn, n int) ([]protocol.Event, error) {
	if conn == nil {
		return nil, fmt.Errorf("engine: nil conn")
	}
	if n < 0 {
		return nil, fmt.Errorf("engine: negative event count")
	}
	events := make([]protocol.Event, 0, n)
	for len(events) < n {
		req, err := conn.ReadRequest()
		if err != nil {
			return events, fmt.Errorf("engine: drain events: %w", err)
		}
		if req.Method != protocol.MethodEventNotify {
			return events, fmt.Errorf("engine: expected %s notification, got %q", protocol.MethodEventNotify, req.Method)
		}
		var ev protocol.Event
		if err := json.Unmarshal(req.Params, &ev); err != nil {
			return events, fmt.Errorf("engine: decode event: %w", err)
		}
		events = append(events, ev)
	}
	return events, nil
}
