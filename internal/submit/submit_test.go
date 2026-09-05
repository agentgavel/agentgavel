package submit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func sampleEntry() map[string]any {
	return map[string]any{
		"run_id":          "signed-example-1",
		"framework":       "Example Framework",
		"adapter":         "example",
		"adapter_version": "1.0.0",
		"provenance":      "ratified",
		"tab":             "opt-in",
		"sample":          false,
		"gsi":             float64(900),
		"grade":           "AA",
		"pillars": map[string]any{
			"auditability": float64(90),
			"chokepoint":   float64(90),
			"governance":   float64(90),
			"resilience":   float64(90),
		},
		"catastrophic": []any{},
		"na":           []any{},
		"fingerprint": map[string]any{
			"hash": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
		"generated_at": "2026-09-05T12:00:00Z",
	}
}

func TestCanonicalJSONSortsKeys(t *testing.T) {
	raw := []byte(`{"b":1,"a":{"d":2,"c":3},"z":[true,false]}`)
	got, err := CanonicalJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"a":{"c":3,"d":2},"b":1,"z":[true,false]}`
	if string(got) != want {
		t.Fatalf("canonical = %s, want %s", got, want)
	}
}

func TestSignVerifyGolden(t *testing.T) {
	regPath := filepath.Join("..", "..", "dashboard", "keys", "registry.json")
	reg, err := LoadRegistry(regPath)
	if err != nil {
		t.Fatalf("LoadRegistry: %v", err)
	}
	priv, err := LoadPrivateKeyFile(filepath.Join("testdata", "example-framework-test-1.priv.b64"))
	if err != nil {
		t.Fatalf("LoadPrivateKeyFile: %v", err)
	}
	entry := sampleEntry()
	if err := Sign(priv, "example-framework-test-1", entry); err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if entry["signature"] == nil || entry["key_id"] != "example-framework-test-1" {
		t.Fatalf("missing key_id/signature after Sign: %#v", entry)
	}
	if err := Verify(reg, entry); err != nil {
		t.Fatalf("Verify: %v", err)
	}

	// Persist a golden for inspectability.
	goldPath := filepath.Join("testdata", "golden-signed-entry.json")
	out, err := MarshalEntry(entry)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(goldPath, append(out, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyRejectsTamper(t *testing.T) {
	reg, err := LoadRegistry(filepath.Join("..", "..", "dashboard", "keys", "registry.json"))
	if err != nil {
		t.Fatal(err)
	}
	priv, err := LoadPrivateKeyFile(filepath.Join("testdata", "example-framework-test-1.priv.b64"))
	if err != nil {
		t.Fatal(err)
	}
	entry := sampleEntry()
	if err := Sign(priv, "example-framework-test-1", entry); err != nil {
		t.Fatal(err)
	}
	entry["gsi"] = float64(1)
	if err := Verify(reg, entry); err == nil {
		t.Fatal("expected Verify to fail after tamper")
	}
}

func TestVerifyRejectsRevokedKey(t *testing.T) {
	reg, err := LoadRegistry(filepath.Join("..", "..", "dashboard", "keys", "registry.json"))
	if err != nil {
		t.Fatal(err)
	}
	entry := sampleEntry()
	entry["key_id"] = "example-framework-revoked-1"
	entry["signature"] = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="
	if err := Verify(reg, entry); err == nil {
		t.Fatal("expected Verify to reject revoked key")
	}
}

func TestLoadRegistryShape(t *testing.T) {
	reg, err := LoadRegistry(filepath.Join("..", "..", "dashboard", "keys", "registry.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(reg) < 1 {
		t.Fatal("registry empty")
	}
	raw, err := os.ReadFile(filepath.Join("..", "..", "dashboard", "keys", "registry.json"))
	if err != nil {
		t.Fatal(err)
	}
	var arr []json.RawMessage
	if err := json.Unmarshal(raw, &arr); err != nil {
		t.Fatal(err)
	}
}
