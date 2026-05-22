# Talking to a gateway over HTTP

Start the paired server's HTTP gateway, then point mcpc at it.

```bash
# Terminal 1 — run the gateway (from the mcpkit repo)
go run ./cmd/mcpkit-gateway -addr :8080

# Terminal 2 — drive it with mcpc
mcpc -http http://localhost:8080/rpc ping
mcpc -http http://localhost:8080/rpc tools
mcpc -http http://localhost:8080/rpc call calculate expression="7*6"
```

The HTTP transport turns each request into a `POST /rpc` and reads the JSON-RPC
response from the body. Notifications receive `204 No Content`. Operational
endpoints `/healthz` and `/metrics` are served by the gateway itself.
