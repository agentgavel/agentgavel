"""ExportLedger honesty: empty entries ↔ CapabilityReport.ledger=false (T15.19)."""

from __future__ import annotations

import pytest

from adapters.openclaw.adapter import OpenClawAdapter
from adapters.openclaw.events import empty_ledger


def test_handshake_ledger_false_matches_empty_export() -> None:
    adapter = OpenClawAdapter()
    report = adapter.handshake("1.0")
    assert report["ledger"] is False
    assert report["observability"] is False

    sid = adapter.start_session({})["id"]
    ledger = adapter.export_ledger(sid)
    assert ledger == empty_ledger(sid)
    assert ledger["entries"] == []
    # Honesty: do not claim ledger while ExportLedger cannot provide hash-chain.
    assert report["ledger"] is (len(ledger["entries"]) > 0)


def test_export_ledger_unknown_session_raises() -> None:
    adapter = OpenClawAdapter()
    with pytest.raises(KeyError, match="unknown session"):
        adapter.export_ledger("missing")


def test_export_ledger_shape_keys() -> None:
    adapter = OpenClawAdapter()
    sid = adapter.start_session({})["id"]
    ledger = adapter.export_ledger(sid)
    assert set(ledger.keys()) == {"session_id", "entries"}
    assert ledger["session_id"] == sid
    assert isinstance(ledger["entries"], list)
