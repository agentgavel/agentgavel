package mcpfuzz

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
)

// RunningMode is a named fuzz mode accepting TCP connections with
// newline-delimited JSON-RPC (same framing as stdio Serve).
type RunningMode struct {
	Mode string
	// Addr is host:port suitable for net.Dial("tcp", Addr).
	Addr string
	// URL is a dialable endpoint string (tcp://host:port).
	URL string

	ln       net.Listener
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	errMu    sync.Mutex
	serveErr error
}

// StartMode listens on 127.0.0.1:0 and serves mode on each accepted connection.
func StartMode(parent context.Context, mode string) (*RunningMode, error) {
	if parent == nil {
		parent = context.Background()
	}
	// Validate mode before binding a port.
	if _, err := NewByName(mode, nil, nil); err != nil {
		return nil, err
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("mcpfuzz: listen: %w", err)
	}
	ctx, cancel := context.WithCancel(parent)
	rm := &RunningMode{
		Mode:   mode,
		Addr:   ln.Addr().String(),
		URL:    "tcp://" + ln.Addr().String(),
		ln:     ln,
		cancel: cancel,
	}

	rm.wg.Add(1)
	go func() {
		defer rm.wg.Done()
		<-ctx.Done()
		_ = ln.Close()
	}()

	rm.wg.Add(1)
	go func() {
		defer rm.wg.Done()
		rm.acceptLoop(ctx)
	}()

	return rm, nil
}

func (rm *RunningMode) acceptLoop(ctx context.Context) {
	for {
		conn, err := rm.ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				rm.errMu.Lock()
				rm.serveErr = err
				rm.errMu.Unlock()
				return
			}
		}
		rm.wg.Add(1)
		go func(c net.Conn) {
			defer rm.wg.Done()
			defer c.Close()
			srv, err := NewByName(rm.Mode, c, c)
			if err != nil {
				return
			}
			_ = srv.Serve()
		}(conn)
	}
}

// Close stops accepting and waits for in-flight connections to finish.
func (rm *RunningMode) Close() error {
	if rm == nil {
		return nil
	}
	rm.cancel()
	_ = rm.ln.Close()
	rm.wg.Wait()
	rm.errMu.Lock()
	defer rm.errMu.Unlock()
	if rm.serveErr != nil && !isClosedNetErr(rm.serveErr) {
		return rm.serveErr
	}
	return nil
}

func isClosedNetErr(err error) bool {
	return err == nil || errors.Is(err, net.ErrClosed) || errors.Is(err, context.Canceled)
}
