package mcpfuzz

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sync/atomic"
)

// Server is a minimal MCP JSON-RPC server over stdio (newline-delimited).
// Modes embed or wrap Server to advertise misbehavior for SEC-007.
type Server struct {
	Name    string
	Version string
	Tools   []Tool
	in      io.Reader
	out     io.Writer
	// OnCall optionally handles tools/call; default returns a fixed noop result.
	OnCall func(name string, args map[string]any) (any, error)
	// RenameCallAs, when set (tool-renamer mode), is the name that differs from Tools[0].Name.
	RenameCallAs string
	// BackendID is an alternate identity for masquerade mode.
	BackendID string
	closed    atomic.Bool
}

// Tool is an MCP tool descriptor subset.
type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	InputSchema map[string]any `json:"inputSchema"`
}

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewSkeleton returns a server that lists a single noop tool.
func NewSkeleton(in io.Reader, out io.Writer) *Server {
	return &Server{
		Name:    "agentgavel-mcpfuzz-skeleton",
		Version: "0.0.0-dev",
		in:      in,
		out:     out,
		Tools: []Tool{{
			Name:        "noop",
			Description: "No-op tool for skeleton tests",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		}},
		OnCall: func(name string, args map[string]any) (any, error) {
			return map[string]any{
				"content": []map[string]any{{
					"type": "text",
					"text": "noop ok",
				}},
			}, nil
		},
	}
}

// Serve reads JSON-RPC requests until EOF or Close.
func (s *Server) Serve() error {
	sc := bufio.NewScanner(s.in)
	// MCP tool schemas can be large in later modes; allow bigger lines.
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 1024*1024)
	for sc.Scan() {
		if s.closed.Load() {
			return nil
		}
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var req rpcRequest
		if err := json.Unmarshal(line, &req); err != nil {
			_ = s.write(rpcResponse{
				JSONRPC: "2.0",
				Error:   &rpcError{Code: -32700, Message: "parse error"},
			})
			continue
		}
		resp := s.handle(req)
		if resp == nil {
			continue // notification
		}
		if err := s.write(*resp); err != nil {
			return err
		}
	}
	return sc.Err()
}

// Close signals Serve to stop after the current request.
func (s *Server) Close() {
	s.closed.Store(true)
}

func (s *Server) handle(req rpcRequest) *rpcResponse {
	switch req.Method {
	case "initialize":
		return &rpcResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]any{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]any{"tools": map[string]any{}},
				"serverInfo": map[string]any{
					"name":    s.Name,
					"version": s.Version,
				},
			},
		}
	case "notifications/initialized":
		return nil
	case "tools/list":
		return &rpcResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  map[string]any{"tools": s.Tools},
		}
	case "tools/call":
		var p struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &p); err != nil {
			return &rpcResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &rpcError{Code: -32602, Message: "invalid params"},
			}
		}
		fn := s.OnCall
		if fn == nil {
			fn = func(string, map[string]any) (any, error) {
				return map[string]any{"content": []map[string]any{{"type": "text", "text": "ok"}}}, nil
			}
		}
		res, err := fn(p.Name, p.Arguments)
		if err != nil {
			return &rpcResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &rpcError{Code: -32000, Message: err.Error()},
			}
		}
		return &rpcResponse{JSONRPC: "2.0", ID: req.ID, Result: res}
	default:
		if len(req.ID) == 0 {
			return nil
		}
		return &rpcResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &rpcError{Code: -32601, Message: fmt.Sprintf("method not found: %s", req.Method)},
		}
	}
}

func (s *Server) write(resp rpcResponse) error {
	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = s.out.Write(append(b, '\n'))
	return err
}
