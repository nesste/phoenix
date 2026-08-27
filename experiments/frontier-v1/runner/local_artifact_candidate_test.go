package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/nesste/phoenix/internal/buildmanifest"
	"github.com/nesste/phoenix/internal/surface"
	"github.com/nesste/phoenix/internal/world"
)

type localArtifactCandidate struct {
	V           int    `json:"v"`
	Status      string `json:"status"`
	SourceLimit string `json:"source_limit"`
	Frozen      bool   `json:"frozen"`
	Gates       struct {
		MayOpenValidation bool `json:"may_open_validation"`
		MayOpenHeldOut    bool `json:"may_open_held_out"`
	} `json:"gates"`
	Artifacts struct {
		Runtime struct {
			ProtocolRuntime          map[string]any    `json:"protocol_runtime"`
			TrialLimits              map[string]any    `json:"trial_limits"`
			DefaultRuntimeExecutable string            `json:"default_runtime_executable"`
			ModelInvocation          []string          `json:"model_invocation_template"`
			PhoenixInvocation        []string          `json:"phoenix_server_invocation_template"`
			SystemPrompts            map[string]string `json:"system_prompts"`
			Files                    map[string]string `json:"files_lf_normalized_utf8_sha256"`
			Frozen                   bool              `json:"frozen"`
		} `json:"runtime_invocation_and_exact_per_arm_system_prompts"`
		ArmA struct {
			Files  map[string]string `json:"files_lf_normalized_utf8_sha256"`
			Tools  []candidateTool   `json:"tools"`
			Frozen bool              `json:"frozen"`
		} `json:"arm_a_schemas"`
		WorldBuild struct {
			WorldDefinition          string `json:"world_definition"`
			WorldDefinitionRawDigest string `json:"world_definition_raw_sha256"`
			CanonicalWorldDigest     string `json:"world_definition_canonical_json_sha256"`
			AuthoringWorld           string `json:"authoring_world_equivalent"`
			ManifestPath             string `json:"world_build_manifest"`
			ManifestRawDigest        string `json:"world_build_manifest_raw_sha256"`
			BuildDigest              string `json:"world_build_digest"`
			Frozen                   bool   `json:"frozen"`
		} `json:"world_definition_and_world_build_digest"`
		Grader struct {
			Digest                   string `json:"digest"`
			DigestImplementation     string `json:"digest_implementation"`
			DigestImplementationHash string `json:"digest_implementation_lf_normalized_utf8_sha256"`
			Frozen                   bool   `json:"frozen"`
		} `json:"grader_digest"`
	} `json:"artifacts"`
	Remaining []string `json:"remaining_after_acceptance"`
}

type candidateTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Pointer     string `json:"args_schema_pointer"`
}

func TestLocalArtifactCandidateMatchesImplementation(t *testing.T) {
	repositoryRoot := filepath.Join("..", "..", "..")
	candidate := readLocalArtifactCandidate(t)
	if candidate.V != 1 || candidate.Status != "review_candidate" || candidate.SourceLimit != "authoring_only" ||
		candidate.Frozen || candidate.Gates.MayOpenValidation || candidate.Gates.MayOpenHeldOut {
		t.Fatalf("local artifact candidate opened a gate or is not review-only: %#v", candidate)
	}
	if !reflect.DeepEqual(candidate.Remaining, []string{"schedule digest"}) {
		t.Fatalf("remaining artifacts = %#v", candidate.Remaining)
	}
	if candidate.Artifacts.Runtime.Frozen || candidate.Artifacts.ArmA.Frozen ||
		candidate.Artifacts.WorldBuild.Frozen || candidate.Artifacts.Grader.Frozen {
		t.Fatal("review candidate marked a local artifact frozen before independent acceptance")
	}

	verifyCandidateFilesAtCommit(t, repositoryRoot, "54e256e9f4347d34844a3ef9b2156600360dcd13", candidate.Artifacts.Runtime.Files)
	verifyCandidateFiles(t, repositoryRoot, candidate.Artifacts.ArmA.Files)
	verifyRuntimeCandidate(t, candidate)
	verifyInvocationTemplates(t, candidate)
	verifyArmACandidate(t, repositoryRoot, candidate)
	verifyWorldBuildCandidate(t, repositoryRoot, candidate)
	verifyGraderCandidate(t, repositoryRoot, candidate)
}

func verifyInvocationTemplates(t *testing.T, candidate localArtifactCandidate) {
	t.Helper()
	temporary := t.TempDir()
	episodePath := filepath.Join(temporary, "episodes.db")
	stateEventsPath := filepath.Join(temporary, "state-events.json")
	request := runtimeRequest{
		Arm: "C", Goal: "candidate goal", Roots: map[string]string{
			"repo": "h_repo", "tests": "h_tests", "git": "h_git", "episodes": "h_episodes",
		},
		PhoenixPath: "<freshly built Phoenix executable>",
		WorldPath:   "worlds/dev-repo/world.json", SchemaPath: "spec/world.schema.json",
		EpisodePath: episodePath, WorldBuild: "<world-build digest>", StateEventsPath: stateEventsPath,
		BudgetUSD: "0.15",
	}
	modelArguments, err := prepareRuntimeInvocation(
		request, "<exact arm system prompt>", []string{"<arm-specific comma-separated allowlist>"},
	)
	if err != nil {
		t.Fatal(err)
	}
	modelInvocation := append([]string{candidate.Artifacts.Runtime.DefaultRuntimeExecutable}, modelArguments...)
	replaceArgument(modelInvocation, runtimePromptForArm(request.Goal, request.Roots, request.Arm), "<arm-specific user prompt>")
	replaceArgument(modelInvocation, filepath.Join(temporary, "mcp.json"), "<per-attempt mcp.json>")
	wantModel := candidate.Artifacts.Runtime.ModelInvocation
	if len(wantModel) < 2 {
		t.Fatal("candidate model invocation omits the Arm B suffix marker")
	}
	baseModel := wantModel[:len(wantModel)-2]
	if !reflect.DeepEqual(modelInvocation, baseModel) {
		t.Fatalf("candidate model invocation = %#v, implementation = %#v", baseModel, modelInvocation)
	}

	armBRequest := request
	armBRequest.Arm = "B"
	armBRequest.ArmBDocument = "experiments/frontier-v1/arms/arm-b.md"
	armBArguments, err := prepareRuntimeInvocation(
		armBRequest, "<exact arm system prompt>", []string{"<arm-specific comma-separated allowlist>"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if got := armBArguments[len(armBArguments)-2:]; !reflect.DeepEqual(got, []string{
		"--append-system-prompt-file", armBRequest.ArmBDocument,
	}) {
		t.Fatalf("Arm B invocation suffix = %#v", got)
	}
	if !reflect.DeepEqual(wantModel[len(wantModel)-2:], []string{
		"[Arm B only] --append-system-prompt-file", armBRequest.ArmBDocument,
	}) {
		t.Fatalf("candidate Arm B invocation suffix = %#v", wantModel[len(wantModel)-2:])
	}

	var mcpConfig struct {
		Servers map[string]struct {
			Command string   `json:"command"`
			Args    []string `json:"args"`
		} `json:"mcpServers"`
	}
	readJSONForTest(t, filepath.Join(temporary, "mcp.json"), &mcpConfig)
	server := mcpConfig.Servers["phoenix"]
	phoenixInvocation := append([]string{server.Command}, server.Args...)
	replaceArgument(phoenixInvocation, episodePath, "<per-attempt episodes.db>")
	replaceArgument(phoenixInvocation, filepath.Join(temporary, "roots.json"), "<per-attempt roots.json>")
	replaceArgument(phoenixInvocation, stateEventsPath, "<per-attempt state-events.json>")
	replaceArgument(phoenixInvocation, armBRequest.Arm, "<A|B|C|D>")
	phoenixInvocation[len(phoenixInvocation)-2] = "[when declared] --state-events"
	if !reflect.DeepEqual(phoenixInvocation, candidate.Artifacts.Runtime.PhoenixInvocation) {
		t.Fatalf("candidate Phoenix invocation = %#v, implementation = %#v", candidate.Artifacts.Runtime.PhoenixInvocation, phoenixInvocation)
	}
}

func replaceArgument(arguments []string, old, replacement string) {
	for index, argument := range arguments {
		if argument == old {
			arguments[index] = replacement
		}
	}
}

func readLocalArtifactCandidate(t *testing.T) localArtifactCandidate {
	t.Helper()
	var candidate localArtifactCandidate
	readJSONForTest(t, filepath.Join("..", "artifacts", "pre-validation-local-candidate.json"), &candidate)
	return candidate
}

func verifyCandidateFiles(t *testing.T, repositoryRoot string, files map[string]string) {
	t.Helper()
	for path, expected := range files {
		actual, err := digestLFNormalizedFile(filepath.Join(repositoryRoot, filepath.FromSlash(path)))
		if err != nil {
			t.Fatal(err)
		}
		if actual != expected {
			t.Fatalf("candidate file %s digest = %s, want %s", path, actual, expected)
		}
	}
}

func verifyCandidateFilesAtCommit(t *testing.T, repositoryRoot, commit string, files map[string]string) {
	t.Helper()
	for path, expected := range files {
		command := exec.Command("git", "-C", repositoryRoot, "show", commit+":"+path)
		contents, err := command.Output()
		if err != nil {
			t.Fatalf("read candidate file %s at %s: %v", path, commit, err)
		}
		normalized := strings.ReplaceAll(strings.ReplaceAll(string(contents), "\r\n", "\n"), "\r", "\n")
		actual := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(normalized)))
		if actual != expected {
			t.Fatalf("candidate file %s at %s digest = %s, want %s", path, commit, actual, expected)
		}
	}
}

func verifyRuntimeCandidate(t *testing.T, candidate localArtifactCandidate) {
	t.Helper()
	var protocol struct {
		Runtime       map[string]any    `json:"runtime"`
		SystemPrompts map[string]string `json:"system_prompts"`
		Trial         struct {
			TimeoutSeconds     int     `json:"timeout_seconds"`
			MaxCostUSDPerTrial float64 `json:"max_cost_usd_per_trial"`
			RetryPolicy        struct {
				Maximum int `json:"maximum_infrastructure_retries"`
			} `json:"retry_policy"`
		} `json:"trial"`
	}
	readJSONForTest(t, filepath.Join("..", "protocol.json"), &protocol)
	if !reflect.DeepEqual(candidate.Artifacts.Runtime.ProtocolRuntime, protocol.Runtime) {
		t.Fatal("candidate runtime contract differs from protocol runtime")
	}
	if !reflect.DeepEqual(candidate.Artifacts.Runtime.SystemPrompts, protocol.SystemPrompts) {
		t.Fatal("candidate system prompts differ from protocol prompts")
	}
	for _, arm := range phase1Arms {
		actual, err := systemPromptForArm(arm)
		if err != nil || actual != candidate.Artifacts.Runtime.SystemPrompts[arm] {
			t.Fatalf("arm %s system prompt = %q, %v", arm, actual, err)
		}
	}
	wantLimits := map[string]any{
		"timeout_seconds":                float64(protocol.Trial.TimeoutSeconds),
		"max_cost_usd_per_trial":         protocol.Trial.MaxCostUSDPerTrial,
		"maximum_infrastructure_retries": float64(protocol.Trial.RetryPolicy.Maximum),
	}
	if !reflect.DeepEqual(candidate.Artifacts.Runtime.TrialLimits, wantLimits) {
		t.Fatalf("candidate trial limits = %#v, want %#v", candidate.Artifacts.Runtime.TrialLimits, wantLimits)
	}
}

func verifyArmACandidate(t *testing.T, repositoryRoot string, candidate localArtifactCandidate) {
	t.Helper()
	definition, err := world.Load(
		filepath.Join(repositoryRoot, "spec", "world.schema.json"),
		filepath.Join(repositoryRoot, "worlds", "dev-repo", "world.json"),
	)
	if err != nil {
		t.Fatal(err)
	}
	names, err := surface.FlatToolNames(definition)
	if err != nil {
		t.Fatal(err)
	}
	want := make([]candidateTool, 0, len(names))
	for _, root := range definition.Roots {
		for verb := range definition.HandleTypes[root.Type].Verbs {
			want = append(want, candidateTool{
				Name:        root.Name + "_" + verb,
				Description: fmt.Sprintf("Invoke %s.%s on %s.", root.Type, verb, root.Label),
				Pointer:     fmt.Sprintf("/handle_types/%s/verbs/%s/args_schema", root.Type, verb),
			})
		}
	}
	sort.Slice(want, func(i, j int) bool { return want[i].Name < want[j].Name })
	if !reflect.DeepEqual(candidate.Artifacts.ArmA.Tools, want) {
		t.Fatalf("candidate Arm A tools = %#v, want %#v", candidate.Artifacts.ArmA.Tools, want)
	}
	for index, name := range names {
		if candidate.Artifacts.ArmA.Tools[index].Name != name {
			t.Fatalf("flat tool %d = %s, candidate has %s", index, name, candidate.Artifacts.ArmA.Tools[index].Name)
		}
	}
}

func verifyWorldBuildCandidate(t *testing.T, repositoryRoot string, candidate localArtifactCandidate) {
	t.Helper()
	artifact := candidate.Artifacts.WorldBuild
	if actual, err := digestLFNormalizedFile(filepath.Join(repositoryRoot, filepath.FromSlash(artifact.WorldDefinition))); err != nil || actual != artifact.WorldDefinitionRawDigest {
		t.Fatalf("world definition digest = %s, %v", actual, err)
	}
	production, err := digestJSONFile(filepath.Join(repositoryRoot, filepath.FromSlash(artifact.WorldDefinition)))
	if err != nil || production != artifact.CanonicalWorldDigest {
		t.Fatalf("canonical production world digest = %s, %v", production, err)
	}
	authoring, err := digestJSONFile(filepath.Join(repositoryRoot, filepath.FromSlash(artifact.AuthoringWorld)))
	if err != nil || authoring != production {
		t.Fatalf("canonical authoring world digest = %s, production = %s, %v", authoring, production, err)
	}

	manifestPath := filepath.Join(repositoryRoot, filepath.FromSlash(artifact.ManifestPath))
	if actual, err := digestLFNormalizedFile(manifestPath); err != nil || actual != artifact.ManifestRawDigest {
		t.Fatalf("world-build manifest file digest = %s, %v", actual, err)
	}
	var manifest buildmanifest.Manifest
	readJSONForTest(t, manifestPath, &manifest)
	encoded, err := json.Marshal(manifest.Build)
	if err != nil {
		t.Fatal(err)
	}
	recomputed := fmt.Sprintf("sha256:%x", sha256.Sum256(encoded))
	if manifest.Digest != recomputed || manifest.Digest != artifact.BuildDigest {
		t.Fatalf("world-build digest = %s, recomputed %s, candidate %s", manifest.Digest, recomputed, artifact.BuildDigest)
	}
}

func verifyGraderCandidate(t *testing.T, repositoryRoot string, candidate localArtifactCandidate) {
	t.Helper()
	artifact := candidate.Artifacts.Grader
	actualImplementation, err := digestLFNormalizedFile(filepath.Join(repositoryRoot, filepath.FromSlash(artifact.DigestImplementation)))
	if err != nil || actualImplementation != artifact.DigestImplementationHash {
		t.Fatalf("grader digest implementation = %s, %v", actualImplementation, err)
	}
	command := exec.Command("go", "run", "./cmd/corpusctl", "grader-digest", "--repo-root", "../../..")
	command.Dir = filepath.Join(repositoryRoot, "experiments", "frontier-v1", "corpusctl")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("compute grader digest: %v: %s", err, output)
	}
	if actual := strings.TrimSpace(string(output)); actual != artifact.Digest {
		t.Fatalf("grader digest = %s, candidate requires %s", actual, artifact.Digest)
	}
}

func TestCompletedArtifactFreezeGateState(t *testing.T) {
	var freeze struct {
		Status string `json:"status"`
		Gates  struct {
			MayOpenValidation         bool   `json:"may_open_validation"`
			MayOpenHeldOut            bool   `json:"may_open_held_out"`
			ValidationExecutionStatus string `json:"validation_execution_status"`
		} `json:"gates"`
		Remaining []string `json:"remaining"`
	}
	readJSONForTest(t, filepath.Join("..", "pre-validation-artifacts.json"), &freeze)
	if freeze.Status != "complete" || freeze.Gates.MayOpenValidation || freeze.Gates.MayOpenHeldOut ||
		freeze.Gates.ValidationExecutionStatus != "indeterminate" || len(freeze.Remaining) != 0 {
		t.Fatalf("completed artifact freeze gate state = %#v", freeze)
	}
}
