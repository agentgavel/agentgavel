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
			want: []string{"SEC-002", "SEC-005", "SEC-006", "SEC-008", "SEC-009", "SEC-010", "GOV-001"},
			pen:  false,
		},
		{
			name: "no tenancy ledger",
			c: CapabilityReport{
				HITL: true, Observability: true, ContextMode: "attestation",
			},
			want: []string{"SEC-008", "SEC-009", "SEC-010", "GOV-001"},
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
