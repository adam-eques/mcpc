package client

import (
	"context"
	"encoding/json"

	"github.com/adam-eques/mcpc/jsonrpc"
	"github.com/adam-eques/mcpc/mcp"
	"github.com/adam-eques/mcpc/transport"
)

// fakeServer plays the server role on one end of a pipe so the client can be
// tested without a real MCP server. handlers maps a method to a function that
// returns a result (or an error); a nil handler for a method returns method not
// found. Methods can also delay by reading from block before replying.
type fakeServer struct {
	t        transport.Transport
	handlers map[string]func(params json.RawMessage) (any, *jsonrpc.Error)
	block    chan struct{}
	blockOn  string
}

func newFakeServer() (*fakeServer, transport.Transport) {
	clientEnd, serverEnd := transport.Pipe()
	fs := &fakeServer{
		t:        serverEnd,
		handlers: map[string]func(json.RawMessage) (any, *jsonrpc.Error){},
		block:    make(chan struct{}),
	}
	fs.handlers[mcp.MethodInitialize] = func(json.RawMessage) (any, *jsonrpc.Error) {
		return mcp.InitializeResult{
			ProtocolVersion: mcp.ProtocolVersion,
			Capabilities:    mcp.ServerCapabilities{Tools: &mcp.ToolsCapability{}},
			ServerInfo:      mcp.Implementation{Name: "fake", Version: "9.9"},
		}, nil
	}
	fs.handlers[mcp.MethodPing] = func(json.RawMessage) (any, *jsonrpc.Error) {
		return mcp.PingResult{}, nil
	}
	return fs, clientEnd
}

func (fs *fakeServer) run(ctx context.Context) {
	for {
		frame, err := fs.t.Receive(ctx)
		if err != nil {
			return
		}
		var env envelope
		if json.Unmarshal(frame, &env) != nil {
			continue
		}
		if env.Method == "" || env.ID == nil {
			continue // notification: nothing to answer
		}
		if fs.blockOn != "" && env.Method == fs.blockOn {
			select {
			case <-fs.block:
			case <-ctx.Done():
				return
			}
		}
		h, ok := fs.handlers[env.Method]
		var resp *jsonrpc.Response
		if !ok {
			resp = jsonrpc.NewErrorResponse(env.ID, jsonrpc.MethodNotFound(env.Method))
		} else {
			result, rpcErr := h(env.Params)
			if rpcErr != nil {
				resp = jsonrpc.NewErrorResponse(env.ID, rpcErr)
			} else {
				raw, _ := json.Marshal(result)
				resp = jsonrpc.NewResponse(env.ID, raw)
			}
		}
		out, _ := json.Marshal(resp)
		if err := fs.t.Send(ctx, out); err != nil {
			return
		}
	}
}
