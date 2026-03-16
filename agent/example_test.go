package agent_test

import (
	"context"
	"fmt"

	"github.com/adam-eques/mcpc/agent"
	"github.com/adam-eques/mcpc/mcp"
)

// stubCaller returns a fixed result so the example is deterministic.
type stubCaller struct{}

func (stubCaller) CallTool(_ context.Context, name string, _ any) (*mcp.CallToolResult, error) {
	return mcp.TextResult("result of " + name), nil
}

// ExampleRun shows how to execute a two-step plan.
func ExampleRun() {
	plan, _ := agent.ParsePlan([]byte(`{"steps":[{"tool":"uuid"},{"tool":"hash"}]}`))
	report, _ := agent.Run(context.Background(), stubCaller{}, plan)
	for _, r := range report.Results {
		fmt.Println(r.Output)
	}
	// Output:
	// result of uuid
	// result of hash
}
