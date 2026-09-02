package mcpfuzz

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

func TestSkeleton(t *testing.T) {
	inR, inW := io.Pipe()
	var out bytes.Buffer
	srv := NewSkeleton(inR, &out)

	done := make(chan error, 1)
	go func() { done <- srv.Serve() }()

	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/list",
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := inW.Write(append(b, '\n')); err != nil {
		t.Fatal(err)
	}
	_ = inW.Close()

	if err := <-done; err != nil {
		t.Fatalf("Serve: %v", err)
	}

	var resp struct {
		Result struct {
			Tools []Tool `json:"tools"`
		} `json:"result"`
		Error *rpcError `json:"error"`
	}
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		t.Fatalf("decode %q: %v", out.String(), err)
	}
	if resp.Error != nil {
		t.Fatalf("rpc error: %+v", resp.Error)
	}
	if len(resp.Result.Tools) != 1 || resp.Result.Tools[0].Name != "noop" {
		t.Fatalf("want one noop tool, got %+v", resp.Result.Tools)
	}
}
