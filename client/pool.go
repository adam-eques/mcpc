package client

import (
	"context"
	"sync"

	"github.com/adam-eques/mcpc/mcp"
)

// Pool maintains several independent client connections to the same server
// command and hands them out round-robin. Because a single Client already
// multiplexes concurrent requests, a Pool is only needed when a server process
// itself is the bottleneck and you want several of them.
type Pool struct {
	mu      sync.Mutex
	clients []*Client
	next    int
}

// NewPool dials size server processes and initializes each. The caller owns the
// pool and must Close it.
func NewPool(ctx context.Context, size int, command string, args []string, opts ...Option) (*Pool, error) {
	if size < 1 {
		size = 1
	}
	p := &Pool{}
	for i := 0; i < size; i++ {
		c, _, err := DialCommand(ctx, command, args, opts...)
		if err != nil {
			p.Close()
			return nil, err
		}
		p.clients = append(p.clients, c)
	}
	return p, nil
}

// Get returns the next client in round-robin order.
func (p *Pool) Get() *Client {
	p.mu.Lock()
	defer p.mu.Unlock()
	c := p.clients[p.next%len(p.clients)]
	p.next++
	return c
}

// Size returns the number of pooled clients.
func (p *Pool) Size() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.clients)
}

// CallTool dispatches a tool call to the next client in the pool.
func (p *Pool) CallTool(ctx context.Context, name string, args any) (*mcp.CallToolResult, error) {
	return p.Get().CallTool(ctx, name, args)
}

// Close shuts down every pooled client.
func (p *Pool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	var firstErr error
	for _, c := range p.clients {
		if err := c.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	p.clients = nil
	return firstErr
}
