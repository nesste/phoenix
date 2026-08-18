package verb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var objectArgs = []byte(`{"type":"object","additionalProperties":false}`)

var objectResult = []byte(`{
  "type":"object",
  "additionalProperties":false,
  "required":["message"],
  "properties":{"message":{"type":"string"},"token":{"type":"string"}}
}`)

func TestRegistryRejectsIncompleteOrUntypedDefinitions(t *testing.T) {
	valid := Definition{
		HandleType: "repo", Name: "inspect", ArgsSchema: objectArgs, ResultSchema: objectResult,
		Handler: HandlerFunc(func(context.Context, Request) (any, error) {
			return map[string]any{"message": "ok"}, nil
		}),
	}
	tests := []struct {
		name string
		edit func(*Definition)
		want string
	}{
		{name: "handle type", edit: func(definition *Definition) { definition.HandleType = "" }, want: "handle type"},
		{name: "argument schema", edit: func(definition *Definition) { definition.ArgsSchema = nil }, want: "argument schema"},
		{name: "untyped result", edit: func(definition *Definition) { definition.ResultSchema = []byte(`{}`) }, want: "result schema must declare a type"},
		{name: "handler", edit: func(definition *Definition) { definition.Handler = nil }, want: "handler"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registry := NewRegistry()
			definition := valid
			test.edit(&definition)
			if err := registry.Register(definition); err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Register error = %v, want it to contain %q", err, test.want)
			}
		})
	}

	registry := NewRegistry()
	if err := registry.Register(valid); err != nil {
		t.Fatalf("Register(valid): %v", err)
	}
	if err := registry.Register(valid); err == nil || !strings.Contains(err.Error(), "already registered") {
		t.Fatalf("duplicate Register error = %v, want duplicate rejection", err)
	}
}

func TestExecutorValidatesTypedResultsAndRedactsBeforeReturn(t *testing.T) {
	registry := NewRegistry()
	called := false
	err := registry.Register(Definition{
		HandleType: "repo", Name: "inspect", ArgsSchema: objectArgs, ResultSchema: objectResult,
		Handler: HandlerFunc(func(context.Context, Request) (any, error) {
			called = true
			return map[string]any{"message": "credential sk-live", "token": "secondary"}, nil
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	executor := NewExecutor(registry, Options{Timeout: time.Second, MaxResultBytes: 1024, Secrets: []string{"sk-live"}})

	invalid := executor.Execute(context.Background(), Request{HandleType: "repo", Verb: "inspect", Args: map[string]any{"unknown": true}})
	if invalid.Status != StatusFail || invalid.Error.Code != "invalid_arguments" {
		t.Fatalf("invalid arguments result = %#v", invalid)
	}
	if called {
		t.Fatal("handler ran with invalid arguments")
	}

	result := executor.Execute(context.Background(), Request{HandleType: "repo", Verb: "inspect", Args: map[string]any{}})
	if result.Status != StatusOK || result.Error != nil {
		t.Fatalf("Execute status = %q, error = %#v", result.Status, result.Error)
	}
	value := result.Value.(map[string]any)
	if got, want := value["message"], "credential [REDACTED]"; got != want {
		t.Fatalf("redacted message = %q, want %q", got, want)
	}
	if got, want := value["token"], "[REDACTED]"; got != want {
		t.Fatalf("sensitive field = %q, want %q", got, want)
	}
}

func TestExecutorReturnsStructuredFailuresWithoutRawErrors(t *testing.T) {
	registry := NewRegistry()
	registerTestVerb(t, registry, "fails", objectResult, func(context.Context, Request) (any, error) {
		return nil, errors.New("runtime stack: secret filesystem detail")
	})
	executor := NewExecutor(registry, Options{Timeout: time.Second, MaxResultBytes: 1024})

	result := executor.Execute(context.Background(), Request{HandleType: "repo", Verb: "fails", Args: map[string]any{}})
	if result.Status != StatusFail || result.Error == nil || result.Error.Code != "execution_failed" {
		t.Fatalf("failure result = %#v", result)
	}
	encoded := fmt.Sprintf("%#v", result)
	if strings.Contains(encoded, "runtime stack") || strings.Contains(encoded, "filesystem detail") {
		t.Fatalf("failure leaked raw error: %s", encoded)
	}
}

func TestExecutorRejectsHandlerOutputOutsideDeclaredResultType(t *testing.T) {
	registry := NewRegistry()
	registerTestVerb(t, registry, "bad_result", objectResult, func(context.Context, Request) (any, error) {
		return map[string]any{"unexpected": true}, nil
	})
	result := NewExecutor(registry, Options{}).Execute(context.Background(), Request{
		HandleType: "repo", Verb: "bad_result", Args: map[string]any{},
	})
	if result.Status != StatusFail || result.Error.Code != "invalid_result" {
		t.Fatalf("invalid handler result = %#v", result)
	}
}

func TestExecutorEnforcesTimeoutAndResultSize(t *testing.T) {
	registry := NewRegistry()
	registerTestVerb(t, registry, "slow", objectResult, func(ctx context.Context, _ Request) (any, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	})
	registerTestVerb(t, registry, "large", objectResult, func(context.Context, Request) (any, error) {
		return map[string]any{"message": strings.Repeat("x", 256)}, nil
	})
	executor := NewExecutor(registry, Options{Timeout: 20 * time.Millisecond, MaxResultBytes: 64})

	timedOut := executor.Execute(context.Background(), Request{HandleType: "repo", Verb: "slow", Args: map[string]any{}})
	if timedOut.Status != StatusFail || timedOut.Error.Code != "timeout" {
		t.Fatalf("timeout result = %#v", timedOut)
	}
	large := executor.Execute(context.Background(), Request{HandleType: "repo", Verb: "large", Args: map[string]any{}})
	if large.Status != StatusFail || large.Error.Code != "output_limit" {
		t.Fatalf("large result = %#v", large)
	}
}

func TestExternalRunnerUsesDigestAllowlistAndBoundsOutput(t *testing.T) {
	_, runner := externalTestRunner(t, 32)
	result, err := runner.Run(context.Background(), Command{
		Executable: "helper",
		Args:       []string{"-test.run=TestExternalHelperProcess", "--", "echo"},
		Stdin:      []byte("typed input"),
	})
	if err != nil {
		t.Fatalf("Run(helper): %v", err)
	}
	if got, want := string(result.Stdout), "typed input"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}

	if _, err := runner.Run(context.Background(), Command{
		Executable: "helper", Args: []string{"-test.run=TestExternalHelperProcess", "--", "large"},
	}); !errors.Is(err, ErrOutputLimit) {
		t.Fatalf("Run(large) error = %v, want ErrOutputLimit", err)
	}
	timeoutContext, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if _, err := runner.Run(timeoutContext, Command{
		Executable: "helper", Args: []string{"-test.run=TestExternalHelperProcess", "--", "slow"},
	}); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Run(slow) error = %v, want deadline exceeded", err)
	}
}

func TestExternalRunnerRejectsUnknownAndSubstitutedExecutables(t *testing.T) {
	executable, runner := externalTestRunner(t, 32)
	if _, err := runner.Run(context.Background(), Command{Executable: executable}); !errors.Is(err, ErrExecutableNotAllowed) {
		t.Fatalf("Run(arbitrary path) error = %v, want ErrExecutableNotAllowed", err)
	}
	if _, err := NewExternalRunner([]Executable{{ID: "wrong", Path: executable, Digest: "sha256:" + strings.Repeat("0", 64)}}, 32); !errors.Is(err, ErrExecutableDigest) {
		t.Fatalf("NewExternalRunner(wrong digest) error = %v, want ErrExecutableDigest", err)
	}

	copyPath := filepath.Join(t.TempDir(), "helper-copy")
	copyExecutable(t, executable, copyPath)
	copyDigest, err := DigestExecutable(copyPath)
	if err != nil {
		t.Fatal(err)
	}
	substitutionRunner, err := NewExternalRunner([]Executable{{ID: "copy", Path: copyPath, Digest: copyDigest}}, 32)
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.OpenFile(copyPath, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.Write([]byte("substitution")); err != nil {
		file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := substitutionRunner.Run(context.Background(), Command{Executable: "copy"}); !errors.Is(err, ErrExecutableDigest) {
		t.Fatalf("Run(substituted executable) error = %v, want ErrExecutableDigest", err)
	}
}

func TestExternalHelperProcess(t *testing.T) {
	separator := -1
	for index, argument := range os.Args {
		if argument == "--" {
			separator = index
			break
		}
	}
	if separator < 0 || separator+1 >= len(os.Args) {
		return
	}
	switch os.Args[separator+1] {
	case "echo":
		_, _ = io.Copy(os.Stdout, os.Stdin)
	case "large":
		_, _ = io.WriteString(os.Stdout, strings.Repeat("x", 128))
	case "slow":
		time.Sleep(10 * time.Second)
	default:
		os.Exit(2)
	}
	os.Exit(0)
}

func registerTestVerb(t *testing.T, registry *Registry, name string, resultSchema []byte, handler HandlerFunc) {
	t.Helper()
	if err := registry.Register(Definition{
		HandleType: "repo", Name: name, ArgsSchema: objectArgs, ResultSchema: resultSchema, Handler: handler,
	}); err != nil {
		t.Fatal(err)
	}
}

func copyExecutable(t *testing.T, source, target string) {
	t.Helper()
	input, err := os.Open(source)
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	output, err := os.OpenFile(target, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o700)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(output, input); err != nil {
		output.Close()
		t.Fatal(err)
	}
	if err := output.Close(); err != nil {
		t.Fatal(err)
	}
}

func externalTestRunner(t *testing.T, limit int) (string, *ExternalRunner) {
	t.Helper()
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	digest, err := DigestExecutable(executable)
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewExternalRunner([]Executable{{ID: "helper", Path: executable, Digest: digest}}, limit)
	if err != nil {
		t.Fatal(err)
	}
	return executable, runner
}
