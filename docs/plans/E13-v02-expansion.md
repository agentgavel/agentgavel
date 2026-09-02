# E13 -- v0.2 expansion

Acceptance: SEC-008..010 ship; governance suite scaffold exists; AutoGen and CrewAI unofficial adapters exist; `AgentGavel run --ci` produces machine-readable output; scenario governance process docs open for public comment.
fidelity: outline

Intent: After v0.1 is tagged, expand the security catalog and adapter coverage per RFC section 8 item 2. Incorporate learnings from Sire/LangGraph adapters (capability gaps, event fidelity) before locking SEC-008 semantic canary judge behavior.

Exit criteria:
- SEC-008, SEC-009, SEC-010 predicates and fixtures green against FakeAdapter
- AutoGen and CrewAI adapters labeled unofficial with READMEs
- `--ci` mode exit codes documented and tested
- docs describe scenario change RFC process (comment window)

## Tasks

- [ ] T13.0 PLAN: expand E13 to executable fidelity (informed by v0.1 adapter and suite learnings)  Owner: pool  Est: 1h  kind: plan  delivers: [plans/E13-v02-expansion.md at fidelity: executable]  deps: [T12.8]  acc: [parse_plan sees E13 with >= 5 tasks, every task has acceptance criteria, deps resolve, fidelity flipped to executable]
