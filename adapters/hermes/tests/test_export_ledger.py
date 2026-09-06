"""ExportLedger shape + CapabilityReport honesty (T15.26 / Sire T10.4 pattern)."""

from __future__ import annotations

import pytest

from adapters.hermes.adapter import HermesAdapter
from adapters.hermes.events import empty_ledger


def test_handshake_matches_empty_export_ledger_reality() -> None:
    """CapabilityReport.ledger stays false while ExportLedger entries are empty."""
    adapter = HermesAdapter()
    report = adapter.handshake("1.0")
    assert report["provenance"] == "unofficial"
    assert report["ledger"] is False
    # Observability is independent: Events sink is wired (T15.26).
    assert report["observability"] is True

    session = adapter.start_session({})
    sid = session["id"]
    ledger = adapter.export_ledger(sid)
    assert ledger == empty_ledger(sid)
    assert ledger["entries"] == []
    # Honesty: do not claim ledger while ExportLedger cannot provide entries.
    assert report["ledger"] is (len(ledger["entries"]) > 0)


def test_export_ledger_shape_and_unknown_session() -> None:
    adapter = HermesAdapter()
    sid = adapter.start_session({})["id"]
    ledger = adapter.export_ledger(sid)
    assert set(ledger.keys()) == {"session_id", "entries"}
    assert ledger["session_id"] == sid
    assert ledger["entries"] == []
    with pytest.raises(KeyError):
        adapter.export_ledger("missing")


def test_ledger_true_only_when_entries_exist_is_documented_by_handshake() -> None:
    """Today's adapter reports ledger=False; lock the empty-path acceptance rule."""
    adapter = HermesAdapter()
    sid = adapter.start_session({})["id"]
    ledger = adapter.export_ledger(sid)
    report = adapter.handshake("1.0")
    has_entries = len(ledger["entries"]) > 0
    assert report["ledger"] is has_entries
    assert has_entries is False
