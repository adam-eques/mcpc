# FAQ

**Does mcpc only work with mcpkit?**
No. It speaks standard MCP over stdio or HTTP, so it works with any compliant
server. mcpkit is simply the paired reference server used in the examples and the
optional live integration test.

**Why no dependencies?**
The client is small enough to build on the standard library, which keeps it
trivial to audit and to vendor. The protocol types are the same ones a server
uses, so there is nothing exotic to pull in.

**Can I use it as a library?**
Yes — see [library.md](library.md). The CLI is a thin wrapper over the `client`
package.

**How are concurrent calls handled?**
A single connection multiplexes them. The background read loop matches each
response to the goroutine that is waiting on that request id, so you can fan out
calls freely.

**How do I run the live integration test?**
Build mcpkit and set `MCPKIT_BIN` to it before running `go test ./client -run
TestLiveServer`.
