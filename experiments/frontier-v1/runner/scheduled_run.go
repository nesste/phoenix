package main

import (
	"fmt"
	"math"
	"path/filepath"
	"strconv"
)

const (
	maximumInfrastructureRetries = 1
	infrastructureCapacityFactor = 1.10
)

func runScheduledCases(
	config runConfig,
	schedule launchSchedule,
	scheduleDigest string,
	runBudgetUSD float64,
	runtime runtimeDriver,
	grader gradeDriver,
) (scheduledSummary, error) {
	perTrialCap, err := strconv.ParseFloat(config.budgetUSD, 64)
	if err != nil || perTrialCap <= 0 || runBudgetUSD <= 0 {
		return scheduledSummary{}, fmt.Errorf("scheduled run budgets must be positive decimal USD values")
	}
	summary := scheduledSummary{
		V: 1, Tranche: schedule.Tranche, ScheduleDigest: scheduleDigest, WorldBuild: config.worldBuild,
		Configuration: scheduledRunConfiguration(config, schedule, perTrialCap),
		Status:        "complete", RunBudgetUSD: runBudgetUSD, Assigned: len(schedule.Entries), Results: []assignedTrialResult{},
	}
	groupSize := len(schedule.Arms)
	if groupSize == 0 || len(schedule.Entries)%groupSize != 0 {
		return scheduledSummary{}, fmt.Errorf("schedule entries must contain complete pairing-key arm groups")
	}
	for start := 0; start < len(schedule.Entries); start += groupSize {
		group := schedule.Entries[start : start+groupSize]
		pairingCapacity := perTrialCap * float64(groupSize) * infrastructureCapacityFactor
		if summary.SpentUSD+pairingCapacity > runBudgetUSD+1e-9 {
			if err := recordBudgetStop(config, schedule.Entries[start:], &summary); err != nil {
				return scheduledSummary{}, err
			}
			break
		}
		for offset, entry := range group {
			result, runErr := runAssignedTrial(config, entry, runtime, grader)
			if runErr != nil {
				if err := recordSafetyStop(config, schedule.Entries[start+offset:], &summary, runErr); err != nil {
					return scheduledSummary{}, err
				}
				return summary, nil
			}
			addAssignedResult(&summary, result)
		}
	}
	return summary, nil
}

func scheduledRunConfiguration(config runConfig, schedule launchSchedule, perTrialCap float64) scheduledConfiguration {
	prompts := make(map[string]string, len(phase1Arms))
	allowedTools := make(map[string][]string, len(phase1Arms))
	for _, arm := range phase1Arms {
		prompt, _ := systemPromptForArm(arm)
		prompts[arm] = prompt
		allowedTools[arm] = allowedToolsForArm(arm, config.flatToolNames)
	}
	armBDocument := config.armBDocument
	if relative, err := filepath.Rel(config.repositoryRoot, armBDocument); err == nil {
		armBDocument = filepath.ToSlash(relative)
	}
	return scheduledConfiguration{
		Runtime: "claude-code", RuntimeVersion: pinnedRuntimeVersion, Model: pinnedModel,
		Effort: "low", ServiceTier: "standard", AccessMode: "surface_only", Transport: "stdio",
		OutputFormat: "stream-json", BuiltinTools: "empty", StrictMCPConfig: true,
		SessionPersistence: false, ModelSeedSupport: false, MaxTurns: 12,
		TimeoutSeconds: int(config.timeout.Seconds()), MaxCostUSDPerTrial: perTrialCap,
		MaximumInfrastructureRetry: maximumInfrastructureRetries,
		ScheduleSeed:               schedule.Seed, Repetitions: schedule.Repetitions,
		SystemPrompts: prompts, AllowedTools: allowedTools, ArmBDocument: armBDocument,
		ArmBDocumentDigest: config.armBDocumentDigest, GraderDigest: config.graderDigest,
		GraderBoundary: config.graderBoundary, GraderAdapterSHA256: config.graderAdapterHash,
		TokenAccounting: "Claude result usage: input, cache-creation input, cache-read input, and output tokens; total is their sum",
	}
}

func runAssignedTrial(config runConfig, entry scheduleEntry, runtime runtimeDriver, grader gradeDriver) (assignedTrialResult, error) {
	result := newAssignedResult(entry)
	entryConfig := config
	entryConfig.arm = entry.Arm
	stem := scheduledStem(entry)
	for attemptIndex := 0; attemptIndex <= maximumInfrastructureRetries; attemptIndex++ {
		workspace, runtimeResult, err := executeCaseAttempt(entryConfig, entry.CaseID, runtime)
		if err != nil {
			return assignedTrialResult{}, err
		}
		runtimeName := fmt.Sprintf("%s.attempt-%02d.runtime.jsonl", stem, attemptIndex+1)
		runtimePath, err := writeRuntimeEvidence(entryConfig, runtimeName, runtimeResult.RawOutput)
		if err != nil {
			workspace.cleanup()
			return assignedTrialResult{}, err
		}
		record := attemptRecord{
			Attempt: attemptIndex + 1, Status: "completed", RuntimePath: runtimePath,
			FirstModelToken: runtimeResult.FirstModelToken,
			Metrics:         runtimeResult.Metrics, Failure: runtimeResult.Failure,
		}
		result.TotalCostUSD += runtimeResult.Metrics.CostUSD
		result.TotalTokens += runtimeResult.Metrics.TotalTokens
		if runtimeResult.Failure != nil {
			workspace.cleanup()
			result.Attempts = append(result.Attempts, recordFailedAttempt(record, attemptIndex))
			if runtimeResult.Failure.RetryEligible && attemptIndex < maximumInfrastructureRetries {
				continue
			}
			finishRuntimeFailure(&result, runtimeResult.Failure)
			return writeAssignedResult(entryConfig, stem, result)
		}
		caseResult, err := finishCase(
			entryConfig, stem, entry.CaseID, workspace.sandbox, workspace.stateDir, workspace.roots, runtimeResult, grader,
		)
		workspace.cleanup()
		if err != nil {
			return assignedTrialResult{}, err
		}
		result.Attempts = append(result.Attempts, record)
		result.TrialPath = caseResult.TrialPath
		result.GradePath = caseResult.GradePath
		finishGrade(&result, caseResult.Status)
		return writeAssignedResult(entryConfig, stem, result)
	}
	return assignedTrialResult{}, fmt.Errorf("infrastructure retry loop ended without a terminal result")
}

func newAssignedResult(entry scheduleEntry) assignedTrialResult {
	return assignedTrialResult{
		V: 1, LaunchIndex: entry.LaunchIndex, FamilyBlock: entry.FamilyBlock, PairingIndex: entry.PairingIndex,
		FamilyID: entry.FamilyID, CaseID: entry.CaseID, Repetition: entry.Repetition, Arm: entry.Arm,
		Status: "pending", Attempts: []attemptRecord{},
	}
}

func recordFailedAttempt(record attemptRecord, attemptIndex int) attemptRecord {
	if record.Failure.RetryEligible && attemptIndex < maximumInfrastructureRetries {
		record.Status = "retryable_infrastructure"
	} else if record.Failure.RetryEligible {
		record.Status = "unresolved_infrastructure"
	} else {
		record.Status = "terminal_failure"
	}
	return record
}

func finishRuntimeFailure(result *assignedTrialResult, failure *runtimeFailure) {
	result.ITTSuccess = false
	result.CapHit = failure.CapHit
	result.Termination = failure.Code
	if failure.RetryEligible {
		result.Status = "unresolved"
		result.Termination = "infrastructure_unresolved"
		return
	}
	result.Status = "fail"
}

func finishGrade(result *assignedTrialResult, grade string) {
	result.ITTSuccess = grade == "pass"
	if grade == "pass" {
		result.Status = "pass"
		result.Termination = "graded_pass"
		return
	}
	result.Status = "fail"
	if grade == "indeterminate" {
		result.Termination = "manual_required"
	} else {
		result.Termination = "graded_fail"
	}
}

func writeAssignedResult(config runConfig, stem string, result assignedTrialResult) (assignedTrialResult, error) {
	path := filepath.Join(config.outputDir, stem+".assignment.json")
	relative, err := filepath.Rel(config.repositoryRoot, path)
	if err != nil {
		return assignedTrialResult{}, err
	}
	result.AssignmentPath = filepath.ToSlash(relative)
	if err := writeJSON(path, result); err != nil {
		return assignedTrialResult{}, err
	}
	return result, nil
}

func recordBudgetStop(config runConfig, entries []scheduleEntry, summary *scheduledSummary) error {
	summary.Status = "indeterminate"
	summary.StopReason = "run_budget"
	for _, entry := range entries {
		result := newAssignedResult(entry)
		result.Status = "budget_stopped"
		result.Termination = "budget_stop"
		written, err := writeAssignedResult(config, scheduledStem(entry), result)
		if err != nil {
			return err
		}
		summary.Results = append(summary.Results, written)
		summary.BudgetStopped++
		summary.Unresolved++
	}
	return nil
}

func recordSafetyStop(config runConfig, entries []scheduleEntry, summary *scheduledSummary, cause error) error {
	summary.Status = "indeterminate"
	summary.StopReason = "safety_stop: " + cause.Error()
	for _, entry := range entries {
		result := newAssignedResult(entry)
		result.Status = "safety_stopped"
		result.Termination = "safety_stop"
		written, err := writeAssignedResult(config, scheduledStem(entry), result)
		if err != nil {
			return err
		}
		summary.Results = append(summary.Results, written)
		summary.SafetyStopped++
		summary.Unresolved++
	}
	return nil
}

func addAssignedResult(summary *scheduledSummary, result assignedTrialResult) {
	summary.Results = append(summary.Results, result)
	summary.Launched++
	summary.SpentUSD = roundUSD(summary.SpentUSD + result.TotalCostUSD)
	if result.CapHit {
		summary.CapHits++
	}
	switch result.Status {
	case "pass":
		summary.Passes++
	case "unresolved":
		summary.Unresolved++
	default:
		summary.Failures++
	}
}

func roundUSD(value float64) float64 {
	return math.Round(value*1e9) / 1e9
}

func scheduledStem(entry scheduleEntry) string {
	return fmt.Sprintf("%06d-%s-r%02d-%s", entry.LaunchIndex, entry.CaseID, entry.Repetition+1, entry.Arm)
}
