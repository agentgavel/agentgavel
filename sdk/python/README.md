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

Subclass `agentgavel_adapter.Adapter`, override lifecycle hooks (at least
`handshake`), then call `serve()` to run the JSON-RPC 2.0 stdio loop. Push
events with `emit({...})` from framework hooks; the SDK buffers and sends
them as `Event` notifications.

## Examples

- `python -m agentgavel_adapter.examples.fake` — minimal FakeAdapter for Go engine Handshake tests.
