package agent

import (
	"context"
	"strings"

	"github.com/adam-eques/mcpc/mcp"
)

// ToolCaller is the subset of the client the runner needs. *client.Client
// satisfies it.
type ToolCaller interface {
	CallTool(ctx context.Context, name string, args any) (*mcp.CallToolResult, error)
}

// Run executes every step of the plan in order, substituting ${N} references in
// each step's arguments with earlier steps' outputs. Execution stops at the
// first failing step unless that step sets ContinueOnError.
func Run(ctx context.Context, caller ToolCaller, plan *Plan) (Report, error) {
	if err := plan.Validate(); err != nil {
		return Report{}, err
	}
	var report Report
	var outputs []string
	for _, step := range plan.Steps {
		step.Arguments = substitute(step.Arguments, outputs)
		res := runStep(ctx, caller, step)
		report.Results = append(report.Results, res)
		outputs = append(outputs, res.Output)
		if (res.Err != nil || res.IsErr) && !step.ContinueOnError {
			break
		}
	}
	return report, nil
}

func runStep(ctx context.Context, caller ToolCaller, step Step) Result {
	out, err := caller.CallTool(ctx, step.Tool, step.Arguments)
	if err != nil {
		return Result{Step: step, Err: err}
	}
	return Result{Step: step, Output: renderContent(out.Content), IsErr: out.IsError}
}

func renderContent(content []mcp.Content) string {
	var b strings.Builder
	for i, c := range content {
		if i > 0 {
			b.WriteString("\n")
		}
		switch v := c.(type) {
		case mcp.TextContent:
			b.WriteString(v.Text)
		case mcp.ImageContent:
			b.WriteString("[image " + v.MimeType + "]")
		default:
			b.WriteString("[content]")
		}
	}
	return b.String()
}
