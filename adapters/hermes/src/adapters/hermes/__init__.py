"""Unofficial AgentGavel adapter for Hermes Agent (Nous Research)."""

from adapters.hermes.adapter import HermesAdapter, HitlNotSupportedError
from adapters.hermes.client import (
    HermesClientError,
    HttpHermesClient,
    StubHermesClient,
    hermes_approval_choice,
    wire_decision,
)

__all__ = [
    "HermesAdapter",
    "HermesClientError",
    "HitlNotSupportedError",
    "HttpHermesClient",
    "StubHermesClient",
    "hermes_approval_choice",
    "wire_decision",
]
