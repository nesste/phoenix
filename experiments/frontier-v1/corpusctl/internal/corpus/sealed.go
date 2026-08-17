package corpus

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type withheldRegistry struct {
	Tranche string                `json:"tranche"`
	Labels  []withheldLabelDigest `json:"labels"`
}

type withheldLabelDigest struct {
	CaseID        string `json:"case_id"`
	Class         string `json:"class"`
	LabelDigest   string `json:"label_digest"`
	GradingScript string `json:"grading_script"`
}

// BuildSealedManifest builds an outcome-free manifest from independently
// supplied label digests. Full validation or held-out labels must not exist in
// the implementation workspace.
func BuildSealedManifest(root, tranche, worldSource, registrySource string) ([]byte, error) {
	if tranche != "validation" && tranche != "held_out" {
		return nil, fmt.Errorf("sealed tranche must be validation or held_out, got %s", tranche)
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
	fixtures, fixtureEntries, err := validateFixtures(schemas, fixtureDocuments, tranche)
	if err != nil {
		return nil, err
	}
	if err := rejectCrossTrancheTemplates(root, tranche, fixtures); err != nil {
		return nil, err
	}
	registryDocument, err := loadDocument(root, filepath.ToSlash(registrySource))
	if err != nil {
		return nil, fmt.Errorf("load withheld-label registry: %w", err)
	}
	if err := schemas.validate("withheld-labels", registryDocument); err != nil {
		return nil, err
	}
	var registry withheldRegistry
	if err := registryDocument.Into(&registry); err != nil {
		return nil, err
	}
	if registry.Tranche != tranche {
		return nil, fmt.Errorf("withheld-label registry tranche %s does not match %s", registry.Tranche, tranche)
	}
	labels, labelDigests, err := validateWithheldDigests(root, tranche, registry.Labels)
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
			return nil, fmt.Errorf("%s corpus has no %s case", tranche, class)
		}
	}
	result := manifest{
		V:         1,
		Tranche:   tranche,
		Sealed:    true,
		Generator: fmt.Sprintf("go run ./cmd/corpusctl seal --repo-root ../../.. --tranche %s --world-source %s --label-digests %s --write", tranche, filepath.ToSlash(worldSource), filepath.ToSlash(registrySource)),
		WorldRefs: []pinnedRef{{Digest: worldDocument.Digest, Source: filepath.ToSlash(worldSource)}},
		Counts: counts{
			Cases: len(caseEntries), Families: len(familyCoverage), Fixtures: len(fixtureEntries), Labels: len(labelDigests),
		},
		ClassCoverage:  classCoverage,
		FamilyCoverage: familyCoverage,
		Grading: grading{
			LabelledCases: len(labelDigests), CasesWithGradingScript: len(labelDigests), ChecksByKind: map[string]int{},
		},
		Fixtures: fixtureEntries,
		Cases:    caseEntries,
	}
	return encodeManifest(schemas, result)
}

func rejectWithheldLabelContent(root string, schemas schemaSet) error {
	return filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.EqualFold(filepath.Ext(entry.Name()), ".json") {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		document, err := loadDocument(root, filepath.ToSlash(relative))
		if err != nil || schemas.validate("label", document) != nil {
			return nil
		}
		var candidate label
		if err := document.Into(&candidate); err != nil {
			return err
		}
		if strings.HasPrefix(candidate.CaseID, "validation_") || strings.HasPrefix(candidate.CaseID, "held_out_") {
			return fmt.Errorf("sealed label content for %s exists in the implementation workspace at %s", candidate.CaseID, filepath.ToSlash(relative))
		}
		return nil
	})
}

// WriteSealedManifest writes a validated outcome-free manifest.
func WriteSealedManifest(root, tranche, worldSource, registrySource string) (string, error) {
	data, err := BuildSealedManifest(root, tranche, worldSource, registrySource)
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

func validateWithheldDigests(root, tranche string, entries []withheldLabelDigest) (map[string]label, map[string]string, error) {
	graderDigest, err := GraderDigest(root)
	if err != nil {
		return nil, nil, err
	}
	labels := make(map[string]label, len(entries))
	digests := make(map[string]string, len(entries))
	for _, entry := range entries {
		if !strings.HasPrefix(entry.CaseID, tranche+"_") {
			return nil, nil, fmt.Errorf("withheld label %s belongs to another tranche", entry.CaseID)
		}
		if !IsDigest(entry.LabelDigest) {
			return nil, nil, fmt.Errorf("withheld label %s has invalid label digest %s", entry.CaseID, entry.LabelDigest)
		}
		if _, exists := labels[entry.CaseID]; exists {
			return nil, nil, fmt.Errorf("duplicate withheld label digest for %s", entry.CaseID)
		}
		if entry.GradingScript != graderDigest {
			return nil, nil, fmt.Errorf("withheld label %s pins grader %s, expected %s", entry.CaseID, entry.GradingScript, graderDigest)
		}
		labels[entry.CaseID] = label{CaseID: entry.CaseID, Class: entry.Class, GradingScript: entry.GradingScript}
		digests[entry.CaseID] = entry.LabelDigest
	}
	return labels, digests, nil
}

func rejectCrossTrancheTemplates(root, tranche string, sealed map[string]fixture) error {
	usedTemplates := map[string]string{}
	usedFileSets := map[string]string{}
	for _, other := range []string{"authoring", "validation", "held_out"} {
		if other == tranche {
			continue
		}
		matches, err := filepath.Glob(filepath.Join(root, filepath.FromSlash(experimentDir+"/fixtures/"+other+"/*.json")))
		if err != nil {
			return err
		}
		for _, match := range matches {
			relative, err := filepath.Rel(root, match)
			if err != nil {
				return err
			}
			document, err := loadDocument(root, filepath.ToSlash(relative))
			if err != nil {
				return err
			}
			var item fixture
			if err := document.Into(&item); err != nil {
				return err
			}
			usedTemplates[item.Template] = other
			filesDigest, err := digestFixtureFiles(item.Files)
			if err != nil {
				return fmt.Errorf("fingerprint fixture %s: %w", item.FixtureID, err)
			}
			usedFileSets[filesDigest] = other
		}
	}
	for _, item := range sealed {
		if other, reused := usedTemplates[item.Template]; reused {
			return fmt.Errorf("sealed fixture %s reuses a generating template from the %s tranche", item.FixtureID, other)
		}
		filesDigest, err := digestFixtureFiles(item.Files)
		if err != nil {
			return fmt.Errorf("fingerprint fixture %s: %w", item.FixtureID, err)
		}
		if other, reused := usedFileSets[filesDigest]; reused {
			return fmt.Errorf("sealed fixture %s duplicates fixture content from the %s tranche", item.FixtureID, other)
		}
	}
	return nil
}

func digestFixtureFiles(files map[string]string) (string, error) {
	data, err := json.Marshal(files)
	if err != nil {
		return "", err
	}
	value, err := DecodeJSON(data)
	if err != nil {
		return "", err
	}
	return Digest(value)
}
