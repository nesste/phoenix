package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"
)

// outcomeFreePartialFields is the exact protocol-v5 section 9 field set of the
// partial summary. The proposal admits launch index, completed pairing keys,
// per-arm assignment counts, artifact digests, timestamps, and cumulative
// spend, and nothing else; the identity fields naming the tranche and the
// schedule carry no outcome either.
// outcomeFreeDigestFields is the admitted key set of the nested
// artifact_digests object. Every entry is an identity, never a measurement.
var outcomeFreeDigestFields = map[string]bool{
	"world_build": true, "grader_digest": true,
	"arm_b_document_digest": true, "process_event_log_sha256": true,
}

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
		}, "requires a written lapse explanation"},
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

// TestProcessTerminationStopsTheRun pins the section 9 host requirement that a
// supervised custodian process actually stops when signalled: recording the
// termination must not swallow it, or the log would assert a termination that
// did not happen while the run kept launching trials.
func TestProcessTerminationStopsTheRun(t *testing.T) {
	output := t.TempDir()
	progress := scheduledProgress{
		log:       newProcessEventLog(output, "validation", "sha256:schedule", testBuild),
		resumes:   1,
		groupSize: len(phase1Arms),
		live:      &liveProgress{},
	}
	// Publish through the same handle production uses, so the wiring is
	// covered and not just the handler.
	progress.live.set(progress.events(12, 3, 0.4))

	signals := make(chan os.Signal, 1)
	exited := make(chan int, 1)
	stop := watchProcessTerminationOn(progress.log, progress.live.get,
		func(code int) { exited <- code }, signals, func() {})
	defer stop()

	signals <- os.Interrupt
	select {
	case code := <-exited:
		if code != 130 {
			t.Fatalf("interrupt exit code = %d, want 130", code)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("termination handler did not stop the process")
	}
	events := readProcessEvents(t, output)
	if len(events) != 1 || events[0].Event != processEventTermination ||
		events[0].NextLaunchIndex != 12 || events[0].CompletedPairingKeys != 3 ||
		events[0].SpentUSD != 0.4 || events[0].Resumes != 1 {
		t.Fatalf("termination event = %#v", events)
	}
	if events[0].Principal.Exposed {
		t.Fatalf("termination principal = %#v", events[0].Principal)
	}
}

// TestProcessTerminationOnSIGTERMExitsWith143 pins the other frozen signal.
func TestProcessTerminationOnSIGTERMExitsWith143(t *testing.T) {
	output := t.TempDir()
	log := newProcessEventLog(output, "validation", "sha256:schedule", testBuild)
	signals := make(chan os.Signal, 1)
	exited := make(chan int, 1)
	stop := watchProcessTerminationOn(log, func() processEventProgress { return processEventProgress{} },
		func(code int) { exited <- code }, signals, func() {})
	defer stop()

	signals <- syscall.SIGTERM
	select {
	case code := <-exited:
		if code != 143 {
			t.Fatalf("terminate exit code = %d, want 143", code)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("termination handler did not stop the process")
	}
}

// TestValidationResumeCorroboratesTheCounterAgainstTheProcessEventLog covers
// the downward attack on the single-resume bound: the checkpoint's counter is
// one integer, so it is checked against the append-only log's start entries.
func TestValidationResumeCorroboratesTheCounterAgainstTheProcessEventLog(t *testing.T) {
	control, resume := interruptedValidationResume(t, maximumScheduledResumes)
	resume.checkpoint.Resumes = 0
	err := requireResumeAuthorization(control, resume)
	if err == nil || !strings.Contains(err.Error(), "the process-event log attests") {
		t.Fatalf("laundered counter error = %v", err)
	}

	control, resume = interruptedValidationResume(t, 0)
	if err := os.Remove(filepath.Join(control.outputDir, scheduledProcessEventLogName)); err != nil {
		t.Fatal(err)
	}
	if err := requireResumeAuthorization(control, resume); err == nil ||
		!strings.Contains(err.Error(), "read process-event log") {
		t.Fatalf("deleted log error = %v", err)
	}
}

// TestScheduledResumePersistsTheResumeCounter exercises the round trip end to
// end: a resumed run writes resumes=1 to the durable checkpoint, and a further
// resume is refused.
func TestScheduledResumePersistsTheResumeCounter(t *testing.T) {
	schedule, output, config := interruptedAuthoringRun(t)
	driver := &scriptedRuntime{results: repeatedRuntimeResults(len(schedule.Entries))}
	if _, err := runScheduledCases(config, schedule, "sha256:schedule", 10, driver, fakeGrader{}); err != nil {
		t.Fatal(err)
	}
	var checkpoint scheduledCheckpoint
	if err := decodeStrict(filepath.Join(output, scheduledCheckpointName), &checkpoint); err != nil {
		t.Fatal(err)
	}
	if checkpoint.Resumes != 1 {
		t.Fatalf("checkpoint resumes = %d, want 1", checkpoint.Resumes)
	}
	summary, err := readScheduledPartialSummary(output)
	if err != nil || summary.Resumes != 1 {
		t.Fatalf("partial summary resumes = %#v, %v", summary, err)
	}
	if err := verifyCheckpointProgress(scheduledCheckpoint{Resumes: 2}, len(phase1Arms), scheduledResume{}); err == nil {
		t.Fatal("a checkpoint claiming two resumes must be refused")
	}
}

func TestValidationResumeAnchorsItsWindowToTheProcessEventLogAndTheClock(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(map[string]any)
		now     time.Time
		message string
	}{
		{"interruption back-dated before the last process event", func(a map[string]any) {
			a["interrupted_at"] = "2026-08-19T00:00:00Z"
			a["cause_classified_at"] = "2026-08-19T00:05:00Z"
			a["conditions_verified_at"] = "2026-08-19T00:10:00Z"
			a["resume_authorized_at"] = "2026-08-19T00:12:00Z"
		}, resumeFixtureBase.Add(12 * time.Minute), "precedes the last process event"},
		{"authorization stamped in the future", func(map[string]any) {},
			resumeFixtureBase.Add(-time.Hour), "stamped in the future"},
		{"authorization stale at the moment of resume", func(map[string]any) {},
			resumeFixtureBase.Add(48 * time.Hour), "stale"},
		// The branch that actually closes the back-dating attack: an
		// attestation authored at resume time for a weeks-old interruption
		// passes the self-consistent checks and fails against the log.
		{"window measured from the last process event", func(a map[string]any) {
			a["interrupted_at"] = "2026-08-23T10:00:00Z"
			a["cause_classified_at"] = "2026-08-23T10:05:00Z"
			a["conditions_verified_at"] = "2026-08-23T10:10:00Z"
			a["resume_authorized_at"] = "2026-08-23T10:12:00Z"
		}, mustParseFixtureTime("2026-08-23T10:12:00Z"), "after the last process event"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			control, resume := interruptedValidationResume(t, 0)
			control.now = func() time.Time { return test.now }
			writeTestAttestation(t, control, test.mutate)
			err := requireResumeAuthorization(control, resume)
			if err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("%s error = %v", test.name, err)
			}
		})
	}
}

// TestValidationResumeAcceptsAnHonestlySkewedInterruptionStamp covers the
// signalled-interruption case: the termination entry is stamped after the
// signal arrives, so an honest custodian attesting the interruption instant
// itself writes a value slightly below the log's last entry.
func TestValidationResumeAcceptsAnHonestlySkewedInterruptionStamp(t *testing.T) {
	control, resume := interruptedValidationResume(t, 0)
	writeTestAttestation(t, control, func(a map[string]any) {
		a["interrupted_at"] = resumeFixtureBase.Add(-2 * time.Second).Format(time.RFC3339)
	})
	if err := requireResumeAuthorization(control, resume); err != nil {
		t.Fatalf("honestly skewed interruption stamp = %v", err)
	}

	control, resume = interruptedValidationResume(t, 0)
	writeTestAttestation(t, control, func(a map[string]any) {
		a["interrupted_at"] = resumeFixtureBase.Add(-24 * time.Hour).Format(time.RFC3339)
		a["cause_classified_at"] = resumeFixtureBase.Add(-23 * time.Hour).Format(time.RFC3339)
		a["conditions_verified_at"] = resumeFixtureBase.Add(-22 * time.Hour).Format(time.RFC3339)
	})
	if err := requireResumeAuthorization(control, resume); err == nil ||
		!strings.Contains(err.Error(), "precedes the last process event") {
		t.Fatalf("day-early interruption stamp = %v", err)
	}
}

// TestPromptResumeNeedsNoLapseExplanation keeps the lapse explanation a marker
// of a real lapse rather than a formality demanded on every resume.
func TestPromptResumeNeedsNoLapseExplanation(t *testing.T) {
	control, resume := interruptedValidationResume(t, 0)
	writeTestAttestation(t, control, func(a map[string]any) { delete(a, "lapse_explanation") })
	if err := requireResumeAuthorization(control, resume); err != nil {
		t.Fatalf("prompt resume without a lapse explanation = %v", err)
	}
}

func mustParseFixtureTime(value string) time.Time {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	return parsed
}

func TestScheduledStopKindNeverCarriesTheRawFailureText(t *testing.T) {
	cases := map[string]string{
		"":           "",
		"run_budget": "run_budget",
		`safety_stop: custody-safe grade status "fail" is inconsistent with its checks`: "safety_stop",
		"something else": "other",
	}
	for reason, want := range cases {
		if actual := scheduledStopKind(reason); actual != want {
			t.Fatalf("stop kind for %q = %q, want %q", reason, actual, want)
		}
	}
}

func TestPostClosureGateRefusesValidationAfterANonCompletedExecution(t *testing.T) {
	root := openValidationGateTestRoot(t, filepath.Join("..", "..", ".."))
	gate := closureGateFixture(t, root)
	digest := writeTestClosureReview(t, root, "an independently reviewed closure record\n")
	fields := writeTestReviewRecord(t, root, "fields-review.md", "the review of the commit that set the fields\n")

	for _, status := range []string{"indeterminate", "budget_stopped", "lapsed"} {
		err := gate(status, "", "", "", "")
		if err == nil || !strings.Contains(err.Error(), "set by a reviewed payload commit") {
			t.Fatalf("post-closure gate error for %s = %v", status, err)
		}
	}
	for _, missing := range []struct {
		name                            string
		review, verdict, digest, fields string
	}{
		{"unaccepted verdict", "docs/reviews/closure.md", "REVISE", digest, fields},
		{"unpinned record", "docs/reviews/closure.md", "ACCEPT", "", fields},
		{"unnamed fields review", "docs/reviews/closure.md", "ACCEPT", digest, ""},
	} {
		err := gate("indeterminate", missing.review, missing.verdict, missing.digest, missing.fields)
		if err == nil || !strings.Contains(err.Error(), "set by a reviewed payload commit") {
			t.Fatalf("%s = %v", missing.name, err)
		}
	}
	if err := gate("indeterminate", "docs/reviews/closure.md", "ACCEPT", digest, fields); err != nil {
		t.Fatalf("reviewed, pinned and attributed closure record = %v", err)
	}
	if err := gate("complete", "", "", "", ""); err != nil {
		t.Fatalf("completed execution = %v", err)
	}
}

// TestPostClosureGateVerifiesTheClosureRecordIdentity covers the record side
// of the gate. The runner verifies identity, not prose: the named record must
// be exactly the bytes the independent reviewer accepted. That refuses a
// record edited after the review that blessed it, which reading the record's
// current text cannot detect.
func TestPostClosureGateVerifiesTheClosureRecordIdentity(t *testing.T) {
	root := openValidationGateTestRoot(t, filepath.Join("..", "..", ".."))
	gate := closureGateFixture(t, root)
	const record = "docs/reviews/closure.md"
	digest := writeTestClosureReview(t, root, "the reviewed closure record\n")
	fields := writeTestReviewRecord(t, root, "fields-review.md", "the review of the commit that set the fields\n")

	for _, refused := range []struct {
		name, review, digest, message string
	}{
		{"outside docs/reviews", "README.md", digest, "committed review under docs/reviews"},
		{"escaping docs/reviews", "docs/reviews/../README.md", digest, "committed review under docs/reviews"},
		{"absent record", "docs/reviews/absent.md", digest, "closure review record"},
		{"failing its pin", record, "sha256:" + strings.Repeat("0", 64), "does not match the pinned"},
	} {
		err := gate("indeterminate", refused.review, "ACCEPT", refused.digest, fields)
		if err == nil || !strings.Contains(err.Error(), refused.message) {
			t.Fatalf("%s = %v", refused.name, err)
		}
	}
	if err := gate("indeterminate", record, "ACCEPT", digest, fields); err != nil {
		t.Fatalf("pinned closure record = %v", err)
	}

	// The decisive property the deleted verdict parser could not provide: a
	// record edited after the review that accepted it no longer matches.
	writeTestClosureReview(t, root, "the reviewed closure record, quietly amended\n")
	if err := gate("indeterminate", record, "ACCEPT", digest, fields); err == nil ||
		!strings.Contains(err.Error(), "does not match the pinned") {
		t.Fatalf("amended closure record = %v", err)
	}
	writeTestClosureReview(t, root, "")
	if err := gate("indeterminate", record, "ACCEPT", digest, fields); err == nil ||
		!strings.Contains(err.Error(), "non-empty regular file") {
		t.Fatalf("emptied closure record = %v", err)
	}
}

// TestPostClosureGateRequiresTheFieldsReviewToBeAResolvableRecord pins the
// sixth review's N1 repair: the gate must name the independent review of the
// commit that set the closure fields, so a chair gate patch setting them
// directly leaves a refusal rather than a silent gap.
func TestPostClosureGateRequiresTheFieldsReviewToBeAResolvableRecord(t *testing.T) {
	root := openValidationGateTestRoot(t, filepath.Join("..", "..", ".."))
	gate := closureGateFixture(t, root)
	digest := writeTestClosureReview(t, root, "the reviewed closure record\n")
	for _, refused := range []struct{ name, fields, message string }{
		{"outside docs/reviews", "README.md", "name their independent review under docs/reviews"},
		{"escaping docs/reviews", "docs/reviews/../README.md", "name their independent review under docs/reviews"},
		{"absent", "docs/reviews/absent.md", "closure gate fields review record"},
	} {
		err := gate("indeterminate", "docs/reviews/closure.md", "ACCEPT", digest, refused.fields)
		if err == nil || !strings.Contains(err.Error(), refused.message) {
			t.Fatalf("%s = %v", refused.name, err)
		}
	}
	empty := writeTestReviewRecord(t, root, "empty-review.md", "")
	if err := gate("indeterminate", "docs/reviews/closure.md", "ACCEPT", digest, empty); err == nil ||
		!strings.Contains(err.Error(), "non-empty regular file") {
		t.Fatalf("empty fields review = %v", err)
	}
}

// closureGateFixture returns a setter that rewrites the post-closure fields of
// a test gate document and reports what the runner's gate check then says.
func closureGateFixture(t *testing.T, root string) func(status, review, verdict, digest, fieldsReview string) error {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(preValidationArtifactsPath))
	var document map[string]any
	readJSONForTest(t, path, &document)
	gates := document["gates"].(map[string]any)
	return func(status, review, verdict, digest, fieldsReview string) error {
		gates["validation_execution_status"] = status
		gates["validation_execution_closure_review"] = review
		gates["validation_execution_closure_review_verdict"] = verdict
		gates["validation_execution_closure_review_lf_normalized_utf8_sha256"] = digest
		gates["validation_execution_closure_fields_payload_review"] = fieldsReview
		if err := writeJSON(path, document); err != nil {
			t.Fatal(err)
		}
		return requireValidationGate(root)
	}
}

// writeTestClosureReview writes a closure record and returns its pinned
// digest, computed with the same LF-normalizing rule the gate uses.
func writeTestClosureReview(t *testing.T, root, contents string) string {
	t.Helper()
	path := filepath.Join(root, "docs", "reviews", "closure.md")
	writeTestReviewRecord(t, root, "closure.md", contents)
	if contents == "" {
		return ""
	}
	digest, err := digestLFNormalizedFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return digest
}

// writeTestReviewRecord writes a document under docs/reviews and returns its
// repository-relative path.
func writeTestReviewRecord(t *testing.T, root, name, contents string) string {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, "docs", "reviews"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "reviews", name), []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	return "docs/reviews/" + name
}

// resumeFixtureBase is the instant the fixture's interrupted custodian process
// last wrote to its process-event log. The attested timestamps and the
// authorization clock are anchored to it.
var resumeFixtureBase = time.Date(2026, 8, 20, 0, 0, 0, 0, time.UTC)

// assertPartialSummaryFieldsAreOutcomeFree checks the written JSON object
// against the admitted section 9 field set in both directions, and descends
// into every nested object so a later field added inside artifact_digests
// cannot ship an outcome past a top-level-only guard.
func assertPartialSummaryFieldsAreOutcomeFree(t *testing.T, path string) {
	t.Helper()
	var raw map[string]json.RawMessage
	readJSONForTest(t, path, &raw)
	assertExactFieldSet(t, "partial summary", raw, outcomeFreePartialFields)

	var digests map[string]json.RawMessage
	if err := json.Unmarshal(raw["artifact_digests"], &digests); err != nil {
		t.Fatal(err)
	}
	assertExactFieldSet(t, "artifact_digests", digests, outcomeFreeDigestFields)

	var counts map[string]int
	if err := json.Unmarshal(raw["arm_assignment_counts"], &counts); err != nil {
		t.Fatal(err)
	}
	for arm := range counts {
		if !slices.Contains(phase1Arms, arm) {
			t.Fatalf("arm_assignment_counts carries the non-arm key %q", arm)
		}
	}
}

func assertExactFieldSet(t *testing.T, label string, actual map[string]json.RawMessage, admitted map[string]bool) {
	t.Helper()
	for field := range actual {
		if !admitted[field] {
			t.Fatalf("%s carries the non-admitted field %q", label, field)
		}
	}
	for field := range admitted {
		if _, present := actual[field]; !present {
			t.Fatalf("%s is missing the required field %q", label, field)
		}
	}
}

// interruptedAuthoringRun leaves a two-pairing-key schedule after its first
// completed pairing key, with the checkpoint, partial summary, and
// process-event log an interrupted custodian process would have written.
func interruptedAuthoringRun(t *testing.T) (launchSchedule, string, runConfig) {
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
	first := schedule.Entries[:len(phase1Arms)]
	spent := 0.0
	for _, entry := range first {
		writeTestAssignment(t, config, entry, "pass", 0.01)
		writeTestTrial(t, output, entry, testBuild)
		spent = roundUSD(spent + 0.01)
	}
	writeTestCheckpoint(t, config, "sha256:schedule", testBuild, 10, spent, len(first), 1)
	log := newProcessEventLog(output, "authoring", "sha256:schedule", testBuild)
	if err := log.record(processEventStart, processEventProgress{}, "fresh scheduled execution"); err != nil {
		t.Fatal(err)
	}
	return schedule, output, config
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
		// The validation path always carries these two, so the outcome-free
		// guard must see the shape it is meant to guard.
		graderDigest: "sha256:grader", armBDocumentDigest: "sha256:armb",
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
	log.clock = func() time.Time { return resumeFixtureBase }
	for start := 0; start <= resumes; start++ {
		if err := log.record(processEventStart, processEventProgress{resumes: start}, "scheduled execution"); err != nil {
			t.Fatal(err)
		}
	}
	control := resumeAuthorization{
		tranche: "validation", outputRelative: "results", outputDir: output,
		attestationPath: filepath.Join(repository, "resume-attestation.json"),
		scheduleDigest:  "sha256:schedule", worldBuild: testBuild,
		now: func() time.Time { return resumeFixtureBase.Add(12 * time.Minute) },
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
