package client

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/adam-eques/mcpc/internal/log"
	"github.com/adam-eques/mcpc/internal/metrics"
	"github.com/adam-eques/mcpc/jsonrpc"
	"github.com/adam-eques/mcpc/mcp"
	"github.com/adam-eques/mcpc/transport"
)

// Client speaks the Model Context Protocol to a single server over a transport.
type Client struct {
	t       transport.Transport
	log     *log.Logger
	metrics *metrics.Registry

	reqTimeout time.Duration
	clientInfo mcp.Implementation
	protocol   string
	onNotify   NotificationHandler
	maxRetries int

	nextID  atomic.Int64
	mu      sync.Mutex
	pending map[string]chan *jsonrpc.Response

	started atomic.Bool
	closed  atomic.Bool
	done    chan struct{}
	readErr atomic.Value // error

	sessMu      sync.RWMutex
	initialized bool
	serverInfo  mcp.Implementation
	serverCaps  mcp.ServerCapabilities
}

// New returns a Client over t. Call Start before issuing requests.
func New(t transport.Transport, opts ...Option) *Client {
	c := &Client{
		t:          t,
		log:        log.Discard(),
		metrics:    metrics.New(),
		reqTimeout: 30 * time.Second,
		clientInfo: mcp.Implementation{Name: "mcpc", Version: "0.1.0"},
		protocol:   mcp.ProtocolVersion,
		pending:    make(map[string]chan *jsonrpc.Response),
		done:       make(chan struct{}),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Start launches the background read loop. It returns once the loop is running;
// the loop stops when the transport closes, ctx is cancelled, or Close is called.
func (c *Client) Start(ctx context.Context) error {
	if !c.started.CompareAndSwap(false, true) {
		return ErrAlreadyStarted
	}
	go c.readLoop(ctx)
	return nil
}

func (c *Client) readLoop(ctx context.Context) {
	defer close(c.done)
	for {
		frame, err := c.t.Receive(ctx)
		if err != nil {
			c.readErr.Store(err)
			c.failAll(err)
			return
		}
		if len(frame) > 0 {
			c.handleFrame(frame)
		}
	}
}

// envelope decodes just enough to tell a response from a server-initiated
// request or notification.
type envelope struct {
	ID     *jsonrpc.ID     `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
	Result json.RawMessage `json:"result"`
	Error  *jsonrpc.Error  `json:"error"`
}

func (c *Client) handleFrame(frame []byte) {
	var env envelope
	if err := json.Unmarshal(frame, &env); err != nil {
		c.log.Warn("dropping malformed frame", "err", err)
		return
	}
	if env.Method != "" {
		if env.ID != nil {
			c.handleServerRequest(&env)
		} else {
			c.log.Debug("server notification", "method", env.Method)
			if c.onNotify != nil {
				c.onNotify(env.Method, env.Params)
			}
		}
		return
	}
	if env.ID == nil {
		return
	}
	c.deliver(env.ID.String(), &jsonrpc.Response{
		JSONRPC: jsonrpc.Version,
		ID:      env.ID,
		Result:  env.Result,
		Error:   env.Error,
	})
}

// handleServerRequest answers the few requests a server may send to a client.
// Unknown methods get a method-not-found error so the server is not left waiting.
func (c *Client) handleServerRequest(env *envelope) {
	var resp *jsonrpc.Response
	switch env.Method {
	case mcp.MethodPing:
		resp = jsonrpc.NewResponse(env.ID, json.RawMessage(`{}`))
	default:
		resp = jsonrpc.NewErrorResponse(env.ID, jsonrpc.MethodNotFound(env.Method))
	}
	frame, err := json.Marshal(resp)
	if err != nil {
		return
	}
	if err := c.t.Send(context.Background(), frame); err != nil {
		c.log.Warn("failed to answer server request", "err", err)
	}
}

func (c *Client) register(key string, ch chan *jsonrpc.Response) {
	c.mu.Lock()
	c.pending[key] = ch
	c.mu.Unlock()
}

func (c *Client) unregister(key string) {
	c.mu.Lock()
	delete(c.pending, key)
	c.mu.Unlock()
}

func (c *Client) deliver(key string, resp *jsonrpc.Response) {
	c.mu.Lock()
	ch, ok := c.pending[key]
	if ok {
		delete(c.pending, key)
	}
	c.mu.Unlock()
	if ok {
		ch <- resp // buffered, never blocks
	}
}

func (c *Client) failAll(err error) {
	c.mu.Lock()
	pending := c.pending
	c.pending = make(map[string]chan *jsonrpc.Response)
	c.mu.Unlock()
	for _, ch := range pending {
		ch <- jsonrpc.NewErrorResponse(nil, jsonrpc.InternalError(err.Error()))
	}
}

// Metrics returns the client's metrics registry.
func (c *Client) Metrics() *metrics.Registry { return c.metrics }

// ServerInfo returns the server identity learned during initialize.
func (c *Client) ServerInfo() mcp.Implementation {
	c.sessMu.RLock()
	defer c.sessMu.RUnlock()
	return c.serverInfo
}

// Capabilities returns the server capabilities learned during initialize.
func (c *Client) Capabilities() mcp.ServerCapabilities {
	c.sessMu.RLock()
	defer c.sessMu.RUnlock()
	return c.serverCaps
}

// Close shuts down the transport and unblocks any in-flight calls.
func (c *Client) Close() error {
	c.closed.Store(true)
	return c.t.Close()
}
