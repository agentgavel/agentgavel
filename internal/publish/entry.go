// Package publish writes Unratified leaderboard entries for the static dashboard.
package publish

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/agentgavel/agentgavel/internal/report"
)

const (
	TabUnratified = "unratified"
	TabOptIn      = "opt-in"

	IndexFileName = "index.json"
	dataDirName   = "data"
)

// Entry is one dashboard/data/<run-id>.json row (mirrors schema.json).
type Entry struct {
	RunID          string             `json:"run_id"`
	Framework      string             `json:"framework"`
	Adapter        string             `json:"adapter"`
	AdapterVersion string             `json:"adapter_version"`
	Provenance     string             `json:"provenance"`
	Tab            string             `json:"tab"`
	Sample         bool               `json:"sample"`
	GSI            float64            `json:"gsi"`
	Grade          string             `json:"grade"`
	Pillars        map[string]float64 `json:"pillars"`
	Catastrophic   []string           `json:"catastrophic"`
	NA             []string           `json:"na"`
	Fingerprint    map[string]string  `json:"fingerprint"`
	GeneratedAt    string             `json:"generated_at"`
}

// FromDocument builds an Unratified (non-sample) entry from a report scorecard.
// AdapterVersion is taken from fingerprint["adapter.version"] when present.
func FromDocument(doc report.Document, framework, adapter string) Entry {
	fp := doc.Fingerprint
	if fp == nil {
		fp = map[string]string{}
	} else {
		copied := make(map[string]string, len(fp))
		for k, v := range fp {
			copied[k] = v
		}
		fp = copied
	}
	pillars := doc.PillarScores
	if pillars == nil {
		pillars = map[string]float64{}
	} else {
		copied := make(map[string]float64, len(pillars))
		for k, v := range pillars {
			copied[k] = v
		}
		pillars = copied
	}
	cat := doc.Catastrophic
	if cat == nil {
		cat = []string{}
	} else {
		cat = append([]string(nil), cat...)
	}
	na := doc.NA
	if na == nil {
		na = []string{}
	} else {
		na = append([]string(nil), na...)
	}
	return Entry{
		RunID:          doc.RunID,
		Framework:      framework,
		Adapter:        adapter,
		AdapterVersion: fp["adapter.version"],
		Provenance:     doc.Provenance,
		Tab:            TabUnratified,
		Sample:         false,
		GSI:            doc.GSI,
		Grade:          doc.Grade,
		Pillars:        pillars,
		Catastrophic:   cat,
		NA:             na,
		Fingerprint:    fp,
		GeneratedAt:    time.Now().UTC().Format(time.RFC3339),
	}
}

// Validate checks required fields and ADR 006 / ADR 007 enums.
func Validate(e Entry) error {
	if strings.TrimSpace(e.RunID) == "" {
		return fmt.Errorf("run_id is required")
	}
	if strings.TrimSpace(e.Framework) == "" {
		return fmt.Errorf("framework is required")
	}
	if strings.TrimSpace(e.Adapter) == "" {
		return fmt.Errorf("adapter is required")
	}
	if strings.TrimSpace(e.AdapterVersion) == "" {
		return fmt.Errorf("adapter_version is required")
	}
	switch e.Provenance {
	case "ratified", "provisional", "unofficial":
	default:
		return fmt.Errorf("provenance %q is invalid (want ratified|provisional|unofficial)", e.Provenance)
	}
	switch e.Tab {
	case TabOptIn, TabUnratified:
	default:
		return fmt.Errorf("tab %q is invalid (want opt-in|unratified)", e.Tab)
	}
	if e.Tab == TabOptIn && !e.Sample {
		return fmt.Errorf("tab=opt-in requires sample=true until v1.0 signatures (ADR 006)")
	}
	if strings.TrimSpace(e.Grade) == "" {
		return fmt.Errorf("grade is required")
	}
	if e.Pillars == nil {
		return fmt.Errorf("pillars is required")
	}
	if e.Catastrophic == nil {
		return fmt.Errorf("catastrophic is required")
	}
	if e.NA == nil {
		return fmt.Errorf("na is required")
	}
	if e.Fingerprint == nil {
		return fmt.Errorf("fingerprint is required")
	}
	if strings.TrimSpace(e.GeneratedAt) == "" {
		return fmt.Errorf("generated_at is required")
	}
	return nil
}

// Write writes e to <dashboard>/data/<run_id>.json and rewrites
// <dashboard>/data/index.json (sorted, deduplicated filenames).
func Write(dashboard string, e Entry) (string, error) {
	if err := Validate(e); err != nil {
		return "", err
	}
	if dashboard == "" {
		dashboard = "dashboard"
	}
	dataDir := filepath.Join(dashboard, dataDirName)
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", dataDir, err)
	}
	filename := e.RunID + ".json"
	entryPath := filepath.Join(dataDir, filename)
	body, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal entry: %w", err)
	}
	body = append(body, '\n')
	if err := os.WriteFile(entryPath, body, 0o644); err != nil {
		return "", fmt.Errorf("write %s: %w", entryPath, err)
	}
	indexPath := filepath.Join(dataDir, IndexFileName)
	entries, err := readIndex(indexPath)
	if err != nil {
		return "", err
	}
	entries = append(entries, filename)
	entries = dedupeSorted(entries)
	indexBody, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal index: %w", err)
	}
	indexBody = append(indexBody, '\n')
	if err := os.WriteFile(indexPath, indexBody, 0o644); err != nil {
		return "", fmt.Errorf("write %s: %w", indexPath, err)
	}
	abs, err := filepath.Abs(entryPath)
	if err != nil {
		return entryPath, nil
	}
	return abs, nil
}

func readIndex(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var entries []string
	if err := json.Unmarshal(b, &entries); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return entries, nil
}

func dedupeSorted(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}
