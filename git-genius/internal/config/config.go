package config

import (
	"encoding/json"
	"os"
)

const Dir = ".git/.genius"
const File = Dir + "/config.json"

type Config struct {
	Branch string `json:"branch"`
	Remote string `json:"remote"`
}

func Load() Config {
	data, err := os.ReadFile(File)
	if err != nil {
		return Config{Branch: "main", Remote: "origin"}
	}
	var c Config
	json.Unmarshal(data, &c)
	return c
}

func Save(c Config) {
	os.MkdirAll(Dir, 0700)
	data, _ := json.MarshalIndent(c, "", "  ")
	os.WriteFile(File, data, 0600)
}
