package corpus

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate test source")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(source), "../../../../.."))
}

func writeTrialFixture(t *testing.T, root, body string) string {
	t.Helper()
	// A temp dir under root, not the OS temp dir, so filepath.Rel below
	// always succeeds even when the OS temp dir lives on another volume
	// (as it commonly does on Windows).
	dir, err := os.MkdirTemp(root, "trial_test_")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	path := filepath.Join(dir, "trial.json")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	relative, err := filepath.Rel(root, path)
	if err != nil {
		t.Fatal(err)
	}
	return filepath.ToSlash(relative)
}

func validTrialJSON() string {
	digest := "sha256:" + repeat64("a")
	return `{
  "v": 1,
  "case_id": "authoring_direct_001",
  "world_build": "` + digest + `",
  "acts": [
    { "seq": 0, "handle": "repo", "handle_type": "tests", "verb": "run", "args": {}, "status": "ok", "command": ["go", "test", "./..."], "exit_code": 1, "output": "FAIL" }
  ],
  "final_message": "The suite failed.",
  "end_state": { "files": [{ "path": "go.sum", "present": true, "content": "x" }] }
}`
}

func repeat64(character string) string {
	result := ""
	for len(result) < 64 {
		result += character
	}
	return result[:64]
}

func TestLoadTrialValid(t *testing.T) {
	root := repoRoot(t)
	path := writeTrialFixture(t, root, validTrialJSON())
	_, trial, err := LoadTrial(root, path)
	if err != nil {
		t.Fatal(err)
	}
	if trial.CaseID != "authoring_direct_001" {
		t.Fatalf("case_id = %q", trial.CaseID)
	}
	if len(trial.Acts) != 1 || trial.Acts[0].HandleType != "tests" || trial.Acts[0].Verb != "run" {
		t.Fatalf("unexpected acts: %+v", trial.Acts)
	}
}

func TestLoadTrialEmptyActsIsLegal(t *testing.T) {
	root := repoRoot(t)
	digest := "sha256:" + repeat64("b")
	body := `{
  "v": 1,
  "case_id": "authoring_direct_001",
  "world_build": "` + digest + `",
  "acts": [],
  "final_message": "done",
  "end_state": { "files": [] }
}`
	path := writeTrialFixture(t, root, body)
	if _, _, err := LoadTrial(root, path); err != nil {
		t.Fatal(err)
	}
}

func TestLoadTrialRejectsUnknownTopLevelField(t *testing.T) {
	root := repoRoot(t)
	digest := "sha256:" + repeat64("c")
	body := `{
  "v": 1,
  "case_id": "authoring_direct_001",
  "world_build": "` + digest + `",
  "acts": [],
  "final_message": "done",
  "end_state": { "files": [] },
  "expected_outcome": {}
}`
	path := writeTrialFixture(t, root, body)
	if _, _, err := LoadTrial(root, path); err == nil {
		t.Fatal("expected schema rejection of unknown top-level field")
	}
}

func TestLoadTrialRejectsUnknownActField(t *testing.T) {
	root := repoRoot(t)
	digest := "sha256:" + repeat64("d")
	body := `{
  "v": 1,
  "case_id": "authoring_direct_001",
  "world_build": "` + digest + `",
  "acts": [{ "seq": 0, "handle": "repo", "handle_type": "tests", "verb": "run", "args": {}, "status": "ok", "unexpected": true }],
  "final_message": "done",
  "end_state": { "files": [] }
}`
	path := writeTrialFixture(t, root, body)
	if _, _, err := LoadTrial(root, path); err == nil {
		t.Fatal("expected schema rejection of unknown act field")
	}
}

func TestLoadTrialRejectsDuplicateKeys(t *testing.T) {
	root := repoRoot(t)
	digest := "sha256:" + repeat64("e")
	body := `{
  "v": 1,
  "case_id": "authoring_direct_001",
  "case_id": "authoring_direct_002",
  "world_build": "` + digest + `",
  "acts": [],
  "final_message": "done",
  "end_state": { "files": [] }
}`
	path := writeTrialFixture(t, root, body)
	if _, _, err := LoadTrial(root, path); err == nil {
		t.Fatal("expected duplicate key rejection")
	}
}

func TestLoadTrialRejectsOutOfRangeExitCode(t *testing.T) {
	root := repoRoot(t)
	digest := "sha256:" + repeat64("f")
	body := `{
  "v": 1,
  "case_id": "authoring_direct_001",
  "world_build": "` + digest + `",
  "acts": [{ "seq": 0, "handle": "repo", "handle_type": "tests", "verb": "run", "args": {}, "status": "ok", "exit_code": 999 }],
  "final_message": "done",
  "end_state": { "files": [] }
}`
	path := writeTrialFixture(t, root, body)
	if _, _, err := LoadTrial(root, path); err == nil {
		t.Fatal("expected exit_code range rejection")
	}
}

func TestLoadTrialRejectsNonContiguousSequenceAndDuplicateFiles(t *testing.T) {
	root := repoRoot(t)
	digest := "sha256:" + repeat64("1")
	nonContiguous := `{
  "v": 1,
  "case_id": "authoring_direct_001",
  "world_build": "` + digest + `",
  "acts": [{ "seq": 2, "handle": "h_tests", "handle_type": "tests", "verb": "run", "args": {}, "status": "ok" }],
  "final_message": "done",
  "end_state": { "files": [] }
}`
	if _, _, err := LoadTrial(root, writeTrialFixture(t, root, nonContiguous)); err == nil {
		t.Fatal("expected non-contiguous sequence rejection")
	}
	duplicateFiles := `{
  "v": 1,
  "case_id": "authoring_direct_001",
  "world_build": "` + digest + `",
  "acts": [],
  "final_message": "done",
  "end_state": { "files": [{"path":"a.txt","present":true},{"path":"a.txt","present":false}] }
}`
	if _, _, err := LoadTrial(root, writeTrialFixture(t, root, duplicateFiles)); err == nil {
		t.Fatal("expected duplicate end-state path rejection")
	}
}

func TestGradeAuthoringLabelEndToEnd(t *testing.T) {
	root := repoRoot(t)
	labelPath := "experiments/frontier-v1/labels/authoring/authoring_direct_001.json"
	_, label, err := LoadLabel(root, labelPath)
	if err != nil {
		t.Fatal(err)
	}
	digest := "sha256:" + repeat64("2")
	body := `{
  "v": 1,
  "case_id": "authoring_direct_001",
  "world_build": "` + digest + `",
  "acts": [{ "seq": 0, "handle": "h_tests", "handle_type": "tests", "verb": "run", "args": {}, "status": "ok", "output": "The complete test suite failed." }],
  "final_message": "The complete test suite failed and is not green.",
  "end_state": { "files": [] }
}`
	_, trial, err := LoadTrial(root, writeTrialFixture(t, root, body))
	if err != nil {
		t.Fatal(err)
	}
	result, err := Grade(label, trial)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != StatusPass {
		t.Fatalf("expected pass, got %+v", result)
	}
	if _, err := EncodeGradeResult(root, result); err != nil {
		t.Fatal(err)
	}
}
