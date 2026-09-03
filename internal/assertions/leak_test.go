package assertions

import (
	"encoding/base64"
	"net/url"
	"testing"

	"github.com/agentgavel/agentgavel/internal/protocol"
)

func TestCredLeak(t *testing.T) {
	key := "agavel_test_static_key_DO_NOT_USE_9f3c2a1b"
	b64 := base64.StdEncoding.EncodeToString([]byte(key))
	// Force a URL-encoded form that differs from raw.
	urlKey := "agavel test+key/1"
	urlEnc := url.QueryEscape(urlKey)

	events := []protocol.Event{
		{ContextSnapshot: &protocol.ContextSnapshot{Role: "tool", Content: "header " + b64}},
		{ToolInvocation: &protocol.ToolInvocation{Error: "echo " + urlEnc}},
	}
	hits := ScanCredentialLeaks(events, []string{key, urlKey})
	var forms []string
	for _, h := range hits {
		forms = append(forms, h.Form)
	}
	if !contains(forms, "base64") {
		t.Fatalf("missing base64 hit: %#v", hits)
	}
	if !contains(forms, "url") {
		t.Fatalf("missing url hit: %#v", hits)
	}
}

func contains(ss []string, want string) bool {
	for _, s := range ss {
		if s == want {
			return true
		}
	}
	return false
}
