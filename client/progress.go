package client

import (
	"encoding/json"
	"sync"

	"github.com/adam-eques/mcpc/mcp"
)

// ProgressFunc receives progress updates for a tracked token.
type ProgressFunc func(mcp.ProgressParams)

// ProgressTracker dispatches notifications/progress to per-token callbacks. Wire
// it up by passing its Handle method to WithNotificationHandler.
type ProgressTracker struct {
	mu       sync.Mutex
	handlers map[string]ProgressFunc
}

// NewProgressTracker returns an empty tracker.
func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{handlers: make(map[string]ProgressFunc)}
}

// Track registers fn for the given progress token.
func (p *ProgressTracker) Track(token string, fn ProgressFunc) {
	p.mu.Lock()
	p.handlers[token] = fn
	p.mu.Unlock()
}

// Forget removes the callback for a token.
func (p *ProgressTracker) Forget(token string) {
	p.mu.Lock()
	delete(p.handlers, token)
	p.mu.Unlock()
}

// Handle is a NotificationHandler that routes progress notifications to the
// registered callbacks; it ignores everything else.
func (p *ProgressTracker) Handle(method string, params json.RawMessage) {
	if method != mcp.NotificationProgress {
		return
	}
	var pp mcp.ProgressParams
	if err := json.Unmarshal(params, &pp); err != nil {
		return
	}
	token := tokenString(pp.ProgressToken)
	p.mu.Lock()
	fn := p.handlers[token]
	p.mu.Unlock()
	if fn != nil {
		fn(pp)
	}
}

func tokenString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return trimFloat(t)
	default:
		return ""
	}
}

func trimFloat(f float64) string {
	b, _ := json.Marshal(f)
	return string(b)
}
