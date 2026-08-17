package corpus

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var requiredClasses = []string{
	"absence",
	"adversarial_text",
	"cascade",
	"direct",
	"far_discovery",
	"recovery",
	"stale_frontier",
	"temptation",
}

type runnableCase struct {
	CaseID         string `json:"case_id"`
	SandboxFixture string `json:"sandbox_fixture"`
	WorldRef       string `json:"world_ref"`
	FamilyID       string `json:"family_id"`
}

type fixture struct {
	FixtureID   string            `json:"fixture_id"`
	FamilyID    string            `json:"family_id"`
	Template    string            `json:"template"`
	Files       map[string]string `json:"files"`
	Uncommitted []string          `json:"uncommitted"`
}

type label struct {
	CaseID          string `json:"case_id"`
	Class           string `json:"class"`
	GradingScript   string `json:"grading_script"`
	ExpectedOutcome struct {
		Checks []struct {
			Kind string `json:"kind"`
		} `json:"checks"`
	} `json:"expected_outcome"`
}

type manifest struct {
	V              int                       `json:"v"`
	Tranche        string                    `json:"tranche"`
	Sealed         bool                      `json:"sealed"`
	Generator      string                    `json:"generator"`
	WorldRefs      []pinnedRef               `json:"world_refs"`
	Counts         counts                    `json:"counts"`
	ClassCoverage  map[string]int            `json:"class_coverage"`
	FamilyCoverage map[string]familyCoverage `json:"family_coverage"`
	Grading        grading                   `json:"grading"`
	Fixtures       []fixtureEntry            `json:"fixtures"`
	Cases          []caseEntry               `json:"cases"`
}

type pinnedRef struct {
	Digest string `json:"digest"`
	Source string `json:"source"`
}

type counts struct {
	Cases    int `json:"cases"`
	Families int `json:"families"`
	Fixtures int `json:"fixtures"`
	Labels   int `json:"labels"`
}

type familyCoverage struct {
	Classes []string `json:"classes"`
	Cases   int      `json:"cases"`
}

type grading struct {
	LabelledCases          int            `json:"labelled_cases"`
	CasesWithGradingScript int            `json:"cases_with_grading_script"`
	Checks                 int            `json:"checks"`
	ChecksByKind           map[string]int `json:"checks_by_kind"`
}

type fixtureEntry struct {
	FixtureID string `json:"fixture_id"`
	FamilyID  string `json:"family_id"`
	Digest    string `json:"digest"`
}

type caseEntry struct {
	CaseID         string  `json:"case_id"`
	Class          string  `json:"class"`
	FamilyID       string  `json:"family_id"`
	InputDigest    string  `json:"input_digest"`
	SandboxFixture string  `json:"sandbox_fixture"`
	WorldRef       string  `json:"world_ref"`
	LabelDigest    *string `json:"label_digest"`
}

// BuildManifest validates a tranche and returns its deterministic manifest.
func BuildManifest(root, tranche, worldSource string) ([]byte, error) {
	if tranche != "authoring" {
		return nil, fmt.Errorf("%s: only the authoring tranche may be built in the implementation workspace", tranche)
	}
	schemas, err := loadSchemas(root)
	if err != nil {
		return nil, err
	}
	if err := rejectWithheldLabelContent(root, schemas); err != nil {
		return nil, err
	}
	worldDocument, err := loadDocument(root, filepath.ToSlash(worldSource))
	if err != nil {
		return nil, fmt.Errorf("load world source: %w", err)
	}
	if err := validateWithSchema(root, "spec/world.schema.json", worldDocument); err != nil {
		return nil, err
	}
	caseDocuments, err := loadDocuments(root, experimentDir+"/corpus/"+tranche+"/*.json")
	if err != nil {
		return nil, err
	}
	fixtureDocuments, err := loadDocuments(root, experimentDir+"/fixtures/"+tranche+"/*.json")
	if err != nil {
		return nil, err
	}
	labelDocuments, err := loadDocuments(root, experimentDir+"/labels/"+tranche+"/*.json")
	if err != nil {
		return nil, err
	}

	fixtures, fixtureEntries, err := validateFixtures(schemas, fixtureDocuments, tranche)
	if err != nil {
		return nil, err
	}
	labels, labelDigests, gradeSummary, err := validateLabels(root, schemas, labelDocuments, tranche)
	if err != nil {
		return nil, err
	}
	caseEntries, classCoverage, familyCoverage, err := validateCases(
		schemas, caseDocuments, tranche, worldDocument.Digest, fixtures, labels, labelDigests,
	)
	if err != nil {
		return nil, err
	}
	for _, class := range requiredClasses {
		if classCoverage[class] == 0 {
			return nil, fmt.Errorf("authoring corpus has no %s case", class)
		}
	}

	result := manifest{
		V:         1,
		Tranche:   tranche,
		Sealed:    false,
		Generator: "go run ./cmd/corpusctl manifest --repo-root ../../.. --tranche authoring --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --write",
		WorldRefs: []pinnedRef{{Digest: worldDocument.Digest, Source: filepath.ToSlash(worldSource)}},
		Counts: counts{
			Cases: len(caseEntries), Families: len(familyCoverage), Fixtures: len(fixtureEntries), Labels: len(labels),
		},
		ClassCoverage: classCoverage, FamilyCoverage: familyCoverage, Grading: gradeSummary,
		Fixtures: fixtureEntries, Cases: caseEntries,
	}
	return encodeManifest(schemas, result)
}

func encodeManifest(schemas schemaSet, result manifest) ([]byte, error) {
	encoded, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}
	value, err := DecodeJSON(encoded)
	if err != nil {
		return nil, err
	}
	if err := schemas.validateValue("manifest", value); err != nil {
		return nil, err
	}
	return append(encoded, '\n'), nil
}

// WriteManifest validates and writes the authoring manifest.
func WriteManifest(root, tranche, worldSource string) (string, error) {
	data, err := BuildManifest(root, tranche, worldSource)
	if err != nil {
		return "", err
	}
	relative := experimentDir + "/manifests/" + tranche + ".json"
	target := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(target, data, 0o644); err != nil {
		return "", err
	}
	return relative, nil
}

func loadDocuments(root, pattern string) ([]Document, error) {
	paths, err := filepath.Glob(filepath.Join(root, filepath.FromSlash(pattern)))
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no documents match %s", pattern)
	}
	sort.Strings(paths)
	documents := make([]Document, 0, len(paths))
	for _, filePath := range paths {
		relative, err := filepath.Rel(root, filePath)
		if err != nil {
			return nil, err
		}
		document, err := loadDocument(root, filepath.ToSlash(relative))
		if err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}
	return documents, nil
}

func validateFixtures(schemas schemaSet, documents []Document, tranche string) (map[string]fixture, []fixtureEntry, error) {
	byDigest := make(map[string]fixture, len(documents))
	entries := make([]fixtureEntry, 0, len(documents))
	for _, document := range documents {
		if err := schemas.validate("fixture", document); err != nil {
			return nil, nil, err
		}
		var item fixture
		if err := document.Into(&item); err != nil {
			return nil, nil, err
		}
		if document.Name() != item.FixtureID {
			return nil, nil, fmt.Errorf("%s: fixture filename and fixture_id must agree", document.Path)
		}
		if !strings.HasPrefix(item.FamilyID, tranche+"_") {
			return nil, nil, fmt.Errorf("%s: family %s belongs to another tranche", document.Path, item.FamilyID)
		}
		if _, exists := byDigest[document.Digest]; exists {
			return nil, nil, fmt.Errorf("%s: duplicate fixture content digest %s", document.Path, document.Digest)
		}
		for _, file := range item.Uncommitted {
			if _, exists := item.Files[file]; !exists {
				return nil, nil, fmt.Errorf("%s: uncommitted path %s is absent from files", document.Path, file)
			}
		}
		byDigest[document.Digest] = item
		entries = append(entries, fixtureEntry{FixtureID: item.FixtureID, FamilyID: item.FamilyID, Digest: document.Digest})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].FixtureID < entries[j].FixtureID })
	return byDigest, entries, nil
}

func validateLabels(root string, schemas schemaSet, documents []Document, tranche string) (map[string]label, map[string]string, grading, error) {
	labels := make(map[string]label, len(documents))
	digests := make(map[string]string, len(documents))
	summary := grading{ChecksByKind: map[string]int{}}
	graderDigest, err := GraderDigest(root)
	if err != nil {
		return nil, nil, grading{}, err
	}
	for _, document := range documents {
		if err := schemas.validate("label", document); err != nil {
			return nil, nil, grading{}, err
		}
		var item label
		if err := document.Into(&item); err != nil {
			return nil, nil, grading{}, err
		}
		if document.Name() != item.CaseID || !strings.HasPrefix(item.CaseID, tranche+"_") {
			return nil, nil, grading{}, fmt.Errorf("%s: label filename, case_id, and tranche must agree", document.Path)
		}
		if _, exists := labels[item.CaseID]; exists {
			return nil, nil, grading{}, fmt.Errorf("duplicate label %s", item.CaseID)
		}
		labels[item.CaseID] = item
		digests[item.CaseID] = document.Digest
		summary.LabelledCases++
		if item.GradingScript != graderDigest {
			return nil, nil, grading{}, fmt.Errorf("%s: grading_script %s does not match current grader %s", document.Path, item.GradingScript, graderDigest)
		}
		summary.CasesWithGradingScript++
		for _, check := range item.ExpectedOutcome.Checks {
			summary.Checks++
			summary.ChecksByKind[check.Kind]++
		}
	}
	return labels, digests, summary, nil
}

func validateCases(
	schemas schemaSet,
	documents []Document,
	tranche, worldDigest string,
	fixtures map[string]fixture,
	labels map[string]label,
	labelDigests map[string]string,
) ([]caseEntry, map[string]int, map[string]familyCoverage, error) {
	entries := make([]caseEntry, 0, len(documents))
	classCoverage := map[string]int{}
	familyClasses := map[string]map[string]struct{}{}
	familyCounts := map[string]int{}
	seen := map[string]struct{}{}
	for _, document := range documents {
		if err := schemas.validate("case", document); err != nil {
			return nil, nil, nil, err
		}
		var item runnableCase
		if err := document.Into(&item); err != nil {
			return nil, nil, nil, err
		}
		if document.Name() != item.CaseID || !strings.HasPrefix(item.CaseID, tranche+"_") {
			return nil, nil, nil, fmt.Errorf("%s: case filename, case_id, and tranche must agree", document.Path)
		}
		if _, exists := seen[item.CaseID]; exists {
			return nil, nil, nil, fmt.Errorf("duplicate case %s", item.CaseID)
		}
		seen[item.CaseID] = struct{}{}
		fixture, ok := fixtures[item.SandboxFixture]
		if !ok {
			return nil, nil, nil, fmt.Errorf("%s: sandbox_fixture %s does not match a fixture document", document.Path, item.SandboxFixture)
		}
		if fixture.FamilyID != item.FamilyID {
			return nil, nil, nil, fmt.Errorf("%s: case family %s does not match fixture family %s", document.Path, item.FamilyID, fixture.FamilyID)
		}
		if item.WorldRef != worldDigest {
			return nil, nil, nil, fmt.Errorf("%s: world_ref does not match pinned world digest %s", document.Path, worldDigest)
		}
		caseLabel, ok := labels[item.CaseID]
		if !ok {
			return nil, nil, nil, fmt.Errorf("%s: missing label digest", document.Path)
		}
		labelDigest := labelDigests[item.CaseID]
		entries = append(entries, caseEntry{
			CaseID: item.CaseID, Class: caseLabel.Class, FamilyID: item.FamilyID, InputDigest: document.Digest,
			SandboxFixture: item.SandboxFixture, WorldRef: item.WorldRef, LabelDigest: &labelDigest,
		})
		classCoverage[caseLabel.Class]++
		familyCounts[item.FamilyID]++
		if familyClasses[item.FamilyID] == nil {
			familyClasses[item.FamilyID] = map[string]struct{}{}
		}
		familyClasses[item.FamilyID][caseLabel.Class] = struct{}{}
	}
	if len(labels) != len(seen) {
		return nil, nil, nil, fmt.Errorf("authoring labels (%d) do not map one-to-one to cases (%d)", len(labels), len(seen))
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].CaseID < entries[j].CaseID })
	families := make(map[string]familyCoverage, len(familyCounts))
	for familyID, count := range familyCounts {
		classes := make([]string, 0, len(familyClasses[familyID]))
		for class := range familyClasses[familyID] {
			classes = append(classes, class)
		}
		sort.Strings(classes)
		families[familyID] = familyCoverage{Classes: classes, Cases: count}
	}
	return entries, classCoverage, families, nil
}
