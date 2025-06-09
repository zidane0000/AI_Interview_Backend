// Entry point for the AI Interview Backend application
// Responsible for initializing configuration, database, router, and starting the server
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zidane0000/AI_Interview_Backend/api"
	"github.com/zidane0000/AI_Interview_Backend/config"
	"github.com/zidane0000/AI_Interview_Backend/data"
)

// gracefulShutdown handles graceful shutdown of the application
func gracefulShutdown(server *http.Server, timeout time.Duration) {
	// Create a channel to receive OS signals
	quit := make(chan os.Signal, 1)

	// Register the channel to receive specific signals
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal
	sig := <-quit
	log.Printf("Received signal: %v. Starting graceful shutdown...", sig)

	// Create a deadline to wait for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		os.Exit(1) // Exit with error code 1
	}

	// Additional cleanup operations
	log.Println("Performing cleanup operations...")
	// Close database connections if available
	if data.GlobalStore != nil {
		if err := data.GlobalStore.Close(); err != nil {
			log.Printf("Error closing database connections: %v", err)
			os.Exit(2) // Exit with error code 2 for database cleanup failure
		}
	}

	log.Println("Graceful shutdown completed successfully")
}

func main() {
	// Load configuration
	fmt.Println("Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// TODO: Initialize logging with proper configuration
	// TODO: Add structured logging with levels (debug, info, warn, error)
	// TODO: Add log rotation and file output options

	// Initialize hybrid store (auto-detects memory vs database backend)
	fmt.Println("Initializing data store...")
	err = data.InitGlobalStore()
	if err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}

	// Log the backend being used
	if data.GlobalStore.GetBackend() == data.BackendDatabase {
		fmt.Println("Using PostgreSQL database backend")
	} else {
		fmt.Println("Using in-memory store backend (set DATABASE_URL for database mode)")
	}
	// TODO: Add store health checks
	// if err := data.GlobalStore.Health(); err != nil {
	//     log.Fatalf("store health check failed: %v", err)
	// }

	// Initialize AI service client (global client is already initialized with .env support)
	fmt.Println("AI client initialized with .env configuration")

	// TODO: Add AI service validation
	// if err := ai.Client.TestConnection(); err != nil {
	//     log.Fatalf("failed to connect to AI service: %v", err)
	// }
	// TODO: Initialize file upload directory and permissions
	// if err := os.MkdirAll(cfg.UploadPath, 0755); err != nil {
	//     log.Fatalf("failed to create upload directory: %v", err)
	// }
	// Set up router
	router := api.SetupRouter()
	// TODO: Add HTTPS support with TLS configuration
	// TODO: Add health check endpoints
	// TODO: Add metrics and monitoring endpoints
	// TODO: Add API documentation serving (Swagger/OpenAPI)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	fmt.Printf("Server successfully started on port %s\n", cfg.Port)
	fmt.Printf("Frontend can now connect to: http://localhost:%s\n", cfg.Port)

	// Start graceful shutdown handler (this will block until shutdown signal)
	gracefulShutdown(server, cfg.ShutdownTimeout)
}
