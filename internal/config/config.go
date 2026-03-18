// Package config defines the mcpc client configuration and loads it from an
// optional JSON file layered with environment-variable overrides.
package config

// Config is the top-level client configuration.
type Config struct {
	// Transport selects how to reach the server: "stdio" or "http".
	Transport string     `json:"transport"`
	Server    ServerSpec `json:"server"`
	HTTP      HTTPSpec   `json:"http"`
	Timeout   int        `json:"timeoutSeconds"`
	Log       LogConfig  `json:"log"`
	Client    ClientInfo `json:"client"`
}

// ServerSpec describes the server process for the stdio transport.
type ServerSpec struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// HTTPSpec describes the gateway endpoint for the http transport.
type HTTPSpec struct {
	Endpoint string `json:"endpoint"`
}

// LogConfig controls logging.
type LogConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// ClientInfo is advertised to the server during initialize.
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Default returns a configuration that launches `mcpkit` over stdio.
func Default() Config {
	return Config{
		Transport: "stdio",
		Server:    ServerSpec{Command: "mcpkit"},
		HTTP:      HTTPSpec{Endpoint: "http://localhost:8080/rpc"},
		Timeout:   30,
		Log:       LogConfig{Level: "info", Format: "text"},
		Client:    ClientInfo{Name: "mcpc", Version: "0.1.0"},
	}
}
