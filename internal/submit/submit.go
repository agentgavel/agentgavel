// Package submit implements ADR 013 Opt-in Ed25519 sign/verify over canonical JSON.
package submit

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// UnsignedKeys are stripped before canonical encoding (ADR 013).
var UnsignedKeys = map[string]struct{}{
	"signature": {},
	"key_id":    {},
	"sample":    {},
}

// Key is one row in dashboard/keys/registry.json.
type Key struct {
	KeyID        string `json:"key_id"`
	Framework    string `json:"framework"`
	Alg          string `json:"alg"`
	PublicKeyB64 string `json:"public_key_b64"`
	AddedAt      string `json:"added_at"`
	Status       string `json:"status"`
}

// Registry is the on-disk maintainer key list.
type Registry []Key

// LoadRegistry reads registry.json (a JSON array).
func LoadRegistry(path string) (Registry, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read registry %s: %w", path, err)
	}
	var reg Registry
	if err := json.Unmarshal(b, &reg); err != nil {
		return nil, fmt.Errorf("parse registry %s: %w", path, err)
	}
	for i, k := range reg {
		if strings.TrimSpace(k.KeyID) == "" {
			return nil, fmt.Errorf("registry[%d]: key_id required", i)
		}
		if k.Alg != "" && k.Alg != "ed25519" {
			return nil, fmt.Errorf("registry key %s: unsupported alg %q", k.KeyID, k.Alg)
		}
		switch k.Status {
		case "active", "revoked":
		default:
			return nil, fmt.Errorf("registry key %s: status %q want active|revoked", k.KeyID, k.Status)
		}
	}
	return reg, nil
}

// FindActive returns the active key with keyID, or an error.
func (r Registry) FindActive(keyID string) (Key, error) {
	for _, k := range r {
		if k.KeyID != keyID {
			continue
		}
		if k.Status != "active" {
			return Key{}, fmt.Errorf("key %s is %s (want active)", keyID, k.Status)
		}
		return k, nil
	}
	return Key{}, fmt.Errorf("key %s not found in registry", keyID)
}

// CanonicalJSON returns UTF-8 JSON with object keys sorted at every level
// and no insignificant whitespace (ADR 013 / JCS subset).
func CanonicalJSON(raw []byte) ([]byte, error) {
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, fmt.Errorf("canonical json: %w", err)
	}
	return marshalCanonical(v)
}

func marshalCanonical(v any) ([]byte, error) {
	switch t := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		buf := []byte{'{'}
		for i, k := range keys {
			if i > 0 {
				buf = append(buf, ',')
			}
			kb, err := json.Marshal(k)
			if err != nil {
				return nil, err
			}
			buf = append(buf, kb...)
			buf = append(buf, ':')
			vb, err := marshalCanonical(t[k])
			if err != nil {
				return nil, err
			}
			buf = append(buf, vb...)
		}
		buf = append(buf, '}')
		return buf, nil
	case []any:
		buf := []byte{'['}
		for i, el := range t {
			if i > 0 {
				buf = append(buf, ',')
			}
			eb, err := marshalCanonical(el)
			if err != nil {
				return nil, err
			}
			buf = append(buf, eb...)
		}
		buf = append(buf, ']')
		return buf, nil
	default:
		return json.Marshal(t)
	}
}

// PayloadForSigning strips signature/key_id/sample and returns canonical bytes.
func PayloadForSigning(entry map[string]any) ([]byte, error) {
	stripped := make(map[string]any, len(entry))
	for k, v := range entry {
		if _, skip := UnsignedKeys[k]; skip {
			continue
		}
		stripped[k] = v
	}
	raw, err := json.Marshal(stripped)
	if err != nil {
		return nil, err
	}
	return CanonicalJSON(raw)
}

// EntryMapFromJSON unmarshals an entry object into a generic map.
func EntryMapFromJSON(raw []byte) (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("entry json: %w", err)
	}
	if m == nil {
		return nil, fmt.Errorf("entry json: null object")
	}
	return m, nil
}

// ParsePrivateKeyB64 accepts base64 of a 64-byte Ed25519 private key or 32-byte seed.
func ParsePrivateKeyB64(s string) (ed25519.PrivateKey, error) {
	b, err := base64.StdEncoding.DecodeString(strings.TrimSpace(s))
	if err != nil {
		return nil, fmt.Errorf("private key base64: %w", err)
	}
	switch len(b) {
	case ed25519.PrivateKeySize:
		return ed25519.PrivateKey(b), nil
	case ed25519.SeedSize:
		return ed25519.NewKeyFromSeed(b), nil
	default:
		return nil, fmt.Errorf("private key length %d want %d or %d", len(b), ed25519.PrivateKeySize, ed25519.SeedSize)
	}
}

// LoadPrivateKeyFile reads a one-line base64 private key from path.
func LoadPrivateKeyFile(path string) (ed25519.PrivateKey, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}
	return ParsePrivateKeyB64(string(b))
}

// Sign sets key_id and signature on entry (mutates the map).
func Sign(priv ed25519.PrivateKey, keyID string, entry map[string]any) error {
	if strings.TrimSpace(keyID) == "" {
		return fmt.Errorf("key_id required")
	}
	payload, err := PayloadForSigning(entry)
	if err != nil {
		return err
	}
	sig := ed25519.Sign(priv, payload)
	entry["key_id"] = keyID
	entry["signature"] = base64.StdEncoding.EncodeToString(sig)
	return nil
}

// Verify checks entry against an active registry key (ADR 013).
func Verify(reg Registry, entry map[string]any) error {
	keyID, _ := entry["key_id"].(string)
	sigB64, _ := entry["signature"].(string)
	if strings.TrimSpace(keyID) == "" || strings.TrimSpace(sigB64) == "" {
		return fmt.Errorf("entry missing key_id or signature")
	}
	k, err := reg.FindActive(keyID)
	if err != nil {
		return err
	}
	framework, _ := entry["framework"].(string)
	if framework != k.Framework {
		return fmt.Errorf("framework %q does not match key %s framework %q", framework, keyID, k.Framework)
	}
	pubBytes, err := base64.StdEncoding.DecodeString(k.PublicKeyB64)
	if err != nil {
		return fmt.Errorf("key %s public_key_b64: %w", keyID, err)
	}
	if len(pubBytes) != ed25519.PublicKeySize {
		return fmt.Errorf("key %s public key length %d", keyID, len(pubBytes))
	}
	sig, err := base64.StdEncoding.DecodeString(sigB64)
	if err != nil {
		return fmt.Errorf("signature base64: %w", err)
	}
	payload, err := PayloadForSigning(entry)
	if err != nil {
		return err
	}
	if !ed25519.Verify(ed25519.PublicKey(pubBytes), payload, sig) {
		return fmt.Errorf("signature verification failed for key %s", keyID)
	}
	return nil
}

// MarshalEntry returns indented JSON for writing to disk (pretty for humans).
func MarshalEntry(entry map[string]any) ([]byte, error) {
	return json.MarshalIndent(entry, "", "  ")
}
