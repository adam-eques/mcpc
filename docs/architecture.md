# Architecture

mcpc is a layered, dependency-free MCP client. Each layer depends only on the
ones beneath it.

```
            ┌───────────────────────────────────────────────┐
 cmd/mcpc   │  ping · tools · describe · call · run · repl   │
            └───────────────────────┬───────────────────────┘
                                    │
 agent/     plan parser + runner ───┤ drives tool calls, chains outputs
                                    │
            ┌───────────────────────▼───────────────────────┐
 client/    │  Client: read loop, pending map, retry, pool   │
            └───────────────────────┬───────────────────────┘
                                    │
 transport/  Command (subprocess) · HTTP · Pipe
                                    │
 mcp/        protocol types        jsonrpc/  JSON-RPC 2.0 core
```

## The read loop

A `Client` runs a single background goroutine that receives frames and
demultiplexes them:

- A **response** (has an `id`, no `method`) is delivered to the waiting caller
  through a per-request channel kept in a `pending` map keyed by id.
- A **notification** (has a `method`, no `id`) is passed to the registered
  `NotificationHandler` — used for progress and log messages.
- A server-initiated **request** (has both) is answered inline; `ping` is
  handled, anything else returns method-not-found so the server is never left
  waiting.

This lets many goroutines issue calls concurrently on one connection. Each call
allocates an id, registers a channel, sends the frame, and then selects on its
channel, its context, and the loop's shutdown — so timeouts and cancellation are
first-class.

## Cancellation

When a call's context is cancelled or times out, the client sends
`notifications/cancelled` for that request id so the server can stop the work,
then returns the context error. This mirrors the server's per-request
cancellation support.

## Design choices

- **No third-party dependencies** — standard library only, so it builds anywhere.
- **Transport-agnostic core** — the same `Client` runs over a spawned process, an
  HTTP gateway, or an in-memory pipe.
- **Stdout discipline** — logs go to stderr; stdout is reserved for command
  output.
