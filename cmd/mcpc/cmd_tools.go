package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/adam-eques/mcpc/internal/config"
	"github.com/adam-eques/mcpc/internal/log"
)

func cmdTools(ctx context.Context, cfg config.Config, logger *log.Logger) error {
	c, _, err := connect(ctx, cfg, logger)
	if err != nil {
		return err
	}
	defer c.Close()
	tools, err := c.ListTools(ctx)
	if err != nil {
		return err
	}
	if len(tools) == 0 {
		fmt.Println("(no tools)")
		return nil
	}
	for _, t := range tools {
		fmt.Printf("%-16s %s\n", t.Name, t.Description)
	}
	return nil
}

func cmdDescribe(ctx context.Context, cfg config.Config, logger *log.Logger, rest []string) error {
	if len(rest) != 1 {
		return fmt.Errorf("usage: mcpc describe <tool>")
	}
	c, _, err := connect(ctx, cfg, logger)
	if err != nil {
		return err
	}
	defer c.Close()
	tools, err := c.ListTools(ctx)
	if err != nil {
		return err
	}
	for _, t := range tools {
		if t.Name == rest[0] {
			fmt.Printf("%s — %s\n\n", t.Name, t.Description)
			var pretty any
			json.Unmarshal(t.InputSchema, &pretty)
			out, _ := json.MarshalIndent(pretty, "", "  ")
			fmt.Println(string(out))
			return nil
		}
	}
	return fmt.Errorf("tool %q not found", rest[0])
}
