package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/farhan-nahid/email-service/initializers"
	"github.com/farhan-nahid/email-service/routes"
	"github.com/farhan-nahid/email-service/utils"
	"github.com/gin-gonic/gin"
)

// init() runs before main(), initializing environment variables and database connection
func init() {
	initializers.LoadEnvVariables() // Load environment variables from the .env file or system
	initializers.ConnectToDatabase() // Uncomment if you need database connection initialization
}

func main() {
	// Create a new Gin router instance
	router := gin.Default()

	router.GET("/health-check", func(c *gin.Context) {utils.SuccessResponse(c, http.StatusOK, nil, "Service is up and running")})
	routes.EmailRoute(router) // Register email routes

	// Define the HTTP server configuration
	server := &http.Server{
		Addr:   ":" + os.Getenv("PORT"), // Server will listen on port 8080
		Handler: router,  // Use Gin router as the handler
	}

	// Channel to capture OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt) // Listen for interrupt signals (e.g., Ctrl+C)

	// Start the server in a separate goroutine
	go func() {
		log.Println("Starting server on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Log fatal error if server fails to start and is not already shutting down
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for an interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Ensure context cancellation to free up resources

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		// Log fatal error if the server cannot shut down gracefully
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
