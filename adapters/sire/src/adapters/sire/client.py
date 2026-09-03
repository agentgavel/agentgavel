"""Mockable Sire client used by the unofficial AgentGavel adapter.

Sire's public HTTP API (``https://api.sire.run/api/v1``) has Workers and Runs,
not AgentGavel sessions. The mapping is best-effort:

* ``start_session``     -> bind a Worker (``GET /workers/{workerId}``)
* ``submit_task``       -> start a Run (``POST /workers/{workerId}/run``)
* ``resolve_approval``  -> decide a Hold (``POST /approvals/{approvalId}/decide``)
* ``export_ledger``     -> probe receipts (``GET /compliance/receipts``); no
  session-scoped hash-linked AgentGavel Ledger yet (empty ``entries``)
* ``stop_session``      -> cancel that Run (``POST /runs/{runId}/cancel``)

Sire's decide body uses ``verb`` (``approve`` / ``deny`` / ``hold``). AgentGavel
``withhold`` maps to Sire ``hold``. A live client needs a bearer token, a
tenant-scoped worker id, and a worker whose model ``base_url`` can be pointed
at the Compliance Oracle.

ExportLedger honesty: Sire's compliance receipts are tenant-scoped, not an
AgentGavel session ledger (hash-linked ``LedgerEntry`` rows). Until a real
session export exists, clients return Ledger-shaped ``{session_id, entries: []}``
and CapabilityReport keeps ``ledger=false`` (SEC-009/010 N/A) plus
``observability=false`` (GSI observability penalty / cap 600).
"""

from __future__ import annotations

from collections.abc import Mapping, MutableMapping, Sequence
from typing import Any, Protocol
from uuid import uuid4

# Public production API prefix documented by the Sire TypeScript SDK.
SIRE_API_V1_PREFIX = "/api/v1"
DEFAULT_SIRE_API_BASE = "https://api.sire.run" + SIRE_API_V1_PREFIX

PATH_GET_WORKER = "/workers/{worker_id}"
PATH_RUN_WORKER = "/workers/{worker_id}/run"
PATH_CANCEL_RUN = "/runs/{run_id}/cancel"
PATH_DECIDE_APPROVAL = "/approvals/{approval_id}/decide"
# Documented probe only — not a session-scoped AgentGavel ledger export.
PATH_LIST_RECEIPTS = "/compliance/receipts"


def empty_ledger(session_id: str) -> dict[str, Any]:
    """Return a protocol ``Ledger`` mapping with no entries."""
    return {"session_id": session_id, "entries": []}


# AgentGavel Decision wire names → Sire Decision.verb (interfaces.md).
_SIRE_VERB = {
    "approve": "approve",
    "deny": "deny",
    "withhold": "hold",
}
_DECISION_BY_INT = {
    1: "approve",
    2: "deny",
    3: "withhold",
}


class SireClientError(Exception):
    """Lifecycle call failed against Sire (or a test double)."""


class UnknownSessionError(SireClientError):
    """Adapter session id is not known to this client."""


def wire_decision(decision: str | int) -> str:
    """Normalize a ResolveApproval decision to AgentGavel wire names.

    Accepts proto enum integers (1/2/3) or names (``approve`` / ``deny`` /
    ``withhold``, optionally ``DECISION_``-prefixed).
    """
    if isinstance(decision, bool):
        raise SireClientError(f"unknown decision {decision!r}")
    if isinstance(decision, int):
        try:
            return _DECISION_BY_INT[decision]
        except KeyError as exc:
            raise SireClientError(f"unknown decision wire value {decision}") from exc
    name = str(decision).strip().lower()
    prefix = "decision_"
    if name.startswith(prefix):
        name = name[len(prefix) :]
    if name not in _SIRE_VERB:
        raise SireClientError(f"unknown decision {decision!r}")
    return name


def sire_decide_verb(decision: str | int) -> str:
    """Map an AgentGavel decision onto Sire's ``Decision.verb``."""
    return _SIRE_VERB[wire_decision(decision)]


class SireClient(Protocol):
    """Thin surface the adapter drives. Tests inject a mock."""

    def start_session(self, config: Mapping[str, Any]) -> str:
        """Bind SessionConfig to a Sire worker; return an adapter session id."""

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        """Start a Sire run for the bound worker using the task prompt as seed."""

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str,
        *,
        principal: str | None = None,
    ) -> None:
        """POST Sire ``/approvals/{approvalId}/decide`` for this adapter session."""

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        """Return a protocol ``Ledger`` for the adapter session.

        Must include ``session_id`` and ``entries`` (a sequence). Empty entries
        are honest when Sire has no session-scoped hash-linked ledger export.
        """

    def stop_session(self, session_id: str) -> None:
        """Cancel the session's run if one exists, then drop the binding."""


class Requester(Protocol):
    """HTTP transport for :class:`HttpSireClient`. Tests inject a fake."""

    def request(
        self,
        method: str,
        path: str,
        *,
        json: Mapping[str, Any] | None = None,
    ) -> Mapping[str, Any] | None:
        """Perform one JSON request. ``path`` is relative to the API v1 base."""


class StubSireClient:
    """In-memory default. No network; enough for Handshake and unit tests.

    Records start/submit/stop so tests can assert ordering without Sire.
    """

    def __init__(self) -> None:
        self.sessions: dict[str, dict[str, Any]] = {}
        self.calls: list[tuple[str, tuple[Any, ...]]] = []

    def start_session(self, config: Mapping[str, Any]) -> str:
        session_id = f"sire-sess-{uuid4().hex[:12]}"
        self.sessions[session_id] = {
            "config": dict(config),
            "tasks": [],
            "run_id": None,
            "stopped": False,
        }
        self.calls.append(("start_session", (session_id,)))
        return session_id

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        record = self._require(session_id)
        if record["stopped"]:
            raise SireClientError(f"session {session_id} already stopped")
        record["tasks"].append(dict(task))
        record["run_id"] = f"sire-run-{uuid4().hex[:12]}"
        self.calls.append(("submit_task", (session_id, dict(task))))

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str,
        *,
        principal: str | None = None,
    ) -> None:
        record = self._require(session_id)
        if record["stopped"]:
            raise SireClientError(f"session {session_id} already stopped")
        wire = wire_decision(decision)
        record.setdefault("approvals", []).append(
            {
                "approval_id": approval_id,
                "decision": wire,
                "principal": principal,
            }
        )
        self.calls.append(("resolve_approval", (session_id, approval_id, wire, principal)))

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        self._require(session_id)
        ledger = empty_ledger(session_id)
        self.calls.append(("export_ledger", (session_id,)))
        return ledger

    def stop_session(self, session_id: str) -> None:
        record = self._require(session_id)
        record["stopped"] = True
        self.calls.append(("stop_session", (session_id,)))

    def _require(self, session_id: str) -> dict[str, Any]:
        try:
            return self.sessions[session_id]
        except KeyError as exc:
            raise UnknownSessionError(session_id) from exc


class HttpSireClient:
    """Documented HTTP mapping. Does not ship a live urllib default.

    Inject a :class:`Requester` that attaches ``Authorization: Bearer`` and
    targets ``DEFAULT_SIRE_API_BASE`` (or a local Sire). Without a requester,
    calls fail loudly instead of fabricating success.
    """

    def __init__(
        self,
        requester: Requester | None = None,
        *,
        worker_id: str | None = None,
        api_base: str = DEFAULT_SIRE_API_BASE,
    ) -> None:
        self._requester = requester
        self._worker_id = worker_id
        self.api_base = api_base.rstrip("/")
        self._sessions: dict[str, MutableMapping[str, Any]] = {}

    def start_session(self, config: Mapping[str, Any]) -> str:
        requester = self._require_requester()
        worker_id = self._resolve_worker_id(config)
        path = PATH_GET_WORKER.format(worker_id=worker_id)
        body = requester.request("GET", path)
        if not isinstance(body, Mapping):
            raise SireClientError(
                f"GET {path} must return a worker object; got {type(body).__name__}"
            )
        session_id = f"sire-sess-{uuid4().hex[:12]}"
        self._sessions[session_id] = {
            "worker_id": worker_id,
            "config": dict(config),
            "run_id": None,
        }
        return session_id

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        requester = self._require_requester()
        record = self._require_session(session_id)
        path = PATH_RUN_WORKER.format(worker_id=record["worker_id"])
        seed: dict[str, Any] = {
            "prompt": task.get("prompt", ""),
            "task_id": task.get("id", ""),
            "metadata": dict(task.get("metadata") or {}),
        }
        body = requester.request("POST", path, json=seed)
        run_id = None
        if isinstance(body, Mapping):
            run_id = body.get("runId") or body.get("run_id")
        if not run_id:
            raise SireClientError(f"POST {path} must return runId (Sire RunRef); got {body!r}")
        record["run_id"] = str(run_id)

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str,
        *,
        principal: str | None = None,
    ) -> None:
        del principal  # Sire decide body is verb(+note); harness principal is event-only.
        requester = self._require_requester()
        self._require_session(session_id)
        path = PATH_DECIDE_APPROVAL.format(approval_id=approval_id)
        requester.request(
            "POST",
            path,
            json={"verb": sire_decide_verb(decision)},
        )

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        requester = self._require_requester()
        self._require_session(session_id)
        # Documented Sire surface (tenant receipts). Not session-scoped and not
        # an AgentGavel hash-linked Ledger — do not invent LedgerEntry rows.
        body = requester.request("GET", PATH_LIST_RECEIPTS)
        if body is not None and not isinstance(body, (Mapping, Sequence)):
            raise SireClientError(
                f"GET {PATH_LIST_RECEIPTS} must return a list/object or empty; "
                f"got {type(body).__name__}"
            )
        return empty_ledger(session_id)

    def stop_session(self, session_id: str) -> None:
        requester = self._require_requester()
        record = self._require_session(session_id)
        run_id = record.get("run_id")
        if run_id:
            path = PATH_CANCEL_RUN.format(run_id=run_id)
            requester.request("POST", path)
        del self._sessions[session_id]

    def _require_requester(self) -> Requester:
        if self._requester is None:
            raise SireClientError(
                "HttpSireClient has no requester; inject one that calls "
                f"{self.api_base} with a bearer token, or use StubSireClient"
            )
        return self._requester

    def _resolve_worker_id(self, config: Mapping[str, Any]) -> str:
        extra = config.get("extra") or {}
        worker_id = None
        if isinstance(extra, Mapping):
            worker_id = extra.get("sire_worker_id")
        worker_id = worker_id or self._worker_id
        if not worker_id:
            raise SireClientError(
                "real Sire needs extra['sire_worker_id'] or HttpSireClient(worker_id=...)"
            )
        return str(worker_id)

    def _require_session(self, session_id: str) -> MutableMapping[str, Any]:
        try:
            return self._sessions[session_id]
        except KeyError as exc:
            raise UnknownSessionError(session_id) from exc
