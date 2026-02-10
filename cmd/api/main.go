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

	"github.com/Filipefr15/cadprev_apis/internal/api"
	"github.com/Filipefr15/cadprev_apis/internal/handler"
	"github.com/Filipefr15/cadprev_apis/internal/repository"
	"github.com/Filipefr15/cadprev_apis/internal/usecase"
)

func main() {
	// Initialize dependencies using constructor-based dependency injection

	// Infrastructure layer - API client
	externalAPIURL := getEnv("EXTERNAL_API_URL", "https://apicadprev.trabalho.gov.br")
	apiClient := api.NewStubAPIClient(externalAPIURL)
	log.Printf("Initialized stub API client with base URL: %s", externalAPIURL)

	// Infrastructure layer - Repository
	dataRepository := repository.NewInMemoryRepository()
	log.Println("Initialized in-memory repository")

	// Use case layer - Business logic orchestration
	ingestUseCase := usecase.NewIngestUseCase(apiClient, dataRepository)
	log.Println("Initialized ingest use case")

	// Handler layer - HTTP interface
	dataHandler := handler.NewDataHandler(ingestUseCase)
	log.Println("Initialized data handler")

	// Setup HTTP server
	mux := http.NewServeMux()
	dataHandler.RegisterRoutes(mux)

	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown handling
	go func() {
		log.Printf("Starting server on port %s...", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Println("Server is ready to handle requests")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server is shutting down...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
