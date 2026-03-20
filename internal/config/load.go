package config

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

// Load returns the default configuration merged with an optional JSON file and
// then MCPC_* environment overrides. A path of "" skips the file.
func Load(path string) (Config, error) {
	cfg := Default()
	if path != "" {
		raw, err := os.ReadFile(path)
		if err != nil {
			return cfg, err
		}
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return cfg, err
		}
	}
	applyEnv(&cfg)
	return cfg, nil
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("MCPC_TRANSPORT"); v != "" {
		cfg.Transport = v
	}
	if v := os.Getenv("MCPC_SERVER"); v != "" {
		fields := strings.Fields(v)
		cfg.Server.Command = fields[0]
		cfg.Server.Args = fields[1:]
	}
	if v := os.Getenv("MCPC_HTTP_ENDPOINT"); v != "" {
		cfg.HTTP.Endpoint = v
	}
	if v := os.Getenv("MCPC_LOG_LEVEL"); v != "" {
		cfg.Log.Level = v
	}
	if v := os.Getenv("MCPC_TIMEOUT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Timeout = n
		}
	}
}
