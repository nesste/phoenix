package buildmanifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildHashesAndSortsArtifacts(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "bin/phoenix", "daemon")
	writeFixture(t, root, "spec/z.json", "z")
	writeFixture(t, root, "spec/a.json", "a")

	manifest, err := Build(root, Input{
		Executable: "bin/phoenix",
		Schemas:    []string{"spec/z.json", "spec/a.json"},
		GOOS:       "linux",
		GOARCH:     "amd64",
	})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	if !strings.HasPrefix(manifest.Digest, "sha256:") || len(manifest.Digest) != len("sha256:")+64 {
		t.Fatalf("manifest digest = %q, want sha256 digest", manifest.Digest)
	}
	if got, want := manifest.Build.Executable.Path, "bin/phoenix"; got != want {
		t.Fatalf("executable path = %q, want %q", got, want)
	}
	if got, want := len(manifest.Build.Schemas), 2; got != want {
		t.Fatalf("schema count = %d, want %d", got, want)
	}
	if got, want := manifest.Build.Schemas[0].Path, "spec/a.json"; got != want {
		t.Fatalf("first schema = %q, want %q", got, want)
	}
	if manifest.Build.RegisteredVerbs == nil || manifest.Build.AuthoredRules == nil {
		t.Fatal("empty component sets must encode as arrays, not null")
	}
	if manifest.Build.WorldDefinition != nil || manifest.Build.ActiveWeights != nil {
		t.Fatal("not-yet-built singleton components must encode as null")
	}
}

func TestBuildRejectsArtifactOutsideRepository(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside")
	if err := os.WriteFile(outside, []byte("outside"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := Build(root, Input{Executable: outside, GOOS: "linux", GOARCH: "amd64"})
	if err == nil || !strings.Contains(err.Error(), "outside repository root") {
		t.Fatalf("Build error = %v, want outside-root rejection", err)
	}
}

func writeFixture(t *testing.T, root, name, contents string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}
