/*
 * Payment Registration System - Main Server
 * --------------------------------------------------
 * This file is the entry point for the Payment Registration System API.
 * It initializes the server and runs it with the specified configuration.
 *
 * Created: Oct. 19, 2024
 * License: GNU General Public License v3.0
 */

package main

import (
	"flag"
	"log"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/cmd/server"
	_ "github.com/GabrielEValenzuela/Payment-Registration-System/src/docs"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/config"
)

// @title			Payment Registration System API
// @version		1.0
// @description	This API manages payment registration and processing for the Database Management Course at UNLP.
// @termsOfService	http://swagger.io/terms/
// @contact.name	@marventu94, @GabrielEValenzuela
// @contact.email	https://github.com/GabrielEValenzuela/Payment-Registration-System
// @license.name	GNU General Public License v3.0
// @license.url	https://www.gnu.org/licenses/gpl-3.0.html
// @host			go-app.localhost
// @BasePath		/v1
func main() {

	// Parse command-line flags
	configPath := flag.String("config", "./config.yml", "path to the configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("❌ Failed to load configuration: ", err)
		// We considered as error, but we use the default configuration
	}

	// Create and run the server
	srv := server.NewServer(cfg)
	if err := srv.Run(); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
