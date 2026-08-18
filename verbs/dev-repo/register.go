// Package devrepo provides the closed Phase 1 starter verb set for a local
// development repository.
package devrepo

import (
	"context"
	"fmt"

	"github.com/nesste/phoenix/internal/verb"
)

type EpisodePointer struct {
	EpisodeID string `json:"episode_id"`
	ActID     string `json:"act_id"`
}

type EpisodeRecaller interface {
	Recall(context.Context, string, string) ([]EpisodePointer, error)
}

type Config struct {
	Runner           verb.CommandRunner
	GoExecutable     string
	GitExecutable    string
	Episodes         EpisodeRecaller
	MaxReadBytes     int64
	MaxSearchMatches int
}

func Register(registry *verb.Registry, config Config) error {
	definitions, err := Definitions(config)
	if err != nil {
		return err
	}
	for _, definition := range definitions {
		if err := registry.Register(definition); err != nil {
			return err
		}
	}
	return nil
}

func Definitions(config Config) ([]verb.Definition, error) {
	if config.Runner == nil {
		return nil, fmt.Errorf("command runner is required")
	}
	if config.GoExecutable == "" || config.GitExecutable == "" {
		return nil, fmt.Errorf("fixed Go and Git executable ids are required")
	}
	if config.MaxReadBytes <= 0 {
		config.MaxReadBytes = 32 * 1024
	}
	if config.MaxSearchMatches <= 0 {
		config.MaxSearchMatches = 20
	}

	return []verb.Definition{
		definition("repo", "read", pathArgs, readResult, readHandler(config.MaxReadBytes)),
		definition("repo", "edit", editArgs, editResult, editHandler(config.MaxReadBytes)),
		definition("repo", "build", emptyArgs, commandResultSchema, buildHandler(config)),
		definition("repo", "find", queryArgs, findResult, findHandler(config)),
		definition("repo", "status", emptyArgs, repoStatusResult, repoStatusHandler(config)),
		definition("tests", "run", emptyArgs, testsRunResult, testsRunHandler(config)),
		definition("tests", "list", emptyArgs, testsListResult, testsListHandler(config)),
		definition("tests", "focus", focusArgs, testsFocusResult, testsFocusHandler(config)),
		definition("git", "status", emptyArgs, gitStatusResult, gitStatusHandler(config)),
		definition("git", "diff", emptyArgs, gitDiffResult, gitDiffHandler(config)),
		definition("git", "commit", commitArgs, gitCommitResult, gitCommitHandler(config)),
		definition("episodes", "recall", queryArgs, recallResult, recallHandler(config.Episodes)),
	}, nil
}

func definition(handleType, name string, args, result []byte, handler verb.Handler) verb.Definition {
	return verb.Definition{HandleType: handleType, Name: name, ArgsSchema: args, ResultSchema: result, Handler: handler}
}
