package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type frozenFileSet struct {
	Commit    string            `json:"candidate_commit"`
	Review    string            `json:"review_record"`
	Verdict   string            `json:"review_verdict"`
	SetDigest string            `json:"artifact_set_lf_normalized_utf8_sha256"`
	Files     map[string]string `json:"files"`
	Frozen    bool              `json:"frozen"`
}

func TestPreValidationFreezeMatchesAcceptedCandidates(t *testing.T) {
	repositoryRoot := filepath.Join("..", "..", "..")
	var freeze struct {
		Status string `json:"status"`
		Gates  struct {
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
			Runner   frozenFileSet `json:"scheduled_runner"`
			Analysis frozenFileSet `json:"analysis_implementation_and_report_template"`
		} `json:"artifacts"`
		Remaining []string `json:"remaining"`
	}
	readJSONForTest(t, filepath.Join("..", "pre-validation-artifacts.json"), &freeze)

	if freeze.Status != "partial" || freeze.Gates.MayOpenValidation || freeze.Gates.MayOpenHeldOut {
		t.Fatalf("pre-validation freeze has open or non-partial gate state: %#v", freeze)
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
	verifyFrozenFileSet(t, repositoryRoot, freeze.Artifacts.Runner,
		"b4df919070bb9a6d2912662b4a59674b0e25a332", 9)
	verifyFrozenFileSet(t, repositoryRoot, freeze.Artifacts.Analysis,
		"62946f4a1a03ea89636c5b3243f3b4d53b166682", 9)
	verifyAcceptedLocalArtifacts(t, repositoryRoot)

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
	)
	if !reflect.DeepEqual(freeze.Remaining, wantRemaining) {
		t.Fatalf("remaining freeze set = %#v, want %#v", freeze.Remaining, wantRemaining)
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
