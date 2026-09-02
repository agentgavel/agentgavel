# E14 -- v0.3 reliability, rubber-stamp, leaderboard

Acceptance: Reliability suite ships; `AgentGavel rubber-stamp` runs SEC-002 and SEC-006 only; static dashboard on GitHub Pages shows opt-in and unratified tabs per ADR 006.
fidelity: outline

Intent: Add fault-recovery scenarios, a fast rubber-stamp CLI, and the public leaderboard surface. Depends on v0.2 CI mode and scenario governance so published scores have a change-control process.

Exit criteria:
- Reliability suite scenarios defined and tested
- rubber-stamp command green in CI
- dashboard/ static site publishes sample scorecards with provenance labels

## Tasks

- [ ] T14.0 PLAN: expand E14 to executable fidelity (informed by v0.2)  Owner: pool  Est: 1h  kind: plan  delivers: [plans/E14-v03-reliability-leaderboard.md at fidelity: executable]  deps: [T13.0]  acc: [parse_plan sees E14 with >= 5 tasks, fidelity executable]
  - Note: trigger is completion of E13 planning+implementation milestone; when E13 is expanded, update this dep to E13 exit tasks.
