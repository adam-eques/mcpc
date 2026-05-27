# Using mcpc as a library

The `client` package is a reusable MCP client you can embed in your own Go
programs.

## Connect and call

```go
ctx := context.Background()
c, info, err := client.DialCommand(ctx, "mcpkit", nil)
if err != nil {
    log.Fatal(err)
}
defer c.Close()
fmt.Println("connected to", info.ServerInfo.Name)

tools, _ := c.ListTools(ctx)
res, _ := c.CallTool(ctx, "calculate", map[string]string{"expression": "2+2"})
```

## Transports

- `client.DialCommand(ctx, name, args, opts...)` — spawn a stdio server.
- `client.DialHTTP(ctx, endpoint, opts...)` — talk to a gateway.
- `client.New(transport, opts...)` — bring your own transport, then `Start` and
  `Initialize` yourself.

## Options

`WithClientInfo`, `WithProtocolVersion`, `WithRequestTimeout`, `WithLogger`,
`WithMetrics`, `WithMaxRetries`, and `WithNotificationHandler`.

## Reliability helpers

- `CallToolRetry` retries transient failures with exponential backoff.
- `ProgressTracker` routes `notifications/progress` to per-token callbacks.
- `Pool` maintains several server processes for CPU-bound workloads.

## Concurrency

A single `Client` is safe for concurrent use; issue calls from as many
goroutines as you like and the read loop correlates each response to its caller.
