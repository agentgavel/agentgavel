# rubber-stamp

Fast standalone check that an LLM-context claim cannot satisfy your
approval gate. Runs **SEC-002** and **SEC-006** only. Verdict contract:
[ADR 011](../adr/011-rubber-stamp-verdict.md) (UC-021).

Run from the repository root unless noted.

## What it checks

| Scenario | Threat |
| --- | --- |
| SEC-002 | Approval forgery (rubber-stamp via tool output) |
| SEC-006 | Timeout auto-approve |

It does not run the rest of the security suite. Use `AgentGavel run` for
full-suite scoring.

## Usage

Prerequisites: Go toolchain; from repo root set `export GOWORK=off` if a
parent `go.work` lists unrelated modules.

```bash
make build
# or: go build -ldflags "-X main.version=$(tr -d '[:space:]' < VERSION)" -o AgentGavel ./cmd/AgentGavel
./AgentGavel version
```

Flags (no `--suite` or `--scenarios`; the filter is fixed):

| Flag | Meaning |
| --- | --- |
| `--adapter` | Adapter command or path (required); may include args |
| `--out` / `--root` | Results root (writes `results/<run-id>/`) |
| `--run-id` | Run id under `results/` (default: generated) |
| `--seeds` | Deterministic seed count (default: 25) |
| `--fingerprint` | Reload pins from a prior `fingerprint.json` or `summary.json` |

Stdout is the absolute path to `summary.json` only (same shape as
`run --ci`). No GSI or grade is printed.

### FakeAdapter (harness proof)

Copy-paste from repo root after `make build`:

```bash
go build -o /tmp/fakeadapter ./internal/engine/testdata/fakeadapter
./AgentGavel rubber-stamp \
  --adapter /tmp/fakeadapter \
  --seeds 25 \
  --out /tmp/ag-rubber-stamp \
  --run-id smoke-rubber-stamp
echo "exit=$?"
# stdout: /tmp/ag-rubber-stamp/results/smoke-rubber-stamp/summary.json
```

Expect `summary.json` scenario keys exactly `SEC-002` and `SEC-006`.

Automated mirror:

```bash
go test ./cmd/AgentGavel -count=1 -run RubberStamp
```

## Exit codes

Same map as `run --ci` (pass→0, fail→1, catastrophic→2), with one
fail-closed rule for this gate ([ADR 011](../adr/011-rubber-stamp-verdict.md)):

| Outcome | Exit code |
| --- | --- |
| All scored scenarios pass | `0` |
| Any Fail | `1` |
| Catastrophic | `2` (wins over Fail if both) |
| Both SEC-002 and SEC-006 are N/A | `1` |

When every selected scenario is N/A (for example `hitl=false`), the
command exits `1` and writes a single stderr line beginning
`rubber-stamp: not_applicable` that names the missing capability. The
`summary.json` still records the N/A rows so the reason is inspectable.

This differs from `run --ci`, which exits `0` on an all-N/A suite:
`run` reports what it observed; `rubber-stamp` certifies one property and
must fail closed when it cannot observe it. Frameworks without HITL
cannot pass `rubber-stamp` (RFC §3.2).

Usage / flag validation errors also exit `2` before scoring.

## CI snippet

GitHub Actions step (after checkout and a Go build of the CLI):

```yaml
- name: rubber-stamp (SEC-002 + SEC-006)
  env:
    GOWORK: "off"
  run: |
    go build -o AgentGavel ./cmd/AgentGavel
    go build -o /tmp/fakeadapter ./internal/engine/testdata/fakeadapter
    ./AgentGavel rubber-stamp \
      --adapter /tmp/fakeadapter \
      --seeds 25 \
      --out "$RUNNER_TEMP/ag-rubber-stamp" \
      --run-id ci-rubber-stamp
```

Replace `/tmp/fakeadapter` with your real adapter command when gating a
product. Non-zero exit fails the job; treat that as the gate.

## What it does not do

- Does **not** compute or print a GSI.
- Does **not** assign or print a grade.
- Is **not** a substitute for `AgentGavel run` (full suite / scorecard).
  Use `run` (and `report`) when you need the full security suite or a
  rendered scorecard; use `rubber-stamp` only as the fast approval-gate
  check.
