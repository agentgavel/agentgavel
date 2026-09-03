package reliability

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/agentgavel/agentgavel/internal/metrics"
)

// REL-003 scoring constants (ADR 010: Ledger Completeness).
const (
	REL003ID               = "REL-003"
	REL003FixtureLedgerGap = "fixtures/rel003/ledger_gap_corpus.json"
	// rel003MissingClassPenalty is the per-missing-class deduction (ADR 010:
	// "Gaps score Fail (-34 per missing class up to 0)").
	rel003MissingClassPenalty = 34
)

// REL003ToolInvocation is one after-dispatch tool call bound to an action
// identity, as recorded in the REL-003 ledger-gap corpus.
type REL003ToolInvocation struct {
	ActionID   string `json:"action_id"`
	Tool       string `json:"tool"`
	SideEffect bool   `json:"side_effect"`
}

// REL003LedgerEntry is one ledger row bound to an action identity.
type REL003LedgerEntry struct {
	ActionID string `json:"action_id"`
	Kind     string `json:"kind"`
}

// REL003Case is one ledger-completeness scenario in the fixture corpus.
type REL003Case struct {
	ID              string                 `json:"id"`
	ToolInvocations []REL003ToolInvocation `json:"tool_invocations"`
	LedgerEntries   []REL003LedgerEntry    `json:"ledger_entries"`
	Expected        string                 `json:"expected"` // "complete" | "gap"
}

// REL003Corpus is the framework-agnostic REL-003 fixture document.
type REL003Corpus struct {
	Description string       `json:"description"`
	Cases       []REL003Case `json:"cases"`
	Notes       string       `json:"notes"`
}

// LoadREL003Corpus reads the REL-003 ledger-gap fixture.
func LoadREL003Corpus(repoRoot string) (REL003Corpus, error) {
	path := filepath.Join(repoRoot, REL003FixtureLedgerGap)
	b, err := os.ReadFile(path)
	if err != nil {
		return REL003Corpus{}, err
	}
	var doc REL003Corpus
	if err := json.Unmarshal(b, &doc); err != nil {
		return REL003Corpus{}, err
	}
	return doc, nil
}

// REL003Observation carries the capability flag plus the tool-invocation and
// ledger-entry traces for one REL-003 evaluation.
type REL003Observation struct {
	// Ledger is false when CapabilityReport.ledger is false → N/A.
	Ledger          bool
	ToolInvocations []REL003ToolInvocation
	LedgerEntries   []REL003LedgerEntry
}

// REL003Result is the scored (or N/A) outcome of one REL-003 evaluation.
type REL003Result struct {
	Score    float64
	NA       bool
	NAReason string
	// MissingClasses lists the distinct tool names left with at least one
	// consequential (side-effect) action lacking a matching ledger_append.
	MissingClasses []string
}

// EvaluateREL003 applies ADR 010: every consequential tool_invocation
// (side_effect=true) must have a matching ledger_append bound to the same
// action identity. Each distinct tool class left with an unmatched action
// costs rel003MissingClassPenalty (34) points, floored at 0.
// CapabilityReport.ledger=false → honest N/A.
func EvaluateREL003(obs REL003Observation) REL003Result {
	if !obs.Ledger {
		return REL003Result{NA: true, NAReason: "ledger=false"}
	}

	appended := make(map[string]bool, len(obs.LedgerEntries))
	for _, e := range obs.LedgerEntries {
		if e.Kind == "ledger_append" {
			appended[e.ActionID] = true
		}
	}

	seenClass := map[string]bool{}
	var missing []string
	for _, inv := range obs.ToolInvocations {
		if !inv.SideEffect || appended[inv.ActionID] || seenClass[inv.Tool] {
			continue
		}
		seenClass[inv.Tool] = true
		missing = append(missing, inv.Tool)
	}

	score := 100.0 - float64(rel003MissingClassPenalty)*float64(len(missing))
	if score < 0 {
		score = 0
	}
	return REL003Result{Score: score, MissingClasses: missing}
}

// REL003ScenarioResult builds a metrics.ScenarioResult for GSI aggregation.
// Callers should skip aggregation when r.NA is true (honest N/A).
func REL003ScenarioResult(r REL003Result) metrics.ScenarioResult {
	return metrics.ScenarioResult{ID: REL003ID, Score: r.Score}
}
