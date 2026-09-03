"""Integration: ResolveApproval over JSON-RPC emits a gate_decision Event."""

from __future__ import annotations

import json
import os
import select
import threading
import time
from typing import Any

from agentgavel_adapter import (
    METHOD_EVENT_NOTIFY,
    METHOD_RESOLVE_APPROVAL,
    METHOD_START_SESSION,
    StdioConn,
)

from adapters.sire.adapter import SireAdapter
from adapters.sire.client import StubSireClient


def _run_with_adapter(adapter: SireAdapter):
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

    thread = threading.Thread(target=run_adapter, name="sire-adapter-serve", daemon=True)
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
    """Read Event notifications left on the pipe after Call returns."""
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


def test_resolve_approval_emits_gate_decision_over_transport() -> None:
    client = StubSireClient()
    adapter = SireAdapter(client=client)
    engine, engine_r, engine_w, thread, errors = _run_with_adapter(adapter)
    try:
        started = engine.call(METHOD_START_SESSION, {"run_mode": "oracle"})
        sid = started["id"]
        assert isinstance(sid, str) and sid
        result = engine.call(
            METHOD_RESOLVE_APPROVAL,
            {
                "session_id": sid,
                "approval_id": "appr-1",
                "decision": "deny",
                "principal": "harness",
            },
        )
        assert result == {}
        events = _drain_event_notifies(engine_r)
        assert events, "expected Event notify with gate_decision"
        event = events[0]
        assert event["session_id"] == sid
        assert event["seq"] == 1
        gate = event["gate_decision"]
        assert gate["approval_id"] == "appr-1"
        assert gate["source"] == "harness"
        assert gate["decision"] == "deny"
        assert gate["genuine_hitl"] is True
        assert gate["principal"] == "harness"
        assert ("resolve_approval", (sid, "appr-1", "deny", "harness")) in client.calls
        assert adapter.emitted[0]["gate_decision"]["decision"] == "deny"
    finally:
        _shutdown(engine_r, engine_w, thread, errors)
