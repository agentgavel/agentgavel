package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestReportPublishWritesEntryAndIndex(t *testing.T) {
	bin := buildAgentGavel(t)
	root, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	dash := t.TempDir()
	cmd := exec.Command(bin, "report",
		"--root", root,
		"--publish",
		"--dashboard", dash,
		"--framework", "FakeAdapter",
		"--adapter-name", "fakeadapter",
		"sample-catastrophic",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("report --publish: %v\n%s", err, out)
	}
	var written string
	if err := json.Unmarshal(out, &written); err != nil {
		t.Fatalf("stdout not JSON path: %v\n%s", err, out)
	}
	entryPath := filepath.Join(dash, "data", "sample-catastrophic.json")
	absEntry, _ := filepath.Abs(entryPath)
	if written != absEntry {
		t.Fatalf("written path = %q, want %q", written, absEntry)
	}

	raw, err := os.ReadFile(entryPath)
	if err != nil {
		t.Fatal(err)
	}
	var entry map[string]any
	if err := json.Unmarshal(raw, &entry); err != nil {
		t.Fatal(err)
	}
	if entry["tab"] != "unratified" {
		t.Fatalf("tab = %v, want unratified", entry["tab"])
	}
	if entry["sample"] != false {
		t.Fatalf("sample = %v, want false", entry["sample"])
	}
	if entry["provenance"] != "unofficial" {
		t.Fatalf("provenance = %v, want unofficial (from summary.json)", entry["provenance"])
	}
	if entry["grade"] != "F" {
		t.Fatalf("grade = %v, want F", entry["grade"])
	}
	if _, ok := entry["gsi"].(float64); !ok {
		t.Fatalf("gsi missing or wrong type: %v", entry["gsi"])
	}
	pillars, ok := entry["pillars"].(map[string]any)
	if !ok || len(pillars) == 0 {
		t.Fatalf("pillars = %v", entry["pillars"])
	}
	na, ok := entry["na"].([]any)
	if !ok || len(na) == 0 {
		t.Fatalf("na = %v, want SEC-008 present", entry["na"])
	}
	cat, ok := entry["catastrophic"].([]any)
	if !ok || len(cat) == 0 {
		t.Fatalf("catastrophic = %v, want SEC-004", entry["catastrophic"])
	}
	fp, ok := entry["fingerprint"].(map[string]any)
	if !ok || fp["adapter.version"] != "0.1.0" {
		t.Fatalf("fingerprint = %v", entry["fingerprint"])
	}
	if entry["framework"] != "FakeAdapter" || entry["adapter"] != "fakeadapter" {
		t.Fatalf("framework/adapter = %v/%v", entry["framework"], entry["adapter"])
	}
	if entry["adapter_version"] != "0.1.0" {
		t.Fatalf("adapter_version = %v, want 0.1.0", entry["adapter_version"])
	}

	idxRaw, err := os.ReadFile(filepath.Join(dash, "data", "index.json"))
	if err != nil {
		t.Fatal(err)
	}
	var idx []string
	if err := json.Unmarshal(idxRaw, &idx); err != nil {
		t.Fatal(err)
	}
	if len(idx) != 1 || idx[0] != "sample-catastrophic.json" {
		t.Fatalf("index = %#v, want [sample-catastrophic.json] exactly once", idx)
	}
}

func TestReportPublishIdempotentIndex(t *testing.T) {
	bin := buildAgentGavel(t)
	root, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	dash := t.TempDir()
	args := []string{
		"report",
		"--root", root,
		"--publish",
		"--dashboard", dash,
		"--framework", "FakeAdapter",
		"--adapter-name", "fakeadapter",
		"sample-catastrophic",
	}
	for i := 0; i < 2; i++ {
		cmd := exec.Command(bin, args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("publish #%d: %v\n%s", i+1, err, out)
		}
	}
	idxRaw, err := os.ReadFile(filepath.Join(dash, "data", "index.json"))
	if err != nil {
		t.Fatal(err)
	}
	var idx []string
	if err := json.Unmarshal(idxRaw, &idx); err != nil {
		t.Fatal(err)
	}
	if len(idx) != 1 || idx[0] != "sample-catastrophic.json" {
		t.Fatalf("index after 2 publishes = %#v, want filename exactly once", idx)
	}
}

func TestReportPublishRejectsOptIn(t *testing.T) {
	bin := buildAgentGavel(t)
	root, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(bin, "report",
		"--root", root,
		"--publish",
		"--dashboard", t.TempDir(),
		"--framework", "FakeAdapter",
		"--adapter-name", "fakeadapter",
		"--tab", "opt-in",
		"sample-catastrophic",
	)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected exit 2, got success:\n%s", out)
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("err type %T, want ExitError; out=%s", err, out)
	}
	if exitErr.ExitCode() != 2 {
		t.Fatalf("exit = %d, want 2\n%s", exitErr.ExitCode(), out)
	}
	text := string(out)
	if !strings.Contains(text, "ADR 006") {
		t.Fatalf("stderr missing ADR 006 citation:\n%s", text)
	}
	if !strings.Contains(text, "opt-in") {
		t.Fatalf("stderr missing opt-in:\n%s", text)
	}
}
