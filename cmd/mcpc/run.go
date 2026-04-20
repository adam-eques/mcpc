package main

import (
	"context"
	"fmt"
	"os"

	"github.com/adam-eques/mcpc/agent"
	"github.com/adam-eques/mcpc/internal/config"
	"github.com/adam-eques/mcpc/internal/log"
)

func cmdRun(ctx context.Context, cfg config.Config, logger *log.Logger, rest []string) error {
	if len(rest) != 1 {
		return fmt.Errorf("usage: mcpc run <plan.json>")
	}
	data, err := os.ReadFile(rest[0])
	if err != nil {
		return err
	}
	plan, err := agent.ParsePlan(data)
	if err != nil {
		return err
	}

	c, _, err := connect(ctx, cfg, logger)
	if err != nil {
		return err
	}
	defer c.Close()

	report, err := agent.Run(ctx, c, plan)
	if err != nil {
		return err
	}
	for i, res := range report.Results {
		label := res.Step.Tool
		if res.Step.Name != "" {
			label = res.Step.Name
		}
		switch {
		case res.Err != nil:
			fmt.Printf("%d. %s: ERROR %v\n", i+1, label, res.Err)
		case res.IsErr:
			fmt.Printf("%d. %s: tool-error\n   %s\n", i+1, label, res.Output)
		default:
			fmt.Printf("%d. %s:\n   %s\n", i+1, label, res.Output)
		}
	}
	if report.Failed() {
		return fmt.Errorf("plan finished with errors")
	}
	return nil
}
