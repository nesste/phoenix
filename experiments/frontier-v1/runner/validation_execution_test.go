package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestFrozenValidationScheduleLoadsOnlyPublicInputs(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	ids, err := caseIDsForTranche(root, "validation")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 120 {
		t.Fatalf("validation case count = %d", len(ids))
	}
	cases, err := loadSelectedCasesForTranche(root, "validation", ids)
	if err != nil {
		t.Fatal(err)
	}
	var schedule launchSchedule
	schedulePath := filepath.Join(root, "experiments", "frontier-v1", "schedules", "validation.json")
	if err := decodeStrict(schedulePath, &schedule); err != nil {
		t.Fatal(err)
	}
	if err := validateScheduleForTranche(schedule, cases, "validation"); err != nil {
		t.Fatal(err)
	}
	if digest, err := digestJSONFile(schedulePath); err != nil || digest != frozenValidationScheduleDigest {
		t.Fatalf("validation schedule digest = %s, %v", digest, err)
	}
	for _, item := range cases {
		if _, err := loadFixtureForTranche(root, "validation", item); err != nil {
			t.Fatalf("public fixture for %s: %v", item.CaseID, err)
		}
	}
}

func TestRepositoryValidationGateRemainsClosed(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	if err := requireValidationGate(root); err == nil || !strings.Contains(err.Error(), "gate is closed") {
		t.Fatalf("gate error = %v", err)
	}
}

func TestClosedValidationGatePrecedesOutputMutation(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	output := filepath.ToSlash(filepath.Join("experiments", "frontier-v1", "results", "must-not-exist-gate-test"))
	resolved := filepath.Join(root, filepath.FromSlash(output))
	_ = os.Remove(resolved)
	_, err := prepareScheduledCLI(root, "validation", output, "", "", "", "0.15", "300", 1, nil)
	if err == nil || !strings.Contains(err.Error(), "gate is closed") {
		t.Fatalf("gate error = %v", err)
	}
	if _, statErr := os.Stat(resolved); !os.IsNotExist(statErr) {
		t.Fatalf("closed gate mutated output: %v", statErr)
	}
}

func TestOpenValidationGateRequiresExactPublicIdentities(t *testing.T) {
	sourceRoot := filepath.Join("..", "..", "..")
	root := openValidationGateTestRoot(t, sourceRoot)
	if err := requireValidationGate(root); err != nil {
		t.Fatal(err)
	}
	registry := filepath.Join(root, "experiments", "frontier-v1", "manifests", "validation-label-digests.json")
	if err := os.WriteFile(registry, []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := requireValidationGate(root); err == nil || !strings.Contains(err.Error(), "label-registry") {
		t.Fatalf("registry identity error = %v", err)
	}
}

func TestValidationRejectsChangedTrialLimitsBeforeOutputOrCustodianContact(t *testing.T) {
	sourceRoot := filepath.Join("..", "..", "..")
	root := openValidationGateTestRoot(t, sourceRoot)
	marker := filepath.Join(t.TempDir(), "custodian-contacted")
	t.Setenv("CUSTODIAN_CONTACT_MARKER", marker)
	custodian := buildCustodianHelper(t)

	tests := []struct {
		name    string
		budget  string
		timeout time.Duration
		message string
	}{
		{name: "lower cost cap", budget: "0.14", timeout: 180 * time.Second, message: "--max-budget-usd exactly 0.15"},
		{name: "higher cost cap", budget: "0.16", timeout: 180 * time.Second, message: "--max-budget-usd exactly 0.15"},
		{name: "whitespace-padded cost cap", budget: " 0.15 ", timeout: 180 * time.Second, message: "--max-budget-usd exactly 0.15"},
		{name: "shorter timeout", budget: "0.15", timeout: 179 * time.Second, message: "--timeout exactly 180s"},
		{name: "longer timeout", budget: "0.15", timeout: 181 * time.Second, message: "--timeout exactly 180s"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := filepath.ToSlash(filepath.Join("results", strings.ReplaceAll(test.name, " ", "-")))
			_, err := prepareScheduledCLI(
				root, "validation", output, "experiments/frontier-v1/schedules/validation.json", "",
				custodian, test.budget, "300", test.timeout, nil,
			)
			if err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("trial-limit error = %v", err)
			}
			if _, statErr := os.Stat(filepath.Join(root, filepath.FromSlash(output))); !os.IsNotExist(statErr) {
				t.Fatalf("invalid trial limit mutated output: %v", statErr)
			}
			if _, statErr := os.Stat(marker); !os.IsNotExist(statErr) {
				t.Fatalf("invalid trial limit contacted custodian: %v", statErr)
			}
		})
	}
}

func TestValidationGraderMustRemainExternal(t *testing.T) {
	root := t.TempDir()
	executable := filepath.Join(root, "grader")
	if err := os.WriteFile(executable, []byte("not executable"), 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := prepareExternalValidationGrader(root, executable); err == nil || !strings.Contains(err.Error(), "outside") {
		t.Fatalf("boundary error = %v", err)
	}
}

func TestExternalValidationGraderHandshakeAndGrade(t *testing.T) {
	repository := t.TempDir()
	executable := buildCustodianHelper(t)
	grader, err := prepareExternalValidationGrader(repository, executable)
	if err != nil {
		t.Fatal(err)
	}
	trialPath := filepath.Join(repository, "evidence", "validation_example.trial.json")
	if err := os.MkdirAll(filepath.Dir(trialPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(trialPath, []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	grade, err := grader.Grade(repository, "validation", "validation_example", trialPath)
	if err != nil {
		t.Fatal(err)
	}
	var result externalGradeResult
	if err := json.Unmarshal(grade, &result); err != nil {
		t.Fatal(err)
	}
	if result.CaseID != "validation_example" || result.Status != "pass" || len(result.Checks) != 1 {
		t.Fatalf("grade = %#v", result)
	}
	if strings.Contains(string(grade), repository) || strings.Contains(string(grade), "labels") {
		t.Fatalf("grade leaks implementation path or label location: %s", grade)
	}
}

func TestExternalGradeConsistencyIsEnforced(t *testing.T) {
	result := externalGradeResult{
		CaseID: "validation_example", Status: "pass",
		Checks: []externalCheckResult{{ID: "check", Kind: "test", Verdict: "fail", Reason: "synthetic"}},
	}
	if err := validateExternalGrade(result); err == nil || !strings.Contains(err.Error(), "inconsistent") {
		t.Fatalf("consistency error = %v", err)
	}
}

func buildCustodianHelper(t *testing.T) string {
	t.Helper()
	directory := t.TempDir()
	source := filepath.Join(directory, "main.go")
	program := `package main
import ("encoding/json"; "fmt"; "os")
func main() {
 if marker := os.Getenv("CUSTODIAN_CONTACT_MARKER"); marker != "" { os.WriteFile(marker, []byte("contacted"), 0600) }
 if len(os.Args) == 2 && os.Args[1] == "describe" {
  fmt.Print(` + "`" + `{"v":1,"tranche":"validation","grader_digest":"` + frozenGraderDigest + `","schedule_digest":"` + frozenValidationScheduleDigest + `","manifest_digest":"` + frozenValidationManifest + `","label_registry_raw_sha256":"` + frozenValidationRegistryRaw + `","private_archive_digest":"` + frozenPrivateArchive + `","cases":120}` + "`" + `); return
 }
 if len(os.Args) == 6 && os.Args[1] == "grade" && os.Args[2] == "--case-id" && os.Args[4] == "--trial" {
  json.NewEncoder(os.Stdout).Encode(map[string]any{"case_id":os.Args[3],"checks":[]any{map[string]any{"id":"synthetic","kind":"test","verdict":"pass","reason":"synthetic pass"}},"status":"pass"}); return
 }
 os.Exit(2)
}`
	if err := os.WriteFile(source, []byte(program), 0o600); err != nil {
		t.Fatal(err)
	}
	name := "custodian"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	executable := filepath.Join(directory, name)
	command := exec.Command("go", "build", "-o", executable, source)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build custodian helper: %v: %s", err, output)
	}
	return executable
}

func openValidationGateTestRoot(t *testing.T, sourceRoot string) string {
	t.Helper()
	root := t.TempDir()
	for _, relative := range []string{
		"experiments/frontier-v1/schedules/validation.json",
		"experiments/frontier-v1/manifests/validation.json",
		"experiments/frontier-v1/manifests/validation-label-digests.json",
	} {
		copyTestFile(t, filepath.Join(sourceRoot, filepath.FromSlash(relative)), filepath.Join(root, filepath.FromSlash(relative)))
	}
	gate := map[string]any{
		"status": "complete",
		"gates":  map[string]any{"may_open_validation": true, "may_open_held_out": false},
		"artifacts": map[string]any{"validation_schedule": map[string]any{
			"path": "experiments/frontier-v1/schedules/validation.json", "canonical_json_sha256": frozenValidationScheduleDigest,
			"source_manifest": "experiments/frontier-v1/manifests/validation.json", "source_manifest_canonical_json_sha256": frozenValidationManifest,
			"frozen": true,
		}},
	}
	if err := writeJSON(filepath.Join(root, filepath.FromSlash(preValidationArtifactsPath)), gate); err != nil {
		t.Fatal(err)
	}
	return root
}

func copyTestFile(t *testing.T, source, target string) {
	t.Helper()
	contents, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, contents, 0o600); err != nil {
		t.Fatal(err)
	}
}
