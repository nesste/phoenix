package main

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestVersion(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run(context.Background(), []string{"version"}, io.NopCloser(strings.NewReader("")), nopWriteCloser{&stdout}, &stderr)

	if exitCode != 0 {
		t.Fatalf("run(version) exit code = %d, want 0; stderr: %s", exitCode, stderr.String())
	}
	if got, want := stdout.String(), "phoenix dev\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestServeStdioHandshake(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	serverReader, clientWriter := io.Pipe()
	clientReader, serverWriter := io.Pipe()
	var stderr bytes.Buffer
	serveDone := make(chan int, 1)
	go func() {
		serveDone <- run(ctx, []string{"serve", "--stdio"}, serverReader, serverWriter, &stderr)
	}()

	client := mcp.NewClient(&mcp.Implementation{Name: "phoenix-smoke-test", Version: "dev"}, nil)
	session, err := client.Connect(ctx, &mcp.IOTransport{Reader: clientReader, Writer: clientWriter}, nil)
	if err != nil {
		t.Fatalf("connect to stdio server: %v; stderr: %s", err, stderr.String())
	}
	result := session.InitializeResult()
	if result == nil || result.ServerInfo == nil {
		t.Fatal("initialize result has no server info")
	}
	if got, want := result.ServerInfo.Name, "phoenix"; got != want {
		t.Fatalf("server name = %q, want %q", got, want)
	}
	if got, want := result.ServerInfo.Version, "dev"; got != want {
		t.Fatalf("server version = %q, want %q", got, want)
	}
	if result.Instructions != "" {
		t.Fatalf("standing instructions = %q, want empty", result.Instructions)
	}
	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	if len(tools.Tools) != 1 || tools.Tools[0].Name != "act" {
		t.Fatalf("served tools = %#v, want only act", tools.Tools)
	}

	if err := session.Close(); err != nil {
		t.Fatalf("close client session: %v", err)
	}
	select {
	case exitCode := <-serveDone:
		if exitCode != 0 {
			t.Fatalf("run(serve --stdio) exit code = %d, want 0; stderr: %s", exitCode, stderr.String())
		}
	case <-ctx.Done():
		t.Fatalf("stdio server did not stop after client close: %v", ctx.Err())
	}
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }
