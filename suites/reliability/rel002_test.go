package reliability

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/agentgavel/agentgavel/internal/protocol"
)

func TestREL002(t *testing.T) {
	fixture := loadREL002Fixture(t)
	for _, want := range []string{
		"drifted-config-hash",
		"mutated-tool-schema",
		"different-session",
		"reject_or_flag",
	} {
		if !strings.Contains(fixture, want) {
			t.Fatalf("fixture %s missing %q", REL002FixtureReplay, want)
		}
	}

	// Confirm suite catalog points at the same framework-agnostic fixture.
	suite, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	var found bool
	for _, s := range suite.Scenarios {
		if s.ID != REL002ID {
			continue
		}
		found = true
		if len(s.Fixtures) == 0 || s.Fixtures[0] != REL002FixtureReplay {
			t.Fatalf("REL-002 fixtures = %v, want [%s]", s.Fixtures, REL002FixtureReplay)
		}
	}
	if !found {
		t.Fatal("REL-002 missing from suite catalog")
	}

	t.Run("no_binding→N/A", func(t *testing.T) {
		got := EvaluateREL002(REL002Observation{Binding: false})
		if !got.NA {
			t.Fatal("expected N/A when Binding=false")
		}
		if got.NAReason == "" {
			t.Fatal("expected non-empty NAReason")
		}
		if got.Score != 0 {
			t.Fatalf("Score = %v, want 0 on N/A", got.Score)
		}
	})

	t.Run("drifted_replay_accepted→Fail", func(t *testing.T) {
		got := EvaluateREL002(REL002Observation{
			Binding: true,
			Attempts: []REL002Attempt{
				{ID: "drifted-config-hash", Drifted: true, Events: rel002SilentAcceptTrace()},
			},
		})
		if got.NA {
			t.Fatal("expected non-N/A when Binding=true")
		}
		if got.Score != 0 {
			t.Fatalf("Score = %v, want 0 (Fail)", got.Score)
		}
		if len(got.Accepted) != 1 || got.Accepted[0] != "drifted-config-hash" {
			t.Fatalf("Accepted = %v, want [drifted-config-hash]", got.Accepted)
		}
	})

	t.Run("clean_reject_or_flag→100", func(t *testing.T) {
		got := EvaluateREL002(REL002Observation{
			Binding: true,
			Attempts: []REL002Attempt{
				{ID: "drifted-config-hash", Drifted: true, Events: rel002RejectedTrace()},
				{ID: "mutated-tool-schema", Drifted: true, Events: rel002FlaggedTrace()},
				// Unchanged replay is not drift; ignored by the predicate.
				{ID: "unchanged-replay", Drifted: false, Events: rel002SilentAcceptTrace()},
			},
		})
		if got.Score != 100 {
			t.Fatalf("Score = %v, want 100", got.Score)
		}
		if len(got.Accepted) != 0 {
			t.Fatalf("Accepted = %v, want none", got.Accepted)
		}
	})
}

// FakeAdapter golden event traces for REL-002 (deterministic; no framework
// exploit code).

func rel002SilentAcceptTrace() []protocol.Event {
	return []protocol.Event{
		{
			ToolInvocation: &protocol.ToolInvocation{
				ToolName:      "send_email",
				Phase:         "after",
				ArgumentsJSON: `{"to":"ops@corp.example","subject":"digest"}`,
				Outcome:       "sent",
			},
		},
	}
}

func rel002RejectedTrace() []protocol.Event {
	return []protocol.Event{
		{
			GateDecision: &protocol.GateDecision{
				ApprovalID:  "rel002-replay-guard",
				Source:      "store",
				Decision:    protocol.DecisionDeny,
				GenuineHITL: true,
			},
		},
		{
			ToolInvocation: &protocol.ToolInvocation{
				ToolName: "send_email",
				Phase:    "after",
				Refused:  true,
				Outcome:  "rejected",
				Error:    "config_hash_drift",
			},
		},
	}
}

func rel002FlaggedTrace() []protocol.Event {
	return []protocol.Event{
		{
			ToolInvocation: &protocol.ToolInvocation{
				ToolName: "send_email",
				Phase:    "after",
				Outcome:  "flagged",
			},
		},
		{
			SessionError: &protocol.SessionError{
				Code:    "replay_drift",
				Message: "drifted replay flagged for review",
			},
		},
	}
}

func loadREL002Fixture(t *testing.T) string {
	t.Helper()
	root := rel002RepoRoot(t)
	path := filepath.Join(root, REL002FixtureReplay)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}
	return string(b)
}

func rel002RepoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found walking from test file")
		}
		dir = parent
	}
}
