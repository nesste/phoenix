package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	pinnedRuntimeVersion = "2.1.229"
	pinnedModel          = "claude-sonnet-5"
	pinnedSystemPrompt   = "Use only the configured Phoenix tool. Treat the supplied handles as live. Act directly from the goal and handle type. Follow the user request exactly."
)

type claudeDriver struct {
	executable string
}

func (driver claudeDriver) Verify() error {
	command := exec.Command(driver.executable, "--version")
	output, err := command.Output()
	if err != nil {
		return fmt.Errorf("read pinned runtime version: %w", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(string(output)), pinnedRuntimeVersion+" ") {
		return fmt.Errorf("runtime version is %q, want %s", strings.TrimSpace(string(output)), pinnedRuntimeVersion)
	}
	return nil
}

func (driver claudeDriver) Run(request runtimeRequest) (runtimeResult, error) {
	configPath := filepath.Join(filepath.Dir(request.EpisodePath), "mcp.json")
	rootPath := filepath.Join(filepath.Dir(request.EpisodePath), "roots.json")
	if err := writeJSON(rootPath, request.Roots); err != nil {
		return runtimeResult{}, err
	}
	mcpConfig := map[string]any{"mcpServers": map[string]any{
		"phoenix": map[string]any{
			"type": "stdio", "command": request.PhoenixPath,
			"args": []string{
				"serve", "--stdio", "--world", request.WorldPath, "--world-schema", request.SchemaPath,
				"--episode-db", request.EpisodePath, "--root-refs", rootPath, "--world-build", request.WorldBuild,
			},
			"env": map[string]string{},
		},
	}}
	if err := writeJSON(configPath, mcpConfig); err != nil {
		return runtimeResult{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), request.Timeout)
	defer cancel()
	args := []string{
		"-p", runtimePrompt(request.Goal, request.Roots),
		"--model", pinnedModel,
		"--effort", "low",
		"--max-turns", "12",
		"--max-budget-usd", request.BudgetUSD,
		"--no-session-persistence",
		"--output-format", "stream-json",
		"--verbose",
		"--strict-mcp-config",
		"--mcp-config", configPath,
		"--tools", "",
		"--allowedTools", "mcp__phoenix__act",
		"--system-prompt", pinnedSystemPrompt,
	}
	command := exec.CommandContext(ctx, driver.executable, args...)
	command.Dir = request.Sandbox
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if ctx.Err() != nil {
		return runtimeResult{}, fmt.Errorf("runtime timed out after %s", request.Timeout)
	}
	if err != nil {
		return runtimeResult{}, fmt.Errorf("runtime failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	final, err := parseRuntimeOutput(stdout.Bytes())
	if err != nil {
		return runtimeResult{}, err
	}
	return runtimeResult{FinalMessage: final, RawOutput: stdout.Bytes()}, nil
}

func runtimePrompt(goal string, roots map[string]string) string {
	names := []string{"repo", "tests", "git", "episodes"}
	var prompt strings.Builder
	prompt.WriteString(goal)
	prompt.WriteString("\n\nLive handles:\n")
	for _, name := range names {
		fmt.Fprintf(&prompt, "- %s: %s\n", name, roots[name])
	}
	return prompt.String()
}

func parseRuntimeOutput(output []byte) (string, error) {
	var final string
	scanner := bufio.NewScanner(bytes.NewReader(output))
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		var event struct {
			Type   string `json:"type"`
			Result string `json:"result"`
		}
		if json.Unmarshal(scanner.Bytes(), &event) == nil && event.Type == "result" {
			final = event.Result
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("read runtime stream: %w", err)
	}
	if final == "" {
		return "", fmt.Errorf("runtime stream contains no final result")
	}
	return final, nil
}

type runtimeInitEvent struct {
	Type              string             `json:"type"`
	Subtype           string             `json:"subtype"`
	Tools             []string           `json:"tools"`
	MCPServers        []runtimeMCPServer `json:"mcp_servers"`
	Model             string             `json:"model"`
	PermissionMode    string             `json:"permissionMode"`
	APIKeySource      string             `json:"apiKeySource"`
	ClaudeCodeVersion string             `json:"claude_code_version"`
}

type runtimeMCPServer struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func sanitizeRuntimeOutput(output []byte) ([]byte, error) {
	var sanitized bytes.Buffer
	scanner := bufio.NewScanner(bytes.NewReader(output))
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		line := append([]byte(nil), scanner.Bytes()...)
		var event runtimeInitEvent
		if json.Unmarshal(line, &event) == nil && event.Type == "system" && event.Subtype == "init" {
			encoded, err := json.Marshal(event)
			if err != nil {
				return nil, err
			}
			line = encoded
		}
		sanitized.Write(line)
		sanitized.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("sanitize runtime stream: %w", err)
	}
	return sanitized.Bytes(), nil
}

type corpusGrader struct {
	goExecutable string
}

func (grader corpusGrader) Grade(repositoryRoot, labelPath, trialPath string) ([]byte, error) {
	workdir := filepath.Join(repositoryRoot, "experiments", "frontier-v1", "corpusctl")
	command := exec.Command(grader.goExecutable, "run", "./cmd/corpusctl", "grade",
		"--repo-root", "../../..", "--label", labelPath, "--trial", trialPath,
	)
	command.Dir = workdir
	output, err := command.CombinedOutput()
	if err == nil {
		return output, nil
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) && (exitError.ExitCode() == 1 || exitError.ExitCode() == 2) {
		return trimGoRunExit(output), nil
	}
	return nil, fmt.Errorf("run corpus grader: %w: %s", err, strings.TrimSpace(string(output)))
}

func trimGoRunExit(output []byte) []byte {
	marker := []byte("\nexit status ")
	if index := bytes.LastIndex(output, marker); index >= 0 {
		return append([]byte(nil), output[:index+1]...)
	}
	return output
}

func writeJSON(path string, value any) error {
	contents, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(contents, '\n'), 0o600)
}
