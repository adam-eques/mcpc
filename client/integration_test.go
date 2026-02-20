package client

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/adam-eques/mcpc/mcp"
)

// TestLiveServer exercises the client against a real MCP server. It is skipped
// unless MCPKIT_BIN points at an mcpkit binary, so the package still builds and
// tests without the paired server present.
func TestLiveServer(t *testing.T) {
	bin := os.Getenv("MCPKIT_BIN")
	if bin == "" {
		t.Skip("set MCPKIT_BIN to the mcpkit binary to run the live integration test")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c, info, err := DialCommand(ctx, bin, nil, WithClientInfo(mcp.Implementation{Name: "itest", Version: "1"}))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if info.ServerInfo.Name == "" {
		t.Fatal("server did not identify itself")
	}

	if err := c.Ping(ctx); err != nil {
		t.Fatalf("ping: %v", err)
	}
	tools, err := c.ListTools(ctx)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	if len(tools) == 0 {
		t.Fatal("expected the server to advertise tools")
	}

	res, err := c.CallTool(ctx, "calculate", map[string]string{"expression": "6*7"})
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	if res.Content[0].(mcp.TextContent).Text != "42" {
		t.Fatalf("unexpected calculate result: %+v", res)
	}
}
