package spike

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestRecordedTask02EvidenceIsAcceptedWithPhase0Closed(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "results.json"))
	if err != nil {
		t.Fatal(err)
	}
	var results struct {
		Status             string `json:"status"`
		SchemaMeasurements []struct {
			Surface     string `json:"surface"`
			WorldSize   int    `json:"world_size"`
			ToolCount   int    `json:"tool_count"`
			SchemaBytes int    `json:"schema_bytes"`
		} `json:"schema_measurements"`
		StandingTokens struct {
			Repetitions int `json:"repetitions"`
			InputTokens struct {
				Empty []int `json:"empty"`
				Act   []int `json:"act"`
				Eval  []int `json:"eval"`
			} `json:"input_tokens"`
			Incremental struct {
				Act  int `json:"act"`
				Eval int `json:"eval"`
			} `json:"incremental_tokens_vs_empty"`
			ProbeCostUSD      *float64 `json:"probe_cost_usd"`
			CostRetentionNote string   `json:"cost_retention_note"`
			Budget            int      `json:"standing_budget_tokens"`
			Passed            bool     `json:"passed"`
		} `json:"standing_token_measurements"`
		MalformedCalls []struct {
			Surface           string   `json:"surface"`
			FreshSessions     int      `json:"fresh_sessions"`
			CompletedSessions int      `json:"completed_sessions"`
			MalformedCalls    int      `json:"malformed_calls"`
			Rate              float64  `json:"malformed_call_rate"`
			MaximumRate       float64  `json:"maximum_rate"`
			SessionIDs        []string `json:"session_ids"`
			Wilson95          struct {
				Lower float64 `json:"lower"`
				Upper float64 `json:"upper"`
			} `json:"wilson_95_interval"`
			Complete bool `json:"complete_for_gate"`
			Passed   bool `json:"passed"`
		} `json:"malformed_call_measurements"`
		Linux struct {
			Status  string `json:"status"`
			Latency struct {
				Values      []float64 `json:"values"`
				Median      float64   `json:"median"`
				P95         float64   `json:"p95_nearest_rank"`
				MedianLimit float64   `json:"median_limit"`
				P95Limit    float64   `json:"p95_limit"`
				Passed      bool      `json:"passed"`
			} `json:"daemon_only_act_latency"`
			Share struct {
				MedianTrialWallMS      float64   `json:"median_successful_trial_wall_time_ms"`
				SourceTrialDurationsMS []float64 `json:"source_trial_durations_ms"`
				DaemonMedianMS         float64   `json:"daemon_median_ms"`
				Ratio                  float64   `json:"ratio"`
				Percent                float64   `json:"percent"`
				MaximumRatio           float64   `json:"maximum_ratio"`
				Passed                 bool      `json:"passed"`
			} `json:"daemon_share_of_trial_wall_time"`
		} `json:"linux_rerun"`
		Review struct {
			Accepted bool `json:"accepted"`
			Record   struct {
				ReviewerRole string `json:"reviewer_role"`
				Date         string `json:"date"`
				Verdict      string `json:"verdict"`
				Candidate    struct {
					Results  string `json:"results_json"`
					Decision string `json:"decision_0002"`
					Readme   string `json:"surface_spike_readme"`
					Test     string `json:"results_test_go"`
				} `json:"candidate_lf_normalized_utf8_sha256"`
				AcceptedLimitations []string `json:"accepted_limitations"`
				RemainingBlockers   []string `json:"remaining_blockers"`
			} `json:"review_record"`
		} `json:"review"`
		Gate struct {
			Accepted         bool     `json:"accepted"`
			AcceptedOn       string   `json:"accepted_on"`
			Decision         string   `json:"decision"`
			EvidenceComplete bool     `json:"evidence_complete"`
			Missing          []string `json:"missing"`
		} `json:"phase_0_gate"`
	}
	if err := json.Unmarshal(data, &results); err != nil {
		t.Fatal(err)
	}
	if results.Status != "accepted" || !results.Review.Accepted || !results.Gate.Accepted || results.Gate.AcceptedOn != "2026-08-18" || results.Gate.Decision != "docs/decisions/0005-phase-0-review.md" || !results.Gate.EvidenceComplete {
		t.Fatalf("unexpected Task 0.2 acceptance state: status=%s review=%+v gate=%+v", results.Status, results.Review, results.Gate)
	}
	if len(results.Gate.Missing) != 0 {
		t.Fatalf("unexpected remaining Phase 0 blockers: %v", results.Gate.Missing)
	}
	wantReviewTimeBlockers := []string{"independent unopened validation and held_out families", "Task 0.6 Phase 0 review"}
	record := results.Review.Record
	if record.ReviewerRole != "evaluation reviewer, independent of implementation" || record.Date != "2026-08-18" || record.Verdict != "ACCEPT" || len(record.AcceptedLimitations) != 5 || !slices.Equal(record.RemainingBlockers, wantReviewTimeBlockers) {
		t.Fatalf("unexpected Task 0.2 review record: %+v", record)
	}
	if record.Candidate.Results != "sha256:376a96534eac11f08065ffc300bfb8b73fcec0c8c18a37336edf0c3731dbaf95" || record.Candidate.Decision != "sha256:897c84d5891705997efeabb8f135fb127f75f48ed31e7fa63198b4d5449d973e" || record.Candidate.Readme != "sha256:7df2814880a1d185d2a844e51302c34b41596ca9c9a5b51cecf19c63cc8706bd" || record.Candidate.Test != "sha256:2342f8266068e5829c8302dd8eda8284215c4c93abba44ff7a11da3d3585694c" {
		t.Fatalf("acceptance record identifies the wrong candidate: %+v", record.Candidate)
	}
	wantSchemas := map[string][2]int{
		"act/10": {1, 1621}, "act/100": {1, 1621}, "act/1000": {1, 1621},
		"eval/10": {1, 1490}, "eval/100": {1, 1490}, "eval/1000": {1, 1490},
		"flat/10": {10, 3861}, "flat/100": {100, 38601}, "flat/1000": {1000, 386001},
	}
	if len(results.SchemaMeasurements) != len(wantSchemas) {
		t.Fatalf("got %d schema measurements, want %d", len(results.SchemaMeasurements), len(wantSchemas))
	}
	for _, measurement := range results.SchemaMeasurements {
		key := measurement.Surface + "/" + fmt.Sprint(measurement.WorldSize)
		want, ok := wantSchemas[key]
		if !ok || measurement.ToolCount != want[0] || measurement.SchemaBytes != want[1] {
			t.Fatalf("unexpected schema measurement %s: %+v", key, measurement)
		}
	}
	standing := results.StandingTokens
	if standing.Repetitions != 3 || !slices.Equal(standing.InputTokens.Empty, []int{227, 227, 227}) || !slices.Equal(standing.InputTokens.Act, []int{774, 774, 774}) || !slices.Equal(standing.InputTokens.Eval, []int{712, 712, 712}) {
		t.Fatalf("unexpected standing-token samples: %+v", standing)
	}
	actIncrement := standing.InputTokens.Act[0] - standing.InputTokens.Empty[0]
	evalIncrement := standing.InputTokens.Eval[0] - standing.InputTokens.Empty[0]
	if actIncrement != standing.Incremental.Act || evalIncrement != standing.Incremental.Eval || actIncrement != 547 || evalIncrement != 485 || standing.Budget != 600 || !standing.Passed {
		t.Fatalf("standing-token arithmetic regressed: act=%d eval=%d record=%+v", actIncrement, evalIncrement, standing)
	}
	if standing.ProbeCostUSD != nil || !strings.Contains(standing.CostRetentionNote, "not retained") {
		t.Fatalf("standing-token cost retention must be explicit: cost=%v note=%q", standing.ProbeCostUSD, standing.CostRetentionNote)
	}
	if len(results.MalformedCalls) != 2 {
		t.Fatalf("got %d malformed-call samples, want 2", len(results.MalformedCalls))
	}
	seenSurfaces := map[string]bool{}
	seenSessions := map[string]bool{}
	const z = 1.959963984540054
	for _, sample := range results.MalformedCalls {
		if sample.Surface != "act" && sample.Surface != "eval" {
			t.Fatalf("unexpected surface %q", sample.Surface)
		}
		if seenSurfaces[sample.Surface] {
			t.Fatalf("duplicate surface %q", sample.Surface)
		}
		seenSurfaces[sample.Surface] = true
		if sample.FreshSessions != 20 || sample.CompletedSessions != 20 || len(sample.SessionIDs) != 20 || sample.MalformedCalls != 0 || sample.Rate != 0 || sample.MaximumRate != 0.05 || !sample.Complete || !sample.Passed {
			t.Fatalf("incomplete malformed-call evidence for %s: %+v", sample.Surface, sample)
		}
		for _, sessionID := range sample.SessionIDs {
			if sessionID == "" || seenSessions[sessionID] {
				t.Fatalf("missing or duplicate final session ID %q", sessionID)
			}
			seenSessions[sessionID] = true
		}
		wantUpper := z * z / (float64(sample.FreshSessions) + z*z)
		if sample.Wilson95.Lower != 0 || math.Abs(sample.Wilson95.Upper-wantUpper) > 0.0000005 {
			t.Fatalf("incorrect Wilson interval for %s: %+v, want upper %.6f", sample.Surface, sample.Wilson95, wantUpper)
		}
	}
	if len(seenSessions) != 40 {
		t.Fatalf("got %d unique final sessions, want 40", len(seenSessions))
	}
	latency := results.Linux.Latency
	if results.Linux.Status != "success" || len(latency.Values) != 100 || !latency.Passed || !results.Linux.Share.Passed {
		t.Fatalf("Linux evidence is incomplete: %+v", results.Linux)
	}
	slices.Sort(latency.Values)
	median := (latency.Values[49] + latency.Values[50]) / 2
	p95 := latency.Values[94]
	if median != latency.Median || p95 != latency.P95 || median > latency.MedianLimit || p95 > latency.P95Limit {
		t.Fatalf("Linux latency mismatch: median=%.1f/%.1f p95=%.1f/%.1f", median, latency.Median, p95, latency.P95)
	}
	share := results.Linux.Share
	if len(share.SourceTrialDurationsMS) != 2 {
		t.Fatalf("daemon share needs two source trial durations: %+v", share)
	}
	slices.Sort(share.SourceTrialDurationsMS)
	trialMedian := (share.SourceTrialDurationsMS[0] + share.SourceTrialDurationsMS[1]) / 2
	ratio := share.DaemonMedianMS / trialMedian
	if trialMedian != share.MedianTrialWallMS || math.Abs(ratio-share.Ratio) > 0.000000005 || math.Abs(ratio*100-share.Percent) > 0.000005 || share.Ratio > share.MaximumRatio {
		t.Fatalf("daemon share arithmetic regressed: computed=%.8f record=%+v", ratio, share)
	}
}
