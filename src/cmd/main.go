package main

import (
	"flag"
	"fmt"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/cmd/server"
	_ "github.com/GabrielEValenzuela/Payment-Registration-System/src/docs"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/config"
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

	// Parse command-line flags
	configPath := flag.String("config", "config.yaml", "path to the configuration file")
	flag.Parse()
	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Errorf("Failed to load configuration: %v", err)
	}

	// Create and run the server
	srv := server.NewServer(cfg)
	srv.Run()
}
