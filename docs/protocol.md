# Protocol support

mcpc implements the client side of the Model Context Protocol `2025-06-18` and
requests negotiation down to older revisions a server may prefer.

## Handshake

`Client.Initialize` performs the full handshake:

1. Send `initialize` with the client's protocol version and identity.
2. Record the server's identity and capabilities from the result.
3. Send the `notifications/initialized` notification.

After that, the tools, resources and prompts methods are available.

## Methods used

| Method | Client API |
| --- | --- |
| `initialize` | `Initialize` |
| `ping` | `Ping` |
| `tools/list` | `ListTools` |
| `tools/call` | `CallTool`, `CallToolRetry` |
| `resources/list`, `resources/read` | `ListResources`, `ReadResource` |
| `prompts/list`, `prompts/get` | `ListPrompts`, `GetPrompt` |
| `logging/setLevel` | `SetLogLevel` |

## Notifications

Outbound: `notifications/initialized`, `notifications/cancelled`. Inbound
notifications (progress, log messages) are delivered to a `NotificationHandler`;
the `ProgressTracker` helper routes `notifications/progress` to per-token
callbacks.
