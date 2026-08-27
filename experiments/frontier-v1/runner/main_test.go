package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nesste/phoenix/internal/episode"
	"github.com/nesste/phoenix/internal/surface"
)

const testBuild = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func TestRunnerRejectsSealedCaseBeforeFilesystemAccess(t *testing.T) {
	_, err := loadCase(t.TempDir(), "validation_secret")
	if err == nil || !strings.Contains(err.Error(), `tranche "authoring"`) {
		t.Fatalf("sealed case error = %v", err)
	}
}

func TestPinnedAuthoringWorldMatchesProductionWorld(t *testing.T) {
	production, err := digestJSONFile(filepath.Join("..", "..", "..", "worlds", "dev-repo", "world.json"))
	if err != nil {
		t.Fatal(err)
	}
	pinned, err := digestJSONFile(filepath.Join("..", "worlds", "authoring.dev_repo.json"))
	if err != nil {
		t.Fatal(err)
	}
	if production != pinned {
		t.Fatalf("pinned authoring world = %s, production = %s", pinned, production)
	}
}

func TestAuthoringLabelsSharePinnedGraderDigest(t *testing.T) {
	repositoryRoot := filepath.Join("..", "..", "..")
	ids, err := authoringCaseIDs(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	digest, err := graderDigestForCases(repositoryRoot, ids)
	if err != nil {
		t.Fatal(err)
	}
	const want = "sha256:7b7438f84164d69157bbde87fd3ffb2f88e069bb01a97785a737910e29df0e05"
	if digest != want {
		t.Fatalf("authoring grader digest = %s, want %s", digest, want)
	}
}

func TestMaterializeFixtureCreatesInitialCommitAndPendingFiles(t *testing.T) {
	target := t.TempDir()
	item := fixture{
		Files:       map[string]string{"go.mod": "module example.test/fixture\n", "pending.txt": "changed\n"},
		Uncommitted: []string{"pending.txt"},
	}
	if err := materializeFixture(target, item); err != nil {
		t.Fatal(err)
	}
	commandOutput, err := gitOutput(target, "status", "--porcelain")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(commandOutput, "?? pending.txt") {
		t.Fatalf("git status = %q, want pending fixture path", commandOutput)
	}
	if contents, err := os.ReadFile(filepath.Join(target, "pending.txt")); err != nil || string(contents) != "changed\n" {
		t.Fatalf("pending file = %q, %v", contents, err)
	}
}

func TestParseRuntimeOutputFindsFinalStreamEvent(t *testing.T) {
	output := []byte("{\"type\":\"assistant\"}\n{\"type\":\"result\",\"result\":\"suite is not green\"}\n")
	final, err := parseRuntimeOutput(output)
	if err != nil {
		t.Fatal(err)
	}
	if final != "suite is not green" {
		t.Fatalf("final message = %q", final)
	}
}

func TestSystemPromptsArePinnedPerArmAndIntentIsPhoenixOnly(t *testing.T) {
	for _, arm := range []string{"A", "B", "C", "D"} {
		prompt, err := systemPromptForArm(arm)
		if err != nil {
			t.Fatalf("arm %s: %v", arm, err)
		}
		hasIntent := strings.Contains(prompt, phoenixIntentPrompt)
		wantIntent := arm == "C" || arm == "D"
		if hasIntent != wantIntent {
			t.Fatalf("arm %s intent instruction = %t, want %t: %q", arm, hasIntent, wantIntent, prompt)
		}
	}
	if _, err := systemPromptForArm("E"); err == nil {
		t.Fatal("retired arm E prompt was accepted")
	}
	if _, err := systemPromptForArm("D_prime"); err == nil {
		t.Fatal("unfrozen Phase 2 prompt was accepted")
	}
}

func TestSanitizeRuntimeOutputDropsLocalEnvironmentMetadata(t *testing.T) {
	input := []byte("{\"type\":\"system\",\"subtype\":\"init\",\"cwd\":\"C:/Users/example\",\"session_id\":\"secret-session\",\"tools\":[\"mcp__phoenix__act\"],\"mcp_servers\":[{\"name\":\"phoenix\",\"status\":\"connected\"}],\"model\":\"claude-sonnet-5\",\"permissionMode\":\"default\",\"apiKeySource\":\"none\",\"claude_code_version\":\"2.1.229\",\"plugins\":[{\"path\":\"C:/Users/example/plugin\"}]}\n{\"type\":\"assistant\",\"message\":\"unchanged\"}\n")
	got, err := sanitizeRuntimeOutput(input)
	if err != nil {
		t.Fatal(err)
	}
	text := string(got)
	for _, forbidden := range []string{"C:/Users/example", "secret-session", "plugins"} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("sanitized runtime contains %q: %s", forbidden, text)
		}
	}
	for _, retained := range []string{"claude-sonnet-5", "mcp__phoenix__act", `{"type":"assistant","message":"unchanged"}`} {
		if !strings.Contains(text, retained) {
			t.Fatalf("sanitized runtime omitted %q: %s", retained, text)
		}
	}
}

func TestRunCaseMaterializesRunsRecordsAndGradesWithoutHarnessDependency(t *testing.T) {
	repository := t.TempDir()
	createRunnerFixture(t, repository)
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	driver := &fakeRuntime{}
	result, err := runCase(runConfig{
		arm:            "C",
		repositoryRoot: repository, outputDir: output,
		worldPath: filepath.Join(repository, "world.json"), schemaPath: filepath.Join(repository, "schema.json"),
		worldBuild: testBuild, timeout: time.Second, budgetUSD: "0.01",
	}, "authoring_testcase", driver, fakeGrader{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "pass" {
		t.Fatalf("case result = %#v", result)
	}
	if _, err := os.Stat(driver.sandbox); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("sandbox survived trial: %s, %v", driver.sandbox, err)
	}
	var recorded trial
	if err := decodeStrict(filepath.Join(output, "authoring_testcase.trial.json"), &recorded); err != nil {
		t.Fatal(err)
	}
	if len(recorded.Acts) != 1 || recorded.Acts[0].HandleType != "tests" || recorded.Acts[0].Status != "ok" {
		t.Fatalf("recorded acts = %#v", recorded.Acts)
	}
}

func TestNormalizeArmAndDefaultOutputProtectRetainedEvidence(t *testing.T) {
	if arm, err := normalizeArm(" d "); err != nil || arm != "D" {
		t.Fatalf("normalize arm = %q, %v", arm, err)
	}
	if _, err := normalizeArm("D_prime"); err == nil {
		t.Fatal("unfrozen Phase 2 arm was accepted")
	}
}

func TestWriteScheduleCLIUsesFrozenDesignWithoutLaunchingRuntime(t *testing.T) {
	repository := t.TempDir()
	createRunnerFixture(t, repository)
	relative := "experiments/frontier-v1/authoring.schedule.json"
	if code := run([]string{"--repo-root", repository, "--write-schedule", relative}); code != 0 {
		t.Fatalf("write schedule exit = %d", code)
	}
	var schedule launchSchedule
	path := filepath.Join(repository, filepath.FromSlash(relative))
	if err := decodeStrict(path, &schedule); err != nil {
		t.Fatal(err)
	}
	if schedule.Seed != phase1ScheduleSeed || schedule.Repetitions != phase1Repetitions ||
		len(schedule.Entries) != phase1Repetitions*len(phase1Arms) {
		t.Fatalf("written schedule = %#v", schedule)
	}
	if code := run([]string{"--repo-root", repository, "--write-schedule", relative}); code == 0 {
		t.Fatal("schedule writer replaced an existing schedule")
	}
}

func TestScheduledOutputMustBeEmpty(t *testing.T) {
	repository := t.TempDir()
	output := filepath.Join(repository, "results")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := requireEmptyScheduledOutput(repository, "results"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, "existing.json"), []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := requireEmptyScheduledOutput(repository, "results"); err == nil {
		t.Fatal("scheduled runner accepted an output directory containing evidence")
	}
}

func TestFlatRuntimePromptAndAllowedToolsContainNoPhoenixHandles(t *testing.T) {
	roots := map[string]string{"repo": "h_secret_repo", "tests": "h_secret_tests", "git": "h_secret_git", "episodes": "h_secret_episodes"}
	prompt := runtimePromptForArm("inspect the repository", roots, "A")
	if prompt != "inspect the repository" || strings.Contains(prompt, "h_secret") {
		t.Fatalf("flat prompt leaked handles: %q", prompt)
	}
	tools := allowedToolsForArm("A", []string{"repo_status", "tests_run"})
	if len(tools) != 2 || tools[0] != "mcp__phoenix__repo_status" || tools[1] != "mcp__phoenix__tests_run" {
		t.Fatalf("flat allowed tools = %v", tools)
	}
	if got := allowedToolsForArm("D", nil); len(got) != 1 || got[0] != "mcp__phoenix__act" {
		t.Fatalf("Phoenix allowed tools = %v", got)
	}
}

type fakeRuntime struct {
	sandbox string
}

func (runtime *fakeRuntime) Verify() error { return nil }

func (runtime *fakeRuntime) Run(request runtimeRequest) (runtimeResult, error) {
	runtime.sandbox = request.Sandbox
	store, err := episode.Open(request.EpisodePath, episode.Options{})
	if err != nil {
		return runtimeResult{}, err
	}
	episodeID, err := store.StartEpisode(context.Background(), "s_runner", request.WorldBuild)
	if err != nil {
		return runtimeResult{}, err
	}
	started := episode.StartedAct{
		EpisodeID: episodeID, ActID: "a_runner", WorldBuild: request.WorldBuild,
		Request: episode.ActRequest{Handle: request.Roots["tests"], Verb: "run", Args: map[string]any{}},
	}
	if err := store.BeginAct(context.Background(), started); err != nil {
		return runtimeResult{}, err
	}
	envelope := surface.Envelope{
		V: 1, WorldBuild: request.WorldBuild, SessionID: "s_runner", ActID: "a_runner",
		Handle: request.Roots["tests"], Verb: "run", Status: surface.StatusOK,
		Result: map[string]any{"failed": 1}, Text: "suite failed",
		Handles:  surface.HandleDelta{Grant: []surface.HandleGrant{}, Revoke: []string{}},
		Frontier: []surface.FrontierEntry{},
	}
	if err := store.CompleteAct(context.Background(), episode.CompletedAct{EpisodeID: episodeID, ActID: "a_runner", Result: envelope}); err != nil {
		return runtimeResult{}, err
	}
	if err := store.Close(); err != nil {
		return runtimeResult{}, err
	}
	return runtimeResult{FinalMessage: "suite is not green", RawOutput: []byte("runtime evidence\n"), FirstModelToken: true}, nil
}

type fakeGrader struct{}

func (fakeGrader) Grade(string, string, string, string) ([]byte, error) {
	return []byte("{\"status\":\"pass\"}\n"), nil
}

func createRunnerFixture(t *testing.T, repository string) {
	t.Helper()
	world := map[string]any{"v": 1, "id": "test"}
	if err := writeJSON(filepath.Join(repository, "world.json"), world); err != nil {
		t.Fatal(err)
	}
	worldDigest, err := digestJSONFile(filepath.Join(repository, "world.json"))
	if err != nil {
		t.Fatal(err)
	}
	fixturePath := filepath.Join(repository, "experiments", "frontier-v1", "fixtures", "authoring", "fixture.json")
	if err := writeJSON(fixturePath, fixture{
		V: 1, FixtureID: "authoring_family_test_v1", FamilyID: "authoring_family_test",
		Template: "A deterministic runner unit fixture.", Files: map[string]string{"go.mod": "module example.test/runner\n"},
	}); err != nil {
		t.Fatal(err)
	}
	fixtureDigest, err := digestJSONFile(fixturePath)
	if err != nil {
		t.Fatal(err)
	}
	casePath := filepath.Join(repository, "experiments", "frontier-v1", "corpus", "authoring", "authoring_testcase.json")
	if err := writeJSON(casePath, runnableCase{
		CaseID: "authoring_testcase", Goal: "Run the complete test suite and report the result.",
		SandboxFixture: fixtureDigest, WorldRef: worldDigest, FamilyID: "authoring_family_test",
	}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repository, "schema.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func gitOutput(dir string, args ...string) (string, error) {
	command := exec.Command("git", args...)
	command.Dir = dir
	output, err := command.CombinedOutput()
	return string(output), err
}

func TestGradeStatusRejectsUnknownVerdict(t *testing.T) {
	if _, err := gradeStatus(json.RawMessage(`{"status":"maybe"}`)); err == nil {
		t.Fatal("gradeStatus accepted an unknown verdict")
	}
}

func TestTrimGoRunExitPreservesGradeJSON(t *testing.T) {
	got := trimGoRunExit([]byte("{\"status\":\"fail\"}\nexit status 1\n"))
	if string(got) != "{\"status\":\"fail\"}\n" {
		t.Fatalf("trimmed grade = %q", got)
	}
}
