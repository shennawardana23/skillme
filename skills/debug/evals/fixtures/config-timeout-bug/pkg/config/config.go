package config

import "time"

// Config holds server configuration loaded at startup.
type Config struct {
	Host    string
	Port    int
	Timeout *time.Duration
}

// Load reads configuration for the server. Environment-variable parsing is
// omitted from this fixture — only Host and Port are populated.
func Load() *Config {
	return &Config{
		Host: "0.0.0.0",
		Port: 8080,
	}
}
