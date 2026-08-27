package main

import (
	"path/filepath"
	"testing"
)

// TestProtocolV5AmendmentCandidateMatchesHistoricalPayload verifies the
// protocol-v5 amendment payload inventory against the accepted payload
// commit: every listed file's LF-normalized digest and the ordinal-path set
// digest must reproduce from Git at that commit. Two set members (the README
// and the execution-boundary document) were amended at the refreeze to update
// the world-build literal per review P3-2, so this guard pins the payload
// commit rather than the working tree.
func TestProtocolV5AmendmentCandidateMatchesHistoricalPayload(t *testing.T) {
	candidate := loadReviewCandidate(t, "gate-1a-protocol-v5-amendment-candidate.json",
		"24895c4a272167cf8ea1445d3b81b2133ae62c4c", 62)
	root := filepath.Join("..", "..", "..")
	verifyCandidateSetAtCommit(t, root, "0f8d9c72c59bea5da5abf792f493f1d75299b1f8", candidate.Files, candidate.SetDigest)
}
