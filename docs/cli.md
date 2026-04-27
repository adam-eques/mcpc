# Command-line usage

```
mcpc [flags] <command> [args]
```

## Global flags

| Flag | Meaning |
| --- | --- |
| `-server "cmd args"` | Launch this command as a stdio server |
| `-http URL` | Use the HTTP transport against a gateway |
| `-config path` | Load a JSON config file |
| `-timeout N` | Per-request timeout in seconds |
| `-log-level` | debug, info, warn, error |

## Commands

| Command | Description |
| --- | --- |
| `ping` | Check the server is responsive |
| `tools` | List the server's tools |
| `describe <tool>` | Print a tool's JSON input schema |
| `call <tool> [k=v ...]` | Invoke a tool |
| `run <plan.json>` | Execute a scripted plan |
| `repl` | Interactive session |
| `version` | Print the client version |

## Arguments

`call` accepts arguments as `key=value` (string) or `key:=json` (raw JSON):

```bash
mcpc -server mcpkit call calculate expression="2 ^ 10"
mcpc -server mcpkit call rag_search query="golang" k:=3
```

## Examples

```bash
# Against a locally built server
mcpc -server ./mcpkit tools

# Against a running gateway
mcpc -http http://localhost:8080/rpc ping
```
