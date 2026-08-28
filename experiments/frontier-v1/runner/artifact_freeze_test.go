package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type frozenFileSet struct {
	Commit                  string            `json:"candidate_commit"`
	CandidateArtifact       string            `json:"candidate_artifact"`
	CandidateArtifactDigest string            `json:"candidate_artifact_raw_sha256"`
	Report                  string            `json:"candidate_report"`
	ReportCommit            string            `json:"candidate_report_commit"`
	ReportDigest            string            `json:"candidate_report_raw_sha256"`
	Review                  string            `json:"review_record"`
	ReviewCommit            string            `json:"review_record_commit"`
	ReviewDigest            string            `json:"review_record_raw_sha256"`
	Verdict                 string            `json:"review_verdict"`
	ReplacesCommit          string            `json:"replaces_candidate_commit"`
	SetDigest               string            `json:"artifact_set_lf_normalized_utf8_sha256"`
	AcceptedFindings        []string          `json:"accepted_findings"`
	CandidateNotes          []string          `json:"candidate_notes"`
	Files                   map[string]string `json:"files"`
	Frozen                  bool              `json:"frozen"`
}

type frozenValidationSchedule struct {
	Retired              bool     `json:"retired"`
	Retirement           string   `json:"retirement"`
	Path                 string   `json:"path"`
	CanonicalDigest      string   `json:"canonical_json_sha256"`
	RawDigest            string   `json:"raw_sha256"`
	RawBytes             int64    `json:"raw_bytes"`
	Commit               string   `json:"candidate_commit"`
	Report               string   `json:"candidate_report"`
	ReportCommit         string   `json:"candidate_report_commit"`
	Review               string   `json:"review_record"`
	ReviewCommit         string   `json:"review_record_commit"`
	ReviewDigest         string   `json:"review_record_raw_sha256"`
	Verdict              string   `json:"review_verdict"`
	SourceManifest       string   `json:"source_manifest"`
	SourceManifestDigest string   `json:"source_manifest_canonical_json_sha256"`
	Seed                 uint64   `json:"seed"`
	Repetitions          int      `json:"repetitions"`
	Arms                 []string `json:"arms"`
	Cases                int      `json:"cases"`
	Families             int      `json:"families"`
	PairingKeys          int      `json:"pairing_keys"`
	Launches             int      `json:"launches"`
	Frozen               bool     `json:"frozen"`
}

type frozenClassifiedSet struct {
	Path                   string `json:"path"`
	Digest                 string `json:"lf_normalized_utf8_sha256"`
	Commit                 string `json:"candidate_commit"`
	Review                 string `json:"review_record"`
	Verdict                string `json:"review_verdict"`
	LabeledBy              string `json:"labeled_by"`
	CountersignatureRecord string `json:"countersignature_record"`
	ArchivedMessages       int    `json:"archived_messages"`
	Frozen                 bool   `json:"frozen"`
}

// verifyClassifiedSetFreeze enforces the section 5 binding: the chair-labeled,
// evaluator-countersigned classified message set is digest-pinned with both
// attributions recorded.
func verifyClassifiedSetFreeze(t *testing.T, repositoryRoot string, classified frozenClassifiedSet) {
	t.Helper()
	if !classified.Frozen || classified.Commit != "54e256e9f4347d34844a3ef9b2156600360dcd13" ||
		classified.Verdict != "ACCEPT" || classified.LabeledBy != "Raoul Bivolaru (project chair)" ||
		!strings.Contains(classified.CountersignatureRecord, "docs/reviews/2026-08-27-absence-acceptance-payload-review.md") ||
		classified.ArchivedMessages != 50 {
		t.Fatalf("classified-set freeze metadata = %#v", classified)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(classified.Review))); err != nil {
		t.Fatalf("classified-set review record: %v", err)
	}
	classifiedDigest, err := digestLFNormalizedFile(filepath.Join(repositoryRoot, filepath.FromSlash(classified.Path)))
	if err != nil || classifiedDigest != classified.Digest {
		t.Fatalf("classified-set digest = %s, freeze requires %s (%v)", classifiedDigest, classified.Digest, err)
	}
}

type frozenGateState struct {
	MayOpenValidation            bool   `json:"may_open_validation"`
	MayOpenHeldOut               bool   `json:"may_open_held_out"`
	ValidationOpenedOn           string `json:"validation_opened_on"`
	ValidationOpeningRecord      string `json:"validation_opening_decision"`
	PriorValidationOpeningRecord string `json:"prior_validation_opening_decision"`
	ValidationScope              string `json:"validation_scope"`
	ValidationExecutionStatus    string `json:"validation_execution_status"`
	ValidationExecutionClosedOn  string `json:"validation_execution_closed_on"`
	ValidationExecutionDecision  string `json:"validation_execution_decision"`
	ValidationExecutionEvidence  string `json:"validation_execution_evidence"`
	ValidationExecutionCustody   string `json:"validation_execution_custody"`
	ClosureGate                  string `json:"validation_execution_closure_gate"`
	ClosureReview                string `json:"validation_execution_closure_review"`
	ClosureReviewVerdict         string `json:"validation_execution_closure_review_verdict"`
	ClosureReviewDigest          string `json:"validation_execution_closure_review_lf_normalized_utf8_sha256"`
	ClosureFieldsReview          string `json:"validation_execution_closure_fields_payload_review"`
}

func TestPreValidationFreezeMatchesAcceptedCandidates(t *testing.T) {
	repositoryRoot := filepath.Join("..", "..", "..")
	var freeze struct {
		Status      string          `json:"status"`
		SourceLimit string          `json:"source_limit"`
		Gates       frozenGateState `json:"gates"`
		Artifacts   struct {
			ArmB struct {
				Path    string `json:"path"`
				Digest  string `json:"lf_normalized_utf8_sha256"`
				Commit  string `json:"candidate_commit"`
				Review  string `json:"review_record"`
				Verdict string `json:"review_verdict"`
				Frozen  bool   `json:"frozen"`
			} `json:"arm_b_static_document"`
			Runner        frozenFileSet            `json:"scheduled_runner"`
			Analysis      frozenFileSet            `json:"analysis_implementation_and_report_template"`
			Schedule      frozenValidationSchedule `json:"validation_schedule"`
			ClassifiedSet frozenClassifiedSet      `json:"absence_acceptance_classified_set"`
		} `json:"artifacts"`
		Remaining []string `json:"remaining"`
	}
	readJSONForTest(t, filepath.Join("..", "pre-validation-artifacts.json"), &freeze)

	verifyValidationGateState(t, repositoryRoot, freeze.Status, freeze.SourceLimit, freeze.Gates)
	armB := freeze.Artifacts.ArmB
	if !armB.Frozen || armB.Commit != "73adf8c608f0edf06597b569b17faa32e1a3b5b9" || armB.Verdict != "ACCEPT" {
		t.Fatalf("Arm B freeze metadata = %#v", armB)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(armB.Review))); err != nil {
		t.Fatalf("Arm B review record: %v", err)
	}
	actual, err := digestLFNormalizedFile(filepath.Join(repositoryRoot, filepath.FromSlash(armB.Path)))
	if err != nil {
		t.Fatal(err)
	}
	if actual != armB.Digest {
		t.Fatalf("Arm B digest = %s, freeze requires %s", actual, armB.Digest)
	}
	// The local-artifact check runs before the runner-freeze check so that a
	// payload commit, whose deliberate red is the runner set, still exercises
	// it rather than short-circuiting at the first t.Fatalf.
	verifyAcceptedLocalArtifacts(t, repositoryRoot)
	verifyReplacementRunnerFreeze(t, repositoryRoot, freeze.Artifacts.Runner)
	verifyFrozenFileSet(t, repositoryRoot, freeze.Artifacts.Analysis,
		"0f8d9c72c59bea5da5abf792f493f1d75299b1f8", 9)
	verifyFrozenValidationSchedule(t, repositoryRoot, freeze.Artifacts.Schedule)
	verifyClassifiedSetFreeze(t, repositoryRoot, freeze.Artifacts.ClassifiedSet)

	var protocol struct {
		ArtifactFreeze struct {
			BeforeValidation []string `json:"before_validation"`
		} `json:"artifact_freeze"`
	}
	readJSONForTest(t, filepath.Join("..", "protocol.json"), &protocol)
	wantRemaining := withoutStrings(protocol.ArtifactFreeze.BeforeValidation,
		"runtime invocation and exact per-arm system prompts",
		"arm A schemas",
		"arm B static document",
		"world definition and world-build digest",
		"runner digest",
		"grader digest",
		"analysis implementation and report template",
		"schedule digest",
	)
	if !reflect.DeepEqual(freeze.Remaining, wantRemaining) {
		t.Fatalf("remaining freeze set = %#v, want %#v", freeze.Remaining, wantRemaining)
	}
}

func verifyValidationGateState(t *testing.T, repositoryRoot, status, sourceLimit string, gates frozenGateState) {
	t.Helper()
	if status != "complete" || sourceLimit != "public_validation_inputs_only" ||
		gates.MayOpenValidation || gates.MayOpenHeldOut {
		t.Fatalf("pre-validation freeze gate state = %#v", gates)
	}
	expected := [][2]string{
		{gates.ValidationOpenedOn, "2026-08-26"},
		{gates.ValidationOpeningRecord, "docs/decisions/0022-protocol-v4-gate-1a-validation-reopening-after-world-compatibility-repair.md"},
		{gates.PriorValidationOpeningRecord, "docs/decisions/0019-protocol-v4-gate-1a-validation-reopening-after-refactor.md"},
		{gates.ValidationScope, "frozen_validation_schedule_only"},
		{gates.ValidationExecutionStatus, "indeterminate"},
		{gates.ValidationExecutionClosedOn, "2026-08-27"},
		{gates.ValidationExecutionDecision, "docs/decisions/0023-protocol-v4-gate-1a-second-interrupted-execution.md"},
		{gates.ValidationExecutionEvidence, "experiments/frontier-v1/results/scheduled-validation-2"},
		{gates.ValidationExecutionCustody, "experiments/frontier-v1/results/scheduled-validation-2-custody.json"},
	}
	for _, pair := range expected {
		if pair[0] != pair[1] {
			t.Fatalf("pre-validation freeze gate field = %q, want %q", pair[0], pair[1])
		}
	}
	verifyPostClosureGateState(t, gates)
	for _, record := range []string{
		gates.ValidationOpeningRecord, gates.PriorValidationOpeningRecord,
		gates.ValidationExecutionDecision, gates.ValidationExecutionCustody,
	} {
		if _, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(record))); err != nil {
			t.Fatalf("gate provenance record: %v", err)
		}
	}
}

// verifyPostClosureGateState enforces the protocol-v5 section 9 post-closure
// gate in the freeze document: the interrupted execution closed by decision
// 0023 has no independently reviewed closure record, so the fields stay empty
// and the runner refuses a further validation execution on that ground.
func verifyPostClosureGateState(t *testing.T, gates frozenGateState) {
	t.Helper()
	for _, required := range []string{
		"post-closure gate", "without a completed schedule",
		"independently reviewed and committed", "requireValidationGate",
		"pins its LF-normalized digest", "verifies identity, not prose",
		"OBLIGATION on the commit that sets these fields", "a chair gate patch must not set them",
	} {
		if !strings.Contains(gates.ClosureGate, required) {
			t.Fatalf("post-closure gate rule is missing %q: %q", required, gates.ClosureGate)
		}
	}
	if gates.ClosureReview != "" || gates.ClosureReviewVerdict != "" ||
		gates.ClosureReviewDigest != "" || gates.ClosureFieldsReview != "" {
		t.Fatalf("post-closure gate must stay unsatisfied: review=%q verdict=%q digest=%q fields-review=%q",
			gates.ClosureReview, gates.ClosureReviewVerdict, gates.ClosureReviewDigest, gates.ClosureFieldsReview)
	}
}

func verifyReplacementRunnerFreeze(t *testing.T, repositoryRoot string, artifact frozenFileSet) {
	t.Helper()
	wantFindings := []string{
		"Residual: the runner forms no opinion about whether the named closure record is an accepting independent review. It verifies that the named path is a non-empty regular file under docs/reviews, free of traversal, whose LF-normalized digest reproduces the pin, and that the gate names the independent review of the payload commit that set those fields. Whether the closure record is an accepting independent review rests on that human step, which the boundary document binds to an independently reviewed payload commit and forbids to a chair gate patch. The fields-review field is an anchor, not a check: it is not digest-pinned, carries no commit or verdict, and is not required to differ from the closure record, so it is satisfied by any non-empty regular file under docs/reviews. Its value is that the discharge is named on the record and checkable by a later auditor, not that the runner verifies it. Deleting the verdict parser gave up machine refusal of a named record whose own prose reads REVISE, or which states no verdict at all; it gave up nothing against a document that merely lists the available verdicts, which the parser accepted.",
		"Residual: neither supported host exposes the sending principal of a signal to a Go handler, so a termination event records exposed=false with the reason; establishing that no project participant initiated a kill remains the custodian written attestation, as section 9 requires.",
		"Residual: the runner verifies that the attestation claims digest identity and that its recorded schedule and world-build digests equal the live ones; byte-level re-verification of the whole frozen set at resume is the custodian step recorded in the attestation, as section 9 words condition (b).",
		"Residual: the process-event log is unsigned custodian-controlled local state, so a coordinated rewrite that deletes a start entry and re-attests the new digest defeats both the single-resume corroboration and the 72-hour anchor; the runner retains no previously accepted digest and cannot detect it. Detection is the committed attestation digest compared across commits, under independent review.",
		"Residual: the 72-hour anchor is measured against the runner own clock, which the custodian host controls; a clock set to shortly after the last process event passes every timing check, while a clock set before it is refused by the interruption anchor. No trusted time source exists on the execution host.",
		"Residual: measuring the 72 hours from the log last entry rather than from the attested interruption shortens the real window by the time from the last completed pairing-key group to the interruption, bounded by one group; the erosion is disclosed in the boundary document and is conservative in the correct direction.",
		"Residual: requirePostClosureAudit runs inside requireValidationGate after the opening-flag checks, so a closed gate reports gate is closed first; both are refusals and only the surfaced message differs.",
		"Residual: the section 9 attestation rule is scoped to the validation tranche; authoring resumes are development loops and are exempt, pinned by TestAuthoringResumeIsExemptFromTheAttestationRule. The exemption cannot launder a validation execution because validateScheduleForTranche refuses a schedule whose tranche does not match the flag. One consequence: the end-to-end resume-counter round-trip test runs on the authoring path and so does not traverse requireResumeAuthorization.",
		"Residual: nothing in the runner writes validation_execution_status. A custodian who wipes the output directory presents as a fresh run and bypasses the resume rule entirely, so the post-closure gate is the sole control on the retry channel after a future non-completed execution.",
		"Residual: P3: the partial summary committed process-event-log digest never covers the terminal stop or termination entry, because recordStop appends after the last partial-summary write; the custody record hashes the whole log at closure.",
		"Residual: P3: the partial summary carries five non-outcome identity fields beyond section 9 literal list (v, tranche, status, schedule_digest, resumes); none carries an outcome.",
		"Residual: P3: the checkpoint and partial summary are written with os.WriteFile, neither fsynced nor atomically renamed, unlike the fsynced append-only process-event log. Because the termination handler deliberately calls os.Exit, a supervised stop as well as an abrupt kill can truncate a boundary write; the failure is closed - decodeStrict refuses and the tranche closes indeterminate - but a routine Ctrl+C at an unlucky microsecond can burn the tranche. A temp-and-rename in writeJSON would remove the class.",
		"Residual: P3: a validation output directory whose first custodian process wrote its start event but died before the initial checkpoint is refused wholesale by the progress-implies-checkpoint rule and cannot be resumed; the window is milliseconds and the rule is what prevents a phantom start entry from permanently refusing a later genuine resume.",
		"Residual: Carried P2-3: applyScheduledResume returns done=true with a non-nil error on stop paths; callers must check the error first.",
		"Residual: Carried P3 (pre-existing): validateExternalGrade lets a gating manual_required dominate a gating fail, the reverse of the grader overallStatus precedence; unreachable under v5 because no check kind produces manual_required.",
		"Residual: Carried P2-1 and P2-5 remain retired by this payload; runScheduledCases and validateExternalGrade remain at gocyclo 14, both pre-existing and unchanged by this payload.",
		"Residual: neither pre-validation-artifacts.json nor artifact_freeze_test.go appears in any frozen file set, so setting the four post-closure closure fields is not a frozen-byte change and does not itself trigger the close-payload-review-refreeze cycle. The obligation that only an independently reviewed payload commit may set them is written in the frozen gate rule, the runner comment and the boundary document's Authorization boundary section, and anchored by the fields-review field; it is not machine-routed. Bringing either file into a frozen block is a change to the freeze machinery affecting every block and is out of scope for this payload.",
		"Residual: the closure record is digest-pinned and its review is not, because under the freeze workflow a payload commit precedes its own review and cannot pin bytes that do not yet exist. The companion commit, digest and verdict fields are available to the future payload that actually sets the closure fields, whose independent reviewer is the person the obligation binds; they are deferred to it rather than omitted.",
		"Seventh review N1 (P3): verifyClosureFieldsReview checks only prefix, traversal, regular-file and non-empty, so the fields-review field is satisfied by any review record including a REVISE one, an unrelated one, or the closure record itself. The project's frozenFileSet pattern carries review_record_commit, review_record_raw_sha256 and review_verdict alongside a named review and machine-checks them; the gate field carries none. Deferred to the future payload that sets the closure fields, which is by construction an independently reviewed commit.",
		"Seventh review N4 (P3): the frozen justification for the fields-review field - that the runner can refuse a gate that names none, which is what setting the fields outside a reviewed payload commit would leave behind - describes the inattentive case rather than the characteristic one, since an actor seeking to open the gate cannot leave the field empty. The field's actual value is that the discharge is named on the record and checkable by a later auditor. Not corrected here because validation_gate.go and validation-execution-boundary.md are frozen inventory files whose digests the accepting review pins; the reviewer's wording is recorded for the payload that next edits them.",
	}
	wantNotes := []string{
		"experiments/frontier-v1/artifacts/execution-resilience-candidate.md",
	}
	if artifact.CandidateArtifact != "experiments/frontier-v1/artifacts/gate-1a-execution-resilience-candidate.json" ||
		artifact.CandidateArtifactDigest != "sha256:e87ea77d058f9100252e1391f96de8398f2189aaab213063f5f84c87b431ac3c" ||
		artifact.Review != "docs/reviews/2026-08-28-execution-resilience-payload-review-7.md" ||
		artifact.ReviewCommit != "fce272f041e4b20557b2f8ffff27690f42db68de" ||
		artifact.ReviewDigest != "sha256:b6df59842575c424f83171dc7ef78694afd436da12f3b04ac36c3671006e80b6" ||
		artifact.ReplacesCommit != "54e256e9f4347d34844a3ef9b2156600360dcd13" ||
		!reflect.DeepEqual(artifact.AcceptedFindings, wantFindings) ||
		!reflect.DeepEqual(artifact.CandidateNotes, wantNotes) {
		t.Fatalf("replacement runner provenance = %#v", artifact)
	}
	verifyRawFileDigest(t, filepath.Join(repositoryRoot, filepath.FromSlash(artifact.CandidateArtifact)), artifact.CandidateArtifactDigest)
	verifyRawGitFileDigest(t, repositoryRoot, artifact.ReviewCommit, artifact.Review, artifact.ReviewDigest)
	verifyRawFileDigest(t, filepath.Join(repositoryRoot, filepath.FromSlash(artifact.Review)), artifact.ReviewDigest)
	for _, note := range artifact.CandidateNotes {
		if _, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(note))); err != nil {
			t.Fatalf("candidate note %s: %v", note, err)
		}
	}
	verifyFrozenFileSet(t, repositoryRoot, artifact,
		"27edb8664a060d0a7690039e8e29dba94facada0", 24)
}

func verifyRawFileDigest(t *testing.T, path string, want string) {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if actual := fmt.Sprintf("sha256:%x", sha256.Sum256(contents)); actual != want {
		t.Fatalf("%s raw digest = %s, freeze requires %s", path, actual, want)
	}
}

func verifyRawGitFileDigest(t *testing.T, repositoryRoot, commit, path, want string) {
	t.Helper()
	contents, err := exec.Command("git", "-C", repositoryRoot, "show", commit+":"+path).Output()
	if err != nil {
		t.Fatalf("read %s at %s: %v", path, commit, err)
	}
	if actual := fmt.Sprintf("sha256:%x", sha256.Sum256(contents)); actual != want {
		t.Fatalf("%s at %s raw digest = %s, freeze requires %s", path, commit, actual, want)
	}
}

func verifyFrozenValidationSchedule(t *testing.T, repositoryRoot string, artifact frozenValidationSchedule) {
	t.Helper()
	want := frozenValidationSchedule{
		Retired:              true,
		Retirement:           "Retired unopened under protocol v5, which removes arm E: the A-D schedule construction refuses this five-arm schedule and its manifest, registry, and archive identities. The bytes below are preserved as custody evidence and must not be run, relabeled, or resealed. A new disjoint sealed v5 tranche and A-D schedule must be independently generated and frozen before any validation execution.",
		Path:                 "experiments/frontier-v1/schedules/validation.json",
		CanonicalDigest:      "sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813",
		RawDigest:            "sha256:6b264a8daff60d9b507f11d760f1bbf19dec556f2aebd700dc5fbdc8a8c56aec",
		RawBytes:             391950,
		Commit:               "44fccf51cb1684a1b71da9128cab30ad2b2fb6af",
		Report:               "docs/reviews/2026-08-21-protocol-v4-validation-schedule-candidate.md",
		ReportCommit:         "bb473e22b1e2a5b21221a0b7b57e2aec4836c52b",
		Review:               "docs/reviews/2026-08-21-protocol-v4-validation-schedule-review.md",
		ReviewCommit:         "a301bdc781bd4b93f12418b577b7446d01c1b7e4",
		ReviewDigest:         "sha256:cd557f3dc79b7a2bbc2acbba0c40969b832178b1c23b4ca4c4442d602c601ca0",
		Verdict:              "ACCEPT",
		SourceManifest:       "experiments/frontier-v1/manifests/validation.json",
		SourceManifestDigest: "sha256:39acbad5e45ad65302659cd0875bdfe589165ede9b60b6448ac4b09ccfb1e0c6",
		Seed:                 20260817,
		Repetitions:          3,
		Arms:                 []string{"A", "B", "C", "D", "E"},
		Cases:                120,
		Families:             24,
		PairingKeys:          360,
		Launches:             1800,
		Frozen:               true,
	}
	if !reflect.DeepEqual(artifact, want) {
		t.Fatalf("validation schedule freeze metadata = %#v", artifact)
	}
	verifyValidationScheduleFiles(t, repositoryRoot, artifact)
}

func verifyValidationScheduleFiles(t *testing.T, repositoryRoot string, artifact frozenValidationSchedule) {
	t.Helper()
	for _, path := range []string{artifact.Path, artifact.Report, artifact.Review, artifact.SourceManifest} {
		if _, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(path))); err != nil {
			t.Fatalf("validation schedule artifact %s: %v", path, err)
		}
	}
	schedulePath := filepath.Join(repositoryRoot, filepath.FromSlash(artifact.Path))
	contents, err := os.ReadFile(schedulePath)
	if err != nil {
		t.Fatal(err)
	}
	if actual := fmt.Sprintf("sha256:%x", sha256.Sum256(contents)); actual != artifact.RawDigest || int64(len(contents)) != artifact.RawBytes {
		t.Fatalf("validation schedule raw identity = %s/%d", actual, len(contents))
	}
	if actual, err := digestJSONFile(schedulePath); err != nil || actual != artifact.CanonicalDigest {
		t.Fatalf("validation schedule canonical digest = %s, %v", actual, err)
	}
	if actual, err := digestJSONFile(filepath.Join(repositoryRoot, filepath.FromSlash(artifact.SourceManifest))); err != nil || actual != artifact.SourceManifestDigest {
		t.Fatalf("validation manifest canonical digest = %s, %v", actual, err)
	}
	reviewContents, err := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(artifact.Review)))
	if err != nil {
		t.Fatal(err)
	}
	if actual := fmt.Sprintf("sha256:%x", sha256.Sum256(reviewContents)); actual != artifact.ReviewDigest {
		t.Fatalf("validation schedule review digest = %s", actual)
	}
}

func verifyAcceptedLocalArtifacts(t *testing.T, repositoryRoot string) {
	t.Helper()
	const (
		candidateCommit = "54e256e9f4347d34844a3ef9b2156600360dcd13"
		candidatePath   = "experiments/frontier-v1/artifacts/pre-validation-local-candidate.json"
		candidateDigest = "sha256:a70ce35b7cfbb686cf7922c99394f6c3ef344235b03f8cf3f88bccb8ce892967"
		reviewPath      = "docs/reviews/2026-08-27-absence-acceptance-payload-review.md"
	)
	actualDigest, err := digestLFNormalizedFile(filepath.Join(repositoryRoot, filepath.FromSlash(candidatePath)))
	if err != nil || actualDigest != candidateDigest {
		t.Fatalf("accepted local candidate digest = %s, %v", actualDigest, err)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(reviewPath))); err != nil {
		t.Fatalf("accepted local review record: %v", err)
	}

	var freezeDocument map[string]any
	readJSONForTest(t, filepath.Join("..", "pre-validation-artifacts.json"), &freezeDocument)
	var candidateDocument map[string]any
	readJSONForTest(t, filepath.Join("..", "artifacts", "pre-validation-local-candidate.json"), &candidateDocument)
	freezeArtifacts := freezeDocument["artifacts"].(map[string]any)
	candidateArtifacts := candidateDocument["artifacts"].(map[string]any)
	for _, key := range []string{
		"runtime_invocation_and_exact_per_arm_system_prompts",
		"arm_a_schemas",
		"world_definition_and_world_build_digest",
		"grader_digest",
	} {
		accepted := freezeArtifacts[key].(map[string]any)
		if accepted["candidate_commit"] != candidateCommit || accepted["candidate_artifact"] != candidatePath ||
			accepted["candidate_artifact_lf_normalized_utf8_sha256"] != candidateDigest ||
			accepted["review_record"] != reviewPath || accepted["review_verdict"] != "ACCEPT" || accepted["frozen"] != true {
			t.Fatalf("accepted local artifact %s metadata = %#v", key, accepted)
		}
		copied := make(map[string]any, len(accepted)-4)
		for field, value := range accepted {
			switch field {
			case "candidate_commit", "candidate_artifact", "candidate_artifact_lf_normalized_utf8_sha256", "review_record", "review_verdict":
				continue
			default:
				copied[field] = value
			}
		}
		copied["frozen"] = false
		if !reflect.DeepEqual(copied, candidateArtifacts[key]) {
			t.Fatalf("accepted local artifact %s differs from reviewed candidate", key)
		}
	}
}

func verifyFrozenFileSet(t *testing.T, repositoryRoot string, artifact frozenFileSet, commit string, fileCount int) {
	t.Helper()
	if !artifact.Frozen || artifact.Commit != commit || artifact.Verdict != "ACCEPT" || len(artifact.Files) != fileCount {
		t.Fatalf("frozen file set metadata = %#v", artifact)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(artifact.Review))); err != nil {
		t.Fatalf("review record: %v", err)
	}
	paths := make([]string, 0, len(artifact.Files))
	for path := range artifact.Files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	var identity strings.Builder
	for _, path := range paths {
		actual, err := digestLFNormalizedFile(filepath.Join(repositoryRoot, filepath.FromSlash(path)))
		if err != nil {
			t.Fatal(err)
		}
		if actual != artifact.Files[path] {
			t.Fatalf("%s digest = %s, freeze requires %s", path, actual, artifact.Files[path])
		}
		fmt.Fprintf(&identity, "%s\t%s\n", path, actual)
	}
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(identity.String())))
	if digest != artifact.SetDigest {
		t.Fatalf("artifact set digest = %s, freeze requires %s", digest, artifact.SetDigest)
	}
}

func readJSONForTest(t *testing.T, path string, target any) {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(contents, target); err != nil {
		t.Fatal(err)
	}
}

func withoutStrings(values []string, omitted ...string) []string {
	omissions := make(map[string]struct{}, len(omitted))
	for _, value := range omitted {
		omissions[value] = struct{}{}
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, skip := omissions[value]; !skip {
			result = append(result, value)
		}
	}
	return result
}
