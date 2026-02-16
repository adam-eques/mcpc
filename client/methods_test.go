package client

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/adam-eques/mcpc/jsonrpc"
	"github.com/adam-eques/mcpc/mcp"
)

func TestResourcesAndPrompts(t *testing.T) {
	fs, end := newFakeServer()
	fs.handlers[mcp.MethodResourcesList] = func(json.RawMessage) (any, *jsonrpc.Error) {
		return mcp.ListResourcesResult{Resources: []mcp.Resource{{URI: "mem://a", Name: "a"}}}, nil
	}
	fs.handlers[mcp.MethodResourcesRead] = func(json.RawMessage) (any, *jsonrpc.Error) {
		return mcp.ReadResourceResult{Contents: []mcp.ResourceContents{{URI: "mem://a", Text: "body"}}}, nil
	}
	fs.handlers[mcp.MethodPromptsList] = func(json.RawMessage) (any, *jsonrpc.Error) {
		return mcp.ListPromptsResult{Prompts: []mcp.Prompt{{Name: "greet"}}}, nil
	}
	fs.handlers[mcp.MethodPromptsGet] = func(json.RawMessage) (any, *jsonrpc.Error) {
		return mcp.GetPromptResult{Messages: []mcp.PromptMessage{{Role: mcp.RoleUser, Content: mcp.Text("hi")}}}, nil
	}
	c, cancel := startClient(t, fs, end)
	defer cancel()
	defer c.Close()
	c.Initialize(context.Background())

	resources, err := c.ListResources(context.Background())
	if err != nil || len(resources) != 1 {
		t.Fatalf("resources: %v %+v", err, resources)
	}
	contents, err := c.ReadResource(context.Background(), "mem://a")
	if err != nil || contents[0].Text != "body" {
		t.Fatalf("read: %v %+v", err, contents)
	}
	prompts, err := c.ListPrompts(context.Background())
	if err != nil || prompts[0].Name != "greet" {
		t.Fatalf("prompts: %v %+v", err, prompts)
	}
	got, err := c.GetPrompt(context.Background(), "greet", map[string]string{"who": "ada"})
	if err != nil || len(got.Messages) != 1 {
		t.Fatalf("get prompt: %v %+v", err, got)
	}
}
