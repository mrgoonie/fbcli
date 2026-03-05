package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds application configuration
type Config struct {
	DefaultPageID   string `yaml:"default_page_id"`
	DefaultPageName string `yaml:"default_page_name"`
	AppID           string `yaml:"app_id"`
	AppSecret       string `yaml:"app_secret"`
	APIVersion      string `yaml:"api_version"`
}

// Dir returns the config directory path (~/.fbcli)
func Dir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".fbcli")
}

// FilePath returns the config file path
func FilePath() string {
	return filepath.Join(Dir(), "config.yaml")
}

// EnsureDir creates the config directory if it doesn't exist
func EnsureDir() error {
	return os.MkdirAll(Dir(), 0700)
}

// Load reads config from disk. Returns empty config if file doesn't exist.
func Load() (*Config, error) {
	cfg := &Config{
		APIVersion: "v24.0",
	}

	data, err := os.ReadFile(FilePath())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	// Override with env vars if set
	if v := os.Getenv("FBCLI_APP_ID"); v != "" {
		cfg.AppID = v
	}
	if v := os.Getenv("FBCLI_APP_SECRET"); v != "" {
		cfg.AppSecret = v
	}

	return cfg, nil
}

// Save writes config to disk
func Save(cfg *Config) error {
	if err := EnsureDir(); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(FilePath(), data, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}
