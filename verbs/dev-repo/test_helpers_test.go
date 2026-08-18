package devrepo

import (
	"context"
	"testing"
	"time"

	"github.com/nesste/phoenix/internal/verb"
)

func newTestExecutor(t *testing.T, config Config) *verb.Executor {
	t.Helper()
	registry := verb.NewRegistry()
	if err := Register(registry, config); err != nil {
		t.Fatal(err)
	}
	return verb.NewExecutor(registry, verb.Options{Timeout: time.Second, MaxResultBytes: 64 * 1024})
}

type fakeRunner struct {
	commands  []verb.Command
	responses []verb.CommandResult
}

func (runner *fakeRunner) Run(_ context.Context, command verb.Command) (verb.CommandResult, error) {
	runner.commands = append(runner.commands, command)
	if len(runner.responses) == 0 {
		return verb.CommandResult{}, nil
	}
	response := runner.responses[0]
	runner.responses = runner.responses[1:]
	return response, nil
}

func (runner *fakeRunner) Commands() []verb.Command {
	return append([]verb.Command(nil), runner.commands...)
}
