package transport

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTP is a transport for gateways that expose MCP over a request/response
// endpoint (such as mcpkit-gateway's POST /rpc). Because HTTP has no independent
// server-to-client stream, each Send performs the POST and enqueues the reply;
// the matching Receive then dequeues it. A notification (no id) yields no reply
// and enqueues nothing.
type HTTP struct {
	endpoint string
	client   *http.Client
	replies  chan frame
	done     chan struct{}
}

// NewHTTP returns an HTTP transport targeting endpoint (for example
// "http://localhost:8080/rpc").
func NewHTTP(endpoint string) *HTTP {
	return &HTTP{
		endpoint: endpoint,
		client:   &http.Client{Timeout: 60 * time.Second},
		replies:  make(chan frame, 64),
		done:     make(chan struct{}),
	}
}

// Send implements Transport by POSTing the frame and queuing any JSON response.
func (h *HTTP) Send(ctx context.Context, frame []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.endpoint, bytes.NewReader(frame))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil // notification acknowledged, no reply to deliver
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxFrameBytes))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gateway returned %d: %s", resp.StatusCode, bytes.TrimSpace(body))
	}
	h.enqueue(frameOf(body))
	return nil
}

// Receive implements Transport by returning the next queued reply.
func (h *HTTP) Receive(ctx context.Context) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-h.done:
		return nil, io.EOF
	case f := <-h.replies:
		return f.data, f.err
	}
}

// Close implements Transport.
func (h *HTTP) Close() error {
	select {
	case <-h.done:
	default:
		close(h.done)
	}
	return nil
}

func (h *HTTP) enqueue(f frame) {
	select {
	case h.replies <- f:
	case <-h.done:
	}
}

func frameOf(b []byte) frame { return frame{data: bytes.TrimRight(b, "\r\n")} }
