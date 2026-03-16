package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/adam-eques/mcpc/mcp"
)

type stubCaller struct {
	results map[string]*mcp.CallToolResult
	errs    map[string]error
	calls   []string
}

func (s *stubCaller) CallTool(_ context.Context, name string, _ any) (*mcp.CallToolResult, error) {
	s.calls = append(s.calls, name)
	if err, ok := s.errs[name]; ok {
		return nil, err
	}
	if res, ok := s.results[name]; ok {
		return res, nil
	}
	return mcp.TextResult("ok"), nil
}

func TestParsePlanObjectAndArray(t *testing.T) {
	obj, err := ParsePlan([]byte(`{"steps":[{"tool":"a"}]}`))
	if err != nil || len(obj.Steps) != 1 {
		t.Fatalf("object plan: %v %+v", err, obj)
	}
	arr, err := ParsePlan([]byte(`[{"tool":"a"},{"tool":"b"}]`))
	if err != nil || len(arr.Steps) != 2 {
		t.Fatalf("array plan: %v %+v", err, arr)
	}
}

func TestValidate(t *testing.T) {
	if err := (&Plan{Steps: []Step{{Tool: ""}}}).Validate(); err == nil {
		t.Fatal("expected validation error for missing tool")
	}
}

func TestRunStopsOnError(t *testing.T) {
	caller := &stubCaller{errs: map[string]error{"b": errors.New("boom")}}
	plan := &Plan{Steps: []Step{{Tool: "a"}, {Tool: "b"}, {Tool: "c"}}}
	report, err := Run(context.Background(), caller, plan)
	if err != nil {
		t.Fatal(err)
	}
	if len(caller.calls) != 2 {
		t.Fatalf("expected to stop after b, called %v", caller.calls)
	}
	if !report.Failed() {
		t.Fatal("report should be marked failed")
	}
}

func TestRunContinueOnError(t *testing.T) {
	caller := &stubCaller{errs: map[string]error{"b": errors.New("boom")}}
	plan := &Plan{Steps: []Step{{Tool: "a"}, {Tool: "b", ContinueOnError: true}, {Tool: "c"}}}
	if _, err := Run(context.Background(), caller, plan); err != nil {
		t.Fatal(err)
	}
	if len(caller.calls) != 3 {
		t.Fatalf("expected all steps to run, called %v", caller.calls)
	}
}

func TestRenderContent(t *testing.T) {
	got := renderContent([]mcp.Content{mcp.Text("hello"), mcp.Text("world")})
	if got != "hello\nworld" {
		t.Fatalf("render=%q", got)
	}
}
