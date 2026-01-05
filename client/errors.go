// Package client implements the client half of the Model Context Protocol. A
// Client runs a background read loop that correlates responses to in-flight
// requests by identifier, dispatches server-initiated notifications, and
// forwards context cancellation to the server as notifications/cancelled.
package client

import "errors"

// ErrClosed is returned when a call is made on a closed client or the transport
// has shut down.
var ErrClosed = errors.New("client: connection closed")

// ErrNotInitialized is returned when a call requiring an initialized session is
// made before Initialize succeeds.
var ErrNotInitialized = errors.New("client: session not initialized")

// ErrAlreadyStarted is returned by Start when the read loop is already running.
var ErrAlreadyStarted = errors.New("client: already started")
