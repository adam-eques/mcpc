package client

import (
	"context"
	"testing"
	"time"
)

func TestDialCommandMissingBinary(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, _, err := DialCommand(ctx, "mcpc-no-such-binary-xyz", nil); err == nil {
		t.Fatal("expected error dialling a missing binary")
	}
}

func TestDialHTTPUnreachable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Port 1 is not listening; initialize must fail rather than hang.
	if _, _, err := DialHTTP(ctx, "http://127.0.0.1:1/rpc"); err == nil {
		t.Fatal("expected error dialling an unreachable endpoint")
	}
}
