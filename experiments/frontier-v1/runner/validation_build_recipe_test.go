package main

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

func TestValidationBuildRecipeMatchesFrozenWorldBuildTarget(t *testing.T) {
	var freeze struct {
		Artifacts struct {
			WorldBuild struct {
				Manifest string `json:"world_build_manifest"`
				Target   struct {
					GOOS         string `json:"goos"`
					GOARCH       string `json:"goarch"`
					CGOEnabled   bool   `json:"cgo_enabled"`
					Version      string `json:"version"`
					Trimpath     bool   `json:"trimpath"`
					Buildvcs     bool   `json:"buildvcs"`
					EmptyBuildID bool   `json:"empty_buildid"`
				} `json:"target"`
			} `json:"world_definition_and_world_build_digest"`
		} `json:"artifacts"`
	}
	readJSONForTest(t, filepath.Join("..", "pre-validation-artifacts.json"), &freeze)
	target := freeze.Artifacts.WorldBuild.Target
	if !target.Trimpath || target.Buildvcs || !target.EmptyBuildID || target.CGOEnabled {
		t.Fatalf("frozen world-build target flags = %#v", target)
	}
	recipe := phoenixBuildRecipeForTranche("validation")
	if recipe.goos != target.GOOS || recipe.goarch != target.GOARCH || recipe.version != target.Version {
		t.Fatalf("validation build recipe = %#v, freeze target requires %#v", recipe, target)
	}
	var manifest struct {
		Build struct {
			Executable struct {
				Path string `json:"path"`
			} `json:"executable"`
		} `json:"build"`
	}
	readJSONForTest(t, filepath.Join("..", "..", "..", filepath.FromSlash(freeze.Artifacts.WorldBuild.Manifest)), &manifest)
	if recipe.output != manifest.Build.Executable.Path {
		t.Fatalf("validation build output = %q, frozen manifest records %q", recipe.output, manifest.Build.Executable.Path)
	}

	command := phoenixBuildCommand(".", "bin/phoenix", recipe)
	arguments := strings.Join(command.Args, " ")
	for _, required := range []string{"-trimpath", "-buildvcs=false", "-buildid= ", "-X main.version=dev"} {
		if !strings.Contains(arguments, required) {
			t.Fatalf("validation build command %q lacks %q", arguments, required)
		}
	}
	environment := strings.Join(command.Env, "\n")
	for _, required := range []string{"GOOS=linux", "GOARCH=amd64", "CGO_ENABLED=0", "GOTOOLCHAIN=go1.26.6"} {
		if !strings.Contains(environment, required) {
			t.Fatalf("validation build environment lacks %q", required)
		}
	}
}

func TestAuthoringBuildRecipeStaysOnHostToolchain(t *testing.T) {
	recipe := phoenixBuildRecipeForTranche("authoring")
	if recipe.goos != runtime.GOOS || recipe.goarch != runtime.GOARCH || recipe.version != "authoring" || recipe.output != "" {
		t.Fatalf("authoring build recipe = %#v", recipe)
	}
}

func TestValidationBuildReproducesFrozenWorldBuildDigest(t *testing.T) {
	if testing.Short() {
		t.Skip("cross-compiling Phoenix is skipped in short mode")
	}
	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	executable, digest, cleanup, err := buildPhoenix(root, phoenixBuildRecipeForTranche("validation"))
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	relative, err := filepath.Rel(root, executable)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.ToSlash(relative) != "bin/phoenix" {
		t.Fatalf("validation executable = %s, freeze manifest requires bin/phoenix", filepath.ToSlash(relative))
	}
	var freeze struct {
		Artifacts struct {
			WorldBuild struct {
				Digest string `json:"world_build_digest"`
			} `json:"world_definition_and_world_build_digest"`
		} `json:"artifacts"`
	}
	readJSONForTest(t, filepath.Join("..", "pre-validation-artifacts.json"), &freeze)
	if digest != freeze.Artifacts.WorldBuild.Digest {
		t.Fatalf("live world-build %s does not reproduce frozen world-build %s", digest, freeze.Artifacts.WorldBuild.Digest)
	}
}

func TestGate1AValidationBuildRecipeCandidateMatchesWorkingTree(t *testing.T) {
	var candidate struct {
		V           int    `json:"v"`
		Status      string `json:"status"`
		SourceLimit string `json:"source_limit"`
		BaseCommit  string `json:"base_commit"`
		Frozen      bool   `json:"frozen"`
		Gates       struct {
			MayOpenValidation bool `json:"may_open_validation"`
			MayOpenHeldOut    bool `json:"may_open_held_out"`
		} `json:"gates"`
		SetDigest          string            `json:"artifact_set_lf_normalized_utf8_sha256"`
		Files              map[string]string `json:"files"`
		ExecutionPerformed bool              `json:"execution_performed"`
	}
	readJSONForTest(t, filepath.Join("..", "artifacts", "gate-1a-validation-build-recipe-candidate.json"), &candidate)
	if candidate.V != 1 || candidate.Status != "review_candidate" || candidate.SourceLimit != "public_validation_inputs_only" ||
		candidate.BaseCommit != "0ca98a1ae3f1ca1dedc7db86e5a676b0f1d8f7cc" || candidate.Frozen || candidate.Gates.MayOpenValidation ||
		candidate.Gates.MayOpenHeldOut || candidate.ExecutionPerformed {
		t.Fatalf("validation build recipe candidate state = %#v", candidate)
	}
	if len(candidate.Files) != 20 {
		t.Fatalf("validation build recipe candidate file count = %d", len(candidate.Files))
	}
	root := filepath.Join("..", "..", "..")
	paths := make([]string, 0, len(candidate.Files))
	for path := range candidate.Files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	var identity strings.Builder
	for _, path := range paths {
		actual, err := digestLFNormalizedFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			t.Fatal(err)
		}
		if actual != candidate.Files[path] {
			t.Fatalf("candidate file %s digest = %s, want %s", path, actual, candidate.Files[path])
		}
		fmt.Fprintf(&identity, "%s\t%s\n", path, actual)
	}
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(identity.String())))
	if digest != candidate.SetDigest {
		t.Fatalf("validation build recipe candidate set digest = %s, want %s", digest, candidate.SetDigest)
	}
}
