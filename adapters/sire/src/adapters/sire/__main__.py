"""CLI entry: ``python -m adapters.sire`` serves stdio; ``--help`` prints usage."""

from __future__ import annotations

import argparse

from adapters.sire.adapter import SireAdapter


def _build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="python -m adapters.sire",
        description=(
            "Unofficial AgentGavel Sire adapter (stdio JSON-RPC). "
            "Provenance is always unofficial (ADR 007)."
        ),
    )
    parser.add_argument(
        "--serve",
        action="store_true",
        help="run the stdio JSON-RPC serve loop (default with no flags)",
    )
    return parser


def main(argv: list[str] | None = None) -> None:
    """Parse argv; ``--help`` prints usage, otherwise start stdio serve()."""
    _build_parser().parse_args(argv)
    # Matches FakeAdapter default: module entry starts the stdio serve loop.
    SireAdapter().serve()


if __name__ == "__main__":
    main()
