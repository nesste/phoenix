package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/raoul/phoenix/internal/buildmanifest"
)

type stringList []string

func (values *stringList) String() string { return fmt.Sprint([]string(*values)) }

func (values *stringList) Set(value string) error {
	*values = append(*values, value)
	return nil
}

func main() {
	var verbs stringList
	var schemas stringList
	var rules stringList
	repoRoot := flag.String("repo-root", ".", "repository root")
	output := flag.String("output", "build/manifest.json", "manifest output path, relative to the repository root")
	executable := flag.String("executable", "bin/phoenix", "Phoenix executable path")
	world := flag.String("world-definition", "", "world definition path")
	weights := flag.String("active-weights", "", "active frontier weights path")
	goos := flag.String("goos", "", "target operating system")
	goarch := flag.String("goarch", "", "target architecture")
	cgoEnabled := flag.Bool("cgo-enabled", false, "whether the executable was built with CGO")
	flag.Var(&verbs, "verb", "registered verb implementation path; repeatable")
	flag.Var(&schemas, "schema", "schema path; repeatable")
	flag.Var(&rules, "rule", "authored rule path; repeatable")
	flag.Parse()

	manifest, err := buildmanifest.Build(*repoRoot, buildmanifest.Input{
		Executable:      *executable,
		RegisteredVerbs: verbs,
		Schemas:         schemas,
		WorldDefinition: *world,
		AuthoredRules:   rules,
		ActiveWeights:   *weights,
		GOOS:            *goos,
		GOARCH:          *goarch,
		CGOEnabled:      *cgoEnabled,
	})
	if err != nil {
		fatal(err)
	}
	encoded, err := buildmanifest.Encode(manifest)
	if err != nil {
		fatal(err)
	}

	target := *output
	if !filepath.IsAbs(target) {
		target = filepath.Join(*repoRoot, target)
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		fatal(fmt.Errorf("create output directory: %w", err))
	}
	if err := os.WriteFile(target, encoded, 0o644); err != nil {
		fatal(fmt.Errorf("write manifest: %w", err))
	}
	fmt.Println(manifest.Digest)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
