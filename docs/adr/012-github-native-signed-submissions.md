# ADR 012: GitHub-Native Signed Opt-in Submissions

## Status
Accepted

## Date
2026-09-05

## Context
E15 (v1.0) must land maintainer-signed Opt-in submissions so ADR 006’s
primary tab is no longer sample-only. Candidates for a “hosted signed-
submission API” included Firebase, a thin custom verify endpoint (Cloud
Run / Worker), and a GitHub-native path that reuses the existing Pages
publish loop (`dashboard/data/` → PR → `pages.yml`).

Trust for Opt-in is cryptographic (signatures bound to registered
maintainer keys), not identity-provider login. Submission volume is low
(framework maintainers). The leaderboard already ships only from `main`
via static files.

## Decision
v1.0 Opt-in submissions are **GitHub-native**:

1. Submitters produce a signed scorecard artifact (fingerprint + entry
   JSON + signature over a fixed canonical encoding).
2. Submission is a pull request (or equivalent GitHub artifact upload that
   becomes a PR) adding `dashboard/data/<run_id>.json` and updating
   `index.json`.
3. CI verifies the signature against a published registry of maintainer
   keys, enforces schema / ADR 006–007 rules, then human or bot merge
   publishes via the existing Pages workflow.
4. There is **no** separate Firebase (or similar) hosted API as the trust
   root or the primary write path in v1.0.

A thin verify helper (local CLI or optional CI action) may exist to check
signatures before open-PR; it does not replace GitHub as the system of
record.

## Consequences
Positive: zero new host; inspectable history; same publish path as today;
signatures remain the spam gate. Negative: submitters need a GitHub
account and PR hygiene; key registry and signature format must be
specified in E15 executable tasks. Firebase Auth/Firestore stay out of
scope unless a later product UI needs them — they are not the Opt-in
trust layer.

## Alternatives considered
- **Firebase / BaaS API** — convenient auth and storage; wrong trust root
  for maintainer keys and adds ops surface for low volume.
- **Standalone verify API that writes storage** — useful later as DX; still
  needs GitHub/Pages or equivalent for the public leaderboard unless we
  rebuild hosting.
