# Security Policy

## Reporting a vulnerability

Please report security issues privately via Telegram
[@dracoeques](https://t.me/dracoeques) rather than opening a public issue.

## Scope

mcpc launches and talks to MCP servers. Treat any server you connect to as you
would any executable: the `stdio` transport runs the configured command as a
child process. Only point mcpc at servers you trust.

## Supported versions

The `main` branch receives security fixes; tagged releases are patched on a
best-effort basis.
