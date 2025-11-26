// Package transport moves raw JSON-RPC frames between an MCP client and a server.
// A frame is a single, complete JSON-RPC message with no embedded newlines.
//
// The client ships three transports: Command spawns a server process and speaks
// newline-delimited JSON over its stdin/stdout, HTTP posts each request to a
// gateway's /rpc endpoint, and Pipe connects two in-memory endpoints for tests.
package transport

import "context"

// Transport is a bidirectional stream of JSON-RPC frames.
type Transport interface {
	// Receive blocks until the next frame arrives, ctx is cancelled, or the peer
	// closes the stream (reported as io.EOF).
	Receive(ctx context.Context) ([]byte, error)

	// Send writes a single frame. Implementations must be safe for concurrent use
	// by multiple goroutines.
	Send(ctx context.Context, frame []byte) error

	// Close releases the underlying resources and unblocks any pending Receive.
	Close() error
}
