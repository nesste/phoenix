package devrepo

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/nesste/phoenix/internal/verb"
)

func buildHandler(config Config) verb.HandlerFunc {
	return func(ctx context.Context, request verb.Request) (any, error) {
		return runSimple(ctx, config, request, config.GoExecutable, []string{"test", "-run=^$", "./..."})
	}
}

func testsRunHandler(config Config) verb.HandlerFunc {
	return func(ctx context.Context, request verb.Request) (any, error) {
		result, err := run(ctx, config, request, config.GoExecutable, []string{"test", "-json", "./..."})
		if err != nil {
			return nil, err
		}
		passed, failed := testCounts(result.Stdout)
		if result.ExitCode != 0 && failed == 0 {
			failed = 1
		}
		return commandMap(result, map[string]any{"passed": passed, "failed": failed}), nil
	}
}

func testsListHandler(config Config) verb.HandlerFunc {
	return func(ctx context.Context, request verb.Request) (any, error) {
		result, err := run(ctx, config, request, config.GoExecutable, []string{"test", "-list", ".", "./..."})
		if err != nil {
			return nil, err
		}
		if result.ExitCode != 0 {
			return nil, verb.NewFailure("test_list_failed", "test listing failed", map[string]any{"exit_code": result.ExitCode})
		}
		return map[string]any{"tests": stringsAsAny(listedTests(result.Stdout))}, nil
	}
}

func testsFocusHandler(config Config) verb.HandlerFunc {
	return func(ctx context.Context, request verb.Request) (any, error) {
		testName := request.Args["test"].(string)
		listed, err := run(ctx, config, request, config.GoExecutable, []string{"test", "-list", ".", "./..."})
		if err != nil {
			return nil, err
		}
		if listed.ExitCode != 0 {
			return nil, verb.NewFailure("test_list_failed", "test listing failed", map[string]any{"exit_code": listed.ExitCode})
		}
		if !containsTest(listedTests(listed.Stdout), testName) {
			return nil, verb.NewFailure("unknown_test", "requested test is not in the live suite", nil)
		}
		pattern := "^" + regexp.QuoteMeta(testName) + "$"
		result, err := run(ctx, config, request, config.GoExecutable, []string{"test", "-json", "./...", "-run", pattern})
		if err != nil {
			return nil, err
		}
		passed := testPassed(result.Stdout, testName) && result.ExitCode == 0
		return commandMap(result, map[string]any{"test": testName, "passed": passed}), nil
	}
}

func containsTest(tests []string, name string) bool {
	for _, test := range tests {
		if test == name {
			return true
		}
	}
	return false
}

func gitStatusHandler(config Config) verb.HandlerFunc {
	return func(ctx context.Context, request verb.Request) (any, error) {
		branch, clean, output, err := status(ctx, config, request)
		if err != nil {
			return nil, err
		}
		return map[string]any{"branch": branch, "clean": clean, "output": output}, nil
	}
}

func repoStatusHandler(config Config) verb.HandlerFunc {
	return func(ctx context.Context, request verb.Request) (any, error) {
		branch, clean, _, err := status(ctx, config, request)
		if err != nil {
			return nil, err
		}
		return map[string]any{"branch": branch, "clean": clean}, nil
	}
}

func status(ctx context.Context, config Config, request verb.Request) (string, bool, string, error) {
	result, err := run(ctx, config, request, config.GitExecutable, []string{"status", "--porcelain=v1", "--branch"})
	if err != nil {
		return "", false, "", err
	}
	if result.ExitCode != 0 {
		return "", false, "", verb.NewFailure("git_failed", "git status failed", map[string]any{"exit_code": result.ExitCode})
	}
	output := string(result.Stdout)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	branch := ""
	if len(lines) > 0 {
		branch = strings.TrimSpace(strings.TrimPrefix(lines[0], "##"))
	}
	return branch, len(lines) <= 1, output, nil
}

func gitDiffHandler(config Config) verb.HandlerFunc {
	return func(ctx context.Context, request verb.Request) (any, error) {
		result, err := run(ctx, config, request, config.GitExecutable, []string{"diff", "--no-ext-diff", "HEAD", "--"})
		if err != nil {
			return nil, err
		}
		if result.ExitCode != 0 {
			return nil, verb.NewFailure("git_failed", "git diff failed", map[string]any{"exit_code": result.ExitCode})
		}
		return map[string]any{"diff": string(result.Stdout)}, nil
	}
}

func gitCommitHandler(config Config) verb.HandlerFunc {
	return func(ctx context.Context, request verb.Request) (any, error) {
		add, err := run(ctx, config, request, config.GitExecutable, []string{"add", "-A", "--", "."})
		if err != nil {
			return nil, err
		}
		if add.ExitCode != 0 {
			return nil, verb.NewFailure("git_failed", "git add failed", map[string]any{"exit_code": add.ExitCode})
		}
		message := request.Args["message"].(string)
		committed, err := run(ctx, config, request, config.GitExecutable, []string{"commit", "--no-verify", "-m", message, "--"})
		if err != nil {
			return nil, err
		}
		if committed.ExitCode != 0 {
			output := strings.ToLower(string(committed.Stdout) + "\n" + string(committed.Stderr))
			if strings.Contains(output, "nothing to commit") || strings.Contains(output, "no changes added to commit") {
				return nil, verb.NewFailure("nothing_to_commit", "no changes are available to commit", nil)
			}
			return nil, verb.NewFailure("git_failed", "git commit failed", map[string]any{"exit_code": committed.ExitCode})
		}
		return map[string]any{
			"committed": true, "stdout": string(committed.Stdout), "stderr": string(committed.Stderr),
		}, nil
	}
}

func runSimple(ctx context.Context, config Config, request verb.Request, executable string, args []string) (any, error) {
	result, err := run(ctx, config, request, executable, args)
	if err != nil {
		return nil, err
	}
	return commandMap(result, nil), nil
}

func run(ctx context.Context, config Config, request verb.Request, executable string, args []string) (verb.CommandResult, error) {
	root, err := repoRoot(request.Resource)
	if err != nil {
		return verb.CommandResult{}, err
	}
	result, err := config.Runner.Run(ctx, verb.Command{Executable: executable, Args: args, Dir: root})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return verb.CommandResult{}, err
		}
		if errors.Is(err, verb.ErrOutputLimit) {
			return verb.CommandResult{}, verb.NewFailure("output_limit", "external command output exceeded its size limit", nil)
		}
		if errors.Is(err, verb.ErrExecutableDigest) {
			return verb.CommandResult{}, verb.NewFailure("executable_changed", "allowlisted executable identity changed", nil)
		}
		return verb.CommandResult{}, verb.NewFailure("command_failed", "allowlisted command could not be executed", nil)
	}
	return result, nil
}

func commandMap(result verb.CommandResult, extra map[string]any) map[string]any {
	value := map[string]any{
		"exit_code": result.ExitCode, "stdout": string(result.Stdout), "stderr": string(result.Stderr),
	}
	for key, item := range extra {
		value[key] = item
	}
	return value
}

type testEvent struct {
	Action string `json:"Action"`
	Test   string `json:"Test"`
}

func testCounts(output []byte) (int, int) {
	final := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		var event testEvent
		if json.Unmarshal(scanner.Bytes(), &event) == nil && event.Test != "" && (event.Action == "pass" || event.Action == "fail") {
			final[event.Test] = event.Action
		}
	}
	passed, failed := 0, 0
	for _, action := range final {
		if action == "pass" {
			passed++
		} else {
			failed++
		}
	}
	return passed, failed
}

func testPassed(output []byte, name string) bool {
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		var event testEvent
		if json.Unmarshal(scanner.Bytes(), &event) == nil && event.Test == name && event.Action == "pass" {
			return true
		}
	}
	return false
}

func listedTests(output []byte) []string {
	var tests []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Test") || strings.HasPrefix(line, "Benchmark") || strings.HasPrefix(line, "Example") {
			tests = append(tests, line)
		}
	}
	return tests
}
