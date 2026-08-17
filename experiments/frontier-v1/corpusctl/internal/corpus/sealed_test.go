package corpus

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildSealedManifestUsesDigestsWithoutLabelContent(t *testing.T) {
	sourceRoot := repoRoot(t)
	root := t.TempDir()
	copyTestGlob(t, sourceRoot, root, "experiments/frontier-v1/schema/*.schema.json")
	copyTestGlob(t, sourceRoot, root, "experiments/frontier-v1/corpusctl/cmd/corpusctl/*.go")
	copyTestGlob(t, sourceRoot, root, "experiments/frontier-v1/corpusctl/internal/corpus/*.go")
	copyTestFile(t, sourceRoot, root, "experiments/frontier-v1/corpusctl/go.mod")
	copyTestFile(t, sourceRoot, root, "experiments/frontier-v1/corpusctl/go.sum")
	copyTestFile(t, sourceRoot, root, "spec/world.schema.json")
	copyTestFile(t, sourceRoot, root, "experiments/frontier-v1/worlds/authoring.dev_repo.json")
	copyTestFile(t, sourceRoot, root, "experiments/frontier-v1/fixtures/authoring/authoring_family_001_v1.json")

	const worldSource = "experiments/frontier-v1/worlds/authoring.dev_repo.json"
	worldDigest, err := DigestFile(filepath.Join(root, filepath.FromSlash(worldSource)))
	if err != nil {
		t.Fatal(err)
	}
	fixturePath := "experiments/frontier-v1/fixtures/validation/validation_family_001_v1.json"
	fixtureValue := map[string]any{
		"v": json.Number("1"), "fixture_id": "validation_family_001_v1", "family_id": "validation_family_001",
		"template": "Independent validation template that is absent from every authoring fixture family.",
		"files":    map[string]any{"README.md": "sealed fixture\n"},
	}
	writeTestJSON(t, root, fixturePath, fixtureValue)
	fixtureDigest, err := DigestFile(filepath.Join(root, filepath.FromSlash(fixturePath)))
	if err != nil {
		t.Fatal(err)
	}
	graderDigest, err := GraderDigest(root)
	if err != nil {
		t.Fatal(err)
	}
	registry := map[string]any{
		"v": json.Number("1"), "tranche": "validation", "labels": []any{},
	}
	labels := registry["labels"].([]any)
	for index, class := range requiredClasses {
		caseID := "validation_" + class + "_001"
		writeTestJSON(t, root, "experiments/frontier-v1/corpus/validation/"+caseID+".json", map[string]any{
			"case_id": caseID, "class": class,
			"goal":            "Evaluate this sealed validation behavior without receiving outcome hints.",
			"sandbox_fixture": fixtureDigest, "world_ref": worldDigest, "family_id": "validation_family_001",
		})
		labels = append(labels, map[string]any{
			"case_id": caseID, "label_digest": fmt.Sprintf("sha256:%064x", index+1), "grading_script": graderDigest,
		})
	}
	registry["labels"] = labels
	const registrySource = "experiments/frontier-v1/manifests/validation-label-digests.json"
	writeTestJSON(t, root, registrySource, registry)

	data, err := BuildSealedManifest(root, "validation", worldSource, registrySource)
	if err != nil {
		t.Fatal(err)
	}
	var result manifest
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	if !result.Sealed || result.Counts.Cases != len(requiredClasses) || result.Counts.Labels != len(requiredClasses) {
		t.Fatalf("unexpected sealed manifest counts: %+v", result.Counts)
	}

	authoringDocument, err := loadDocument(root, "experiments/frontier-v1/fixtures/authoring/authoring_family_001_v1.json")
	if err != nil {
		t.Fatal(err)
	}
	var authoringFixture fixture
	if err := authoringDocument.Into(&authoringFixture); err != nil {
		t.Fatal(err)
	}
	originalTemplate := fixtureValue["template"]
	fixtureValue["template"] = authoringFixture.Template
	writeTestJSON(t, root, fixturePath, fixtureValue)
	if _, err := BuildSealedManifest(root, "validation", worldSource, registrySource); err == nil || !strings.Contains(err.Error(), "authoring tranche") {
		t.Fatalf("expected cross-tranche template rejection, got %v", err)
	}
	fixtureValue["template"] = originalTemplate
	writeTestJSON(t, root, fixturePath, fixtureValue)

	writeTestJSON(t, root, "experiments/frontier-v1/labels/validation/validation_direct_001.json", map[string]any{})
	if _, err := BuildSealedManifest(root, "validation", worldSource, registrySource); err == nil || !strings.Contains(err.Error(), "label content exists") {
		t.Fatalf("expected sealed label-content rejection, got %v", err)
	}
}

func copyTestGlob(t *testing.T, sourceRoot, targetRoot, pattern string) {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(sourceRoot, filepath.FromSlash(pattern)))
	if err != nil {
		t.Fatal(err)
	}
	for _, source := range matches {
		if strings.HasSuffix(source, "_test.go") {
			continue
		}
		relative, err := filepath.Rel(sourceRoot, source)
		if err != nil {
			t.Fatal(err)
		}
		copyTestFile(t, sourceRoot, targetRoot, filepath.ToSlash(relative))
	}
}

func copyTestFile(t *testing.T, sourceRoot, targetRoot, relative string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(sourceRoot, filepath.FromSlash(relative)))
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(targetRoot, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeTestJSON(t *testing.T, root, relative string, value any) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}
