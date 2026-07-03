# mcpc

A **Model Context Protocol (MCP) client — library and CLI — written in Go with
zero third-party dependencies.** It connects to any MCP server over stdio
(launching it as a child process) or HTTP, and is the paired client for
[`mcpkit`](https://github.com/adam-eques/mcpkit).

[![CI](https://github.com/adam-eques/mcpc/actions/workflows/ci.yml/badge.svg)](https://github.com/adam-eques/mcpc/actions/workflows/ci.yml)
![Go](https://img.shields.io/badge/go-1.26-00ADD8)
![Dependencies](https://img.shields.io/badge/dependencies-none-brightgreen)
![License](https://img.shields.io/badge/license-MIT-blue)

## Why

There is an official MCP Go SDK, and for most production clients it's the right
choice. mcpc is a deliberate alternative — a from-scratch, zero-dependency
implementation built to make the interesting part legible: the concurrency.

One connection, many in-flight requests, responses arriving out of order, plus
server-initiated notifications. mcpc solves this with a single background read
loop that demultiplexes responses to their waiting callers by request id — so you
can fan out calls from as many goroutines as you like, with per-call timeouts and
cancellation that propagates to the server as `notifications/cancelled`. It's all
on the standard library, so the read loop and the cancellation plumbing are here
to read, not hidden behind a dependency.

## Highlights

- **Full client protocol** — MCP `2025-06-18`: initialize handshake, tools,
  resources, prompts, logging and cancellation.
- **Three transports** — a **subprocess** transport that launches a stdio server,
  an **HTTP** transport for gateways, and an in-memory **pipe** for tests.
- **Concurrent by design** — a demultiplexing read loop, a `pending` map keyed by
  id, and first-class context timeouts.
- **Reliability** — `CallToolRetry` with exponential backoff, a `ProgressTracker`
  for streaming progress, and a `Pool` of server processes for CPU-bound work.
- **A real CLI** — `ping`, `tools`, `describe`, `call`, `run` and an interactive
  `repl`.
- **A scripted agent runner** — execute a JSON plan of tool calls with output
  chaining (`${1}` feeds one step's result into the next).
- **Tested** — table-driven unit tests against a mock server over the in-memory
  pipe, benchmarks, and an optional live integration test against `mcpkit`.

## Quick start

```bash
# Build the CLI
make build

# List a server's tools (launches mcpkit as a child process)
mcpc -server mcpkit tools

# Call a tool
mcpc -server mcpkit call calculate expression="2 ^ 10 + sqrt(81)"

# Run a scripted plan
mcpc -server mcpkit run examples/plan.json

# Interactive session
mcpc -server mcpkit repl
```

Against a running HTTP gateway instead:

```bash
mcpc -http http://localhost:8080/rpc ping
```

## Library

```go
ctx := context.Background()
c, info, err := client.DialCommand(ctx, "mcpkit", nil)
if err != nil {
    log.Fatal(err)
}
defer c.Close()

tools, _ := c.ListTools(ctx)
res, _ := c.CallTool(ctx, "calculate", map[string]string{"expression": "6*7"})
fmt.Println(res.Content[0].(mcp.TextContent).Text) // 42
```

## Project layout

```
cmd/mcpc/       the command-line client
client/         the reusable client: read loop, retry, progress, pool
agent/          scripted plan parser and runner with output chaining
transport/      subprocess, HTTP and in-memory pipe transports
mcp/            MCP protocol types
jsonrpc/        JSON-RPC 2.0 core
internal/       config, logging, metrics, version
docs/           architecture, protocol, cli, library, transports, agent
examples/       a library example and sample plans
```

## Paired with mcpkit

mcpc is developed alongside the [`mcpkit`](https://github.com/adam-eques/mcpkit)
server. Build mcpkit and point the live integration test at it:

```bash
go build -o /tmp/mcpkit github.com/adam-eques/mcpkit/cmd/mcpkit
MCPKIT_BIN=/tmp/mcpkit go test ./client -run TestLiveServer
```

## Documentation

- [Architecture](docs/architecture.md)
- [Protocol support](docs/protocol.md)
- [Command-line usage](docs/cli.md)
- [Using mcpc as a library](docs/library.md)
- [Transports](docs/transports.md)
- [Scripted plans](docs/agent.md)
- [Configuration](docs/configuration.md)

## License

MIT — see [LICENSE](LICENSE).
