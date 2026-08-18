package verb

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

var (
	ErrExecutableNotAllowed = errors.New("executable is not allowlisted")
	ErrExecutableDigest     = errors.New("executable digest does not match allowlist")
	ErrOutputLimit          = errors.New("external output exceeded its size limit")
)

type Executable struct {
	ID     string
	Path   string
	Digest string
}

type Command struct {
	Executable string
	Args       []string
	Dir        string
	Stdin      []byte
}

type CommandResult struct {
	ExitCode int
	Stdout   []byte
	Stderr   []byte
}

type CommandRunner interface {
	Run(context.Context, Command) (CommandResult, error)
}

type ExternalRunner struct {
	executables map[string]Executable
	maxOutput   int
}

func NewExternalRunner(executables []Executable, maxOutput int) (*ExternalRunner, error) {
	if maxOutput <= 0 {
		return nil, fmt.Errorf("external output limit must be positive")
	}
	runner := &ExternalRunner{executables: make(map[string]Executable, len(executables)), maxOutput: maxOutput}
	for _, executable := range executables {
		if !identifierPattern.MatchString(executable.ID) || !filepath.IsAbs(executable.Path) {
			return nil, fmt.Errorf("%w: executable id and absolute path are required", ErrExecutableNotAllowed)
		}
		if _, exists := runner.executables[executable.ID]; exists {
			return nil, fmt.Errorf("%w: duplicate id %q", ErrExecutableNotAllowed, executable.ID)
		}
		if err := verifyExecutable(executable); err != nil {
			return nil, err
		}
		runner.executables[executable.ID] = executable
	}
	return runner, nil
}

func (runner *ExternalRunner) Run(ctx context.Context, command Command) (CommandResult, error) {
	executable, exists := runner.executables[command.Executable]
	if !exists {
		return CommandResult{}, ErrExecutableNotAllowed
	}
	if err := verifyExecutable(executable); err != nil {
		return CommandResult{}, err
	}

	process := exec.CommandContext(ctx, executable.Path, command.Args...)
	process.Dir = command.Dir
	process.Stdin = bytes.NewReader(command.Stdin)
	outputs := newBoundedOutputs(runner.maxOutput)
	process.Stdout = outputStream{outputs: outputs, stderr: false}
	process.Stderr = outputStream{outputs: outputs, stderr: true}
	err := process.Run()
	if outputs.exceededLimit() {
		return CommandResult{}, ErrOutputLimit
	}
	if ctx.Err() != nil {
		return CommandResult{}, ctx.Err()
	}
	result := CommandResult{ExitCode: 0, Stdout: outputs.stdoutBytes(), Stderr: outputs.stderrBytes()}
	if err == nil {
		return result, nil
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		result.ExitCode = exitError.ExitCode()
		return result, nil
	}
	return CommandResult{}, fmt.Errorf("start allowlisted executable: %w", err)
}

func DigestExecutable(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return "sha256:" + hex.EncodeToString(hash.Sum(nil)), nil
}

func verifyExecutable(executable Executable) error {
	actual, err := DigestExecutable(executable.Path)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrExecutableDigest, err)
	}
	if actual != executable.Digest {
		return ErrExecutableDigest
	}
	return nil
}

type boundedOutputs struct {
	mu        sync.Mutex
	remaining int
	exceeded  bool
	stdout    bytes.Buffer
	stderr    bytes.Buffer
}

func newBoundedOutputs(limit int) *boundedOutputs {
	return &boundedOutputs{remaining: limit}
}

func (outputs *boundedOutputs) write(stderr bool, contents []byte) (int, error) {
	outputs.mu.Lock()
	defer outputs.mu.Unlock()
	original := len(contents)
	if len(contents) > outputs.remaining {
		contents = contents[:outputs.remaining]
		outputs.exceeded = true
	}
	outputs.remaining -= len(contents)
	if stderr {
		_, _ = outputs.stderr.Write(contents)
	} else {
		_, _ = outputs.stdout.Write(contents)
	}
	return original, nil
}

func (outputs *boundedOutputs) exceededLimit() bool {
	outputs.mu.Lock()
	defer outputs.mu.Unlock()
	return outputs.exceeded
}

func (outputs *boundedOutputs) stdoutBytes() []byte {
	outputs.mu.Lock()
	defer outputs.mu.Unlock()
	return append([]byte(nil), outputs.stdout.Bytes()...)
}

func (outputs *boundedOutputs) stderrBytes() []byte {
	outputs.mu.Lock()
	defer outputs.mu.Unlock()
	return append([]byte(nil), outputs.stderr.Bytes()...)
}

type outputStream struct {
	outputs *boundedOutputs
	stderr  bool
}

func (stream outputStream) Write(contents []byte) (int, error) {
	return stream.outputs.write(stream.stderr, contents)
}
