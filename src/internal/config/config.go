package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App     AppConfig
	SQLDb   SQLConfig
	NoSQLDb NoSQLConfig
}

type AppConfig struct {
	Port string
}

type SQLConfig struct {
	// Configs
}

type NoSQLConfig struct {
	// Configs
}

func LoadConfig(configPath string) (*Config, error) {
	// Set config file path
	viper.SetConfigFile(configPath)

	// Set default values for application
	viper.SetDefault("app.port", "8080")

	// Set default values for SQL connection

	// Set default values for NoSQL connection

	// Read in environment variables that match
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshal config
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
