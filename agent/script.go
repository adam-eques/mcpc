// Package agent runs scripted plans of tool calls against an MCP server. It is
// the reusable core behind the `mcpc run` command: parse a plan of steps, invoke
// each tool in order, and collect the results. It is deliberately model-free —
// an LLM agent would produce the same plan shape, so this is the execution layer
// such an agent would drive.
package agent

import (
	"encoding/json"
	"fmt"
)

// Step is a single tool invocation in a plan.
type Step struct {
	// Name is an optional label used in output.
	Name string `json:"name,omitempty"`
	// Tool is the tool to call.
	Tool string `json:"tool"`
	// Arguments are passed verbatim to the tool.
	Arguments map[string]any `json:"arguments,omitempty"`
	// ContinueOnError keeps the plan running if this step fails.
	ContinueOnError bool `json:"continueOnError,omitempty"`
}

// Plan is an ordered list of steps.
type Plan struct {
	Steps []Step `json:"steps"`
}

// ParsePlan decodes a plan from JSON. It accepts either an object with a "steps"
// array or a bare array of steps.
func ParsePlan(data []byte) (*Plan, error) {
	trimmed := trimSpace(data)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("agent: empty plan")
	}
	if trimmed[0] == '[' {
		var steps []Step
		if err := json.Unmarshal(data, &steps); err != nil {
			return nil, err
		}
		return &Plan{Steps: steps}, nil
	}
	var p Plan
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Validate checks that every step names a tool.
func (p *Plan) Validate() error {
	for i, s := range p.Steps {
		if s.Tool == "" {
			return fmt.Errorf("agent: step %d has no tool", i)
		}
	}
	return nil
}

func trimSpace(b []byte) []byte {
	i := 0
	for i < len(b) && (b[i] == ' ' || b[i] == '\t' || b[i] == '\n' || b[i] == '\r') {
		i++
	}
	return b[i:]
}
