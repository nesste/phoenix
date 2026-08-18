package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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
	defaults := defaultServePaths()
	flags := flag.NewFlagSet("phoenix serve", flag.ContinueOnError)
	flags.SetOutput(stderr)
	stdio := flags.Bool("stdio", false, "serve Model Context Protocol over stdin and stdout")
	worldPath := flags.String("world", defaults.world, "world definition path")
	schemaPath := flags.String("world-schema", defaults.schema, "world schema path")
	episodePath := flags.String("episode-db", defaults.episodes, "episode database path; empty disables logging")
	rootRefsPath := flags.String("root-refs", "", "isolated runner root-reference JSON path")
	worldBuild := flags.String("world-build", "", "complete world-build digest; defaults to the world definition digest")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 || !*stdio {
		fmt.Fprintln(stderr, "usage: phoenix serve --stdio")
		return 2
	}

	assembled, err := assembleSurface(serveOptions{
		serverVersion: version, worldPath: *worldPath, schemaPath: *schemaPath,
		episodePath: *episodePath, rootRefsPath: *rootRefsPath, worldBuild: *worldBuild,
		warning: stderr,
	})
	if err != nil {
		fmt.Fprintf(stderr, "configure surface: %v\n", err)
		return 1
	}
	defer assembled.Close()
	transport := &mcp.IOTransport{Reader: stdin, Writer: stdout}
	if err := assembled.server.Server().Run(ctx, transport); err != nil {
		fmt.Fprintf(stderr, "serve stdio: %v\n", err)
		return 1
	}
	return 0
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: phoenix <version|serve --stdio>")
}
