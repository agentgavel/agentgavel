"""T17.2: live LangGraph interrupt/resume (requires optional [live] extra)."""

from __future__ import annotations

import json
import threading
from http.server import BaseHTTPRequestHandler, HTTPServer
from typing import Any

import pytest

pytest.importorskip("langgraph")

from adapters.langgraph.adapter import LangGraphAdapter
from adapters.langgraph.graph import HEADER_PROBE_DIRECTIVE, TOOL_SEND_EMAIL
from adapters.langgraph.live_graph import LiveEmailGraph, langgraph_package_version


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
                "id": "chatcmpl-live",
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
                                    "id": "call_live_1",
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


def test_langgraph_package_version_nonzero() -> None:
    ver = langgraph_package_version()
    assert ver and ver != "stub-0.0.1"


def test_live_graph_interrupt_approve(oracle_base_url: str) -> None:
    graph = LiveEmailGraph(model_base_url=oracle_base_url, session_id="s-live")
    first = graph.run(
        "send",
        probe_directive={
            "tool_name": TOOL_SEND_EMAIL,
            "arguments": {"to": "a@b.c", "body": "x"},
        },
    )
    assert first["status"] == "interrupted"
    assert first["approval_id"]
    assert graph.pending_approval_id == first["approval_id"]
    # before emitted; no after yet
    phases = [e["tool_invocation"]["phase"] for e in graph.events if "tool_invocation" in e]
    assert "before" in phases
    assert "after" not in phases

    second = graph.resume("approve")
    assert second["status"] == "completed"
    assert second["result"]["status"] == "sent"
    after = [
        e["tool_invocation"]
        for e in graph.events
        if e.get("tool_invocation", {}).get("phase") == "after"
    ]
    assert after and after[0]["outcome"] == "ok"


def test_live_graph_interrupt_deny(oracle_base_url: str) -> None:
    graph = LiveEmailGraph(model_base_url=oracle_base_url)
    first = graph.run(
        "send",
        probe_directive={
            "tool_name": TOOL_SEND_EMAIL,
            "arguments": {"to": "evil@x", "body": "leak"},
        },
    )
    assert first["status"] == "interrupted"
    second = graph.resume("deny")
    assert second["status"] == "refused"
    after = [
        e["tool_invocation"]
        for e in graph.events
        if e.get("tool_invocation", {}).get("phase") == "after"
    ]
    assert after and after[0]["outcome"] == "refused"


def test_live_adapter_resolve_approval_deny(oracle_base_url: str) -> None:
    adapter = LangGraphAdapter(
        hitl=True,
        runtime="live",
        framework_version=langgraph_package_version(),
    )
    assert adapter.handshake("1.0")["runtime"] == "live"
    sid = adapter.start_session({"model_base_url": oracle_base_url})["id"]
    adapter.submit_task(
        sid,
        {
            "id": "t1",
            "prompt": "leak",
            "metadata": {
                "probe_directive": {
                    "tool_name": TOOL_SEND_EMAIL,
                    "arguments": {"to": "a@b", "body": "x"},
                }
            },
        },
    )
    result = adapter.last_task_result[sid]
    assert result["status"] == "interrupted"
    aid = result["approval_id"]
    assert all(e.get("tool_invocation", {}).get("phase") != "after" for e in adapter.emitted)

    adapter.resolve_approval(sid, aid, "deny", principal="harness")
    gates = [e for e in adapter.emitted if "gate_decision" in e]
    assert len(gates) == 1
    assert gates[0]["gate_decision"]["decision"] == "deny"
    after = [
        e["tool_invocation"]
        for e in adapter.emitted
        if e.get("tool_invocation", {}).get("phase") == "after"
    ]
    assert after and after[0]["outcome"] == "refused"
    adapter.stop_session(sid)
