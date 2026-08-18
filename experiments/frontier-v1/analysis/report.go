package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func renderNewReport(root, templatePath, outputPath string, report analysisReport) error {
	resolvedOutput, err := newOutputPath(root, outputPath)
	if err != nil {
		return err
	}
	resolvedTemplate := resolvePath(root, templatePath)
	if err := requireInside(root, resolvedTemplate); err != nil {
		return err
	}
	contents, err := os.ReadFile(resolvedTemplate)
	if err != nil {
		return err
	}
	parsed, err := template.New("report").Funcs(template.FuncMap{
		"decimal": func(value float64) string { return fmt.Sprintf("%.6f", value) },
		"optional": func(value *float64) string {
			if value == nil {
				return "not estimable"
			}
			return fmt.Sprintf("%.6f", *value)
		},
	}).Parse(string(contents))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(resolvedOutput), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(resolvedOutput, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()
	return parsed.Execute(file, report)
}
