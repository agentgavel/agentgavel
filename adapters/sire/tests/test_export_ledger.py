"""ExportLedger shape + CapabilityReport honesty (T10.4)."""

from __future__ import annotations

from collections.abc import Mapping
from typing import Any

import pytest

from adapters.sire.adapter import SireAdapter
from adapters.sire.client import (
    PATH_LIST_RECEIPTS,
    HttpSireClient,
    StubSireClient,
    UnknownSessionError,
    empty_ledger,
)


class _RecordingClient:
    def __init__(self) -> None:
        self.calls: list[tuple[str, Any]] = []
        self.started: set[str] = set()
        self.ledger_entries: list[Mapping[str, Any]] = []

    def start_session(self, config: Mapping[str, Any]) -> str:
        del config
        session_id = "sess-mock-1"
        self.started.add(session_id)
        self.calls.append(("start_session", session_id))
        return session_id

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        del task
        if session_id not in self.started:
            raise UnknownSessionError(session_id)
        self.calls.append(("submit_task", session_id))

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

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        if session_id not in self.started:
            raise UnknownSessionError(session_id)
        self.calls.append(("export_ledger", session_id))
        return {"session_id": session_id, "entries": list(self.ledger_entries)}

    def stop_session(self, session_id: str) -> None:
        if session_id not in self.started:
            raise UnknownSessionError(session_id)
        self.calls.append(("stop_session", session_id))
        self.started.discard(session_id)


class _FakeRequester:
    def __init__(self, receipts: Any = None) -> None:
        self.calls: list[tuple[str, str, Mapping[str, Any] | None]] = []
        self.receipts = receipts

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
        if method == "GET" and path == PATH_LIST_RECEIPTS:
            return self.receipts
        raise AssertionError(f"unexpected {method} {path}")


def test_handshake_matches_empty_export_ledger_reality() -> None:
    """CapabilityReport.ledger/observability stay false while entries are empty."""
    adapter = SireAdapter()
    report = adapter.handshake("1.0")
    assert report["provenance"] == "unofficial"
    assert report["hitl"] is True
    assert report["ledger"] is False
    assert report["observability"] is False

    session = adapter.start_session({"run_mode": "oracle"})
    sid = session["id"]
    ledger = adapter.export_ledger(sid)
    assert ledger == empty_ledger(sid)
    assert ledger["entries"] == []
    # Honesty: do not claim ledger while ExportLedger cannot provide entries.
    assert report["ledger"] is (len(ledger["entries"]) > 0)


def test_export_ledger_goes_through_injected_client() -> None:
    client = _RecordingClient()
    adapter = SireAdapter(client=client)
    adapter.start_session({"run_mode": "oracle"})
    ledger = adapter.export_ledger("sess-mock-1")
    assert ledger == {"session_id": "sess-mock-1", "entries": []}
    assert ("export_ledger", "sess-mock-1") in client.calls


def test_export_ledger_unknown_session_raises() -> None:
    adapter = SireAdapter(client=_RecordingClient())
    with pytest.raises(UnknownSessionError):
        adapter.export_ledger("missing")


def test_stub_export_ledger_records_call_and_shape() -> None:
    stub = StubSireClient()
    adapter = SireAdapter(client=stub)
    sid = adapter.start_session({"run_mode": "oracle"})["id"]
    ledger = adapter.export_ledger(sid)
    assert set(ledger.keys()) == {"session_id", "entries"}
    assert ledger["session_id"] == sid
    assert ledger["entries"] == []
    assert ("export_ledger", (sid,)) in stub.calls


def test_http_client_probes_compliance_receipts_but_returns_empty_ledger() -> None:
    transport = _FakeRequester(receipts=[{"id": "rcpt_1"}])
    client = HttpSireClient(transport, worker_id="wrk_1")
    sid = client.start_session({"run_mode": "oracle"})
    ledger = client.export_ledger(sid)
    assert transport.calls[-1] == ("GET", PATH_LIST_RECEIPTS, None)
    # Tenant receipts are not invented into AgentGavel LedgerEntry rows.
    assert ledger == empty_ledger(sid)


def test_ledger_true_only_when_entries_exist_is_documented_by_handshake() -> None:
    """If a future client returns real entries, handshake must flip ledger=True.

    Today's adapter still reports ledger=False; this test locks the empty path
    and the acceptance rule (fields match what ExportLedger can provide).
    """
    client = _RecordingClient()
    adapter = SireAdapter(client=client)
    adapter.start_session({})
    ledger = adapter.export_ledger("sess-mock-1")
    report = adapter.handshake("1.0")
    has_entries = len(ledger["entries"]) > 0
    assert report["ledger"] is has_entries
    assert has_entries is False
