package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Fingerprint identifies the exact configuration a run was executed under,
// so results can be compared apples-to-apples across runs.
type Fingerprint struct {
	ScenarioVersion  string
	FrameworkVersion string
	ConfigHash       string
	AdapterVersion   string
	Model            string
	SeedSet          []int64
}

// Hash returns a stable hex-encoded SHA-256 digest of the fingerprint fields.
// Seeds are sorted before hashing so equivalent seed sets always produce the
// same hash regardless of input order.
func (f Fingerprint) Hash() string {
	seeds := make([]int64, len(f.SeedSet))
	copy(seeds, f.SeedSet)
	sort.Slice(seeds, func(i, j int) bool { return seeds[i] < seeds[j] })

	seedStrs := make([]string, len(seeds))
	for i, s := range seeds {
		seedStrs[i] = fmt.Sprintf("%d", s)
	}

	parts := []string{
		"scenario.version=" + f.ScenarioVersion,
		"framework.version=" + f.FrameworkVersion,
		"config.hash=" + f.ConfigHash,
		"adapter.version=" + f.AdapterVersion,
		"model=" + f.Model,
		"seed.set=" + strings.Join(seedStrs, ","),
	}

	sum := sha256.Sum256([]byte(strings.Join(parts, "\n")))
	return hex.EncodeToString(sum[:])
}

// Fields returns the fingerprint as a flat map, suitable for embedding in a
// RunArtifact.
func (f Fingerprint) Fields() map[string]string {
	seeds := make([]int64, len(f.SeedSet))
	copy(seeds, f.SeedSet)
	sort.Slice(seeds, func(i, j int) bool { return seeds[i] < seeds[j] })

	seedStrs := make([]string, len(seeds))
	for i, s := range seeds {
		seedStrs[i] = fmt.Sprintf("%d", s)
	}

	return map[string]string{
		"scenario.version":  f.ScenarioVersion,
		"framework.version": f.FrameworkVersion,
		"config.hash":       f.ConfigHash,
		"adapter.version":   f.AdapterVersion,
		"model":             f.Model,
		"seed.set":          strings.Join(seedStrs, ","),
		"hash":              f.Hash(),
	}
}

// FingerprintFromFields rebuilds a Fingerprint from Fields()/summary JSON keys.
// Unknown or empty keys are left zero; seed.set is required for a useful pin.
func FingerprintFromFields(fields map[string]string) (Fingerprint, error) {
	var f Fingerprint
	if fields == nil {
		return f, fmt.Errorf("fingerprint: empty fields")
	}
	f.ScenarioVersion = fields["scenario.version"]
	f.FrameworkVersion = fields["framework.version"]
	f.ConfigHash = fields["config.hash"]
	f.AdapterVersion = fields["adapter.version"]
	f.Model = fields["model"]
	seedStr := strings.TrimSpace(fields["seed.set"])
	if seedStr == "" {
		return f, fmt.Errorf("fingerprint: missing seed.set")
	}
	seeds, err := ParseSeedSet(seedStr)
	if err != nil {
		return f, err
	}
	f.SeedSet = seeds
	return f, nil
}

// ParseSeedSet parses a comma-separated seed.set string (e.g. "0,1,2").
func ParseSeedSet(s string) ([]int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("fingerprint: empty seed.set")
	}
	parts := strings.Split(s, ",")
	out := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("fingerprint: parse seed %q: %w", p, err)
		}
		out = append(out, n)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("fingerprint: empty seed.set")
	}
	return out, nil
}

// LoadFingerprintFile loads pins from fingerprint.json (flat Fields map) or
// summary.json (nested "fingerprint" object).
func LoadFingerprintFile(path string) (Fingerprint, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Fingerprint{}, fmt.Errorf("fingerprint: read %s: %w", path, err)
	}
	var flat map[string]string
	if err := json.Unmarshal(b, &flat); err == nil && strings.TrimSpace(flat["seed.set"]) != "" {
		return FingerprintFromFields(flat)
	}
	var wrapped struct {
		Fingerprint map[string]string `json:"fingerprint"`
	}
	if err := json.Unmarshal(b, &wrapped); err != nil {
		return Fingerprint{}, fmt.Errorf("fingerprint: parse %s: %w", path, err)
	}
	if len(wrapped.Fingerprint) == 0 {
		return Fingerprint{}, fmt.Errorf("fingerprint: no fingerprint fields in %s", path)
	}
	return FingerprintFromFields(wrapped.Fingerprint)
}

// WriteFingerprintFile writes Fields() as results/.../fingerprint.json.
func WriteFingerprintFile(dir string, fields map[string]string) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	if fields == nil {
		fields = map[string]string{}
	}
	path := filepath.Join(dir, "fingerprint.json")
	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return "", err
	}
	return path, nil
}
