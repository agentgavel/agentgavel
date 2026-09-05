# E15 -- v1.0 public submission and harness red-team

Acceptance: Public signed Opt-in submission process live (GitHub-native per
ADR 012); adversarial red-team bounty for the harness itself announced;
ratification and provisional paths operational per ADR 007; at least one
non-author adapter reaches provisional or ratified.
fidelity: executable

## Learnings from v0.3 (bind into tasks)

- Pages + `agentgavel.dev` already ship from `dashboard/` via `pages.yml`.
  Opt-in trust must reuse that path (ADR 012) — no new host.
- `report --publish` is the sole writer of `dashboard/data/*.json` and
  `index.json` today, and rejects `--tab opt-in` (ADR 006 addendum).
  v1.0 flips that for **signed** Opt-in only (ADR 013).
- `scripts/check-dashboard.sh` enforces `tab=opt-in ⇒ sample=true`. Flip to
  `sample=true OR verified signature` in the same wave as the verifier.
- Samples stay: Opt-in demo rows with `sample: true` need no key. Real
  Opt-in rows use Ed25519 + `dashboard/keys/registry.json` (ADR 013).
- Sire remains author-affiliated: cannot reach `ratified` via the
  provisional path alone (ADR 007). Prefer LangGraph (or another
  non-author §8.1 adapter) for the first provisional/ratified badge.
- Soft rates, FakeAdapter oracle-first, and no per-framework exploit code
  still bind. Bounty targets the **harness** (engine, oracle, scoring,
  dashboard verify), not adapter exploit kits.

## Locked decisions

- Submission host: GitHub-native (ADR 012).
- Signature: Ed25519 + JCS-ish canonical JSON; registry under
  `dashboard/keys/` (ADR 013).
- No Firebase / BaaS as trust root in v1.0.

## File-touch map

| Task | Creates / modifies |
| ---- | ---- |
| T15.1 | `dashboard/keys/registry.json`, `dashboard/keys/README.md`, schema notes |
| T15.2 | `internal/submit/{canonical.go,verify.go,*_test.go}`, golden vectors |
| T15.3 | `cmd/AgentGavel/report.go` (`--sign`), `verify_entry.go`, tests |
| T15.4 | `.github/workflows/ci.yml` (opt-in verify job) and/or `scripts/verify-opt-in.sh` |
| T15.5 | `scripts/check-dashboard.sh`, `dashboard/data/schema.json` (`key_id`, `signature`) |
| T15.6 | `internal/publish`, `cmd/AgentGavel/report.go` Opt-in when signed |
| T15.7 | `docs/manual/opt-in-submission.md` |
| T15.8 | `SECURITY.md`, `docs/manual/harness-bounty.md` |
| T15.9 | `docs/manual/leaderboard-pages.md`, `README.md`, ADR 006 addendum note |
| T15.10 | `docs/manual/adapter-ratification.md` |
| T15.11 | adapter README provenance + registry/outreach record (LangGraph) |
| T15.12 | `docs/manual/v1.0-smoke.md` |
| T15.13 | repo-wide `make test` / `make lint` / `check-dashboard` |
| T15.14 | git tag `v1.0.0` + GitHub release assets |

## Wave 1 -- planning (Wave 28)

- [x] T15.0 PLAN: expand E15 to executable fidelity (informed by v0.3 leaderboard operations)  Owner: pool  Est: 1h  kind: plan  delivers: [plans/E15-v10-public-process.md at fidelity: executable]  deps: [T14.16]  acc: [parse_plan sees E15 with >= 5 tasks, fidelity executable]  completed: 2026-09-05
  - Prescribed: key registry, canonical signed payload, CI verify, Opt-in CI rule flip, bounty docs (ADR 012 + ADR 013).

## Wave 2 -- signature contract + CLI (Wave 29, 3 agents)

- [ ] T15.1 Add maintainer key registry under dashboard/keys (ADR 013)  Owner: pool  Est: 45m  kind: agent  verifies: [UC-032]  lane: agent  acc: [dashboard/keys/registry.json parses as a JSON array; each object has key_id, framework, alg=ed25519, public_key_b64, added_at, status in {active,revoked}; README documents how maintainers register a key via PR; at least one test/fixture active key exists for FakeAdapter or Example Framework]  deps: [T15.0]

- [ ] T15.2 Implement Ed25519 canonical encode + verify library (ADR 013)  Owner: pool  Est: 90m  kind: agent  verifies: [UC-032]  acc: [go test ./internal/submit -count=1 passes: CanonicalJSON sorts keys, Verify rejects tampered fields and revoked keys, Verify accepts a golden signed entry against the fixture public key]  deps: [T15.1]
  - Package `internal/submit` (or `internal/sign`): `CanonicalJSON([]byte) ([]byte, error)`, `Sign(priv, entry)`, `Verify(registry, entry)`.
  - Signature covers entry without `signature`/`key_id`/`sample`.

- [ ] T15.3 Wire AgentGavel report --sign and verify-entry CLI  Owner: pool  Est: 75m  kind: agent  verifies: [UC-032]  acc: [go test ./cmd/AgentGavel -run 'Sign|VerifyEntry' -count=1 passes: report --sign with a test key writes key_id+signature; verify-entry exits 0 on valid and 1 on tamper; help text lists both]  deps: [T15.2, T14.11]

## Wave 3 -- CI + Opt-in rule flip (Wave 30, 3 agents)

- [ ] T15.4 Add CI job that verifies Opt-in signatures on dashboard/data PRs  Owner: pool  Est: 60m  kind: agent  verifies: [UC-032]  acc: [ci.yml (or scripts/verify-opt-in.sh invoked by CI) fails when an opt-in entry has sample=false and a bad signature, and passes on the committed tree; job runs on pull_request and push to main]  deps: [T15.2, T14.21]

- [ ] T15.5 Flip check-dashboard Opt-in rule to sample OR verified signature (ADR 013)  Owner: pool  Est: 60m  kind: agent  verifies: [UC-022, UC-032]  acc: [bash scripts/check-dashboard.sh exits 0 on committed tree; exits 1 for tab=opt-in sample=false without valid signature; exits 0 for tab=opt-in sample=false with valid signature against registry; schema.json documents optional key_id and signature]  deps: [T15.2, T14.21]

- [ ] T15.6 Allow report --publish --tab opt-in only when entry is signed  Owner: pool  Est: 75m  kind: agent  verifies: [UC-022, UC-032]  acc: [go test ./cmd/AgentGavel -run ReportPublish -count=1 passes: --tab opt-in without sign exits non-zero citing ADR 013; with --sign (or pre-signed stdin) writes tab=opt-in sample=false and updates index.json; unratified path unchanged]  deps: [T15.3, T14.11]

## Wave 4 -- docs + bounty (Wave 31, 3 agents)

- [ ] T15.7 Document signed Opt-in PR submission workflow  Owner: pool  Est: 45m  kind: agent  delivers: [docs/manual/opt-in-submission.md]  verifies: [UC-032]  lane: agent  acc: [doc cites ADR 012 and ADR 013; lists generate scorecard, report --sign, open PR adding dashboard/data entry + index, CI verify, merge→Pages; states samples need no signature]  deps: [T15.3, T15.4]

- [ ] T15.8 Publish harness red-team bounty scope and disclosure policy  Owner: pool  Est: 60m  kind: agent  delivers: [SECURITY.md, docs/manual/harness-bounty.md]  verifies: [UC-033]  lane: agent  acc: [SECURITY.md links disclosure path; harness-bounty.md defines in-scope (engine, oracle, scoring, dashboard verify, CI signature checks), out-of-scope (per-framework exploit kits, social engineering), and safe-harbor / coordinated disclosure]  deps: [T15.0]

- [ ] T15.9 Update leaderboard Pages manual and README for v1.0 Opt-in signatures  Owner: pool  Est: 45m  kind: agent  delivers: [docs/manual/leaderboard-pages.md, README.md]  verifies: [UC-022, UC-032]  lane: agent  acc: [leaderboard-pages.md states ADR 006 addendum expired for signed Opt-in, links opt-in-submission.md and ADR 013; README mentions report --sign / verify-entry and SECURITY.md]  deps: [T15.5, T15.7, T14.12]

## Wave 5 -- ratification ops (Wave 32, 2 agents)

- [ ] T15.10 Document adapter ratification and provisional ops checklist (ADR 007)  Owner: pool  Est: 45m  kind: agent  delivers: [docs/manual/adapter-ratification.md]  verifies: [UC-034]  lane: agent  acc: [doc covers outreach attempt, 30-day comment window, provisional checklist (contract honesty, no oracle special-casing, event completeness), 180-day expiry, and that author-affiliated adapters cannot skip to ratified]  deps: [T15.0]

- [ ] T15.11 Grant provisional provenance to one non-author adapter (LangGraph preferred)  Owner: pool  Est: 90m  kind: human  verifies: [UC-034]  acc: [adapters/langgraph README (or ratification record) shows provenance=provisional with dated checklist reference; dashboard sample or live entry for that adapter uses provenance=provisional; Sire remains unofficial or provisional-only via external review]  deps: [T15.10]  blocked: Needs independent review sign-off (founder or external reviewer)

## Wave 6 -- quality + ship (Wave 33, 3 agents)

- [ ] T15.12 Add docs/manual/v1.0-smoke.md for sign, verify, Opt-in publish, bounty links  Owner: pool  Est: 45m  kind: agent  delivers: [docs/manual/v1.0-smoke.md]  verifies: [UC-032, UC-033]  lane: agent  acc: [smoke doc has copy-paste commands for report --sign, verify-entry, check-dashboard, and local Pages serve, each with expected exit codes]  deps: [T15.6, T15.7, T15.8]

- [ ] T15.13 Full make test and make lint green on clean tree (v1.0 gate)  Owner: pool  Est: 30m  kind: agent  verifies: [infrastructure]  acc: [make test && make lint && bash scripts/check-dashboard.sh all exit 0 on a clean checkout of the integration branch]  deps: [T15.5, T15.6, T15.4]

- [ ] T15.14 Cut v1.0.0 GitHub release with GoReleaser assets  Owner: pool  Est: 45m  kind: human  verifies: [infrastructure]  acc: [git tag v1.0.0 exists; gh release view v1.0.0 shows darwin/linux amd64/arm64 binaries + checksums; release notes cite signed Opt-in (ADR 012/013) and harness bounty]  deps: [T15.13, T15.12, T15.9]  blocked: Founder cuts tag after T15.13 green (same pattern as T14.16)
