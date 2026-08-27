// Command runner executes an approved frontier-v1 tranche against a fresh
// sandbox, Phoenix process, and pinned runtime session for every case.
package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/nesste/phoenix/internal/buildmanifest"
	"github.com/nesste/phoenix/internal/episode"
	"github.com/nesste/phoenix/internal/surface"
	"github.com/nesste/phoenix/internal/world"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	flags := flag.NewFlagSet("frontier-v1-runner", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	repositoryRoot := flags.String("repo-root", ".", "Phoenix repository root")
	tranche := flags.String("tranche", "authoring", "experiment tranche: authoring or validation")
	caseID := flags.String("case", "all", "case id or all")
	outputDir := flags.String("output-dir", "", "repository-relative result directory; defaults to retained C evidence or an arm-specific probe directory")
	arm := flags.String("arm", "C", "experiment arm: A, B, C, or D")
	armBDocument := flags.String("arm-b-document", "", "frozen Arm B static document; required for Arm B")
	writeSchedulePath := flags.String("write-schedule", "", "write a deterministic authoring A-D schedule and exit")
	schedulePath := flags.String("schedule", "", "execute a previously written A-D schedule")
	validationGrader := flags.String("validation-grader", "", "absolute path to the external validation custodian grader")
	runtimePath := flags.String("runtime", "claude", "pinned runtime executable")
	budget := flags.String("max-budget-usd", "0.15", "maximum model cost per case")
	runBudget := flags.String("run-budget-usd", "75", "hard scheduled-run budget in USD; checked at pairing-key boundaries")
	timeout := flags.Duration("timeout", 180*time.Second, "per-case timeout")
	resumeAttestation := flags.String("resume-attestation", "", "custodian attestation authorizing the single permitted resume of an interrupted validation tranche")
	if err := flags.Parse(args); err != nil || flags.NArg() != 0 {
		return 2
	}
	root, err := filepath.Abs(*repositoryRoot)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	selectedTranche, err := normalizeTranche(*tranche)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	ids, err := selectedCaseIDsForTranche(root, selectedTranche, *caseID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if *writeSchedulePath != "" {
		if selectedTranche != "authoring" {
			fmt.Fprintln(os.Stderr, "validation schedule generation is frozen and disabled")
			return 2
		}
		if *schedulePath != "" {
			fmt.Fprintln(os.Stderr, "--write-schedule and --schedule are mutually exclusive")
			return 2
		}
		return writeScheduleCLI(root, *writeSchedulePath, ids)
	}
	if *schedulePath != "" {
		return runScheduleCLI(root, selectedTranche, *outputDir, *schedulePath, *armBDocument,
			*validationGrader, *runtimePath, *budget, *runBudget, *resumeAttestation, *timeout, ids)
	}
	if *resumeAttestation != "" {
		fmt.Fprintln(os.Stderr, "--resume-attestation applies to scheduled execution only")
		return 2
	}
	if selectedTranche != "authoring" {
		fmt.Fprintln(os.Stderr, "validation permits scheduled execution only")
		return 2
	}
	return runProbeCLI(root, *outputDir, *arm, *armBDocument, *runtimePath, *budget, *timeout, ids)
}

func runProbeCLI(root, output, arm, armBDocument, runtimePath, budget string, timeout time.Duration, ids []string) int {
	config, cleanup, err := prepareRunConfig(root, output, arm, armBDocument, "authoring", timeout, budget)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer cleanup()
	driver := claudeDriver{executable: runtimePath}
	if err := driver.Verify(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	summary, err := runCases(config, ids, driver, corpusGrader{goExecutable: "go"})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err := writeJSON(filepath.Join(config.outputDir, "summary.json"), summary); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Printf("authoring: %d passed, %d failed\n", summary.Passed, summary.Failed)
	if summary.Failed != 0 {
		return 1
	}
	return 0
}

func runScheduleCLI(
	root, tranche, output, schedulePath, armBDocument, validationGrader, runtimePath, budget, runBudget, resumeAttestation string,
	timeout time.Duration,
	ids []string,
) int {
	prepared, err := prepareScheduledCLI(root, tranche, output, schedulePath, armBDocument,
		validationGrader, budget, runBudget, resumeAttestation, timeout, ids)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer prepared.cleanup()
	driver := claudeDriver{executable: runtimePath}
	if err := driver.Verify(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	stopWatch := watchProcessTermination(
		newProcessEventLog(prepared.config.outputDir, tranche, prepared.scheduleDigest, prepared.config.worldBuild),
		func() processEventProgress { return processEventProgress{} },
	)
	defer stopWatch()
	summary, err := runScheduledCases(
		prepared.config, prepared.schedule, prepared.scheduleDigest, prepared.runBudgetUSD,
		driver, prepared.grader,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err := writeJSON(filepath.Join(prepared.config.outputDir, "scheduled-summary.json"), summary); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Printf("scheduled %s: %d pass, %d fail, %d unresolved, %d budget-stopped\n", tranche, summary.Passes, summary.Failures, summary.Unresolved, summary.BudgetStopped)
	if summary.Status == "indeterminate" {
		return 1
	}
	return 0
}

type preparedScheduledCLI struct {
	config         runConfig
	schedule       launchSchedule
	scheduleDigest string
	runBudgetUSD   float64
	grader         gradeDriver
	cleanup        func()
}

type preparedScheduledInputs struct {
	cases             []runnableCase
	schedule          launchSchedule
	scheduleDigest    string
	runBudgetUSD      float64
	grader            gradeDriver
	graderDigest      string
	graderBoundary    *validationGraderDescription
	graderAdapterHash string
}

func prepareScheduledCLI(
	root, tranche, output, schedulePath, armBDocument, validationGrader, budget, runBudget, resumeAttestation string,
	timeout time.Duration,
	ids []string,
) (preparedScheduledCLI, error) {
	if tranche == "validation" {
		if err := requireValidationGate(root); err != nil {
			return preparedScheduledCLI{}, err
		}
		if err := requireFrozenValidationTrialLimits(budget, timeout); err != nil {
			return preparedScheduledCLI{}, err
		}
	}
	if output == "" {
		output = filepath.ToSlash(filepath.Join("experiments", "frontier-v1", "results", "scheduled-"+tranche))
	}
	inputs, err := prepareScheduledInputs(root, tranche, schedulePath, validationGrader, runBudget, ids)
	if err != nil {
		return preparedScheduledCLI{}, err
	}
	preflightConfig := runConfig{
		repositoryRoot: root,
		worldPath:      filepath.Join(root, "worlds", "dev-repo", "world.json"),
	}
	if err := preflightSelectedCasesForTranche(preflightConfig, tranche, inputs.cases); err != nil {
		return preparedScheduledCLI{}, err
	}
	if err := requireScheduledOutputReady(root, output, inputs.schedule, inputs.scheduleDigest); err != nil {
		return preparedScheduledCLI{}, err
	}
	config, cleanup, err := prepareRunConfig(root, output, "B", armBDocument, tranche, timeout, budget)
	if err != nil {
		return preparedScheduledCLI{}, err
	}
	config.graderDigest = inputs.graderDigest
	config.graderBoundary = inputs.graderBoundary
	config.graderAdapterHash = inputs.graderAdapterHash
	if resumeAttestation != "" {
		config.resumeAttestation = resolveRepositoryPath(root, resumeAttestation)
	}
	if tranche == "validation" {
		if err := requireFrozenWorldBuild(root, config.worldBuild); err != nil {
			cleanup()
			return preparedScheduledCLI{}, err
		}
	}
	return preparedScheduledCLI{
		config: config, schedule: inputs.schedule, scheduleDigest: inputs.scheduleDigest,
		runBudgetUSD: inputs.runBudgetUSD, grader: inputs.grader, cleanup: cleanup,
	}, nil
}

func prepareScheduledInputs(root, tranche, schedulePath, validationGrader, runBudget string, ids []string) (preparedScheduledInputs, error) {
	cases, err := loadSelectedCasesForTranche(root, tranche, ids)
	if err != nil {
		return preparedScheduledInputs{}, err
	}
	resolvedSchedule := resolveRepositoryPath(root, schedulePath)
	if err := ensureInside(root, resolvedSchedule); err != nil {
		return preparedScheduledInputs{}, err
	}
	var schedule launchSchedule
	if err := decodeStrict(resolvedSchedule, &schedule); err != nil {
		return preparedScheduledInputs{}, err
	}
	if err := validateScheduleForTranche(schedule, cases, tranche); err != nil {
		return preparedScheduledInputs{}, err
	}
	digest, err := digestJSONFile(resolvedSchedule)
	if err != nil {
		return preparedScheduledInputs{}, err
	}
	totalBudget, err := strconv.ParseFloat(runBudget, 64)
	if err != nil || totalBudget <= 0 {
		return preparedScheduledInputs{}, fmt.Errorf("--run-budget-usd must be a positive decimal value")
	}
	inputs := preparedScheduledInputs{
		cases: cases, schedule: schedule, scheduleDigest: digest, runBudgetUSD: totalBudget,
		grader: corpusGrader{goExecutable: "go"},
	}
	if tranche == "validation" {
		if digest != frozenValidationScheduleDigest {
			return preparedScheduledInputs{}, fmt.Errorf("validation schedule digest %s does not match frozen custody contract", digest)
		}
		if totalBudget != 300 {
			return preparedScheduledInputs{}, fmt.Errorf("validation requires the frozen 300 USD run budget")
		}
		external, err := prepareExternalValidationGrader(root, validationGrader)
		if err != nil {
			return preparedScheduledInputs{}, err
		}
		inputs.grader = external
		inputs.graderDigest = frozenGraderDigest
		inputs.graderBoundary = &external.description
		inputs.graderAdapterHash = external.adapterHash
		return inputs, nil
	}
	inputs.graderDigest, err = verifiedGraderDigest(root, ids)
	if err != nil {
		return preparedScheduledInputs{}, err
	}
	return inputs, nil
}

func writeScheduleCLI(root, path string, ids []string) int {
	cases, err := loadSelectedCases(root, ids)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	schedule, err := generatePhase1Schedule(cases)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	resolved := resolveRepositoryPath(root, path)
	if err := ensureInside(root, resolved); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if _, err := os.Stat(resolved); err == nil {
		fmt.Fprintln(os.Stderr, "schedule path already exists")
		return 1
	} else if !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err := writeJSON(resolved, schedule); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	digest, err := digestJSONFile(resolved)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Printf("schedule: %s (%d launches)\n", digest, len(schedule.Entries))
	return 0
}

func selectedCaseIDsForTranche(root, tranche, selection string) ([]string, error) {
	if selection == "all" {
		return caseIDsForTranche(root, tranche)
	}
	if tranche == "validation" {
		return nil, fmt.Errorf("validation execution requires --case all")
	}
	if _, err := loadCaseForTranche(root, tranche, selection); err != nil {
		return nil, err
	}
	return []string{selection}, nil
}

func loadSelectedCases(root string, ids []string) ([]runnableCase, error) {
	return loadSelectedCasesForTranche(root, "authoring", ids)
}

func loadSelectedCasesForTranche(root, tranche string, ids []string) ([]runnableCase, error) {
	cases := make([]runnableCase, 0, len(ids))
	for _, id := range ids {
		item, err := loadCaseForTranche(root, tranche, id)
		if err != nil {
			return nil, err
		}
		cases = append(cases, item)
	}
	return cases, nil
}

func preflightSelectedCasesForTranche(config runConfig, tranche string, cases []runnableCase) error {
	worldRef, err := digestJSONFile(config.worldPath)
	if err != nil {
		return err
	}
	for _, item := range cases {
		if item.WorldRef != worldRef {
			return fmt.Errorf("case %s world_ref %s does not match production world %s", item.CaseID, item.WorldRef, worldRef)
		}
		if _, err := loadFixtureForTranche(config.repositoryRoot, tranche, item); err != nil {
			return fmt.Errorf("case %s fixture: %w", item.CaseID, err)
		}
	}
	return nil
}

func graderDigestForCases(root string, ids []string) (string, error) {
	graderDigest := ""
	for _, id := range ids {
		path := filepath.Join(root, "experiments", "frontier-v1", "labels", "authoring", id+".json")
		contents, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		var label struct {
			GradingScript string `json:"grading_script"`
		}
		if err := json.Unmarshal(contents, &label); err != nil {
			return "", err
		}
		if label.GradingScript == "" {
			return "", fmt.Errorf("label %s has no grading_script digest", id)
		}
		if graderDigest != "" && graderDigest != label.GradingScript {
			return "", fmt.Errorf("selected authoring labels do not share one grader digest")
		}
		graderDigest = label.GradingScript
	}
	return graderDigest, nil
}

func verifiedGraderDigest(root string, ids []string) (string, error) {
	declared, err := graderDigestForCases(root, ids)
	if err != nil {
		return "", err
	}
	command := exec.Command("go", "run", "./cmd/corpusctl", "grader-digest", "--repo-root", "../../..")
	command.Dir = filepath.Join(root, "experiments", "frontier-v1", "corpusctl")
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("compute grader digest: %w: %s", err, strings.TrimSpace(string(output)))
	}
	actual := strings.TrimSpace(string(output))
	if actual != declared {
		return "", fmt.Errorf("current grader digest %s does not match selected labels %s", actual, declared)
	}
	return actual, nil
}

func resolveRepositoryPath(root, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, filepath.FromSlash(path))
}

func requireEmptyScheduledOutput(root, output string) error {
	path := resolveRepositoryPath(root, output)
	if err := ensureInside(root, path); err != nil {
		return err
	}
	entries, err := os.ReadDir(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(entries) != 0 {
		return fmt.Errorf("scheduled output directory must not already contain evidence")
	}
	return nil
}

func prepareRunConfig(root, output, arm, armBDocument, tranche string, timeout time.Duration, budget string) (runConfig, func(), error) {
	if timeout <= 0 || strings.TrimSpace(budget) == "" {
		return runConfig{}, func() {}, fmt.Errorf("positive timeout and budget are required")
	}
	arm, err := normalizeArm(arm)
	if err != nil {
		return runConfig{}, func() {}, err
	}
	armBDocument, err = resolveArmBDocument(arm, armBDocument)
	if err != nil {
		return runConfig{}, func() {}, err
	}
	armBDocumentDigest := ""
	if arm == "B" {
		armBDocumentDigest, err = digestLFNormalizedFile(armBDocument)
		if err != nil {
			return runConfig{}, func() {}, err
		}
	}
	outputPath := resolveOutputPath(root, output, arm)
	if err := ensureInside(root, outputPath); err != nil {
		return runConfig{}, func() {}, err
	}
	if err := os.MkdirAll(outputPath, 0o755); err != nil {
		return runConfig{}, func() {}, err
	}
	phoenixPath, build, cleanup, err := buildPhoenix(root, phoenixBuildRecipeForTranche(tranche))
	if err != nil {
		return runConfig{}, func() {}, err
	}
	worldPath := filepath.Join(root, "worlds", "dev-repo", "world.json")
	schemaPath := filepath.Join(root, "spec", "world.schema.json")
	flatToolNames, err := loadFlatToolNames(arm, schemaPath, worldPath)
	if err != nil {
		cleanup()
		return runConfig{}, func() {}, err
	}
	return runConfig{
		arm: arm, armBDocument: armBDocument, armBDocumentDigest: armBDocumentDigest, flatToolNames: flatToolNames,
		tranche:        tranche,
		repositoryRoot: root, outputDir: outputPath,
		worldPath:   worldPath,
		schemaPath:  schemaPath,
		phoenixPath: phoenixPath, worldBuild: build,
		timeout: timeout, budgetUSD: budget,
	}, cleanup, nil
}

func resolveArmBDocument(arm, path string) (string, error) {
	if arm != "B" {
		return path, nil
	}
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("arm B requires --arm-b-document")
	}
	resolved, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve arm B document: %w", err)
	}
	info, err := os.Stat(resolved)
	if err != nil || !info.Mode().IsRegular() {
		return "", fmt.Errorf("arm B document must be a regular file")
	}
	return resolved, nil
}

func digestLFNormalizedFile(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	normalized := strings.ReplaceAll(strings.ReplaceAll(string(contents), "\r\n", "\n"), "\r", "\n")
	digest := sha256.Sum256([]byte(normalized))
	return "sha256:" + hex.EncodeToString(digest[:]), nil
}

func resolveOutputPath(root, output, arm string) string {
	if output == "" {
		output = "experiments/frontier-v1/results/authoring"
		if arm != "C" {
			output = filepath.ToSlash(filepath.Join("experiments", "frontier-v1", "results", "arm-probes", arm))
		}
	}
	if filepath.IsAbs(output) {
		return output
	}
	return filepath.Join(root, filepath.FromSlash(output))
}

func loadFlatToolNames(arm, schemaPath, worldPath string) ([]string, error) {
	if arm != "A" && arm != "B" {
		return nil, nil
	}
	definition, err := world.Load(schemaPath, worldPath)
	if err != nil {
		return nil, err
	}
	return surface.FlatToolNames(definition)
}

func normalizeArm(arm string) (string, error) {
	arm = strings.ToUpper(strings.TrimSpace(arm))
	switch arm {
	case "A", "B", "C", "D":
		return arm, nil
	default:
		return "", fmt.Errorf("unsupported experiment arm %q", arm)
	}
}

// phoenixBuildRecipe fixes the Phoenix build identity for one tranche. The
// validation recipe must reproduce the frozen world-build target in
// pre-validation-artifacts.json; authoring keeps a host-runnable binary.
type phoenixBuildRecipe struct {
	goos    string
	goarch  string
	version string
	// output is the repository-relative executable path recorded in the
	// world-build manifest; empty selects a fresh temporary directory.
	output string
}

func phoenixBuildRecipeForTranche(tranche string) phoenixBuildRecipe {
	if tranche == "validation" {
		return phoenixBuildRecipe{goos: "linux", goarch: "amd64", version: "dev", output: "bin/phoenix"}
	}
	return phoenixBuildRecipe{goos: runtime.GOOS, goarch: runtime.GOARCH, version: "authoring"}
}

func phoenixExecutablePath(root string, recipe phoenixBuildRecipe) (string, func(), error) {
	if recipe.output != "" {
		executable := filepath.Join(root, filepath.FromSlash(recipe.output))
		if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
			return "", nil, err
		}
		return executable, func() { _ = os.Remove(executable) }, nil
	}
	tempRoot := filepath.Join(root, "build")
	if err := os.MkdirAll(tempRoot, 0o755); err != nil {
		return "", nil, err
	}
	dir, err := os.MkdirTemp(tempRoot, "authoring-runner-")
	if err != nil {
		return "", nil, err
	}
	name := "phoenix"
	if recipe.goos == "windows" {
		name += ".exe"
	}
	return filepath.Join(dir, name), func() { _ = os.RemoveAll(dir) }, nil
}

func buildPhoenix(root string, recipe phoenixBuildRecipe) (string, string, func(), error) {
	executable, cleanup, err := phoenixExecutablePath(root, recipe)
	if err != nil {
		return "", "", func() {}, err
	}
	command := phoenixBuildCommand(root, executable, recipe)
	if output, err := command.CombinedOutput(); err != nil {
		cleanup()
		return "", "", func() {}, fmt.Errorf("build Phoenix runner executable: %w: %s", err, strings.TrimSpace(string(output)))
	}
	relative, err := filepath.Rel(root, executable)
	if err != nil {
		cleanup()
		return "", "", func() {}, err
	}
	manifest, err := buildmanifest.Build(root, buildmanifest.Input{
		Executable: filepath.ToSlash(relative),
		RegisteredVerbs: []string{
			"internal/verb/exec.go", "internal/verb/external.go", "internal/verb/registry.go",
			"verbs/dev-repo/commands.go", "verbs/dev-repo/episodes.go", "verbs/dev-repo/register.go",
			"verbs/dev-repo/repository.go", "verbs/dev-repo/schemas.go",
		},
		Schemas:         []string{"spec/result.schema.json", "spec/episode.schema.json", "spec/world.schema.json"},
		WorldDefinition: "worlds/dev-repo/world.json",
		AuthoredRules:   []string{"verbs/dev-repo/activations.go", "verbs/dev-repo/refusals.go", "verbs/dev-repo/transitions.go"},
		GOOS:            recipe.goos, GOARCH: recipe.goarch, CGOEnabled: false,
	})
	if err != nil {
		cleanup()
		return "", "", func() {}, err
	}
	return executable, manifest.Digest, cleanup, nil
}

func phoenixBuildCommand(root, executable string, recipe phoenixBuildRecipe) *exec.Cmd {
	command := exec.Command("go", "build", "-trimpath", "-buildvcs=false",
		"-ldflags=-s -w -buildid= -X main.version="+recipe.version, "-o", executable, "./cmd/phoenix")
	command.Dir = root
	command.Env = append(os.Environ(), "CGO_ENABLED=0", "GOTOOLCHAIN=go1.26.6", "GOOS="+recipe.goos, "GOARCH="+recipe.goarch)
	return command
}

func runCases(config runConfig, ids []string, runtime runtimeDriver, grader gradeDriver) (authoringSummary, error) {
	summary := authoringSummary{V: 1, Tranche: "authoring", Arm: config.arm, WorldBuild: config.worldBuild, Cases: []caseResult{}}
	for _, id := range ids {
		result, err := runCase(config, id, runtime, grader)
		if err != nil {
			return authoringSummary{}, fmt.Errorf("%s: %w", id, err)
		}
		summary.Cases = append(summary.Cases, result)
		if result.Status == "pass" {
			summary.Passed++
		} else {
			summary.Failed++
		}
		fmt.Printf("%s: %s\n", id, result.Status)
	}
	return summary, nil
}

func runCase(config runConfig, caseID string, runtime runtimeDriver, grader gradeDriver) (caseResult, error) {
	attempt, result, err := executeCaseAttempt(config, caseID, runtime)
	if err != nil {
		return caseResult{}, err
	}
	defer attempt.cleanup()
	if _, err := writeRuntimeEvidence(config, caseID+".runtime.jsonl", result.RawOutput); err != nil {
		return caseResult{}, err
	}
	if result.Failure != nil {
		return caseResult{}, fmt.Errorf("runtime %s: %s", result.Failure.Code, result.Failure.Message)
	}
	return finishCase(config, caseID, caseID, attempt.sandbox, attempt.stateDir, attempt.roots, result, grader)
}

type caseAttempt struct {
	sandbox  string
	stateDir string
	roots    map[string]string
}

func (attempt caseAttempt) cleanup() {
	_ = os.RemoveAll(attempt.sandbox)
	_ = os.RemoveAll(attempt.stateDir)
}

func executeCaseAttempt(config runConfig, caseID string, runtime runtimeDriver) (caseAttempt, runtimeResult, error) {
	tranche := effectiveTranche(config)
	item, err := loadCaseForTranche(config.repositoryRoot, tranche, caseID)
	if err != nil {
		return caseAttempt{}, runtimeResult{}, err
	}
	worldRef, err := digestJSONFile(config.worldPath)
	if err != nil {
		return caseAttempt{}, runtimeResult{}, err
	}
	if item.WorldRef != worldRef {
		return caseAttempt{}, runtimeResult{}, fmt.Errorf("case world_ref %s does not match production world %s", item.WorldRef, worldRef)
	}
	fixture, err := loadFixtureForTranche(config.repositoryRoot, tranche, item)
	if err != nil {
		return caseAttempt{}, runtimeResult{}, err
	}
	sandbox, err := os.MkdirTemp("", caseID+"-")
	if err != nil {
		return caseAttempt{}, runtimeResult{}, err
	}
	attempt := caseAttempt{sandbox: sandbox}
	if err := materializeFixture(sandbox, fixture); err != nil {
		attempt.cleanup()
		return caseAttempt{}, runtimeResult{}, err
	}
	stateDir, err := os.MkdirTemp("", caseID+"-state-")
	if err != nil {
		attempt.cleanup()
		return caseAttempt{}, runtimeResult{}, err
	}
	attempt.stateDir = stateDir
	roots, err := opaqueRoots()
	if err != nil {
		attempt.cleanup()
		return caseAttempt{}, runtimeResult{}, err
	}
	attempt.roots = roots
	stateEventsPath := ""
	if len(item.StateChanges) > 0 {
		stateEventsPath = filepath.Join(stateDir, "state-events.json")
		if err := writeJSON(stateEventsPath, map[string]any{"v": 1, "events": item.StateChanges}); err != nil {
			attempt.cleanup()
			return caseAttempt{}, runtimeResult{}, err
		}
	}
	result, err := runtime.Run(runtimeRequest{
		Arm: config.arm, ArmBDocument: config.armBDocument,
		FlatToolNames: config.flatToolNames,
		Sandbox:       sandbox, Goal: item.Goal, Roots: roots,
		PhoenixPath: config.phoenixPath, WorldPath: config.worldPath, SchemaPath: config.schemaPath,
		EpisodePath: filepath.Join(stateDir, "episodes.db"), WorldBuild: config.worldBuild,
		StateEventsPath: stateEventsPath,
		Timeout:         config.timeout, BudgetUSD: config.budgetUSD,
	})
	if err != nil {
		attempt.cleanup()
		return caseAttempt{}, runtimeResult{}, err
	}
	return attempt, result, nil
}

func finishCase(
	config runConfig,
	stem, caseID, sandbox, stateDir string,
	roots map[string]string,
	runtimeResult runtimeResult,
	grader gradeDriver,
) (caseResult, error) {
	record, err := readTrialEpisode(filepath.Join(stateDir, "episodes.db"))
	if err != nil {
		return caseResult{}, err
	}
	acts, orientations, err := episodeEvidence(record, roots)
	if err != nil {
		return caseResult{}, err
	}
	end, err := captureEndState(sandbox)
	if err != nil {
		return caseResult{}, err
	}
	trialRecord := trial{
		V: 1, CaseID: caseID, WorldBuild: config.worldBuild,
		Acts: acts, Orientations: orientations, FinalMessage: runtimeResult.FinalMessage, EndState: end,
	}
	trialPath := filepath.Join(config.outputDir, stem+".trial.json")
	if err := writeJSON(trialPath, trialRecord); err != nil {
		return caseResult{}, err
	}
	trialRelative, err := filepath.Rel(config.repositoryRoot, trialPath)
	if err != nil {
		return caseResult{}, err
	}
	grade, err := grader.Grade(config.repositoryRoot, effectiveTranche(config), caseID, filepath.ToSlash(trialRelative))
	if err != nil {
		return caseResult{}, err
	}
	gradePath := filepath.Join(config.outputDir, stem+".grade.json")
	if err := os.WriteFile(gradePath, grade, 0o644); err != nil {
		return caseResult{}, err
	}
	status, err := gradeStatus(grade)
	if err != nil {
		return caseResult{}, err
	}
	gradeRelative, _ := filepath.Rel(config.repositoryRoot, gradePath)
	return caseResult{
		CaseID: caseID, TrialPath: filepath.ToSlash(trialRelative),
		GradePath: filepath.ToSlash(gradeRelative), Status: status,
	}, nil
}

func writeRuntimeEvidence(config runConfig, name string, output []byte) (string, error) {
	sanitized, err := sanitizeRuntimeOutput(output)
	if err != nil {
		return "", err
	}
	path := filepath.Join(config.outputDir, name)
	if err := os.WriteFile(path, sanitized, 0o644); err != nil {
		return "", err
	}
	relative, err := filepath.Rel(config.repositoryRoot, path)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(relative), nil
}

func readTrialEpisode(path string) (episode.Record, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return episode.Record{Acts: []episode.ActRecord{}}, nil
	}
	store, err := episode.Open(path, episode.Options{})
	if err != nil {
		return episode.Record{}, err
	}
	defer store.Close()
	record, err := store.Latest(context.Background())
	if errors.Is(err, episode.ErrNotFound) {
		return episode.Record{Acts: []episode.ActRecord{}}, nil
	}
	return record, err
}

func episodeEvidence(record episode.Record, roots map[string]string) ([]act, []orientation, error) {
	handleTypes := make(map[string]string, len(roots))
	for name, ref := range roots {
		handleTypes[ref] = name
	}
	acts := make([]act, 0, len(record.Acts))
	orientations := []orientation{}
	for _, recorded := range record.Acts {
		encoded, err := json.Marshal(recorded.Result)
		if err != nil {
			return nil, nil, err
		}
		var envelope surface.Envelope
		if err := json.Unmarshal(encoded, &envelope); err != nil {
			return nil, nil, fmt.Errorf("decode act %s result: %w", recorded.ActID, err)
		}
		for _, grant := range envelope.Handles.Grant {
			handleTypes[grant.Ref] = grant.Type
		}
		if recorded.Request.Intent != "" {
			orientations = append(orientations, orientation{
				Seq: len(orientations), Handle: handleTypes[recorded.Request.Handle], Intent: recorded.Request.Intent,
				Matched: len(envelope.Frontier) > 0, Calls: len(envelope.Frontier),
			})
			continue
		}
		output, err := json.Marshal(envelope)
		if err != nil {
			return nil, nil, err
		}
		item := act{
			Seq: len(acts), Handle: handleTypes[recorded.Request.Handle], HandleType: handleTypes[recorded.Request.Handle],
			Verb: recorded.Request.Verb, Args: recorded.Request.Args, Status: string(envelope.Status), Output: string(output),
		}
		if item.HandleType == "" {
			item.HandleType = "unknown"
			item.Handle = "unknown"
		}
		if result, ok := envelope.Result.(map[string]any); ok {
			if value, ok := result["exit_code"].(float64); ok {
				exitCode := int(value)
				item.ExitCode = &exitCode
			}
		}
		acts = append(acts, item)
	}
	return acts, orientations, nil
}

func opaqueRoots() (map[string]string, error) {
	roots := make(map[string]string, 4)
	for _, name := range []string{"repo", "tests", "git", "episodes"} {
		value := make([]byte, 18)
		if _, err := rand.Read(value); err != nil {
			return nil, err
		}
		roots[name] = "h_" + base64.RawURLEncoding.EncodeToString(value)
	}
	return roots, nil
}

func gradeStatus(contents []byte) (string, error) {
	var grade struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(contents, &grade); err != nil {
		return "", fmt.Errorf("decode grade: %w", err)
	}
	if grade.Status != "pass" && grade.Status != "fail" && grade.Status != "indeterminate" {
		return "", fmt.Errorf("unknown grade status %q", grade.Status)
	}
	return grade.Status, nil
}

func ensureInside(root, target string) error {
	relative, err := filepath.Rel(root, target)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
		return fmt.Errorf("output directory must stay inside the repository")
	}
	return nil
}
