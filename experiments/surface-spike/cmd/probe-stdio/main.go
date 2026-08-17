package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/raoul/phoenix/experiments/surface-spike/internal/spike"
)

func main() {
	executable := flag.String("exe", "./bin/phoenix-spike.exe", "path to the spike server executable")
	surface := flag.String("surface", "act", "surface to probe: act or eval")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := mcp.NewClient(&mcp.Implementation{Name: "stdio-probe", Version: "0.1.0"}, nil)
	session, err := client.Connect(ctx, &mcp.CommandTransport{Command: exec.Command(*executable, "--surface", *surface)}, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	params := &mcp.CallToolParams{Name: "act", Arguments: map[string]any{"handle": spike.RepoHandle, "verb": "status", "args": map[string]any{}}}
	if *surface == "eval" {
		params = &mcp.CallToolParams{Name: "eval", Arguments: map[string]any{"expression": "repo.status()"}}
	}
	started := time.Now()
	result, err := session.CallTool(ctx, params)
	if err != nil {
		log.Fatal(err)
	}
	output := struct {
		Surface   string              `json:"surface"`
		LatencyUS int64               `json:"latency_us"`
		Result    *mcp.CallToolResult `json:"result"`
	}{Surface: *surface, LatencyUS: time.Since(started).Microseconds(), Result: result}
	encoded, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(encoded))
}
