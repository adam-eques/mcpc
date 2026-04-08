package main

import (
	"context"
	"fmt"

	"github.com/adam-eques/mcpc/internal/config"
	"github.com/adam-eques/mcpc/internal/log"
)

func cmdPing(ctx context.Context, cfg config.Config, logger *log.Logger) error {
	c, info, err := connect(ctx, cfg, logger)
	if err != nil {
		return err
	}
	defer c.Close()
	if err := c.Ping(ctx); err != nil {
		return err
	}
	fmt.Printf("pong from %s %s\n", info.ServerInfo.Name, info.ServerInfo.Version)
	return nil
}
