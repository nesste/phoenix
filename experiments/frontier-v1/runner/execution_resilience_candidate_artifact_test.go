package main

import (
	"path/filepath"
	"testing"
)

// TestExecutionResilienceCandidateMatchesHistoricalPayload verifies the
// protocol-v5 section 9 execution-resilience payload inventory against the
// accepted payload commit: every listed file's LF-normalized digest and the
// ordinal-path set digest must reproduce from Git at that commit. The
// candidate note gained the seventh review's corrections at the refreeze, so
// this guard pins the payload commit rather than the working tree; the
// refrozen set is pinned separately in pre-validation-artifacts.
func TestExecutionResilienceCandidateMatchesHistoricalPayload(t *testing.T) {
	candidate := loadReviewCandidate(t, "gate-1a-execution-resilience-candidate.json",
		"3fc001e247332a5dfdb9d99fc54fc0a9681e8ea1", 16)
	root := filepath.Join("..", "..", "..")
	verifyCandidateSetAtCommit(t, root, "27edb8664a060d0a7690039e8e29dba94facada0", candidate.Files, candidate.SetDigest)
}
