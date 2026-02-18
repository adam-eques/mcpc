package client

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/adam-eques/mcpc/mcp"
)

// TestPoolLive exercises the Pool against a real server; skipped without
// MCPKIT_BIN so the package still builds and tests without the paired server.
func TestPoolLive(t *testing.T) {
	bin := os.Getenv("MCPKIT_BIN")
	if bin == "" {
		t.Skip("set MCPKIT_BIN to run the pool integration test")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := NewPool(ctx, 3, bin, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	if pool.Size() != 3 {
		t.Fatalf("size=%d", pool.Size())
	}

	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := pool.CallTool(ctx, "calculate", map[string]string{"expression": "1+1"})
			if err != nil {
				t.Errorf("call: %v", err)
				return
			}
			if res.Content[0].(mcp.TextContent).Text != "2" {
				t.Errorf("unexpected result: %+v", res)
			}
		}()
	}
	wg.Wait()
}
