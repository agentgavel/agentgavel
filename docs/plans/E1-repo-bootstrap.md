# E1 -- Repository bootstrap

Acceptance: Go module builds `AgentGavel version`; CI workflow exists as a stub that compiles; LICENSE and README state neutrality rules from RFC section 0.
fidelity: executable

## Wave 1

- [x] T1.1 Initialize Go module github.com/agentgavel/gavel and cmd/AgentGavel stub  Owner: pool  Est: 45m  kind: agent  verifies: [infrastructure]  acc: [go build -o /tmp/AgentGavel ./cmd/AgentGavel succeeds and ./tmp/AgentGavel version prints a semver placeholder]  completed: 2026-09-02
  - deps: []
  - Acceptance: `go.mod` present; stub main with `version` and `help` via stdlib `flag`; no third-party CLI libs.
  - Decision rationale: docs/adr/001-go-engine-sidecar-adapters.md, docs/adr/008-repo-module-boundaries.md

- [x] T1.2 Add LICENSE (Apache-2.0), README with RFC abstract and neutrality, CONTRIBUTING stub  Owner: pool  Est: 45m  kind: agent  verifies: [infrastructure]  delivers: [root README + LICENSE]  acc: [README contains Hard vs Soft and unofficial adapter provenance language; LICENSE file exists]  completed: 2026-09-02
  - deps: []
  - Note: content task mixed with engineering; keep marketing claims out.

- [ ] T1.3 Add Makefile targets: build, test, lint, fmt  Owner: pool  Est: 30m  kind: agent  verifies: [infrastructure]  acc: [make build && make test && make fmt are no-ops or green on empty packages]
  - deps: [T1.1]

- [ ] T1.4 Add GitHub Actions CI: go test, go vet, golangci-lint placeholder  Owner: pool  Est: 60m  kind: agent  verifies: [infrastructure]  acc: [.github/workflows/ci.yml runs on PR and executes go test ./...]
  - deps: [T1.1, T1.3]

- [x] T1.5 Create docs/devlog.md seed and .gitignore for binaries, scratch, Python venvs  Owner: pool  Est: 30m  kind: agent  verifies: [infrastructure]  acc: [git check-ignore -q .claude/scratch/foo && docs/devlog.md exists]  completed: 2026-09-02
  - deps: []

- [ ] T1.6 Run gofmt and golangci-lint (or go vet) on bootstrap; fix findings  Owner: pool  Est: 30m  kind: agent  verifies: [infrastructure]  acc: [gofmt -l . is empty; go vet ./... exits 0]
  - deps: [T1.1, T1.3]
