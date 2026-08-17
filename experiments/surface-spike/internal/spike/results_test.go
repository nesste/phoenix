package spike

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestRecordedTask02EvidenceIsCompletePendingReview(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "results.json"))
	if err != nil {
		t.Fatal(err)
	}
	var results struct {
		Status         string `json:"status"`
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
				Ratio        float64 `json:"ratio"`
				MaximumRatio float64 `json:"maximum_ratio"`
				Passed       bool    `json:"passed"`
			} `json:"daemon_share_of_trial_wall_time"`
		} `json:"linux_rerun"`
		Gate struct {
			Accepted         bool     `json:"accepted"`
			EvidenceComplete bool     `json:"evidence_complete"`
			Missing          []string `json:"missing"`
		} `json:"phase_0_gate"`
	}
	if err := json.Unmarshal(data, &results); err != nil {
		t.Fatal(err)
	}
	if results.Status != "review_candidate" || results.Gate.Accepted || !results.Gate.EvidenceComplete || len(results.Gate.Missing) != 1 || results.Gate.Missing[0] != "independent Task 0.2 evidence review" {
		t.Fatalf("unexpected Task 0.2 review state: status=%s gate=%+v", results.Status, results.Gate)
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
	if results.Linux.Share.Ratio > results.Linux.Share.MaximumRatio {
		t.Fatalf("daemon share %.8f exceeds limit %.2f", results.Linux.Share.Ratio, results.Linux.Share.MaximumRatio)
	}
}
