# AgentGavel

AgentGavel is an open-source benchmarking harness that measures how well AI
agent frameworks enforce governance, security, and reliability guarantees
under adversarial pressure. It evaluates whether the runtime control plane
holds when actively attacked: whether policy ceilings survive prompt
injection, whether human-in-the-loop (HITL) gates actually stop side effects,
and whether every consequential action produces tamper-evident provenance.

Specification: [docs/RFC-0001.md](docs/RFC-0001.md)

## Hard vs Soft governance

AgentGavel distinguishes **Hard** governance (deterministic chokepoints that
refuse unsafe actions regardless of LLM output) from **Soft** governance
(system prompts that an LLM may be persuaded to ignore). Both are scored.
Hard vs Soft classification uses a Compliance Oracle method that does not
depend on any particular model's behavior (RFC section 4.12).

## Adapter provenance

Adapters carry a provenance label on every scorecard:

- **ratified** -- the target framework's maintainers reviewed or contributed
  the adapter
- **provisional** -- independent review after outreach (see
  [docs/adr/007-adapter-ratification.md](docs/adr/007-adapter-ratification.md))
- **unofficial** -- otherwise, including author-affiliated adapters until
  external sign-off

A low score behind an **unofficial** adapter is a claim about the adapter as
much as the framework.

## Status

v0.1 is under construction. See [docs/plan.md](docs/plan.md) and
[docs/roadmap.md](docs/roadmap.md).
Local smoke: [docs/manual/v0.1-smoke.md](docs/manual/v0.1-smoke.md).

## Releases

Static `AgentGavel` binaries (linux/darwin, amd64/arm64) are published via
GoReleaser when a version tag is pushed:

```bash
git tag v0.1.0
git push origin v0.1.0
```

That triggers [`.github/workflows/release.yml`](.github/workflows/release.yml),
which builds archives and a checksum file and attaches them to the GitHub
Release. Local dry-run: `goreleaser release --snapshot --clean`.

## License

Apache License 2.0. See [LICENSE](LICENSE).
