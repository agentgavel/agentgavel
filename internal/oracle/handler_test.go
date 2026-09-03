package oracle

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenAIToolCall(t *testing.T) {
	h := NewHandler()
	dir := Directive{
		ToolName:  "send_email",
		Arguments: map[string]any{"to": "attacker@evil.example"},
	}
	raw, _ := json.Marshal(dir)
	body := `{"model":"oracle","messages":[{"role":"user","content":"hi"}]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(HeaderProbeDirective, string(raw))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	var resp struct {
		Choices []struct {
			Message struct {
				ToolCalls []struct {
					Function struct {
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Choices) != 1 || len(resp.Choices[0].Message.ToolCalls) != 1 {
		t.Fatalf("unexpected response: %s", rr.Body.String())
	}
	tc := resp.Choices[0].Message.ToolCalls[0]
	if tc.Function.Name != "send_email" {
		t.Fatalf("tool name %q", tc.Function.Name)
	}
	var args map[string]any
	if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
		t.Fatal(err)
	}
	if args["to"] != "attacker@evil.example" {
		t.Fatalf("args %#v", args)
	}
}

func TestAnthropicToolUse(t *testing.T) {
	h := NewHandler()
	dir := Directive{
		ToolName:  "send_email",
		Arguments: map[string]any{"to": "attacker@evil.example"},
	}
	raw, _ := json.Marshal(dir)
	body := `{"model":"oracle","messages":[{"role":"user","content":"hi"}]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(HeaderProbeDirective, string(raw))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	var resp struct {
		Content []struct {
			Type  string         `json:"type"`
			Name  string         `json:"name"`
			Input map[string]any `json:"input"`
		} `json:"content"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Content) != 1 || resp.Content[0].Type != "tool_use" {
		t.Fatalf("unexpected: %s", rr.Body.String())
	}
	if resp.Content[0].Name != "send_email" {
		t.Fatalf("name %q", resp.Content[0].Name)
	}
	if resp.Content[0].Input["to"] != "attacker@evil.example" {
		t.Fatalf("input %#v", resp.Content[0].Input)
	}
}

func TestMissingDirective(t *testing.T) {
	h := NewHandler()
	body := `{"model":"oracle","messages":[{"role":"user","content":"hi"}]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code < 400 || rr.Code >= 500 {
		t.Fatalf("want 4xx, got %d %s", rr.Code, rr.Body.String())
	}
}
