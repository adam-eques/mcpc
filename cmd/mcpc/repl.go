package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/adam-eques/mcpc/internal/config"
	"github.com/adam-eques/mcpc/internal/log"
)

// cmdRepl runs an interactive loop. Type a tool name followed by key=value
// arguments, "tools" to list, "help" for usage, or "quit" to exit.
func cmdRepl(ctx context.Context, cfg config.Config, logger *log.Logger) error {
	c, info, err := connect(ctx, cfg, logger)
	if err != nil {
		return err
	}
	defer c.Close()

	fmt.Printf("connected to %s %s — type 'help' or 'quit'\n", info.ServerInfo.Name, info.ServerInfo.Version)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("mcpc> ")
		if !scanner.Scan() {
			fmt.Println()
			return nil
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		switch fields[0] {
		case "quit", "exit":
			return nil
		case "help":
			fmt.Println("commands: tools | describe <tool> | <tool> [k=v ...] | quit")
		case "tools":
			tools, err := c.ListTools(ctx)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				continue
			}
			for _, t := range tools {
				fmt.Printf("  %-16s %s\n", t.Name, t.Description)
			}
		default:
			args, err := parseArgs(fields[1:])
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				continue
			}
			res, err := c.CallTool(ctx, fields[0], args)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				continue
			}
			printResult(res)
		}
	}
}
