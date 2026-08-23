package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRequireFrozenWorldBuildRejectsMismatch(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	err := requireFrozenWorldBuild(root, "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	if err == nil || !strings.Contains(err.Error(), "does not match frozen world-build") {
		t.Fatalf("mismatch error = %v", err)
	}
}

func TestRequireFrozenWorldBuildAcceptsFrozenDigest(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	document, err := readValidationGateDocument(root)
	if err != nil {
		t.Fatal(err)
	}
	frozen := document.Artifacts.WorldDefinitionAndWorldBuildDigest.WorldBuildDigest
	if frozen == "" {
		t.Fatal("freeze document is missing world_build_digest")
	}
	if err := requireFrozenWorldBuild(root, frozen); err != nil {
		t.Fatal(err)
	}
}

func TestRequireFrozenWorldBuildRejectsMissingFreeze(t *testing.T) {
	root := t.TempDir()
	gate := map[string]any{
		"status": "complete",
		"gates":  map[string]any{"may_open_validation": true, "may_open_held_out": false},
		"artifacts": map[string]any{
			"world_definition_and_world_build_digest": map[string]any{
				"frozen": false, "world_build_digest": "",
			},
		},
	}
	if err := writeJSON(filepath.Join(root, filepath.FromSlash(preValidationArtifactsPath)), gate); err != nil {
		t.Fatal(err)
	}
	if err := requireFrozenWorldBuild(root, "sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d"); err == nil || !strings.Contains(err.Error(), "frozen world-build digest") {
		t.Fatalf("missing freeze error = %v", err)
	}
}

func TestRequireFrozenWorldBuildRejectsEmptyLiveDigest(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	if err := requireFrozenWorldBuild(root, ""); err == nil || !strings.Contains(err.Error(), "live world-build digest") {
		t.Fatalf("empty live digest error = %v", err)
	}
}
