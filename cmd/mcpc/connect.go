package main

import (
	"context"
	"fmt"
	"time"

	"github.com/adam-eques/mcpc/client"
	"github.com/adam-eques/mcpc/internal/config"
	"github.com/adam-eques/mcpc/internal/log"
	"github.com/adam-eques/mcpc/mcp"
)

// connect dials the server described by cfg and completes the handshake.
func connect(ctx context.Context, cfg config.Config, logger *log.Logger) (*client.Client, *mcp.InitializeResult, error) {
	opts := []client.Option{
		client.WithLogger(logger),
		client.WithRequestTimeout(time.Duration(cfg.Timeout) * time.Second),
		client.WithClientInfo(mcp.Implementation{Name: cfg.Client.Name, Version: cfg.Client.Version}),
	}
	switch cfg.Transport {
	case "http":
		return client.DialHTTP(ctx, cfg.HTTP.Endpoint, opts...)
	case "stdio", "":
		if cfg.Server.Command == "" {
			return nil, nil, fmt.Errorf("no server command configured; set -server or a config file")
		}
		return client.DialCommand(ctx, cfg.Server.Command, cfg.Server.Args, opts...)
	default:
		return nil, nil, fmt.Errorf("unknown transport %q", cfg.Transport)
	}
}
