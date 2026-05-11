# Documentation index

- [Architecture](architecture.md) — layers, the read loop and cancellation.
- [Protocol support](protocol.md) — handshake, methods and notifications.
- [Command-line usage](cli.md) — the `mcpc` commands and flags.
- [Using mcpc as a library](library.md) — embed the client in your Go program.
- [Transports](transports.md) — stdio, HTTP and the in-memory pipe.
- [Scripted plans](agent.md) — the `agent` runner and output chaining.
- [Configuration](configuration.md) — defaults, files and environment.

New to the codebase? Read the architecture overview, then the `client` package's
read loop in `client/client.go`.
