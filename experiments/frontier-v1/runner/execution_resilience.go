package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const scheduledPartialSummaryName = "scheduled-summary.partial.json"

// maximumScheduledResumes encodes the protocol-v5 section 9 rule that at most
// one resume is permitted per tranche; a second interruption closes the
// tranche indeterminate.
const maximumScheduledResumes = 1

// resumeWindow is the section 9 backstop between interruption and resume. It
// is a ceiling, not an allowance: the promptness duty requires the attempt as
// soon as the no-inspection, digest-identity, and cause conditions verify.
const resumeWindow = 72 * time.Hour

// frozenResumeCauses is the frozen outcome-uncorrelated cause list. OOM is
// struck deliberately: transcript-heavy failing trials make memory pressure
// outcome-correlated. process_kill is admissible only when the custodian's
// process-event log establishes that no project participant, account, or
// agent initiated the kill.
var frozenResumeCauses = map[string]bool{
	"host_restart":     true,
	"power_loss":       true,
	"hardware_failure": true,
	"process_kill":     true,
}

type partialArtifactDigests struct {
	WorldBuild         string `json:"world_build"`
	GraderDigest       string `json:"grader_digest,omitempty"`
	ArmBDocumentDigest string `json:"arm_b_document_digest,omitempty"`
	ProcessEventLog    string `json:"process_event_log_sha256"`
}

// scheduledPartialSummary is the outcome-free progress record written at every
// pairing-key checkpoint. It carries launch index, the completed pairing-key
// list, per-arm assignment counts, artifact digests, timestamps, and
// cumulative spend, and nothing else: no success counts, no grade tallies, and
// no per-arm outcome field, so reading it cannot steer a resume decision.
type scheduledPartialSummary struct {
	V                    int                    `json:"v"`
	Tranche              string                 `json:"tranche"`
	Status               string                 `json:"status"`
	ScheduleDigest       string                 `json:"schedule_digest"`
	ArtifactDigests      partialArtifactDigests `json:"artifact_digests"`
	StartedAt            string                 `json:"started_at"`
	UpdatedAt            string                 `json:"updated_at"`
	NextLaunchIndex      int                    `json:"next_launch_index"`
	CompletedPairingKeys []string               `json:"completed_pairing_keys"`
	ArmAssignmentCounts  map[string]int         `json:"arm_assignment_counts"`
	SpentUSD             float64                `json:"spent_usd"`
	Resumes              int                    `json:"resumes"`
}

func scheduledPairingKey(entry scheduleEntry) string {
	return fmt.Sprintf("%s/r%02d", entry.CaseID, entry.Repetition+1)
}

// writeScheduledPartialSummary rewrites the partial summary for the launches
// retained before nextLaunchIndex. The first write stamps started_at; later
// writes preserve it.
func writeScheduledPartialSummary(
	config runConfig,
	schedule launchSchedule,
	scheduleDigest string,
	spentUSD float64,
	nextLaunchIndex, resumes int,
) error {
	path := filepath.Join(config.outputDir, scheduledPartialSummaryName)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	started := now
	var existing scheduledPartialSummary
	if err := decodeStrict(path, &existing); err == nil && existing.StartedAt != "" {
		started = existing.StartedAt
	}
	logDigest, err := digestProcessEventLog(config.outputDir)
	if err != nil {
		return err
	}
	keys := make([]string, 0, nextLaunchIndex)
	counts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range schedule.Entries[:nextLaunchIndex] {
		counts[entry.Arm]++
		if key := scheduledPairingKey(entry); !seen[key] {
			seen[key] = true
			keys = append(keys, key)
		}
	}
	return writeJSON(path, scheduledPartialSummary{
		V: 1, Tranche: schedule.Tranche, Status: "in_progress", ScheduleDigest: scheduleDigest,
		ArtifactDigests: partialArtifactDigests{
			WorldBuild: config.worldBuild, GraderDigest: config.graderDigest,
			ArmBDocumentDigest: config.armBDocumentDigest, ProcessEventLog: logDigest,
		},
		StartedAt: started, UpdatedAt: now, NextLaunchIndex: nextLaunchIndex,
		CompletedPairingKeys: keys, ArmAssignmentCounts: counts,
		SpentUSD: roundUSD(spentUSD), Resumes: resumes,
	})
}

// resumeAttestation is the custodian's written record required before any
// resume of an interrupted validation tranche. Every field is a precommitted
// section 9 condition; the runner refuses the resume when one fails, which
// closes the tranche indeterminate rather than abandoning it by choice.
type resumeAttestation struct {
	V                      int               `json:"v"`
	Tranche                string            `json:"tranche"`
	OutputDir              string            `json:"output_dir"`
	Custodian              string            `json:"custodian"`
	Cause                  string            `json:"cause"`
	ResumeIndex            int               `json:"resume_index"`
	InterruptedAt          string            `json:"interrupted_at"`
	CauseClassifiedAt      string            `json:"cause_classified_at"`
	ConditionsVerifiedAt   string            `json:"conditions_verified_at"`
	ResumeAuthorizedAt     string            `json:"resume_authorized_at"`
	NoOutcomeInspection    bool              `json:"no_outcome_inspection_between_interruption_and_resume"`
	ClassifiedBeforeAny    bool              `json:"cause_classified_before_inspection_and_resume_decision"`
	FrozenBytesIdentical   bool              `json:"frozen_bytes_digest_identical_at_resume"`
	KillPrincipal          bool              `json:"kill_principal_established_as_non_participant"`
	ProcessEventLog        string            `json:"process_event_log"`
	ProcessEventLogSHA256  string            `json:"process_event_log_sha256"`
	PromptnessStatement    string            `json:"promptness_statement"`
	LapseExplanation       string            `json:"lapse_explanation,omitempty"`
	FrozenDigestsVerified  map[string]string `json:"frozen_digests_verified"`
	IndependentReviewOfLog string            `json:"process_event_log_review,omitempty"`
}

type resumeAuthorization struct {
	tranche         string
	outputRelative  string
	outputDir       string
	attestationPath string
	scheduleDigest  string
	worldBuild      string
}

// requireResumeAuthorization enforces the section 9 mandatory-resume rule for
// a validation tranche whose output directory already carries a checkpoint
// written by an earlier custodian process. Authoring resumes are development
// loops and are exempt; validation resumes are not.
func requireResumeAuthorization(control resumeAuthorization, resume scheduledResume) error {
	if !resume.hasCheckpoint || control.tranche != "validation" {
		return nil
	}
	if resume.checkpoint.Resumes >= maximumScheduledResumes {
		return fmt.Errorf("validation tranche already used its single resume; a second interruption closes it indeterminate")
	}
	if control.attestationPath == "" {
		return fmt.Errorf("resuming an interrupted validation tranche requires --resume-attestation")
	}
	var attestation resumeAttestation
	if err := decodeStrict(control.attestationPath, &attestation); err != nil {
		return err
	}
	if err := verifyResumeAttestationFields(control, attestation, resume.checkpoint.Resumes+1); err != nil {
		return err
	}
	return verifyResumeAttestationTiming(attestation)
}

func verifyResumeAttestationFields(control resumeAuthorization, attestation resumeAttestation, wantIndex int) error {
	if attestation.V != 1 || attestation.Tranche != "validation" || attestation.OutputDir != control.outputRelative {
		return fmt.Errorf("resume attestation does not describe this validation output directory")
	}
	if attestation.Custodian == "" || attestation.PromptnessStatement == "" {
		return fmt.Errorf("resume attestation requires a custodian and a promptness statement")
	}
	if attestation.ResumeIndex != wantIndex {
		return fmt.Errorf("resume attestation records resume %d, this is resume %d", attestation.ResumeIndex, wantIndex)
	}
	if !frozenResumeCauses[attestation.Cause] {
		return fmt.Errorf("interruption cause %q is not on the frozen outcome-uncorrelated cause list", attestation.Cause)
	}
	if !attestation.NoOutcomeInspection || !attestation.ClassifiedBeforeAny || !attestation.FrozenBytesIdentical {
		return fmt.Errorf("resume attestation does not establish no-inspection, pre-decision cause classification, and digest identity")
	}
	if attestation.Cause == "process_kill" && !attestation.KillPrincipal {
		return fmt.Errorf("a process kill resumes only when the process-event log establishes a non-participant principal")
	}
	if err := verifyResumeFrozenDigests(control, attestation); err != nil {
		return err
	}
	return verifyAttestedProcessEventLog(control.outputDir, attestation)
}

func verifyResumeFrozenDigests(control resumeAuthorization, attestation resumeAttestation) error {
	for field, want := range map[string]string{
		"schedule_digest": control.scheduleDigest,
		"world_build":     control.worldBuild,
	} {
		if want == "" {
			continue
		}
		if attestation.FrozenDigestsVerified[field] != want {
			return fmt.Errorf("resume attestation records %s %q, live %q", field, attestation.FrozenDigestsVerified[field], want)
		}
	}
	return nil
}

func verifyAttestedProcessEventLog(outputDir string, attestation resumeAttestation) error {
	if attestation.ProcessEventLog != scheduledProcessEventLogName {
		return fmt.Errorf("resume attestation must cite the %s process-event log", scheduledProcessEventLogName)
	}
	digest, err := digestProcessEventLog(outputDir)
	if err != nil {
		return err
	}
	if digest != attestation.ProcessEventLogSHA256 {
		return fmt.Errorf("process-event log digest %s does not match the attested %s", digest, attestation.ProcessEventLogSHA256)
	}
	return nil
}

func verifyResumeAttestationTiming(attestation resumeAttestation) error {
	stamps := map[string]time.Time{}
	for field, value := range map[string]string{
		"interrupted_at":         attestation.InterruptedAt,
		"cause_classified_at":    attestation.CauseClassifiedAt,
		"conditions_verified_at": attestation.ConditionsVerifiedAt,
		"resume_authorized_at":   attestation.ResumeAuthorizedAt,
	} {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return fmt.Errorf("resume attestation %s is not an RFC3339 timestamp", field)
		}
		stamps[field] = parsed
	}
	ordered := []string{"interrupted_at", "cause_classified_at", "conditions_verified_at", "resume_authorized_at"}
	for index := 1; index < len(ordered); index++ {
		if stamps[ordered[index]].Before(stamps[ordered[index-1]]) {
			return fmt.Errorf("resume attestation timestamps are out of order at %s", ordered[index])
		}
	}
	elapsed := stamps["resume_authorized_at"].Sub(stamps["interrupted_at"])
	if elapsed > resumeWindow {
		return fmt.Errorf("resume authorized %s after the interruption, past the 72-hour backstop", elapsed)
	}
	if stamps["resume_authorized_at"].After(stamps["conditions_verified_at"]) && attestation.LapseExplanation == "" {
		return fmt.Errorf("a resume later than the moment its conditions verified requires a written lapse explanation")
	}
	return nil
}

// sortedPartialPairingKeys is used by tests and the closure tooling to compare
// completed pairing-key lists independent of schedule order.
func sortedPartialPairingKeys(summary scheduledPartialSummary) []string {
	keys := append([]string(nil), summary.CompletedPairingKeys...)
	sort.Strings(keys)
	return keys
}

// readScheduledPartialSummary loads the partial summary from a scheduled
// output directory.
func readScheduledPartialSummary(outputDir string) (scheduledPartialSummary, error) {
	var summary scheduledPartialSummary
	path := filepath.Join(outputDir, scheduledPartialSummaryName)
	if _, err := os.Stat(path); err != nil {
		return scheduledPartialSummary{}, err
	}
	if err := decodeStrict(path, &summary); err != nil {
		return scheduledPartialSummary{}, err
	}
	return summary, nil
}

// scheduledProgress bundles the section 9 progress records written at every
// pairing-key boundary: the durable resume checkpoint, the outcome-free
// partial summary, and the append-only process-event log.
type scheduledProgress struct {
	log       *processEventLog
	resumes   int
	groupSize int
}

// recordStop appends the terminal process event for this custodian process.
// The detail carries the scheduled run status and stop reason only, the same
// non-outcome facts the completion criterion already reports.
func (progress scheduledProgress) recordStop(nextLaunchIndex int, summary scheduledSummary) error {
	completed := nextLaunchIndex
	if progress.groupSize > 0 {
		completed = nextLaunchIndex / progress.groupSize
	}
	detail := summary.Status
	if summary.StopReason != "" {
		detail += ": " + summary.StopReason
	}
	return progress.log.record(processEventStop, progress.events(nextLaunchIndex, completed, summary.SpentUSD), detail)
}

// scheduledResumeControl describes this execution to the section 9 resume
// check: which tranche, which output directory, and the live frozen digests a
// custodian attestation must have verified.
func scheduledResumeControl(config runConfig, scheduleDigest string) (resumeAuthorization, error) {
	relative, err := filepath.Rel(config.repositoryRoot, config.outputDir)
	if err != nil {
		return resumeAuthorization{}, err
	}
	return resumeAuthorization{
		tranche:         effectiveTranche(config),
		outputRelative:  filepath.ToSlash(relative),
		outputDir:       config.outputDir,
		attestationPath: config.resumeAttestation,
		scheduleDigest:  scheduleDigest,
		worldBuild:      config.worldBuild,
	}, nil
}

func (progress scheduledProgress) events(nextLaunchIndex, completedPairingKeys int, spentUSD float64) processEventProgress {
	return processEventProgress{
		nextLaunchIndex: nextLaunchIndex, completedPairingKeys: completedPairingKeys,
		spentUSD: spentUSD, resumes: progress.resumes,
	}
}

func (progress scheduledProgress) recordCheckpoint(
	config runConfig,
	schedule launchSchedule,
	scheduleDigest string,
	summary scheduledSummary,
	nextLaunchIndex, completedPairingKeys int,
) error {
	if err := writeScheduledCheckpoint(config, scheduleDigest, summary, nextLaunchIndex, completedPairingKeys, progress.resumes); err != nil {
		return err
	}
	if err := progress.log.record(processEventCheckpoint,
		progress.events(nextLaunchIndex, completedPairingKeys, summary.SpentUSD), ""); err != nil {
		return err
	}
	return writeScheduledPartialSummary(config, schedule, scheduleDigest, summary.SpentUSD, nextLaunchIndex, progress.resumes)
}
