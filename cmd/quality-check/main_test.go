package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestUnformattedFiles(t *testing.T) {
	directory := t.TempDir()
	formatted := filepath.Join(directory, "formatted.go")
	unformatted := filepath.Join(directory, "unformatted.go")
	if err := os.WriteFile(formatted, []byte("package fixture\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(unformatted, []byte("package fixture\nfunc f( ){ }\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	files, err := unformattedFiles([]string{directory})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(files), 1; got != want {
		t.Fatalf("unformatted count = %d, want %d", got, want)
	}
	if got, want := files[0], filepath.ToSlash(unformatted); got != want {
		t.Fatalf("unformatted file = %q, want %q", got, want)
	}
}

func TestRunRejectsMissingCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := run([]string{"no-output"}, &stdout, &stderr); err == nil {
		t.Fatal("run(no-output) error = nil, want an error")
	}
}

func TestCompareFiles(t *testing.T) {
	directory := t.TempDir()
	first := filepath.Join(directory, "first")
	second := filepath.Join(directory, "second")
	if err := os.WriteFile(first, []byte("same\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(second, []byte("same\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := compareFiles(first, second); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(second, []byte("different\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := compareFiles(first, second); err == nil {
		t.Fatal("compareFiles accepted different files")
	}
}
