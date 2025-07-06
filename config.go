package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	HatenaID   string `json:"hatena_id"`
	APIKey     string `json:"api_key"`
	BlogDomain string `json:"blog_domain"`
}

func loadConfig(configFile, hatenaID, apiKey, blogDomain string) (*Config, error) {
	config := &Config{
		HatenaID:   hatenaID,
		APIKey:     apiKey,
		BlogDomain: blogDomain,
	}

	if configFile != "" {
		return loadConfigFromFile(configFile)
	}

	defaultConfigPath := getDefaultConfigPath()
	if fileExists(defaultConfigPath) {
		fileConfig, err := loadConfigFromFile(defaultConfigPath)
		if err != nil {
			return nil, err
		}
		if config.HatenaID == "" {
			config.HatenaID = fileConfig.HatenaID
		}
		if config.APIKey == "" {
			config.APIKey = fileConfig.APIKey
		}
		if config.BlogDomain == "" {
			config.BlogDomain = fileConfig.BlogDomain
		}
	}

	return config, nil
}

func loadConfigFromFile(configFile string) (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

func saveConfig(config *Config, configFile string) error {
	if configFile == "" {
		configFile = getDefaultConfigPath()
	}

	err := ensureConfigDir()
	if err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	err = os.WriteFile(configFile, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func getDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./hatena-blog-org.json"
	}
	return filepath.Join(homeDir, ".config", "hatena-blog-org", "config.json")
}

func getConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(homeDir, ".config", "hatena-blog-org")
}

func ensureConfigDir() error {
	configDir := getConfigDir()
	return os.MkdirAll(configDir, 0755)
}
