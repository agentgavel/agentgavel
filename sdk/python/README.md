# agentgavel-adapter

Python SDK for AgentGavel sidecar adapters.

Agent frameworks implement callbacks on a base `Adapter` class so the Go
harness can score them without importing their dependency trees. The engine
talks to adapters as child processes over a versioned wire protocol; this
package is the callback surface for Python targets.

Install from the repo root:

```bash
pip install -e sdk/python
```

Subclass `agentgavel_adapter.Adapter` and override the session lifecycle
hooks. Transport (JSON-RPC over stdio) ships in a later release.
