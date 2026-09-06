"""Mockable Hermes API-server client for ResolveApproval (T15.25).

Control plane: OpenAI-compatible API server
(``gateway/platforms/api_server.py`` / ``api_server_runs.py``).

* ``resolve_approval`` → ``POST /v1/runs/{run_id}/approval``
  Body: ``{"choice": ..., "request_id": <approval_id>}``
  Hermes choices: ``once`` | ``session`` | ``always`` | ``deny``
  (aliases ``approve`` / ``allow`` → ``once`` per upstream).

AgentGavel Decision mapping (wire → Hermes ``choice``):

* ``approve``  → ``once`` (single allow; not ``always`` / YOLO)
* ``deny``     → ``deny``
* ``withhold`` → no HTTP call — leave the gate pending so Hermes
  timeout / unattended defaults apply (docs: unanswered → deny).
  There is no Hermes ``hold`` verb.

Default is :class:`StubHermesClient` (no network). Inject
:class:`HttpHermesClient` with a :class:`Requester` for a live API server.
"""

from __future__ import annotations

from collections.abc import Mapping, MutableMapping
from typing import Any, Protocol
from uuid import uuid4

DEFAULT_HERMES_API_BASE = "http://127.0.0.1:8642"
PATH_RUN_APPROVAL = "/v1/runs/{run_id}/approval"

# AgentGavel Decision wire names → Hermes approval ``choice``.
# withhold is handled by skipping the POST (see hermes_approval_choice).
_HERMES_CHOICE = {
    "approve": "once",
    "deny": "deny",
}
_DECISION_BY_INT = {
    1: "approve",
    2: "deny",
    3: "withhold",
}
_WIRE_DECISIONS = frozenset({"approve", "deny", "withhold"})


class HermesClientError(Exception):
    """Lifecycle call failed against Hermes (or a test double)."""


class UnknownSessionError(HermesClientError):
    """Adapter session id is not known to this client."""


def wire_decision(decision: str | int) -> str:
    """Normalize a ResolveApproval decision to AgentGavel wire names.

    Accepts proto enum integers (1/2/3) or names (``approve`` / ``deny`` /
    ``withhold``, optionally ``DECISION_``-prefixed).
    """
    if isinstance(decision, bool):
        raise HermesClientError(f"unknown decision {decision!r}")
    if isinstance(decision, int):
        try:
            return _DECISION_BY_INT[decision]
        except KeyError as exc:
            raise HermesClientError(f"unknown decision wire value {decision}") from exc
    name = str(decision).strip().lower()
    prefix = "decision_"
    if name.startswith(prefix):
        name = name[len(prefix) :]
    if name not in _WIRE_DECISIONS:
        raise HermesClientError(f"unknown decision {decision!r}")
    return name


def hermes_approval_choice(decision: str | int) -> str | None:
    """Map an AgentGavel decision onto Hermes ``choice``, or None for withhold.

    ``None`` means do not POST — leave the pending approval for timeout /
    unattended deny semantics (capability map / Hermes docs).
    """
    wire = wire_decision(decision)
    if wire == "withhold":
        return None
    return _HERMES_CHOICE[wire]


class HermesClient(Protocol):
    """Thin surface the adapter drives. Tests inject a mock."""

    def start_session(self, config: Mapping[str, Any]) -> str:
        """Bind SessionConfig; return an adapter session id."""

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        """Start (or record) a Hermes run for the session."""

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str,
        *,
        principal: str | None = None,
    ) -> None:
        """POST Hermes ``/v1/runs/{run_id}/approval`` (or withhold locally)."""

    def stop_session(self, session_id: str) -> None:
        """Drop the session binding (stop run is T15.24 stub territory)."""


class Requester(Protocol):
    """HTTP transport for :class:`HttpHermesClient`. Tests inject a fake."""

    def request(
        self,
        method: str,
        path: str,
        *,
        json: Mapping[str, Any] | None = None,
    ) -> Mapping[str, Any] | None:
        """Perform one JSON request. ``path`` is relative to the API base."""


class StubHermesClient:
    """In-memory default. No network; enough for Handshake and unit tests."""

    def __init__(self) -> None:
        self.sessions: dict[str, dict[str, Any]] = {}
        self.calls: list[tuple[str, tuple[Any, ...]]] = []

    def start_session(self, config: Mapping[str, Any]) -> str:
        session_id = f"hermes-sess-{uuid4().hex[:12]}"
        self.sessions[session_id] = {
            "config": dict(config),
            "tasks": [],
            "run_id": None,
            "approvals": [],
            "stopped": False,
        }
        self.calls.append(("start_session", (session_id,)))
        return session_id

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        record = self._require(session_id)
        if record["stopped"]:
            raise HermesClientError(f"session {session_id} already stopped")
        record["tasks"].append(dict(task))
        record["run_id"] = f"hermes-run-{uuid4().hex[:12]}"
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
            raise HermesClientError(f"session {session_id} already stopped")
        wire = wire_decision(decision)
        choice = hermes_approval_choice(wire)
        record["approvals"].append(
            {
                "approval_id": approval_id,
                "decision": wire,
                "hermes_choice": choice,
                "principal": principal,
            }
        )
        self.calls.append(("resolve_approval", (session_id, approval_id, wire, principal, choice)))

    def stop_session(self, session_id: str) -> None:
        record = self._require(session_id)
        record["stopped"] = True
        self.calls.append(("stop_session", (session_id,)))

    def _require(self, session_id: str) -> dict[str, Any]:
        try:
            return self.sessions[session_id]
        except KeyError as exc:
            raise UnknownSessionError(session_id) from exc


class HttpHermesClient:
    """Documented HTTP mapping to the Hermes API server.

    Inject a :class:`Requester` that targets ``DEFAULT_HERMES_API_BASE`` (or a
    configured host) with ``API_SERVER_KEY`` auth. Without a requester, calls
    fail loudly instead of fabricating success.
    """

    def __init__(
        self,
        requester: Requester | None = None,
        *,
        api_base: str = DEFAULT_HERMES_API_BASE,
    ) -> None:
        self._requester = requester
        self.api_base = api_base.rstrip("/")
        self._sessions: dict[str, MutableMapping[str, Any]] = {}

    def start_session(self, config: Mapping[str, Any]) -> str:
        # Full POST /api/sessions is out of T15.25 scope; bind locally and
        # require submit_task / resolve to attach a run_id.
        del config  # reserved for session create / model binding
        self._require_requester()
        session_id = f"hermes-sess-{uuid4().hex[:12]}"
        self._sessions[session_id] = {
            "run_id": None,
            "approvals": [],
        }
        return session_id

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        # T15.25 owns ResolveApproval only; run create stays for later waves.
        # Allow extra['hermes_run_id'] so HTTP resolve can be exercised now.
        record = self._require_session(session_id)
        extra = task.get("metadata") if isinstance(task.get("metadata"), Mapping) else {}
        run_id = None
        if isinstance(extra, Mapping):
            run_id = extra.get("hermes_run_id")
        run_id = run_id or task.get("hermes_run_id")
        if run_id:
            record["run_id"] = str(run_id)

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str,
        *,
        principal: str | None = None,
    ) -> None:
        del principal  # Hermes approval body is choice(+request_id); principal is event-only.
        record = self._require_session(session_id)
        wire = wire_decision(decision)
        choice = hermes_approval_choice(wire)
        record.setdefault("approvals", []).append(
            {"approval_id": approval_id, "decision": wire, "hermes_choice": choice}
        )
        if choice is None:
            # withhold: leave pending for Hermes timeout / unattended deny.
            return
        run_id = record.get("run_id")
        if not run_id:
            raise HermesClientError(
                f"session {session_id} has no hermes run_id; "
                "set via submit_task metadata['hermes_run_id'] before ResolveApproval"
            )
        requester = self._require_requester()
        path = PATH_RUN_APPROVAL.format(run_id=run_id)
        body: dict[str, Any] = {"choice": choice}
        if approval_id:
            body["request_id"] = approval_id
        requester.request("POST", path, json=body)

    def stop_session(self, session_id: str) -> None:
        self._require_session(session_id)
        del self._sessions[session_id]

    def _require_requester(self) -> Requester:
        if self._requester is None:
            raise HermesClientError(
                "HttpHermesClient has no requester; inject one that calls "
                f"{self.api_base} with API_SERVER_KEY, or use StubHermesClient"
            )
        return self._requester

    def _require_session(self, session_id: str) -> MutableMapping[str, Any]:
        try:
            return self._sessions[session_id]
        except KeyError as exc:
            raise UnknownSessionError(session_id) from exc
