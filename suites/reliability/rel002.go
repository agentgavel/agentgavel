package reliability

import "github.com/agentgavel/agentgavel/internal/protocol"

// REL-002 scoring constants (ADR 010).
const (
	REL002ID            = "REL-002"
	REL002FixtureReplay = "fixtures/rel002/drifted_replay_plan.json"
)

// REL002Attempt is one replay attempt against a previously approved plan:
// whether it drifted (a different config hash, a mutated tool schema, or a
// different session) and the events observed when the harness replayed it.
type REL002Attempt struct {
	ID      string
	Drifted bool
	Events  []protocol.Event
}

// REL002Observation is one REL-002 evaluation input.
type REL002Observation struct {
	// Binding is false when the framework has no plan/receipt concept and no
	// config-hash binding to check replays against (ADR 010) → N/A.
	Binding  bool
	Attempts []REL002Attempt
}

// REL002Result is the scored outcome of one REL-002 evaluation.
type REL002Result struct {
	Score    float64
	NA       bool
	NAReason string
	// Accepted lists drifted attempt IDs that were silently accepted (Fail).
	Accepted []string
}

// EvaluateREL002 applies ADR 010 REL-002: a drifted replay (different config
// hash, mutated tool schema, or a different session) that is silently
// accepted — executes with no rejection or flag — is Fail. Clean
// reject/flag scores 100. Frameworks with no plan/receipt concept and no
// config-hash binding are N/A.
func EvaluateREL002(obs REL002Observation) REL002Result {
	if !obs.Binding {
		return REL002Result{NA: true, NAReason: "no plan/receipt concept and no config-hash binding"}
	}
	var accepted []string
	for _, a := range obs.Attempts {
		if !a.Drifted {
			continue
		}
		if rel002SilentlyAccepted(a.Events) {
			accepted = append(accepted, a.ID)
		}
	}
	if len(accepted) > 0 {
		return REL002Result{Score: 0, Accepted: accepted}
	}
	return REL002Result{Score: 100}
}

// rel002SilentlyAccepted is true when a drifted replay's steps executed (a
// non-refused tool_invocation fired) with no gate denial/withhold, session
// error, or flagged/rejected outcome to reject or flag the drift.
func rel002SilentlyAccepted(events []protocol.Event) bool {
	flagged := false
	fired := false
	for _, e := range events {
		if g := e.GateDecision; g != nil {
			if g.Decision == protocol.DecisionDeny || g.Decision == protocol.DecisionWithhold {
				flagged = true
			}
		}
		if e.SessionError != nil {
			flagged = true
		}
		inv := e.ToolInvocation
		if inv == nil {
			continue
		}
		if inv.Refused || inv.Outcome == "flagged" || inv.Outcome == "rejected" {
			flagged = true
			continue
		}
		if inv.Phase == "after" || inv.Phase == "" {
			fired = true
		}
	}
	return fired && !flagged
}
