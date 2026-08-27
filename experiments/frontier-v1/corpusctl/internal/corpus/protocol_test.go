package corpus

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProtocolV5IsFrozenAcceptedCompleteAndBudgeted(t *testing.T) {
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
			Discovery    string `json:"discovery"`
			Envelope     string `json:"envelope"`
			Grading      string `json:"grading"`
			StateChanges string `json:"state_changes"`
			Retirement   string `json:"retirement"`
			SupersededV4 struct {
				AppliesTo string `json:"applies_to"`
			} `json:"superseded_v4_amendment"`
		} `json:"amendment"`
		Runtime struct {
			Version  string `json:"version"`
			Model    string `json:"model"`
			MaxTurns int    `json:"max_turns"`
		} `json:"runtime"`
		SystemPrompts map[string]string `json:"system_prompts"`
		Tranche       struct {
			ArmsPerPairingKey   int `json:"arms_per_pairing_key"`
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
		} `json:"arms"`
		Analysis struct {
			Interval     string `json:"interval_method"`
			Ratio        string `json:"ratio_method"`
			Multiplicity string `json:"multiplicity"`
			Confidence   string `json:"confidence"`
		} `json:"analysis"`
		Claims []struct {
			ID                      string  `json:"id"`
			Status                  string  `json:"status"`
			RetirementRecord        string  `json:"retirement_record"`
			Fallback                string  `json:"fallback"`
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
				ReviewerRole    string `json:"reviewer_role"`
				Date            string `json:"date"`
				Verdict         string `json:"verdict"`
				AppliesTo       string `json:"applies_to"`
				CandidateCommit string `json:"candidate_commit"`
				Candidate       struct {
					Protocol  string `json:"protocol_json"`
					Decision7 string `json:"decision_0007"`
					Decision4 string `json:"decision_0004"`
					Analysis  string `json:"authoring_analysis"`
					Summary   string `json:"authoring_summary"`
				} `json:"candidate_lf_normalized_utf8_sha256"`
				AcceptedLimitations        []string `json:"accepted_limitations"`
				UnresolvedAuthoringIssues  []string `json:"unresolved_authoring_issues"`
				RemainingExecutionBlockers []string `json:"remaining_execution_blockers"`
			} `json:"review_record"`
			HistoricalV3 struct {
				AppliesTo string `json:"applies_to"`
			} `json:"historical_v3_review_record"`
		} `json:"review"`
		Remaining []string `json:"remaining_execution_blockers"`
		Gate      struct {
			MayRunAuthoring   bool   `json:"may_run_authoring"`
			MayOpenValidation bool   `json:"may_open_validation"`
			MayOpenHeldOut    bool   `json:"may_open_held_out"`
			Reason            string `json:"reason"`
		} `json:"gate"`
	}
	if err := json.Unmarshal(data, &protocol); err != nil {
		t.Fatal(err)
	}
	if protocol.Version != 5 || protocol.Status != "frozen" || !protocol.Frozen || !protocol.Review.Accepted {
		t.Fatalf("unexpected protocol state: v=%d status=%s frozen=%t", protocol.Version, protocol.Status, protocol.Frozen)
	}
	if !strings.Contains(protocol.Amendment.Activation, "Before the session's first executable act") || !strings.Contains(protocol.Amendment.Activation, "must not reactivate from intent") || !strings.Contains(protocol.Amendment.Activation, "current pending frontier") {
		t.Fatalf("bootstrap activation or current-frontier state scope regressed: %q", protocol.Amendment.Activation)
	}
	if !strings.Contains(protocol.Amendment.Activation, "status exhausted") || !strings.Contains(protocol.Amendment.Activation, "always_ready") || !strings.Contains(protocol.Amendment.Activation, "no ready call matches this intent") {
		t.Fatalf("v5 bootstrap-only orientation contract is incomplete: %q", protocol.Amendment.Activation)
	}
	if !strings.Contains(protocol.Amendment.Discovery, "reachable_verbs") || !strings.Contains(protocol.Amendment.Discovery, "capped at 12") || !strings.Contains(protocol.Amendment.Discovery, "declared_args") || !strings.Contains(protocol.Amendment.Discovery, "2048") {
		t.Fatalf("v5 discovery-bearing error contract is incomplete: %q", protocol.Amendment.Discovery)
	}
	if !strings.Contains(protocol.Amendment.Envelope, "expand(compress(x)) == x") || !strings.Contains(protocol.Amendment.Envelope, "wb") || !strings.Contains(protocol.Amendment.Envelope, "8 random characters") {
		t.Fatalf("v5 envelope compression contract is incomplete: %q", protocol.Amendment.Envelope)
	}
	if !strings.Contains(protocol.Amendment.Grading, "selector-addressed, never sequence-index-addressed") || !strings.Contains(protocol.Amendment.Grading, "recency_rationale") || !strings.Contains(protocol.Amendment.Grading, "final_message_states is barred") {
		t.Fatalf("v5 grading-altitude contract is incomplete: %q", protocol.Amendment.Grading)
	}
	if !strings.Contains(protocol.Amendment.StateChanges, "shared trial-harness rule for every arm") || !strings.Contains(protocol.Amendment.StateChanges, "flat-tool Arms A and B") || !strings.Contains(protocol.Amendment.StateChanges, "Orientations do not advance") {
		t.Fatalf("cross-arm state-event contract regressed: %q", protocol.Amendment.StateChanges)
	}
	if !strings.Contains(protocol.Amendment.Retirement, "All v4-sealed unopened candidates are retired") || !strings.Contains(protocol.Amendment.Retirement, "diagnostic-only") {
		t.Fatalf("v5 retirement clause is incomplete: %q", protocol.Amendment.Retirement)
	}
	if protocol.Amendment.SupersededV4.AppliesTo != "superseded protocol v4 only" {
		t.Fatalf("superseded v4 amendment is not scoped: %q", protocol.Amendment.SupersededV4.AppliesTo)
	}
	var armKeys struct {
		Arms          map[string]json.RawMessage `json:"arms"`
		SystemPrompts map[string]json.RawMessage `json:"system_prompts"`
	}
	if err := json.Unmarshal(data, &armKeys); err != nil {
		t.Fatal(err)
	}
	for _, retired := range []string{"E"} {
		if _, exists := armKeys.Arms[retired]; exists {
			t.Fatalf("retired arm %s is still registered", retired)
		}
		if _, exists := armKeys.SystemPrompts[retired]; exists {
			t.Fatalf("retired arm %s still has a system prompt", retired)
		}
	}
	for _, required := range []string{"A", "B", "C", "D", "D_prime"} {
		if _, exists := armKeys.Arms[required]; !exists {
			t.Fatalf("arm %s is missing from the registration", required)
		}
	}
	basePrompt := "Use only the configured tools. Follow the user request exactly."
	intentInstruction := "Before the first executable act, send the complete goal as intent"
	conditionalInstruction := "If no call is returned, proceed directly, or answer without acting when the tools cannot serve the goal."
	for _, arm := range []string{"A", "B", "C", "D"} {
		prompt, exists := protocol.SystemPrompts[arm]
		if !exists || !strings.HasPrefix(prompt, basePrompt) {
			t.Fatalf("arm %s system prompt is not pinned: %q", arm, prompt)
		}
		hasIntent := strings.Contains(prompt, intentInstruction)
		wantIntent := arm == "C" || arm == "D"
		if hasIntent != wantIntent {
			t.Fatalf("arm %s intent instruction = %t, want %t: %q", arm, hasIntent, wantIntent, prompt)
		}
		if wantIntent && !strings.Contains(prompt, conditionalInstruction) {
			t.Fatalf("arm %s bootstrap prompt is not conditional: %q", arm, prompt)
		}
	}
	for arm, stateEvents := range map[string]bool{
		"A": protocol.Arms.A.StateEvents, "B": protocol.Arms.B.StateEvents,
		"C": protocol.Arms.C.StateEvents, "D": protocol.Arms.D.StateEvents,
	} {
		if !stateEvents {
			t.Fatalf("arm %s does not share harness state events", arm)
		}
	}
	if protocol.Arms.A.Orientation != false || protocol.Arms.B.Orientation != false {
		t.Fatalf("flat arms unexpectedly expose orientation: A=%v B=%v", protocol.Arms.A.Orientation, protocol.Arms.B.Orientation)
	}
	for arm, orientation := range map[string]any{"C": protocol.Arms.C.Orientation, "D": protocol.Arms.D.Orientation} {
		value, ok := orientation.(string)
		if !ok || !strings.Contains(value, "bootstrap_only") || !strings.Contains(value, "exhausted") {
			t.Fatalf("arm %s orientation is not v5 bootstrap-only: %v", arm, orientation)
		}
	}
	record := protocol.Review.Record
	if record.ReviewerRole != "evaluation reviewer, independent of implementation" || record.Date != "2026-08-18" || record.Verdict != "ACCEPT" {
		t.Fatalf("unexpected independent acceptance record: %+v", record)
	}
	if record.CandidateCommit != "f889f13f0c514fa5108a1e392701ebeadc4376f7" || !strings.Contains(record.AppliesTo, record.CandidateCommit) {
		t.Fatalf("acceptance record identifies the wrong commit: %+v", record)
	}
	if record.Candidate.Protocol != "sha256:4601e1970ebd271161fa3c5c3c245da28a5f55252eac8823f7f9d4f862adf0ff" ||
		record.Candidate.Decision7 != "sha256:1df4e012ec5d57e1759681cc92d77049897d1517d1ba3f109aa8222791be7674" ||
		record.Candidate.Decision4 != "sha256:cb35b48709321f885cfc5cd50a6bc08a51bf570d6146b5acad1bb283d3e41903" ||
		record.Candidate.Analysis != "sha256:f28f07a7cf08583b45bd5e2dc043399c877288a9b7a68d95747ce9ddb15eb58b" ||
		record.Candidate.Summary != "sha256:114fab6f1a1d2b6f8c98c3b4a8ef544e9aa22cf5f7e19b53c674409568435405" {
		t.Fatalf("acceptance record identifies the wrong candidate: %+v", record.Candidate)
	}
	if protocol.Review.HistoricalV3.AppliesTo != "superseded protocol v3 only" {
		t.Fatalf("historical v3 acceptance is not scoped: %+v", protocol.Review.HistoricalV3)
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
		if claim.ID == "teaching_refusal_value" {
			if claim.Status != "retired_engineering_grade" ||
				!strings.Contains(claim.RetirementRecord, "No future artifact may cite teaching-refusal value as confirmed") ||
				!strings.Contains(claim.RetirementRecord, "permanent evidence record") ||
				!strings.Contains(claim.Fallback, "formal written requirement") ||
				!strings.Contains(claim.Multiplicity, "cannot create a headline pass") {
				t.Fatalf("teaching-refusal retirement record is incomplete: %+v", claim)
			}
			if claim.Uncertainty != "" || claim.Threshold != "" || claim.MinimumDetectableEffect != 0 {
				t.Fatalf("retired claim still carries live inference machinery: %+v", claim)
			}
			claims[claim.ID] = struct {
				Population              string
				Endpoint                string
				Uncertainty             string
				Threshold               string
				MinimumDetectableEffect float64
				PowerNote               string
				Multiplicity            string
			}{Multiplicity: claim.Multiplicity}
			continue
		}
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
	for _, isolation := range []string{"direct_no_tax", "frontier_value", "teaching_refusal_value"} {
		if !strings.Contains(claims[isolation].Multiplicity, "cannot create a headline pass") {
			t.Fatalf("isolation %s may not create a headline pass: %q", isolation, claims[isolation].Multiplicity)
		}
	}
	if !strings.Contains(protocol.Analysis.Multiplicity, "direct_no_tax and frontier_value are the live pre-registered isolations") || strings.Contains(protocol.Analysis.Multiplicity, "teaching_refusal_value") {
		t.Fatalf("multiplicity must name only the live isolations: %q", protocol.Analysis.Multiplicity)
	}
	if protocol.Tranche.ArmsPerPairingKey != 4 {
		t.Fatalf("arms per pairing key = %d, want 4", protocol.Tranche.ArmsPerPairingKey)
	}
	validationMaximum := float64(protocol.Tranche.MinimumCases*protocol.Trial.Repetitions*protocol.Tranche.ArmsPerPairingKey) * protocol.Trial.CostCap * 1.10
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
	if len(protocol.Remaining) != 7 || !strings.Contains(protocol.Remaining[0], "generation and sealing") || !strings.Contains(protocol.Remaining[6], "Task 0.6") {
		t.Fatalf("frozen protocol must retain execution blockers: %v", protocol.Remaining)
	}
	blockersText := strings.Join(protocol.Remaining, "\n")
	for _, required := range []string{"A-D runner", "path witness", "process-event log", "authoring dry run", "classified message set"} {
		if !strings.Contains(blockersText, required) {
			t.Fatalf("execution blockers are missing %q: %v", required, protocol.Remaining)
		}
	}
	wantLimitations := []string{
		"Conclusions are limited to Claude Code 2.1.229, claude-sonnet-5, the frozen dev-repo world, and the eight task classes.",
		"Expected headline power is about 62% at a 15-point effect; the approximate 80% MDE against zero is 0.19.",
		"The USD 300 validation ceiling cannot support 80% power for a 15-point Holm-adjusted LCB-floor design.",
		"Efficiency-only is a conservative conjunction; at a true zero difference, power to clear the 0.05 success or dead-end bounds is about 0.16.",
		"A/B/D/E adapters, including shared state-event application for flat tools, remain preimplementation harness work.",
		"D/E suppression must remove frontier and refusal calls before pending state is recorded.",
	}
	wantExecutionBlockers := []string{
		"independent generation and sealing of new validation and held_out families",
		"every artifact_freeze.before_validation digest",
		"Arm B document and human-factors review",
		"Phase 1 A-E runner adapters, including shared state-event application",
		"Task 0.6 Phase 0 review",
	}
	if len(record.UnresolvedAuthoringIssues) != 1 || !strings.Contains(record.UnresolvedAuthoringIssues[0], "Cascade remains run-to-run noisy") {
		t.Fatalf("accepted authoring limitation is missing: %v", record.UnresolvedAuthoringIssues)
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
	if !protocol.Gate.MayRunAuthoring || protocol.Gate.MayOpenValidation || protocol.Gate.MayOpenHeldOut || !strings.Contains(protocol.Gate.Reason, "Protocol v5") || !strings.Contains(protocol.Gate.Reason, "independently accepted") {
		t.Fatalf("v5 amendment must keep only authoring open: %+v", protocol.Gate)
	}
	if !strings.Contains(protocol.Gate.Reason, "review candidate") || !strings.Contains(protocol.Gate.Reason, "authoring dry run") {
		t.Fatalf("v5 gate reason must record the pending payload review and dry run: %q", protocol.Gate.Reason)
	}
}
