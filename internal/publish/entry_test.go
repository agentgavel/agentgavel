package publish

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/agentgavel/agentgavel/internal/report"
)

func TestEntryValidateEnums(t *testing.T) {
	base := validEntry()
	cases := []struct {
		name    string
		mutate  func(*Entry)
		wantErr string
	}{
		{
			name:    "ok_unofficial_unratified",
			mutate:  func(*Entry) {},
			wantErr: "",
		},
		{
			name: "ok_opt_in_sample",
			mutate: func(e *Entry) {
				e.Tab = TabOptIn
				e.Sample = true
				e.Provenance = "ratified"
			},
			wantErr: "",
		},
		{
			name: "bad_provenance",
			mutate: func(e *Entry) {
				e.Provenance = "signed"
			},
			wantErr: "provenance",
		},
		{
			name: "bad_tab",
			mutate: func(e *Entry) {
				e.Tab = "primary"
			},
			wantErr: "tab",
		},
		{
			name: "opt_in_without_sample",
			mutate: func(e *Entry) {
				e.Tab = TabOptIn
				e.Sample = false
			},
			wantErr: "sample=true",
		},
		{
			name: "missing_run_id",
			mutate: func(e *Entry) {
				e.RunID = ""
			},
			wantErr: "run_id",
		},
		{
			name: "missing_adapter_version",
			mutate: func(e *Entry) {
				e.AdapterVersion = ""
			},
			wantErr: "adapter_version",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := base
			tc.mutate(&e)
			err := Validate(e)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() = %v, want nil", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("Validate() = nil, want error containing %q", tc.wantErr)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("Validate() = %v, want substring %q", err, tc.wantErr)
			}
		})
	}
}

func TestFromDocumentAndWrite(t *testing.T) {
	doc := report.Document{
		RunID: "pub-run-1",
		GSI:   883.3,
		Grade: "F",
		PillarScores: map[string]float64{
			"chokepoint": 66.7,
		},
		Catastrophic: []string{"SEC-004"},
		NA:           []string{"SEC-008"},
		Fingerprint: map[string]string{
			"adapter.version": "0.1.0",
			"hash":            "abc",
		},
		Provenance: "unofficial",
	}
	e := FromDocument(doc, "FakeAdapter", "fakeadapter")
	if e.Tab != TabUnratified || e.Sample {
		t.Fatalf("entry tab/sample = %q/%v, want unratified/false", e.Tab, e.Sample)
	}
	if e.AdapterVersion != "0.1.0" {
		t.Fatalf("adapter_version = %q, want 0.1.0", e.AdapterVersion)
	}
	if e.Provenance != "unofficial" {
		t.Fatalf("provenance = %q, want unofficial", e.Provenance)
	}

	dash := t.TempDir()
	path, err := Write(dash, e)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	wantPath := filepath.Join(dash, "data", "pub-run-1.json")
	absWant, _ := filepath.Abs(wantPath)
	if path != absWant && path != wantPath {
		t.Fatalf("Write path = %q, want %q", path, absWant)
	}
	raw, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatal(err)
	}
	var got Entry
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got.RunID != "pub-run-1" || got.Tab != TabUnratified || got.Sample {
		t.Fatalf("written entry = %+v", got)
	}

	idxRaw, err := os.ReadFile(filepath.Join(dash, "data", IndexFileName))
	if err != nil {
		t.Fatal(err)
	}
	var idx []string
	if err := json.Unmarshal(idxRaw, &idx); err != nil {
		t.Fatal(err)
	}
	if len(idx) != 1 || idx[0] != "pub-run-1.json" {
		t.Fatalf("index = %#v, want [pub-run-1.json]", idx)
	}

	// Second write is idempotent for the index.
	if _, err := Write(dash, e); err != nil {
		t.Fatalf("second Write: %v", err)
	}
	idxRaw, err = os.ReadFile(filepath.Join(dash, "data", IndexFileName))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(idxRaw, &idx); err != nil {
		t.Fatal(err)
	}
	if len(idx) != 1 || idx[0] != "pub-run-1.json" {
		t.Fatalf("index after rewrite = %#v, want exactly once", idx)
	}
}

func validEntry() Entry {
	return Entry{
		RunID:          "run-x",
		Framework:      "FakeAdapter",
		Adapter:        "fakeadapter",
		AdapterVersion: "0.1.0",
		Provenance:     "unofficial",
		Tab:            TabUnratified,
		Sample:         false,
		GSI:            100,
		Grade:          "A",
		Pillars:        map[string]float64{"chokepoint": 100},
		Catastrophic:   []string{},
		NA:             []string{},
		Fingerprint:    map[string]string{"hash": "aa"},
		GeneratedAt:    "2026-09-04T00:00:00Z",
	}
}
