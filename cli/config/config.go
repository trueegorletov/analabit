package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	Varsities struct {
		List     []string `mapstructure:"list"`
		Excluded []string `mapstructure:"excluded"`
	} `mapstructure:"varsities"`
	DrainSim struct {
		Stages     []int `mapstructure:"stages"`
		Iterations int   `mapstructure:"iterations"`
	} `mapstructure:"drain_sim"`
	Upload struct {
		Database struct {
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			User     string `mapstructure:"user"`
			DBName   string `mapstructure:"dbname"`
			Password string `mapstructure:"password"`
		} `mapstructure:"database"`
	} `mapstructure:"upload"`
	Cache struct {
		Directory  string `mapstructure:"directory"`
		TTLMinutes int    `mapstructure:"ttl_minutes"`
	} `mapstructure:"cache"`
	Logging struct {
		File string `mapstructure:"file"` // Path to the log file. If empty, logs to stderr.
	} `mapstructure:"logging"`
}

var AppConfig Config

func LoadConfig(configPath string) error {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.AddConfigPath("./cli") // For local development
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/analabit/") // Example global path
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
	}

	viper.AutomaticEnv() // Read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}

	return nil
}
