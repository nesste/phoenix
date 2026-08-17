package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/raoul/phoenix/experiments/frontier-v1/corpusctl/internal/corpus"
)

const defaultWorldSource = "experiments/frontier-v1/worlds/authoring.dev_repo.json"

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "digest":
		digest(os.Args[2:])
	case "grade":
		grade(os.Args[2:])
	case "grader-digest":
		graderDigest(os.Args[2:])
	case "validate":
		manifest(os.Args[2:], false)
	case "manifest":
		manifest(os.Args[2:], true)
	case "seal":
		seal(os.Args[2:])
	default:
		usage()
	}
}

func grade(arguments []string) {
	flags := flag.NewFlagSet("grade", flag.ExitOnError)
	repoRoot := flags.String("repo-root", "../../..", "path to the repository root")
	labelPath := flags.String("label", "", "repository-relative label path")
	trialPath := flags.String("trial", "", "repository-relative trial-record path")
	flags.Parse(arguments)
	if *labelPath == "" || *trialPath == "" {
		log.Fatal("grade requires --label and --trial")
	}
	root, err := filepath.Abs(*repoRoot)
	if err != nil {
		log.Fatal(err)
	}
	_, label, err := corpus.LoadLabel(root, filepath.ToSlash(*labelPath))
	if err != nil {
		log.Fatal(err)
	}
	_, trial, err := corpus.LoadTrial(root, filepath.ToSlash(*trialPath))
	if err != nil {
		log.Fatal(err)
	}
	result, err := corpus.Grade(label, trial)
	if err != nil {
		log.Fatal(err)
	}
	data, err := corpus.EncodeGradeResult(root, result)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = os.Stdout.Write(data)
	switch result.Status {
	case corpus.StatusPass:
		return
	case corpus.StatusFail:
		os.Exit(1)
	default:
		os.Exit(2)
	}
}

func graderDigest(arguments []string) {
	flags := flag.NewFlagSet("grader-digest", flag.ExitOnError)
	repoRoot := flags.String("repo-root", "../../..", "path to the repository root")
	flags.Parse(arguments)
	root, err := filepath.Abs(*repoRoot)
	if err != nil {
		log.Fatal(err)
	}
	digest, err := corpus.GraderDigest(root)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(digest)
}

func digest(arguments []string) {
	flags := flag.NewFlagSet("digest", flag.ExitOnError)
	flags.Parse(arguments)
	if flags.NArg() != 1 {
		log.Fatal("usage: corpusctl digest <json-file>")
	}
	value, err := corpus.DigestFile(flags.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(value)
}

func manifest(arguments []string, write bool) {
	flags := flag.NewFlagSet("manifest", flag.ExitOnError)
	repoRoot := flags.String("repo-root", "../../..", "path to the repository root")
	tranche := flags.String("tranche", "authoring", "corpus tranche")
	worldSource := flags.String("world-source", defaultWorldSource, "repository-relative world definition")
	writeFile := flags.Bool("write", false, "write the generated manifest")
	flags.Parse(arguments)
	root, err := filepath.Abs(*repoRoot)
	if err != nil {
		log.Fatal(err)
	}
	data, err := corpus.BuildManifest(root, *tranche, *worldSource)
	if err != nil {
		log.Fatal(err)
	}
	if write && *writeFile {
		path, err := corpus.WriteManifest(root, *tranche, *worldSource)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("wrote %s\n", path)
		return
	}
	if write {
		_, _ = os.Stdout.Write(data)
		return
	}
	target := filepath.Join(root, filepath.FromSlash("experiments/frontier-v1/manifests/"+*tranche+".json"))
	existing, err := os.ReadFile(target)
	if err != nil {
		log.Fatal(err)
	}
	if !bytes.Equal(existing, data) {
		log.Fatalf("%s is stale; regenerate it with the manifest command", target)
	}
	fmt.Printf("valid %s\n", filepath.ToSlash(target))
}

func seal(arguments []string) {
	flags := flag.NewFlagSet("seal", flag.ExitOnError)
	repoRoot := flags.String("repo-root", "../../..", "path to the repository root")
	tranche := flags.String("tranche", "", "validation or held_out")
	worldSource := flags.String("world-source", "", "repository-relative world definition")
	labelDigests := flags.String("label-digests", "", "repository-relative withheld-label digest registry")
	writeFile := flags.Bool("write", false, "write the generated sealed manifest")
	flags.Parse(arguments)
	if *tranche == "" || *worldSource == "" || *labelDigests == "" {
		log.Fatal("seal requires --tranche, --world-source, and --label-digests")
	}
	root, err := filepath.Abs(*repoRoot)
	if err != nil {
		log.Fatal(err)
	}
	if *writeFile {
		path, err := corpus.WriteSealedManifest(root, *tranche, *worldSource, *labelDigests)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("wrote %s\n", path)
		return
	}
	data, err := corpus.BuildSealedManifest(root, *tranche, *worldSource, *labelDigests)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = os.Stdout.Write(data)
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: corpusctl <digest|grade|grader-digest|manifest|seal|validate> [options]")
	os.Exit(2)
}
