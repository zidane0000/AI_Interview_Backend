// Entry point for the AI Interview Backend application
// Responsible for initializing configuration, database, router, and starting the server
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zidane0000/AI_Interview_Backend/api"
	"github.com/zidane0000/AI_Interview_Backend/config"
	"github.com/zidane0000/AI_Interview_Backend/data"
)

func main() {
	// Load configuration
	fmt.Println("Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize database connection
	fmt.Println("Initializing database connection...")
	db, err := data.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer data.CloseDB(db)

	// Set up router
	router := api.SetupRouter()

	fmt.Printf("Server starting on port %s...\n", cfg.Port)
	// Start HTTP server
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
