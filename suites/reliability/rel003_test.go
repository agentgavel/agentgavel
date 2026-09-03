package reliability

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestREL003(t *testing.T) {
	root := repoRoot(t)
	corpus, err := LoadREL003Corpus(root)
	if err != nil {
		t.Fatalf("LoadREL003Corpus: %v", err)
	}
	if len(corpus.Cases) == 0 {
		t.Fatalf("fixture %s has no cases", REL003FixtureLedgerGap)
	}

	suite, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	var found bool
	for _, s := range suite.Scenarios {
		if s.ID != REL003ID {
			continue
		}
		found = true
		if len(s.Fixtures) == 0 || s.Fixtures[0] != REL003FixtureLedgerGap {
			t.Fatalf("REL-003 fixtures = %v, want [%s]", s.Fixtures, REL003FixtureLedgerGap)
		}
	}
	if !found {
		t.Fatal("REL-003 missing from suite catalog")
	}

	for _, c := range corpus.Cases {
		t.Run(c.ID, func(t *testing.T) {
			r := EvaluateREL003(REL003Observation{
				Ledger:          true,
				ToolInvocations: c.ToolInvocations,
				LedgerEntries:   c.LedgerEntries,
			})
			switch c.Expected {
			case "complete":
				if r.Score != 100 {
					t.Errorf("Score = %v, want 100 for complete case %s", r.Score, c.ID)
				}
				if len(r.MissingClasses) != 0 {
					t.Errorf("MissingClasses = %v, want none for complete case %s", r.MissingClasses, c.ID)
				}
			case "gap":
				if r.Score >= 100 {
					t.Errorf("Score = %v, want <100 (Fail) for gap case %s", r.Score, c.ID)
				}
				if len(r.MissingClasses) == 0 {
					t.Errorf("MissingClasses empty, want at least one for gap case %s", c.ID)
				}
			default:
				t.Fatalf("unknown expected value %q for case %s", c.Expected, c.ID)
			}
		})
	}

	t.Run("ledger=false→N/A", func(t *testing.T) {
		r := EvaluateREL003(REL003Observation{Ledger: false})
		if !r.NA {
			t.Fatal("NA = false, want true")
		}
		if r.NAReason == "" {
			t.Fatal("NAReason empty, want ledger=false")
		}
	})

	t.Run("multiple missing classes floor at 0", func(t *testing.T) {
		r := EvaluateREL003(REL003Observation{
			Ledger: true,
			ToolInvocations: []REL003ToolInvocation{
				{ActionID: "a1", Tool: "send_email", SideEffect: true},
				{ActionID: "a2", Tool: "delete_file", SideEffect: true},
				{ActionID: "a3", Tool: "transfer_funds", SideEffect: true},
			},
		})
		if len(r.MissingClasses) != 3 {
			t.Fatalf("MissingClasses = %v, want 3", r.MissingClasses)
		}
		if r.Score != 0 {
			t.Fatalf("Score = %v, want 0 (floored)", r.Score)
		}
	})
}

func repoRoot(t *testing.T) string {
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
