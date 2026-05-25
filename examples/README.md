# Examples

These examples assume the paired server, `mcpkit`, is available — either on your
`PATH` or runnable with `go run github.com/adam-eques/mcpkit/cmd/mcpkit`.

## Library

Embed the client in a Go program:

```bash
go run ./examples/library mcpkit
```

## Scripted plans

Run a plan of tool calls (index → search → hash):

```bash
mcpc -server mcpkit run examples/plan.json
```

Output chaining — the id created in step 1 flows into steps 2 and 3 via `${1}`:

```bash
mcpc -server mcpkit run examples/chaining-plan.json
```

## Config file

`mcpc.json` launches the server with `go run` and enables debug logging:

```bash
mcpc -config examples/mcpc.json tools
```
