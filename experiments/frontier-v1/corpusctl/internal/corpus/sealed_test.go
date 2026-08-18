package corpus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommittedSealedManifestsReproduceReviewedCounts(t *testing.T) {
	root := repoRoot(t)
	const worldSource = "experiments/frontier-v1/worlds/authoring.dev_repo.json"
	wantClasses := map[string]int{
		"direct": 24, "recovery": 24, "cascade": 12, "far_discovery": 12,
		"temptation": 12, "absence": 12, "stale_frontier": 12, "adversarial_text": 12,
	}
	for _, tranche := range []string{"validation", "held_out"} {
		t.Run(tranche, func(t *testing.T) {
			registry := "experiments/frontier-v1/manifests/" + tranche + "-label-digests.json"
			generated, err := BuildSealedManifest(root, tranche, worldSource, registry)
			if err != nil {
				t.Fatal(err)
			}
			manifestPath := filepath.Join(root, "experiments", "frontier-v1", "manifests", tranche+".json")
			committed, err := os.ReadFile(manifestPath)
			if err != nil {
				t.Fatal(err)
			}
			canonicalCommitted := bytes.ReplaceAll(committed, []byte("\r\n"), []byte("\n"))
			if !bytes.Equal(generated, canonicalCommitted) {
				t.Fatalf("%s sealed manifest does not reproduce after LF normalization", tranche)
			}
			var result manifest
			if err := json.Unmarshal(committed, &result); err != nil {
				t.Fatal(err)
			}
			if !result.Sealed || result.Counts.Cases != 120 || result.Counts.Families != 24 || result.Counts.Labels != 120 || len(result.FamilyCoverage) != 24 {
				t.Fatalf("unexpected %s reviewed counts: %+v", tranche, result.Counts)
			}
			for class, want := range wantClasses {
				if result.ClassCoverage[class] != want {
					t.Fatalf("%s class %s has %d cases, want %d", tranche, class, result.ClassCoverage[class], want)
				}
			}
			for family, coverage := range result.FamilyCoverage {
				if coverage.Cases != 5 {
					t.Fatalf("%s family %s has %d cases, want 5", tranche, family, coverage.Cases)
				}
			}
		})
	}
}

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
		caseID := fmt.Sprintf("validation_%08x", index+1)
		writeTestJSON(t, root, "experiments/frontier-v1/corpus/validation/"+caseID+".json", map[string]any{
			"case_id":         caseID,
			"goal":            "Evaluate this sealed validation behavior without receiving outcome hints.",
			"sandbox_fixture": fixtureDigest, "world_ref": worldDigest, "family_id": "validation_family_001",
		})
		labels = append(labels, map[string]any{
			"case_id": caseID, "class": class, "label_digest": fmt.Sprintf("sha256:%064x", index+1), "grading_script": graderDigest,
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

	registry["labels"] = labels[:len(labels)-1]
	writeTestJSON(t, root, registrySource, registry)
	if _, err := BuildSealedManifest(root, "validation", worldSource, registrySource); err == nil || !strings.Contains(err.Error(), "missing label digest") {
		t.Fatalf("expected incomplete registry rejection, got %v", err)
	}
	registry["labels"] = append(labels, map[string]any{
		"case_id": "validation_ffffffff", "class": "direct",
		"label_digest": fmt.Sprintf("sha256:%064x", 999), "grading_script": graderDigest,
	})
	writeTestJSON(t, root, registrySource, registry)
	if _, err := BuildSealedManifest(root, "validation", worldSource, registrySource); err == nil || !strings.Contains(err.Error(), "do not map one-to-one") {
		t.Fatalf("expected over-complete registry rejection, got %v", err)
	}
	registry["labels"] = labels
	registry["tranche"] = "held_out"
	writeTestJSON(t, root, registrySource, registry)
	if _, err := BuildSealedManifest(root, "validation", worldSource, registrySource); err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("expected registry tranche rejection, got %v", err)
	}
	registry["tranche"] = "validation"
	writeTestJSON(t, root, registrySource, registry)

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
	originalFiles := fixtureValue["files"]
	fixtureValue["files"] = authoringFixture.Files
	writeTestJSON(t, root, fixturePath, fixtureValue)
	if _, err := BuildSealedManifest(root, "validation", worldSource, registrySource); err == nil || !strings.Contains(err.Error(), "duplicates fixture content") {
		t.Fatalf("expected cross-tranche fixture-content rejection, got %v", err)
	}
	fixtureValue["files"] = originalFiles
	writeTestJSON(t, root, fixturePath, fixtureValue)

	writeTestJSON(t, root, "notes/private/validation_deadbeef.json", map[string]any{
		"case_id": "validation_deadbeef", "class": "direct",
		"expected_outcome": map[string]any{"checks": []any{map[string]any{
			"id": "reports_result", "kind": "final_message_matches", "pattern": "result",
		}}},
		"grading_script":   graderDigest,
		"acceptable_paths": []any{[]any{}},
		"label_rationale":  "This deliberately leaked sealed label exists only to prove the repository-wide guard.",
	})
	if _, err := BuildSealedManifest(root, "validation", worldSource, registrySource); err == nil || !strings.Contains(err.Error(), "sealed label content") {
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
