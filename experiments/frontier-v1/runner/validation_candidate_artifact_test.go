package main

import (
	"crypto/sha256"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

type reviewCandidateDocument struct {
	V           int    `json:"v"`
	Status      string `json:"status"`
	SourceLimit string `json:"source_limit"`
	BaseCommit  string `json:"base_commit"`
	Frozen      bool   `json:"frozen"`
	Gates       struct {
		MayOpenValidation bool `json:"may_open_validation"`
		MayOpenHeldOut    bool `json:"may_open_held_out"`
	} `json:"gates"`
	SetDigest          string            `json:"artifact_set_lf_normalized_utf8_sha256"`
	Files              map[string]string `json:"files"`
	ExecutionPerformed bool              `json:"execution_performed"`
}

func loadReviewCandidate(t *testing.T, name, baseCommit string, fileCount int) reviewCandidateDocument {
	t.Helper()
	var candidate reviewCandidateDocument
	readJSONForTest(t, filepath.Join("..", "artifacts", name), &candidate)
	if candidate.V != 1 || candidate.Status != "review_candidate" || candidate.SourceLimit != "public_validation_inputs_only" ||
		candidate.BaseCommit != baseCommit || candidate.Frozen || candidate.Gates.MayOpenValidation ||
		candidate.Gates.MayOpenHeldOut || candidate.ExecutionPerformed {
		t.Fatalf("candidate %s state = %#v", name, candidate)
	}
	if len(candidate.Files) != fileCount {
		t.Fatalf("candidate %s file count = %d", name, len(candidate.Files))
	}
	return candidate
}

func TestGate1ATestTransitionCandidateMatchesHistoricalPayload(t *testing.T) {
	candidate := loadReviewCandidate(t, "gate-1a-test-transition-candidate.json",
		"94761cf43b938f643173b486c7f1fca5149e8a08", 16)
	root := filepath.Join("..", "..", "..")
	verifyCandidateSetAtCommit(t, root, "71f9648789decf4cd56ef8a24bc840b0dda7efd9", candidate.Files, candidate.SetDigest)
}

func TestGate1AAuthoringRepairCandidateMatchesHistoricalPayload(t *testing.T) {
	candidate := loadReviewCandidate(t, "gate-1a-authoring-repair-candidate.json",
		"0ffb8f9b44edca25582ec35541d25257d6a36dd0", 19)
	root := filepath.Join("..", "..", "..")
	verifyCandidateSetAtCommit(t, root, "610fa588488579cf5551a627795d4dc7f771f053", candidate.Files, candidate.SetDigest)
}

func TestGate1AComplexityRefactorCandidateMatchesHistoricalPayload(t *testing.T) {
	candidate := loadReviewCandidate(t, "gate-1a-complexity-refactor-candidate.json",
		"8881d5f552aedfd9e6283ee1cc33e468fe036856", 20)
	root := filepath.Join("..", "..", "..")
	verifyCandidateSetAtCommit(t, root, "8ed9202c80d8f591c5d0db8a7e8a349022f952f8", candidate.Files, candidate.SetDigest)
}

func TestGate1AValidationWorldCompatibilityCandidateMatchesPayload(t *testing.T) {
	candidate := loadReviewCandidate(t, "gate-1a-validation-world-compatibility-candidate.json",
		"acc9730dfa553e56dc20af992176717a9c4f6044", 21)
	root := filepath.Join("..", "..", "..")
	verifyCandidateSetAtCommit(t, root, "215b30ae89933c532468b452be6238a6028a740e", candidate.Files, candidate.SetDigest)
}

func verifyCandidateSetInWorkingTree(t *testing.T, repositoryRoot string, files map[string]string, wantSetDigest string) {
	t.Helper()
	var identity strings.Builder
	for _, path := range sortedCandidatePaths(files) {
		actual, err := digestLFNormalizedFile(filepath.Join(repositoryRoot, filepath.FromSlash(path)))
		if err != nil {
			t.Fatal(err)
		}
		if actual != files[path] {
			t.Fatalf("candidate file %s digest = %s, want %s", path, actual, files[path])
		}
		fmt.Fprintf(&identity, "%s\t%s\n", path, actual)
	}
	verifyCandidateSetDigest(t, identity.String(), wantSetDigest)
}

func verifyCandidateSetAtCommit(t *testing.T, repositoryRoot, commit string, files map[string]string, wantSetDigest string) {
	t.Helper()
	var identity strings.Builder
	for _, path := range sortedCandidatePaths(files) {
		contents, err := exec.Command("git", "-C", repositoryRoot, "show", commit+":"+path).Output()
		if err != nil {
			t.Fatalf("read %s at %s: %v", path, commit, err)
		}
		normalized := strings.ReplaceAll(strings.ReplaceAll(string(contents), "\r\n", "\n"), "\r", "\n")
		actual := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(normalized)))
		if actual != files[path] {
			t.Fatalf("candidate file %s at %s digest = %s, want %s", path, commit, actual, files[path])
		}
		fmt.Fprintf(&identity, "%s\t%s\n", path, actual)
	}
	verifyCandidateSetDigest(t, identity.String(), wantSetDigest)
}

func sortedCandidatePaths(files map[string]string) []string {
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

func verifyCandidateSetDigest(t *testing.T, identity, want string) {
	t.Helper()
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(identity)))
	if digest != want {
		t.Fatalf("candidate set digest = %s, want %s", digest, want)
	}
}
