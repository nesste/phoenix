package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/raoul/phoenix/experiments/surface-spike/internal/spike"
)

type measurement struct {
	Surface     string `json:"surface"`
	WorldSize   int    `json:"world_size"`
	ToolCount   int    `json:"tool_count"`
	SchemaBytes int    `json:"schema_bytes"`
}

func main() {
	var measurements []measurement
	for _, surface := range []string{"act", "eval", "flat"} {
		for _, size := range []int{10, 100, 1000} {
			toolCount, schemaBytes, err := measure(surface, size)
			if err != nil {
				log.Fatal(err)
			}
			measurements = append(measurements, measurement{Surface: surface, WorldSize: size, ToolCount: toolCount, SchemaBytes: schemaBytes})
		}
	}
	encoded, err := json.MarshalIndent(measurements, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(encoded))
}

func measure(surface string, size int) (int, int, error) {
	ctx := context.Background()
	server, err := spike.NewServer(surface, size)
	if err != nil {
		return 0, 0, err
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "schema-measure", Version: "0.1.0"}, nil)
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		return 0, 0, err
	}
	defer serverSession.Close()
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		return 0, 0, err
	}
	defer clientSession.Close()
	tools, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		return 0, 0, err
	}
	encoded, err := json.Marshal(tools.Tools)
	if err != nil {
		return 0, 0, err
	}
	return len(tools.Tools), len(encoded), nil
}
