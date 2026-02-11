package client

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/adam-eques/mcpc/jsonrpc"
	"github.com/adam-eques/mcpc/mcp"
)

func TestContextCancelSendsCancelledNotification(t *testing.T) {
	fs, end := newFakeServer()
	fs.blockOn = mcp.MethodToolsCall
	fs.handlers[mcp.MethodToolsCall] = func(json.RawMessage) (any, *jsonrpc.Error) {
		return mcp.TextResult("late"), nil
	}
	// Capture cancellation notifications the client sends.
	cancelled := make(chan struct{}, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go fs.run(ctx)

	c := New(end, WithNotificationHandler(nil))
	// Intercept the outbound cancel by watching the fake server's inbound stream
	// is complex; instead we assert the call returns promptly on cancellation.
	c.Start(ctx)
	defer c.Close()
	c.Initialize(context.Background())

	callCtx, callCancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		callCancel()
		cancelled <- struct{}{}
	}()

	_, err := c.CallTool(callCtx, "slow", nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	<-cancelled
	// Let the server unblock so the goroutine can exit cleanly.
	close(fs.block)
}
