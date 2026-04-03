// Command mcpc is a command-line Model Context Protocol client. It connects to a
// server over stdio (launching it as a child process) or HTTP and exposes
// subcommands to inspect and invoke tools.
//
// Usage:
//
//	mcpc [flags] <command> [args]
//
// Commands:
//
//	ping                     check the server is responsive
//	tools                    list the server's tools
//	describe <tool>          print a tool's input schema
//	call <tool> [k=v ...]    invoke a tool
//	run <plan.json>          execute a scripted plan of tool calls
//	repl                     start an interactive session
//	version                  print the client version
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/adam-eques/mcpc/internal/config"
	"github.com/adam-eques/mcpc/internal/log"
	"github.com/adam-eques/mcpc/internal/version"
)

func main() {
	configPath := flag.String("config", "", "path to a JSON config file")
	server := flag.String("server", "", `server command for stdio, e.g. "mcpkit" or "go run ./cmd/mcpkit"`)
	httpEndpoint := flag.String("http", "", "gateway endpoint for the http transport")
	timeout := flag.Int("timeout", 0, "per-request timeout in seconds")
	logLevel := flag.String("log-level", "", "log level: debug, info, warn, error")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(2)
	}
	if args[0] == "version" {
		version.FromBuildInfo()
		fmt.Println("mcpc", version.String())
		return
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		fail(err)
	}
	applyFlags(&cfg, *server, *httpEndpoint, *timeout, *logLevel)

	logger := log.New(log.Options{
		Level:  log.ParseLevel(cfg.Log.Level),
		Format: log.Format(cfg.Log.Format),
		Writer: os.Stderr,
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := dispatch(ctx, cfg, logger, args[0], args[1:]); err != nil {
		fail(err)
	}
}

func dispatch(ctx context.Context, cfg config.Config, logger *log.Logger, cmd string, rest []string) error {
	switch cmd {
	case "ping":
		return cmdPing(ctx, cfg, logger)
	case "tools":
		return cmdTools(ctx, cfg, logger)
	case "describe":
		return cmdDescribe(ctx, cfg, logger, rest)
	case "call":
		return cmdCall(ctx, cfg, logger, rest)
	case "run":
		return cmdRun(ctx, cfg, logger, rest)
	case "repl":
		return cmdRepl(ctx, cfg, logger)
	default:
		return fmt.Errorf("unknown command %q (try: ping, tools, describe, call, run, repl)", cmd)
	}
}

func applyFlags(cfg *config.Config, server, httpEndpoint string, timeout int, logLevel string) {
	if server != "" {
		fields := strings.Fields(server)
		cfg.Transport = "stdio"
		cfg.Server.Command = fields[0]
		cfg.Server.Args = fields[1:]
	}
	if httpEndpoint != "" {
		cfg.Transport = "http"
		cfg.HTTP.Endpoint = httpEndpoint
	}
	if timeout > 0 {
		cfg.Timeout = timeout
	}
	if logLevel != "" {
		cfg.Log.Level = logLevel
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `mcpc — a Model Context Protocol client

Usage:
  mcpc [flags] <command> [args]

Commands:
  ping                     check the server is responsive
  tools                    list the server's tools
  describe <tool>          print a tool's input schema
  call <tool> [k=v ...]    invoke a tool (use k:=json for non-string values)
  run <plan.json>          execute a scripted plan of tool calls
  repl                     start an interactive session
  version                  print the client version

Flags:
`)
	flag.PrintDefaults()
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "mcpc:", err)
	os.Exit(1)
}
