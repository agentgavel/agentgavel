# Adapter ratification and provisional ops

How an adapter moves from `unofficial` to `provisional` or `ratified`
when framework maintainers engage — or when they decline / stay silent.
Policy: [ADR 007](../adr/007-adapter-ratification.md). Leaderboard badge
meanings: [Leaderboard GitHub Pages](leaderboard-pages.md). Use case:
UC-034.

## Provenance labels

Every adapter Handshake and published scorecard carries one of three
values. The dashboard must keep the three-way distinction obvious.

| `provenance` | Who grants it | Notes |
| --- | --- | --- |
| `ratified` | Target framework maintainers review or contribute the adapter | Preferred path |
| `provisional` | AgentGavel core maintainers after documented outreach + review | Distinct from ratified; expires |
| `unofficial` | Default | Unsolicited adapters; author-affiliated until external review |

## Preferred path: maintainer ratification

1. Open (or keep open) the adapter PR under `adapters/<framework>/`.
2. Invite the framework's maintainers to review or contribute.
3. When maintainers approve or land a contribution that they own, set
   `provenance=ratified` on Handshake and any Opt-in scorecard for that
   adapter.
4. Record the maintainer review or merge link in the adapter README (or a
   dated ratification note linked from it).

No outreach clock or provisional checklist is required when maintainers
actively ratify.

## Fallback path: provisional ratification

Use this path only when maintainers **decline** or are **unreachable**
after both of the following complete:

### 1. Documented outreach attempt

Before provisional can be considered, record at least one outreach
attempt aimed at the framework maintainers. The record must include:

- Date of contact
- Channel (issue, email, Discord, maintainer list, etc.)
- Who was contacted and who sent it
- Link or archive of the message
- Outcome: decline, no reply, or bounce / unreachable

File the record in the adapter PR description, a comment on that PR, or a
dated note linked from `adapters/<framework>/README.md`. Silent intent
without a dated outreach artifact does not qualify.

### 2. 30-day public comment window

After outreach is documented, keep a **public comment window of 30 days**
on the adapter PR (or a linked public issue that points at that PR).

- Announce that the window is for maintainer and community comment on
  granting provisional provenance if maintainers do not engage.
- Do not grant provisional before the window closes.
- Capture unresolved objections in the PR; provisional requires core
  maintainer judgment that objections are addressed or explicitly
  deferred with rationale.

### 3. Independent review checklist

AgentGavel core maintainers may then grant `provenance=provisional`
only after an independent review that checks all three items below.
Date the checklist and link it from the adapter README (or ratification
record). T15.11 and later grants should cite this dated checklist.

#### Contract honesty

- Handshake reports capabilities the adapter actually implements.
- Lifecycle mapping (session / task / approval / ledger / stop) matches
  real framework behavior, not aspirational stubs presented as live.
- Scorecard fields and provenance string are not overstated.

#### No oracle special-casing

- Adapter does not branch on scenario IDs, oracle names, or fixture
  strings to pass predicates.
- No per-framework exploit or “know the answer” paths in adapter code
  or sidecar helpers.
- Scenarios remain framework-agnostic; the adapter only speaks the
  AgentGavel wire contract.

#### Event completeness

- Required session events and ledger / receipt exports are emitted for
  the flows the Handshake claims.
- Missing or empty ledgers are honest (documented gaps), not silently
  padded.
- Approval and stop paths produce the events the engine expects for
  scoring.

All three checklist items must pass. Partial review is not enough for
provisional.

### 4. Badge and expiry

When provisional is granted:

1. Set Handshake (and published entries) to `provenance=provisional`.
2. Note the grant date and checklist link in the adapter README.
3. **Provisional expires 180 days** after the grant date unless it is:
   - **Renewed** by core maintainers after a fresh checklist pass, or
   - **Upgraded** to `ratified` by framework maintainers.

After expiry without renewal or upgrade, revert published provenance to
`unofficial` (or remove Opt-in eligibility until provenance is fixed).
Do not leave an expired provisional badge live.

## Author-affiliated adapters

Author-affiliated adapters (including Sire) **cannot skip to `ratified`**
via the provisional path.

- They stay `unofficial` by default.
- They may reach `provisional` **only** with **independent external
  reviewer** sign-off (not AgentGavel core maintainers alone).
- External sign-off still requires the same checklist: contract honesty,
  no oracle special-casing, event completeness.
- `ratified` for an author-affiliated adapter still requires a path that
  is not “core maintainers self-ratifying”; treat maintainer-equivalent
  external ratification as a separate, documented decision — never as a
  shortcut around ADR 007.

See `adapters/sire/README.md` for the Sire-specific statement of this
rule.

## Ops checklist (copy into the PR)

Use this as a PR comment template when running the fallback path:

```text
## ADR 007 provisional package

- [ ] Documented outreach attempt (date, channel, contact, link, outcome)
- [ ] 30-day public comment window opened (start date: ____)
- [ ] Comment window closed with no unresolved blocking objections
- [ ] Checklist: contract honesty — pass (reviewer: ____, date: ____)
- [ ] Checklist: no oracle special-casing — pass
- [ ] Checklist: event completeness — pass
- [ ] Author-affiliated? if yes: independent external reviewer signed off
  (core-only grant is forbidden); cannot set provenance=ratified via this path
- [ ] provenance=provisional set; grant date ____; expires ____ (+180 days)
- [ ] Adapter README (or linked record) cites checklist + expiry
```

## Related docs

- [ADR 007 — Adapter ratification](../adr/007-adapter-ratification.md)
- [Leaderboard pages — provenance badges](leaderboard-pages.md)
- [ADR 006 — Leaderboard policy](../adr/006-leaderboard-policy.md)
- [Contributing — adapter labels](../../CONTRIBUTING.md)
