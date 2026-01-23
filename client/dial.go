package client

import (
	"context"

	"github.com/adam-eques/mcpc/mcp"
	"github.com/adam-eques/mcpc/transport"
)

// DialCommand launches an MCP server process, starts a client over its stdio,
// and completes the initialize handshake. The returned client owns the process
// and terminates it on Close.
func DialCommand(ctx context.Context, name string, args []string, opts ...Option) (*Client, *mcp.InitializeResult, error) {
	t, err := transport.StartCommand(name, args, transport.CommandOptions{})
	if err != nil {
		return nil, nil, err
	}
	return dial(ctx, t, opts...)
}

// DialHTTP connects to a gateway endpoint (for example
// "http://localhost:8080/rpc") and completes the initialize handshake.
func DialHTTP(ctx context.Context, endpoint string, opts ...Option) (*Client, *mcp.InitializeResult, error) {
	return dial(ctx, transport.NewHTTP(endpoint), opts...)
}

func dial(ctx context.Context, t transport.Transport, opts ...Option) (*Client, *mcp.InitializeResult, error) {
	c := New(t, opts...)
	if err := c.Start(ctx); err != nil {
		c.Close()
		return nil, nil, err
	}
	info, err := c.Initialize(ctx)
	if err != nil {
		c.Close()
		return nil, nil, err
	}
	return c, info, nil
}
