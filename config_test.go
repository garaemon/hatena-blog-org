package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigFromFile(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "test_config.json")
	configContent := `{
		"hatena_id": "testuser",
		"api_key": "testapi",
		"blog_domain": "testblog.example.com"
	}`

	err := os.WriteFile(tmpFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	defer os.Remove(tmpFile)

	config, err := loadConfigFromFile(tmpFile)
	if err != nil {
		t.Fatalf("loadConfigFromFile failed: %v", err)
	}

	if config.HatenaID != "testuser" {
		t.Errorf("Expected HatenaID to be 'testuser', got '%s'", config.HatenaID)
	}
	if config.APIKey != "testapi" {
		t.Errorf("Expected APIKey to be 'testapi', got '%s'", config.APIKey)
	}
	if config.BlogDomain != "testblog.example.com" {
		t.Errorf("Expected BlogDomain to be 'testblog.example.com', got '%s'", config.BlogDomain)
	}
}

func TestLoadConfigFromFileNotFound(t *testing.T) {
	_, err := loadConfigFromFile("/nonexistent/config.json")
	if err == nil {
		t.Error("Expected error for nonexistent config file")
	}
}

func TestLoadConfigFromFileInvalidJSON(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "invalid_config.json")
	invalidContent := `{
		"hatena_id": "testuser",
		"api_key": "testapi"
		"blog_domain": "testblog.example.com"
	}`

	err := os.WriteFile(tmpFile, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	defer os.Remove(tmpFile)

	_, err = loadConfigFromFile(tmpFile)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestSaveConfig(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "save_test_config.json")
	defer os.Remove(tmpFile)

	config := &Config{
		HatenaID:   "testuser",
		APIKey:     "testapi",
		BlogDomain: "testblog.example.com",
	}

	err := saveConfig(config, tmpFile)
	if err != nil {
		t.Fatalf("saveConfig failed: %v", err)
	}

	if !fileExists(tmpFile) {
		t.Error("Config file should exist after saving")
	}

	loadedConfig, err := loadConfigFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if loadedConfig.HatenaID != config.HatenaID {
		t.Errorf("Expected HatenaID to be '%s', got '%s'", config.HatenaID, loadedConfig.HatenaID)
	}
	if loadedConfig.APIKey != config.APIKey {
		t.Errorf("Expected APIKey to be '%s', got '%s'", config.APIKey, loadedConfig.APIKey)
	}
	if loadedConfig.BlogDomain != config.BlogDomain {
		t.Errorf("Expected BlogDomain to be '%s', got '%s'", config.BlogDomain, loadedConfig.BlogDomain)
	}
}

func TestValidateConfig(t *testing.T) {
	validConfig := &Config{
		HatenaID:   "testuser",
		APIKey:     "testapi",
		BlogDomain: "testblog.example.com",
	}

	err := validateConfig(validConfig)
	if err != nil {
		t.Errorf("Valid config should not return error: %v", err)
	}

	invalidConfigs := []*Config{
		{APIKey: "testapi", BlogDomain: "testblog.example.com"},
		{HatenaID: "testuser", BlogDomain: "testblog.example.com"},
		{HatenaID: "testuser", APIKey: "testapi"},
	}

	for i, config := range invalidConfigs {
		err := validateConfig(config)
		if err == nil {
			t.Errorf("Invalid config %d should return error", i)
		}
	}
}

func TestGetDefaultConfigPath(t *testing.T) {
	path := getDefaultConfigPath()
	if path == "" {
		t.Error("Default config path should not be empty")
	}
}

func TestGetConfigDir(t *testing.T) {
	dir := getConfigDir()
	if dir == "" {
		t.Error("Config directory should not be empty")
	}
}

func TestLoadConfig(t *testing.T) {
	config, err := loadConfig("", "testuser", "testapi", "testblog.example.com")
	if err != nil {
		t.Fatalf("loadConfig failed: %v", err)
	}

	if config.HatenaID != "testuser" {
		t.Errorf("Expected HatenaID to be 'testuser', got '%s'", config.HatenaID)
	}
	if config.APIKey != "testapi" {
		t.Errorf("Expected APIKey to be 'testapi', got '%s'", config.APIKey)
	}
	if config.BlogDomain != "testblog.example.com" {
		t.Errorf("Expected BlogDomain to be 'testblog.example.com', got '%s'", config.BlogDomain)
	}
}
