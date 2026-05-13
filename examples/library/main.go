// Command example-library shows how to embed the mcpc client in a Go program.
// It launches an mcpkit server, lists its tools, and calls the calculator.
//
// Run it with the server command as the first argument:
//
//	go run ./examples/library mcpkit
//	go run ./examples/library "go run github.com/adam-eques/mcpkit/cmd/mcpkit"
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/adam-eques/mcpc/client"
	"github.com/adam-eques/mcpc/mcp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: example-library <server-command>")
		os.Exit(2)
	}
	if err := run(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(serverCmd string) error {
	fields := strings.Fields(serverCmd)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c, info, err := client.DialCommand(ctx, fields[0], fields[1:],
		client.WithClientInfo(mcp.Implementation{Name: "example-library", Version: "1.0"}))
	if err != nil {
		return err
	}
	defer c.Close()
	fmt.Printf("connected to %s %s\n", info.ServerInfo.Name, info.ServerInfo.Version)

	tools, err := c.ListTools(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("server exposes %d tools\n", len(tools))

	res, err := c.CallTool(ctx, "calculate", map[string]string{"expression": "2 ^ 16"})
	if err != nil {
		return err
	}
	fmt.Println("2 ^ 16 =", res.Content[0].(mcp.TextContent).Text)
	return nil
}
