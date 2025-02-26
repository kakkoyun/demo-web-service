package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/kakkoyun/demo-web-service/config"
	"github.com/kakkoyun/demo-web-service/handlers"
)

func main() {
	// Initialize structured logger
	logger := setupLogger()

	// Get build info
	buildInfo := getBuildInfo()

	// Log version information
	logger.Info("Starting application",
		"version", buildInfo.Version,
		"module", buildInfo.Module,
		"goVersion", buildInfo.GoVersion,
	)

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
	// Add version endpoint
	mux.HandleFunc("GET /api/version", versionHandler)

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

// VersionInfo stores application version information
type VersionInfo struct {
	Version   string `json:"version"`
	Module    string `json:"module"`
	GoVersion string `json:"goVersion"`
}

// getBuildInfo retrieves the build information from the binary
func getBuildInfo() VersionInfo {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return VersionInfo{
			Version:   "dev",
			Module:    "unknown",
			GoVersion: "unknown",
		}
	}

	// Extract the main module info
	var versionInfo VersionInfo
	versionInfo.Module = info.Main.Path
	versionInfo.Version = info.Main.Version
	versionInfo.GoVersion = info.GoVersion

	// If version isn't set (common in development builds), use a default
	if versionInfo.Version == "" {
		versionInfo.Version = "dev"
	}

	return versionInfo
}

// versionHandler returns the application version information
func versionHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Version information requested", "remote_addr", r.RemoteAddr)

	buildInfo := getBuildInfo()

	handlers.JSONResponse(w, http.StatusOK, buildInfo)
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
