package corpus

import (
	"bytes"
	"os"
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
