package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestOracleListenHealth(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "AgentGavel")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Env = append(os.Environ(), "GOWORK=off")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}

	cmd := exec.Command(bin, "oracle", "--listen", "127.0.0.1:0")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	}()

	type result struct {
		n   int
		err error
		buf []byte
	}
	ch := make(chan result, 1)
	go func() {
		buf := make([]byte, 256)
		n, err := stdout.Read(buf)
		ch <- result{n, err, buf}
	}()
	select {
	case r := <-ch:
		if r.err != nil && r.err != io.EOF {
			errOut, _ := io.ReadAll(stderr)
			t.Fatalf("read addr: %v stderr=%s", r.err, errOut)
		}
		addr := strings.TrimSpace(string(r.buf[:r.n]))
		if addr == "" {
			t.Fatal("empty addr")
		}
		resp, err := http.Get("http://" + addr + "/healthz")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("health %d", resp.StatusCode)
		}
	case <-time.After(15 * time.Second):
		errOut, _ := io.ReadAll(stderr)
		t.Fatalf("timeout waiting for listen addr; stderr=%s", errOut)
	}
}
