package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func buildAgentGavel(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "AgentGavel")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Env = append(os.Environ(), "GOWORK=off")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	return bin
}

func TestReportPrintsGSIGradePillarsCatastrophic(t *testing.T) {
	bin := buildAgentGavel(t)
	root, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(bin, "report", "--root", root, "sample-catastrophic")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("report: %v\n%s", err, out)
	}
	text := string(out)
	needles := []string{
		"GSI:",
		"Grade: F",
		"Pillars:",
		"chokepoint:",
		"governance:",
		"auditability:",
		"resilience:",
		"Catastrophic flags:",
		"SEC-004",
		"sample-catastrophic",
	}
	for _, n := range needles {
		if !strings.Contains(text, n) {
			t.Fatalf("missing %q in:\n%s", n, text)
		}
	}
}

func TestReportJSON(t *testing.T) {
	bin := buildAgentGavel(t)
	root, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(bin, "report", "--root", root, "--json", "sample-catastrophic")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("report --json: %v\n%s", err, out)
	}
	text := string(out)
	for _, n := range []string{`"grade": "F"`, `"SEC-004"`, `"gsi":`, `"pillars"`} {
		if !strings.Contains(text, n) {
			t.Fatalf("missing %q in:\n%s", n, text)
		}
	}
}

func TestReportScorecardFixture(t *testing.T) {
	bin := buildAgentGavel(t)
	root, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(bin, "report", "--root", root, "sample-clean")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("report: %v\n%s", err, out)
	}
	text := string(out)
	for _, n := range []string{"GSI: 1000.0", "Grade: AAA", "Catastrophic flags:", "(none)", "Pillars:"} {
		if !strings.Contains(text, n) {
			t.Fatalf("missing %q in:\n%s", n, text)
		}
	}
}

func TestReportMissingRun(t *testing.T) {
	bin := buildAgentGavel(t)
	cmd := exec.Command(bin, "report", "--root", t.TempDir(), "does-not-exist")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected error, got: %s", out)
	}
	if !strings.Contains(string(out), "summary.json") && !strings.Contains(string(out), "report:") {
		t.Fatalf("unexpected stderr: %s", out)
	}
}
