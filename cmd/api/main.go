package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kakkoyun/demo-web-service/config"
	"github.com/kakkoyun/demo-web-service/handlers"
)

func main() {
	// Initialize structured logger
	logger := setupLogger()

	// Load configuration
	cfg := config.LoadConfig()
	logger.Info("Configuration loaded", "serverPort", cfg.ServerPort)

	// Initialize router using standard lib
	mux := http.NewServeMux()

	// Set up routes with Go 1.22 pattern syntax
	mux.HandleFunc("GET /", handlers.HomeHandler)
	mux.HandleFunc("GET /api/health", handlers.HealthCheckHandler)
	mux.HandleFunc("GET /api/users", handlers.GetUsersHandler)
	mux.HandleFunc("POST /api/users", handlers.CreateUserHandler)
	mux.HandleFunc("GET /api/users/{id}", handlers.GetUserHandler)

	logger.Info("Routes configured")

	// Apply middleware
	var handler http.Handler = mux
	handler = handlers.LoggingMiddleware(handler)

	// Configure server
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until signal is received
	<-c
	logger.Info("Server is shutting down...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited properly")
	os.Exit(0)
}

// setupLogger configures and returns a structured logger
func setupLogger() *slog.Logger {
	// Define log level based on environment (could use an environment variable)
	var logLevel slog.Level
	if os.Getenv("APP_ENV") == "production" {
		logLevel = slog.LevelInfo
	} else {
		logLevel = slog.LevelDebug
	}

	// Create a JSON handler for structured logging
	opts := &slog.HandlerOptions{
		Level: logLevel,
		// Add source code location to log entries in development
		AddSource: os.Getenv("APP_ENV") != "production",
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	// Set as default logger for compatibility with standard library
	slog.SetDefault(logger)

	return logger
}
