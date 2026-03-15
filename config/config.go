package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds user preferences
type Config struct {
	GroupID string `json:"group_id"`
	Version string `json:"version"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		GroupID: "com.example",
		Version: "1.0.0",
	}
}

// ConfigDir returns the configuration directory
func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".config", "fabric-cli")
}

// ConfigPath returns the full path to the config file
func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.json")
}

// Load reads the configuration from disk
func Load() (*Config, error) {
	path := ConfigPath()

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set defaults for empty values
	if cfg.GroupID == "" {
		cfg.GroupID = DefaultConfig().GroupID
	}
	if cfg.Version == "" {
		cfg.Version = DefaultConfig().Version
	}

	return &cfg, nil
}

// Save writes the configuration to disk
func (c *Config) Save() error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	path := ConfigPath()
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
