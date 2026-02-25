package client

import (
	"context"
	"errors"
	"time"

	"github.com/adam-eques/mcpc/jsonrpc"
	"github.com/adam-eques/mcpc/mcp"
)

// CallToolRetry is like CallTool but retries when the transport or server
// returns a transient error, backing off between attempts. Protocol errors
// (invalid params, method not found) are not retried because retrying cannot
// help. Context cancellation stops retrying immediately.
func (c *Client) CallToolRetry(ctx context.Context, name string, args any) (*mcp.CallToolResult, error) {
	var lastErr error
	attempts := c.maxRetries + 1
	for attempt := 0; attempt < attempts; attempt++ {
		res, err := c.CallTool(ctx, name, args)
		if err == nil {
			return res, nil
		}
		lastErr = err
		if !isRetryable(err) {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff(attempt)):
		}
		c.log.Debug("retrying tool call", "tool", name, "attempt", attempt+1, "err", err)
	}
	return nil, lastErr
}

// isRetryable reports whether an error is worth retrying. Protocol-level errors
// with a defined JSON-RPC code are treated as permanent.
func isRetryable(err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	var rpcErr *jsonrpc.Error
	if errors.As(err, &rpcErr) {
		return rpcErr.Code == jsonrpc.CodeInternalError
	}
	return true // transport-level failure
}

func backoff(attempt int) time.Duration {
	d := time.Duration(100<<attempt) * time.Millisecond
	if d > 2*time.Second {
		return 2 * time.Second
	}
	return d
}
