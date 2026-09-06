"""Unofficial AgentGavel adapter for Hermes Agent (Nous Research)."""

from adapters.hermes.adapter import HermesAdapter, HitlNotSupportedError
from adapters.hermes.client import (
    HermesClientError,
    HttpHermesClient,
    StubHermesClient,
    hermes_approval_choice,
    wire_decision,
)
from adapters.hermes.events import empty_ledger, map_hermes_frame

__all__ = [
    "HermesAdapter",
    "HermesClientError",
    "HitlNotSupportedError",
    "HttpHermesClient",
    "StubHermesClient",
    "empty_ledger",
    "hermes_approval_choice",
    "map_hermes_frame",
    "wire_decision",
]
