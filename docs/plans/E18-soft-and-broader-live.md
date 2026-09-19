# E18 -- Soft model-mode and broader live adapter set

Acceptance: Soft rates with >=25 seeds and Wilson intervals run in an
opt-in job against at least one live adapter; additional §8.1 adapters
can be promoted to live under ADR 015 without changing the wire protocol.
fidelity: executable

Intent: After live Hard/oracle baselines exist (E16 FakeAdapter + E17
LangGraph live), schedule expensive Soft campaigns and expand live coverage
beyond Sire/LangGraph. Soft stays opt-in (cost / latency).

Exit criteria: documented model-mode runbook; one Soft scorecard published
or explicitly deferred with cost rationale; ADR/scenario governance followed
for any new SEC/REL IDs.

Learnings from E16/E17:
- Oracle Hard runs are cheap locally; Soft needs real model spend + seeds≥25.
- Env-gated live bootstrap (Sire / LangGraph / Hermes / OpenClaw) is the
  pattern — reuse for any new §8.1 promotion.
- Cloud adapters need publicly reachable Oracle URLs (Sire Wave 43 lesson).

## Wave 49 -- Plan + Soft runbook

- [x] T18.0 PLAN: expand E18 to executable fidelity (informed by E16/E17 live costs)  Owner: pool  Est: 1h  kind: plan  delivers: [plans/E18-soft-and-broader-live.md at fidelity: executable]  deps: [T17.0]  acc: [parse_plan sees E18 with >= 5 tasks, every task has acceptance criteria, deps resolve, fidelity flipped to executable]  completed: 2026-09-19

- [x] T18.1 Document Soft model-mode ops (seeds≥25, Wilson, cost caveats) in benchmark-ops  Owner: pool  Est: 60m  kind: agent  verifies: [UC-042]  lane: agent  deps: [T18.0]  acc: [benchmark-ops.md Soft section states --mode model, seeds≥25, Wilson intervals, and that Soft is opt-in not CI-default]  completed: 2026-09-19

- [x] T18.2 Wire/verify CLI Soft path exits cleanly when model endpoint missing (fail closed)  Owner: pool  Est: 60m  kind: agent  verifies: [UC-042]  deps: [T18.1]  acc: [AgentGavel run --mode model without model URL exits non-zero with clear stderr; unit or CLI test covers]  completed: 2026-09-19

## Wave 50 -- First Soft evidence (opt-in)

- [ ] T18.3 Human: provision Soft model credentials / endpoint for one adapter (LangGraph live or FakeAdapter+model)  Owner: founder  Est: 45m  kind: human  verifies: [UC-042]  deps: [T18.1]  delivers: [scratch pointer; never commit secrets]  blocked: needs model API access

- [ ] T18.4 Run Soft security suite ≥25 seeds against one live/harness target; record summary + Wilson  Owner: pool  Est: 120m  kind: agent  verifies: [UC-042]  deps: [T18.2, T18.3]  acc: [summary has Soft rates with intervals; runtime honest; path cited in scratch/devlog]

- [ ] T18.5 Publish or defer Soft Unratified with explicit cost rationale in Progress Log  Owner: pool  Est: 45m  kind: agent  verifies: [UC-038, UC-042]  deps: [T18.4]  acc: [dashboard sample OR plan Progress Log defers Soft publish with cost reason]

## Wave 51 -- Broader live §8.1 promotions

- [x] T18.6 Promote one additional stub adapter (CrewAI or ADK) to optional live bootstrap pattern  Owner: pool  Est: 120m  kind: agent  verifies: [UC-043]  deps: [T17.3]  acc: [adapter has stub default + env live path; Handshake runtime flips only on proven live; pytest covers]  completed: 2026-09-19

- [x] T18.7 CI ruff matrix covers all Python adapter trees touched by live work  Owner: pool  Est: 30m  kind: agent  verifies: [infrastructure]  deps: [T18.6]  acc: [ci.yml ruff matrix includes langgraph hermes openclaw sdk]  completed: 2026-09-19

- [ ] T18.8 make test + make lint + check-dashboard green  Owner: pool  Est: 45m  kind: agent  verifies: [infrastructure]  deps: [T18.2, T18.7]  acc: [all three exit 0]

- [ ] T18.9 Update plan.md frontier when Soft evidence or deferral lands  Owner: pool  Est: 15m  kind: agent  verifies: [infrastructure]  deps: [T18.5, T18.8]  acc: [Progress Log notes Soft outcome; E18 checkbox count accurate]
