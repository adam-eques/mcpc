package client

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/adam-eques/mcpc/jsonrpc"
	"github.com/adam-eques/mcpc/mcp"
	"github.com/adam-eques/mcpc/transport"
)

func TestServerNotificationDispatched(t *testing.T) {
	clientEnd, serverEnd := transport.Pipe()
	got := make(chan string, 1)
	c := New(clientEnd, WithNotificationHandler(func(method string, _ json.RawMessage) {
		got <- method
	}))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c.Start(ctx)
	defer c.Close()

	// The server pushes a notification with no id.
	note := jsonrpc.NewNotification(mcp.NotificationMessage, json.RawMessage(`{"level":"info","data":"hi"}`))
	frame, _ := json.Marshal(note)
	serverEnd.Send(ctx, frame)

	select {
	case method := <-got:
		if method != mcp.NotificationMessage {
			t.Fatalf("method=%s", method)
		}
	case <-time.After(time.Second):
		t.Fatal("notification not dispatched")
	}
}

func TestServerPingAnswered(t *testing.T) {
	clientEnd, serverEnd := transport.Pipe()
	c := New(clientEnd)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c.Start(ctx)
	defer c.Close()

	// The server sends a ping request; the client must answer it.
	req := jsonrpc.NewRequest(jsonrpc.Int64ID(99), mcp.MethodPing, nil)
	frame, _ := json.Marshal(req)
	serverEnd.Send(ctx, frame)

	reply, err := serverEnd.Receive(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var resp jsonrpc.Response
	json.Unmarshal(reply, &resp)
	if resp.Error != nil {
		t.Fatalf("ping answered with error: %v", resp.Error)
	}
}
