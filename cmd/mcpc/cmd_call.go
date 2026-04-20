package main

import (
	"context"
	"fmt"

	"github.com/adam-eques/mcpc/internal/config"
	"github.com/adam-eques/mcpc/internal/log"
)

func cmdCall(ctx context.Context, cfg config.Config, logger *log.Logger, rest []string) error {
	if len(rest) == 0 {
		return fmt.Errorf("usage: mcpc call <tool> [key=value ...]")
	}
	args, err := parseArgs(rest[1:])
	if err != nil {
		return err
	}
	c, _, err := connect(ctx, cfg, logger)
	if err != nil {
		return err
	}
	defer c.Close()

	res, err := c.CallTool(ctx, rest[0], args)
	if err != nil {
		return err
	}
	printResult(res)
	if res.IsError {
		return fmt.Errorf("tool reported an error")
	}
	return nil
}
