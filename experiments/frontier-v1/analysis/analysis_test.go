package main

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestArmReportsUseAssignedITTDenominator(t *testing.T) {
	observations := []observation{
		{Arm: "C", ITTSuccess: true},
		{Arm: "C", ITTSuccess: false, Unresolved: true},
		{Arm: "C", ITTSuccess: false, CompleteCase: true},
	}
	report := armReports(observations)["C"]
	if report.Assigned != 3 || report.Successes != 1 || report.SuccessRate != 1.0/3.0 {
		t.Fatalf("ITT arm report = %#v", report)
	}
	if report.Unresolved != 1 || report.UnresolvedRate != 1.0/3.0 {
		t.Fatalf("unresolved accounting = %#v", report)
	}
}

func TestIndeterminateRulesUseUnroundedAssignedRates(t *testing.T) {
	arms := map[string]armReport{}
	for _, arm := range phase1Arms {
		arms[arm] = armReport{Assigned: 100, Unresolved: 4, UnresolvedRate: 0.04}
	}
	c := arms["C"]
	c.Unresolved, c.UnresolvedRate = 6, 0.06
	arms["C"] = c
	b := arms["B"]
	b.Unresolved, b.UnresolvedRate = 3, 0.03
	arms["B"] = b
	reasons := indeterminateReasons(scheduledSummary{Status: "complete"}, arms, 0)
	if len(reasons) != 2 {
		t.Fatalf("reasons = %#v, want arm threshold and C-B imbalance", reasons)
	}
}

func TestExactSignFlipRejectsConsistentHarm(t *testing.T) {
	values := make([]pairedValue, 8)
	for index := range values {
		values[index] = pairedValue{FamilyID: integerString(index), CaseID: integerString(index), X: 0, Y: 1, Complete: true}
	}
	result := inferDifference(values)
	if result.Method != "exact family-mean sign-flip permutation" || result.Point != -1 {
		t.Fatalf("inference = %#v", result)
	}
	if result.PWorse == nil || *result.PWorse >= 0.05 {
		t.Fatalf("p_worse = %#v, want < 0.05", result.PWorse)
	}
	if result.LowerBound != -1.5 || result.UpperBound != -0.5 {
		t.Fatalf("bounds = [%f,%f], want [-1.5,-0.5]", result.LowerBound, result.UpperBound)
	}
}

func TestHierarchicalBootstrapIsDeterministicAtTwentyFamilies(t *testing.T) {
	values := make([]pairedValue, 20)
	for index := range values {
		values[index] = pairedValue{FamilyID: integerString(index), CaseID: integerString(index), X: 1, Y: 0, Complete: true}
	}
	first := inferDifference(values)
	second := inferDifference(values)
	if first.Method != "10000-replicate paired hierarchical bootstrap" || first != second {
		t.Fatalf("bootstrap inference is not pinned: first=%#v second=%#v", first, second)
	}
	if first.LowerBound != 1 || first.UpperBound != 1 {
		t.Fatalf("constant-effect bounds = [%f,%f]", first.LowerBound, first.UpperBound)
	}
	if first.PWorse != nil {
		t.Fatalf("bootstrap tail probability was mislabeled as a harm p-value: %f", *first.PWorse)
	}
}

func TestPairedSuccessRatioRequiresTenSuccessesAndHonorsCapImbalance(t *testing.T) {
	values := []pairedValue{{FamilyID: "family", CaseID: "case", X: 10, Y: 20}}
	lowSuccess := pairedSuccessRatio(values, 9, 20, false)
	if lowSuccess.Determinate || lowSuccess.Reason == "" {
		t.Fatalf("low-success ratio = %#v", lowSuccess)
	}
	imbalanced := pairedSuccessRatio(values, 20, 20, true)
	if imbalanced.Determinate || imbalanced.Reason == "" {
		t.Fatalf("cap-imbalanced ratio = %#v", imbalanced)
	}
}

func TestHeadlineCapabilityPassIsNotVetoedByCapImbalance(t *testing.T) {
	var observations []observation
	for family := 0; family < 24; family++ {
		caseID := "case-" + integerString(family)
		familyID := "family-" + integerString(family)
		observations = append(observations,
			observation{FamilyID: familyID, CaseID: caseID, Arm: "C", ITTSuccess: true, CompleteCase: true, CapHit: family == 0, TimeoutOrCap: family == 0, Tokens: 10},
			observation{FamilyID: familyID, CaseID: caseID, Arm: "B", ITTSuccess: false, CompleteCase: true, Tokens: 20},
		)
	}
	arms := armReports(observations)
	claim := headlineClaim(observations, arms, false)
	if claim.Decision != "capability_pass" {
		t.Fatalf("headline decision = %s: %s", claim.Decision, claim.Reason)
	}
	if claim.CostRatio == nil || claim.CostRatio.Determinate {
		t.Fatalf("cost ratio should remain separately indeterminate: %#v", claim.CostRatio)
	}
}

func TestHeadlineClearFailureIsNotMadeIndeterminateByCapImbalance(t *testing.T) {
	var observations []observation
	for family := 0; family < 24; family++ {
		caseID := "case-" + integerString(family)
		familyID := "family-" + integerString(family)
		observations = append(observations,
			observation{FamilyID: familyID, CaseID: caseID, Arm: "C", CompleteCase: true, CapHit: family == 0, TimeoutOrCap: family == 0, Tokens: 10},
			observation{FamilyID: familyID, CaseID: caseID, Arm: "B", ITTSuccess: true, CompleteCase: true, Tokens: 20},
		)
	}
	claim := headlineClaim(observations, armReports(observations), false)
	if claim.LowerBound == nil || *claim.LowerBound > -0.05 {
		t.Fatalf("test setup did not produce an independent capability failure: %#v", claim.LowerBound)
	}
	if claim.Decision != "fail" {
		t.Fatalf("clear failure became %s under cap imbalance: %s", claim.Decision, claim.Reason)
	}
}

func TestComponentContrastsApplyForcedDecisions(t *testing.T) {
	var observations []observation
	for family := 0; family < 8; family++ {
		caseID := "case-" + integerString(family)
		familyID := "family-" + integerString(family)
		observations = append(observations,
			observation{FamilyID: familyID, CaseID: caseID, Class: "recovery", Arm: "C", CompleteCase: true},
			observation{FamilyID: familyID, CaseID: caseID, Class: "recovery", Arm: "D", ITTSuccess: true, CompleteCase: true, Recovery: false},
		)
	}
	arms := armReports(observations)
	frontier := frontierClaim(observations, arms, false)
	if frontier.Decision != "remove_frontier" {
		t.Fatalf("frontier decision = %#v", frontier)
	}
}

func TestSensitivityKeepsITTAndCompleteCaseSeparate(t *testing.T) {
	values := []pairedValue{
		{FamilyID: "f1", CaseID: "c1", X: 1, Y: 0, Complete: true},
		{FamilyID: "f2", CaseID: "c2", X: 0, Y: 1, Complete: false},
	}
	report := sensitivity(values)
	if report.CompleteCase == nil || *report.CompleteCase != 1 {
		t.Fatalf("complete-case sensitivity = %#v", report.CompleteCase)
	}
	if math.Abs(report.OneVotePerCase) > 1e-12 || math.Abs(report.OneVotePerFamily) > 1e-12 {
		t.Fatalf("ITT sensitivities = %#v", report)
	}
}

func TestSafetyStopDiscardsWholeFourArmPairing(t *testing.T) {
	results := make([]assignedTrialResult, 0, phase1Repetitions*len(phase1Arms))
	for repetition := 0; repetition < phase1Repetitions; repetition++ {
		for _, arm := range phase1Arms {
			result := assignedTrialResult{FamilyID: "family", CaseID: "case", Repetition: repetition, Arm: arm, Status: "pass", ITTSuccess: true}
			if repetition == 1 && arm == "D" {
				result.Status = "safety_stopped"
				result.ITTSuccess = false
			}
			results = append(results, result)
		}
	}
	summary := scheduledSummary{V: 1, Tranche: "authoring", Assigned: len(results), Results: results}
	manifest := trancheManifest{V: 1, Tranche: "authoring", Cases: []manifestCase{{CaseID: "case", Class: "direct", FamilyID: "family"}}}
	observations, incomplete, err := loadObservations(t.TempDir(), summary, manifest)
	if err != nil {
		t.Fatal(err)
	}
	if incomplete != 1 {
		t.Fatalf("incomplete pairings = %d", incomplete)
	}
	for _, item := range observations {
		if item.Repetition == 1 && (item.ITTSuccess || item.CompleteCase || !item.Unresolved) {
			t.Fatalf("pair member was not discarded under ITT: %#v", item)
		}
	}
}

func TestFrontierLinkageAllowsOmittedStateAndCanonicalNumericArgs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "runtime.jsonl")
	events := []any{
		map[string]any{"type": "user", "message": map[string]any{"content": []any{map[string]any{
			"type": "tool_result", "content": `{"frontier":[{"call":{"handle":"h_1","verb":"focus","args":{"count":1},"state":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}}]}`,
		}}}},
		map[string]any{"type": "assistant", "message": map[string]any{"content": []any{map[string]any{
			"type": "tool_use", "name": "mcp__phoenix__act", "input": map[string]any{"handle": "h_1", "verb": "focus", "args": map[string]any{"count": 1}},
		}}}},
	}
	var contents []byte
	for _, event := range events {
		line, err := json.Marshal(event)
		if err != nil {
			t.Fatal(err)
		}
		contents = append(contents, line...)
		contents = append(contents, '\n')
	}
	if err := os.WriteFile(path, contents, 0o600); err != nil {
		t.Fatal(err)
	}
	shown, taken, passed, err := frontierEvidence(path, true)
	if err != nil {
		t.Fatal(err)
	}
	if shown != 1 || taken != 1 || passed != 1 {
		t.Fatalf("frontier evidence = shown %d, taken %d, passed %d", shown, taken, passed)
	}
}

func TestDeadEndRecognizesFailedHelpRequestButNotPassingAbsenceReport(t *testing.T) {
	if !asksForHelp("I cannot proceed without the deployment target. Could you provide it?") {
		t.Fatal("expected explicit request for missing input to be a help request")
	}
	if !asksForHelp("Please confirm the exact test name.") {
		t.Fatal("documented please-confirm request was not classified as help seeking")
	}
	failed := assignedTrialResult{Termination: "graded_fail"}
	if !isDeadEnd(failed, observation{Acts: 1, HelpRequest: true}) {
		t.Fatal("failed help request should be a dead end")
	}
	if isDeadEnd(assignedTrialResult{Termination: "graded_pass"}, observation{ITTSuccess: true, Acts: 0, HelpRequest: true}) {
		t.Fatal("passing absence report should not be a dead end")
	}
}

func TestReportTemplateRendersAndRefusesOverwrite(t *testing.T) {
	templateContents, err := os.ReadFile("report-template.md.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "template.md.tmpl"), templateContents, 0o600); err != nil {
		t.Fatal(err)
	}
	report := analysisReport{
		V: 1, Tranche: "authoring", Status: "analysis_complete", Arms: map[string]armReport{"A": {}},
		Claims: []claimReport{
			{ID: "headline_value", Intervention: "C", Comparator: "B", Decision: "fail", Reason: "synthetic test"},
			{ID: "direct_no_tax", Intervention: "C", Comparator: "A", PWorse: floatPointer(0.25), Decision: "no_harm_signal", Reason: "synthetic test"},
		},
	}
	if err := renderNewReport(root, "template.md.tmpl", "report.md", report); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(filepath.Join(root, "report.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(contents), "headline_value") || !strings.Contains(string(contents), "Gate 1A") || !strings.Contains(string(contents), "versus -0.100000") {
		t.Fatalf("rendered report omitted required sections:\n%s", contents)
	}
	if strings.Count(string(contents), "sign-flip harm p-value") != 1 {
		t.Fatalf("report should label only the sign-flip p-value:\n%s", contents)
	}
	if err := renderNewReport(root, "template.md.tmpl", "report.md", report); err == nil {
		t.Fatal("report renderer overwrote an existing output")
	}
}

func floatPointer(value float64) *float64 { return &value }
