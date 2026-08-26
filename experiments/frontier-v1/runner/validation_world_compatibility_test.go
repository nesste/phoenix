package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestFrozenValidationCasesMatchProductionWorld(t *testing.T) {
	repositoryRoot := filepath.Join("..", "..", "..")
	ids, err := selectedCaseIDsForTranche(repositoryRoot, "validation", "all")
	if err != nil {
		t.Fatal(err)
	}
	cases, err := loadSelectedCasesForTranche(repositoryRoot, "validation", ids)
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) != 120 {
		t.Fatalf("validation case count = %d, want 120", len(cases))
	}
	config := runConfig{
		repositoryRoot: repositoryRoot,
		worldPath:      filepath.Join(repositoryRoot, "worlds", "dev-repo", "world.json"),
	}
	if err := preflightSelectedCasesForTranche(config, "validation", cases); err != nil {
		t.Fatal(err)
	}
}

func TestValidationWorldPreflightRejectsMismatchedWorld(t *testing.T) {
	repositoryRoot := filepath.Join("..", "..", "..")
	config := runConfig{
		repositoryRoot: repositoryRoot,
		worldPath:      filepath.Join(repositoryRoot, "worlds", "dev-repo", "world.json"),
	}
	cases := []runnableCase{{
		CaseID:   "validation_mismatch",
		WorldRef: "sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b",
	}}
	err := preflightSelectedCasesForTranche(config, "validation", cases)
	if err == nil || !strings.Contains(err.Error(), "does not match production world") {
		t.Fatalf("mismatched world preflight error = %v", err)
	}
}
