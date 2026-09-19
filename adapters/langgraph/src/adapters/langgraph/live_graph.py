"""Live LangGraph email graph (T17.2).

Requires the optional ``langgraph`` package (``pip install
'agentgavel-adapter-langgraph[live]'``). Uses real ``StateGraph`` +
``interrupt()`` / ``Command(resume=...)`` with an in-memory checkpointer.

The stub graph in :mod:`adapters.langgraph.graph` remains the CI default.
"""

from __future__ import annotations

import json
import urllib.error
import urllib.request
from collections.abc import Callable, Mapping, MutableMapping
from typing import Any, Literal, TypedDict
from uuid import uuid4

from adapters.langgraph.attestation import context_attestation_payload
from adapters.langgraph.events import build_tool_invocation, make_event
from adapters.langgraph.graph import (
    HEADER_PROBE_DIRECTIVE,
    TOOL_READ_EMAIL,
    invoke_tool_node,
)
from adapters.langgraph.interrupt import DEFAULT_GATED_TOOLS

EventSink = Callable[[MutableMapping[str, Any]], None]


class LiveGraphState(TypedDict, total=False):
    prompt: str
    probe_directive: dict[str, Any]
    model: str
    tool_name: str
    arguments: dict[str, Any]
    call_id: str
    approval_id: str
    gated: bool
    decision: str
    status: str
    result: dict[str, Any]


def require_langgraph() -> None:
    """Fail closed when the live package is missing."""
    try:
        import langgraph  # noqa: F401
    except ImportError as exc:
        raise RuntimeError(
            "langgraph package required for runtime=live; "
            "install with: pip install 'agentgavel-adapter-langgraph[live]'"
        ) from exc


def langgraph_package_version() -> str:
    """Return installed langgraph version or raise if missing."""
    require_langgraph()
    from importlib.metadata import PackageNotFoundError, version

    try:
        return version("langgraph")
    except PackageNotFoundError as exc:
        raise RuntimeError("langgraph installed but metadata missing") from exc


class LiveEmailGraph:
    """Oracle → tool graph built on real LangGraph interrupt/resume."""

    def __init__(
        self,
        *,
        model_base_url: str,
        session_id: str = "live-sess",
        on_event: EventSink | None = None,
        gated_tools: frozenset[str] | None = None,
        thread_id: str | None = None,
    ) -> None:
        require_langgraph()
        if not model_base_url or not str(model_base_url).strip():
            raise ValueError("model_base_url is required")
        self.model_base_url = str(model_base_url).rstrip("/")
        self.session_id = session_id
        self._on_event = on_event
        self._gated = gated_tools if gated_tools is not None else DEFAULT_GATED_TOOLS
        self.thread_id = thread_id or f"lg-thread-{uuid4().hex[:12]}"
        self.events: list[MutableMapping[str, Any]] = []
        self._seq = 0
        self._graph = self._compile()
        self._config: dict[str, Any] = {"configurable": {"thread_id": self.thread_id}}
        self._pending_approval: str | None = None

    def run(
        self,
        prompt: str,
        *,
        probe_directive: Mapping[str, Any] | None = None,
        model: str = "oracle",
    ) -> Mapping[str, Any]:
        """Invoke until complete or interrupted for HITL."""
        directive = (
            dict(probe_directive)
            if probe_directive is not None
            else {
                "tool_name": TOOL_READ_EMAIL,
                "arguments": {"mailbox": "inbox"},
            }
        )
        self._record_context_attestation(prompt)
        result = self._graph.invoke(
            {
                "prompt": prompt,
                "probe_directive": directive,
                "model": model,
            },
            self._config,
        )
        return self._normalize_result(result)

    def resume(self, decision: str) -> Mapping[str, Any]:
        """Resume after ResolveApproval with approve|deny|withhold."""
        from langgraph.types import Command

        if self._pending_approval is None:
            raise RuntimeError("no pending live interrupt to resume")
        result = self._graph.invoke(Command(resume=decision), self._config)
        self._pending_approval = None
        return self._normalize_result(result)

    @property
    def pending_approval_id(self) -> str | None:
        return self._pending_approval

    def _compile(self) -> Any:
        from langgraph.checkpoint.memory import InMemorySaver
        from langgraph.graph import END, START, StateGraph
        from langgraph.types import interrupt

        gated_tools = self._gated
        complete = self._complete_tool_call
        record_tool = self._record_tool

        def oracle_node(state: LiveGraphState) -> dict[str, Any]:
            tool_name, arguments, call_id = complete(
                state["prompt"],
                state.get("probe_directive") or {},
                state.get("model") or "oracle",
            )
            approval_id = f"lg-appr-{uuid4().hex[:12]}"
            is_gated = tool_name in gated_tools
            # Intent observed before any side effect (and before interrupt pause).
            record_tool(tool_name, call_id, "before", arguments=arguments)
            return {
                "tool_name": tool_name,
                "arguments": arguments,
                "call_id": call_id,
                "approval_id": approval_id,
                "gated": is_gated,
            }

        def route_after_oracle(state: LiveGraphState) -> Literal["hitl_gate", "run_tool"]:
            return "hitl_gate" if state.get("gated") else "run_tool"

        def hitl_gate(state: LiveGraphState) -> dict[str, Any]:
            # interrupt lives alone in this node so oracle_node (before event)
            # does not re-run on Command(resume=...).
            decision = interrupt(
                {
                    "approval_id": state["approval_id"],
                    "tool_name": state["tool_name"],
                    "arguments": dict(state.get("arguments") or {}),
                    "call_id": state["call_id"],
                    "session_id": self.session_id,
                }
            )
            return {"decision": str(decision).strip().lower()}

        def run_tool(state: LiveGraphState) -> dict[str, Any]:
            tool_name = state["tool_name"]
            arguments = dict(state.get("arguments") or {})
            call_id = state["call_id"]
            decision = str(state.get("decision") or "approve").strip().lower()

            if state.get("gated") and decision != "approve":
                record_tool(tool_name, call_id, "after", outcome="refused")
                return {
                    "status": "refused",
                    "result": {"refused": True, "decision": decision},
                }

            result = invoke_tool_node(tool_name, arguments)
            record_tool(tool_name, call_id, "after", outcome="ok")
            return {"status": "completed", "result": dict(result)}

        builder = StateGraph(LiveGraphState)
        builder.add_node("oracle", oracle_node)
        builder.add_node("hitl_gate", hitl_gate)
        builder.add_node("run_tool", run_tool)
        builder.add_edge(START, "oracle")
        builder.add_conditional_edges(
            "oracle",
            route_after_oracle,
            {"hitl_gate": "hitl_gate", "run_tool": "run_tool"},
        )
        builder.add_edge("hitl_gate", "run_tool")
        builder.add_edge("run_tool", END)
        return builder.compile(checkpointer=InMemorySaver())

    def _normalize_result(self, result: Mapping[str, Any] | Any) -> Mapping[str, Any]:
        interrupts = None
        if isinstance(result, Mapping):
            interrupts = result.get("__interrupt__")
        if interrupts:
            first = interrupts[0] if isinstance(interrupts, (list, tuple)) else interrupts
            value = getattr(first, "value", first)
            if not isinstance(value, Mapping):
                value = {}
            approval_id = str(value.get("approval_id") or "")
            self._pending_approval = approval_id or None
            return {
                "tool_name": value.get("tool_name"),
                "arguments": value.get("arguments") or {},
                "status": "interrupted",
                "approval_id": approval_id,
                "call_id": value.get("call_id"),
                "events": list(self.events),
            }

        if not isinstance(result, Mapping):
            result = {}
        status = str(result.get("status") or "completed")
        return {
            "tool_name": result.get("tool_name"),
            "arguments": result.get("arguments") or {},
            "status": status,
            "result": result.get("result"),
            "approval_id": result.get("approval_id"),
            "call_id": result.get("call_id"),
            "decision": result.get("decision"),
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
                HEADER_PROBE_DIRECTIVE: json.dumps(dict(directive)),
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
        if "tool_name" in directive and directive["tool_name"]:
            tool_name = str(directive["tool_name"])
            if isinstance(directive.get("arguments"), dict):
                arguments = dict(directive["arguments"])
        return tool_name, arguments, call_id

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
    ) -> None:
        self._seq += 1
        inv = build_tool_invocation(
            tool_name,
            tool_id,
            phase,
            arguments=arguments,
            outcome=outcome,
            refused=(outcome == "refused"),
        )
        event = make_event(self.session_id, self._seq, tool_invocation_payload=inv)
        self.events.append(event)
        if self._on_event is not None:
            self._on_event(event)
