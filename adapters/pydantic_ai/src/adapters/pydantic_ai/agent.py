"""Minimal email tool agent for SEC fixtures (T13.18).

Dependency choice: this module is an in-process, pydantic-ai-*compatible* stub.
We deliberately do **not** depend on the ``pydantic-ai`` PyPI package — it pulls
LLM providers and tooling extras unsuitable for CI. The stub still exposes tool
functions (``read_email`` / ``send_email``), accepts a model ``base_url`` aimed
at the Compliance Oracle, and records ``tool_invocation`` events suitable for
AgentGavel observation.
"""

from __future__ import annotations

import json
import time
import urllib.error
import urllib.request
from collections.abc import Callable, Mapping, MutableMapping
from typing import Any

TOOL_READ_EMAIL = "read_email"
TOOL_SEND_EMAIL = "send_email"

HEADER_PROBE_DIRECTIVE = "X-AgentGavel-Probe-Directive"

EventSink = Callable[[MutableMapping[str, Any]], None]

_TOOL_NODES: dict[str, Callable[[Mapping[str, Any]], Mapping[str, Any]]] = {}


def _register(
    name: str,
) -> Callable[[Callable[..., Mapping[str, Any]]], Callable[..., Mapping[str, Any]]]:
    def deco(fn: Callable[..., Mapping[str, Any]]) -> Callable[..., Mapping[str, Any]]:
        _TOOL_NODES[name] = fn
        return fn

    return deco


@_register(TOOL_READ_EMAIL)
def read_email_node(args: Mapping[str, Any]) -> Mapping[str, Any]:
    """SEC fixture tool: read an inbox message."""
    mailbox = str(args.get("mailbox", "inbox"))
    return {
        "mailbox": mailbox,
        "subject": str(args.get("subject", "Weekly digest")),
        "body": str(args.get("body", "Contact list attached.")),
    }


@_register(TOOL_SEND_EMAIL)
def send_email_node(args: Mapping[str, Any]) -> Mapping[str, Any]:
    """SEC fixture tool: send an outbound email."""
    return {
        "to": str(args.get("to", "")),
        "body": str(args.get("body", "")),
        "status": "sent",
    }


def tool_nodes() -> Mapping[str, Callable[[Mapping[str, Any]], Mapping[str, Any]]]:
    """Return the registered tool callables (read-only view)."""
    return dict(_TOOL_NODES)


class MinimalEmailAgent:
    """One-pass agent run: Oracle completion → tool → recorded events."""

    def __init__(
        self,
        *,
        model_base_url: str,
        session_id: str = "agent-sess",
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
        tool_name, arguments, call_id = self._complete_tool_call(prompt, directive, model)
        result = self._invoke_tool(tool_name, arguments, call_id)
        return {
            "tool_name": tool_name,
            "arguments": arguments,
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
        self._record_tool(tool_name, call_id, "before")
        result = node(arguments)
        self._record_tool(tool_name, call_id, "after", outcome="ok")
        return result

    def _record_tool(
        self,
        tool_name: str,
        tool_id: str,
        phase: str,
        *,
        outcome: str | None = None,
    ) -> None:
        self._seq += 1
        inv: dict[str, Any] = {
            "tool_name": tool_name,
            "tool_id": tool_id,
            "phase": phase,
        }
        if outcome is not None:
            inv["outcome"] = outcome
        event: MutableMapping[str, Any] = {
            "session_id": self.session_id,
            "seq": self._seq,
            "unix_ms": int(time.time() * 1000),
            "tool_invocation": inv,
        }
        self.events.append(event)
        if self._on_event is not None:
            self._on_event(event)
