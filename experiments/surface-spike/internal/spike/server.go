package spike

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	RepoHandle  = "h_repo_demo"
	TestsHandle = "h_tests_demo"
)

var expressionPattern = regexp.MustCompile(`^([a-z_]+)\.([a-z_]+)\((.*)\)$`)

type ActInput struct {
	Handle string         `json:"handle" jsonschema:"opaque live handle reference"`
	Verb   string         `json:"verb" jsonschema:"verb attached to the handle"`
	Args   map[string]any `json:"args,omitempty" jsonschema:"fully bound verb arguments"`
}

type EvalInput struct {
	Expression string `json:"expression" jsonschema:"single world expression such as repo.status()"`
}

type Call struct {
	Handle string         `json:"handle"`
	Verb   string         `json:"verb"`
	Args   map[string]any `json:"args"`
}

type FrontierEntry struct {
	Call       Call    `json:"call"`
	Why        string  `json:"why"`
	Provenance string  `json:"provenance"`
	Score      float64 `json:"score"`
}

type Refusal struct {
	What    string `json:"what"`
	Why     string `json:"why"`
	Instead *Call  `json:"instead"`
}

type Envelope struct {
	Version  int             `json:"v"`
	Handle   string          `json:"handle"`
	Verb     string          `json:"verb"`
	Status   string          `json:"status"`
	Result   map[string]any  `json:"result,omitempty"`
	Text     string          `json:"text"`
	Frontier []FrontierEntry `json:"frontier"`
	Refusal  *Refusal        `json:"refusal,omitempty"`
}

func NewServer(surface string, flatVerbCount int) (*mcp.Server, error) {
	server := mcp.NewServer(&mcp.Implementation{Name: "phoenix-surface-spike", Version: "0.1.0"}, nil)
	switch surface {
	case "act":
		mcp.AddTool(server, &mcp.Tool{
			Name:        "act",
			Description: "Invoke one verb on a live Phoenix handle. Results include up to three ready-to-run next calls.",
		}, actTool)
	case "eval":
		mcp.AddTool(server, &mcp.Tool{
			Name:        "eval",
			Description: "Evaluate one expression against the live Phoenix object graph. Results include ready-to-run next expressions.",
		}, evalTool)
	case "flat":
		if flatVerbCount < 1 {
			return nil, fmt.Errorf("flat verb count must be positive")
		}
		addFlatTools(server, flatVerbCount)
	default:
		return nil, fmt.Errorf("unknown surface %q", surface)
	}
	return server, nil
}

func actTool(_ context.Context, _ *mcp.CallToolRequest, input ActInput) (*mcp.CallToolResult, Envelope, error) {
	envelope := Execute(input)
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: envelope.Text}},
		IsError: envelope.Status != "ok",
	}, envelope, nil
}

func evalTool(_ context.Context, _ *mcp.CallToolRequest, input EvalInput) (*mcp.CallToolResult, Envelope, error) {
	call, ok := parseExpression(input.Expression)
	if !ok {
		envelope := refused("world", "eval", "malformed expression", nil)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: envelope.Text}},
			IsError: true,
		}, envelope, nil
	}
	envelope := Execute(call)
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: envelope.Text}},
		IsError: envelope.Status != "ok",
	}, envelope, nil
}

func Execute(input ActInput) Envelope {
	if input.Args == nil {
		input.Args = map[string]any{}
	}
	switch input.Handle + "." + input.Verb {
	case RepoHandle + ".status":
		return success(input, map[string]any{"branch": "main", "clean": true}, "repo is clean on main", next(TestsHandle, "run", "verify the suite before changing code"))
	case RepoHandle + ".read":
		path, _ := input.Args["path"].(string)
		if path == "" {
			return refused(input.Handle, input.Verb, "path is required", &Call{Handle: RepoHandle, Verb: "status", Args: map[string]any{}})
		}
		return success(input, map[string]any{"path": path, "content": "package demo"}, "read "+path, next(TestsHandle, "run", "check behavior after inspecting the file"))
	case RepoHandle + ".build":
		return success(input, map[string]any{"exit_code": 0}, "build passed", next(TestsHandle, "run", "run tests against the successful build"))
	case TestsHandle + ".run":
		return success(input, map[string]any{"passed": 12, "failed": 0}, "12 tests passed", next(RepoHandle, "status", "inspect the repository state after verification"))
	case TestsHandle + ".focus":
		test, _ := input.Args["test"].(string)
		if test == "" {
			return refused(input.Handle, input.Verb, "test is required", &Call{Handle: TestsHandle, Verb: "list", Args: map[string]any{}})
		}
		return success(input, map[string]any{"test": test, "passed": true}, test+" passed", next(RepoHandle, "status", "inspect the repository after the focused test"))
	case TestsHandle + ".list":
		return success(input, map[string]any{"tests": []string{"auth_test", "repo_test"}}, "2 tests available", FrontierEntry{
			Call:       Call{Handle: TestsHandle, Verb: "focus", Args: map[string]any{"test": "auth_test"}},
			Why:        "run the relevant test with its argument already bound",
			Provenance: "authored",
			Score:      1,
		})
	default:
		if input.Handle != RepoHandle && input.Handle != TestsHandle {
			return Envelope{Version: 1, Handle: input.Handle, Verb: input.Verb, Status: "absent", Text: "absent: handle is not reachable", Frontier: []FrontierEntry{}}
		}
		return refused(input.Handle, input.Verb, "verb is not attached to this handle", nil)
	}
}

func success(input ActInput, result map[string]any, text string, frontier ...FrontierEntry) Envelope {
	return Envelope{Version: 1, Handle: input.Handle, Verb: input.Verb, Status: "ok", Result: result, Text: render(text, frontier), Frontier: frontier}
}

func refused(handle, verb, why string, instead *Call) Envelope {
	refusal := &Refusal{What: handle + "." + verb, Why: why, Instead: instead}
	text := "refused " + refusal.What + ": " + why
	if instead != nil {
		text += "\ninstead: " + instead.Handle + "." + instead.Verb
	}
	return Envelope{Version: 1, Handle: handle, Verb: verb, Status: "refused", Text: text, Frontier: []FrontierEntry{}, Refusal: refusal}
}

func next(handle, verb, why string) FrontierEntry {
	return FrontierEntry{Call: Call{Handle: handle, Verb: verb, Args: map[string]any{}}, Why: why, Provenance: "authored", Score: 1}
}

func render(text string, frontier []FrontierEntry) string {
	var builder strings.Builder
	builder.WriteString(text)
	for _, entry := range frontier {
		fmt.Fprintf(&builder, "\nnext: %s.%s — %s", entry.Call.Handle, entry.Call.Verb, entry.Why)
	}
	return builder.String()
}

func parseExpression(expression string) (ActInput, bool) {
	match := expressionPattern.FindStringSubmatch(strings.TrimSpace(expression))
	if match == nil || strings.TrimSpace(match[3]) != "" {
		return ActInput{}, false
	}
	handles := map[string]string{"repo": RepoHandle, "tests": TestsHandle}
	handle, ok := handles[match[1]]
	if !ok {
		return ActInput{Handle: match[1], Verb: match[2], Args: map[string]any{}}, true
	}
	return ActInput{Handle: handle, Verb: match[2], Args: map[string]any{}}, true
}

func addFlatTools(server *mcp.Server, count int) {
	for index := 0; index < count; index++ {
		name := fmt.Sprintf("capability_%04d", index)
		mcp.AddTool(server, &mcp.Tool{
			Name:        name,
			Description: "Invoke one conventional capability with its complete upfront schema.",
		}, func(_ context.Context, _ *mcp.CallToolRequest, input struct {
			Target string `json:"target" jsonschema:"target object"`
			Value  string `json:"value,omitempty" jsonschema:"operation value"`
		}) (*mcp.CallToolResult, map[string]any, error) {
			return nil, map[string]any{"target": input.Target, "ok": true}, nil
		})
	}
}
