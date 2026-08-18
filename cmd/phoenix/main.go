package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nesste/phoenix/internal/frontier"
	"github.com/nesste/phoenix/internal/surface"
	"github.com/nesste/phoenix/internal/teach"
	"github.com/nesste/phoenix/internal/verb"
	"github.com/nesste/phoenix/internal/world"
)

var version = "dev"

func main() {
	os.Exit(run(context.Background(), os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(ctx context.Context, args []string, stdin io.ReadCloser, stdout io.WriteCloser, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	switch args[0] {
	case "version":
		if len(args) != 1 {
			fmt.Fprintln(stderr, "phoenix version takes no arguments")
			return 2
		}
		fmt.Fprintf(stdout, "phoenix %s\n", version)
		return 0
	case "serve":
		return runServe(ctx, args[1:], stdin, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		printUsage(stderr)
		return 2
	}
}

func runServe(ctx context.Context, args []string, stdin io.ReadCloser, stdout io.WriteCloser, stderr io.Writer) int {
	flags := flag.NewFlagSet("phoenix serve", flag.ContinueOnError)
	flags.SetOutput(stderr)
	stdio := flags.Bool("stdio", false, "serve Model Context Protocol over stdin and stdout")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 || !*stdio {
		fmt.Fprintln(stderr, "usage: phoenix serve --stdio")
		return 2
	}

	server, err := defaultSurface(version)
	if err != nil {
		fmt.Fprintf(stderr, "configure surface: %v\n", err)
		return 1
	}
	transport := &mcp.IOTransport{Reader: stdin, Writer: stdout}
	if err := server.Server().Run(ctx, transport); err != nil {
		fmt.Fprintf(stderr, "serve stdio: %v\n", err)
		return 1
	}
	return 0
}

const unassembledWorldBuild = "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

func defaultSurface(serverVersion string) (*surface.MCP, error) {
	definition := &world.Definition{
		V: 1, ID: "unassembled", Roots: []world.Root{},
		HandleTypes: map[string]world.HandleType{}, Transitions: []world.Transition{},
	}
	frontierEngine, err := frontier.New(definition)
	if err != nil {
		return nil, err
	}
	teachingEngine, err := teach.New(definition)
	if err != nil {
		return nil, err
	}
	admission, err := surface.New(surface.Config{
		WorldBuild: unassembledWorldBuild,
		Graph:      world.NewGraph(definition, emptyResolver{}),
		Executor:   verb.NewExecutor(verb.NewRegistry(), verb.Options{}),
		Frontier:   frontierEngine,
		Teacher:    teachingEngine,
	})
	if err != nil {
		return nil, err
	}
	return surface.NewMCP(serverVersion, admission), nil
}

type emptyResolver struct{}

func (emptyResolver) Resolve(context.Context, world.Resource) (any, error) {
	return map[string]any{}, nil
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: phoenix <version|serve --stdio>")
}
