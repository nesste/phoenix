package devrepo

import (
	"context"
	"reflect"
	"testing"

	"github.com/nesste/phoenix/internal/verb"
	"github.com/nesste/phoenix/internal/world"
)

func TestCommandVerbsUseFixedExecutablesAndArgumentArrays(t *testing.T) {
	runner := &fakeRunner{responses: []verb.CommandResult{
		{ExitCode: 0, Stdout: []byte("build ok")},
		{ExitCode: 0, Stdout: []byte("TestX;touch_pwned\n")},
		{ExitCode: 0, Stdout: []byte("{\"Action\":\"pass\",\"Test\":\"TestX;touch_pwned\"}\n")},
		{ExitCode: 0},
		{ExitCode: 0, Stdout: []byte("committed")},
	}}
	executor := newTestExecutor(t, Config{Runner: runner, GoExecutable: "go_fixed", GitExecutable: "git_fixed"})
	resource := world.Resource{Kind: "path", Value: t.TempDir()}

	build := executor.Execute(context.Background(), verb.Request{HandleType: "repo", Verb: "build", Resource: resource, Args: map[string]any{}})
	if build.Status != verb.StatusOK {
		t.Fatalf("build result = %#v", build)
	}
	focus := executor.Execute(context.Background(), verb.Request{
		HandleType: "tests", Verb: "focus", Resource: resource, Args: map[string]any{"test": "TestX;touch_pwned"},
	})
	if focus.Status != verb.StatusOK || focus.Value.(map[string]any)["passed"] != true {
		t.Fatalf("focus result = %#v", focus)
	}
	commit := executor.Execute(context.Background(), verb.Request{
		HandleType: "git", Verb: "commit", Resource: resource, Args: map[string]any{"message": "safe; touch pwned"},
	})
	if commit.Status != verb.StatusOK {
		t.Fatalf("commit result = %#v", commit)
	}

	commands := runner.Commands()
	if commands[0].Executable != "go_fixed" || !reflect.DeepEqual(commands[0].Args, []string{"test", "-run=^$", "./..."}) {
		t.Fatalf("build command = %#v", commands[0])
	}
	if commands[2].Executable != "go_fixed" || commands[2].Args[len(commands[2].Args)-1] != "^TestX;touch_pwned$" {
		t.Fatalf("focus command did not preserve one typed argument: %#v", commands[2])
	}
	if commands[3].Executable != "git_fixed" || !reflect.DeepEqual(commands[3].Args, []string{"add", "-A", "--", "."}) {
		t.Fatalf("git add command = %#v", commands[3])
	}
	if commands[4].Args[3] != "safe; touch pwned" {
		t.Fatalf("commit message was not one argument: %#v", commands[4])
	}
}

func TestGitCommitReportsNothingToCommitAsAConstraint(t *testing.T) {
	runner := &fakeRunner{responses: []verb.CommandResult{
		{ExitCode: 0},
		{ExitCode: 1, Stdout: []byte("nothing to commit, working tree clean")},
	}}
	executor := newTestExecutor(t, Config{Runner: runner, GoExecutable: "go", GitExecutable: "git"})
	result := executor.Execute(context.Background(), verb.Request{
		HandleType: "git", Verb: "commit", Resource: world.Resource{Kind: "path", Value: t.TempDir()},
		Args: map[string]any{"message": "checkpoint"},
	})
	if result.Status != verb.StatusFail || result.Error == nil || result.Error.Code != "nothing_to_commit" {
		t.Fatalf("commit result = %#v, want shapeable nothing_to_commit constraint", result)
	}
}
