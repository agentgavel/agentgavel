package protocol

import "testing"

func TestCapabilityNAMapping(t *testing.T) {
	cases := []struct {
		name string
		c    CapabilityReport
		want []string
		pen  bool
	}{
		{
			name: "full",
			c: CapabilityReport{
				HITL: true, Tenancy: true, Ledger: true, Observability: true,
				PolicyCeiling: true, ContextMode: "raw",
			},
			want: nil,
			pen:  false,
		},
		{
			name: "no hitl",
			c:    CapabilityReport{Observability: true, ContextMode: "raw"},
			want: []string{"SEC-002", "SEC-005", "SEC-006", "SEC-008", "SEC-009", "SEC-010", "GOV-001", "REL-001", "REL-002", "REL-003"},
			pen:  false,
		},
		{
			name: "no tenancy ledger",
			c: CapabilityReport{
				HITL: true, Observability: true, ContextMode: "attestation",
			},
			want: []string{"SEC-008", "SEC-009", "SEC-010", "GOV-001", "REL-002", "REL-003"},
			pen:  false,
		},
		{
			name: "no observability",
			c: CapabilityReport{
				HITL: true, Tenancy: true, Ledger: true, PolicyCeiling: true, ContextMode: "raw",
			},
			want: nil,
			pen:  true,
		},
		{
			name: "no context",
			c: CapabilityReport{
				HITL: true, Tenancy: true, Ledger: true, Observability: true,
				PolicyCeiling: true, ContextMode: "none",
			},
			want: []string{"SEC-004"},
			pen:  false,
		},
		{
			name: "no policy ceiling",
			c: CapabilityReport{
				HITL: true, Tenancy: true, Ledger: true, Observability: true, ContextMode: "raw",
			},
			want: []string{"GOV-001"},
			pen:  false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ScenarioNA(tc.c)
			for _, id := range tc.want {
				if _, ok := got[id]; !ok {
					t.Errorf("missing N/A for %s in %#v", id, got)
				}
			}
			if len(got) != len(tc.want) {
				t.Errorf("got %#v want %v", got, tc.want)
			}
			if ObservabilityPenalty(tc.c) != tc.pen {
				t.Errorf("penalty=%v want %v", ObservabilityPenalty(tc.c), tc.pen)
			}
		})
	}
}

// TestScenarioNAReliability is T14.17: hitl=false drives REL-001 N/A and
// ledger=false drives REL-002/REL-003 N/A (ADR 010), matching the SEC-009/010
// convention of reusing Ledger for the receipt/binding concept.
func TestScenarioNAReliability(t *testing.T) {
	t.Run("hitl=false→REL-001", func(t *testing.T) {
		got := ScenarioNA(CapabilityReport{Ledger: true})
		if _, ok := got["REL-001"]; !ok {
			t.Fatalf("expected REL-001 N/A when hitl=false, got %#v", got)
		}
	})
	t.Run("ledger=false→REL-002,REL-003", func(t *testing.T) {
		got := ScenarioNA(CapabilityReport{HITL: true})
		for _, id := range []string{"REL-002", "REL-003"} {
			if _, ok := got[id]; !ok {
				t.Fatalf("expected %s N/A when ledger=false, got %#v", id, got)
			}
		}
	})
	t.Run("hitl=true ledger=true→no REL N/A", func(t *testing.T) {
		got := ScenarioNA(CapabilityReport{
			HITL: true, Tenancy: true, Ledger: true, PolicyCeiling: true, ContextMode: "raw",
		})
		for _, id := range []string{"REL-001", "REL-002", "REL-003"} {
			if _, ok := got[id]; ok {
				t.Fatalf("unexpected N/A for %s: %#v", id, got)
			}
		}
	})
}
