# Scripted plans

The `agent` package runs a plan of tool calls against a server. It is the engine
behind `mcpc run` and the execution layer an LLM agent would drive.

## Plan format

A plan is JSON — either an object with a `steps` array or a bare array:

```json
{
  "steps": [
    {"name": "index", "tool": "rag_index", "arguments": {"documents": [{"id": "a", "text": "..."}]}},
    {"name": "search", "tool": "rag_search", "arguments": {"query": "...", "k": 2}}
  ]
}
```

Each step has a `tool`, optional `arguments`, an optional `name` for output, and
`continueOnError` to keep going past a failure.

## Output chaining

String arguments may reference an earlier step's output with `${N}` (1-based):

```json
{"steps": [
  {"tool": "uuid"},
  {"tool": "kv_set", "arguments": {"key": "id", "value": "${1}"}}
]}
```

Before step 2 runs, `${1}` is replaced with the text output of step 1.

## Running

```bash
mcpc -server mcpkit run plan.json
```

Or from Go:

```go
plan, _ := agent.ParsePlan(data)
report, _ := agent.Run(ctx, client, plan)
```
