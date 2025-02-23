package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"api-gateway/config"
	_ "api-gateway/docs" // Import swagger docs
	"api-gateway/internal/server"
)

// @title           Go Server API
// @version         1.0
// @description     A RESTful API server implemented in Go using Gin framework

// @contact.name   API Support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create new server instance
	srv, err := server.New(cfg)
	if err != nil {
		panic(err)
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the server
	fmt.Printf("Starting server on port %s...\n", cfg.Server.Port)
	if err := srv.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	// Wait for interrupt signal
	<-quit
	fmt.Println("\nShutdown signal received...")

	// Gracefully shutdown the server
	if err := srv.Stop(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exited properly")
}
