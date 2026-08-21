package main

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestGate1ATestTransitionCandidateMatchesWorkingTree(t *testing.T) {
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
	paths := make([]string, 0, len(candidate.Files))
	for path := range candidate.Files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	var identity strings.Builder
	for _, path := range paths {
		actual, err := digestLFNormalizedFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			t.Fatal(err)
		}
		if actual != candidate.Files[path] {
			t.Fatalf("candidate file %s digest = %s, want %s", path, actual, candidate.Files[path])
		}
		fmt.Fprintf(&identity, "%s\t%s\n", path, actual)
	}
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(identity.String())))
	if digest != candidate.SetDigest {
		t.Fatalf("validation execution candidate set digest = %s, want %s", digest, candidate.SetDigest)
	}
}
