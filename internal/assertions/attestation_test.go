package assertions

import (
	"testing"

	"github.com/agentgavel/agentgavel/internal/protocol"
)

func TestCredLeakAttestation(t *testing.T) {
	cred := "agavel_test_static_key_DO_NOT_USE_9f3c2a1b"
	events := []protocol.Event{{
		ContextAttestation: &protocol.ContextAttestation{
			Algorithm:   "sha256",
			NgramHashes: []string{HashNgram(cred)},
		},
	}}
	hits := ScanCredentialAttestations(events, []string{cred})
	if len(hits) != 1 || hits[0].Form != "attestation" {
		t.Fatalf("hits %#v", hits)
	}
	miss := ScanCredentialAttestations(events, []string{"other_secret"})
	if len(miss) != 0 {
		t.Fatalf("unexpected %#v", miss)
	}
}
