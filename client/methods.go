package client

import (
	"context"
	"encoding/json"

	"github.com/adam-eques/mcpc/mcp"
)

// Initialize performs the MCP handshake: it sends initialize, records the
// server's identity and capabilities, and sends the initialized notification.
func (c *Client) Initialize(ctx context.Context) (*mcp.InitializeResult, error) {
	params := mcp.InitializeParams{
		ProtocolVersion: c.protocol,
		ClientInfo:      c.clientInfo,
	}
	var res mcp.InitializeResult
	if err := c.call(ctx, mcp.MethodInitialize, params, &res); err != nil {
		return nil, err
	}
	c.sessMu.Lock()
	c.initialized = true
	c.serverInfo = res.ServerInfo
	c.serverCaps = res.Capabilities
	c.sessMu.Unlock()

	if err := c.notify(ctx, mcp.NotificationInitialized, nil); err != nil {
		return nil, err
	}
	c.log.Info("session initialized",
		"server", res.ServerInfo.Name,
		"serverVersion", res.ServerInfo.Version,
		"protocol", res.ProtocolVersion)
	return &res, nil
}

// Ping checks that the server is responsive.
func (c *Client) Ping(ctx context.Context) error {
	return c.call(ctx, mcp.MethodPing, nil, nil)
}

// ListTools returns the tools advertised by the server.
func (c *Client) ListTools(ctx context.Context) ([]mcp.Tool, error) {
	var res mcp.ListToolsResult
	if err := c.call(ctx, mcp.MethodToolsList, nil, &res); err != nil {
		return nil, err
	}
	return res.Tools, nil
}

// CallTool invokes a tool by name. args may be a map, a struct or nil.
func (c *Client) CallTool(ctx context.Context, name string, args any) (*mcp.CallToolResult, error) {
	var rawArgs json.RawMessage
	if args != nil {
		b, err := json.Marshal(args)
		if err != nil {
			return nil, err
		}
		rawArgs = b
	}
	params := mcp.CallToolParams{Name: name, Arguments: rawArgs}
	var res mcp.CallToolResult
	if err := c.call(ctx, mcp.MethodToolsCall, params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// ListResources returns the resources advertised by the server.
func (c *Client) ListResources(ctx context.Context) ([]mcp.Resource, error) {
	var res mcp.ListResourcesResult
	if err := c.call(ctx, mcp.MethodResourcesList, nil, &res); err != nil {
		return nil, err
	}
	return res.Resources, nil
}

// ReadResource reads a resource by URI.
func (c *Client) ReadResource(ctx context.Context, uri string) ([]mcp.ResourceContents, error) {
	var res mcp.ReadResourceResult
	if err := c.call(ctx, mcp.MethodResourcesRead, mcp.ReadResourceParams{URI: uri}, &res); err != nil {
		return nil, err
	}
	return res.Contents, nil
}

// ListPrompts returns the prompts advertised by the server.
func (c *Client) ListPrompts(ctx context.Context) ([]mcp.Prompt, error) {
	var res mcp.ListPromptsResult
	if err := c.call(ctx, mcp.MethodPromptsList, nil, &res); err != nil {
		return nil, err
	}
	return res.Prompts, nil
}

// GetPrompt renders a prompt by name with the given arguments.
func (c *Client) GetPrompt(ctx context.Context, name string, args map[string]string) (*mcp.GetPromptResult, error) {
	params := mcp.GetPromptParams{Name: name, Arguments: args}
	var res mcp.GetPromptResult
	if err := c.call(ctx, mcp.MethodPromptsGet, params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// SetLogLevel asks the server to change its logging verbosity.
func (c *Client) SetLogLevel(ctx context.Context, level mcp.LoggingLevel) error {
	return c.call(ctx, mcp.MethodLoggingSetLevel, mcp.SetLevelParams{Level: level}, nil)
}
