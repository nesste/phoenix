package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestProtocolV5AmendmentCandidateMatchesWorkingTree verifies the protocol-v5
// amendment payload inventory against the working tree: every listed file's
// LF-normalized digest and the ordinal-path set digest must match the
// committed candidate document. At refreeze this guard becomes a historical
// check pinned to the accepted payload commit.
func TestProtocolV5AmendmentCandidateMatchesWorkingTree(t *testing.T) {
	candidate := loadReviewCandidate(t, "gate-1a-protocol-v5-amendment-candidate.json",
		"24895c4a272167cf8ea1445d3b81b2133ae62c4c", 62)
	root := filepath.Join("..", "..", "..")
	var identity strings.Builder
	for _, path := range sortedCandidatePaths(candidate.Files) {
		contents, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			t.Fatalf("read candidate file %s: %v", path, err)
		}
		normalized := strings.ReplaceAll(strings.ReplaceAll(string(contents), "\r\n", "\n"), "\r", "\n")
		actual := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(normalized)))
		if actual != candidate.Files[path] {
			t.Fatalf("candidate file %s digest = %s, want %s", path, actual, candidate.Files[path])
		}
		fmt.Fprintf(&identity, "%s\t%s\n", path, actual)
	}
	verifyCandidateSetDigest(t, identity.String(), candidate.SetDigest)
}
