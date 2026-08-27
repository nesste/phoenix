package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// outcomeFreePartialFields is the exact protocol-v5 section 9 field set of the
// partial summary. The proposal admits launch index, completed pairing keys,
// per-arm assignment counts, artifact digests, timestamps, and cumulative
// spend, and nothing else; the identity fields naming the tranche and the
// schedule carry no outcome either.
var outcomeFreePartialFields = map[string]bool{
	"v": true, "tranche": true, "status": true, "schedule_digest": true,
	"artifact_digests": true, "started_at": true, "updated_at": true,
	"next_launch_index": true, "completed_pairing_keys": true,
	"arm_assignment_counts": true, "spent_usd": true, "resumes": true,
}

func TestScheduledRunWritesAnOutcomeFreePartialSummary(t *testing.T) {
	schedule, output := runResilienceSchedule(t)
	assertPartialSummaryFieldsAreOutcomeFree(t, filepath.Join(output, scheduledPartialSummaryName))
	summary, err := readScheduledPartialSummary(output)
	if err != nil {
		t.Fatal(err)
	}
	if summary.V != 1 || summary.Status != "in_progress" || summary.Resumes != 0 ||
		summary.NextLaunchIndex != len(schedule.Entries) || len(sortedPartialPairingKeys(summary)) != 2 {
		t.Fatalf("partial summary = %#v", summary)
	}
	for _, arm := range phase1Arms {
		if summary.ArmAssignmentCounts[arm] != 2 {
			t.Fatalf("arm %s assignment count = %d, want 2", arm, summary.ArmAssignmentCounts[arm])
		}
	}
	if summary.ArtifactDigests.WorldBuild != testBuild ||
		!strings.HasPrefix(summary.ArtifactDigests.ProcessEventLog, "sha256:") ||
		summary.StartedAt == "" || summary.UpdatedAt == "" {
		t.Fatalf("partial identity = %#v", summary)
	}
}

func TestScheduledRunAppendsProcessEventsForEveryLifecyclePoint(t *testing.T) {
	schedule, output := runResilienceSchedule(t)
	events := readProcessEvents(t, output)
	wantKinds := []string{
		processEventStart, processEventCheckpoint, processEventCheckpoint,
		processEventCheckpoint, processEventStop,
	}
	if len(events) != len(wantKinds) {
		t.Fatalf("process events = %#v", events)
	}
	for index, want := range wantKinds {
		event := events[index]
		if event.Event != want || event.V != 1 || event.RecordedAt == "" ||
			event.Resumes != 0 || event.ScheduleDigest != "sha256:schedule" {
			t.Fatalf("process event %d = %#v, want kind %q", index, event, want)
		}
	}
	if !events[0].Principal.Exposed || events[0].Principal.PID != os.Getpid() ||
		events[0].Principal.Source != "host_process_owner" {
		t.Fatalf("start principal = %#v", events[0].Principal)
	}
	last := events[len(events)-1]
	if last.Detail != "complete" || last.NextLaunchIndex != len(schedule.Entries) {
		t.Fatalf("stop event = %#v", last)
	}
}

func TestProcessTerminationRecordsAnUnexposedSendingPrincipal(t *testing.T) {
	output := t.TempDir()
	log := newProcessEventLog(output, "validation", "sha256:schedule", testBuild)
	if err := log.recordTermination(processEventProgress{nextLaunchIndex: 8, resumes: 1}, "terminated"); err != nil {
		t.Fatal(err)
	}
	events := readProcessEvents(t, output)
	if len(events) != 1 || events[0].Event != processEventTermination || events[0].Resumes != 1 {
		t.Fatalf("termination events = %#v", events)
	}
	principal := events[0].Principal
	if principal.Exposed || principal.Source != "signal_without_sender_identity" || principal.Reason == "" {
		t.Fatalf("termination principal = %#v", principal)
	}
}

func TestValidationResumeRequiresCustodianAttestation(t *testing.T) {
	control, resume := interruptedValidationResume(t, 0)
	control.attestationPath = ""
	err := requireResumeAuthorization(control, resume)
	if err == nil || !strings.Contains(err.Error(), "--resume-attestation") {
		t.Fatalf("unattested resume error = %v", err)
	}
}

func TestAuthoringResumeIsExemptFromTheAttestationRule(t *testing.T) {
	control, resume := interruptedValidationResume(t, 0)
	control.tranche = "authoring"
	control.attestationPath = ""
	if err := requireResumeAuthorization(control, resume); err != nil {
		t.Fatalf("authoring resume = %v", err)
	}
}

func TestValidationResumeRefusesASecondResume(t *testing.T) {
	control, resume := interruptedValidationResume(t, maximumScheduledResumes)
	err := requireResumeAuthorization(control, resume)
	if err == nil || !strings.Contains(err.Error(), "already used its single resume") {
		t.Fatalf("second-resume error = %v", err)
	}
}

func TestValidationResumeAcceptsOnlyFrozenOutcomeUncorrelatedCauses(t *testing.T) {
	for _, cause := range []string{"host_restart", "power_loss", "hardware_failure"} {
		t.Run(cause, func(t *testing.T) {
			control, resume := interruptedValidationResume(t, 0)
			writeTestAttestation(t, control, func(attestation map[string]any) {
				attestation["cause"] = cause
			})
			if err := requireResumeAuthorization(control, resume); err != nil {
				t.Fatalf("cause %s = %v", cause, err)
			}
		})
	}
	for _, cause := range []string{"oom", "out_of_memory", "operator_choice", ""} {
		t.Run("refused-"+cause, func(t *testing.T) {
			control, resume := interruptedValidationResume(t, 0)
			writeTestAttestation(t, control, func(attestation map[string]any) {
				attestation["cause"] = cause
			})
			err := requireResumeAuthorization(control, resume)
			if err == nil || !strings.Contains(err.Error(), "frozen outcome-uncorrelated cause list") {
				t.Fatalf("cause %q error = %v", cause, err)
			}
		})
	}
}

func TestValidationResumeRefusesAKillWithoutAnEstablishedPrincipal(t *testing.T) {
	control, resume := interruptedValidationResume(t, 0)
	writeTestAttestation(t, control, func(attestation map[string]any) {
		attestation["cause"] = "process_kill"
		attestation["kill_principal_established_as_non_participant"] = false
	})
	err := requireResumeAuthorization(control, resume)
	if err == nil || !strings.Contains(err.Error(), "non-participant principal") {
		t.Fatalf("unestablished kill error = %v", err)
	}
	writeTestAttestation(t, control, func(attestation map[string]any) {
		attestation["cause"] = "process_kill"
	})
	if err := requireResumeAuthorization(control, resume); err != nil {
		t.Fatalf("established kill = %v", err)
	}
}

func TestValidationResumeRefusesBrokenConditionsAndTiming(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(map[string]any)
		message string
	}{
		{"inspected outcomes", func(a map[string]any) {
			a["no_outcome_inspection_between_interruption_and_resume"] = false
		}, "no-inspection"},
		{"cause classified after the decision", func(a map[string]any) {
			a["cause_classified_before_inspection_and_resume_decision"] = false
		}, "no-inspection"},
		{"frozen bytes changed", func(a map[string]any) {
			a["frozen_bytes_digest_identical_at_resume"] = false
		}, "no-inspection"},
		{"different world build", func(a map[string]any) {
			a["frozen_digests_verified"] = map[string]any{"schedule_digest": "sha256:schedule", "world_build": "sha256:other"}
		}, "world_build"},
		{"past the backstop", func(a map[string]any) {
			a["interrupted_at"] = "2026-08-20T00:00:00Z"
			a["cause_classified_at"] = "2026-08-20T00:05:00Z"
			a["conditions_verified_at"] = "2026-08-20T00:10:00Z"
			a["resume_authorized_at"] = "2026-08-24T00:00:00Z"
		}, "72-hour backstop"},
		{"timestamps out of order", func(a map[string]any) {
			a["cause_classified_at"] = "2026-08-19T00:00:00Z"
		}, "out of order"},
		{"lapse without explanation", func(a map[string]any) {
			a["resume_authorized_at"] = "2026-08-20T06:00:00Z"
			delete(a, "lapse_explanation")
		}, "lapse explanation"},
		{"wrong resume index", func(a map[string]any) { a["resume_index"] = 2 }, "this is resume 1"},
		{"no custodian", func(a map[string]any) { a["custodian"] = "" }, "custodian"},
		{"stale process-event log digest", func(a map[string]any) {
			a["process_event_log_sha256"] = "sha256:" + strings.Repeat("0", 64)
		}, "does not match the attested"},
		{"foreign output directory", func(a map[string]any) { a["output_dir"] = "results/elsewhere" }, "output directory"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			control, resume := interruptedValidationResume(t, 0)
			writeTestAttestation(t, control, test.mutate)
			err := requireResumeAuthorization(control, resume)
			if err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("%s error = %v", test.name, err)
			}
		})
	}
}

func TestPostClosureGateRefusesValidationAfterANonCompletedExecution(t *testing.T) {
	root := openValidationGateTestRoot(t, filepath.Join("..", "..", ".."))
	gate := closureGateFixture(t, root)

	for _, status := range []string{"indeterminate", "budget_stopped", "lapsed"} {
		err := gate(status, "", "")
		if err == nil || !strings.Contains(err.Error(), "independently reviewed closure record") {
			t.Fatalf("post-closure gate error for %s = %v", status, err)
		}
	}
	if err := gate("indeterminate", "docs/reviews/closure.md", "REVISE"); err == nil ||
		!strings.Contains(err.Error(), "independently reviewed closure record") {
		t.Fatalf("unaccepted closure review = %v", err)
	}
	if err := gate("indeterminate", "docs/reviews/closure.md", "ACCEPT"); err == nil ||
		!strings.Contains(err.Error(), "closure review record") {
		t.Fatalf("missing closure review file = %v", err)
	}
	writeTestClosureReview(t, root)
	if err := gate("indeterminate", "docs/reviews/closure.md", "ACCEPT"); err != nil {
		t.Fatalf("reviewed closure record = %v", err)
	}
	if err := gate("complete", "", ""); err != nil {
		t.Fatalf("completed execution = %v", err)
	}
}

// closureGateFixture returns a setter that rewrites the post-closure fields of
// a test gate document and reports what the runner's gate check then says.
func closureGateFixture(t *testing.T, root string) func(status, review, verdict string) error {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(preValidationArtifactsPath))
	var document map[string]any
	readJSONForTest(t, path, &document)
	gates := document["gates"].(map[string]any)
	return func(status, review, verdict string) error {
		gates["validation_execution_status"] = status
		gates["validation_execution_closure_review"] = review
		gates["validation_execution_closure_review_verdict"] = verdict
		if err := writeJSON(path, document); err != nil {
			t.Fatal(err)
		}
		return requireValidationGate(root)
	}
}

func writeTestClosureReview(t *testing.T, root string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, "docs", "reviews"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "reviews", "closure.md"), []byte("closure\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

// assertPartialSummaryFieldsAreOutcomeFree checks the written JSON object
// against the admitted section 9 field set in both directions.
func assertPartialSummaryFieldsAreOutcomeFree(t *testing.T, path string) {
	t.Helper()
	var raw map[string]json.RawMessage
	readJSONForTest(t, path, &raw)
	for field := range raw {
		if !outcomeFreePartialFields[field] {
			t.Fatalf("partial summary carries the non-admitted field %q", field)
		}
	}
	for field := range outcomeFreePartialFields {
		if _, present := raw[field]; !present {
			t.Fatalf("partial summary is missing the required field %q", field)
		}
	}
}

// runResilienceSchedule executes a complete two-pairing-key authoring schedule
// and returns it with the output directory holding the section 9 records.
func runResilienceSchedule(t *testing.T) (launchSchedule, string) {
	t.Helper()
	schedule := distinctPairingSchedule(2)
	repository := t.TempDir()
	createRunnerFixture(t, repository)
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	config := runConfig{
		repositoryRoot: repository, outputDir: output, worldBuild: testBuild,
		worldPath: filepath.Join(repository, "world.json"), schemaPath: filepath.Join(repository, "schema.json"),
		timeout: time.Second, budgetUSD: "0.15",
	}
	driver := &scriptedRuntime{results: repeatedRuntimeResults(len(schedule.Entries))}
	if _, err := runScheduledCases(config, schedule, "sha256:schedule", 10, driver, fakeGrader{}); err != nil {
		t.Fatal(err)
	}
	return schedule, output
}

// distinctPairingSchedule gives each pairing key its own repetition so the
// partial summary's completed pairing-key list is exercised with more than one
// member; the shared helper reuses one (case_id, repetition) pair.
func distinctPairingSchedule(keys int) launchSchedule {
	schedule := testPairingSchedule(keys)
	for index := range schedule.Entries {
		schedule.Entries[index].Repetition = schedule.Entries[index].PairingIndex
	}
	return schedule
}

func repeatedRuntimeResults(count int) []runtimeResult {
	results := make([]runtimeResult, 0, count)
	for index := 0; index < count; index++ {
		results = append(results, runtimeResult{Metrics: runtimeMetrics{CostUSD: 0.01, TotalTokens: 4}})
	}
	return results
}

func readProcessEvents(t *testing.T, outputDir string) []processEvent {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(outputDir, scheduledProcessEventLogName))
	if err != nil {
		t.Fatal(err)
	}
	var events []processEvent
	for _, line := range strings.Split(strings.TrimSpace(string(contents)), "\n") {
		var event processEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("decode process event %q: %v", line, err)
		}
		events = append(events, event)
	}
	return events
}

// interruptedValidationResume builds a validation output directory that a
// previous custodian process left with a checkpoint and a process-event log,
// plus a passing attestation the individual tests then mutate.
func interruptedValidationResume(t *testing.T, resumes int) (resumeAuthorization, scheduledResume) {
	t.Helper()
	repository := t.TempDir()
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	log := newProcessEventLog(output, "validation", "sha256:schedule", testBuild)
	if err := log.record(processEventStart, processEventProgress{}, "fresh scheduled execution"); err != nil {
		t.Fatal(err)
	}
	control := resumeAuthorization{
		tranche: "validation", outputRelative: "results", outputDir: output,
		attestationPath: filepath.Join(repository, "resume-attestation.json"),
		scheduleDigest:  "sha256:schedule", worldBuild: testBuild,
	}
	writeTestAttestation(t, control, func(map[string]any) {})
	resume := scheduledResume{hasCheckpoint: true, checkpoint: scheduledCheckpoint{V: 1, Resumes: resumes}}
	return control, resume
}

func writeTestAttestation(t *testing.T, control resumeAuthorization, mutate func(map[string]any)) {
	t.Helper()
	digest, err := digestProcessEventLog(control.outputDir)
	if err != nil {
		t.Fatal(err)
	}
	attestation := map[string]any{
		"v": 1, "tranche": "validation", "output_dir": control.outputRelative,
		"custodian": "Raoul Bivolaru (custodian)", "cause": "host_restart", "resume_index": 1,
		"interrupted_at":         "2026-08-20T00:00:00Z",
		"cause_classified_at":    "2026-08-20T00:05:00Z",
		"conditions_verified_at": "2026-08-20T00:10:00Z",
		"resume_authorized_at":   "2026-08-20T00:12:00Z",
		"no_outcome_inspection_between_interruption_and_resume":  true,
		"cause_classified_before_inspection_and_resume_decision": true,
		"frozen_bytes_digest_identical_at_resume":                true,
		"kill_principal_established_as_non_participant":          true,
		"process_event_log":        scheduledProcessEventLogName,
		"process_event_log_sha256": digest,
		"promptness_statement":     "resume attempted as soon as the host returned and the conditions verified",
		"lapse_explanation":        "two minutes elapsed while the frozen digests were recomputed",
		"frozen_digests_verified": map[string]any{
			"schedule_digest": control.scheduleDigest, "world_build": control.worldBuild,
		},
	}
	mutate(attestation)
	if err := writeJSON(control.attestationPath, attestation); err != nil {
		t.Fatal(err)
	}
}
