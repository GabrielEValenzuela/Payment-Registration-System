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
	DSN   string
	Clean bool
}

type NoSQLConfig struct {
	URI      string
	Database string
	Clean    bool
}

func LoadConfig(configPath string) (*Config, error) {
	// Set config file path
	viper.SetConfigFile(configPath)

	// Set default values for application
	viper.SetDefault("app.port", "8080")

	// Set default values for SQL connection
	viper.SetDefault("sqldb.dsn", "root:password@tcp(localhost:3306)/payment_registration_system?charset=utf8mb4&parseTime=True&loc=Local")
	viper.SetDefault("sqldb.clean", false)

	// Set default values for NoSQL connection
	viper.SetDefault("nosqldb.uri", "mongodb://localhost:27017")
	viper.SetDefault("nosqldb.database", "payment_registration_system")
	viper.SetDefault("nosqldb.clean", false)

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
