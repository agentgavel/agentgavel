package protocol

import (
	"testing"
)

func TestCodec(t *testing.T) {
	t.Run("HandshakeRequest", func(t *testing.T) {
		in := HandshakeRequest{
			EngineProtocolVersion: "1.0",
			EngineVersion:         "0.0.0-dev",
		}
		out, err := RoundTripJSON(in)
		if err != nil {
			t.Fatal(err)
		}
		if out != in {
			t.Fatalf("got %#v want %#v", out, in)
		}
	})
	t.Run("EventToolInvocation", func(t *testing.T) {
		in := Event{
			SessionID: "s1",
			Seq:       1,
			UnixMs:    1000,
			ToolInvocation: &ToolInvocation{
				ToolName:      "send_email",
				Phase:         "before",
				ArgumentsJSON: `{"to":"attacker@evil.example"}`,
			},
		}
		out, err := RoundTripJSON(in)
		if err != nil {
			t.Fatal(err)
		}
		if out.Kind() != "tool_invocation" {
			t.Fatalf("kind %q", out.Kind())
		}
		if out.ToolInvocation == nil || out.ToolInvocation.ToolName != "send_email" {
			t.Fatalf("got %#v", out.ToolInvocation)
		}
		if out.ToolInvocation.ArgumentsJSON != in.ToolInvocation.ArgumentsJSON {
			t.Fatalf("args %q", out.ToolInvocation.ArgumentsJSON)
		}
	})
}
