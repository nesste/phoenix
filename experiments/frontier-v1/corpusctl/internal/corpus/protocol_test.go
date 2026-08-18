package corpus

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProtocolV4IsUnfrozenCompleteAndBudgeted(t *testing.T) {
	root := repoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, "experiments", "frontier-v1", "protocol.json"))
	if err != nil {
		t.Fatal(err)
	}
	var protocol struct {
		Version   int    `json:"v"`
		Status    string `json:"status"`
		Frozen    bool   `json:"frozen"`
		Amendment struct {
			Activation   string `json:"activation"`
			StateChanges string `json:"state_changes"`
		} `json:"amendment"`
		Runtime struct {
			Version  string `json:"version"`
			Model    string `json:"model"`
			MaxTurns int    `json:"max_turns"`
		} `json:"runtime"`
		SystemPrompts map[string]string `json:"system_prompts"`
		Tranche       struct {
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
					Families   int            `json:"families"`
					Counts     map[string]int `json:"counts_across_group"`
					Constraint string         `json:"constraint"`
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
			RetryPolicy  struct {
				Ineligible string `json:"ineligible"`
			} `json:"retry_policy"`
		} `json:"trial"`
		Arms struct {
			A struct {
				Orientation any  `json:"orientation"`
				StateEvents bool `json:"state_events"`
			} `json:"A"`
			B struct {
				Documentation string `json:"documentation"`
				Orientation   any    `json:"orientation"`
				StateEvents   bool   `json:"state_events"`
			} `json:"B"`
			C struct {
				Orientation any  `json:"orientation"`
				StateEvents bool `json:"state_events"`
			} `json:"C"`
			D struct {
				Orientation any  `json:"orientation"`
				StateEvents bool `json:"state_events"`
			} `json:"D"`
			E struct {
				Orientation any  `json:"orientation"`
				StateEvents bool `json:"state_events"`
			} `json:"E"`
		} `json:"arms"`
		Analysis struct {
			Interval     string `json:"interval_method"`
			Ratio        string `json:"ratio_method"`
			Multiplicity string `json:"multiplicity"`
			Confidence   string `json:"confidence"`
		} `json:"analysis"`
		Claims []struct {
			ID                      string  `json:"id"`
			Population              string  `json:"population"`
			Endpoint                string  `json:"endpoint"`
			Uncertainty             string  `json:"uncertainty"`
			Threshold               string  `json:"threshold"`
			MinimumDetectableEffect float64 `json:"minimum_detectable_effect"`
			PowerNote               string  `json:"power_note"`
			Multiplicity            string  `json:"multiplicity"`
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
			Record   struct {
				ReviewerRole string `json:"reviewer_role"`
				Date         string `json:"date"`
				Verdict      string `json:"verdict"`
				Candidate    struct {
					Protocol string `json:"protocol_json"`
					Decision string `json:"decision_0004"`
				} `json:"candidate_lf_normalized_utf8_sha256"`
				AcceptedLimitations        []string `json:"accepted_limitations"`
				RemainingExecutionBlockers []string `json:"remaining_execution_blockers"`
			} `json:"review_record"`
		} `json:"review"`
		Unresolved []string `json:"unresolved_before_freeze"`
		Gate       struct {
			MayRunAuthoring   bool   `json:"may_run_authoring"`
			MayOpenValidation bool   `json:"may_open_validation"`
			MayOpenHeldOut    bool   `json:"may_open_held_out"`
			Reason            string `json:"reason"`
		} `json:"gate"`
	}
	if err := json.Unmarshal(data, &protocol); err != nil {
		t.Fatal(err)
	}
	if protocol.Version != 4 || protocol.Status != "authoring_amendment" || protocol.Frozen || protocol.Review.Accepted {
		t.Fatalf("unexpected protocol state: v=%d status=%s frozen=%t", protocol.Version, protocol.Status, protocol.Frozen)
	}
	if !strings.Contains(protocol.Amendment.Activation, "Before the session's first executable act") || !strings.Contains(protocol.Amendment.Activation, "must not reactivate from intent") || !strings.Contains(protocol.Amendment.Activation, "current pending frontier") {
		t.Fatalf("bootstrap activation or current-frontier state scope regressed: %q", protocol.Amendment.Activation)
	}
	if !strings.Contains(protocol.Amendment.StateChanges, "shared trial-harness rule for every arm") || !strings.Contains(protocol.Amendment.StateChanges, "flat-tool Arms A and B") || !strings.Contains(protocol.Amendment.StateChanges, "Orientations do not advance") {
		t.Fatalf("cross-arm state-event contract regressed: %q", protocol.Amendment.StateChanges)
	}
	basePrompt := "Use only the configured tools. Follow the user request exactly."
	intentInstruction := "Before the first executable act, send the complete goal as intent"
	for _, arm := range []string{"A", "B", "C", "D", "E"} {
		prompt, exists := protocol.SystemPrompts[arm]
		if !exists || !strings.HasPrefix(prompt, basePrompt) {
			t.Fatalf("arm %s system prompt is not pinned: %q", arm, prompt)
		}
		hasIntent := strings.Contains(prompt, intentInstruction)
		wantIntent := arm == "C" || arm == "D" || arm == "E"
		if hasIntent != wantIntent {
			t.Fatalf("arm %s intent instruction = %t, want %t: %q", arm, hasIntent, wantIntent, prompt)
		}
	}
	for arm, stateEvents := range map[string]bool{
		"A": protocol.Arms.A.StateEvents, "B": protocol.Arms.B.StateEvents,
		"C": protocol.Arms.C.StateEvents, "D": protocol.Arms.D.StateEvents, "E": protocol.Arms.E.StateEvents,
	} {
		if !stateEvents {
			t.Fatalf("arm %s does not share harness state events", arm)
		}
	}
	if protocol.Arms.A.Orientation != false || protocol.Arms.B.Orientation != false {
		t.Fatalf("flat arms unexpectedly expose orientation: A=%v B=%v", protocol.Arms.A.Orientation, protocol.Arms.B.Orientation)
	}
	for arm, orientation := range map[string]any{"C": protocol.Arms.C.Orientation, "D": protocol.Arms.D.Orientation, "E": protocol.Arms.E.Orientation} {
		value, ok := orientation.(string)
		if !ok || !strings.Contains(value, "bootstrap_only") {
			t.Fatalf("arm %s orientation is not bootstrap-only: %v", arm, orientation)
		}
	}
	record := protocol.Review.Record
	if record.ReviewerRole != "evaluation reviewer, independent of implementation" || record.Date != "2026-08-17" || record.Verdict != "ACCEPT" {
		t.Fatalf("unexpected independent acceptance record: %+v", record)
	}
	if record.Candidate.Protocol != "sha256:4d13c140742ba462d68af865d4ed12b08e6c2e1aba4f652fed4e35adf782a47d" || record.Candidate.Decision != "sha256:399d337ef957ff1a43e7b384a03a7ddfe1e3a360b9f5337ae1d662338dfe3849" {
		t.Fatalf("acceptance record identifies the wrong candidate: %+v", record.Candidate)
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
	if !strings.Contains(allocation.Mix.Constraint, "at least 1 cascade case") || !strings.Contains(allocation.Mix.Constraint, "G=16") {
		t.Fatalf("mix allocation must force cascade across all eight mix families and G=16: %q", allocation.Mix.Constraint)
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
	if !strings.Contains(protocol.Analysis.Interval, "G>=20") || !strings.Contains(protocol.Analysis.Interval, "hierarchical bootstrap") || !strings.Contains(protocol.Analysis.Interval, "G<20") || !strings.Contains(protocol.Analysis.Interval, "sign-flip permutation") || !strings.Contains(protocol.Analysis.Ratio, "fewer than 10 successful trials") || !strings.Contains(protocol.Analysis.Multiplicity, "Single Phase 1 primary claim") {
		t.Fatalf("analysis package is incomplete: %+v", protocol.Analysis)
	}
	if !strings.Contains(protocol.Analysis.Confidence, "directional claims and the efficiency-only success guardrail") || strings.Contains(protocol.Analysis.Confidence, "noninferiority claims") {
		t.Fatalf("confidence statement misclassifies a harm gate as noninferiority: %q", protocol.Analysis.Confidence)
	}
	if !strings.Contains(protocol.Trial.RetryPolicy.Ineligible, "turn-limit hit") || !strings.Contains(protocol.Trial.RetryPolicy.Ineligible, "cost-cap hit") {
		t.Fatalf("turn-limit and cost-cap hits must be retry-ineligible: %q", protocol.Trial.RetryPolicy.Ineligible)
	}
	armBDoc := protocol.Arms.B.Documentation
	for _, required := range []string{"(1) name every starter-world verb", "(2) include one worked path per authoring class", "(3) include a recovery recipe", "(4) use authoring evidence only", "(5) receive a human-factors completeness review", "(6) freeze by digest", "same verb implementations, typed payloads, side effects, limits, and surface-only access as Arm A"} {
		if !strings.Contains(armBDoc, required) {
			t.Fatalf("Arm B documentation is missing %q: %q", required, armBDoc)
		}
	}
	claims := map[string]struct {
		Population              string
		Endpoint                string
		Uncertainty             string
		Threshold               string
		MinimumDetectableEffect float64
		PowerNote               string
		Multiplicity            string
	}{}
	for _, claim := range protocol.Claims {
		if claim.ID == "" || claim.Uncertainty == "" || claim.Threshold == "" || claim.MinimumDetectableEffect <= 0 || claim.ForcedFailure == "" {
			t.Fatalf("incomplete claim: %+v", claim)
		}
		if _, exists := claims[claim.ID]; exists {
			t.Fatalf("duplicate claim %s", claim.ID)
		}
		claims[claim.ID] = struct {
			Population              string
			Endpoint                string
			Uncertainty             string
			Threshold               string
			MinimumDetectableEffect float64
			PowerNote               string
			Multiplicity            string
		}{claim.Population, claim.Endpoint, claim.Uncertainty, claim.Threshold, claim.MinimumDetectableEffect, claim.PowerNote, claim.Multiplicity}
	}
	for _, required := range []string{"direct_no_tax", "headline_value", "frontier_value", "teaching_refusal_value", "learning_value"} {
		if _, exists := claims[required]; !exists {
			t.Fatalf("missing claim %s", required)
		}
	}
	direct := claims["direct_no_tax"]
	if direct.MinimumDetectableEffect != 0.39 || !strings.Contains(direct.PowerNote, "80% harm-detection MDE is 0.39") || !strings.Contains(direct.PowerNote, "0.30 effect size is not an 80% MDE") {
		t.Fatalf("direct harm-gate sensitivity is mislabeled: %+v", direct)
	}
	headline := claims["headline_value"]
	if !strings.Contains(headline.Threshold, "lower bound for C-B ITT success is greater than 0") || !strings.Contains(headline.Threshold, "paired-success C/B token-ratio") || !strings.Contains(headline.Threshold, "dead-end-rate") || !strings.Contains(headline.Threshold, "timeout-or-cap-hit-rate") {
		t.Fatalf("headline and efficiency-only predicates regressed: %+v", headline)
	}
	frontier := claims["frontier_value"]
	if !strings.Contains(frontier.Population, "16 generating families") || !strings.Contains(frontier.Uncertainty, "sign-flip permutation test because G=16") || frontier.MinimumDetectableEffect != 0.26 {
		t.Fatalf("frontier allocation or sensitivity regressed: %+v", frontier)
	}
	teaching := claims["teaching_refusal_value"]
	if !strings.Contains(teaching.Population, "all assigned recovery trials") || !strings.Contains(teaching.Population, "D teaching refusals and E plain typed errors are both eligible") || !strings.Contains(teaching.Endpoint, "first state-sensitive act fails or is declined") {
		t.Fatalf("teaching-refusal ITT population or trigger regressed: %+v", teaching)
	}
	for _, isolation := range []string{"direct_no_tax", "frontier_value", "teaching_refusal_value"} {
		if !strings.Contains(claims[isolation].Multiplicity, "cannot create a headline pass") {
			t.Fatalf("isolation %s may not create a headline pass: %q", isolation, claims[isolation].Multiplicity)
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
	if protocol.Phase2.OnlineUpdates || len(protocol.Phase2.Arms) != 3 || protocol.Phase2.Learning == "" || protocol.Phase2.Replication == "" {
		t.Fatalf("Phase 2 contract is incomplete: %+v", protocol.Phase2)
	}
	wantPhase2Frozen := []string{"daemon", "verbs", "authored_rules", "candidate_weights", "training_episode_digest"}
	if len(protocol.Phase2.Frozen) != len(wantPhase2Frozen) {
		t.Fatalf("unexpected Phase 2 frozen artifacts: %v", protocol.Phase2.Frozen)
	}
	for i, want := range wantPhase2Frozen {
		if protocol.Phase2.Frozen[i] != want {
			t.Fatalf("Phase 2 frozen artifact %d is %q, want %q", i, protocol.Phase2.Frozen[i], want)
		}
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
	if len(protocol.Unresolved) != 2 || !strings.Contains(protocol.Unresolved[0], "protocol v4") || !strings.Contains(protocol.Unresolved[1], "regeneration") {
		t.Fatalf("amended protocol must retain review and resealing blockers: %v", protocol.Unresolved)
	}
	wantLimitations := []string{
		"Conclusions are limited to Claude Code 2.1.229, claude-sonnet-5, the frozen dev-repo world, and the eight task classes.",
		"Expected headline power is about 62% at a 15-point effect; the approximate 80% MDE against zero is 0.19.",
		"The USD 300 validation ceiling cannot support 80% power for a 15-point Holm-adjusted LCB-floor design.",
		"Efficiency-only is a conservative conjunction; at a true zero difference, power to clear the 0.05 success or dead-end bounds is about 0.16.",
	}
	wantExecutionBlockers := []string{
		"independent unopened validation and held_out families",
		"Task 0.6 Phase 0 review",
	}
	for label, pair := range map[string][2][]string{
		"accepted limitations":         {record.AcceptedLimitations, wantLimitations},
		"remaining execution blockers": {record.RemainingExecutionBlockers, wantExecutionBlockers},
	} {
		if len(pair[0]) != len(pair[1]) {
			t.Fatalf("unexpected %s: %v", label, pair[0])
		}
		for i, want := range pair[1] {
			if pair[0][i] != want {
				t.Fatalf("%s item %d is %q, want %q", label, i, pair[0][i], want)
			}
		}
	}
	if !protocol.Gate.MayRunAuthoring || protocol.Gate.MayOpenValidation || protocol.Gate.MayOpenHeldOut || !strings.Contains(protocol.Gate.Reason, "Protocol v4") || !strings.Contains(protocol.Gate.Reason, "independently accepted") {
		t.Fatalf("v4 amendment must keep only authoring open: %+v", protocol.Gate)
	}
}
