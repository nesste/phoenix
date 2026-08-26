package main

import (
	"path/filepath"
	"testing"
)

func TestFrozenValidationCasesMatchProductionWorld(t *testing.T) {
	repositoryRoot := filepath.Join("..", "..", "..")
	document, err := readValidationGateDocument(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	if document.Gates.MayOpenValidation || document.Gates.MayOpenHeldOut {
		t.Fatal("compatibility candidate must be tested with both outcome gates closed")
	}

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
