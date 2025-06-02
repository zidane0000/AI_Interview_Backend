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

	// TODO: Initialize logging with proper configuration
	// TODO: Add structured logging with levels (debug, info, warn, error)
	// TODO: Add log rotation and file output options

	// Initialize database connection
	fmt.Println("Initializing database connection...")
	db, err := data.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer data.CloseDB(db)

	// TODO: Run database migrations on startup
	// TODO: Add database health checks
	// TODO: Initialize connection pooling with proper settings

	// TODO: Initialize AI service client
	// aiClient := ai.NewClient(cfg.AIProvider, cfg.AIAPIKey, cfg.AIBaseURL)
	// if err := aiClient.TestConnection(); err != nil {
	//     log.Fatalf("failed to connect to AI service: %v", err)
	// }

	// TODO: Initialize business services with dependency injection
	// interviewService := business.NewInterviewService(db, aiClient)
	// evaluationService := business.NewEvaluationService(db, aiClient)
	// chatService := business.NewChatService(db, aiClient)

	// TODO: Initialize file upload directory and permissions
	// if err := os.MkdirAll(cfg.UploadPath, 0755); err != nil {
	//     log.Fatalf("failed to create upload directory: %v", err)
	// }

	// Set up router
	// TODO: Pass services to router for dependency injection
	router := api.SetupRouter()

	// TODO: Add graceful shutdown handling
	// TODO: Add HTTPS support with TLS configuration
	// TODO: Add health check endpoints
	// TODO: Add metrics and monitoring endpoints
	// TODO: Add API documentation serving (Swagger/OpenAPI)

	fmt.Printf("Server starting on port %s...\n", cfg.Port)
	fmt.Printf("Frontend can now connect to: http://localhost:%s\n", cfg.Port)

	// TODO: Implement graceful shutdown with signal handling
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// server := &http.Server{Addr: ":" + cfg.Port, Handler: router}
	// go func() { log.Fatal(server.ListenAndServe()) }()
	// <-c
	// server.Shutdown(context.Background())

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
