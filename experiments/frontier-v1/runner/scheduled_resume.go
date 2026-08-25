package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const scheduledCheckpointName = "scheduled-checkpoint.json"
const scheduledSummaryName = "scheduled-summary.json"

type scheduledResume struct {
	results   []assignedTrialResult
	nextIndex int
	stopKind  string
}

func requireScheduledOutputReady(root, output string, schedule launchSchedule, scheduleDigest string) error {
	path := resolveRepositoryPath(root, output)
	if err := ensureInside(root, path); err != nil {
		return err
	}
	_, err := loadScheduledResume(path, "", schedule, scheduleDigest, 0)
	return err
}

func loadScheduledResume(
	outputDir, worldBuild string,
	schedule launchSchedule,
	scheduleDigest string,
	runBudgetUSD float64,
) (scheduledResume, error) {
	entries, err := os.ReadDir(outputDir)
	if errors.Is(err, os.ErrNotExist) {
		return scheduledResume{}, nil
	}
	if err != nil {
		return scheduledResume{}, err
	}
	if len(entries) == 0 {
		return scheduledResume{}, nil
	}

	groupSize := len(schedule.Arms)
	if groupSize == 0 || len(schedule.Entries)%groupSize != 0 {
		return scheduledResume{}, fmt.Errorf("schedule entries must contain complete pairing-key arm groups")
	}

	inventory, err := classifyScheduledOutputEntries(outputDir, entries, schedule)
	if err != nil {
		return scheduledResume{}, err
	}
	resume, err := reconstructScheduledResume(schedule, inventory.assignments)
	if err != nil {
		return scheduledResume{}, err
	}
	if err := verifyResumeWorldBuild(outputDir, worldBuild, schedule, resume); err != nil {
		return scheduledResume{}, err
	}
	if inventory.hasCheckpoint {
		if err := verifyScheduledCheckpoint(
			filepath.Join(outputDir, scheduledCheckpointName),
			worldBuild, scheduleDigest, runBudgetUSD, groupSize, resume,
		); err != nil {
			return scheduledResume{}, err
		}
	}
	return resume, nil
}

type scheduledOutputInventory struct {
	assignments   map[int]assignedTrialResult
	hasCheckpoint bool
}

// classifyScheduledOutputEntries sorts a non-empty scheduled output directory
// into assignment records, matching evidence, and the checkpoint. Anything
// unrecognized, duplicated, orphaned, or missing its checkpoint is refused.
func classifyScheduledOutputEntries(outputDir string, entries []os.DirEntry, schedule launchSchedule) (scheduledOutputInventory, error) {
	inventory := scheduledOutputInventory{assignments: map[int]assignedTrialResult{}}
	assignmentStems := map[string]struct{}{}
	var evidence []string
	for _, entry := range entries {
		if entry.IsDir() {
			return scheduledOutputInventory{}, fmt.Errorf("scheduled output directory contains unrecognized files")
		}
		name := entry.Name()
		switch {
		case name == scheduledSummaryName:
			return scheduledOutputInventory{}, fmt.Errorf("scheduled output directory already contains a finished summary")
		case name == scheduledCheckpointName:
			inventory.hasCheckpoint = true
		case strings.HasSuffix(name, ".assignment.json"):
			result, err := loadAssignmentFile(filepath.Join(outputDir, name), name, schedule)
			if err != nil {
				return scheduledOutputInventory{}, err
			}
			if _, exists := inventory.assignments[result.LaunchIndex]; exists {
				return scheduledOutputInventory{}, fmt.Errorf("scheduled resume has duplicate assignment records")
			}
			inventory.assignments[result.LaunchIndex] = result
			assignmentStems[strings.TrimSuffix(name, ".assignment.json")] = struct{}{}
		case isScheduledEvidenceName(name):
			evidence = append(evidence, name)
		default:
			return scheduledOutputInventory{}, fmt.Errorf("scheduled output directory contains unrecognized files")
		}
	}
	for _, name := range evidence {
		stem := scheduledEvidenceStem(name)
		if _, ok := assignmentStems[stem]; !ok || stem == "" {
			return scheduledOutputInventory{}, fmt.Errorf("scheduled resume has evidence without a matching assignment")
		}
	}
	if (len(inventory.assignments) > 0 || len(evidence) > 0) && !inventory.hasCheckpoint {
		return scheduledOutputInventory{}, fmt.Errorf("scheduled resume is missing a checkpoint")
	}
	return inventory, nil
}

func loadAssignmentFile(path, name string, schedule launchSchedule) (assignedTrialResult, error) {
	var result assignedTrialResult
	if err := decodeStrict(path, &result); err != nil {
		return assignedTrialResult{}, err
	}
	if result.LaunchIndex < 0 || result.LaunchIndex >= len(schedule.Entries) {
		return assignedTrialResult{}, fmt.Errorf("scheduled resume does not match the selected schedule")
	}
	entry := schedule.Entries[result.LaunchIndex]
	if name != scheduledStem(entry)+".assignment.json" || !assignmentMatchesEntry(result, entry) {
		return assignedTrialResult{}, fmt.Errorf("scheduled resume does not match the selected schedule")
	}
	if !isTerminalAssignmentStatus(result.Status) {
		return assignedTrialResult{}, fmt.Errorf("scheduled resume assignment is not terminal")
	}
	return result, nil
}

func reconstructScheduledResume(schedule launchSchedule, assignments map[int]assignedTrialResult) (scheduledResume, error) {
	resume := scheduledResume{results: []assignedTrialResult{}}
	for index, entry := range schedule.Entries {
		result, ok := assignments[index]
		if !ok {
			for later := index + 1; later < len(schedule.Entries); later++ {
				if _, exists := assignments[later]; exists {
					return scheduledResume{}, fmt.Errorf("scheduled resume has a gap in assignment records")
				}
			}
			resume.nextIndex = index
			return resume, nil
		}
		if !assignmentMatchesEntry(result, entry) {
			return scheduledResume{}, fmt.Errorf("scheduled resume does not match the selected schedule")
		}
		if result.Status == "budget_stopped" || result.Status == "safety_stopped" {
			if resume.stopKind == "" {
				resume.stopKind = result.Status
			} else if resume.stopKind != result.Status {
				return scheduledResume{}, fmt.Errorf("scheduled resume has mixed stop records")
			}
			resume.results = append(resume.results, result)
			continue
		}
		if resume.stopKind != "" {
			return scheduledResume{}, fmt.Errorf("scheduled resume has a gap in assignment records")
		}
		if !isLaunchedAssignmentStatus(result.Status) {
			return scheduledResume{}, fmt.Errorf("scheduled resume assignment is not terminal")
		}
		resume.results = append(resume.results, result)
	}
	resume.nextIndex = len(schedule.Entries)
	return resume, nil
}

func verifyResumeWorldBuild(outputDir, worldBuild string, schedule launchSchedule, resume scheduledResume) error {
	if worldBuild == "" {
		return nil
	}
	for _, result := range resume.results {
		if !isLaunchedAssignmentStatus(result.Status) {
			continue
		}
		entry := schedule.Entries[result.LaunchIndex]
		trialPath := filepath.Join(outputDir, scheduledStem(entry)+".trial.json")
		if _, err := os.Stat(trialPath); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("scheduled resume launched assignment is missing world-build evidence")
		} else if err != nil {
			return err
		}
		var recorded trial
		if err := decodeStrict(trialPath, &recorded); err != nil {
			return err
		}
		if recorded.WorldBuild != worldBuild {
			return fmt.Errorf("scheduled resume world-build does not match the live world-build")
		}
	}
	return nil
}

func verifyScheduledCheckpoint(
	path, worldBuild, scheduleDigest string,
	runBudgetUSD float64,
	groupSize int,
	resume scheduledResume,
) error {
	var checkpoint scheduledCheckpoint
	if err := decodeStrict(path, &checkpoint); err != nil {
		return err
	}
	if checkpoint.V != 1 || checkpoint.ScheduleDigest != scheduleDigest {
		return fmt.Errorf("scheduled resume checkpoint does not match assignment records")
	}
	if worldBuild != "" && checkpoint.WorldBuild != worldBuild {
		return fmt.Errorf("scheduled resume world-build does not match the live world-build")
	}
	if runBudgetUSD > 0 && checkpoint.RunBudgetUSD > 0 && checkpoint.RunBudgetUSD != runBudgetUSD {
		return fmt.Errorf("scheduled resume checkpoint does not match assignment records")
	}
	if checkpoint.NextLaunchIndex < 0 || checkpoint.NextLaunchIndex > resume.nextIndex ||
		checkpoint.NextLaunchIndex%groupSize != 0 ||
		checkpoint.CompletedPairingKeys != checkpoint.NextLaunchIndex/groupSize {
		return fmt.Errorf("scheduled resume checkpoint does not match assignment records")
	}
	if roundUSD(checkpoint.SpentUSD) != spentBefore(resume.results, checkpoint.NextLaunchIndex) {
		return fmt.Errorf("scheduled resume checkpoint does not match assignment records")
	}
	return nil
}

func writeScheduledCheckpoint(
	config runConfig,
	scheduleDigest string,
	summary scheduledSummary,
	nextLaunchIndex, completedPairingKeys int,
) error {
	return writeJSON(filepath.Join(config.outputDir, scheduledCheckpointName), scheduledCheckpoint{
		V: 1, ScheduleDigest: scheduleDigest, WorldBuild: summary.WorldBuild,
		RunBudgetUSD: summary.RunBudgetUSD, SpentUSD: summary.SpentUSD,
		NextLaunchIndex: nextLaunchIndex, CompletedPairingKeys: completedPairingKeys,
	})
}

func assignmentMatchesEntry(result assignedTrialResult, entry scheduleEntry) bool {
	return result.LaunchIndex == entry.LaunchIndex &&
		result.FamilyBlock == entry.FamilyBlock &&
		result.PairingIndex == entry.PairingIndex &&
		result.FamilyID == entry.FamilyID &&
		result.CaseID == entry.CaseID &&
		result.Repetition == entry.Repetition &&
		result.Arm == entry.Arm
}

func isLaunchedAssignmentStatus(status string) bool {
	return status == "pass" || status == "fail" || status == "unresolved"
}

func isTerminalAssignmentStatus(status string) bool {
	return isLaunchedAssignmentStatus(status) || status == "budget_stopped" || status == "safety_stopped"
}

func isScheduledEvidenceName(name string) bool {
	return strings.HasSuffix(name, ".trial.json") ||
		strings.HasSuffix(name, ".grade.json") ||
		strings.Contains(name, ".attempt-") && strings.HasSuffix(name, ".runtime.jsonl")
}

func scheduledEvidenceStem(name string) string {
	switch {
	case strings.HasSuffix(name, ".trial.json"):
		return strings.TrimSuffix(name, ".trial.json")
	case strings.HasSuffix(name, ".grade.json"):
		return strings.TrimSuffix(name, ".grade.json")
	case strings.HasSuffix(name, ".runtime.jsonl"):
		if index := strings.LastIndex(name, ".attempt-"); index >= 0 {
			return name[:index]
		}
	}
	return ""
}

func spentBefore(results []assignedTrialResult, nextIndex int) float64 {
	spent := 0.0
	for _, result := range results {
		if result.LaunchIndex < nextIndex && isLaunchedAssignmentStatus(result.Status) {
			spent = roundUSD(spent + result.TotalCostUSD)
		}
	}
	return spent
}
