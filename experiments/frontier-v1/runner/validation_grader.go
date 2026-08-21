package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	frozenGraderDigest             = "sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc"
	frozenValidationScheduleDigest = "sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813"
	frozenValidationManifest       = "sha256:57ccc0c754f7c2beb74063dacc4bad26ab8b598ac2da9b6a8d637afd391fd490"
	frozenValidationRegistryRaw    = "sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816"
	frozenPrivateArchive           = "sha256:5a320a8742185e0471c6add885d1861950bbde8c7c312afbe907ae99762922e4"
	maximumCustodianOutput         = 4 * 1024 * 1024
	custodianCommandTimeout        = 30 * time.Second
)

type validationGraderDescription struct {
	V                    int    `json:"v"`
	Tranche              string `json:"tranche"`
	GraderDigest         string `json:"grader_digest"`
	ScheduleDigest       string `json:"schedule_digest"`
	ManifestDigest       string `json:"manifest_digest"`
	LabelRegistryRawHash string `json:"label_registry_raw_sha256"`
	PrivateArchiveDigest string `json:"private_archive_digest"`
	Cases                int    `json:"cases"`
}

type externalValidationGrader struct {
	executable  string
	description validationGraderDescription
	adapterHash string
}

type externalGradeResult struct {
	CaseID string                `json:"case_id"`
	Checks []externalCheckResult `json:"checks"`
	Status string                `json:"status"`
}

type externalCheckResult struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Verdict string `json:"verdict"`
	Reason  string `json:"reason"`
}

func prepareExternalValidationGrader(repositoryRoot, executable string) (externalValidationGrader, error) {
	if !filepath.IsAbs(executable) {
		return externalValidationGrader{}, fmt.Errorf("validation grader must be an absolute external executable")
	}
	if err := ensureOutside(repositoryRoot, executable); err != nil {
		return externalValidationGrader{}, err
	}
	info, err := os.Stat(executable)
	if err != nil || !info.Mode().IsRegular() {
		return externalValidationGrader{}, fmt.Errorf("validation grader must be a regular external executable")
	}
	contents, err := os.ReadFile(executable)
	if err != nil {
		return externalValidationGrader{}, err
	}
	grader := externalValidationGrader{
		executable:  executable,
		adapterHash: fmt.Sprintf("sha256:%x", sha256.Sum256(contents)),
	}
	output, err := runCustodianCommand(executable, "describe")
	if err != nil {
		return externalValidationGrader{}, err
	}
	if err := decodeStrictBytes(output, &grader.description); err != nil {
		return externalValidationGrader{}, fmt.Errorf("decode validation grader description: %w", err)
	}
	if err := validateGraderDescription(grader.description); err != nil {
		return externalValidationGrader{}, err
	}
	if actual, err := digestRawFile(executable); err != nil || actual != grader.adapterHash {
		return externalValidationGrader{}, fmt.Errorf("external validation grader identity changed during handshake")
	}
	return grader, nil
}

func validateGraderDescription(description validationGraderDescription) error {
	want := validationGraderDescription{
		V: 1, Tranche: "validation", GraderDigest: frozenGraderDigest,
		ScheduleDigest: frozenValidationScheduleDigest, ManifestDigest: frozenValidationManifest,
		LabelRegistryRawHash: frozenValidationRegistryRaw, PrivateArchiveDigest: frozenPrivateArchive,
		Cases: 120,
	}
	if description != want {
		return fmt.Errorf("validation grader description does not match the frozen custody contract")
	}
	return nil
}

func (grader externalValidationGrader) Grade(repositoryRoot, tranche, caseID, trialPath string) ([]byte, error) {
	if tranche != "validation" || !strings.HasPrefix(caseID, "validation_") {
		return nil, fmt.Errorf("external validation grader rejects tranche %q case %q", tranche, caseID)
	}
	resolvedTrial := resolveRepositoryPath(repositoryRoot, trialPath)
	if err := ensureInside(repositoryRoot, resolvedTrial); err != nil {
		return nil, err
	}
	if info, err := os.Stat(resolvedTrial); err != nil || !info.Mode().IsRegular() {
		return nil, fmt.Errorf("validation trial must be a regular repository evidence file")
	}
	if actual, err := digestRawFile(grader.executable); err != nil || actual != grader.adapterHash {
		return nil, fmt.Errorf("external validation grader identity changed after handshake")
	}
	output, err := runCustodianCommand(grader.executable, "grade", "--case-id", caseID, "--trial", resolvedTrial)
	if err != nil {
		return nil, err
	}
	var result externalGradeResult
	if err := decodeStrictBytes(output, &result); err != nil {
		return nil, fmt.Errorf("decode custody-safe grade: %w", err)
	}
	if result.CaseID != caseID {
		return nil, fmt.Errorf("custody-safe grade case_id %q does not match %q", result.CaseID, caseID)
	}
	if err := validateExternalGrade(result); err != nil {
		return nil, err
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return append(encoded, '\n'), nil
}

func validateExternalGrade(result externalGradeResult) error {
	if len(result.Checks) == 0 {
		return fmt.Errorf("custody-safe grade must contain at least one check")
	}
	wantStatus := "pass"
	seen := make(map[string]struct{}, len(result.Checks))
	for _, check := range result.Checks {
		if check.ID == "" || check.Kind == "" || check.Reason == "" {
			return fmt.Errorf("custody-safe grade contains an incomplete check")
		}
		if _, duplicate := seen[check.ID]; duplicate {
			return fmt.Errorf("custody-safe grade contains duplicate check %q", check.ID)
		}
		seen[check.ID] = struct{}{}
		switch check.Verdict {
		case "pass":
		case "fail":
			if wantStatus != "indeterminate" {
				wantStatus = "fail"
			}
		case "manual_required":
			wantStatus = "indeterminate"
		default:
			return fmt.Errorf("custody-safe grade contains unknown verdict %q", check.Verdict)
		}
	}
	if result.Status != wantStatus {
		return fmt.Errorf("custody-safe grade status %q is inconsistent with its checks", result.Status)
	}
	return nil
}

func runCustodianCommand(executable string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), custodianCommandTimeout)
	defer cancel()
	command := exec.CommandContext(ctx, executable, args...)
	var output cappedBuffer
	command.Stdout = &output
	command.Stderr = io.Discard
	err := command.Run()
	if ctx.Err() != nil {
		return nil, fmt.Errorf("external validation grader timed out")
	}
	if errors.Is(output.err, errCustodianOutputLimit) {
		return nil, output.err
	}
	if err != nil {
		return nil, fmt.Errorf("external validation grader failed")
	}
	return output.Bytes(), nil
}

var errCustodianOutputLimit = errors.New("external validation grader exceeded output limit")

type cappedBuffer struct {
	bytes.Buffer
	err error
}

func (buffer *cappedBuffer) Write(value []byte) (int, error) {
	if buffer.Len()+len(value) > maximumCustodianOutput {
		buffer.err = errCustodianOutputLimit
		return 0, buffer.err
	}
	return buffer.Buffer.Write(value)
}

func decodeStrictBytes(contents []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("trailing JSON data")
	}
	return nil
}

func ensureOutside(root, target string) error {
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return fmt.Errorf("resolve implementation repository boundary: %w", err)
	}
	resolvedTarget, err := filepath.EvalSymlinks(target)
	if err != nil {
		return fmt.Errorf("resolve validation grader boundary: %w", err)
	}
	relative, err := filepath.Rel(resolvedRoot, resolvedTarget)
	if err != nil {
		return fmt.Errorf("compare validation grader boundary: %w", err)
	}
	if relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) && !filepath.IsAbs(relative) {
		return fmt.Errorf("validation grader must remain outside the implementation repository")
	}
	return nil
}

func digestRawFile(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("sha256:%x", sha256.Sum256(contents)), nil
}
