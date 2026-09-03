"""Protocol Event builders for CrewAI adapter hooks (T13.21).

Keeps ``tool_invocation`` / ``context_attestation`` shapes aligned with
``proto/adapter.proto`` and ``internal/protocol``.
"""

from __future__ import annotations

import json
import time
from collections.abc import Mapping, MutableMapping
from typing import Any


def build_tool_invocation(
    tool_name: str,
    tool_id: str,
    phase: str,
    *,
    arguments: Mapping[str, Any] | None = None,
    outcome: str | None = None,
    error: str | None = None,
    refused: bool = False,
) -> dict[str, Any]:
    """Build a ``tool_invocation`` payload (before | after)."""
    if phase not in ("before", "after"):
        raise ValueError(f"phase must be before|after, got {phase!r}")
    inv: dict[str, Any] = {
        "tool_name": tool_name,
        "tool_id": tool_id,
        "phase": phase,
    }
    if arguments is not None:
        inv["arguments_json"] = json.dumps(dict(arguments), separators=(",", ":"))
    if outcome is not None:
        inv["outcome"] = outcome
    if error is not None:
        inv["error"] = error
    if refused:
        inv["refused"] = True
    return inv


def make_event(
    session_id: str,
    seq: int,
    *,
    unix_ms: int | None = None,
    tool_invocation_payload: Mapping[str, Any] | None = None,
    context_attestation_payload: Mapping[str, Any] | None = None,
) -> MutableMapping[str, Any]:
    """Assemble a protocol Event with exactly one kind payload."""
    kinds = [
        ("tool_invocation", tool_invocation_payload),
        ("context_attestation", context_attestation_payload),
    ]
    present = [(k, v) for k, v in kinds if v is not None]
    if len(present) != 1:
        raise ValueError("make_event requires exactly one event kind payload")
    kind, payload = present[0]
    event: MutableMapping[str, Any] = {
        "session_id": session_id,
        "seq": seq,
        "unix_ms": int(time.time() * 1000) if unix_ms is None else unix_ms,
        kind: dict(payload),
    }
    return event
