// Package engine orchestrates scenarios, seeds, and adapter sessions.
package engine

import (
	"context"

	"github.com/agentgavel/gavel/internal/protocol"
)

// Scenario is one adversarial test with deterministic scoring hooks.
type Scenario interface {
	ID() string
	Setup(ctx context.Context, sess *Session) error
	Probe(ctx context.Context, sess *Session) error
	Observe(ctx context.Context, sess *Session) ([]protocol.Event, error)
	Predicate(ctx context.Context, sess *Session, events []protocol.Event) (Result, error)
	Score(ctx context.Context, sess *Session, result Result) (float64, error)
}

// Result is the predicate outcome for a scenario seed.
type Result struct {
	HardEligible bool
	Pass         bool
	Fail         bool
	Leaks        []string
	Catastrophic bool
	CatastrophicKind string
	Notes        string
}

// Session holds per-seed adapter state.
type Session struct {
	ID     string
	Seed   int64
	Mode   string // oracle | model
	Client *protocol.Client
	Config protocol.SessionConfig
}

// Runner executes scenarios across seeds.
type Runner struct {
	Scheduler *SeedScheduler
}

// Ensure Scenario method set is referenced for compile-time interface checks in tests.
var _ Scenario = (*noopScenario)(nil)

type noopScenario struct{}

func (n *noopScenario) ID() string { return "NOOP" }
func (n *noopScenario) Setup(context.Context, *Session) error {
	return nil
}
func (n *noopScenario) Probe(context.Context, *Session) error { return nil }
func (n *noopScenario) Observe(context.Context, *Session) ([]protocol.Event, error) {
	return nil, nil
}
func (n *noopScenario) Predicate(context.Context, *Session, []protocol.Event) (Result, error) {
	return Result{Pass: true, HardEligible: true}, nil
}
func (n *noopScenario) Score(context.Context, *Session, Result) (float64, error) {
	return 100, nil
}
