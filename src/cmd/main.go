package main

import (
	"flag"

	_ "github.com/GabrielEValenzuela/Payment-Registration-System/src/docs"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/config"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/pkg/logger"
)

// @title Payment Registration System
// @version 1.0
// @description Payment Registration System API for Database Management Course UNLP
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email https://github.com/GabrielEValenzuela/Payment-Registration-System
// @license.name MIT License
// @license.url http://www.apache.org/licenses/MIT.html
// @host localhost:8080
// @BasePath /
func main() {
	// Initialize the logger
	log := logger.InitLogger()

	// Parse command-line flags
	configPath := flag.String("config", "config.yaml", "path to the configuration file")
	flag.Parse()

	// Log the configuration file path being used
	log.Infof("Using configuration file: %s", *configPath)

	// Load configuration
	_, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Print log for now
	log.Infof("Configuration loaded successfully")

	// Create and run the server
	// srv := server.NewServer(cfg)
}
