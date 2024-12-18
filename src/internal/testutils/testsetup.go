package testutils

import (
	"flag"
	"log"
	"sync"

	config "github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/config"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/pkg/logger"
)

var (
	DSN        string
	initOnce   sync.Once // Ensures initialization is done only once
	configPath *string
)

func InitTestSetup() {
	initOnce.Do(func() {
		// Initialize the logger
		logger.InitLogger(false, "")

		// Define the flag for the configuration file (only once)
		if configPath == nil {
			configPath = flag.String("config", "../../../../../config.yml", "path to the configuration file")
		}

		// Parse flags if not already done
		if !flag.Parsed() {
			flag.Parse()
		}

		// Load the configuration
		cfg, err := config.LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}

		// Extract the DSN
		DSN = cfg.SQLDb.DSN
	})
}
