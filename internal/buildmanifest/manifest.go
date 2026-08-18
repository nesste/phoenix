// Package buildmanifest creates deterministic, content-addressed descriptions
// of complete Phoenix world builds.
package buildmanifest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const Version = 1

type Input struct {
	Executable      string
	RegisteredVerbs []string
	Schemas         []string
	WorldDefinition string
	AuthoredRules   []string
	ActiveWeights   string
	GOOS            string
	GOARCH          string
	CGOEnabled      bool
}

type Artifact struct {
	Path   string `json:"path"`
	Digest string `json:"digest"`
}

type BuildRecord struct {
	GOOS            string     `json:"goos"`
	GOARCH          string     `json:"goarch"`
	CGOEnabled      bool       `json:"cgo_enabled"`
	Executable      Artifact   `json:"executable"`
	RegisteredVerbs []Artifact `json:"registered_verbs"`
	Schemas         []Artifact `json:"schemas"`
	WorldDefinition *Artifact  `json:"world_definition"`
	AuthoredRules   []Artifact `json:"authored_rules"`
	ActiveWeights   *Artifact  `json:"active_weights"`
}

type Manifest struct {
	V      int         `json:"v"`
	Digest string      `json:"digest"`
	Build  BuildRecord `json:"build"`
}

func Build(repoRoot string, input Input) (Manifest, error) {
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return Manifest{}, fmt.Errorf("resolve repository root: %w", err)
	}
	if input.Executable == "" {
		return Manifest{}, fmt.Errorf("executable is required")
	}
	if input.GOOS == "" || input.GOARCH == "" {
		return Manifest{}, fmt.Errorf("goos and goarch are required")
	}

	executable, err := hashArtifact(root, input.Executable)
	if err != nil {
		return Manifest{}, fmt.Errorf("executable: %w", err)
	}
	verbs, err := hashArtifacts(root, input.RegisteredVerbs)
	if err != nil {
		return Manifest{}, fmt.Errorf("registered verbs: %w", err)
	}
	schemas, err := hashArtifacts(root, input.Schemas)
	if err != nil {
		return Manifest{}, fmt.Errorf("schemas: %w", err)
	}
	rules, err := hashArtifacts(root, input.AuthoredRules)
	if err != nil {
		return Manifest{}, fmt.Errorf("authored rules: %w", err)
	}
	world, err := hashOptionalArtifact(root, input.WorldDefinition)
	if err != nil {
		return Manifest{}, fmt.Errorf("world definition: %w", err)
	}
	weights, err := hashOptionalArtifact(root, input.ActiveWeights)
	if err != nil {
		return Manifest{}, fmt.Errorf("active weights: %w", err)
	}

	record := BuildRecord{
		GOOS:            input.GOOS,
		GOARCH:          input.GOARCH,
		CGOEnabled:      input.CGOEnabled,
		Executable:      executable,
		RegisteredVerbs: verbs,
		Schemas:         schemas,
		WorldDefinition: world,
		AuthoredRules:   rules,
		ActiveWeights:   weights,
	}
	canonical, err := json.Marshal(record)
	if err != nil {
		return Manifest{}, fmt.Errorf("encode build record: %w", err)
	}

	return Manifest{V: Version, Digest: digest(canonical), Build: record}, nil
}

func Encode(manifest Manifest) ([]byte, error) {
	encoded, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode manifest: %w", err)
	}
	return append(encoded, '\n'), nil
}

func hashArtifacts(root string, names []string) ([]Artifact, error) {
	artifacts := make([]Artifact, 0, len(names))
	for _, name := range names {
		artifact, err := hashArtifact(root, name)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, artifact)
	}
	sort.Slice(artifacts, func(i, j int) bool { return artifacts[i].Path < artifacts[j].Path })
	return artifacts, nil
}

func hashOptionalArtifact(root, name string) (*Artifact, error) {
	if name == "" {
		return nil, nil
	}
	artifact, err := hashArtifact(root, name)
	if err != nil {
		return nil, err
	}
	return &artifact, nil
}

func hashArtifact(root, name string) (Artifact, error) {
	path := name
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return Artifact{}, fmt.Errorf("resolve %q: %w", name, err)
	}
	relative, err := filepath.Rel(root, abs)
	if err != nil {
		return Artifact{}, fmt.Errorf("relativize %q: %w", name, err)
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
		return Artifact{}, fmt.Errorf("artifact %q is outside repository root", name)
	}

	contents, err := os.ReadFile(abs)
	if err != nil {
		return Artifact{}, fmt.Errorf("read %q: %w", filepath.ToSlash(relative), err)
	}
	return Artifact{Path: filepath.ToSlash(relative), Digest: digest(contents)}, nil
}

func digest(contents []byte) string {
	sum := sha256.Sum256(contents)
	return "sha256:" + hex.EncodeToString(sum[:])
}
