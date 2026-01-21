package client

import (
	"encoding/json"
	"time"

	"github.com/adam-eques/mcpc/internal/log"
	"github.com/adam-eques/mcpc/mcp"
)

// NotificationHandler is invoked for each server-initiated notification, such as
// progress updates or log messages.
type NotificationHandler func(method string, params json.RawMessage)

// Option customises a Client.
type Option func(*Client)

// WithClientInfo sets the implementation name and version advertised during
// initialize.
func WithClientInfo(info mcp.Implementation) Option {
	return func(c *Client) { c.clientInfo = info }
}

// WithProtocolVersion overrides the protocol version requested at initialize.
func WithProtocolVersion(v string) Option {
	return func(c *Client) {
		if v != "" {
			c.protocol = v
		}
	}
}

// WithRequestTimeout sets the per-call timeout. Zero disables it.
func WithRequestTimeout(d time.Duration) Option {
	return func(c *Client) { c.reqTimeout = d }
}

// WithLogger sets the structured logger.
func WithLogger(l *log.Logger) Option {
	return func(c *Client) {
		if l != nil {
			c.log = l
		}
	}
}

// WithNotificationHandler registers a handler for server notifications.
func WithNotificationHandler(h NotificationHandler) Option {
	return func(c *Client) { c.onNotify = h }
}
