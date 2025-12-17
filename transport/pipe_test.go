package transport

import (
	"context"
	"errors"
	"io"
	"testing"
)

func TestPipe(t *testing.T) {
	client, server := Pipe()
	ctx := context.Background()
	if err := client.Send(ctx, []byte(`{"ping":1}`)); err != nil {
		t.Fatal(err)
	}
	got, err := server.Receive(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"ping":1}` {
		t.Fatalf("got=%s", got)
	}
	client.Close()
	if _, err := server.Receive(ctx); !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF, got %v", err)
	}
}
