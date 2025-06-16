package main

import (
	"log"
	"micro_geoip/internal/api"
	"micro_geoip/internal/config"
	"micro_geoip/internal/geoip"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize GeoIP service
	geoipService, err := geoip.NewService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize GeoIP service: %v", err)
	}

	// Start the API server
	server := api.NewServer(cfg, geoipService)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}