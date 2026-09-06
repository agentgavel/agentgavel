"""T15.26: Hermes Events — tool_invocation / gate_decision emission + N/A paths."""

from __future__ import annotations

import json
import os
import select
import threading
import time
from typing import Any

import pytest

from agentgavel_adapter import METHOD_EVENT_NOTIFY, METHOD_START_SESSION, StdioConn

from adapters.hermes.adapter import HermesAdapter, HitlNotSupportedError
from adapters.hermes.events import (
    assert_tool_invocation_order,
    build_gate_decision,
    map_hermes_frame,
    normalize_decision,
)


def test_handshake_observability_true_ledger_false() -> None:
    report = HermesAdapter().handshake("1.0")
    assert report["observability"] is True
    assert report["ledger"] is False
    assert report["hitl"] is False


def test_ingest_tool_started_completed_emits_ordered_tool_invocation() -> None:
    adapter = HermesAdapter()
    sid = adapter.start_session({})["id"]
    adapter.ingest_hermes_event(
        sid,
        {
            "type": "tool.started",
            "tool_name": "send_email",
            "tool_id": "call-1",
            "arguments": {"to": "a@b.c"},
        },
    )
    adapter.ingest_hermes_event(
        sid,
        {
            "type": "tool.completed",
            "tool_name": "send_email",
            "tool_id": "call-1",
            "outcome": "sent",
        },
    )
    tool_ev = [e for e in adapter.emitted if "tool_invocation" in e]
    assert len(tool_ev) == 2
    assert tool_ev[0]["seq"] == 1
    assert tool_ev[1]["seq"] == 2
    assert tool_ev[0]["tool_invocation"]["phase"] == "before"
    assert tool_ev[1]["tool_invocation"]["phase"] == "after"
    assert tool_ev[0]["tool_invocation"]["tool_name"] == "send_email"
    assert json.loads(tool_ev[0]["tool_invocation"]["arguments_json"]) == {"to": "a@b.c"}
    assert tool_ev[1]["tool_invocation"]["outcome"] == "sent"
    adapter.assert_emitted_tool_order()
    assert_tool_invocation_order(adapter.emitted)


def test_ingest_approval_resolved_emits_gate_decision() -> None:
    adapter = HermesAdapter()
    sid = adapter.start_session({})["id"]
    adapter.ingest_hermes_event(
        sid,
        {
            "type": "approval.resolved",
            "approval_id": "appr-9",
            "decision": "deny",
            "source": "store",
            "genuine_hitl": False,
        },
    )
    gates = [e for e in adapter.emitted if "gate_decision" in e]
    assert len(gates) == 1
    gate = gates[0]["gate_decision"]
    assert gate["approval_id"] == "appr-9"
    assert gate["decision"] == "deny"
    assert gate["source"] == "store"
    assert gate["genuine_hitl"] is False


def test_ingest_unknown_frame_is_na_noop() -> None:
    adapter = HermesAdapter()
    sid = adapter.start_session({})["id"]
    recorded = adapter.ingest_hermes_event(
        sid,
        {"type": "assistant.delta", "text": "hello"},
    )
    assert recorded == []
    assert adapter.emitted == []


def test_resolve_approval_still_na_when_hitl_false() -> None:
    """N/A path: hitl=false → no gate_decision from ResolveApproval (T15.25)."""
    adapter = HermesAdapter()
    sid = adapter.start_session({})["id"]
    with pytest.raises(HitlNotSupportedError):
        adapter.resolve_approval(sid, "appr-1", "approve")
    assert not any("gate_decision" in e for e in adapter.emitted)


def test_map_hermes_frame_tool_and_gate_shapes() -> None:
    started = map_hermes_frame({"event": "hermes.tool.progress", "name": "terminal", "id": "t1"})
    assert started[0]["tool_invocation"]["phase"] == "before"
    done = map_hermes_frame(
        {"type": "tool.completed", "data": {"tool_name": "terminal", "tool_id": "t1", "refused": True}}
    )
    assert done[0]["tool_invocation"]["phase"] == "after"
    assert done[0]["tool_invocation"]["refused"] is True
    gate = map_hermes_frame(
        {"type": "approval.respond", "approval_id": "a1", "decision": "allow_once"}
    )
    assert gate[0]["gate_decision"]["decision"] == "approve"
    assert map_hermes_frame({"type": "run.completed"}) == []


def test_normalize_decision_hermes_synonyms() -> None:
    assert normalize_decision("allow_always") == "approve"
    assert normalize_decision("reject") == "deny"
    assert normalize_decision(2) == "deny"
    assert build_gate_decision("x", "withhold")["decision"] == "withhold"


def test_ingest_unknown_session_raises() -> None:
    adapter = HermesAdapter()
    with pytest.raises(KeyError):
        adapter.ingest_hermes_event("missing", {"type": "tool.started", "tool_name": "x"})


def _run_with_adapter(adapter: HermesAdapter):
    engine_r_fd, adapter_w_fd = os.pipe()
    adapter_r_fd, engine_w_fd = os.pipe()

    engine_r = os.fdopen(engine_r_fd, "rb", buffering=0)
    engine_w = os.fdopen(engine_w_fd, "wb", buffering=0)
    adapter_r = os.fdopen(adapter_r_fd, "rb", buffering=0)
    adapter_w = os.fdopen(adapter_w_fd, "wb", buffering=0)

    errors: list[BaseException] = []

    def run_adapter() -> None:
        try:
            adapter.serve(reader=adapter_r, writer=adapter_w)
        except BaseException as exc:  # noqa: BLE001 — surface in main thread
            errors.append(exc)
        finally:
            adapter_r.close()
            adapter_w.close()

    thread = threading.Thread(target=run_adapter, name="hermes-adapter-serve", daemon=True)
    thread.start()
    engine = StdioConn(engine_r, engine_w)
    return engine, engine_r, engine_w, thread, errors


def _shutdown(engine_r, engine_w, thread, errors) -> None:
    engine_w.close()
    engine_r.close()
    thread.join(timeout=5.0)
    assert not thread.is_alive(), "adapter serve loop did not exit"
    assert errors == [], f"adapter thread errors: {errors!r}"


def _drain_event_notifies(engine_r, *, timeout: float = 1.0) -> list[dict[str, Any]]:
    events: list[dict[str, Any]] = []
    deadline = time.monotonic() + timeout
    while time.monotonic() < deadline:
        remaining = deadline - time.monotonic()
        ready, _, _ = select.select([engine_r], [], [], remaining)
        if not ready:
            break
        line = engine_r.readline()
        if not line:
            break
        msg = json.loads(line)
        if msg.get("method") == METHOD_EVENT_NOTIFY:
            params = msg.get("params")
            if isinstance(params, dict):
                events.append(params)
        deadline = time.monotonic() + 0.05
    return events


def test_ingest_emits_tool_invocation_over_transport() -> None:
    adapter = HermesAdapter()
    engine, engine_r, engine_w, thread, errors = _run_with_adapter(adapter)
    try:
        started = engine.call(METHOD_START_SESSION, {})
        sid = started["id"]
        assert isinstance(sid, str) and sid
        adapter.ingest_hermes_event(
            sid,
            {"type": "tool.started", "tool_name": "read_email", "tool_id": "c1"},
        )
        adapter.ingest_hermes_event(
            sid,
            {"type": "tool.completed", "tool_name": "read_email", "tool_id": "c1", "outcome": "ok"},
        )
        events = _drain_event_notifies(engine_r)
        tools = [e for e in events if "tool_invocation" in e]
        assert len(tools) == 2
        assert tools[0]["tool_invocation"]["phase"] == "before"
        assert tools[1]["tool_invocation"]["phase"] == "after"
        assert_tool_invocation_order(events)
    finally:
        _shutdown(engine_r, engine_w, thread, errors)
