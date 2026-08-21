// Command scheduletool creates and verifies the outcome-free Phase 1
// validation schedule. It reads only public sealed-corpus artifacts and never
// executes a trial, opens a label, or changes an outcome gate.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/gowebpki/jcs"
)

const (
	validationCaseCount   = 120
	validationFamilyCount = 24
)

type sealedManifest struct {
	V       int            `json:"v"`
	Tranche string         `json:"tranche"`
	Sealed  bool           `json:"sealed"`
	Counts  manifestCounts `json:"counts"`
	Cases   []manifestCase `json:"cases"`
}

type manifestCounts struct {
	Cases    int `json:"cases"`
	Families int `json:"families"`
}

type manifestCase struct {
	CaseID         string `json:"case_id"`
	FamilyID       string `json:"family_id"`
	InputDigest    string `json:"input_digest"`
	SandboxFixture string `json:"sandbox_fixture"`
}

type publicCase struct {
	CaseID         string `json:"case_id"`
	FamilyID       string `json:"family_id"`
	SandboxFixture string `json:"sandbox_fixture"`
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("validation-schedule", flag.ContinueOnError)
	flags.SetOutput(stderr)
	repositoryRoot := flags.String("repo-root", ".", "Phoenix repository root")
	writePath := flags.String("write", "", "write the deterministic validation schedule and exit")
	verifyPath := flags.String("verify", "", "verify an existing deterministic validation schedule and exit")
	if err := flags.Parse(args); err != nil || flags.NArg() != 0 {
		return 2
	}
	if (*writePath == "") == (*verifyPath == "") {
		fmt.Fprintln(stderr, "exactly one of --write or --verify is required")
		return 2
	}
	root, err := filepath.Abs(*repositoryRoot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	cases, manifestDigest, err := loadValidationCases(root)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	expected, err := generateSchedule("validation", cases, phase1ScheduleSeed, phase1Repetitions, phase1Arms)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	requestedPath := *writePath
	if requestedPath == "" {
		requestedPath = *verifyPath
	}
	resolved := resolveRepositoryPath(root, requestedPath)
	if err := ensureInside(root, resolved); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if *writePath != "" {
		if _, err := os.Stat(resolved); err == nil {
			fmt.Fprintln(stderr, "schedule path already exists")
			return 1
		} else if !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if err := writeJSON(resolved, expected); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	} else {
		var actual launchSchedule
		if err := decodeStrict(resolved, &actual); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if !reflect.DeepEqual(actual, expected) {
			fmt.Fprintln(stderr, "schedule does not match the deterministic validation construction")
			return 1
		}
	}

	scheduleDigest, err := digestJSONFile(resolved)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "schedule: %s (%d launches); manifest: %s\n", scheduleDigest, len(expected.Entries), manifestDigest)
	return 0
}

func loadValidationCases(root string) ([]scheduleCase, string, error) {
	manifestPath := filepath.Join(root, "experiments", "frontier-v1", "manifests", "validation.json")
	var manifest sealedManifest
	if err := decodeJSON(manifestPath, &manifest); err != nil {
		return nil, "", err
	}
	if manifest.V != 1 || manifest.Tranche != "validation" || !manifest.Sealed {
		return nil, "", fmt.Errorf("validation manifest must be sealed version 1")
	}
	if manifest.Counts.Cases != validationCaseCount || manifest.Counts.Families != validationFamilyCount || len(manifest.Cases) != validationCaseCount {
		return nil, "", fmt.Errorf("validation manifest counts must be %d cases in %d families", validationCaseCount, validationFamilyCount)
	}

	caseDirectory := filepath.Join(root, "experiments", "frontier-v1", "corpus", "validation")
	paths, err := filepath.Glob(filepath.Join(caseDirectory, "validation_*.json"))
	if err != nil {
		return nil, "", err
	}
	if len(paths) != validationCaseCount {
		return nil, "", fmt.Errorf("validation corpus has %d case files, want %d", len(paths), validationCaseCount)
	}

	seenCases := map[string]struct{}{}
	familyCounts := map[string]int{}
	cases := make([]scheduleCase, 0, len(manifest.Cases))
	for _, row := range manifest.Cases {
		if !strings.HasPrefix(row.CaseID, "validation_") || !strings.HasPrefix(row.FamilyID, "validation_family_") {
			return nil, "", fmt.Errorf("manifest contains non-validation case or family")
		}
		if _, duplicate := seenCases[row.CaseID]; duplicate {
			return nil, "", fmt.Errorf("manifest case %q is duplicated", row.CaseID)
		}
		seenCases[row.CaseID] = struct{}{}
		path := filepath.Join(caseDirectory, row.CaseID+".json")
		digest, err := digestJSONFile(path)
		if err != nil {
			return nil, "", err
		}
		if digest != row.InputDigest {
			return nil, "", fmt.Errorf("case %s digest %s does not match manifest %s", row.CaseID, digest, row.InputDigest)
		}
		var item publicCase
		if err := decodeJSON(path, &item); err != nil {
			return nil, "", err
		}
		if item.CaseID != row.CaseID || item.FamilyID != row.FamilyID || item.SandboxFixture != row.SandboxFixture {
			return nil, "", fmt.Errorf("case %s does not match its manifest row", row.CaseID)
		}
		familyCounts[row.FamilyID]++
		cases = append(cases, scheduleCase{CaseID: row.CaseID, FamilyID: row.FamilyID})
	}
	if len(familyCounts) != validationFamilyCount {
		return nil, "", fmt.Errorf("validation corpus has %d families, want %d", len(familyCounts), validationFamilyCount)
	}
	for familyID, count := range familyCounts {
		if count != 5 {
			return nil, "", fmt.Errorf("family %s has %d cases, want 5", familyID, count)
		}
	}
	sort.Slice(cases, func(left, right int) bool { return cases[left].CaseID < cases[right].CaseID })
	manifestDigest, err := digestJSONFile(manifestPath)
	return cases, manifestDigest, err
}

func decodeJSON(path string, target any) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(contents))
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("decode %s: trailing JSON data", path)
	}
	return nil
}

func decodeStrict(path string, target any) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("decode %s: trailing JSON data", path)
	}
	return nil
}

func writeJSON(path string, value any) error {
	contents, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(contents, '\n'), 0o644)
}

func digestJSONFile(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	canonical, err := jcs.Transform(contents)
	if err != nil {
		return "", fmt.Errorf("canonicalize %s: %w", path, err)
	}
	sum := sha256.Sum256(canonical)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func resolveRepositoryPath(root, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, filepath.FromSlash(path))
}

func ensureInside(root, target string) error {
	relative, err := filepath.Rel(root, target)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
		return fmt.Errorf("schedule path must stay inside the repository")
	}
	return nil
}
