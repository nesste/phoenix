package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nesste/phoenix/internal/surface"
	"github.com/nesste/phoenix/internal/world"
)

func TestArmBDocumentCoversFlatToolsAndAuthoringPaths(t *testing.T) {
	repositoryRoot := filepath.Join("..", "..", "..")
	documentPath := filepath.Join("..", "arms", "arm-b.md")
	contents, err := os.ReadFile(documentPath)
	if err != nil {
		t.Fatal(err)
	}
	text := string(contents)
	assertArmBToolCoverage(t, repositoryRoot, text)
	assertArmBPathCoverage(t, text)
	assertArmBTerminology(t, text)
}

func assertArmBToolCoverage(t *testing.T, repositoryRoot, text string) {
	t.Helper()
	toolSection, _, ok := strings.Cut(text, "## Worked paths")
	if !ok {
		t.Fatal("Arm B document has no worked-path boundary")
	}
	definition, err := world.Load(
		filepath.Join(repositoryRoot, "spec", "world.schema.json"),
		filepath.Join(repositoryRoot, "worlds", "dev-repo", "world.json"),
	)
	if err != nil {
		t.Fatal(err)
	}
	toolNames, err := surface.FlatToolNames(definition)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range toolNames {
		rowPrefix := "| `" + name + "` |"
		if strings.Count(toolSection, rowPrefix) != 1 {
			t.Errorf("Arm B tool table contains %q %d times, want once", name, strings.Count(toolSection, rowPrefix))
		}
	}
}

func assertArmBPathCoverage(t *testing.T, text string) {
	t.Helper()
	labels, err := filepath.Glob(filepath.Join("..", "labels", "authoring", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range labels {
		var label struct {
			Class           string     `json:"class"`
			AcceptablePaths [][]string `json:"acceptable_paths"`
		}
		encoded, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(encoded, &label); err != nil {
			t.Fatal(err)
		}
		covered := false
		for _, acceptable := range label.AcceptablePaths {
			if len(acceptable) == 0 && strings.Contains(text, "No call, or one `repo_find` check") {
				covered = true
				break
			}
			flat := make([]string, len(acceptable))
			for index, action := range acceptable {
				flat[index] = "`" + strings.ReplaceAll(action, ".", "_") + "`"
			}
			if strings.Contains(text, strings.Join(flat, " → ")) {
				covered = true
				break
			}
		}
		if !covered {
			t.Errorf("Arm B document contains no acceptable path for %s", label.Class)
		}
	}
}

func assertArmBTerminology(t *testing.T, text string) {
	t.Helper()
	for _, forbidden := range []string{"Phoenix handle", "orientation", "frontier", "validation", "held_out"} {
		if strings.Contains(text, forbidden) {
			t.Errorf("Arm B document leaks treatment or sealed-tranche term %q", forbidden)
		}
	}
}
