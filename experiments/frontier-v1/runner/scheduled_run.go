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
	resume, err := loadScheduledResume(config.outputDir, config.worldBuild, schedule, scheduleDigest, runBudgetUSD)
	if err != nil {
		return scheduledSummary{}, err
	}
	progress, err := openScheduledProgress(config, schedule, scheduleDigest, resume)
	if err != nil {
		return scheduledSummary{}, err
	}
	stopWatch := watchProcessTermination(progress.log, progress.live.get)
	defer stopWatch()
	done, err := applyScheduledResume(config, schedule, resume, &summary)
	if err != nil {
		return scheduledSummary{}, err
	}
	if done {
		return summary, progress.recordStop(resume.nextIndex, summary)
	}
	if resume.nextIndex == 0 {
		if err := progress.recordCheckpoint(config, schedule, scheduleDigest, summary, 0, 0); err != nil {
			return scheduledSummary{}, err
		}
	}
	runErr := runRemainingSchedule(config, schedule, scheduleDigest, resume.nextIndex, perTrialCap, runBudgetUSD, runtime, grader, progress, &summary)
	if runErr != nil {
		_ = progress.recordStop(summary.Launched, scheduledSummary{Status: "error"})
		return scheduledSummary{}, runErr
	}
	if err := progress.recordStop(len(schedule.Entries), summary); err != nil {
		return scheduledSummary{}, err
	}
	return summary, nil
}

// openScheduledProgress authorizes a resume under the protocol-v5 section 9
// mandatory-resume rule and opens this process's section of the append-only
// process-event log. A validation resume without a verified custodian
// attestation never reaches a launch.
func openScheduledProgress(
	config runConfig,
	schedule launchSchedule,
	scheduleDigest string,
	resume scheduledResume,
) (scheduledProgress, error) {
	control, err := scheduledResumeControl(config, scheduleDigest)
	if err != nil {
		return scheduledProgress{}, err
	}
	if err := requireResumeAuthorization(control, resume); err != nil {
		return scheduledProgress{}, err
	}
	progress := scheduledProgress{
		log:       newProcessEventLog(config.outputDir, effectiveTranche(config), scheduleDigest, config.worldBuild),
		resumes:   resume.checkpoint.Resumes,
		groupSize: len(schedule.Arms),
		live:      &liveProgress{},
	}
	detail := "fresh scheduled execution"
	if resume.hasCheckpoint {
		progress.resumes++
		detail = fmt.Sprintf("resume %d after an interrupted execution", progress.resumes)
	}
	completed := resume.nextIndex
	if progress.groupSize > 0 {
		completed = resume.nextIndex / progress.groupSize
	}
	spent := spentBefore(resume.results, resume.nextIndex)
	opening := progress.events(resume.nextIndex, completed, spent)
	progress.live.set(opening)
	return progress, progress.log.record(processEventStart, opening, detail)
}

// applyScheduledResume replays retained assignment records into the summary
// and reports whether the run is already terminal: a replayed budget or
// safety stop, or a schedule with no remaining launches.
func applyScheduledResume(config runConfig, schedule launchSchedule, resume scheduledResume, summary *scheduledSummary) (bool, error) {
	for _, result := range resume.results {
		replayAssignedResult(summary, result)
	}
	if resume.stopKind == "budget_stopped" {
		return true, recordBudgetStop(config, schedule.Entries[resume.nextIndex:], summary)
	}
	if resume.stopKind == "safety_stopped" {
		remaining := schedule.Entries[resume.nextIndex:]
		if len(remaining) > 0 {
			return true, recordSafetyStop(config, remaining, summary, fmt.Errorf("resumed incomplete safety stop"))
		}
		summary.Status = "indeterminate"
		if summary.StopReason == "" {
			summary.StopReason = "safety_stop: resumed incomplete safety stop"
		}
		return true, nil
	}
	return resume.nextIndex >= len(schedule.Entries), nil
}

// runRemainingSchedule executes the pairing-key groups at and after
// nextIndex. A budget or safety stop records the untouched remainder and
// ends the run; the returned error covers evidence-write failures only.
func runRemainingSchedule(
	config runConfig,
	schedule launchSchedule,
	scheduleDigest string,
	nextIndex int,
	perTrialCap, runBudgetUSD float64,
	runtime runtimeDriver,
	grader gradeDriver,
	progress scheduledProgress,
	summary *scheduledSummary,
) error {
	groupSize := len(schedule.Arms)
	for start := 0; start < len(schedule.Entries); start += groupSize {
		if start+groupSize <= nextIndex {
			continue
		}
		group := schedule.Entries[start : start+groupSize]
		if start >= nextIndex {
			pairingCapacity := perTrialCap * float64(groupSize) * infrastructureCapacityFactor
			if summary.SpentUSD+pairingCapacity > runBudgetUSD+1e-9 {
				return recordBudgetStop(config, schedule.Entries[start:], summary)
			}
		}
		for offset, entry := range group {
			if start+offset < nextIndex {
				continue
			}
			result, runErr := runAssignedTrial(config, entry, runtime, grader)
			if runErr != nil {
				return recordSafetyStop(config, schedule.Entries[start+offset:], summary, runErr)
			}
			addAssignedResult(summary, result)
		}
		if err := progress.recordCheckpoint(config, schedule, scheduleDigest, *summary, start+groupSize, (start+groupSize)/groupSize); err != nil {
			return err
		}
	}
	return nil
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

func replayAssignedResult(summary *scheduledSummary, result assignedTrialResult) {
	switch result.Status {
	case "budget_stopped":
		summary.Results = append(summary.Results, result)
		summary.BudgetStopped++
		summary.Unresolved++
		summary.Status = "indeterminate"
		if summary.StopReason == "" {
			summary.StopReason = "run_budget"
		}
	case "safety_stopped":
		summary.Results = append(summary.Results, result)
		summary.SafetyStopped++
		summary.Unresolved++
		summary.Status = "indeterminate"
		if summary.StopReason == "" {
			summary.StopReason = "safety_stop: resumed incomplete safety stop"
		}
	default:
		addAssignedResult(summary, result)
	}
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
