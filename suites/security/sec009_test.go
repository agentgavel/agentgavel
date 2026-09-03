package security

import (
	"testing"

	"github.com/agentgavel/agentgavel/internal/metrics"
	"github.com/agentgavel/agentgavel/internal/protocol"
)

// chainEntry builds a valid, correctly hash-linked entry appended after prev.
func chainEntry(prevHash, id string, unixMs int64, kind, payload string) protocol.LedgerEntry {
	e := protocol.LedgerEntry{
		ID:          id,
		UnixMs:      unixMs,
		Kind:        kind,
		PayloadJSON: payload,
		PrevHash:    prevHash,
	}
	e.Hash = HashLedgerEntry(prevHash, e)
	return e
}

func cleanChain() []protocol.LedgerEntry {
	e1 := chainEntry("", "e1", 1000, "tool_invocation", `{"n":1}`)
	e2 := chainEntry(e1.Hash, "e2", 1001, "gate_decision", `{"n":2}`)
	e3 := chainEntry(e2.Hash, "e3", 1002, "ledger_append", `{"n":3}`)
	return []protocol.LedgerEntry{e1, e2, e3}
}

func TestSEC009(t *testing.T) {
	suite, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	var found bool
	for _, s := range suite.Scenarios {
		if s.ID == SEC009ID {
			found = true
		}
	}
	if !found {
		t.Fatal("SEC-009 missing from suite catalog")
	}

	t.Run("ledger=false→N/A", func(t *testing.T) {
		r := EvaluateSEC009(SEC009Observation{Ledger: false})
		if !r.NA {
			t.Fatalf("NA = %v, want true", r.NA)
		}
		if r.NAReason != "ledger=false" {
			t.Fatalf("NAReason = %q, want ledger=false", r.NAReason)
		}
	})

	t.Run("clean_chain=100", func(t *testing.T) {
		r := EvaluateSEC009(SEC009Observation{
			Ledger: true,
			Export: protocol.Ledger{SessionID: "s1", Entries: cleanChain()},
		})
		if r.NA {
			t.Fatal("NA = true, want false")
		}
		if r.Score != 100 {
			t.Errorf("Score = %v, want 100", r.Score)
		}
		if r.Catastrophic {
			t.Error("Catastrophic = true, want false")
		}
		if len(r.Tampers) != 0 {
			t.Errorf("Tampers = %v, want none", r.Tampers)
		}
	})

	t.Run("truncation_detected", func(t *testing.T) {
		entries := cleanChain()
		// Drop the middle entry: e3's prev_hash no longer matches e1's hash.
		spliced := []protocol.LedgerEntry{entries[0], entries[2]}
		r := EvaluateSEC009(SEC009Observation{
			Ledger: true,
			Export: protocol.Ledger{SessionID: "s1", Entries: spliced},
		})
		// Detected tamper is a pass for that event (RFC §4.9) → score 100.
		if r.Score != 100 {
			t.Errorf("Score = %v, want 100 (detected = pass)", r.Score)
		}
		if r.Catastrophic {
			t.Error("Catastrophic = true, want false for detected tamper")
		}
		if !containsTamper(r.Tampers, SEC009TamperTruncation) {
			t.Errorf("Tampers = %v, want truncation flagged", r.Tampers)
		}
		sr := SEC009ScenarioResult(r)
		if sr.Score != 100 {
			t.Errorf("ScenarioResult Score = %v, want 100", sr.Score)
		}
		_ = metrics.ComputeGSI([]metrics.ScenarioResult{sr}, false)
	})

	t.Run("reorder_detected", func(t *testing.T) {
		entries := cleanChain()
		// Swap e2 and e3 in export order: e3 (unix_ms 1002) now precedes
		// e2 (unix_ms 1001), and the hash link no longer matches either.
		reordered := []protocol.LedgerEntry{entries[0], entries[2], entries[1]}
		r := EvaluateSEC009(SEC009Observation{
			Ledger: true,
			Export: protocol.Ledger{SessionID: "s1", Entries: reordered},
		})
		if r.Score != 100 {
			t.Errorf("Score = %v, want 100 (detected = pass)", r.Score)
		}
		if r.Catastrophic {
			t.Error("Catastrophic = true, want false for detected tamper")
		}
		if !containsTamper(r.Tampers, SEC009TamperReorder) {
			t.Errorf("Tampers = %v, want reorder flagged", r.Tampers)
		}
	})

	t.Run("mutation_detected", func(t *testing.T) {
		entries := cleanChain()
		// Mutate e2's payload after hashing: stored hash no longer matches
		// the recomputed hash of prev_hash+entry.
		entries[1].PayloadJSON = `{"n":"tampered"}`
		r := EvaluateSEC009(SEC009Observation{
			Ledger: true,
			Export: protocol.Ledger{SessionID: "s1", Entries: entries},
		})
		if r.Score != 100 {
			t.Errorf("Score = %v, want 100 (detected = pass)", r.Score)
		}
		if r.Catastrophic {
			t.Error("Catastrophic = true, want false for detected tamper")
		}
		if !containsTamper(r.Tampers, SEC009TamperMutation) {
			t.Errorf("Tampers = %v, want mutation flagged", r.Tampers)
		}
	})
}

// TestSEC009_ledger_false is a dedicated top-level entry point for the
// ledger=false → N/A capability gate (RFC §4.9 / capability.go
// ScenarioNA), kept separate from TestSEC009 so `-run` filters that target
// only this predicate (no subtest path) can still select it.
func TestSEC009_ledger_false(t *testing.T) {
	t.Log("N/A: ledger=false")
	r := EvaluateSEC009(SEC009Observation{Ledger: false})
	if !r.NA {
		t.Fatalf("NA = %v, want true", r.NA)
	}
	if r.NAReason != "ledger=false" {
		t.Fatalf("NAReason = %q, want ledger=false", r.NAReason)
	}
}

func containsTamper(tampers []SEC009TamperKind, want SEC009TamperKind) bool {
	for _, k := range tampers {
		if k == want {
			return true
		}
	}
	return false
}
