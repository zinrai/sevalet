package main

import (
	"flag"
	"log"
	"os"

	"github.com/zinrai/sevalet/internal/api"
	"github.com/zinrai/sevalet/internal/config"
)

func main() {
	// Parse command-line arguments
	configPath := flag.String("config", "configs/app.yaml", "Path to config file")
	flag.Parse()

	// Load configuration file
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and start server
	server := api.NewServer(cfg)
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
		os.Exit(1)
	}
}
