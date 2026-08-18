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
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	pinnedRuntimeVersion = "2.1.229"
	pinnedModel          = "claude-sonnet-5"
	baseSystemPrompt     = "Use only the configured tools. Follow the user request exactly."
	phoenixIntentPrompt  = "Treat the supplied Phoenix handles as live. Before the first executable act, send the complete goal as intent on any live handle, then execute a returned call."
)

var providerStatusPattern = regexp.MustCompile(`\b(?:429|5[0-9]{2})\b`)

func systemPromptForArm(arm string) (string, error) {
	switch arm {
	case "A", "B":
		return baseSystemPrompt, nil
	case "C", "D", "E":
		return baseSystemPrompt + " " + phoenixIntentPrompt, nil
	default:
		return "", fmt.Errorf("unsupported experiment arm %q", arm)
	}
}

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
	systemPrompt, err := systemPromptForArm(request.Arm)
	if err != nil {
		return runtimeResult{}, err
	}
	allowedTools := allowedToolsForArm(request.Arm, request.FlatToolNames)
	if len(allowedTools) == 0 {
		return runtimeResult{}, fmt.Errorf("arm %s has no allowed tools", request.Arm)
	}
	if request.Arm == "B" && request.ArmBDocument == "" {
		return runtimeResult{}, fmt.Errorf("arm B document is required")
	}
	args, err := prepareRuntimeInvocation(request, systemPrompt, allowedTools)
	if err != nil {
		return runtimeResult{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), request.Timeout)
	defer cancel()
	command := exec.CommandContext(ctx, driver.executable, args...)
	command.Dir = request.Sandbox
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	started := time.Now()
	err = command.Run()
	elapsed := time.Since(started)
	observation, parseErr := inspectRuntimeOutput(stdout.Bytes())
	observation.metrics.WallTimeMS = elapsed.Milliseconds()
	result := runtimeResult{
		FinalMessage: observation.finalMessage, RawOutput: stdout.Bytes(),
		FirstModelToken: observation.modelOutputSeen, Metrics: observation.metrics,
	}
	return finishRuntimeInvocation(result, observation, parseErr, err, ctx.Err(), stderr.String(), request.Timeout)
}

func prepareRuntimeInvocation(request runtimeRequest, systemPrompt string, allowedTools []string) ([]string, error) {
	configPath := filepath.Join(filepath.Dir(request.EpisodePath), "mcp.json")
	rootPath := filepath.Join(filepath.Dir(request.EpisodePath), "roots.json")
	if err := writeJSON(rootPath, request.Roots); err != nil {
		return nil, err
	}
	mcpConfig := map[string]any{"mcpServers": map[string]any{
		"phoenix": map[string]any{
			"type": "stdio", "command": request.PhoenixPath,
			"args": []string{
				"serve", "--stdio", "--world", request.WorldPath, "--world-schema", request.SchemaPath,
				"--episode-db", request.EpisodePath, "--root-refs", rootPath, "--world-build", request.WorldBuild,
				"--arm", request.Arm,
			},
			"env": map[string]string{},
		},
	}}
	if request.StateEventsPath != "" {
		server := mcpConfig["mcpServers"].(map[string]any)["phoenix"].(map[string]any)
		server["args"] = append(server["args"].([]string), "--state-events", request.StateEventsPath)
	}
	if err := writeJSON(configPath, mcpConfig); err != nil {
		return nil, err
	}
	args := []string{
		"-p", runtimePromptForArm(request.Goal, request.Roots, request.Arm),
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
		"--allowedTools", strings.Join(allowedTools, ","),
		"--system-prompt", systemPrompt,
	}
	if request.Arm == "B" {
		args = append(args, "--append-system-prompt-file", request.ArmBDocument)
	}
	return args, nil
}

func finishRuntimeInvocation(
	result runtimeResult,
	observation runtimeObservation,
	parseErr, runErr, contextErr error,
	stderr string,
	timeout time.Duration,
) (runtimeResult, error) {
	if contextErr != nil {
		result.Failure = terminalRuntimeFailure("timeout", fmt.Sprintf("runtime timed out after %s", timeout), observation.modelOutputSeen, true)
		return result, nil
	}
	if runErr != nil {
		message := strings.TrimSpace(stderr)
		if message == "" {
			message = runErr.Error()
		}
		result.Failure = classifyRuntimeFailure(runErr, message, observation)
		return result, nil
	}
	if failure := resultFailure(observation); failure != nil {
		message := observation.finalMessage
		if message == "" {
			message = failure.Message
		}
		result.Failure = classifyRuntimeFailure(errors.New("runtime result error"), message, observation)
		return result, nil
	}
	if !observation.modelOutputSeen && observation.mcpFailed {
		result.Failure = retryableRuntimeFailure("mcp_connection", "Phoenix MCP server did not connect before model output")
		return result, nil
	}
	if parseErr != nil {
		result.Failure = terminalRuntimeFailure("malformed_output", parseErr.Error(), observation.modelOutputSeen, false)
	}
	return result, nil
}

func runtimePromptForArm(goal string, roots map[string]string, arm string) string {
	if arm == "A" || arm == "B" {
		return goal
	}
	names := []string{"repo", "tests", "git", "episodes"}
	var prompt strings.Builder
	prompt.WriteString(goal)
	prompt.WriteString("\n\nLive handles:\n")
	for _, name := range names {
		fmt.Fprintf(&prompt, "- %s: %s\n", name, roots[name])
	}
	return prompt.String()
}

func allowedToolsForArm(arm string, flatToolNames []string) []string {
	if arm == "C" || arm == "D" || arm == "E" {
		return []string{"mcp__phoenix__act"}
	}
	tools := make([]string, len(flatToolNames))
	for index, name := range flatToolNames {
		tools[index] = "mcp__phoenix__" + name
	}
	return tools
}

func parseRuntimeOutput(output []byte) (string, error) {
	observation, err := inspectRuntimeOutput(output)
	if err != nil {
		return "", err
	}
	return observation.finalMessage, nil
}

type runtimeObservation struct {
	finalMessage    string
	resultSubtype   string
	resultError     bool
	modelOutputSeen bool
	mcpFailed       bool
	metrics         runtimeMetrics
}

func inspectRuntimeOutput(output []byte) (runtimeObservation, error) {
	var observation runtimeObservation
	scanner := bufio.NewScanner(bytes.NewReader(output))
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		var event map[string]any
		decoder := json.NewDecoder(bytes.NewReader(scanner.Bytes()))
		decoder.UseNumber()
		if decoder.Decode(&event) != nil {
			continue
		}
		typeName, _ := event["type"].(string)
		if typeName == "assistant" {
			observation.modelOutputSeen = true
		}
		if typeName == "system" {
			observation.mcpFailed = observation.mcpFailed || hasFailedMCPServer(event)
		}
		if typeName == "result" {
			observation.finalMessage, _ = event["result"].(string)
			observation.resultSubtype, _ = event["subtype"].(string)
			observation.resultError, _ = event["is_error"].(bool)
			observation.metrics = resultMetrics(event)
		}
	}
	if err := scanner.Err(); err != nil {
		return observation, fmt.Errorf("read runtime stream: %w", err)
	}
	if observation.finalMessage == "" {
		return observation, fmt.Errorf("runtime stream contains no final result")
	}
	return observation, nil
}

func resultMetrics(event map[string]any) runtimeMetrics {
	metrics := runtimeMetrics{
		Turns: intValue(event["num_turns"]), CostUSD: floatValue(event["total_cost_usd"]),
		APITimeMS: int64(intValue(event["duration_api_ms"])),
	}
	if usage, ok := event["usage"].(map[string]any); ok {
		metrics.InputTokens = int64(intValue(usage["input_tokens"]))
		metrics.CacheCreationInputTokens = int64(intValue(usage["cache_creation_input_tokens"]))
		metrics.CacheReadInputTokens = int64(intValue(usage["cache_read_input_tokens"]))
		metrics.OutputTokens = int64(intValue(usage["output_tokens"]))
	}
	if metrics.InputTokens == 0 && metrics.OutputTokens == 0 {
		sumModelUsage(event["modelUsage"], &metrics)
	}
	metrics.TotalTokens = metrics.InputTokens + metrics.CacheCreationInputTokens + metrics.CacheReadInputTokens + metrics.OutputTokens
	return metrics
}

func sumModelUsage(value any, metrics *runtimeMetrics) {
	models, ok := value.(map[string]any)
	if !ok {
		return
	}
	var modelCost float64
	for _, raw := range models {
		usage, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		metrics.InputTokens += int64(intValue(usage["inputTokens"]))
		metrics.CacheCreationInputTokens += int64(intValue(usage["cacheCreationInputTokens"]))
		metrics.CacheReadInputTokens += int64(intValue(usage["cacheReadInputTokens"]))
		metrics.OutputTokens += int64(intValue(usage["outputTokens"]))
		modelCost += floatValue(usage["costUSD"])
	}
	if metrics.CostUSD == 0 {
		metrics.CostUSD = modelCost
	}
}

func intValue(value any) int {
	switch typed := value.(type) {
	case json.Number:
		parsed, _ := strconv.ParseInt(string(typed), 10, 64)
		return int(parsed)
	case float64:
		return int(typed)
	default:
		return 0
	}
}

func floatValue(value any) float64 {
	switch typed := value.(type) {
	case json.Number:
		parsed, _ := strconv.ParseFloat(string(typed), 64)
		return parsed
	case float64:
		return typed
	default:
		return 0
	}
}

func hasFailedMCPServer(event map[string]any) bool {
	servers, ok := event["mcp_servers"].([]any)
	if !ok {
		return false
	}
	for _, raw := range servers {
		server, ok := raw.(map[string]any)
		if !ok || server["name"] != "phoenix" {
			continue
		}
		status, _ := server["status"].(string)
		status = strings.ToLower(status)
		return strings.Contains(status, "fail") || strings.Contains(status, "error") || status == "disconnected"
	}
	return false
}

func classifyRuntimeFailure(runErr error, message string, observation runtimeObservation) *runtimeFailure {
	beforeToken := !observation.modelOutputSeen
	if failure := resultFailure(observation); failure != nil {
		if failure.CapHit {
			failure.Message = message
			return failure
		}
	}
	if beforeToken && isLaunchFailure(runErr) {
		return retryableRuntimeFailure("runtime_launch", message)
	}
	if beforeToken && isTransientProviderFailure(message) {
		return retryableRuntimeFailure("provider_transient", message)
	}
	if beforeToken && (observation.mcpFailed || isMCPConnectionFailure(message)) {
		return retryableRuntimeFailure("mcp_connection", message)
	}
	if failure := resultFailure(observation); failure != nil {
		failure.Message = message
		return failure
	}
	return terminalRuntimeFailure("runtime_error", message, observation.modelOutputSeen, false)
}

func resultFailure(observation runtimeObservation) *runtimeFailure {
	subtype := strings.ToLower(observation.resultSubtype)
	switch {
	case strings.Contains(subtype, "max_turn"):
		return terminalRuntimeFailure("turn_limit", "runtime reached the turn limit", observation.modelOutputSeen, true)
	case strings.Contains(subtype, "budget") || strings.Contains(subtype, "cost"):
		return terminalRuntimeFailure("cost_cap", "runtime reached the cost cap", observation.modelOutputSeen, true)
	case observation.resultError:
		return terminalRuntimeFailure("agent_error", "runtime returned an error result", observation.modelOutputSeen, false)
	default:
		return nil
	}
}

func retryableRuntimeFailure(code, message string) *runtimeFailure {
	return &runtimeFailure{Code: code, Message: message, BeforeFirstModelToken: true, RetryEligible: true}
}

func terminalRuntimeFailure(code, message string, modelOutputSeen, capHit bool) *runtimeFailure {
	return &runtimeFailure{
		Code: code, Message: message, BeforeFirstModelToken: !modelOutputSeen,
		RetryEligible: false, CapHit: capHit,
	}
}

func isLaunchFailure(err error) bool {
	var execError *exec.Error
	var pathError *os.PathError
	return errors.As(err, &execError) || errors.As(err, &pathError)
}

func isTransientProviderFailure(message string) bool {
	lower := strings.ToLower(message)
	if providerStatusPattern.MatchString(lower) {
		return true
	}
	for _, marker := range []string{"rate limit", "internal server error", "service unavailable", "overloaded"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func isMCPConnectionFailure(message string) bool {
	lower := strings.ToLower(message)
	return strings.Contains(lower, "mcp connection") || strings.Contains(lower, "failed to connect to mcp")
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
