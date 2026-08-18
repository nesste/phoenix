package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAssignedTrialRetriesEligibleInfrastructureWithFreshIsolation(t *testing.T) {
	repository := t.TempDir()
	createRunnerFixture(t, repository)
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	driver := &scriptedRuntime{results: []runtimeResult{{
		Failure: &runtimeFailure{
			Code: "provider_transient", Message: "provider status 503",
			BeforeFirstModelToken: true, RetryEligible: true,
		},
	}, {Metrics: runtimeMetrics{InputTokens: 10, OutputTokens: 5, TotalTokens: 15, CostUSD: 0.01}}}}
	result, err := runAssignedTrial(runConfig{
		armBDocument: filepath.Join(repository, "arm-b.md"), repositoryRoot: repository, outputDir: output,
		worldPath: filepath.Join(repository, "world.json"), schemaPath: filepath.Join(repository, "schema.json"),
		worldBuild: testBuild, timeout: time.Second, budgetUSD: "0.15",
	}, scheduleEntry{
		LaunchIndex: 1, FamilyBlock: 0, PairingIndex: 0,
		FamilyID: "authoring_family_test", CaseID: "authoring_testcase", Repetition: 0, Arm: "C",
	}, driver, fakeGrader{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "pass" || !result.ITTSuccess || len(result.Attempts) != 2 || result.TotalTokens != 15 {
		t.Fatalf("assigned result = %#v", result)
	}
	if result.Attempts[0].Status != "retryable_infrastructure" || result.Attempts[1].Status != "completed" {
		t.Fatalf("attempts = %#v", result.Attempts)
	}
	if len(driver.sandboxes) != 2 || driver.sandboxes[0] == driver.sandboxes[1] || driver.roots[0] == driver.roots[1] {
		t.Fatalf("retry isolation sandboxes=%v roots=%v", driver.sandboxes, driver.roots)
	}
	for _, sandbox := range driver.sandboxes {
		if _, err := os.Stat(sandbox); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("attempt sandbox survived: %s, %v", sandbox, err)
		}
	}
	if _, err := os.Stat(filepath.Join(repository, filepath.FromSlash(result.AssignmentPath))); err != nil {
		t.Fatalf("assignment record: %v", err)
	}
}

func TestAssignedTrialDoesNotRetryCapHit(t *testing.T) {
	repository := t.TempDir()
	createRunnerFixture(t, repository)
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	driver := &scriptedRuntime{results: []runtimeResult{{
		Metrics: runtimeMetrics{CostUSD: 0.15},
		Failure: &runtimeFailure{Code: "cost_cap", Message: "cost cap", CapHit: true},
	}}}
	result, err := runAssignedTrial(runConfig{
		repositoryRoot: repository, outputDir: output,
		worldPath: filepath.Join(repository, "world.json"), schemaPath: filepath.Join(repository, "schema.json"),
		worldBuild: testBuild, timeout: time.Second, budgetUSD: "0.15",
	}, scheduleEntry{FamilyID: "authoring_family_test", CaseID: "authoring_testcase", Arm: "C"}, driver, fakeGrader{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "fail" || result.Termination != "cost_cap" || !result.CapHit || len(result.Attempts) != 1 {
		t.Fatalf("cap result = %#v", result)
	}
}

func TestAssignedTrialRecordsUnresolvedAfterRetryExhaustion(t *testing.T) {
	repository := t.TempDir()
	createRunnerFixture(t, repository)
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	failure := &runtimeFailure{
		Code: "mcp_connection", Message: "connection failed",
		BeforeFirstModelToken: true, RetryEligible: true,
	}
	driver := &scriptedRuntime{results: []runtimeResult{{Failure: failure}, {Failure: failure}}}
	result, err := runAssignedTrial(runConfig{
		repositoryRoot: repository, outputDir: output,
		worldPath: filepath.Join(repository, "world.json"), schemaPath: filepath.Join(repository, "schema.json"),
		worldBuild: testBuild, timeout: time.Second, budgetUSD: "0.15",
	}, scheduleEntry{FamilyID: "authoring_family_test", CaseID: "authoring_testcase", Arm: "C"}, driver, fakeGrader{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "unresolved" || result.ITTSuccess || result.Termination != "infrastructure_unresolved" ||
		len(result.Attempts) != 2 || result.Attempts[1].Status != "unresolved_infrastructure" {
		t.Fatalf("unresolved result = %#v", result)
	}
}

func TestScheduledBudgetStopOccursBeforeWholePairingKey(t *testing.T) {
	schedule := launchSchedule{Arms: append([]string(nil), phase1Arms...)}
	for index, arm := range phase1Arms {
		schedule.Entries = append(schedule.Entries, scheduleEntry{
			LaunchIndex: index, PairingIndex: 0, FamilyID: "family", CaseID: "authoring_case", Arm: arm,
		})
	}
	repository := t.TempDir()
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	driver := &scriptedRuntime{}
	summary, err := runScheduledCases(runConfig{
		repositoryRoot: repository, outputDir: output, budgetUSD: "0.15",
	}, schedule, "sha256:schedule", 0.74, driver, fakeGrader{})
	if err != nil {
		t.Fatal(err)
	}
	if summary.Status != "indeterminate" || summary.StopReason != "run_budget" ||
		summary.Launched != 0 || summary.BudgetStopped != len(phase1Arms) ||
		summary.Unresolved != len(phase1Arms) || driver.calls != 0 {
		t.Fatalf("budget summary = %#v, runtime calls = %d", summary, driver.calls)
	}
}

func TestScheduledSummaryPinsRunConfigurationAndPrompts(t *testing.T) {
	configuration := scheduledRunConfiguration(runConfig{
		repositoryRoot: "C:/repo", armBDocument: "C:/repo/experiments/frontier-v1/arms/arm-b.md",
		flatToolNames: []string{"repo_status", "tests_run"}, timeout: 180 * time.Second,
	}, launchSchedule{Seed: phase1ScheduleSeed, Repetitions: phase1Repetitions}, 0.15)
	if configuration.RuntimeVersion != pinnedRuntimeVersion || configuration.Model != pinnedModel ||
		configuration.MaxTurns != 12 || configuration.TimeoutSeconds != 180 ||
		configuration.MaximumInfrastructureRetry != 1 || configuration.SystemPrompts["A"] != baseSystemPrompt ||
		configuration.SystemPrompts["C"] != baseSystemPrompt+" "+phoenixIntentPrompt ||
		len(configuration.AllowedTools["A"]) != 2 || configuration.AllowedTools["C"][0] != "mcp__phoenix__act" {
		t.Fatalf("scheduled configuration = %#v", configuration)
	}
}

type scriptedRuntime struct {
	results   []runtimeResult
	calls     int
	sandboxes []string
	roots     []string
}

func (runtime *scriptedRuntime) Verify() error { return nil }

func (runtime *scriptedRuntime) Run(request runtimeRequest) (runtimeResult, error) {
	runtime.sandboxes = append(runtime.sandboxes, request.Sandbox)
	runtime.roots = append(runtime.roots, request.Roots["repo"])
	index := runtime.calls
	runtime.calls++
	if index >= len(runtime.results) {
		return runtimeResult{}, errors.New("unexpected runtime call")
	}
	result := runtime.results[index]
	if result.Failure != nil {
		return result, nil
	}
	completed, err := (&fakeRuntime{}).Run(request)
	if err != nil {
		return runtimeResult{}, err
	}
	completed.Metrics = result.Metrics
	return completed, nil
}
