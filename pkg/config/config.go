package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	ConfigFileName = "config"
	DefaultConfigDir = ".devctl"
)

var configPath string

// InitConfig initializes the configuration system with the given path.
// If path is empty, it defaults to ~/.devctl
func InitConfig(path string) error {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		path = filepath.Join(home, DefaultConfigDir)
	}
	
	configPath = path
	if err := os.MkdirAll(configPath, 0700); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}
	
	viper.AddConfigPath(configPath)
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType("yaml")
	
	// Try to read existing config, but don't error if it doesn't exist
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}
	
	return nil
}

func FetchConfig(cmd string) (map[string]interface{}, error) {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file doesn't exist yet, return empty map
			return make(map[string]interface{}), nil
		}
		return nil, err
	}
	config := viper.GetStringMap(cmd)
	return config, nil
}

func SetConfigValue(cmd, key string, value interface{}) {
	viper.Set(fmt.Sprintf("%s.%s", cmd, key), value)
}

func WriteConfig() error {
	configFile := filepath.Join(configPath, ConfigFileName+".yaml")
	return viper.WriteConfigAs(configFile)
}

func DeleteConfigValue(cmd, key string) error {
	// First read the current config to ensure we have the latest state
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}
	
	// Get all settings and make a deep copy to modify
	allSettings := viper.AllSettings()
	
	// Navigate to the command section and delete the key
	if cmdSection, exists := allSettings[cmd].(map[string]interface{}); exists {
		delete(cmdSection, key)
		// Reset viper and rebuild the configuration
		viper.Reset()
		viper.AddConfigPath(configPath)
		viper.SetConfigName(ConfigFileName)
		viper.SetConfigType("yaml")
		
		// Set all the modified settings
		for k, v := range allSettings {
			viper.Set(k, v)
		}
	}
	
	return WriteConfig()
}
