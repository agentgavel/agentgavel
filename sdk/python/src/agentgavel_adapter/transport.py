"""JSON-RPC 2.0 newline-delimited stdio transport (ADR 002).

Mirrors ``internal/protocol/stdio.go``: request/response framing, notifications
for adapter-pushed Events, and an adapter-side serve loop.
"""

from __future__ import annotations

import json
import sys
import threading
from collections import deque
from typing import Any, BinaryIO, Callable, Mapping, MutableMapping, TextIO

# Method names match internal/protocol/stdio.go.
METHOD_HANDSHAKE = "Handshake"
METHOD_START_SESSION = "StartSession"
METHOD_SUBMIT_TASK = "SubmitTask"
METHOD_RESOLVE_APPROVAL = "ResolveApproval"
METHOD_EVENTS_SUBSCRIBE = "Events"
METHOD_EXPORT_LEDGER = "ExportLedger"
METHOD_STOP_SESSION = "StopSession"
METHOD_EVENT_NOTIFY = "Event"

JSONRPC_VERSION = "2.0"

# JSON-RPC 2.0 reserved error codes.
PARSE_ERROR = -32700
INVALID_REQUEST = -32600
METHOD_NOT_FOUND = -32601
INVALID_PARAMS = -32602
INTERNAL_ERROR = -32603


class RPCError(Exception):
    """Raised when a peer returns a JSON-RPC error object."""

    def __init__(self, code: int, message: str) -> None:
        super().__init__(message)
        self.code = code
        self.message = message


def _as_binary(stream: BinaryIO | TextIO) -> BinaryIO:
    """Normalize text or binary streams to a binary interface."""
    buffer = getattr(stream, "buffer", None)
    if buffer is not None:
        return buffer
    return stream  # type: ignore[return-value]


class StdioConn:
    """Newline-delimited JSON-RPC 2.0 connection over a reader/writer pair.

    Read and write locks are separate so ``emit()`` / ``notify`` can flush
    Events while the serve loop is blocked in ``read_request``.
    """

    def __init__(self, reader: BinaryIO | TextIO, writer: BinaryIO | TextIO) -> None:
        self._in = _as_binary(reader)
        self._out = _as_binary(writer)
        self._read_mu = threading.Lock()
        self._write_mu = threading.Lock()
        self._id_mu = threading.Lock()
        self._next_id = 0

    def call(self, method: str, params: Mapping[str, Any] | None = None) -> Any:
        """Send a request and wait for the matching response (engine side)."""
        with self._id_mu:
            self._next_id += 1
            req_id = self._next_id
        payload: dict[str, Any] = {
            "jsonrpc": JSONRPC_VERSION,
            "id": req_id,
            "method": method,
        }
        if params is not None:
            payload["params"] = params
        with self._write_mu:
            self._write_locked(payload)
        with self._read_mu:
            while True:
                line = self._readline_locked()
                if line is None:
                    raise EOFError("stdio closed while waiting for response")
                msg = json.loads(line)
                if "method" in msg and "id" not in msg:
                    # Peer notification; ignore on the client call path.
                    continue
                if msg.get("id") != req_id:
                    continue
                if "error" in msg and msg["error"] is not None:
                    err = msg["error"]
                    raise RPCError(int(err.get("code", INTERNAL_ERROR)), str(err.get("message", "")))
                return msg.get("result")

    def handshake(
        self,
        engine_protocol_version: str,
        *,
        engine_version: str | None = None,
    ) -> Mapping[str, Any]:
        """Convenience Call for METHOD_HANDSHAKE."""
        params: dict[str, Any] = {"engine_protocol_version": engine_protocol_version}
        if engine_version is not None:
            params["engine_version"] = engine_version
        result = self.call(METHOD_HANDSHAKE, params)
        if not isinstance(result, Mapping):
            raise TypeError(f"Handshake result must be an object, got {type(result)!r}")
        return result

    def notify(self, method: str, params: Mapping[str, Any] | None = None) -> None:
        """Send a JSON-RPC notification (no id)."""
        payload: dict[str, Any] = {"jsonrpc": JSONRPC_VERSION, "method": method}
        if params is not None:
            payload["params"] = params
        with self._write_mu:
            self._write_locked(payload)

    def read_request(self) -> dict[str, Any] | None:
        """Read the next request object, or None on EOF."""
        with self._read_mu:
            line = self._readline_locked()
            if line is None:
                return None
            msg = json.loads(line)
            if not isinstance(msg, dict):
                raise ValueError("JSON-RPC message must be an object")
            return msg

    def reply(self, req_id: Any, result: Any) -> None:
        """Write a success response for id."""
        with self._write_mu:
            self._write_locked(
                {"jsonrpc": JSONRPC_VERSION, "id": req_id, "result": result}
            )

    def reply_error(self, req_id: Any, code: int, message: str) -> None:
        """Write an error response for id."""
        with self._write_mu:
            self._write_locked(
                {
                    "jsonrpc": JSONRPC_VERSION,
                    "id": req_id,
                    "error": {"code": code, "message": message},
                }
            )

    def _write_locked(self, obj: Mapping[str, Any]) -> None:
        data = json.dumps(obj, separators=(",", ":"), ensure_ascii=False).encode("utf-8")
        self._out.write(data + b"\n")
        self._out.flush()

    def _readline_locked(self) -> str | None:
        raw = self._in.readline()
        if not raw:
            return None
        if isinstance(raw, bytes):
            return raw.decode("utf-8").rstrip("\r\n")
        return str(raw).rstrip("\r\n")


class EventBuffer:
    """Thread-safe queue of Event payloads pending emission as notifications."""

    def __init__(self) -> None:
        self._mu = threading.Lock()
        self._pending: deque[MutableMapping[str, Any]] = deque()

    def push(self, event: Mapping[str, Any]) -> None:
        with self._mu:
            self._pending.append(dict(event))

    def drain(self) -> list[MutableMapping[str, Any]]:
        with self._mu:
            items = list(self._pending)
            self._pending.clear()
            return items

    def __len__(self) -> int:
        with self._mu:
            return len(self._pending)


class TransportLoop:
    """Adapter-side serve loop: dispatch engine requests and flush emit() buffer."""

    def __init__(
        self,
        conn: StdioConn,
        *,
        on_handshake: Callable[[Mapping[str, Any]], Mapping[str, Any]],
        event_buffer: EventBuffer | None = None,
    ) -> None:
        self._conn = conn
        self._on_handshake = on_handshake
        self._events = event_buffer if event_buffer is not None else EventBuffer()
        self._closed = False

    @property
    def events(self) -> EventBuffer:
        return self._events

    def emit(self, event: Mapping[str, Any]) -> None:
        """Buffer an Event and flush it as a METHOD_EVENT_NOTIFY notification."""
        if not isinstance(event, Mapping):
            raise TypeError("event must be a mapping")
        self._events.push(event)
        self.flush_events()

    def flush_events(self) -> None:
        """Send all buffered events as JSON-RPC notifications."""
        for event in self._events.drain():
            self._conn.notify(METHOD_EVENT_NOTIFY, event)

    def handle_one(self, req: Mapping[str, Any]) -> bool:
        """Dispatch a single request. Returns False if the peer asked to stop serving."""
        method = req.get("method")
        req_id = req.get("id")
        params = req.get("params") or {}

        if not isinstance(method, str) or not method:
            if req_id is not None:
                self._conn.reply_error(req_id, INVALID_REQUEST, "missing method")
            return True

        # Notifications (no id) — ignore unknown; Events subscribe is engine-initiated with id.
        if req_id is None:
            return True

        if not isinstance(params, Mapping):
            self._conn.reply_error(req_id, INVALID_PARAMS, "params must be an object")
            return True

        try:
            if method == METHOD_HANDSHAKE:
                result = self._on_handshake(params)
                self._conn.reply(req_id, result)
            else:
                # Full callback dispatch is T7.3; Handshake is the T7.2 surface.
                self._conn.reply_error(req_id, METHOD_NOT_FOUND, f"method not found: {method}")
        except Exception as exc:  # noqa: BLE001 — surface as JSON-RPC error to peer
            self._conn.reply_error(req_id, INTERNAL_ERROR, str(exc))

        self.flush_events()
        return True

    def serve(self) -> None:
        """Read requests until EOF."""
        while not self._closed:
            req = self._conn.read_request()
            if req is None:
                break
            if not self.handle_one(req):
                break
        self.flush_events()

    def close(self) -> None:
        self._closed = True


def serve_adapter(
    adapter: Any,
    reader: BinaryIO | TextIO | None = None,
    writer: BinaryIO | TextIO | None = None,
) -> None:
    """Run ``adapter`` on stdio until the engine closes the pipe.

    Wires Handshake dispatch and ``adapter.emit`` buffering onto the transport.
    """
    in_stream = reader if reader is not None else sys.stdin.buffer
    out_stream = writer if writer is not None else sys.stdout.buffer
    conn = StdioConn(in_stream, out_stream)
    buffer = EventBuffer()
    loop = TransportLoop(
        conn,
        on_handshake=lambda params: _dispatch_handshake(adapter, params),
        event_buffer=buffer,
    )
    adapter._attach_transport(loop)  # noqa: SLF001 — intentional SDK hook
    try:
        loop.serve()
    finally:
        adapter._detach_transport(loop)  # noqa: SLF001


def _dispatch_handshake(adapter: Any, params: Mapping[str, Any]) -> Mapping[str, Any]:
    engine_protocol_version = params.get("engine_protocol_version")
    if not isinstance(engine_protocol_version, str) or not engine_protocol_version:
        raise ValueError("engine_protocol_version is required")
    engine_version = params.get("engine_version")
    if engine_version is not None and not isinstance(engine_version, str):
        raise ValueError("engine_version must be a string when set")
    result = adapter.handshake(
        engine_protocol_version,
        engine_version=engine_version,
    )
    if not isinstance(result, Mapping):
        raise TypeError("handshake() must return a mapping (CapabilityReport)")
    return dict(result)
