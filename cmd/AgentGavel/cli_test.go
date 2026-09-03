package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCLI_Help(t *testing.T) {
	bin := buildAgentGavel(t)
	cmd := exec.Command(bin, "help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("help: %v\n%s", err, out)
	}
	text := string(out)
	for _, n := range []string{"Usage:", "Commands:", "version", "run", "report", "help"} {
		if !strings.Contains(text, n) {
			t.Fatalf("missing %q in:\n%s", n, text)
		}
	}
}

func TestCLI_Version(t *testing.T) {
	bin := buildAgentGavel(t)
	cmd := exec.Command(bin, "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version: %v\n%s", err, out)
	}
	text := strings.TrimSpace(string(out))
	if text == "" {
		t.Fatal("expected non-empty version")
	}
	if !strings.Contains(text, "0.0.0-dev") && !strings.Contains(text, ".") {
		t.Fatalf("unexpected version output: %q", text)
	}
}

func TestCLI_RunMissingAdapter(t *testing.T) {
	bin := buildAgentGavel(t)
	cmd := exec.Command(bin, "run", "--suite", "security", "--mode", "oracle")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit, got: %s", out)
	}
	text := string(out)
	if !strings.Contains(text, "adapter") {
		t.Fatalf("expected error mentioning adapter, got: %s", text)
	}
}
