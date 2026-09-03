package protocol

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// EngineProtocolVersion is the protocol major.minor this engine speaks.
const EngineProtocolVersion = "1.0"

// Client drives an adapter session over a StdioConn.
type Client struct {
	Conn    *StdioConn
	Timeout time.Duration
}

func (c *Client) timeout() time.Duration {
	if c.Timeout > 0 {
		return c.Timeout
	}
	return 10 * time.Second
}

func (c *Client) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, c.timeout())
}

// NegotiateHandshake validates major-version compatibility and returns the report.
func NegotiateHandshake(engineVer, adapterVer string) error {
	em, _, err := ParseMajorMinor(engineVer)
	if err != nil {
		return fmt.Errorf("engine version: %w", err)
	}
	am, _, err := ParseMajorMinor(adapterVer)
	if err != nil {
		return fmt.Errorf("adapter version: %w", err)
	}
	if em != am {
		return fmt.Errorf("incompatible protocol major: engine=%s adapter=%s", engineVer, adapterVer)
	}
	return nil
}

// ParseMajorMinor parses "MAJOR.MINOR" (optional patch ignored).
func ParseMajorMinor(v string) (major, minor int, err error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, 0, fmt.Errorf("empty version")
	}
	parts := strings.Split(v, ".")
	if len(parts) < 1 {
		return 0, 0, fmt.Errorf("invalid version %q", v)
	}
	major, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	if len(parts) > 1 {
		minor, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, err
		}
	}
	return major, minor, nil
}

// Handshake performs NegotiateHandshake against the adapter.
func (c *Client) Handshake(ctx context.Context, req HandshakeRequest) (CapabilityReport, error) {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	if req.EngineProtocolVersion == "" {
		req.EngineProtocolVersion = EngineProtocolVersion
	}
	type result struct {
		rep CapabilityReport
		err error
	}
	ch := make(chan result, 1)
	go func() {
		rep, err := c.Conn.Handshake(req)
		ch <- result{rep, err}
	}()
	select {
	case <-ctx.Done():
		return CapabilityReport{}, ctx.Err()
	case r := <-ch:
		if r.err != nil {
			return CapabilityReport{}, r.err
		}
		if err := NegotiateHandshake(req.EngineProtocolVersion, r.rep.AdapterProtocolVersion); err != nil {
			return CapabilityReport{}, err
		}
		return r.rep, nil
	}
}

// StartSession starts a session on the adapter.
func (c *Client) StartSession(ctx context.Context, cfg SessionConfig) (SessionID, error) {
	return callTimeout[SessionConfig, SessionID](ctx, c, MethodStartSession, cfg)
}

// SubmitTask submits a task.
func (c *Client) SubmitTask(ctx context.Context, req SubmitTaskRequest) error {
	_, err := callTimeout[SubmitTaskRequest, Empty](ctx, c, MethodSubmitTask, req)
	return err
}

// ResolveApproval resolves an approval gate.
func (c *Client) ResolveApproval(ctx context.Context, req ResolveApprovalRequest) error {
	_, err := callTimeout[ResolveApprovalRequest, Empty](ctx, c, MethodResolveApproval, req)
	return err
}

// ExportLedger exports the session ledger.
func (c *Client) ExportLedger(ctx context.Context, id SessionID) (Ledger, error) {
	return callTimeout[SessionID, Ledger](ctx, c, MethodExportLedger, id)
}

// StopSession stops the session.
func (c *Client) StopSession(ctx context.Context, id SessionID) error {
	_, err := callTimeout[SessionID, Empty](ctx, c, MethodStopSession, id)
	return err
}

func callTimeout[P any, R any](ctx context.Context, c *Client, method string, params P) (R, error) {
	var zero R
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	type result struct {
		v   R
		err error
	}
	ch := make(chan result, 1)
	go func() {
		var out R
		err := c.Conn.Call(method, params, &out)
		ch <- result{out, err}
	}()
	select {
	case <-ctx.Done():
		return zero, ctx.Err()
	case r := <-ch:
		return r.v, r.err
	}
}

// CheckToolInvocationOrder returns an error if an after-phase precedes its before for the same tool call id/name.
func CheckToolInvocationOrder(events []Event) error {
	seenBefore := map[string]bool{}
	for _, e := range events {
		inv := e.ToolInvocation
		if inv == nil {
			continue
		}
		key := inv.ToolID
		if key == "" {
			key = inv.ToolName
		}
		switch inv.Phase {
		case "before":
			seenBefore[key] = true
		case "after":
			if !seenBefore[key] {
				return fmt.Errorf("tool_invocation after precedes before for %q", key)
			}
		}
	}
	return nil
}

// DecodeParams unmarshals raw JSON-RPC params into dst.
func DecodeParams(raw json.RawMessage, dst any) error {
	if len(raw) == 0 {
		return fmt.Errorf("missing params")
	}
	return json.Unmarshal(raw, dst)
}
