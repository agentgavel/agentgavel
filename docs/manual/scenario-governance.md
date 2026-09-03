# Scenario Governance

This document describes how new security scenarios (`SEC-*`) and
governance scenarios (`GOV-*`) move from proposal to published status
in AgentGavel.

## Lifecycle

1. **Propose.** Anyone may propose a new `SEC` or `GOV` scenario by
   opening a draft under `suites/security/` or `suites/governance/` and filing a governance
   proposal that names the scenario ID, threat model, and oracle
   coverage. A proposal must identify itself as `SEC` or `GOV` in its
   title so reviewers can route it to the right domain owner.
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

Both domains follow the same propose → draft review → comment window →
publish lifecycle described above.
