package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	ConfigFileName = "config"
)

var configPath string

func InitConfig(path string) error {
	configPath = path
	if err := os.MkdirAll(configPath, 0700); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType("yaml")
	return nil
}

func FetchConfig(cmd string) (map[string]interface{}, error) {
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	config := viper.GetStringMap(cmd)
	return config, nil
}

func SetConfigValue(cmd, key string, value interface{}) {
	viper.Set(fmt.Sprintf("%s.%s", cmd, key), value)
}
