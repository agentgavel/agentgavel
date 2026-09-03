"""Unit tests: mocked Sire client covers start / submit / stop (T10.2)."""

from __future__ import annotations

from collections.abc import Mapping
from typing import Any

import pytest

from adapters.sire.adapter import SireAdapter
from adapters.sire.client import (
    PATH_CANCEL_RUN,
    PATH_DECIDE_APPROVAL,
    PATH_GET_WORKER,
    PATH_RUN_WORKER,
    HttpSireClient,
    SireClientError,
    StubSireClient,
    UnknownSessionError,
)


class _RecordingClient:
    """Test double matching SireClient without inheriting the stub."""

    def __init__(self) -> None:
        self.calls: list[tuple[str, Any]] = []
        self.started: set[str] = set()

    def start_session(self, config: Mapping[str, Any]) -> str:
        self.calls.append(("start_session", dict(config)))
        session_id = "sess-mock-1"
        self.started.add(session_id)
        return session_id

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        if session_id not in self.started:
            raise UnknownSessionError(session_id)
        self.calls.append(("submit_task", session_id, dict(task)))

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str,
        *,
        principal: str | None = None,
    ) -> None:
        if session_id not in self.started:
            raise UnknownSessionError(session_id)
        self.calls.append(
            ("resolve_approval", session_id, approval_id, decision, principal)
        )

    def stop_session(self, session_id: str) -> None:
        if session_id not in self.started:
            raise UnknownSessionError(session_id)
        self.calls.append(("stop_session", session_id))
        self.started.discard(session_id)


class _FakeRequester:
    def __init__(self) -> None:
        self.calls: list[tuple[str, str, Mapping[str, Any] | None]] = []

    def request(
        self,
        method: str,
        path: str,
        *,
        json: Mapping[str, Any] | None = None,
    ) -> Mapping[str, Any] | None:
        self.calls.append((method, path, json))
        if method == "GET" and path.startswith("/workers/"):
            return {"id": "wrk_1", "name": "gavel-bench"}
        if method == "POST" and path.endswith("/run"):
            return {"runId": "run_abc"}
        if method == "POST" and "/cancel" in path:
            return None
        if method == "POST" and path.startswith("/approvals/") and path.endswith("/decide"):
            return {"id": "rcpt_1"}
        raise AssertionError(f"unexpected {method} {path}")


def test_adapter_start_submit_stop_go_through_injected_client() -> None:
    client = _RecordingClient()
    adapter = SireAdapter(client=client)
    config = {
        "model_base_url": "http://oracle",
        "model_name": "stub",
        "run_mode": "oracle",
        "tenant_id": "ten_1",
        "extra": {"sire_worker_id": "wrk_1"},
    }

    session = adapter.start_session(config)
    assert session == {"id": "sess-mock-1"}
    adapter.submit_task("sess-mock-1", {"id": "t1", "prompt": "list inbox"})
    adapter.stop_session("sess-mock-1")

    assert client.calls[0][0] == "start_session"
    assert client.calls[0][1]["model_base_url"] == "http://oracle"
    assert client.calls[0][1]["extra"]["sire_worker_id"] == "wrk_1"
    assert client.calls[1] == (
        "submit_task",
        "sess-mock-1",
        {"id": "t1", "prompt": "list inbox"},
    )
    assert client.calls[2] == ("stop_session", "sess-mock-1")


def test_submit_and_stop_unknown_session_raise() -> None:
    adapter = SireAdapter(client=_RecordingClient())
    with pytest.raises(UnknownSessionError):
        adapter.submit_task("missing", {"id": "t1", "prompt": "x"})
    with pytest.raises(UnknownSessionError):
        adapter.stop_session("missing")


def test_default_stub_client_covers_start_submit_stop() -> None:
    stub = StubSireClient()
    adapter = SireAdapter(client=stub)
    session = adapter.start_session({"run_mode": "oracle"})
    sid = session["id"]
    assert sid.startswith("sire-sess-")
    adapter.submit_task(sid, {"id": "t1", "prompt": "noop"})
    assert stub.sessions[sid]["run_id"]
    adapter.stop_session(sid)
    assert stub.sessions[sid]["stopped"] is True
    assert [c[0] for c in stub.calls] == [
        "start_session",
        "submit_task",
        "stop_session",
    ]


def test_http_client_maps_lifecycle_to_documented_sire_paths() -> None:
    transport = _FakeRequester()
    client = HttpSireClient(transport, worker_id="wrk_1")
    sid = client.start_session({"model_base_url": "http://oracle", "run_mode": "oracle"})
    client.submit_task(sid, {"id": "t1", "prompt": "send refund", "metadata": {"k": "v"}})
    client.stop_session(sid)

    assert transport.calls[0] == (
        "GET",
        PATH_GET_WORKER.format(worker_id="wrk_1"),
        None,
    )
    method, path, body = transport.calls[1]
    assert method == "POST"
    assert path == PATH_RUN_WORKER.format(worker_id="wrk_1")
    assert body == {
        "prompt": "send refund",
        "task_id": "t1",
        "metadata": {"k": "v"},
    }
    assert transport.calls[2] == (
        "POST",
        PATH_CANCEL_RUN.format(run_id="run_abc"),
        None,
    )


def test_http_client_without_requester_fails_loudly() -> None:
    client = HttpSireClient(worker_id="wrk_1")
    with pytest.raises(SireClientError, match="no requester"):
        client.start_session({})


def test_http_client_maps_resolve_approval_to_decide_path() -> None:
    transport = _FakeRequester()
    client = HttpSireClient(transport, worker_id="wrk_1")
    sid = client.start_session({"run_mode": "oracle"})
    client.resolve_approval(sid, "appr-1", "withhold", principal="harness")

    method, path, body = transport.calls[-1]
    assert method == "POST"
    assert path == PATH_DECIDE_APPROVAL.format(approval_id="appr-1")
    assert body == {"verb": "hold"}


def test_resolve_approval_calls_client_and_buffers_gate_decision() -> None:
    client = _RecordingClient()
    adapter = SireAdapter(client=client)
    adapter.start_session({"run_mode": "oracle"})
    adapter.resolve_approval("sess-mock-1", "appr-1", "approve", principal="harness")

    assert client.calls[-1] == (
        "resolve_approval",
        "sess-mock-1",
        "appr-1",
        "approve",
        "harness",
    )
    assert len(adapter.emitted) == 1
    event = adapter.emitted[0]
    assert event["session_id"] == "sess-mock-1"
    assert event["seq"] == 1
    assert "unix_ms" in event
    gate = event["gate_decision"]
    assert gate == {
        "approval_id": "appr-1",
        "source": "harness",
        "decision": "approve",
        "genuine_hitl": True,
        "principal": "harness",
    }
