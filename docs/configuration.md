# Configuration

mcpc reads defaults, overlays an optional JSON file (`-config`), then applies
`MCPC_*` environment overrides. Command-line flags win over all of them.

```
defaults < config file < environment < flags
```

## Fields

| Field | Meaning |
| --- | --- |
| `transport` | `stdio` or `http` |
| `server.command`, `server.args` | server process for stdio |
| `http.endpoint` | gateway URL for http |
| `timeoutSeconds` | per-request timeout |
| `log.level`, `log.format` | logging |
| `client.name`, `client.version` | identity sent at initialize |

## Environment variables

| Variable | Effect |
| --- | --- |
| `MCPC_TRANSPORT` | `stdio` or `http` |
| `MCPC_SERVER` | server command line, e.g. `go run ./cmd/mcpkit` |
| `MCPC_HTTP_ENDPOINT` | gateway URL |
| `MCPC_TIMEOUT` | timeout in seconds |
| `MCPC_LOG_LEVEL` | log level |

See `config.example.json` for a complete file.
