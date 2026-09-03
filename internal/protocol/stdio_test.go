package protocol

import (
	"encoding/json"
	"io"
	"testing"
)

func TestStdio(t *testing.T) {
	engineR, adapterW := io.Pipe()
	adapterR, engineW := io.Pipe()

	engine := NewStdioConn(engineR, engineW)
	adapter := NewStdioConn(adapterR, adapterW)

	done := make(chan error, 1)
	go func() {
		req, err := adapter.ReadRequest()
		if err != nil {
			done <- err
			return
		}
		if req.Method != MethodHandshake {
			done <- errUnexpected(req.Method)
			return
		}
		var hs HandshakeRequest
		if err := json.Unmarshal(req.Params, &hs); err != nil {
			done <- err
			return
		}
		rep := CapabilityReport{
			AdapterProtocolVersion: hs.EngineProtocolVersion,
			AdapterName:            "fake",
			AdapterVersion:         "0.0.1",
			Provenance:             "unofficial",
			HITL:                   true,
			Observability:          true,
			ContextMode:            "raw",
		}
		done <- adapter.Reply(req.ID, rep)
	}()

	rep, err := engine.Handshake(HandshakeRequest{
		EngineProtocolVersion: "1.0",
		EngineVersion:         "0.0.0-dev",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if rep.AdapterName != "fake" || rep.AdapterProtocolVersion != "1.0" {
		t.Fatalf("report %#v", rep)
	}
	_ = engineR.Close()
	_ = adapterR.Close()
}

func errUnexpected(method string) error {
	return &unexpectedMethodError{method: method}
}

type unexpectedMethodError struct{ method string }

func (e *unexpectedMethodError) Error() string { return "unexpected method " + e.method }
