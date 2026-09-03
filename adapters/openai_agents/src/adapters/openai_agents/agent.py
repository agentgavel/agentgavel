"""Minimal email-tool agent for SEC fixtures (T13.17).

Dependency choice: this module is an in-process, OpenAI-Agents-*compatible*
stub. We deliberately do **not** depend on the ``openai-agents`` PyPI package —
it pulls a heavier tool/LLM stack than CI needs. The stub still exposes tool
callables (``read_email`` / ``send_email``), accepts a model ``base_url`` aimed
at the Compliance Oracle, and records protocol events (context attestation +
``tool_invocation`` before/after) suitable for AgentGavel observation.
"""

from __future__ import annotations

import json
import urllib.error
import urllib.request
from collections.abc import Callable, Mapping, MutableMapping
from typing import Any

from adapters.openai_agents.attestation import context_attestation_payload
from adapters.openai_agents.events import build_tool_invocation, make_event

TOOL_READ_EMAIL = "read_email"
TOOL_SEND_EMAIL = "send_email"

HEADER_PROBE_DIRECTIVE = "X-AgentGavel-Probe-Directive"

EventSink = Callable[[MutableMapping[str, Any]], None]

_TOOL_NODES: dict[str, Callable[[Mapping[str, Any]], Mapping[str, Any]]] = {}

_ToolNode = Callable[[Mapping[str, Any]], Mapping[str, Any]]


def _register(name: str) -> Callable[[_ToolNode], _ToolNode]:
    def deco(fn: _ToolNode) -> _ToolNode:
        _TOOL_NODES[name] = fn
        return fn

    return deco


@_register(TOOL_READ_EMAIL)
def read_email_tool(args: Mapping[str, Any]) -> Mapping[str, Any]:
    """SEC fixture tool: read an inbox message."""
    mailbox = str(args.get("mailbox", "inbox"))
    return {
        "mailbox": mailbox,
        "subject": str(args.get("subject", "Weekly digest")),
        "body": str(args.get("body", "Contact list attached.")),
    }


@_register(TOOL_SEND_EMAIL)
def send_email_tool(args: Mapping[str, Any]) -> Mapping[str, Any]:
    """SEC fixture tool: send an outbound email."""
    return {
        "to": str(args.get("to", "")),
        "body": str(args.get("body", "")),
        "status": "sent",
    }


def tool_nodes() -> Mapping[str, Callable[[Mapping[str, Any]], Mapping[str, Any]]]:
    """Return the registered tool callables (read-only view)."""
    return dict(_TOOL_NODES)


def invoke_tool(tool_name: str, arguments: Mapping[str, Any]) -> Mapping[str, Any]:
    """Execute a registered tool (used by the agent loop)."""
    node = _TOOL_NODES.get(tool_name)
    if node is None:
        raise KeyError(f"unknown tool: {tool_name!r}")
    return node(arguments)


class MinimalEmailAgent:
    """One-pass agent: Oracle completion → tool → recorded events."""

    def __init__(
        self,
        *,
        model_base_url: str,
        session_id: str = "oa-sess",
        on_event: EventSink | None = None,
    ) -> None:
        if not model_base_url or not str(model_base_url).strip():
            raise ValueError("model_base_url is required")
        self.model_base_url = str(model_base_url).rstrip("/")
        self.session_id = session_id
        self._on_event = on_event
        self.events: list[MutableMapping[str, Any]] = []
        self._seq = 0

    def run(
        self,
        prompt: str,
        *,
        probe_directive: Mapping[str, Any] | None = None,
        model: str = "oracle",
    ) -> Mapping[str, Any]:
        """Run once against the Oracle at ``model_base_url``."""
        directive = (
            dict(probe_directive)
            if probe_directive is not None
            else {
                "tool_name": TOOL_READ_EMAIL,
                "arguments": {"mailbox": "inbox"},
            }
        )
        # ADR 005: attest the prompt before tool dispatch (hosted-safe).
        self._record_context_attestation(prompt)
        tool_name, arguments, call_id = self._complete_tool_call(prompt, directive, model)
        result = self._invoke_tool(tool_name, arguments, call_id)
        return {
            "tool_name": tool_name,
            "arguments": arguments,
            "status": "completed",
            "result": result,
            "events": list(self.events),
        }

    def _complete_tool_call(
        self,
        prompt: str,
        directive: Mapping[str, Any],
        model: str,
    ) -> tuple[str, dict[str, Any], str]:
        url = f"{self.model_base_url}/v1/chat/completions"
        body = json.dumps(
            {
                "model": model,
                "messages": [{"role": "user", "content": prompt}],
            }
        ).encode("utf-8")
        req = urllib.request.Request(
            url,
            data=body,
            method="POST",
            headers={
                "Content-Type": "application/json",
                HEADER_PROBE_DIRECTIVE: json.dumps(directive),
            },
        )
        try:
            with urllib.request.urlopen(req, timeout=10) as resp:
                raw = resp.read().decode("utf-8")
        except urllib.error.HTTPError as exc:
            detail = exc.read().decode("utf-8", errors="replace")
            raise RuntimeError(f"oracle HTTP {exc.code} at {url}: {detail}") from exc
        except urllib.error.URLError as exc:
            raise RuntimeError(f"oracle unreachable at {url}: {exc}") from exc

        payload = json.loads(raw)
        try:
            message = payload["choices"][0]["message"]
            tool_call = message["tool_calls"][0]
            call_id = str(tool_call.get("id") or "call_oracle_1")
            fn = tool_call["function"]
            tool_name = str(fn["name"])
            arguments = json.loads(fn.get("arguments") or "{}")
        except (KeyError, IndexError, TypeError, json.JSONDecodeError) as exc:
            raise RuntimeError(f"oracle response missing tool_calls: {payload!r}") from exc
        if not isinstance(arguments, dict):
            raise RuntimeError("oracle tool arguments must be a JSON object")
        return tool_name, arguments, call_id

    def _invoke_tool(
        self,
        tool_name: str,
        arguments: Mapping[str, Any],
        call_id: str,
    ) -> Mapping[str, Any]:
        node = _TOOL_NODES.get(tool_name)
        if node is None:
            raise KeyError(f"unknown tool: {tool_name!r}")
        self._record_tool(tool_name, call_id, "before", arguments=arguments)
        try:
            result = node(arguments)
        except Exception as exc:
            self._record_tool(
                tool_name,
                call_id,
                "after",
                outcome="error",
                error=str(exc),
            )
            raise
        self._record_tool(tool_name, call_id, "after", outcome="ok")
        return result

    def _record_context_attestation(self, text: str) -> None:
        if not text:
            return
        self._seq += 1
        event = make_event(
            self.session_id,
            self._seq,
            context_attestation_payload=context_attestation_payload(text),
        )
        self.events.append(event)
        if self._on_event is not None:
            self._on_event(event)

    def _record_tool(
        self,
        tool_name: str,
        tool_id: str,
        phase: str,
        *,
        arguments: Mapping[str, Any] | None = None,
        outcome: str | None = None,
        error: str | None = None,
    ) -> None:
        self._seq += 1
        inv = build_tool_invocation(
            tool_name,
            tool_id,
            phase,
            arguments=arguments if phase == "before" else None,
            outcome=outcome,
            error=error,
        )
        event = make_event(
            self.session_id,
            self._seq,
            tool_invocation_payload=inv,
        )
        self.events.append(event)
        if self._on_event is not None:
            self._on_event(event)
