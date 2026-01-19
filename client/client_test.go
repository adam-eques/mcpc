package client

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/adam-eques/mcpc/jsonrpc"
	"github.com/adam-eques/mcpc/mcp"
	"github.com/adam-eques/mcpc/transport"
)

func startClient(t *testing.T, fs *fakeServer, end transport.Transport) (*Client, context.CancelFunc) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	go fs.run(ctx)
	c := New(end)
	if err := c.Start(ctx); err != nil {
		t.Fatal(err)
	}
	return c, cancel
}

func TestInitializeHandshake(t *testing.T) {
	fs, end := newFakeServer()
	c, cancel := startClient(t, fs, end)
	defer cancel()
	defer c.Close()

	res, err := c.Initialize(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if res.ServerInfo.Name != "fake" {
		t.Fatalf("serverInfo=%+v", res.ServerInfo)
	}
	if c.Capabilities().Tools == nil {
		t.Fatal("tools capability not recorded")
	}
}

func TestListAndCallTool(t *testing.T) {
	fs, end := newFakeServer()
	fs.handlers[mcp.MethodToolsList] = func(json.RawMessage) (any, *jsonrpc.Error) {
		return mcp.ListToolsResult{Tools: []mcp.Tool{{Name: "calculate"}}}, nil
	}
	fs.handlers[mcp.MethodToolsCall] = func(params json.RawMessage) (any, *jsonrpc.Error) {
		var p mcp.CallToolParams
		json.Unmarshal(params, &p)
		if p.Name != "calculate" {
			return nil, jsonrpc.InvalidParams("unknown tool")
		}
		return mcp.TextResult("42"), nil
	}
	c, cancel := startClient(t, fs, end)
	defer cancel()
	defer c.Close()

	if _, err := c.Initialize(context.Background()); err != nil {
		t.Fatal(err)
	}
	tools, err := c.ListTools(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(tools) != 1 || tools[0].Name != "calculate" {
		t.Fatalf("tools=%+v", tools)
	}
	res, err := c.CallTool(context.Background(), "calculate", map[string]string{"expression": "6*7"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Content[0].(mcp.TextContent).Text != "42" {
		t.Fatalf("result=%+v", res)
	}
}

func TestServerErrorPropagates(t *testing.T) {
	fs, end := newFakeServer()
	c, cancel := startClient(t, fs, end)
	defer cancel()
	defer c.Close()
	c.Initialize(context.Background())

	_, err := c.ListTools(context.Background())
	var rpcErr *jsonrpc.Error
	if !errors.As(err, &rpcErr) || rpcErr.Code != jsonrpc.CodeMethodNotFound {
		t.Fatalf("expected method-not-found error, got %v", err)
	}
}

func TestConcurrentCalls(t *testing.T) {
	fs, end := newFakeServer()
	fs.handlers[mcp.MethodPing] = func(json.RawMessage) (any, *jsonrpc.Error) {
		return mcp.PingResult{}, nil
	}
	c, cancel := startClient(t, fs, end)
	defer cancel()
	defer c.Close()
	c.Initialize(context.Background())

	errs := make(chan error, 20)
	for i := 0; i < 20; i++ {
		go func() { errs <- c.Ping(context.Background()) }()
	}
	for i := 0; i < 20; i++ {
		if err := <-errs; err != nil {
			t.Fatalf("concurrent ping failed: %v", err)
		}
	}
}

func TestRequestTimeout(t *testing.T) {
	fs, end := newFakeServer()
	fs.blockOn = mcp.MethodToolsList
	fs.handlers[mcp.MethodToolsList] = func(json.RawMessage) (any, *jsonrpc.Error) {
		return mcp.ListToolsResult{}, nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go fs.run(ctx)
	c := New(end, WithRequestTimeout(100*time.Millisecond))
	c.Start(ctx)
	defer c.Close()
	c.Initialize(context.Background())

	_, err := c.ListTools(context.Background())
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
}
