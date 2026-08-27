package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const preValidationArtifactsPath = "experiments/frontier-v1/pre-validation-artifacts.json"

const (
	frozenValidationCostCapUSD = 0.15
	frozenValidationTimeout    = 180 * time.Second
)

type validationGateDocument struct {
	Status string `json:"status"`
	Gates  struct {
		MayOpenValidation bool   `json:"may_open_validation"`
		MayOpenHeldOut    bool   `json:"may_open_held_out"`
		ExecutionStatus   string `json:"validation_execution_status"`
		ClosureReview     string `json:"validation_execution_closure_review"`
		ClosureVerdict    string `json:"validation_execution_closure_review_verdict"`
	} `json:"gates"`
	Artifacts struct {
		Schedule struct {
			Path                 string `json:"path"`
			CanonicalDigest      string `json:"canonical_json_sha256"`
			SourceManifest       string `json:"source_manifest"`
			SourceManifestDigest string `json:"source_manifest_canonical_json_sha256"`
			Frozen               bool   `json:"frozen"`
		} `json:"validation_schedule"`
		WorldDefinitionAndWorldBuildDigest struct {
			Frozen             bool   `json:"frozen"`
			WorldBuildManifest string `json:"world_build_manifest"`
			WorldBuildDigest   string `json:"world_build_digest"`
		} `json:"world_definition_and_world_build_digest"`
	} `json:"artifacts"`
}

func normalizeTranche(value string) (string, error) {
	tranche := strings.ToLower(strings.TrimSpace(value))
	if tranche != "authoring" && tranche != "validation" {
		return "", fmt.Errorf("unsupported tranche %q", value)
	}
	return tranche, nil
}

func effectiveTranche(config runConfig) string {
	if config.tranche == "" {
		return "authoring"
	}
	return config.tranche
}

func requireValidationGate(repositoryRoot string) error {
	document, err := readValidationGateDocument(repositoryRoot)
	if err != nil {
		return err
	}
	if document.Gates.MayOpenHeldOut {
		return fmt.Errorf("validation execution requires the held-out gate to remain closed")
	}
	if !document.Gates.MayOpenValidation {
		return fmt.Errorf("validation gate is closed")
	}
	if err := requirePostClosureAudit(repositoryRoot, document); err != nil {
		return err
	}
	if document.Status != "complete" || !document.Artifacts.Schedule.Frozen ||
		document.Artifacts.Schedule.Path != "experiments/frontier-v1/schedules/validation.json" ||
		document.Artifacts.Schedule.CanonicalDigest != frozenValidationScheduleDigest ||
		document.Artifacts.Schedule.SourceManifest != "experiments/frontier-v1/manifests/validation.json" ||
		document.Artifacts.Schedule.SourceManifestDigest != frozenValidationManifest {
		return fmt.Errorf("validation gate does not name the frozen validation schedule")
	}
	return verifyValidationPublicIdentities(repositoryRoot, document)
}

// requirePostClosureAudit encodes the protocol-v5 section 9 post-closure gate.
// After any validation execution that ended without a completed schedule -
// interruption, budget stop, or lapse - no subsequent validation execution is
// authorized until the closure record has been independently reviewed and
// committed. Repeated close-and-retry would otherwise condition the
// eventually completed tranche on side signals, so each retry must survive an
// independent audit of why the last execution died.
func requirePostClosureAudit(repositoryRoot string, document validationGateDocument) error {
	status := document.Gates.ExecutionStatus
	if status == "" || status == "complete" {
		return nil
	}
	review := document.Gates.ClosureReview
	if review == "" || document.Gates.ClosureVerdict != "ACCEPT" {
		return fmt.Errorf("validation execution closed %s: another execution requires an independently reviewed closure record", status)
	}
	return verifyClosureReviewRecord(repositoryRoot, review)
}

// verifyClosureReviewRecord requires the named closure record to be a real,
// non-empty review document under docs/reviews that states the verdict the
// gate document claims for it. Existence alone would let any path in the
// repository satisfy the post-closure gate.
func verifyClosureReviewRecord(repositoryRoot, review string) error {
	if !strings.HasPrefix(review, "docs/reviews/") || !strings.HasSuffix(review, ".md") {
		return fmt.Errorf("validation closure review record must be a markdown review under docs/reviews")
	}
	path := filepath.Join(repositoryRoot, filepath.FromSlash(review))
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("validation closure review record: %w", err)
	}
	if !info.Mode().IsRegular() || info.Size() == 0 {
		return fmt.Errorf("validation closure review record is not a non-empty regular file")
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("validation closure review record: %w", err)
	}
	if !strings.Contains(string(contents), "ACCEPT") {
		return fmt.Errorf("validation closure review record does not state an ACCEPT verdict")
	}
	return nil
}

func requireFrozenWorldBuild(repositoryRoot, liveDigest string) error {
	document, err := readValidationGateDocument(repositoryRoot)
	if err != nil {
		return err
	}
	frozen := document.Artifacts.WorldDefinitionAndWorldBuildDigest
	if !frozen.Frozen || frozen.WorldBuildDigest == "" {
		return fmt.Errorf("validation requires a frozen world-build digest")
	}
	if strings.TrimSpace(liveDigest) == "" {
		return fmt.Errorf("validation requires a live world-build digest")
	}
	if liveDigest != frozen.WorldBuildDigest {
		return fmt.Errorf("live world-build %s does not match frozen world-build %s", liveDigest, frozen.WorldBuildDigest)
	}
	return nil
}

func requireFrozenValidationTrialLimits(budget string, timeout time.Duration) error {
	costCap, err := strconv.ParseFloat(budget, 64)
	if err != nil || costCap != frozenValidationCostCapUSD {
		return fmt.Errorf("validation requires --max-budget-usd exactly 0.15")
	}
	if timeout != frozenValidationTimeout {
		return fmt.Errorf("validation requires --timeout exactly 180s")
	}
	return nil
}

func readValidationGateDocument(repositoryRoot string) (validationGateDocument, error) {
	var document validationGateDocument
	contents, err := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(preValidationArtifactsPath)))
	if err != nil {
		return validationGateDocument{}, fmt.Errorf("read validation gate: %w", err)
	}
	if err := json.Unmarshal(contents, &document); err != nil {
		return validationGateDocument{}, fmt.Errorf("read validation gate: %w", err)
	}
	return document, nil
}

func verifyValidationPublicIdentities(repositoryRoot string, document validationGateDocument) error {
	if actual, err := digestJSONFile(filepath.Join(repositoryRoot, filepath.FromSlash(document.Artifacts.Schedule.Path))); err != nil || actual != frozenValidationScheduleDigest {
		return fmt.Errorf("frozen validation schedule identity check failed")
	}
	if actual, err := digestJSONFile(filepath.Join(repositoryRoot, filepath.FromSlash(document.Artifacts.Schedule.SourceManifest))); err != nil || actual != frozenValidationManifest {
		return fmt.Errorf("frozen validation manifest identity check failed")
	}
	registry := filepath.Join(repositoryRoot, "experiments", "frontier-v1", "manifests", "validation-label-digests.json")
	if actual, err := digestRawFile(registry); err != nil || actual != frozenValidationRegistryRaw {
		return fmt.Errorf("frozen validation label-registry identity check failed")
	}
	return nil
}
