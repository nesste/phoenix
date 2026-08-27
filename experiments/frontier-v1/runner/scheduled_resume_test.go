package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRequireScheduledOutputReadyAllowsEmptyDirectory(t *testing.T) {
	repository := t.TempDir()
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := requireScheduledOutputReady(repository, "results", launchSchedule{Arms: append([]string(nil), phase1Arms...)}, "sha256:schedule"); err != nil {
		t.Fatal(err)
	}
}

func TestRequireScheduledOutputReadyRejectsUnrecognizedFiles(t *testing.T) {
	repository := t.TempDir()
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, "existing.json"), []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	err := requireScheduledOutputReady(repository, "results", launchSchedule{Arms: append([]string(nil), phase1Arms...)}, "sha256:schedule")
	if err == nil || !strings.Contains(err.Error(), "unrecognized files") {
		t.Fatalf("unrecognized output error = %v", err)
	}
}

func TestScheduledResumeSkipsCompletedPairingKeyAndReconstructsSpend(t *testing.T) {
	schedule := testPairingSchedule(2)
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
	var spent float64
	for _, entry := range schedule.Entries[:len(phase1Arms)] {
		writeTestAssignment(t, config, entry, "pass", 0.01)
		writeTestTrial(t, output, entry, testBuild)
		spent = roundUSD(spent + 0.01)
	}
	writeTestCheckpoint(t, config, "sha256:schedule", testBuild, 10, spent, len(phase1Arms), 1)
	driver := &scriptedRuntime{results: []runtimeResult{
		{Metrics: runtimeMetrics{CostUSD: 0.02, TotalTokens: 4}},
		{Metrics: runtimeMetrics{CostUSD: 0.02, TotalTokens: 4}},
		{Metrics: runtimeMetrics{CostUSD: 0.02, TotalTokens: 4}},
		{Metrics: runtimeMetrics{CostUSD: 0.02, TotalTokens: 4}},
	}}
	summary, err := runScheduledCases(config, schedule, "sha256:schedule", 10, driver, fakeGrader{})
	if err != nil {
		t.Fatal(err)
	}
	if summary.Status != "complete" || summary.Launched != 2*len(phase1Arms) || driver.calls != len(phase1Arms) ||
		summary.SpentUSD != roundUSD(0.04+0.08) || summary.Results[0].Status != "pass" ||
		summary.Results[len(phase1Arms)].Status != "pass" {
		t.Fatalf("resume summary = %#v, runtime calls = %d", summary, driver.calls)
	}
}

func TestScheduledResumeContinuesRemainingArmsOfIncompletePairingKey(t *testing.T) {
	schedule := testPairingSchedule(1)
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
	writeTestAssignment(t, config, schedule.Entries[0], "pass", 0.01)
	writeTestTrial(t, output, schedule.Entries[0], testBuild)
	writeTestAssignment(t, config, schedule.Entries[1], "fail", 0.02)
	writeTestTrial(t, output, schedule.Entries[1], testBuild)
	writeTestCheckpoint(t, config, "sha256:schedule", testBuild, 10, 0, 0, 0)
	driver := &scriptedRuntime{results: []runtimeResult{
		{Metrics: runtimeMetrics{CostUSD: 0.03, TotalTokens: 4}},
		{Metrics: runtimeMetrics{CostUSD: 0.03, TotalTokens: 4}},
	}}
	summary, err := runScheduledCases(config, schedule, "sha256:schedule", 10, driver, fakeGrader{})
	if err != nil {
		t.Fatal(err)
	}
	if summary.Status != "complete" || summary.Launched != len(phase1Arms) || driver.calls != 2 ||
		summary.Passes != 3 || summary.Failures != 1 || summary.SpentUSD != roundUSD(0.09) {
		t.Fatalf("incomplete-key resume = %#v, runtime calls = %d", summary, driver.calls)
	}
}

func TestScheduledResumeBudgetStopUsesReconstructedSpend(t *testing.T) {
	schedule := testPairingSchedule(2)
	repository := t.TempDir()
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	config := runConfig{repositoryRoot: repository, outputDir: output, budgetUSD: "0.15"}
	for _, entry := range schedule.Entries[:len(phase1Arms)] {
		writeTestAssignment(t, config, entry, "pass", 0.14)
	}
	writeTestCheckpoint(t, config, "sha256:schedule", "", 0.74, roundUSD(0.14*float64(len(phase1Arms))), len(phase1Arms), 1)
	driver := &scriptedRuntime{}
	summary, err := runScheduledCases(config, schedule, "sha256:schedule", 0.74, driver, fakeGrader{})
	if err != nil {
		t.Fatal(err)
	}
	if summary.Status != "indeterminate" || summary.StopReason != "run_budget" ||
		summary.Launched != len(phase1Arms) || summary.BudgetStopped != len(phase1Arms) || driver.calls != 0 {
		t.Fatalf("resumed budget summary = %#v, runtime calls = %d", summary, driver.calls)
	}
}

func TestScheduledResumeRejectsLaunchGap(t *testing.T) {
	schedule := testPairingSchedule(1)
	repository := t.TempDir()
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	config := runConfig{repositoryRoot: repository, outputDir: output}
	writeTestAssignment(t, config, schedule.Entries[0], "pass", 0.01)
	writeTestAssignment(t, config, schedule.Entries[2], "pass", 0.01)
	writeTestCheckpoint(t, config, "sha256:schedule", "", 10, 0, 0, 0)
	_, err := loadScheduledResume(output, "", schedule, "sha256:schedule", 10)
	if err == nil || !strings.Contains(err.Error(), "gap") {
		t.Fatalf("gap error = %v", err)
	}
}

func TestScheduledResumeRejectsScheduleIdentityMismatch(t *testing.T) {
	schedule := testPairingSchedule(1)
	repository := t.TempDir()
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	config := runConfig{repositoryRoot: repository, outputDir: output}
	writeTestAssignment(t, config, schedule.Entries[0], "pass", 0.01)
	schedule.Entries[0].CaseID = "other_case"
	_, err := loadScheduledResume(output, "", schedule, "sha256:schedule", 10)
	if err == nil || !strings.Contains(err.Error(), "does not match the selected schedule") {
		t.Fatalf("identity error = %v", err)
	}
}

func TestScheduledResumeRejectsFinishedSummary(t *testing.T) {
	repository := t.TempDir()
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := writeJSON(filepath.Join(output, "scheduled-summary.json"), scheduledSummary{V: 1}); err != nil {
		t.Fatal(err)
	}
	_, err := loadScheduledResume(output, "", testPairingSchedule(1), "sha256:schedule", 10)
	if err == nil || !strings.Contains(err.Error(), "finished summary") {
		t.Fatalf("finished summary error = %v", err)
	}
}

func TestScheduledResumeRejectsWorldBuildMismatch(t *testing.T) {
	schedule := testPairingSchedule(1)
	repository := t.TempDir()
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	config := runConfig{repositoryRoot: repository, outputDir: output}
	writeTestAssignment(t, config, schedule.Entries[0], "pass", 0.01)
	writeTestTrial(t, output, schedule.Entries[0], "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	writeTestCheckpoint(t, config, "sha256:schedule", testBuild, 10, 0.01, 0, 0)
	_, err := loadScheduledResume(output, testBuild, schedule, "sha256:schedule", 10)
	if err == nil || !strings.Contains(err.Error(), "world-build") {
		t.Fatalf("world-build error = %v", err)
	}
}

func TestScheduledResumeWritesCheckpointAfterPairingKey(t *testing.T) {
	schedule := testPairingSchedule(1)
	repository := t.TempDir()
	createRunnerFixture(t, repository)
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	driver := &scriptedRuntime{results: []runtimeResult{
		{Metrics: runtimeMetrics{CostUSD: 0.01, TotalTokens: 2}},
		{Metrics: runtimeMetrics{CostUSD: 0.01, TotalTokens: 2}},
		{Metrics: runtimeMetrics{CostUSD: 0.01, TotalTokens: 2}},
		{Metrics: runtimeMetrics{CostUSD: 0.01, TotalTokens: 2}},
	}}
	summary, err := runScheduledCases(runConfig{
		repositoryRoot: repository, outputDir: output, worldBuild: testBuild,
		worldPath: filepath.Join(repository, "world.json"), schemaPath: filepath.Join(repository, "schema.json"),
		timeout: time.Second, budgetUSD: "0.15",
	}, schedule, "sha256:schedule", 10, driver, fakeGrader{})
	if err != nil {
		t.Fatal(err)
	}
	var checkpoint scheduledCheckpoint
	if err := decodeStrict(filepath.Join(output, scheduledCheckpointName), &checkpoint); err != nil {
		t.Fatal(err)
	}
	if summary.Launched != len(phase1Arms) || checkpoint.NextLaunchIndex != len(phase1Arms) || checkpoint.CompletedPairingKeys != 1 ||
		checkpoint.ScheduleDigest != "sha256:schedule" || checkpoint.WorldBuild != testBuild {
		t.Fatalf("checkpoint = %#v, summary = %#v", checkpoint, summary)
	}
}

func TestScheduledResumeRejectsMissingCheckpoint(t *testing.T) {
	schedule := testPairingSchedule(1)
	repository := t.TempDir()
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	config := runConfig{repositoryRoot: repository, outputDir: output}
	writeTestAssignment(t, config, schedule.Entries[0], "pass", 0.01)
	_, err := loadScheduledResume(output, "", schedule, "sha256:schedule", 10)
	if err == nil || !strings.Contains(err.Error(), "missing a checkpoint") {
		t.Fatalf("missing checkpoint error = %v", err)
	}
}

func TestScheduledResumeRejectsScheduleDigestMismatchBeforeFirstCompletedKey(t *testing.T) {
	schedule := testPairingSchedule(2)
	repository := t.TempDir()
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	config := runConfig{repositoryRoot: repository, outputDir: output, worldBuild: testBuild}
	writeTestAssignment(t, config, schedule.Entries[0], "pass", 0.01)
	writeTestTrial(t, output, schedule.Entries[0], testBuild)
	writeTestCheckpoint(t, config, "sha256:schedule", testBuild, 10, 0, 0, 0)
	_, err := loadScheduledResume(output, testBuild, schedule, "sha256:other-schedule", 10)
	if err == nil || !strings.Contains(err.Error(), "checkpoint does not match") {
		t.Fatalf("pre-checkpoint digest error = %v", err)
	}
}

func TestScheduledResumeRejectsLaunchedAssignmentMissingTrial(t *testing.T) {
	schedule := testPairingSchedule(1)
	repository := t.TempDir()
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	config := runConfig{repositoryRoot: repository, outputDir: output, worldBuild: testBuild}
	writeTestAssignment(t, config, schedule.Entries[0], "fail", 0.01)
	writeTestCheckpoint(t, config, "sha256:schedule", testBuild, 10, 0, 0, 0)
	_, err := loadScheduledResume(output, testBuild, schedule, "sha256:schedule", 10)
	if err == nil || !strings.Contains(err.Error(), "missing world-build evidence") {
		t.Fatalf("missing trial error = %v", err)
	}
}

func TestScheduledResumeAcceptsCheckpointOnlyDirectory(t *testing.T) {
	schedule := testPairingSchedule(1)
	repository := t.TempDir()
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	config := runConfig{repositoryRoot: repository, outputDir: output, worldBuild: testBuild}
	writeTestCheckpoint(t, config, "sha256:schedule", testBuild, 10, 0, 0, 0)
	resume, err := loadScheduledResume(output, testBuild, schedule, "sha256:schedule", 10)
	if err != nil || resume.nextIndex != 0 || len(resume.results) != 0 {
		t.Fatalf("checkpoint-only resume = %#v, err = %v", resume, err)
	}
}

func TestScheduledResumeWritesInitialCheckpointBeforeFirstLaunch(t *testing.T) {
	schedule := testPairingSchedule(1)
	repository := t.TempDir()
	createRunnerFixture(t, repository)
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	summary, err := runScheduledCases(runConfig{
		repositoryRoot: repository, outputDir: output, worldBuild: testBuild,
		worldPath: filepath.Join(repository, "world.json"), schemaPath: filepath.Join(repository, "schema.json"),
		timeout: time.Second, budgetUSD: "0.15",
	}, schedule, "sha256:schedule", 10, &scriptedRuntime{}, fakeGrader{})
	if err != nil {
		t.Fatal(err)
	}
	var checkpoint scheduledCheckpoint
	if err := decodeStrict(filepath.Join(output, scheduledCheckpointName), &checkpoint); err != nil {
		t.Fatal(err)
	}
	if summary.Status != "indeterminate" || summary.SafetyStopped != len(phase1Arms) ||
		checkpoint.NextLaunchIndex != 0 || checkpoint.CompletedPairingKeys != 0 ||
		checkpoint.ScheduleDigest != "sha256:schedule" || checkpoint.WorldBuild != testBuild {
		t.Fatalf("initial checkpoint = %#v, summary = %#v", checkpoint, summary)
	}
}

func TestScheduledResumeRejectsOrphanRuntimeEvidence(t *testing.T) {
	repository := t.TempDir()
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, "000000-authoring_testcase-r01-A.attempt-01.runtime.jsonl"), []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := loadScheduledResume(output, "", testPairingSchedule(1), "sha256:schedule", 10)
	if err == nil || !strings.Contains(err.Error(), "evidence without a matching assignment") {
		t.Fatalf("orphan evidence error = %v", err)
	}
}

func testPairingSchedule(keys int) launchSchedule {
	schedule := launchSchedule{V: 1, Tranche: "authoring", Arms: append([]string(nil), phase1Arms...)}
	for pairing := 0; pairing < keys; pairing++ {
		for index, arm := range phase1Arms {
			schedule.Entries = append(schedule.Entries, scheduleEntry{
				LaunchIndex: pairing*len(phase1Arms) + index, PairingIndex: pairing,
				FamilyID: "authoring_family_test", CaseID: "authoring_testcase", Arm: arm,
			})
		}
	}
	return schedule
}

func writeTestCheckpoint(t *testing.T, config runConfig, digest, worldBuild string, budget, spent float64, next, keys int) {
	t.Helper()
	writeTestCheckpointWithResumes(t, config, digest, worldBuild, budget, spent, next, keys, 0)
}

func writeTestCheckpointWithResumes(
	t *testing.T,
	config runConfig,
	digest, worldBuild string,
	budget, spent float64,
	next, keys, resumes int,
) {
	t.Helper()
	if err := writeScheduledCheckpoint(config, digest, scheduledSummary{
		WorldBuild: worldBuild, RunBudgetUSD: budget, SpentUSD: spent,
	}, next, keys, resumes); err != nil {
		t.Fatal(err)
	}
}

func writeTestTrial(t *testing.T, output string, entry scheduleEntry, worldBuild string) {
	t.Helper()
	if err := writeJSON(filepath.Join(output, scheduledStem(entry)+".trial.json"), trial{
		V: 1, CaseID: entry.CaseID, WorldBuild: worldBuild,
	}); err != nil {
		t.Fatal(err)
	}
}

func writeTestAssignment(t *testing.T, config runConfig, entry scheduleEntry, status string, cost float64) assignedTrialResult {
	t.Helper()
	result := newAssignedResult(entry)
	result.Status = status
	result.ITTSuccess = status == "pass"
	result.Termination = "graded_" + status
	result.TotalCostUSD = cost
	written, err := writeAssignedResult(config, scheduledStem(entry), result)
	if err != nil {
		t.Fatal(err)
	}
	return written
}
