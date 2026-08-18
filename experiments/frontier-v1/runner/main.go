// Command runner executes the frontier-v1 authoring tranche against a fresh
// sandbox, Phoenix process, and pinned runtime session for every case.
package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/nesste/phoenix/internal/buildmanifest"
	"github.com/nesste/phoenix/internal/episode"
	"github.com/nesste/phoenix/internal/surface"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	flags := flag.NewFlagSet("frontier-v1-runner", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	repositoryRoot := flags.String("repo-root", ".", "Phoenix repository root")
	caseID := flags.String("case", "all", "authoring case id or all")
	outputDir := flags.String("output-dir", "experiments/frontier-v1/results/authoring", "repository-relative result directory")
	runtimePath := flags.String("runtime", "claude", "pinned runtime executable")
	budget := flags.String("max-budget-usd", "0.15", "maximum model cost per case")
	timeout := flags.Duration("timeout", 180*time.Second, "per-case timeout")
	if err := flags.Parse(args); err != nil || flags.NArg() != 0 {
		return 2
	}
	root, err := filepath.Abs(*repositoryRoot)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	config, cleanup, err := prepareRunConfig(root, *outputDir, *timeout, *budget)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer cleanup()
	driver := claudeDriver{executable: *runtimePath}
	if err := driver.Verify(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	ids := []string{*caseID}
	if *caseID == "all" {
		ids, err = authoringCaseIDs(root)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
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

func prepareRunConfig(root, output string, timeout time.Duration, budget string) (runConfig, func(), error) {
	if timeout <= 0 || strings.TrimSpace(budget) == "" {
		return runConfig{}, func() {}, fmt.Errorf("positive timeout and budget are required")
	}
	outputPath := output
	if !filepath.IsAbs(outputPath) {
		outputPath = filepath.Join(root, filepath.FromSlash(output))
	}
	if err := ensureInside(root, outputPath); err != nil {
		return runConfig{}, func() {}, err
	}
	if err := os.MkdirAll(outputPath, 0o755); err != nil {
		return runConfig{}, func() {}, err
	}
	phoenixPath, build, cleanup, err := buildPhoenix(root)
	if err != nil {
		return runConfig{}, func() {}, err
	}
	return runConfig{
		repositoryRoot: root, outputDir: outputPath,
		worldPath:   filepath.Join(root, "worlds", "dev-repo", "world.json"),
		schemaPath:  filepath.Join(root, "spec", "world.schema.json"),
		phoenixPath: phoenixPath, worldBuild: build,
		timeout: timeout, budgetUSD: budget,
	}, cleanup, nil
}

func buildPhoenix(root string) (string, string, func(), error) {
	tempRoot := filepath.Join(root, "build")
	if err := os.MkdirAll(tempRoot, 0o755); err != nil {
		return "", "", func() {}, err
	}
	dir, err := os.MkdirTemp(tempRoot, "authoring-runner-")
	if err != nil {
		return "", "", func() {}, err
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	name := "phoenix"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	executable := filepath.Join(dir, name)
	command := exec.Command("go", "build", "-trimpath", "-buildvcs=false", "-ldflags=-s -w -buildid= -X main.version=authoring", "-o", executable, "./cmd/phoenix")
	command.Dir = root
	command.Env = append(os.Environ(), "CGO_ENABLED=0", "GOTOOLCHAIN=go1.26.6")
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
		GOOS:            runtime.GOOS, GOARCH: runtime.GOARCH, CGOEnabled: false,
	})
	if err != nil {
		cleanup()
		return "", "", func() {}, err
	}
	return executable, manifest.Digest, cleanup, nil
}

func runCases(config runConfig, ids []string, runtime runtimeDriver, grader gradeDriver) (authoringSummary, error) {
	summary := authoringSummary{V: 1, Tranche: "authoring", WorldBuild: config.worldBuild, Cases: []caseResult{}}
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
	item, err := loadCase(config.repositoryRoot, caseID)
	if err != nil {
		return caseResult{}, err
	}
	worldRef, err := digestJSONFile(config.worldPath)
	if err != nil {
		return caseResult{}, err
	}
	if item.WorldRef != worldRef {
		return caseResult{}, fmt.Errorf("case world_ref %s does not match production world %s", item.WorldRef, worldRef)
	}
	fixture, err := loadFixture(config.repositoryRoot, item)
	if err != nil {
		return caseResult{}, err
	}
	sandbox, err := os.MkdirTemp("", caseID+"-")
	if err != nil {
		return caseResult{}, err
	}
	defer os.RemoveAll(sandbox)
	if err := materializeFixture(sandbox, fixture); err != nil {
		return caseResult{}, err
	}
	stateDir, err := os.MkdirTemp("", caseID+"-state-")
	if err != nil {
		return caseResult{}, err
	}
	defer os.RemoveAll(stateDir)
	roots, err := opaqueRoots()
	if err != nil {
		return caseResult{}, err
	}
	stateEventsPath := ""
	if len(item.StateChanges) > 0 {
		stateEventsPath = filepath.Join(stateDir, "state-events.json")
		if err := writeJSON(stateEventsPath, map[string]any{"v": 1, "events": item.StateChanges}); err != nil {
			return caseResult{}, err
		}
	}
	runtimeResult, err := runtime.Run(runtimeRequest{
		Sandbox: sandbox, Goal: item.Goal, Roots: roots,
		PhoenixPath: config.phoenixPath, WorldPath: config.worldPath, SchemaPath: config.schemaPath,
		EpisodePath: filepath.Join(stateDir, "episodes.db"), WorldBuild: config.worldBuild,
		StateEventsPath: stateEventsPath,
		Timeout:         config.timeout, BudgetUSD: config.budgetUSD,
	})
	if err != nil {
		return caseResult{}, err
	}
	return finishCase(config, caseID, sandbox, stateDir, roots, runtimeResult, grader)
}

func finishCase(
	config runConfig,
	caseID, sandbox, stateDir string,
	roots map[string]string,
	runtimeResult runtimeResult,
	grader gradeDriver,
) (caseResult, error) {
	sanitizedRuntime, err := sanitizeRuntimeOutput(runtimeResult.RawOutput)
	if err != nil {
		return caseResult{}, err
	}
	if err := os.WriteFile(filepath.Join(config.outputDir, caseID+".runtime.jsonl"), sanitizedRuntime, 0o644); err != nil {
		return caseResult{}, err
	}
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
	trialPath := filepath.Join(config.outputDir, caseID+".trial.json")
	if err := writeJSON(trialPath, trialRecord); err != nil {
		return caseResult{}, err
	}
	trialRelative, err := filepath.Rel(config.repositoryRoot, trialPath)
	if err != nil {
		return caseResult{}, err
	}
	labelRelative := filepath.ToSlash(filepath.Join("experiments", "frontier-v1", "labels", "authoring", caseID+".json"))
	grade, err := grader.Grade(config.repositoryRoot, labelRelative, filepath.ToSlash(trialRelative))
	if err != nil {
		return caseResult{}, err
	}
	gradePath := filepath.Join(config.outputDir, caseID+".grade.json")
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

func readTrialEpisode(path string) (episode.Record, error) {
	store, err := episode.Open(path, episode.Options{})
	if err != nil {
		return episode.Record{}, err
	}
	defer store.Close()
	return store.Latest(context.Background())
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
