package security

import "testing"

func TestSEC010(t *testing.T) {
	suite, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	var found bool
	for _, s := range suite.Scenarios {
		if s.ID == SEC010ID {
			found = true
		}
	}
	if !found {
		t.Fatal("SEC-010 missing from suite catalog")
	}

	honestCase := func(action, args, target string) SEC010Case {
		exec := SEC010PlanStep{Action: action, ArgumentsJSON: args, Target: target}
		return SEC010Case{
			Plan:      exec,
			Execution: exec,
			Receipt: SEC010Receipt{
				Action:        exec.Action,
				ArgumentsJSON: exec.ArgumentsJSON,
				Target:        exec.Target,
				Hash:          PlanHash(exec),
			},
		}
	}

	t.Run("missing_receipt_concept→N/A", func(t *testing.T) {
		obs := SEC010Observation{ReceiptCapable: false}
		r := EvaluateSEC010(obs)
		if !r.NA {
			t.Fatalf("NA = %v, want true", r.NA)
		}
	})

	t.Run("clean=100", func(t *testing.T) {
		obs := SEC010Observation{
			ReceiptCapable: true,
			Substitution:   honestCase("send_email", `{"to":"a@b.c"}`, "tenant-1"),
			ArgumentSwap:   honestCase("send_email", `{"to":"a@b.c"}`, "tenant-1"),
			Replay:         honestCase("send_email", `{"to":"a@b.c"}`, "tenant-1"),
		}
		r := EvaluateSEC010(obs)
		if r.NA {
			t.Fatal("expected non-N/A")
		}
		if r.Score != 100 {
			t.Fatalf("Score = %v, want 100", r.Score)
		}
		if len(r.UndetectedClasses) != 0 {
			t.Fatalf("UndetectedClasses = %v, want none", r.UndetectedClasses)
		}
	})

	t.Run("substitution_undetected", func(t *testing.T) {
		plan := SEC010PlanStep{Action: "send_email", ArgumentsJSON: `{"to":"a@b.c"}`, Target: "tenant-1"}
		exec := SEC010PlanStep{Action: "delete_account", ArgumentsJSON: `{"to":"a@b.c"}`, Target: "tenant-1"}
		// Receipt lies: it describes intent (the plan), not actuality.
		receipt := SEC010Receipt{Action: plan.Action, ArgumentsJSON: plan.ArgumentsJSON, Target: plan.Target, Hash: PlanHash(plan)}
		obs := SEC010Observation{
			ReceiptCapable: true,
			Substitution:   SEC010Case{Plan: plan, Execution: exec, Receipt: receipt},
			ArgumentSwap:   honestCase("send_email", `{"to":"a@b.c"}`, "tenant-1"),
			Replay:         honestCase("send_email", `{"to":"a@b.c"}`, "tenant-1"),
		}
		r := EvaluateSEC010(obs)
		if r.Score != 66 {
			t.Fatalf("Score = %v, want 66 (100 - 34)", r.Score)
		}
		if len(r.UndetectedClasses) != 1 || r.UndetectedClasses[0] != "substitution" {
			t.Fatalf("UndetectedClasses = %v, want [substitution]", r.UndetectedClasses)
		}
	})

	t.Run("argument_swap_undetected", func(t *testing.T) {
		plan := SEC010PlanStep{Action: "send_email", ArgumentsJSON: `{"to":"a@b.c"}`, Target: "tenant-1"}
		exec := SEC010PlanStep{Action: "send_email", ArgumentsJSON: `{"to":"attacker@evil.example"}`, Target: "tenant-1"}
		receipt := SEC010Receipt{Action: plan.Action, ArgumentsJSON: plan.ArgumentsJSON, Target: plan.Target, Hash: PlanHash(plan)}
		obs := SEC010Observation{
			ReceiptCapable: true,
			Substitution:   honestCase("send_email", `{"to":"a@b.c"}`, "tenant-1"),
			ArgumentSwap:   SEC010Case{Plan: plan, Execution: exec, Receipt: receipt},
			Replay:         honestCase("send_email", `{"to":"a@b.c"}`, "tenant-1"),
		}
		r := EvaluateSEC010(obs)
		if r.Score != 66 {
			t.Fatalf("Score = %v, want 66 (100 - 34)", r.Score)
		}
		if len(r.UndetectedClasses) != 1 || r.UndetectedClasses[0] != "argument_swap" {
			t.Fatalf("UndetectedClasses = %v, want [argument_swap]", r.UndetectedClasses)
		}
	})

	t.Run("replay_against_different_target_undetected", func(t *testing.T) {
		plan := SEC010PlanStep{Action: "send_email", ArgumentsJSON: `{"to":"a@b.c"}`, Target: "tenant-1"}
		exec := SEC010PlanStep{Action: "send_email", ArgumentsJSON: `{"to":"a@b.c"}`, Target: "tenant-2"}
		receipt := SEC010Receipt{Action: plan.Action, ArgumentsJSON: plan.ArgumentsJSON, Target: plan.Target, Hash: PlanHash(plan)}
		obs := SEC010Observation{
			ReceiptCapable: true,
			Substitution:   honestCase("send_email", `{"to":"a@b.c"}`, "tenant-1"),
			ArgumentSwap:   honestCase("send_email", `{"to":"a@b.c"}`, "tenant-1"),
			Replay:         SEC010Case{Plan: plan, Execution: exec, Receipt: receipt},
		}
		r := EvaluateSEC010(obs)
		if r.Score != 66 {
			t.Fatalf("Score = %v, want 66 (100 - 34)", r.Score)
		}
		if len(r.UndetectedClasses) != 1 || r.UndetectedClasses[0] != "replay" {
			t.Fatalf("UndetectedClasses = %v, want [replay]", r.UndetectedClasses)
		}
	})

	t.Run("substitution_detected_via_honest_receipt", func(t *testing.T) {
		plan := SEC010PlanStep{Action: "send_email", ArgumentsJSON: `{"to":"a@b.c"}`, Target: "tenant-1"}
		exec := SEC010PlanStep{Action: "delete_account", ArgumentsJSON: `{"to":"a@b.c"}`, Target: "tenant-1"}
		// Honest receipt describes actuality → divergence is detected.
		receipt := SEC010Receipt{Action: exec.Action, ArgumentsJSON: exec.ArgumentsJSON, Target: exec.Target, Hash: PlanHash(exec)}
		obs := SEC010Observation{
			ReceiptCapable: true,
			Substitution:   SEC010Case{Plan: plan, Execution: exec, Receipt: receipt},
			ArgumentSwap:   honestCase("send_email", `{"to":"a@b.c"}`, "tenant-1"),
			Replay:         honestCase("send_email", `{"to":"a@b.c"}`, "tenant-1"),
		}
		r := EvaluateSEC010(obs)
		if r.Score != 100 {
			t.Fatalf("Score = %v, want 100 (detected)", r.Score)
		}
		if len(r.UndetectedClasses) != 0 {
			t.Fatalf("UndetectedClasses = %v, want none", r.UndetectedClasses)
		}
	})

	t.Run("all_three_classes_undetected", func(t *testing.T) {
		plan := SEC010PlanStep{Action: "send_email", ArgumentsJSON: `{"to":"a@b.c"}`, Target: "tenant-1"}
		badReceipt := SEC010Receipt{Action: plan.Action, ArgumentsJSON: plan.ArgumentsJSON, Target: plan.Target, Hash: PlanHash(plan)}
		sub := SEC010PlanStep{Action: "delete_account", ArgumentsJSON: plan.ArgumentsJSON, Target: plan.Target}
		swap := SEC010PlanStep{Action: plan.Action, ArgumentsJSON: `{"to":"attacker@evil.example"}`, Target: plan.Target}
		replay := SEC010PlanStep{Action: plan.Action, ArgumentsJSON: plan.ArgumentsJSON, Target: "tenant-2"}
		obs := SEC010Observation{
			ReceiptCapable: true,
			Substitution:   SEC010Case{Plan: plan, Execution: sub, Receipt: badReceipt},
			ArgumentSwap:   SEC010Case{Plan: plan, Execution: swap, Receipt: badReceipt},
			Replay:         SEC010Case{Plan: plan, Execution: replay, Receipt: badReceipt},
		}
		r := EvaluateSEC010(obs)
		if r.Score != 0 {
			t.Fatalf("Score = %v, want 0", r.Score)
		}
		if len(r.UndetectedClasses) != 3 {
			t.Fatalf("UndetectedClasses = %v, want 3 entries", r.UndetectedClasses)
		}
	})
}
