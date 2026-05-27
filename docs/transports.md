# Transports

A transport moves raw JSON-RPC frames. mcpc ships three.

## Command (stdio)

`transport.StartCommand` launches an MCP server as a child process and speaks
newline-delimited JSON over its stdin/stdout. A background goroutine reads frames
so `Receive` honours context cancellation, and `Close` terminates the child. This
is the usual way to reach a stdio server such as mcpkit.

## HTTP

`transport.NewHTTP` targets a gateway's `/rpc` endpoint. Because HTTP has no
independent server-to-client stream, each `Send` performs the POST and queues the
JSON reply; the matching `Receive` dequeues it. A notification returns `204 No
Content` and enqueues nothing.

## Pipe

`transport.Pipe` connects two in-memory endpoints, used by tests and the client's
own test suite to run a client against a mock server without any process or
socket.

## Writing your own

Implement the three-method `transport.Transport` interface (`Receive`, `Send`,
`Close`) and hand it to `client.New`.
