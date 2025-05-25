package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServerURL    string `json:"server_url"`
	ConnectKey   string `json:"connect_key"`
	PollInterval int    `json:"poll_interval_seconds"`
	LogLevel     string `json:"log_level"`
	// Add more configuration parameters as needed
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
