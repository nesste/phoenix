package main

import (
	"path/filepath"
	"testing"
)

// TestAbsenceAcceptanceCandidateMatchesHistoricalPayload verifies the
// protocol-v5 section 5 absence-acceptance payload inventory against the
// accepted payload commit: every listed file's LF-normalized digest and the
// ordinal-path set digest must reproduce from Git at that commit. The
// classified message set gained its countersignature reference at the
// refreeze, so this guard pins the payload commit rather than the working
// tree; the refrozen set is pinned separately in pre-validation-artifacts.
func TestAbsenceAcceptanceCandidateMatchesHistoricalPayload(t *testing.T) {
	candidate := loadReviewCandidate(t, "gate-1a-absence-acceptance-candidate.json",
		"6e0f35afe6128a971369a710256691f6cc8de311", 18)
	root := filepath.Join("..", "..", "..")
	verifyCandidateSetAtCommit(t, root, "54e256e9f4347d34844a3ef9b2156600360dcd13", candidate.Files, candidate.SetDigest)
}
