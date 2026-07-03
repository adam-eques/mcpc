# Roadmap

Ideas under consideration. Nothing here is a commitment.

## Near term

- Batch requests once the negotiated protocol version supports them.
- Resource subscriptions (`resources/subscribe`) with change notifications.
- A `--json` output mode for `call` and `tools` for scripting.

## Later

- Sampling support (`sampling/createMessage`) so a server can ask the client's
  model to generate text.
- A reconnect strategy for the HTTP transport with resumable sessions.
- A recording transport that captures a session to a file for replay in tests.

## Non-goals

- Bundling a specific LLM. mcpc is the transport and orchestration layer; the
  model lives above it.
- Third-party dependencies. The client stays standard-library-only.
