package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
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
