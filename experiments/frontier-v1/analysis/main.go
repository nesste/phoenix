// Command analysis computes the frozen frontier-v1 Phase 1 analysis from a
// scheduled-runner summary and its retained evidence.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const analysisSeed = uint64(20260817)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	flags := flag.NewFlagSet("frontier-v1-analysis", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	repositoryRoot := flags.String("repo-root", ".", "Phoenix repository root")
	summaryPath := flags.String("summary", "", "scheduled-summary.json to analyze")
	manifestPath := flags.String("manifest", "", "tranche manifest containing case class and family")
	outputPath := flags.String("output", "", "analysis JSON output path")
	reportPath := flags.String("report", "", "optional rendered Markdown report path")
	templatePath := flags.String("template", "experiments/frontier-v1/analysis/report-template.md.tmpl", "Markdown report template")
	if err := flags.Parse(args); err != nil || flags.NArg() != 0 {
		return 2
	}
	if *summaryPath == "" || *manifestPath == "" || *outputPath == "" {
		fmt.Fprintln(os.Stderr, "--summary, --manifest, and --output are required")
		return 2
	}
	root, err := filepath.Abs(*repositoryRoot)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	report, err := analyzeFiles(root, *summaryPath, *manifestPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err := writeNewJSON(root, *outputPath, report); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if *reportPath != "" {
		if err := renderNewReport(root, *templatePath, *reportPath, report); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}
	fmt.Printf("analysis: %s (%d claims)\n", report.Status, len(report.Claims))
	return 0
}

func writeNewJSON(root, path string, value any) error {
	resolved, err := newOutputPath(root, path)
	if err != nil {
		return err
	}
	contents, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(resolved), 0o755); err != nil {
		return err
	}
	return os.WriteFile(resolved, append(contents, '\n'), 0o600)
}

func newOutputPath(root, path string) (string, error) {
	resolved := resolvePath(root, path)
	if err := requireInside(root, resolved); err != nil {
		return "", err
	}
	if _, err := os.Stat(resolved); err == nil {
		return "", fmt.Errorf("refusing to overwrite %s", path)
	} else if !os.IsNotExist(err) {
		return "", err
	}
	return resolved, nil
}
