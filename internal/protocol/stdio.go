package protocol

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

// JSON-RPC 2.0 method names for the adapter contract.
const (
	MethodHandshake       = "Handshake"
	MethodStartSession    = "StartSession"
	MethodSubmitTask      = "SubmitTask"
	MethodResolveApproval = "ResolveApproval"
	MethodEventsSubscribe = "Events"
	MethodExportLedger    = "ExportLedger"
	MethodStopSession     = "StopSession"
	// MethodEventNotify is an adapter->engine notification (no response).
	MethodEventNotify = "Event"
)

// RPCRequest is a JSON-RPC 2.0 request.
type RPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// RPCResponse is a JSON-RPC 2.0 response.
type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

// RPCError is a JSON-RPC error object.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// StdioConn is a newline-delimited JSON-RPC connection over an io pair.
type StdioConn struct {
	in  *bufio.Scanner
	out io.Writer
	mu  sync.Mutex
	id  int64
}

// NewStdioConn wraps reader/writer for JSON-RPC framing.
func NewStdioConn(in io.Reader, out io.Writer) *StdioConn {
	sc := bufio.NewScanner(in)
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 4*1024*1024)
	return &StdioConn{in: sc, out: out}
}

// Call sends a request and waits for the matching response.
func (c *StdioConn) Call(method string, params any, result any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.id++
	idBytes, _ := json.Marshal(c.id)
	var paramsRaw json.RawMessage
	if params != nil {
		b, err := json.Marshal(params)
		if err != nil {
			return err
		}
		paramsRaw = b
	}
	req := RPCRequest{JSONRPC: "2.0", ID: idBytes, Method: method, Params: paramsRaw}
	if err := c.writeLocked(req); err != nil {
		return err
	}
	for c.in.Scan() {
		line := c.in.Bytes()
		var resp RPCResponse
		if err := json.Unmarshal(line, &resp); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
		if len(resp.ID) == 0 {
			continue // notification
		}
		if string(resp.ID) != string(idBytes) {
			continue
		}
		if resp.Error != nil {
			return fmt.Errorf("rpc %s: %s", method, resp.Error.Message)
		}
		if result != nil && len(resp.Result) > 0 {
			if err := json.Unmarshal(resp.Result, result); err != nil {
				return err
			}
		}
		return nil
	}
	if err := c.in.Err(); err != nil {
		return err
	}
	return io.ErrUnexpectedEOF
}

// Notify sends a JSON-RPC notification (no id).
func (c *StdioConn) Notify(method string, params any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var paramsRaw json.RawMessage
	if params != nil {
		b, err := json.Marshal(params)
		if err != nil {
			return err
		}
		paramsRaw = b
	}
	return c.writeLocked(RPCRequest{JSONRPC: "2.0", Method: method, Params: paramsRaw})
}

// ReadRequest reads the next request from the peer.
func (c *StdioConn) ReadRequest() (*RPCRequest, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.in.Scan() {
		if err := c.in.Err(); err != nil {
			return nil, err
		}
		return nil, io.EOF
	}
	var req RPCRequest
	if err := json.Unmarshal(c.in.Bytes(), &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// Reply writes a success response for id.
func (c *StdioConn) Reply(id json.RawMessage, result any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	b, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return c.writeLocked(RPCResponse{JSONRPC: "2.0", ID: id, Result: b})
}

// ReplyError writes an error response.
func (c *StdioConn) ReplyError(id json.RawMessage, code int, msg string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.writeLocked(RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &RPCError{Code: code, Message: msg},
	})
}

func (c *StdioConn) writeLocked(v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = c.out.Write(append(b, '\n'))
	return err
}

// Handshake is a convenience Call for MethodHandshake.
func (c *StdioConn) Handshake(req HandshakeRequest) (CapabilityReport, error) {
	var rep CapabilityReport
	err := c.Call(MethodHandshake, req, &rep)
	return rep, err
}
