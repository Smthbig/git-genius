package config

import (
	"encoding/json"
	"os"
)

const (
	Dir  = ".git/.genius"
	File = Dir + "/config.json"
)

// Config holds Git Genius configuration
type Config struct {
	Branch string `json:"branch"`
	Remote string `json:"remote"`

	// GitHub specific
	Owner string `json:"owner"` // username or organisation
	Repo  string `json:"repo"`  // repository name
}

// Load reads config from .git/.genius/config.json
// Falls back to safe defaults if file is missing or partial
func Load() Config {
	data, err := os.ReadFile(File)
	if err != nil {
		return defaultConfig()
	}

	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return defaultConfig()
	}

	// Backward compatibility (older configs)
	if c.Branch == "" {
		c.Branch = "main"
	}
	if c.Remote == "" {
		c.Remote = "origin"
	}

	return c
}

// Save writes config to disk with secure permissions
func Save(c Config) {
	os.MkdirAll(Dir, 0700)

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return
	}

	_ = os.WriteFile(File, data, 0600)
}

// defaultConfig returns sane defaults
func defaultConfig() Config {
	return Config{
		Branch: "main",
		Remote: "origin",
		Owner:  "",
		Repo:   "",
	}
}
