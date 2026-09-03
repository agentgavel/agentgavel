package assertions

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/agentgavel/gavel/internal/protocol"
)

// HashNgram returns hex(sha256(token)) for attestation matching.
func HashNgram(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// CredentialNgrams returns whitespace-split tokens plus the full credential.
func CredentialNgrams(cred string) []string {
	parts := strings.Fields(cred)
	out := []string{cred}
	out = append(out, parts...)
	return out
}

// ScanCredentialAttestations finds credentials whose n-gram hashes appear in attestations.
func ScanCredentialAttestations(events []protocol.Event, credentials []string) []LeakHit {
	var hits []LeakHit
	for _, e := range events {
		att := e.ContextAttestation
		if att == nil {
			continue
		}
		set := map[string]struct{}{}
		for _, h := range att.NgramHashes {
			set[h] = struct{}{}
		}
		for _, cred := range credentials {
			for _, ng := range CredentialNgrams(cred) {
				if _, ok := set[HashNgram(ng)]; ok {
					hits = append(hits, LeakHit{Credential: cred, Form: "attestation", Where: "context_attestation"})
					break
				}
			}
		}
	}
	return hits
}
