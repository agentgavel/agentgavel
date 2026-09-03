# Scenario Governance

This document describes how new security scenarios (`SEC-*`),
governance scenarios (`GOV-*`), and reliability scenarios (`REL-*`)
move from proposal to published status in AgentGavel.

## Lifecycle

1. **Propose.** Anyone may propose a new `SEC`, `GOV`, or `REL` scenario by
   opening a draft under `suites/security/`, `suites/governance/`, or
   `suites/reliability/` and filing a governance proposal that names the
   scenario ID, threat model, and oracle coverage. A proposal must identify
   itself as `SEC`, `GOV`, or `REL` in its title so reviewers can route it
   to the right domain owner.
2. **Draft review.** Drafts are unpublished scenarios. They run in CI
   for signal, but an unpublished/draft scenario never affects a
   session's published score — draft results are informational only
   until the scenario graduates.
3. **Comment window.** Once a proposal is opened, it enters a
   14-day comment window during which domain owners and contributors
   may raise objections, request changes, or flag oracle gaps. No
   scenario may be published before its comment window closes.
4. **Publish.** After the comment window closes with no unresolved
   objections, the scenario is merged and marked published. From that
   point forward, its results affect the published score for any
   session that exercises it.

## Ownership

- `SEC` scenarios are owned by the security domain reviewers.
- `GOV` scenarios are owned by the governance domain reviewers.
- `REL` scenarios are owned by the reliability domain reviewers.

All three domains follow the same propose → draft review → comment window →
publish lifecycle described above.

## Reliability suite (REL-v0)

REL-v0 (`REL-001` Demotion Latch, `REL-002` Replay / Config Drift,
`REL-003` Ledger Completeness) is published under
[ADR 010](../adr/010-reliability-suite.md) and lives in
`suites/reliability/`. REL-v1 additions or predicate changes require the
same 14-day comment window before they may move published leaderboard
scores.
