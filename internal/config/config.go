package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// AuthConfig holds authentication settings
type AuthConfig struct {
	Enabled       bool   `json:"enabled"`
	Password      string `json:"password"`
	SessionSecret string `json:"sessionSecret,omitempty"`
}

// Config represents the application configuration
type Config struct {
	Port        int        `json:"port"`
	BindAddress string     `json:"bindAddress"`
	Auth        AuthConfig `json:"auth"`
}

var (
	currentConfig *Config
	configMutex   sync.RWMutex
)

// getConfigPath determines where config.json is stored
func getConfigPath() string {
	// 1. If config.json exists in current working dir, use it
	if _, err := os.Stat("config.json"); err == nil {
		return "config.json"
	}
	// 2. Otherwise use user's ~/.config/disk-cleaner/config.json (when run as macOS .app)
	home, err := os.UserHomeDir()
	if err == nil {
		dir := filepath.Join(home, ".config", "disk-cleaner")
		_ = os.MkdirAll(dir, 0755)
		return filepath.Join(dir, "config.json")
	}
	return "config.json"
}

// LoadConfig loads the configuration from disk, creating default if not found
func LoadConfig() (*Config, error) {
	configMutex.Lock()
	defer configMutex.Unlock()

	cfg := &Config{
		Port:        8089,
		BindAddress: "127.0.0.1",
		Auth: AuthConfig{
			Enabled:  false, // Desktop app on localhost does not require password by default
			Password: "admin", // fallback password if user manually enables auth
		},
	}

	configFile := getConfigPath()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Generate random session secret
		bytes := make([]byte, 16)
		_, _ = rand.Read(bytes)
		cfg.Auth.SessionSecret = hex.EncodeToString(bytes)

		data, err := json.MarshalIndent(cfg, "", "  ")
		if err == nil {
			_ = os.WriteFile(configFile, data, 0600)
		}
		currentConfig = cfg
		return cfg, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		currentConfig = cfg
		return cfg, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		currentConfig = cfg
		return cfg, err
	}

	if cfg.Auth.SessionSecret == "" {
		bytes := make([]byte, 16)
		_, _ = rand.Read(bytes)
		cfg.Auth.SessionSecret = hex.EncodeToString(bytes)
		SaveConfig(cfg)
	}

	if cfg.BindAddress == "" {
		cfg.BindAddress = "127.0.0.1"
	}

	currentConfig = cfg
	return cfg, nil
}

// GetConfig returns the cached config
func GetConfig() *Config {
	configMutex.RLock()
	defer configMutex.RUnlock()
	if currentConfig == nil {
		cfg, _ := LoadConfig()
		return cfg
	}
	return currentConfig
}

// SaveConfig saves configuration to file
func SaveConfig(cfg *Config) error {
	configFile := getConfigPath()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(configFile)
	if dir != "" && dir != "." {
		_ = os.MkdirAll(dir, 0755)
	}
	return os.WriteFile(configFile, data, 0600)
}
