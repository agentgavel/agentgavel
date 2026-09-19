"""OpenClaw Gateway health probe (T17.9).

Env-gated live path: ``AGENTGAVEL_OPENCLAW_GATEWAY_URL`` must respond to a
simple HTTP GET (default path ``/health`` or ``/``). Fail closed on probe
failure — never report ``runtime=live`` without a successful probe.

``hitl`` stays false until withhold is mapped (ADR 014 / capability map).
"""

from __future__ import annotations

import json
import os
import urllib.error
import urllib.request
from typing import Any


class GatewayProbeError(RuntimeError):
    """Gateway URL set but probe failed."""


def gateway_url_from_env() -> str:
    return (
        os.environ.get("AGENTGAVEL_OPENCLAW_GATEWAY_URL")
        or os.environ.get("OPENCLAW_GATEWAY_URL")
        or ""
    ).strip()


def probe_gateway(
    base_url: str,
    *,
    timeout_s: float = 5.0,
    opener: Any | None = None,
) -> dict[str, Any]:
    """GET health endpoint; return a small probe report or raise."""
    if not base_url or not str(base_url).strip():
        raise GatewayProbeError("gateway URL is empty")
    root = str(base_url).rstrip("/")
    # Prefer /health; fall back to / if 404.
    last_err: Exception | None = None
    for path in ("/health", "/"):
        url = f"{root}{path}"
        req = urllib.request.Request(url, method="GET")
        try:
            if opener is not None:
                with opener.open(req, timeout=timeout_s) as resp:
                    status = getattr(resp, "status", 200)
                    body = resp.read(512).decode("utf-8", errors="replace")
            else:
                with urllib.request.urlopen(req, timeout=timeout_s) as resp:
                    status = getattr(resp, "status", 200)
                    body = resp.read(512).decode("utf-8", errors="replace")
        except urllib.error.HTTPError as exc:
            if exc.code == 404 and path == "/health":
                last_err = exc
                continue
            raise GatewayProbeError(f"gateway HTTP {exc.code} at {url}") from exc
        except urllib.error.URLError as exc:
            raise GatewayProbeError(f"gateway unreachable at {url}: {exc}") from exc
        if status >= 400:
            raise GatewayProbeError(f"gateway HTTP {status} at {url}")
        version = "probed"
        try:
            payload = json.loads(body) if body.strip().startswith("{") else {}
            if isinstance(payload, dict):
                version = str(
                    payload.get("version")
                    or payload.get("gateway_version")
                    or payload.get("openclaw_version")
                    or "probed"
                )
        except json.JSONDecodeError:
            pass
        return {"url": url, "status": status, "framework_version": version}
    if last_err is not None:
        raise GatewayProbeError(f"gateway probe failed: {last_err}") from last_err
    raise GatewayProbeError(f"gateway probe failed for {root}")
