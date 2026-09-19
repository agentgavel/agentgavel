"""Tests for env-gated live Sire client bootstrap (T16.8 / T16.9)."""

from __future__ import annotations

from collections.abc import Mapping
from typing import Any

import pytest

from adapters.sire.adapter import SireAdapter
from adapters.sire.client import (
    HttpSireClient,
    StubSireClient,
    client_from_env,
)


class _FakeRequester:
    def __init__(self, worker: Mapping[str, Any] | None = None) -> None:
        self.worker = dict(worker or {"id": "w1", "version": "sire-test-1.2.3"})
        self.calls: list[tuple[str, str]] = []

    def request(
        self,
        method: str,
        path: str,
        *,
        json: Mapping[str, Any] | None = None,
    ) -> Mapping[str, Any] | None:
        del json
        self.calls.append((method.upper(), path))
        if method.upper() == "GET" and path.startswith("/workers/"):
            return self.worker
        if method.upper() == "POST" and path.endswith("/run"):
            return {"runId": "run-1"}
        if method.upper() == "GET" and path == "/compliance/receipts":
            return []
        if method.upper() == "POST" and "/approvals/" in path:
            return {"ok": True}
        if method.upper() == "POST" and path.startswith("/runs/"):
            return {"ok": True}
        return {}


def test_client_from_env_defaults_to_stub(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.delenv("AGENTGAVEL_SIRE_TOKEN", raising=False)
    monkeypatch.delenv("SIRE_API_TOKEN", raising=False)
    monkeypatch.delenv("AGENTGAVEL_SIRE_WORKER_ID", raising=False)
    monkeypatch.delenv("SIRE_WORKER_ID", raising=False)
    client = client_from_env()
    assert isinstance(client, StubSireClient)
    assert client.runtime_class == "stub"


def test_client_from_env_live_when_token_and_worker(
    monkeypatch: pytest.MonkeyPatch,
) -> None:
    monkeypatch.setenv("AGENTGAVEL_SIRE_TOKEN", "tok")
    monkeypatch.setenv("AGENTGAVEL_SIRE_WORKER_ID", "worker-1")
    monkeypatch.setenv("AGENTGAVEL_SIRE_API_BASE", "https://example.test/api/v1")
    client = client_from_env()
    assert isinstance(client, HttpSireClient)
    assert client.runtime_class == "live"
    assert client._worker_id == "worker-1"
    assert client.api_base.endswith("/api/v1")


def test_handshake_stub_runtime() -> None:
    report = SireAdapter().handshake("1.0")
    assert report["runtime"] == "stub"
    assert report["ledger"] is False
    assert report["observability"] is False
    assert report["framework_version"] == "stub"


def test_handshake_live_runtime_and_probe() -> None:
    requester = _FakeRequester()
    client = HttpSireClient(requester, worker_id="worker-1")
    adapter = SireAdapter(client=client)
    report = adapter.handshake("1.0")
    assert report["runtime"] == "live"
    assert report["ledger"] is False
    assert report["observability"] is False
    assert report["provenance"] == "unofficial"
    # version probes on start_session, not handshake — still live-unprobed here
    assert report["framework_version"] == "live-unprobed"
    sid = adapter.start_session({"extra": {"sire_worker_id": "worker-1"}})["id"]
    assert sid
    assert client.framework_version == "sire-test-1.2.3"
    # Re-handshake would still show cached client version if we re-read getattr
    report2 = adapter.handshake("1.0")
    assert report2["framework_version"] == "sire-test-1.2.3"
