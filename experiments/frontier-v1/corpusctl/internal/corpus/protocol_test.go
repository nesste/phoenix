package corpus

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProtocolCandidateIsCompleteAndBudgeted(t *testing.T) {
	root := repoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, "experiments", "frontier-v1", "protocol.json"))
	if err != nil {
		t.Fatal(err)
	}
	var protocol struct {
		Version int    `json:"v"`
		Status  string `json:"status"`
		Frozen  bool   `json:"frozen"`
		Runtime struct {
			Version  string `json:"version"`
			Model    string `json:"model"`
			MaxTurns int    `json:"max_turns"`
		} `json:"runtime"`
		Tranche struct {
			MinimumCases        int `json:"minimum_cases_per_tranche"`
			MinimumFamilies     int `json:"minimum_families_per_class"`
			InferentialFamilies int `json:"inferential_subset_minimum_families"`
			FrontierFamilies    int `json:"frontier_subset_minimum_families"`
			Allocation          struct {
				Families       int `json:"families"`
				CasesPerFamily int `json:"cases_per_family"`
				Direct         struct {
					Families  int            `json:"families"`
					PerFamily map[string]int `json:"per_family"`
				} `json:"group_direct"`
				Recovery struct {
					Families  int            `json:"families"`
					PerFamily map[string]int `json:"per_family"`
				} `json:"group_recovery"`
				Mix struct {
					Families int            `json:"families"`
					Counts   map[string]int `json:"counts_across_group"`
				} `json:"group_mix"`
			} `json:"family_allocation"`
			Power struct {
				TargetAtFifteen float64 `json:"target_power_at_15_point_headline"`
				MDEEighty       float64 `json:"eighty_percent_mde_headline_reject_zero"`
				PrimaryClaims   int     `json:"phase_1_primary_claims"`
				RepetitionICC   float64 `json:"repetition_intraclass_correlation"`
			} `json:"power_assumptions"`
		} `json:"tranche_design"`
		Trial struct {
			Repetitions  int     `json:"repetitions_per_case_arm"`
			CostCap      float64 `json:"max_cost_usd_per_trial"`
			LaunchOrder  string  `json:"launch_order"`
			BudgetStop   string  `json:"budget_exhaustion"`
			CapImbalance string  `json:"cap_imbalance"`
			Denominator  string  `json:"unresolved_denominator"`
		} `json:"trial"`
		Analysis struct {
			Interval     string `json:"interval_method"`
			Ratio        string `json:"ratio_method"`
			Multiplicity string `json:"multiplicity"`
		} `json:"analysis"`
		Claims []struct {
			ID                      string  `json:"id"`
			Uncertainty             string  `json:"uncertainty"`
			Threshold               string  `json:"threshold"`
			MinimumDetectableEffect float64 `json:"minimum_detectable_effect"`
			ForcedFailure           string  `json:"forced_failure"`
		} `json:"claims"`
		Budget struct {
			Validation float64 `json:"validation_usd"`
			HeldOut    float64 `json:"held_out_usd"`
			HardTotal  float64 `json:"hard_total_usd"`
		} `json:"budget"`
		Phase2 struct {
			OnlineUpdates bool     `json:"online_updates"`
			Frozen        []string `json:"frozen"`
			Arms          []string `json:"arms"`
			Learning      string   `json:"learning_claim"`
			Replication   string   `json:"headline_replication"`
		} `json:"phase_2_eval"`
		ArtifactFreeze struct {
			BeforeValidation []string `json:"before_validation"`
		} `json:"artifact_freeze"`
		Review struct {
			Accepted bool `json:"accepted"`
		} `json:"review"`
		Unresolved []string `json:"unresolved_before_freeze"`
	}
	if err := json.Unmarshal(data, &protocol); err != nil {
		t.Fatal(err)
	}
	if protocol.Version != 3 || protocol.Status != "review_candidate" || protocol.Frozen || protocol.Review.Accepted {
		t.Fatalf("unexpected protocol state: v=%d status=%s frozen=%t", protocol.Version, protocol.Status, protocol.Frozen)
	}
	if protocol.Runtime.Version == "" || protocol.Runtime.Model == "" || protocol.Runtime.MaxTurns <= 0 {
		t.Fatal("runtime is not fully pinned")
	}
	if protocol.Trial.Repetitions <= 0 || protocol.Trial.CostCap <= 0 || protocol.Tranche.MinimumCases <= 0 {
		t.Fatal("trial size, repetition count, and cost cap must be positive")
	}
	if protocol.Tranche.MinimumFamilies != 8 || protocol.Tranche.InferentialFamilies != 8 || protocol.Tranche.FrontierFamilies != 12 {
		t.Fatalf("unexpected family minima: class=%d inferential=%d frontier=%d", protocol.Tranche.MinimumFamilies, protocol.Tranche.InferentialFamilies, protocol.Tranche.FrontierFamilies)
	}
	allocation := protocol.Tranche.Allocation
	if allocation.Families != 24 || allocation.CasesPerFamily != 5 || allocation.Direct.Families != 8 || allocation.Recovery.Families != 8 || allocation.Mix.Families != 8 {
		t.Fatalf("unexpected family allocation: %+v", allocation)
	}
	classTotals := map[string]int{}
	for class, count := range allocation.Direct.PerFamily {
		classTotals[class] += allocation.Direct.Families * count
	}
	for class, count := range allocation.Recovery.PerFamily {
		classTotals[class] += allocation.Recovery.Families * count
	}
	for class, count := range allocation.Mix.Counts {
		classTotals[class] += count
	}
	wantTotals := map[string]int{"direct": 24, "recovery": 24, "cascade": 12, "far_discovery": 12, "temptation": 12, "absence": 12, "stale_frontier": 12, "adversarial_text": 12}
	for class, want := range wantTotals {
		if classTotals[class] != want {
			t.Fatalf("class %s has %d cases, want %d", class, classTotals[class], want)
		}
	}
	if protocol.Tranche.Power.TargetAtFifteen != 0.62 || protocol.Tranche.Power.MDEEighty != 0.19 || protocol.Tranche.Power.PrimaryClaims != 1 || protocol.Tranche.Power.RepetitionICC != 1 {
		t.Fatalf("unexpected power assumptions: %+v", protocol.Tranche.Power)
	}
	if protocol.Trial.LaunchOrder == "" || protocol.Trial.BudgetStop == "" || protocol.Trial.CapImbalance == "" || protocol.Trial.Denominator == "" {
		t.Fatal("launch, budget-stop, cap-imbalance, and denominator rules must be explicit")
	}
	if !strings.Contains(protocol.Analysis.Interval, "sign-flip permutation") || !strings.Contains(protocol.Analysis.Ratio, "fewer than 10 successful trials") || !strings.Contains(protocol.Analysis.Multiplicity, "Single Phase 1 primary claim") {
		t.Fatalf("analysis package is incomplete: %+v", protocol.Analysis)
	}
	seen := map[string]bool{}
	for _, claim := range protocol.Claims {
		if claim.ID == "" || claim.Uncertainty == "" || claim.Threshold == "" || claim.MinimumDetectableEffect <= 0 || claim.ForcedFailure == "" {
			t.Fatalf("incomplete claim: %+v", claim)
		}
		if seen[claim.ID] {
			t.Fatalf("duplicate claim %s", claim.ID)
		}
		seen[claim.ID] = true
	}
	for _, required := range []string{"direct_no_tax", "headline_value", "frontier_value", "teaching_refusal_value", "learning_value"} {
		if !seen[required] {
			t.Fatalf("missing claim %s", required)
		}
	}
	validationMaximum := float64(protocol.Tranche.MinimumCases*protocol.Trial.Repetitions*5) * protocol.Trial.CostCap * 1.10
	heldOutMaximum := float64(protocol.Tranche.MinimumCases*protocol.Trial.Repetitions*3) * protocol.Trial.CostCap * 1.10
	if validationMaximum > protocol.Budget.Validation+0.000001 {
		t.Fatalf("validation ceiling %.2f does not cover planned maximum %.2f", protocol.Budget.Validation, validationMaximum)
	}
	if heldOutMaximum > protocol.Budget.HeldOut+0.000001 {
		t.Fatalf("held_out ceiling %.2f does not cover planned maximum %.2f", protocol.Budget.HeldOut, heldOutMaximum)
	}
	if protocol.Budget.HardTotal != 565 {
		t.Fatalf("unexpected total budget %.2f", protocol.Budget.HardTotal)
	}
	if protocol.Phase2.OnlineUpdates || len(protocol.Phase2.Frozen) != 5 || len(protocol.Phase2.Arms) != 3 || protocol.Phase2.Learning == "" || protocol.Phase2.Replication == "" {
		t.Fatalf("Phase 2 contract is incomplete: %+v", protocol.Phase2)
	}
	hasScheduleDigest := false
	for _, artifact := range protocol.ArtifactFreeze.BeforeValidation {
		if artifact == "schedule digest" {
			hasScheduleDigest = true
		}
	}
	if !hasScheduleDigest {
		t.Fatal("schedule digest must be frozen explicitly")
	}
	if len(protocol.Unresolved) != 1 || protocol.Unresolved[0] != "independent evaluation reviewer acceptance" {
		t.Fatalf("freeze blockers must name only independent acceptance: %v", protocol.Unresolved)
	}
}
