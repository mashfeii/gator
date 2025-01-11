package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

const DBURL = "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUsername string `json:"current_username"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUsername = username

	return writeToConfig(c)
}

func ReadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	var config Config

	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return &config, nil
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to open home directory: %w", err)
	}

	configPath := path.Join(homeDir, configFileName)

	if _, err := os.Stat(configPath); err != nil {
		return "", fmt.Errorf("failed to find config file: %w", err)
	}

	return configPath, nil
}

func writeToConfig(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	byteInfo, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(configPath, byteInfo, 0o600)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
