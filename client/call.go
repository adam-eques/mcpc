package client

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/adam-eques/mcpc/jsonrpc"
	"github.com/adam-eques/mcpc/mcp"
)

// call issues a request and waits for the correlated response. On context
// cancellation it best-effort notifies the server so the work can be aborted.
func (c *Client) call(ctx context.Context, method string, params any, out any) error {
	start := time.Now()
	err := c.callOnce(ctx, method, params, out)
	c.metrics.Observe(method, time.Since(start), err != nil)
	c.metrics.Inc("requests_total")
	if err != nil {
		c.metrics.Inc("errors_total")
	}
	return err
}

func (c *Client) callOnce(ctx context.Context, method string, params any, out any) error {
	if c.closed.Load() {
		return ErrClosed
	}
	id := c.nextID.Add(1)
	key := strconv.FormatInt(id, 10)

	raw, err := marshalParams(params)
	if err != nil {
		return err
	}
	frame, err := json.Marshal(jsonrpc.NewRequest(jsonrpc.Int64ID(id), method, raw))
	if err != nil {
		return err
	}

	ch := make(chan *jsonrpc.Response, 1)
	c.register(key, ch)
	defer c.unregister(key)

	callCtx := ctx
	if c.reqTimeout > 0 {
		var cancel context.CancelFunc
		callCtx, cancel = context.WithTimeout(ctx, c.reqTimeout)
		defer cancel()
	}

	if err := c.t.Send(callCtx, frame); err != nil {
		return err
	}

	select {
	case <-callCtx.Done():
		c.sendCancelled(id, callCtx.Err())
		return callCtx.Err()
	case <-c.done:
		if e, ok := c.readErr.Load().(error); ok && e != nil {
			return e
		}
		return ErrClosed
	case resp := <-ch:
		if resp.Error != nil {
			return resp.Error
		}
		if out != nil && len(resp.Result) > 0 {
			return json.Unmarshal(resp.Result, out)
		}
		return nil
	}
}

// notify sends a fire-and-forget notification.
func (c *Client) notify(ctx context.Context, method string, params any) error {
	raw, err := marshalParams(params)
	if err != nil {
		return err
	}
	frame, err := json.Marshal(jsonrpc.NewNotification(method, raw))
	if err != nil {
		return err
	}
	return c.t.Send(ctx, frame)
}

// sendCancelled tells the server to abort the request with the given id.
func (c *Client) sendCancelled(id int64, cause error) {
	reason := ""
	if cause != nil {
		reason = cause.Error()
	}
	// Best effort with a fresh context; the call's context is already done.
	_ = c.notify(context.Background(), mcp.NotificationCancelled, mcp.CancelledParams{
		RequestID: id,
		Reason:    reason,
	})
}

func marshalParams(params any) (json.RawMessage, error) {
	if params == nil {
		return nil, nil
	}
	if raw, ok := params.(json.RawMessage); ok {
		return raw, nil
	}
	return json.Marshal(params)
}
