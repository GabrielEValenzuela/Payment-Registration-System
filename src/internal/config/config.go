/*
 * Payment Registration System
 * ----------------------------------------
 * This file is part of the Payment Registration System, responsible for loading
 * and managing application configuration settings.
 *
 * Created: Oct. 19, 2024
 * License: GNU General Public License v3.0
 */

package config

import (
	"github.com/spf13/viper"
)

/*
 * Config
 * ----------------------------------------
 * Represents the application configuration structure.
 * It includes settings for the application, SQL database, NoSQL database,
 * and other relevant configurations.
 */
type Config struct {
	App          AppConfig   // Application-level configuration
	SQLDb        SQLConfig   // SQL database connection settings
	NoSQLDb      NoSQLConfig // NoSQL database connection settings
	IsProduction bool        // Flag indicating if the app runs in production mode
	LogPath      string      // Path for logging
}

/*
 * AppConfig
 * ----------------------------------------
 * Defines the application-specific configurations.
 */
type AppConfig struct {
	Port string // Application's listening port
}

/*
 * SQLConfig
 * ----------------------------------------
 * Defines the SQL database configuration.
 */
type SQLConfig struct {
	DSN   string // Data Source Name for connecting to the SQL database
	Clean bool   // Whether to clean the SQL database on startup
}

/*
 * NoSQLConfig
 * ----------------------------------------
 * Defines the NoSQL database configuration.
 */
type NoSQLConfig struct {
	URI      string // Connection URI for MongoDB
	Database string // Name of the MongoDB database
	Clean    bool   // Whether to clean the NoSQL database on startup
}

/*
 * LoadConfig
 * ----------------------------------------
 * Loads the application configuration from a specified file.
 * It applies default values and environment variable overrides.
 *
 * Parameters:
 * - configPath (string): The path to the configuration file.
 *
 * Returns:
 * - *Config: Pointer to the loaded configuration.
 * - error: Any error encountered while reading or parsing the configuration.
 */
func LoadConfig(configPath string) (*Config, error) {
	// Set config file path
	viper.SetConfigFile(configPath)

	// Set default values for application
	viper.SetDefault("app.port", "8080")
	viper.SetDefault("app.log_path", "payment_system.log")
	viper.SetDefault("app.is_production", false)

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
