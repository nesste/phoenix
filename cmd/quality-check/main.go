// quality-check contains small cross-platform checks used by the Makefile.
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: quality-check <format|no-output> [arguments]")
	}
	switch args[0] {
	case "format":
		if len(args) == 1 {
			return fmt.Errorf("format requires at least one path")
		}
		unformatted, err := unformattedFiles(args[1:])
		if err != nil {
			return err
		}
		if len(unformatted) == 0 {
			return nil
		}
		fmt.Fprintln(stdout, "gofmt required:")
		for _, name := range unformatted {
			fmt.Fprintln(stdout, name)
		}
		return fmt.Errorf("format check failed")
	case "no-output":
		command := args[1:]
		if len(command) > 0 && command[0] == "--" {
			command = command[1:]
		}
		if len(command) == 0 {
			return fmt.Errorf("no-output requires a command")
		}
		return requireNoOutput(command, stdout, stderr)
	default:
		return fmt.Errorf("unknown quality check %q", args[0])
	}
}

func unformattedFiles(paths []string) ([]string, error) {
	var unformatted []string
	for _, root := range paths {
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
				return nil
			}
			contents, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			formatted, err := format.Source(contents)
			if err != nil {
				return fmt.Errorf("format %s: %w", path, err)
			}
			if !bytes.Equal(contents, formatted) {
				unformatted = append(unformatted, filepath.ToSlash(path))
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return unformatted, nil
}

func requireNoOutput(arguments []string, stdout, stderr io.Writer) error {
	command := exec.Command(arguments[0], arguments[1:]...)
	var output bytes.Buffer
	command.Stdout = &output
	command.Stderr = stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", arguments[0], err)
	}
	if len(bytes.TrimSpace(output.Bytes())) == 0 {
		return nil
	}
	_, _ = io.Copy(stdout, &output)
	return fmt.Errorf("command produced findings")
}
