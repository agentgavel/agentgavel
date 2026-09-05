# E15 -- v1.0 public submission and harness red-team

Acceptance: Public signed submission process live; adversarial red-team bounty for the harness itself announced; ratification and provisional paths operational per ADR 007.
fidelity: outline

Intent: Open AgentGavel to external score submissions and stress-test the harness credibility under bounty pressure. Only after leaderboard (E14) and `v0.3.0` exist.

Exit criteria:
- Signed Opt-in submission via GitHub-native PR workflow documented and enforced ([ADR 012](../adr/012-github-native-signed-submissions.md))
- Bounty scope and disclosure policy published
- At least one non-author adapter reaches ratified or provisional

## Locked (pre-T15.0)

- **Submission host:** GitHub-native (signed artifact → PR into
  `dashboard/data/` → CI signature verify → merge → Pages). Not Firebase
  as trust root or primary write path (ADR 012).

## Tasks

- [ ] T15.0 PLAN: expand E15 to executable fidelity (informed by v0.3 leaderboard operations)  Owner: pool  Est: 1h  kind: plan  delivers: [plans/E15-v10-public-process.md at fidelity: executable]  deps: [T14.16]  acc: [parse_plan sees E15 with >= 5 tasks, fidelity executable]
  - Note: waits on v0.3.0 tag so planning incorporates REL-v0 / rubber-stamp / Pages learnings.
  - Must prescribe: maintainer key registry format, canonical signed payload, CI verify job, Opt-in CI rule flip (`sample` no longer required once signature verifies), and bounty docs.
