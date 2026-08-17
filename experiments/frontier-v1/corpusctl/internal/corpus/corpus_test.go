package corpus

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAuthoringManifestIsDeterministicAndCurrent(t *testing.T) {
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate test source")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(source), "../../../../.."))
	const world = "experiments/frontier-v1/worlds/authoring.dev_repo.json"
	first, err := BuildManifest(root, "authoring", world)
	if err != nil {
		t.Fatal(err)
	}
	second, err := BuildManifest(root, "authoring", world)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("manifest generation is not deterministic")
	}
	committed, err := os.ReadFile(filepath.Join(root, "experiments/frontier-v1/manifests/authoring.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, committed) {
		t.Fatal("authoring manifest is stale")
	}
}

func TestExecutableAuthoringFixtures(t *testing.T) {
	root := repoRoot(t)
	for _, fixtureID := range []string{"authoring_family_002_v1", "authoring_family_004_v1"} {
		t.Run(fixtureID, func(t *testing.T) {
			document, err := loadDocument(root, "experiments/frontier-v1/fixtures/authoring/"+fixtureID+".json")
			if err != nil {
				t.Fatal(err)
			}
			var item fixture
			if err := document.Into(&item); err != nil {
				t.Fatal(err)
			}
			workspace := t.TempDir()
			for relative, content := range item.Files {
				target := filepath.Join(workspace, filepath.FromSlash(relative))
				if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			command := exec.Command("go", "test", "./...")
			command.Dir = workspace
			if output, err := command.CombinedOutput(); err != nil {
				t.Fatalf("materialized fixture is not executable: %v\n%s", err, output)
			}
		})
	}
}

func TestRunnableCaseRejectsClassLeakage(t *testing.T) {
	schemas, err := loadSchemas(repoRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	value := map[string]any{
		"case_id":         "authoring_0a10d1ec",
		"goal":            "Run the complete repository test suite and report whether it is green.",
		"sandbox_fixture": "sha256:90121ca6f1fe7e4851ec9e97a224201e905e8b9d360a8bda0cbe3633a499209c",
		"world_ref":       "sha256:afa4896d40a68177e86c51574e940ff16633076434acd57a16e00c22e787f1b0",
		"family_id":       "authoring_family_001",
	}
	if err := schemas.validateValue("case", value); err != nil {
		t.Fatalf("expected answer-free runnable case to validate: %v", err)
	}
	value["class"] = "direct"
	if err := schemas.validateValue("case", value); err == nil {
		t.Fatal("expected runnable case schema to reject evaluator-only class metadata")
	}
}
