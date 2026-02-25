package client

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/adam-eques/mcpc/jsonrpc"
	"github.com/adam-eques/mcpc/mcp"
)

func TestCallToolRetrySucceedsEventually(t *testing.T) {
	fs, end := newFakeServer()
	var calls int
	fs.handlers[mcp.MethodToolsCall] = func(json.RawMessage) (any, *jsonrpc.Error) {
		calls++
		if calls < 3 {
			return nil, jsonrpc.InternalError("temporary")
		}
		return mcp.TextResult("ok"), nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go fs.run(ctx)
	c := New(end, WithMaxRetries(5))
	c.Start(ctx)
	defer c.Close()
	c.Initialize(context.Background())

	res, err := c.CallToolRetry(context.Background(), "x", nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Content[0].(mcp.TextContent).Text != "ok" {
		t.Fatalf("res=%+v", res)
	}
	if calls != 3 {
		t.Fatalf("expected 3 attempts, got %d", calls)
	}
}

func TestRetryDoesNotRetryProtocolErrors(t *testing.T) {
	fs, end := newFakeServer()
	var calls int
	fs.handlers[mcp.MethodToolsCall] = func(json.RawMessage) (any, *jsonrpc.Error) {
		calls++
		return nil, jsonrpc.InvalidParams("bad")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go fs.run(ctx)
	c := New(end, WithMaxRetries(5))
	c.Start(ctx)
	defer c.Close()
	c.Initialize(context.Background())

	if _, err := c.CallToolRetry(context.Background(), "x", nil); err == nil {
		t.Fatal("expected error")
	}
	if calls != 1 {
		t.Fatalf("invalid params must not be retried, got %d calls", calls)
	}
}

func TestMetricsRecorded(t *testing.T) {
	fs, end := newFakeServer()
	c, cancel := startClient(t, fs, end)
	defer cancel()
	defer c.Close()
	c.Initialize(context.Background())
	c.Ping(context.Background())
	if c.Metrics().Snapshot().Counters["requests_total"] < 2 {
		t.Fatal("expected requests recorded in metrics")
	}
}
