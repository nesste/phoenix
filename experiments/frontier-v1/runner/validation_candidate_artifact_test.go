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

func TestGate1ATestTransitionCandidateMatchesHistoricalPayload(t *testing.T) {
	var candidate struct {
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
	readJSONForTest(t, filepath.Join("..", "artifacts", "gate-1a-test-transition-candidate.json"), &candidate)
	if candidate.V != 1 || candidate.Status != "review_candidate" || candidate.SourceLimit != "public_validation_inputs_only" ||
		candidate.BaseCommit != "94761cf43b938f643173b486c7f1fca5149e8a08" || candidate.Frozen || candidate.Gates.MayOpenValidation ||
		candidate.Gates.MayOpenHeldOut || candidate.ExecutionPerformed {
		t.Fatalf("validation execution candidate state = %#v", candidate)
	}
	if len(candidate.Files) != 16 {
		t.Fatalf("validation execution candidate file count = %d", len(candidate.Files))
	}
	root := filepath.Join("..", "..", "..")
	verifyCandidateSetAtCommit(t, root, "71f9648789decf4cd56ef8a24bc840b0dda7efd9", candidate.Files, candidate.SetDigest)
}

func TestGate1AAuthoringRepairCandidateMatchesHistoricalPayload(t *testing.T) {
	var candidate struct {
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
	readJSONForTest(t, filepath.Join("..", "artifacts", "gate-1a-authoring-repair-candidate.json"), &candidate)
	if candidate.V != 1 || candidate.Status != "review_candidate" || candidate.SourceLimit != "public_validation_inputs_only" ||
		candidate.BaseCommit != "0ffb8f9b44edca25582ec35541d25257d6a36dd0" || candidate.Frozen || candidate.Gates.MayOpenValidation ||
		candidate.Gates.MayOpenHeldOut || candidate.ExecutionPerformed {
		t.Fatalf("authoring repair candidate state = %#v", candidate)
	}
	if len(candidate.Files) != 19 {
		t.Fatalf("authoring repair candidate file count = %d", len(candidate.Files))
	}
	root := filepath.Join("..", "..", "..")
	verifyCandidateSetAtCommit(t, root, "610fa588488579cf5551a627795d4dc7f771f053", candidate.Files, candidate.SetDigest)
}

func verifyCandidateSetAtCommit(t *testing.T, repositoryRoot, commit string, files map[string]string, wantSetDigest string) {
	t.Helper()
	paths := make([]string, 0, len(files))
	for path := range files {
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
		if actual != files[path] {
			t.Fatalf("candidate file %s at %s digest = %s, want %s", path, commit, actual, files[path])
		}
		fmt.Fprintf(&identity, "%s\t%s\n", path, actual)
	}
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(identity.String())))
	if digest != wantSetDigest {
		t.Fatalf("candidate set digest at %s = %s, want %s", commit, digest, wantSetDigest)
	}
}
