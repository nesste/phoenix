package main

import (
	"context"
	"flag"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/raoul/phoenix/experiments/surface-spike/internal/spike"
)

func main() {
	surface := flag.String("surface", "act", "surface to serve: act or eval")
	flag.Parse()

	server, err := spike.NewServer(*surface, 0)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("surface spike stopped: %v", err)
	}
}
