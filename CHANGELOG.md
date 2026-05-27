# Changelog

All notable changes to mcpc are documented here.

## 0.0.1 — Scaffolding

- Initialize the Go module, license and build tooling.

## 0.0.2 — JSON-RPC 2.0

- Implement the JSON-RPC 2.0 message and error types the client sends and decodes.

## 0.0.3 — MCP types

- Add the initialize handshake, capabilities, content blocks and primitives.

## 0.0.4 — Transports

- Add subprocess, HTTP and in-memory pipe transports behind one interface.

## 0.0.5 — Observability

- Add a structured logger, a metrics registry and build stamping.

## 0.1.0-rc1 — Client core

- Add a background read loop that correlates responses by id.
- Support concurrent calls, per-call timeouts and cancellation.

## 0.1.0-rc2 — Reliability

- Add dial helpers, retry with backoff, progress tracking and a connection pool.
- Record request counts and latencies in a metrics registry.

## 0.1.0-rc3 — Agent runner

- Execute a JSON plan of tool calls with ${N} output chaining between steps.

## 0.1.0-rc4 — Configuration

- Add layered configuration: defaults, JSON file and MCPC_* environment.

## 0.1.0-rc5 — Command-line client

- Add the mcpc CLI: ping, tools, describe, call, run and an interactive repl.

## 0.1.0-rc6 — Documentation

- Add architecture, protocol, CLI, library, transport and agent guides.
- Add a library example and sample plans.
