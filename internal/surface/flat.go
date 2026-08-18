package surface

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nesste/phoenix/internal/world"
)

// FlatMCP exposes every root verb as a conventional upfront MCP tool while
// executing through the same admission path used by the Phoenix surface.
type FlatMCP struct {
	server    *mcp.Server
	admission *Admission
	sessions  sync.Map
}

type flatTool struct {
	name        string
	root        string
	verb        string
	description string
	schema      json.RawMessage
}

type flatSession struct {
	admission *Session
	roots     map[string]string
}

func NewFlatMCP(version string, admission *Admission, definition *world.Definition) (*FlatMCP, error) {
	if admission == nil {
		return nil, fmt.Errorf("surface admission is required")
	}
	tools, err := flatTools(definition)
	if err != nil {
		return nil, err
	}
	adapter := &FlatMCP{admission: admission}
	adapter.server = mcp.NewServer(
		&mcp.Implementation{Name: "phoenix-flat", Title: "Phoenix flat baseline", Version: version},
		&mcp.ServerOptions{Capabilities: &mcp.ServerCapabilities{Tools: &mcp.ToolCapabilities{}}},
	)
	for _, tool := range tools {
		tool := tool
		adapter.server.AddTool(&mcp.Tool{
			Name: tool.name, Description: tool.description, InputSchema: tool.schema,
		}, func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return adapter.handle(ctx, request, tool)
		})
	}
	return adapter, nil
}

func (adapter *FlatMCP) Server() *mcp.Server {
	return adapter.server
}

func flatTools(definition *world.Definition) ([]flatTool, error) {
	if definition == nil {
		return nil, fmt.Errorf("world definition is required")
	}
	tools := []flatTool{}
	seen := map[string]struct{}{}
	for _, root := range definition.Roots {
		descriptor, exists := definition.HandleTypes[root.Type]
		if !exists {
			return nil, fmt.Errorf("root %q has unknown handle type %q", root.Name, root.Type)
		}
		for verb, declaration := range descriptor.Verbs {
			name := root.Name + "_" + verb
			if _, duplicate := seen[name]; duplicate {
				return nil, fmt.Errorf("flat tool name %q is duplicated", name)
			}
			seen[name] = struct{}{}
			tools = append(tools, flatTool{
				name: name, root: root.Name, verb: verb,
				description: fmt.Sprintf("Invoke %s.%s on %s.", root.Type, verb, root.Label),
				schema:      append(json.RawMessage(nil), declaration.ArgsSchema...),
			})
		}
	}
	sort.Slice(tools, func(i, j int) bool { return tools[i].name < tools[j].name })
	return tools, nil
}

// FlatToolNames returns the deterministic upfront tool set for runner allowlists
// and artifact freezing. It is derived from the same declarations as FlatMCP.
func FlatToolNames(definition *world.Definition) ([]string, error) {
	tools, err := flatTools(definition)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(tools))
	for index, tool := range tools {
		names[index] = tool.name
	}
	return names, nil
}

func (adapter *FlatMCP) handle(ctx context.Context, request *mcp.CallToolRequest, tool flatTool) (*mcp.CallToolResult, error) {
	if request == nil || request.Params == nil || request.Session == nil {
		return toolError("invalid tool arguments"), nil
	}
	arguments, err := decodeFlatArguments(request.Params.Arguments)
	if err != nil {
		return toolError("invalid tool arguments"), nil
	}
	session, err := adapter.session(request.Session)
	if err != nil {
		return toolError("tool session is unavailable"), nil
	}
	handle, exists := session.roots[tool.root]
	if !exists {
		return toolError("tool session is unavailable"), nil
	}
	envelope := adapter.admission.Act(ctx, session.admission, Input{
		Handle: handle, Verb: tool.verb, Args: arguments,
	})
	structured := envelope.Result
	if envelope.Status != StatusOK {
		structured = map[string]any{
			"status": envelope.Status, "error": envelope.Error,
		}
	}
	return &mcp.CallToolResult{
		Content:           []mcp.Content{&mcp.TextContent{Text: envelope.Text}},
		StructuredContent: structured, IsError: false,
	}, nil
}

func decodeFlatArguments(raw json.RawMessage) (map[string]any, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var arguments map[string]any
	if err := decoder.Decode(&arguments); err != nil || arguments == nil {
		return nil, fmt.Errorf("arguments must be a JSON object")
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return nil, fmt.Errorf("trailing JSON data")
	}
	return arguments, nil
}

func (adapter *FlatMCP) session(serverSession *mcp.ServerSession) (*flatSession, error) {
	if existing, ok := adapter.sessions.Load(serverSession); ok {
		return existing.(*flatSession), nil
	}
	createdAdmission, roots, err := adapter.admission.StartSession()
	if err != nil {
		return nil, err
	}
	created := &flatSession{admission: createdAdmission, roots: make(map[string]string, len(roots))}
	for _, root := range roots {
		created.roots[root.Name] = root.Ref
	}
	actual, loaded := adapter.sessions.LoadOrStore(serverSession, created)
	if !loaded {
		go func() {
			_ = serverSession.Wait()
			adapter.sessions.CompareAndDelete(serverSession, created)
		}()
	}
	return actual.(*flatSession), nil
}
