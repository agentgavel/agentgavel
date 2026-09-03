"""Transport tests: Handshake against a Python peer over stdio pipes."""

from __future__ import annotations

import os
import threading

from agentgavel_adapter import Adapter, StdioConn


class _HandshakeAdapter(Adapter):
    """Minimal adapter that only implements Handshake."""

    def handshake(
        self,
        engine_protocol_version: str,
        *,
        engine_version: str | None = None,
    ):
        return {
            "adapter_protocol_version": engine_protocol_version,
            "adapter_name": "fake-python",
            "adapter_version": "0.0.1",
            "provenance": "unofficial",
            "hitl": True,
            "tenancy": False,
            "ledger": False,
            "observability": True,
            "context_mode": "raw",
            "framework_name": "test",
            "framework_version": "0.0.0",
        }


def test_handshake() -> None:
    """Engine peer completes Handshake against an Adapter serving on pipes."""
    engine_r_fd, adapter_w_fd = os.pipe()
    adapter_r_fd, engine_w_fd = os.pipe()

    engine_r = os.fdopen(engine_r_fd, "rb", buffering=0)
    engine_w = os.fdopen(engine_w_fd, "wb", buffering=0)
    adapter_r = os.fdopen(adapter_r_fd, "rb", buffering=0)
    adapter_w = os.fdopen(adapter_w_fd, "wb", buffering=0)

    adapter = _HandshakeAdapter()
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

    try:
        engine = StdioConn(engine_r, engine_w)
        report = engine.handshake("1.0", engine_version="0.0.0-dev")
        assert report["adapter_name"] == "fake-python"
        assert report["adapter_protocol_version"] == "1.0"
        assert report["provenance"] == "unofficial"
        assert report["hitl"] is True
        assert report["observability"] is True
        assert report["context_mode"] == "raw"
    finally:
        engine_w.close()
        engine_r.close()
        thread.join(timeout=5.0)

    assert not thread.is_alive(), "adapter serve loop did not exit after pipe close"
    assert errors == [], f"adapter thread errors: {errors!r}"
