package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gowebpki/jcs"
)

func loadCase(repositoryRoot, caseID string) (runnableCase, error) {
	return loadCaseForTranche(repositoryRoot, "authoring", caseID)
}

func loadCaseForTranche(repositoryRoot, tranche, caseID string) (runnableCase, error) {
	prefix := tranche + "_"
	if tranche != "authoring" && tranche != "validation" || !strings.HasPrefix(caseID, prefix) {
		return runnableCase{}, fmt.Errorf("runner rejects case %q for tranche %q", caseID, tranche)
	}
	path := filepath.Join(repositoryRoot, "experiments", "frontier-v1", "corpus", tranche, caseID+".json")
	var item runnableCase
	if err := decodeStrict(path, &item); err != nil {
		return runnableCase{}, err
	}
	if item.CaseID != caseID {
		return runnableCase{}, fmt.Errorf("case file identity %q does not match %q", item.CaseID, caseID)
	}
	return item, nil
}

func authoringCaseIDs(repositoryRoot string) ([]string, error) {
	return caseIDsForTranche(repositoryRoot, "authoring")
}

func caseIDsForTranche(repositoryRoot, tranche string) ([]string, error) {
	if tranche != "authoring" && tranche != "validation" {
		return nil, fmt.Errorf("unsupported tranche %q", tranche)
	}
	pattern := filepath.Join(repositoryRoot, "experiments", "frontier-v1", "corpus", tranche, tranche+"_*.json")
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(paths))
	for _, path := range paths {
		ids = append(ids, strings.TrimSuffix(filepath.Base(path), ".json"))
	}
	sort.Strings(ids)
	return ids, nil
}

func loadFixtureForTranche(repositoryRoot, tranche string, item runnableCase) (fixture, error) {
	if tranche != "authoring" && tranche != "validation" {
		return fixture{}, fmt.Errorf("unsupported tranche %q", tranche)
	}
	pattern := filepath.Join(repositoryRoot, "experiments", "frontier-v1", "fixtures", tranche, "*.json")
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return fixture{}, err
	}
	for _, path := range paths {
		digest, err := digestJSONFile(path)
		if err != nil {
			return fixture{}, err
		}
		if digest != item.SandboxFixture {
			continue
		}
		var found fixture
		if err := decodeStrict(path, &found); err != nil {
			return fixture{}, err
		}
		if found.FamilyID != item.FamilyID {
			return fixture{}, fmt.Errorf("case family %s does not match fixture family %s", item.FamilyID, found.FamilyID)
		}
		return found, nil
	}
	return fixture{}, fmt.Errorf("%s fixture %s was not found", tranche, item.SandboxFixture)
}

func materializeFixture(target string, item fixture) error {
	uncommitted := make(map[string]struct{}, len(item.Uncommitted))
	for _, path := range item.Uncommitted {
		if err := validateRelativePath(path); err != nil {
			return err
		}
		uncommitted[path] = struct{}{}
	}
	for path, contents := range item.Files {
		if err := validateRelativePath(path); err != nil {
			return err
		}
		if _, pending := uncommitted[path]; pending {
			continue
		}
		if err := writeFixtureFile(target, path, contents); err != nil {
			return err
		}
	}
	for _, args := range [][]string{
		{"init", "-q"},
		{"config", "user.name", "Phoenix Runner"},
		{"config", "user.email", "runner@phoenix.invalid"},
		{"add", "-A"},
		{"commit", "-q", "-m", "fixture"},
	} {
		if err := runGit(target, args...); err != nil {
			return err
		}
	}
	for _, path := range item.Uncommitted {
		if err := writeFixtureFile(target, path, item.Files[path]); err != nil {
			return err
		}
	}
	return nil
}

func captureEndState(root string) (endState, error) {
	state := endState{Files: []fileState{}}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() && entry.Name() == ".git" {
			return fs.SkipDir
		}
		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		state.Files = append(state.Files, fileState{Path: filepath.ToSlash(relative), Present: true, Content: string(contents)})
		return nil
	})
	if err != nil {
		return endState{}, err
	}
	sort.Slice(state.Files, func(i, j int) bool { return state.Files[i].Path < state.Files[j].Path })
	return state, nil
}

func writeFixtureFile(root, relative, contents string) error {
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(contents), 0o644)
}

func validateRelativePath(path string) error {
	clean := filepath.Clean(filepath.FromSlash(path))
	if path == "" || filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return fmt.Errorf("fixture path %q escapes the sandbox", path)
	}
	return nil
}

func runGit(dir string, args ...string) error {
	command := exec.Command("git", args...)
	command.Dir = dir
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, bytes.TrimSpace(output))
	}
	return nil
}

func decodeStrict(path string, target any) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("decode %s: trailing JSON data", path)
	}
	return nil
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
