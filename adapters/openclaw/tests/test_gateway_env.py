"""T17.9: OpenClaw Gateway env probe — stub default, fail-closed live."""

from __future__ import annotations

import json
import threading
from http.server import BaseHTTPRequestHandler, HTTPServer
from typing import Any

import pytest

from adapters.openclaw.adapter import OpenClawAdapter
from adapters.openclaw.gateway import GatewayProbeError, probe_gateway
from adapters.openclaw.runtime_env import adapter_from_env


def test_default_runtime_stub() -> None:
    report = OpenClawAdapter().handshake("1.0")
    assert report["runtime"] == "stub"
    assert report["hitl"] is False
    assert report["framework_version"] == "unprobed"


def test_adapter_from_env_stub_without_url(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.delenv("AGENTGAVEL_OPENCLAW_GATEWAY_URL", raising=False)
    monkeypatch.delenv("OPENCLAW_GATEWAY_URL", raising=False)
    adapter = adapter_from_env()
    assert adapter.handshake("1.0")["runtime"] == "stub"


class _HealthHandler(BaseHTTPRequestHandler):
    def log_message(self, format: str, *args: Any) -> None:  # noqa: A003
        del format, args

    def do_GET(self) -> None:  # noqa: N802
        if self.path not in ("/health", "/"):
            self.send_response(404)
            self.end_headers()
            return
        body = json.dumps({"version": "gw-test-1"}).encode("utf-8")
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)


@pytest.fixture()
def gateway_url() -> Any:
    server = HTTPServer(("127.0.0.1", 0), _HealthHandler)
    thread = threading.Thread(target=server.serve_forever, daemon=True)
    thread.start()
    host, port = server.server_address[:2]
    try:
        yield f"http://{host}:{port}"
    finally:
        server.shutdown()
        thread.join(timeout=5)


def test_probe_gateway_ok(gateway_url: str) -> None:
    report = probe_gateway(gateway_url)
    assert report["framework_version"] == "gw-test-1"


def test_adapter_from_env_live_after_probe(
    monkeypatch: pytest.MonkeyPatch,
    gateway_url: str,
) -> None:
    monkeypatch.setenv("AGENTGAVEL_OPENCLAW_GATEWAY_URL", gateway_url)
    adapter = adapter_from_env()
    report = adapter.handshake("1.0")
    assert report["runtime"] == "live"
    assert report["hitl"] is False  # withhold still unmapped
    assert report["framework_version"] == "gw-test-1"


def test_adapter_from_env_fail_closed_bad_url(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv(
        "AGENTGAVEL_OPENCLAW_GATEWAY_URL",
        "http://127.0.0.1:1",
    )
    with pytest.raises(GatewayProbeError):
        adapter_from_env()
