"""Callback dispatch tests: each JSON-RPC method invokes the Adapter hook."""

from __future__ import annotations

import os
import threading
from typing import Any, Mapping

from agentgavel_adapter import (
    METHOD_EXPORT_LEDGER,
    METHOD_RESOLVE_APPROVAL,
    METHOD_START_SESSION,
    METHOD_STOP_SESSION,
    METHOD_SUBMIT_TASK,
    Adapter,
    StdioConn,
)


class _RecordingAdapter(Adapter):
    """Minimal adapter that records lifecycle hook invocations for tests."""

    def __init__(self) -> None:
        super().__init__()
        self.calls: list[tuple[str, tuple[Any, ...], dict[str, Any]]] = []

    def handshake(
        self,
        engine_protocol_version: str,
        *,
        engine_version: str | None = None,
    ) -> Mapping[str, Any]:
        self.calls.append(
            ("handshake", (engine_protocol_version,), {"engine_version": engine_version})
        )
        return {
            "adapter_protocol_version": engine_protocol_version,
            "adapter_name": "recording",
            "adapter_version": "0.0.1",
            "provenance": "unofficial",
            "hitl": True,
            "tenancy": False,
            "ledger": True,
            "observability": True,
            "context_mode": "raw",
        }

    def start_session(self, config: Mapping[str, Any]) -> Mapping[str, Any]:
        self.calls.append(("start_session", (dict(config),), {}))
        return {"id": "sess-test-1"}

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        self.calls.append(("submit_task", (session_id, dict(task)), {}))

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str | int,
        *,
        principal: str | None = None,
    ) -> None:
        self.calls.append(
            (
                "resolve_approval",
                (session_id, approval_id, decision),
                {"principal": principal},
            )
        )

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        self.calls.append(("export_ledger", (session_id,), {}))
        return {
            "session_id": session_id,
            "entries": [{"id": "e1", "kind": "task", "unix_ms": 1}],
        }

    def stop_session(self, session_id: str) -> None:
        self.calls.append(("stop_session", (session_id,), {}))


def _run_with_adapter(adapter: Adapter):
    """Start adapter.serve on pipes; yield (engine StdioConn, join/cleanup)."""
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

    thread = threading.Thread(target=run_adapter, name="adapter-serve", daemon=True)
    thread.start()
    engine = StdioConn(engine_r, engine_w)
    return engine, engine_r, engine_w, thread, errors


def _shutdown(engine_r, engine_w, thread, errors) -> None:
    engine_w.close()
    engine_r.close()
    thread.join(timeout=5.0)
    assert not thread.is_alive(), "adapter serve loop did not exit"
    assert errors == [], f"adapter thread errors: {errors!r}"


def test_start_session_dispatches_callback() -> None:
    adapter = _RecordingAdapter()
    engine, engine_r, engine_w, thread, errors = _run_with_adapter(adapter)
    try:
        result = engine.call(
            METHOD_START_SESSION,
            {
                "model_base_url": "http://oracle",
                "run_mode": "oracle",
                "model_name": "stub",
            },
        )
        assert result == {"id": "sess-test-1"}
        assert adapter.calls[0][0] == "start_session"
        config = adapter.calls[0][1][0]
        assert config["model_base_url"] == "http://oracle"
        assert config["run_mode"] == "oracle"
        assert config["model_name"] == "stub"
    finally:
        _shutdown(engine_r, engine_w, thread, errors)


def test_submit_task_dispatches_callback() -> None:
    adapter = _RecordingAdapter()
    engine, engine_r, engine_w, thread, errors = _run_with_adapter(adapter)
    try:
        result = engine.call(
            METHOD_SUBMIT_TASK,
            {
                "session_id": "sess-test-1",
                "task": {"id": "t1", "prompt": "hi", "metadata": {"k": "v"}},
            },
        )
        assert result == {}
        assert adapter.calls[0][0] == "submit_task"
        session_id, task = adapter.calls[0][1]
        assert session_id == "sess-test-1"
        assert task == {"id": "t1", "prompt": "hi", "metadata": {"k": "v"}}
    finally:
        _shutdown(engine_r, engine_w, thread, errors)


def test_resolve_approval_dispatches_callback() -> None:
    adapter = _RecordingAdapter()
    engine, engine_r, engine_w, thread, errors = _run_with_adapter(adapter)
    try:
        result = engine.call(
            METHOD_RESOLVE_APPROVAL,
            {
                "session_id": "sess-test-1",
                "approval_id": "a1",
                "decision": "deny",
                "principal": "harness",
            },
        )
        assert result == {}
        assert adapter.calls[0][0] == "resolve_approval"
        session_id, approval_id, decision = adapter.calls[0][1]
        assert session_id == "sess-test-1"
        assert approval_id == "a1"
        assert decision == "deny"
        assert adapter.calls[0][2]["principal"] == "harness"
    finally:
        _shutdown(engine_r, engine_w, thread, errors)


def test_export_ledger_dispatches_callback() -> None:
    adapter = _RecordingAdapter()
    engine, engine_r, engine_w, thread, errors = _run_with_adapter(adapter)
    try:
        result = engine.call(METHOD_EXPORT_LEDGER, {"id": "sess-test-1"})
        assert result["session_id"] == "sess-test-1"
        assert result["entries"][0]["id"] == "e1"
        assert adapter.calls[0][0] == "export_ledger"
        assert adapter.calls[0][1][0] == "sess-test-1"
    finally:
        _shutdown(engine_r, engine_w, thread, errors)


def test_stop_session_dispatches_callback() -> None:
    adapter = _RecordingAdapter()
    engine, engine_r, engine_w, thread, errors = _run_with_adapter(adapter)
    try:
        result = engine.call(METHOD_STOP_SESSION, {"id": "sess-test-1"})
        assert result == {}
        assert adapter.calls[0][0] == "stop_session"
        assert adapter.calls[0][1][0] == "sess-test-1"
    finally:
        # StopSession ends the serve loop; join without relying on pipe close alone.
        thread.join(timeout=5.0)
        engine_w.close()
        engine_r.close()
        assert not thread.is_alive(), "adapter serve loop did not exit after StopSession"
        assert errors == [], f"adapter thread errors: {errors!r}"
