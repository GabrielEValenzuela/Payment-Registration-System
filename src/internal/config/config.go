package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App     AppConfig
	MySQL   MySQLConfig
	MongoDB MongoDBConfig
}

type AppConfig struct {
	Port             string
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	GracefulShutdown time.Duration
}

type MySQLConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type MongoDBConfig struct {
	URI      string
	Database string
	// Other MongoDB-specific configs
}

func LoadConfig(configPath string) (*Config, error) {
	// Set config file path
	viper.SetConfigFile(configPath)

	// Set default values for application
	viper.SetDefault("app.port", "8080")
	viper.SetDefault("app.read_timeout", 15)
	viper.SetDefault("app.write_timeout", 15)
	viper.SetDefault("app.graceful_shutdown", 15)

	// Set default values for SQL connection
	viper.SetDefault("sql.host", "localhost")
	viper.SetDefault("sql.port", "3306")
	viper.SetDefault("sql.max_open_conns", 10)
	viper.SetDefault("sql.max_idle_conns", 10)
	viper.SetDefault("sql.conn_max_lifetime", 5)

	// Set default values for NoSQL connection
	viper.SetDefault("no-sql.uri", "mongodb://localhost:27017")
	viper.SetDefault("no-sql.database", "payment_registration_system")

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
