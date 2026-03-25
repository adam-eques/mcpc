package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.Transport != "stdio" || cfg.Server.Command != "mcpkit" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestLoadFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mcpc.json")
	os.WriteFile(path, []byte(`{"transport":"http","http":{"endpoint":"http://x/rpc"}}`), 0o644)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Transport != "http" || cfg.HTTP.Endpoint != "http://x/rpc" {
		t.Fatalf("file not applied: %+v", cfg)
	}
}

func TestEnvOverride(t *testing.T) {
	t.Setenv("MCPC_SERVER", "go run ./cmd/mcpkit")
	t.Setenv("MCPC_TIMEOUT", "5")
	cfg, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Command != "go" || len(cfg.Server.Args) != 2 {
		t.Fatalf("server not parsed: %+v", cfg.Server)
	}
	if cfg.Timeout != 5 {
		t.Fatalf("timeout=%d", cfg.Timeout)
	}
}
