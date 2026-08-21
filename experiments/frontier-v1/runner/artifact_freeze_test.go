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
	Files                   map[string]string `json:"files"`
	Frozen                  bool              `json:"frozen"`
}

type frozenValidationSchedule struct {
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

func TestPreValidationFreezeMatchesAcceptedCandidates(t *testing.T) {
	repositoryRoot := filepath.Join("..", "..", "..")
	var freeze struct {
		Status      string `json:"status"`
		SourceLimit string `json:"source_limit"`
		Gates       struct {
			MayOpenValidation bool `json:"may_open_validation"`
			MayOpenHeldOut    bool `json:"may_open_held_out"`
		} `json:"gates"`
		Artifacts struct {
			ArmB struct {
				Path    string `json:"path"`
				Digest  string `json:"lf_normalized_utf8_sha256"`
				Commit  string `json:"candidate_commit"`
				Review  string `json:"review_record"`
				Verdict string `json:"review_verdict"`
				Frozen  bool   `json:"frozen"`
			} `json:"arm_b_static_document"`
			Runner   frozenFileSet            `json:"scheduled_runner"`
			Analysis frozenFileSet            `json:"analysis_implementation_and_report_template"`
			Schedule frozenValidationSchedule `json:"validation_schedule"`
		} `json:"artifacts"`
		Remaining []string `json:"remaining"`
	}
	readJSONForTest(t, filepath.Join("..", "pre-validation-artifacts.json"), &freeze)

	if freeze.Status != "complete" || freeze.SourceLimit != "public_validation_inputs_only" ||
		freeze.Gates.MayOpenValidation || freeze.Gates.MayOpenHeldOut {
		t.Fatalf("pre-validation freeze has open or incomplete gate state: %#v", freeze)
	}
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
	verifyReplacementRunnerFreeze(t, repositoryRoot, freeze.Artifacts.Runner)
	verifyFrozenFileSet(t, repositoryRoot, freeze.Artifacts.Analysis,
		"62946f4a1a03ea89636c5b3243f3b4d53b166682", 9)
	verifyAcceptedLocalArtifacts(t, repositoryRoot)
	verifyFrozenValidationSchedule(t, repositoryRoot, freeze.Artifacts.Schedule)

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

func verifyReplacementRunnerFreeze(t *testing.T, repositoryRoot string, artifact frozenFileSet) {
	t.Helper()
	wantFindings := []string{
		"P2-1: The candidate report overstates the marker and output assertions in its nil-ID regression test; its error assertion and guard ordering remain load-bearing, and the independent populated probe established no custodian contact.",
		"P2-2: Numeric-equivalent cap spellings are forwarded verbatim to the external runtime parser; the frozen operator spelling 0.15 is unaffected, and parser rejection remains a loud, auditable failure.",
	}
	if artifact.CandidateArtifact != "experiments/frontier-v1/artifacts/validation-execution-boundary-candidate.json" ||
		artifact.CandidateArtifactDigest != "sha256:14b0761362709780ded9f6fc47e7ff8d49f688e09062d8dcc6d13724b9109d1d" ||
		artifact.Report != "docs/reviews/2026-08-21-protocol-v4-validation-execution-boundary-second-revision-candidate.md" ||
		artifact.ReportCommit != "05ab169acdfabae4ce7b42b4ffb8d6ecb276e935" ||
		artifact.ReportDigest != "sha256:c77207aaebb6876805354e720adc7e72301ab759e2affaa3df3bb7c88526c546" ||
		artifact.Review != "docs/reviews/2026-08-21-protocol-v4-validation-execution-boundary-second-revision-review.md" ||
		artifact.ReviewCommit != "6b158ea0de163b43f8c0fdeb1e8cc3fd40609a52" ||
		artifact.ReviewDigest != "sha256:e1444fcca48e7d5b13e1b0e4cd385898d387d8aba4799fae0e3d49f27326ddaf" ||
		artifact.ReplacesCommit != "b4df919070bb9a6d2912662b4a59674b0e25a332" ||
		!reflect.DeepEqual(artifact.AcceptedFindings, wantFindings) {
		t.Fatalf("replacement runner provenance = %#v", artifact)
	}
	verifyRawFileDigest(t, filepath.Join(repositoryRoot, filepath.FromSlash(artifact.CandidateArtifact)), artifact.CandidateArtifactDigest)
	verifyRawGitFileDigest(t, repositoryRoot, artifact.ReportCommit, artifact.Report, artifact.ReportDigest)
	verifyRawFileDigest(t, filepath.Join(repositoryRoot, filepath.FromSlash(artifact.Review)), artifact.ReviewDigest)
	verifyFrozenFileSetAtCommit(t, repositoryRoot, artifact,
		"74d06da13f62cffa4fd635e048e6331a6e2d95a6", 16)
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

func verifyFrozenFileSetAtCommit(t *testing.T, repositoryRoot string, artifact frozenFileSet, commit string, fileCount int) {
	t.Helper()
	if !artifact.Frozen || artifact.Commit != commit || artifact.Verdict != "ACCEPT" || len(artifact.Files) != fileCount {
		t.Fatalf("historical frozen file set metadata = %#v", artifact)
	}
	paths := make([]string, 0, len(artifact.Files))
	for path := range artifact.Files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	var identity strings.Builder
	for _, path := range paths {
		contents, err := exec.Command("git", "-C", repositoryRoot, "show", commit+":"+path).Output()
		if err != nil {
			t.Fatalf("read %s at %s: %v", path, commit, err)
		}
		normalized := strings.ReplaceAll(strings.ReplaceAll(string(contents), "\r\n", "\n"), "\r", "\n")
		actual := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(normalized)))
		if actual != artifact.Files[path] {
			t.Fatalf("%s at %s digest = %s, freeze requires %s", path, commit, actual, artifact.Files[path])
		}
		fmt.Fprintf(&identity, "%s\t%s\n", path, actual)
	}
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(identity.String())))
	if digest != artifact.SetDigest {
		t.Fatalf("historical artifact set digest = %s, freeze requires %s", digest, artifact.SetDigest)
	}
}

func verifyFrozenValidationSchedule(t *testing.T, repositoryRoot string, artifact frozenValidationSchedule) {
	t.Helper()
	want := frozenValidationSchedule{
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
		SourceManifestDigest: "sha256:57ccc0c754f7c2beb74063dacc4bad26ab8b598ac2da9b6a8d637afd391fd490",
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
		candidateCommit = "795ce71acd51beb977190811c90cc538f3c6b928"
		candidatePath   = "experiments/frontier-v1/artifacts/pre-validation-local-candidate.json"
		candidateDigest = "sha256:19743d85232a2ec83b80bb71fc6b5dc3ee6cc669d319c446218aa248eb8f76b9"
		reviewPath      = "docs/reviews/2026-08-19-pre-validation-local-freeze-review-2.md"
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
