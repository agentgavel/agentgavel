package oracle_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/agentgavel/agentgavel/internal/oracle"
)

func TestProbeDirectiveBindingDocumented(t *testing.T) {
	if oracle.HeaderProbeDirective != "X-AgentGavel-Probe-Directive" {
		t.Fatalf("header constant = %q", oracle.HeaderProbeDirective)
	}
	if !strings.HasPrefix(oracle.SystemPrefixProbeDirective, "AGENTGAVEL_PROBE_DIRECTIVE:") {
		t.Fatalf("system prefix = %q", oracle.SystemPrefixProbeDirective)
	}
	readme := filepath.Join("README.md")
	b, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	body := string(b)
	for _, need := range []string{
		"X-AgentGavel-Probe-Directive",
		"tool_name",
		"arguments",
		"AGENTGAVEL_PROBE_DIRECTIVE:",
	} {
		if !strings.Contains(body, need) {
			t.Errorf("README missing %q", need)
		}
	}
}
