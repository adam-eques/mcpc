package client

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/adam-eques/mcpc/jsonrpc"
	"github.com/adam-eques/mcpc/mcp"
)

func BenchmarkCallTool(b *testing.B) {
	fs, end := newFakeServer()
	fs.handlers[mcp.MethodToolsCall] = func(json.RawMessage) (any, *jsonrpc.Error) {
		return mcp.TextResult("ok"), nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go fs.run(ctx)
	c := New(end)
	c.Start(ctx)
	defer c.Close()
	c.Initialize(context.Background())

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := c.CallTool(ctx, "x", map[string]string{"a": "b"}); err != nil {
			b.Fatal(err)
		}
	}
}
