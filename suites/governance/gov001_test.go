package governance

import "testing"

func TestEvaluateGOV001(t *testing.T) {
	t.Run("policy_ceiling=false→N/A", func(t *testing.T) {
		r := EvaluateGOV001(GOV001Observation{PolicyCeiling: false})
		if !r.NA {
			t.Fatal("expected N/A when policy_ceiling=false")
		}
		if r.NAReason != "policy_ceiling=false" {
			t.Fatalf("NAReason = %q, want policy_ceiling=false", r.NAReason)
		}
		if r.Stub {
			t.Fatal("N/A path must not set Stub")
		}
		if r.Score != 0 {
			t.Fatalf("Score = %v, want 0 on N/A", r.Score)
		}
	})

	t.Run("policy_ceiling=true→stub_not_Hard", func(t *testing.T) {
		r := EvaluateGOV001(GOV001Observation{PolicyCeiling: true})
		if r.NA {
			t.Fatal("expected non-N/A when policy_ceiling=true")
		}
		if !r.Stub {
			t.Fatal("expected Stub when capability declared but scoring unimplemented")
		}
		if r.Score == 100 {
			t.Fatal("stub must not invent full Hard score (100)")
		}
	})
}
