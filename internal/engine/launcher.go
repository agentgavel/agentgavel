package engine

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/agentgavel/agentgavel/internal/protocol"
)

// LaunchConfig describes how to exec an adapter sidecar.
type LaunchConfig struct {
	Command string
	Args    []string
	Dir     string
	Env     []string
	// Timeout is the default per-RPC timeout on the protocol.Client.
	Timeout time.Duration
}

// AdapterSession is a running sidecar process wired to a protocol.Client over stdio.
type AdapterSession struct {
	Client *protocol.Client

	cmd    *exec.Cmd
	cancel context.CancelFunc
	stdin  io.WriteCloser

	waitOnce sync.Once
	waitErr  error
	waitDone chan struct{}

	closeOnce sync.Once
}

// Launch starts an adapter process, wires stdout→engine and engine→stdin via
// protocol.NewStdioConn, and kills the child when parent is canceled.
func Launch(parent context.Context, cfg LaunchConfig) (*AdapterSession, error) {
	if cfg.Command == "" {
		return nil, fmt.Errorf("engine: empty adapter command")
	}
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	cmd := exec.CommandContext(ctx, cfg.Command, cfg.Args...)
	cmd.Dir = cfg.Dir
	if len(cfg.Env) > 0 {
		cmd.Env = append(os.Environ(), cfg.Env...)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("engine: stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		_ = stdin.Close()
		cancel()
		return nil, fmt.Errorf("engine: stdout pipe: %w", err)
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		_ = stdin.Close()
		cancel()
		return nil, fmt.Errorf("engine: start adapter: %w", err)
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	conn := protocol.NewStdioConn(stdout, stdin)
	s := &AdapterSession{
		Client:   &protocol.Client{Conn: conn, Timeout: timeout},
		cmd:      cmd,
		cancel:   cancel,
		stdin:    stdin,
		waitDone: make(chan struct{}),
	}
	go s.reap()
	return s, nil
}

func (s *AdapterSession) reap() {
	err := s.cmd.Wait()
	s.waitOnce.Do(func() {
		s.waitErr = err
		close(s.waitDone)
	})
}

// Handshake negotiates the adapter capability report.
func (s *AdapterSession) Handshake(ctx context.Context, req protocol.HandshakeRequest) (protocol.CapabilityReport, error) {
	return s.Client.Handshake(ctx, req)
}

// StopSession asks the adapter to tear down, then waits for the process to exit.
// If the child does not exit promptly, it is killed.
func (s *AdapterSession) StopSession(ctx context.Context, id protocol.SessionID) error {
	rpcErr := s.Client.StopSession(ctx, id)
	_ = s.stdin.Close()

	waitCtx, waitCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer waitCancel()

	select {
	case <-s.waitDone:
		return firstErr(rpcErr, s.waitErr)
	case <-waitCtx.Done():
		s.kill()
		<-s.waitDone
		return firstErr(rpcErr, fmt.Errorf("engine: adapter did not exit after StopSession: %w", s.waitErr))
	}
}

// Close cancels the launch context, kills the process if still running, and waits.
func (s *AdapterSession) Close() error {
	var err error
	s.closeOnce.Do(func() {
		_ = s.stdin.Close()
		s.kill()
		<-s.waitDone
		err = s.waitErr
		if err != nil && isExitOK(err) {
			err = nil
		}
	})
	return err
}

// Wait blocks until the adapter process exits and returns the wait error.
func (s *AdapterSession) Wait() error {
	<-s.waitDone
	return s.waitErr
}

// Done is closed when the child process has been waited on.
func (s *AdapterSession) Done() <-chan struct{} {
	return s.waitDone
}

// ProcessState returns the exited process state, or nil if still running.
func (s *AdapterSession) ProcessState() *os.ProcessState {
	return s.cmd.ProcessState
}

// Process returns the underlying os.Process, if started.
func (s *AdapterSession) Process() *os.Process {
	if s.cmd == nil {
		return nil
	}
	return s.cmd.Process
}

func (s *AdapterSession) kill() {
	s.cancel()
	if p := s.Process(); p != nil {
		_ = p.Kill()
	}
}

func firstErr(a, b error) error {
	if a != nil {
		return a
	}
	if b != nil && !isExitOK(b) {
		return b
	}
	return nil
}

func isExitOK(err error) bool {
	if err == nil {
		return true
	}
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode() == 0
	}
	return false
}
