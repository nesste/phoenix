package verb

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nesste/phoenix/internal/redact"
	"github.com/nesste/phoenix/internal/world"
)

type Handler interface {
	Execute(context.Context, Request) (any, error)
}

type HandlerFunc func(context.Context, Request) (any, error)

func (function HandlerFunc) Execute(ctx context.Context, request Request) (any, error) {
	return function(ctx, request)
}

type Request struct {
	SessionID  string
	HandleType string
	Verb       string
	Resource   world.Resource
	State      *world.LiveState
	Args       map[string]any
	Reachable  ReachableFinder
}

type ReachableFinder interface {
	FindReachable(query string) []world.ReachableMatch
}

type Status string

const (
	StatusOK   Status = "ok"
	StatusFail Status = "fail"
)

type Failure struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
	// DeclaredArgs carries the attempted verb's declared argument schema on
	// invalid_arguments failures so a failed guess becomes a corrected next
	// call. It is capped at declaredArgsLimit serialized bytes.
	DeclaredArgs json.RawMessage `json:"declared_args,omitempty"`
}

// declaredArgsLimit is the frozen serialized-size cap for declared_args; a
// larger schema is truncated to its top-level property names and types.
const declaredArgsLimit = 2048

func capDeclaredArgs(raw json.RawMessage) json.RawMessage {
	compact := new(bytes.Buffer)
	if err := json.Compact(compact, raw); err != nil {
		return nil
	}
	if compact.Len() <= declaredArgsLimit {
		return json.RawMessage(compact.Bytes())
	}
	var schema map[string]any
	if err := json.Unmarshal(raw, &schema); err != nil {
		return nil
	}
	truncated := map[string]any{"truncated": true}
	if declaredType, ok := schema["type"].(string); ok {
		truncated["type"] = declaredType
	}
	if properties, ok := schema["properties"].(map[string]any); ok {
		names := make(map[string]any, len(properties))
		for name, value := range properties {
			propertyType := ""
			if member, ok := value.(map[string]any); ok {
				propertyType, _ = member["type"].(string)
			}
			names[name] = map[string]any{"type": propertyType}
		}
		truncated["properties"] = names
	}
	encoded, err := json.Marshal(truncated)
	if err != nil {
		return nil
	}
	return json.RawMessage(encoded)
}

type Result struct {
	Status Status
	Value  any
	Error  *Failure
}

type FailureError struct {
	Failure Failure
}

func (failure *FailureError) Error() string {
	return failure.Failure.Message
}

func NewFailure(code, message string, details map[string]any) error {
	return &FailureError{Failure: Failure{Code: code, Message: message, Details: details}}
}

type Options struct {
	Timeout        time.Duration
	MaxResultBytes int
	Secrets        []string
}

type Executor struct {
	registry       *Registry
	timeout        time.Duration
	maxResultBytes int
	redactor       redact.Redactor
}

func NewExecutor(registry *Registry, options Options) *Executor {
	if options.Timeout <= 0 {
		options.Timeout = 30 * time.Second
	}
	if options.MaxResultBytes <= 0 {
		options.MaxResultBytes = 64 * 1024
	}
	return &Executor{
		registry: registry, timeout: options.Timeout, maxResultBytes: options.MaxResultBytes,
		redactor: redact.New(options.Secrets),
	}
}

func (executor *Executor) Execute(ctx context.Context, request Request) Result {
	entry, exists := executor.registry.lookup(request.HandleType, request.Verb)
	if !exists {
		return failed("unregistered_verb", "reachable verb has no registered implementation", nil)
	}
	if request.Args == nil {
		request.Args = map[string]any{}
	}
	if err := entry.args.Validate(request.Args); err != nil {
		result := failed("invalid_arguments", "verb arguments do not match the declared schema", nil)
		result.Error.DeclaredArgs = capDeclaredArgs(entry.definition.ArgsSchema)
		return result
	}

	callContext, cancel := context.WithTimeout(ctx, executor.timeout)
	defer cancel()
	outcomes := make(chan handlerOutcome, 1)
	go callHandler(callContext, entry.definition.Handler, request, outcomes)
	select {
	case <-callContext.Done():
		if ctx.Err() != nil {
			return failed("cancelled", "verb execution was cancelled", nil)
		}
		return failed("timeout", "verb execution exceeded its time limit", nil)
	case outcome := <-outcomes:
		if outcome.err != nil {
			return executor.handlerFailure(outcome.err)
		}
		return executor.success(entry, outcome.value)
	}
}

type handlerOutcome struct {
	value any
	err   error
}

func callHandler(ctx context.Context, handler Handler, request Request, outcomes chan<- handlerOutcome) {
	defer func() {
		if recover() != nil {
			outcomes <- handlerOutcome{err: fmt.Errorf("handler panic")}
		}
	}()
	value, err := handler.Execute(ctx, request)
	outcomes <- handlerOutcome{value: value, err: err}
}

func (executor *Executor) success(entry registered, value any) Result {
	encoded, err := json.Marshal(value)
	if err != nil {
		return failed("invalid_result", "verb result is not valid JSON", nil)
	}
	if len(encoded) > executor.maxResultBytes {
		return failed("output_limit", "verb result exceeded its size limit", map[string]any{"limit_bytes": executor.maxResultBytes})
	}
	var normalized any
	if err := json.Unmarshal(encoded, &normalized); err != nil {
		return failed("invalid_result", "verb result is not valid JSON", nil)
	}
	if err := entry.result.Validate(normalized); err != nil {
		return failed("invalid_result", "verb result does not match the declared schema", nil)
	}
	redacted := executor.redactor.Value(normalized)
	if err := entry.result.Validate(redacted); err != nil {
		return failed("invalid_result", "redacted verb result does not match the declared schema", nil)
	}
	redactedJSON, err := json.Marshal(redacted)
	if err != nil {
		return failed("invalid_result", "redacted verb result is not valid JSON", nil)
	}
	if len(redactedJSON) > executor.maxResultBytes {
		return failed("output_limit", "redacted verb result exceeded its size limit", map[string]any{"limit_bytes": executor.maxResultBytes})
	}
	return Result{Status: StatusOK, Value: redacted}
}

func (executor *Executor) handlerFailure(err error) Result {
	var failure *FailureError
	if errors.As(err, &failure) && identifierPattern.MatchString(failure.Failure.Code) && failure.Failure.Message != "" {
		normalized, normalizeErr := normalizeJSON(map[string]any{
			"code": failure.Failure.Code, "message": failure.Failure.Message, "details": failure.Failure.Details,
		})
		if normalizeErr != nil {
			return failed("execution_failed", "verb execution failed", nil)
		}
		redacted := executor.redactor.Value(normalized).(map[string]any)
		return failed(
			redacted["code"].(string), redacted["message"].(string), optionalMap(redacted["details"]),
		)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return failed("timeout", "verb execution exceeded its time limit", nil)
	}
	if errors.Is(err, context.Canceled) {
		return failed("cancelled", "verb execution was cancelled", nil)
	}
	return failed("execution_failed", "verb execution failed", nil)
}

func normalizeJSON(value any) (any, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var normalized any
	if err := json.Unmarshal(encoded, &normalized); err != nil {
		return nil, err
	}
	return normalized, nil
}

func failed(code, message string, details map[string]any) Result {
	return Result{Status: StatusFail, Error: &Failure{Code: code, Message: message, Details: details}}
}

func optionalMap(value any) map[string]any {
	if value == nil {
		return nil
	}
	result, _ := value.(map[string]any)
	return result
}
