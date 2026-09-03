"""Context attestation helper (ADR 005 / UC-010).

Matches ``internal/assertions`` hashing so SEC-004 leak checks against
``context_attestation`` events agree with the Go scanner: hex(sha256(token))
over whitespace-split n-grams plus the full string when multi-token.
"""

from __future__ import annotations

import hashlib
from collections.abc import Iterable, Mapping

ALGORITHM_SHA256 = "sha256"


def hash_ngram(token: str) -> str:
    """Return hex(sha256(token)) — same contract as Go ``HashNgram``."""
    return hashlib.sha256(token.encode("utf-8")).hexdigest()


def ngrams_for_text(text: str) -> list[str]:
    """Full string plus whitespace-split tokens (mirrors Go ``CredentialNgrams``)."""
    if text == "":
        return []
    parts = text.split()
    return [text, *parts]


def ngram_hashes(text: str) -> list[str]:
    """Hash each attestation n-gram of ``text`` (deduplicated, stable order)."""
    seen: set[str] = set()
    hashes: list[str] = []
    for ng in ngrams_for_text(text):
        h = hash_ngram(ng)
        if h not in seen:
            seen.add(h)
            hashes.append(h)
    return hashes


def attest_texts(texts: Iterable[str]) -> dict[str, object]:
    """Build a protocol ``ContextAttestation`` payload from context fragments."""
    seen: set[str] = set()
    hashes: list[str] = []
    for text in texts:
        for h in ngram_hashes(text):
            if h not in seen:
                seen.add(h)
                hashes.append(h)
    return {
        "ngram_hashes": hashes,
        "algorithm": ALGORITHM_SHA256,
    }


def context_attestation_payload(text: str) -> Mapping[str, object]:
    """Attest a single context fragment (prompt, tool arg, etc.)."""
    return attest_texts([text])
