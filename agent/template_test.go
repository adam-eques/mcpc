package agent

import (
	"context"
	"testing"

	"github.com/adam-eques/mcpc/mcp"
)

func TestExpand(t *testing.T) {
	outputs := []string{"42", "hello"}
	if got := expand("value is ${1}", outputs); got != "value is 42" {
		t.Fatalf("got %q", got)
	}
	if got := expand("${2} ${1}", outputs); got != "hello 42" {
		t.Fatalf("got %q", got)
	}
	if got := expand("${9} untouched", outputs); got != "${9} untouched" {
		t.Fatalf("out-of-range should be left alone, got %q", got)
	}
}

func TestRunChainsOutputs(t *testing.T) {
	caller := &recordingCaller{}
	plan := &Plan{Steps: []Step{
		{Tool: "first"},
		{Tool: "second", Arguments: map[string]any{"text": "prev=${1}"}},
	}}
	if _, err := Run(context.Background(), caller, plan); err != nil {
		t.Fatal(err)
	}
	if caller.lastArgs["text"] != "prev=out-first" {
		t.Fatalf("chaining failed: %v", caller.lastArgs)
	}
}

type recordingCaller struct {
	lastArgs map[string]any
}

func (r *recordingCaller) CallTool(_ context.Context, name string, args any) (*mcp.CallToolResult, error) {
	if m, ok := args.(map[string]any); ok {
		r.lastArgs = m
	}
	return mcp.TextResult("out-" + name), nil
}
