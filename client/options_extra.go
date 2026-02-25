package client

import "github.com/adam-eques/mcpc/internal/metrics"

// WithMetrics attaches a metrics registry that records request counts and
// latencies.
func WithMetrics(m *metrics.Registry) Option {
	return func(c *Client) {
		if m != nil {
			c.metrics = m
		}
	}
}

// WithMaxRetries sets how many additional attempts CallToolRetry makes when a
// tool reports a transient (non-protocol) failure. Zero disables retrying.
func WithMaxRetries(n int) Option {
	return func(c *Client) {
		if n >= 0 {
			c.maxRetries = n
		}
	}
}
