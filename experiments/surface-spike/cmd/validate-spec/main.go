package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type validation struct {
	schema   string
	instance string
}

func main() {
	repoRoot := flag.String("repo-root", "../..", "path to the Phoenix repository root")
	flag.Parse()

	root, err := filepath.Abs(*repoRoot)
	if err != nil {
		log.Fatal(err)
	}
	validations := []validation{
		{schema: "spec/result.schema.json", instance: "spec/examples/result.ok.json"},
		{schema: "spec/result.schema.json", instance: "spec/examples/result.fail.json"},
		{schema: "spec/result.schema.json", instance: "spec/examples/result.refused.json"},
		{schema: "spec/result.schema.json", instance: "spec/examples/result.refused_no_alternative.json"},
		{schema: "spec/result.schema.json", instance: "spec/examples/result.absent.json"},
		{schema: "spec/world.schema.json", instance: "spec/examples/world.dev_repo.json"},
	}
	for _, item := range validations {
		if err := validate(root, item); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("valid %s against %s\n", item.instance, item.schema)
	}
}

func validate(root string, item validation) error {
	schemaPath := filepath.Join(root, filepath.FromSlash(item.schema))
	instancePath := filepath.Join(root, filepath.FromSlash(item.instance))
	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile(schemaPath)
	if err != nil {
		return fmt.Errorf("compile %s: %w", item.schema, err)
	}
	file, err := os.Open(instancePath)
	if err != nil {
		return fmt.Errorf("open %s: %w", item.instance, err)
	}
	defer file.Close()
	instance, err := jsonschema.UnmarshalJSON(file)
	if err != nil {
		return fmt.Errorf("decode %s: %w", item.instance, err)
	}
	if err := schema.Validate(instance); err != nil {
		return fmt.Errorf("validate %s: %w", item.instance, err)
	}
	return nil
}
