// Package report loads run artifacts and renders GSI scorecards.
package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/agentgavel/agentgavel/internal/metrics"
)

const (
	SummaryFileName   = "summary.json"
	ScorecardFileName = "scorecard.json"
)

// Document is the machine-readable scorecard for a run.
type Document struct {
	RunID         string             `json:"run_id"`
	GSI           float64            `json:"gsi"`
	Grade         string             `json:"grade"`
	PillarScores  map[string]float64 `json:"pillars"`
	Catastrophic  []string           `json:"catastrophic"`
	Observability bool               `json:"observability_penalty"`
	NA            []string           `json:"na,omitempty"`
	Fingerprint   map[string]string  `json:"fingerprint,omitempty"`
	Provenance    string             `json:"provenance,omitempty"`
}

// summaryFile is the on-disk results/<run-id>/summary.json shape.
type summaryFile struct {
	RunID                string                     `json:"run_id"`
	Fingerprint          map[string]string          `json:"fingerprint"`
	Scenarios            map[string]json.RawMessage `json:"scenarios"`
	ObservabilityPenalty bool                       `json:"observability_penalty"`
	Provenance           string                     `json:"provenance"`
}

// scenarioEntry is one scenario row inside summary.json.
type scenarioEntry struct {
	Score            float64 `json:"score"`
	NA               bool    `json:"na"`
	Catastrophic     bool    `json:"catastrophic"`
	CatastrophicKind string  `json:"catastrophic_kind"`
}

// ResolveRunDir maps a run-id or filesystem path to a results directory.
// Bare run IDs resolve under root/results/<id>. Absolute paths, relative
// paths that exist, or args containing a path separator are used as-is
// (relative to the process working directory when not absolute).
func ResolveRunDir(root, arg string) (string, error) {
	if arg == "" {
		return "", fmt.Errorf("run id or path is required")
	}
	if filepath.IsAbs(arg) || strings.Contains(arg, string(os.PathSeparator)) || strings.Contains(arg, "/") {
		return arg, nil
	}
	if st, err := os.Stat(arg); err == nil && st.IsDir() {
		return arg, nil
	}
	if root == "" {
		root = "."
	}
	return filepath.Join(root, "results", arg), nil
}

// Load reads scorecard.json when present; otherwise computes a scorecard
// from summary.json scenario rows via metrics.ComputeGSI.
func Load(dir string) (Document, error) {
	scorecardPath := filepath.Join(dir, ScorecardFileName)
	if b, err := os.ReadFile(scorecardPath); err == nil {
		var doc Document
		if err := json.Unmarshal(b, &doc); err != nil {
			return Document{}, fmt.Errorf("parse %s: %w", scorecardPath, err)
		}
		if doc.PillarScores == nil {
			doc.PillarScores = map[string]float64{}
		}
		if doc.Catastrophic == nil {
			doc.Catastrophic = []string{}
		}
		if doc.RunID == "" {
			doc.RunID = filepath.Base(filepath.Clean(dir))
		}
		return doc, nil
	} else if !os.IsNotExist(err) {
		return Document{}, fmt.Errorf("read %s: %w", scorecardPath, err)
	}

	summaryPath := filepath.Join(dir, SummaryFileName)
	b, err := os.ReadFile(summaryPath)
	if err != nil {
		return Document{}, fmt.Errorf("read %s: %w", summaryPath, err)
	}
	var sum summaryFile
	if err := json.Unmarshal(b, &sum); err != nil {
		return Document{}, fmt.Errorf("parse %s: %w", summaryPath, err)
	}
	doc, err := fromSummary(sum)
	if err != nil {
		return Document{}, err
	}
	if doc.RunID == "" {
		doc.RunID = filepath.Base(filepath.Clean(dir))
	}
	return doc, nil
}

func fromSummary(sum summaryFile) (Document, error) {
	var results []metrics.ScenarioResult
	ids := make([]string, 0, len(sum.Scenarios))
	for id := range sum.Scenarios {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		raw := sum.Scenarios[id]
		var entry scenarioEntry
		if err := json.Unmarshal(raw, &entry); err != nil {
			return Document{}, fmt.Errorf("scenario %s: %w", id, err)
		}
		kind := entry.CatastrophicKind
		if entry.Catastrophic && kind == "" {
			kind = id
		}
		results = append(results, metrics.ScenarioResult{
			ID:               id,
			Score:            entry.Score,
			NA:               entry.NA,
			Catastrophic:     entry.Catastrophic,
			CatastrophicKind: kind,
		})
	}
	sc := metrics.ComputeGSI(results, sum.ObservabilityPenalty)
	return Document{
		RunID:         sum.RunID,
		GSI:           sc.GSI,
		Grade:         sc.Grade,
		PillarScores:  sc.PillarScores,
		Catastrophic:  sc.Catastrophic,
		Observability: sc.Observability,
		NA:            sc.NA,
		Fingerprint:   sum.Fingerprint,
		Provenance:    sum.Provenance,
	}, nil
}

// FormatText renders a human-readable scorecard.
func FormatText(doc Document) string {
	var b strings.Builder
	fmt.Fprintf(&b, "AgentGavel Scorecard\n")
	if doc.RunID != "" {
		fmt.Fprintf(&b, "Run ID: %s\n", doc.RunID)
	}
	fmt.Fprintf(&b, "GSI: %.1f\n", doc.GSI)
	fmt.Fprintf(&b, "Grade: %s\n", doc.Grade)
	if doc.Provenance != "" {
		fmt.Fprintf(&b, "Provenance: %s\n", doc.Provenance)
	}
	if doc.Observability {
		fmt.Fprintf(&b, "Observability penalty: applied (GSI capped at 600)\n")
	}
	fmt.Fprintf(&b, "\nPillars:\n")
	for _, pillar := range pillarOrder() {
		score, ok := doc.PillarScores[pillar]
		if !ok {
			continue
		}
		fmt.Fprintf(&b, "  %-13s %.1f\n", pillar+":", score)
	}
	fmt.Fprintf(&b, "\nCatastrophic flags:\n")
	if len(doc.Catastrophic) == 0 {
		fmt.Fprintf(&b, "  (none)\n")
	} else {
		for _, flag := range doc.Catastrophic {
			fmt.Fprintf(&b, "  - %s\n", flag)
		}
	}
	if len(doc.NA) > 0 {
		fmt.Fprintf(&b, "\nN/A scenarios:\n")
		for _, id := range doc.NA {
			fmt.Fprintf(&b, "  - %s\n", id)
		}
	}
	if len(doc.Fingerprint) > 0 {
		fmt.Fprintf(&b, "\nFingerprint:\n")
		keys := make([]string, 0, len(doc.Fingerprint))
		for k := range doc.Fingerprint {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(&b, "  %s: %s\n", k, doc.Fingerprint[k])
		}
	}
	return b.String()
}

// FormatJSON returns indented scorecard JSON.
func FormatJSON(doc Document) (string, error) {
	if doc.PillarScores == nil {
		doc.PillarScores = map[string]float64{}
	}
	if doc.Catastrophic == nil {
		doc.Catastrophic = []string{}
	}
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}

func pillarOrder() []string {
	type pw struct {
		name   string
		weight float64
	}
	var list []pw
	for name, w := range metrics.PillarWeights {
		list = append(list, pw{name, w})
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].weight == list[j].weight {
			return list[i].name < list[j].name
		}
		return list[i].weight > list[j].weight
	})
	out := make([]string, len(list))
	for i, p := range list {
		out[i] = p.name
	}
	return out
}
