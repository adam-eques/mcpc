# Contributing

Thanks for your interest in mcpc.

## Development

```bash
make test    # run the suite
make race    # run with the race detector
make vet     # static checks
make bench   # benchmarks
```

The module has **no third-party dependencies** — please keep it
standard-library-only.

To run the live integration test against the paired server, build mcpkit and
point `MCPKIT_BIN` at it:

```bash
go build -o /tmp/mcpkit github.com/adam-eques/mcpkit/cmd/mcpkit
MCPKIT_BIN=/tmp/mcpkit go test ./client -run TestLiveServer
```

## Branching

Work happens on `dev` and is merged into `main` at milestone boundaries. Keep
commits small with imperative messages ("add", "fix", "refactor").

## Checklist before opening a PR

- [ ] `make vet test` passes
- [ ] New behaviour is covered by tests
- [ ] Exported identifiers are documented
- [ ] No new dependencies
