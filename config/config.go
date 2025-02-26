package config

import (
	"os"
	"strings"
	"time"
)

// Config holds the application configuration
type Config struct {
	ServerPort     string
	AllowedOrigins []string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
}

// LoadConfig loads the configuration from environment variables
// with sensible defaults
func LoadConfig() *Config {
	return &Config{
		ServerPort:     env("SERVER_PORT", "8080"),
		ReadTimeout:    durationEnv("READ_TIMEOUT", "15s"),
		WriteTimeout:   durationEnv("WRITE_TIMEOUT", "15s"),
		IdleTimeout:    durationEnv("IDLE_TIMEOUT", "60s"),
		AllowedOrigins: sliceEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080"),
	}
}

// env gets an environment variable or returns a fallback value
func env(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// durationEnv gets a duration environment variable or returns a fallback value
func durationEnv(key, fallback string) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}

	duration, _ := time.ParseDuration(fallback)
	return duration
}

// sliceEnv gets a slice from a comma-separated environment variable or returns a fallback
func sliceEnv(key, fallback string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return parseSlice(value)
	}
	return parseSlice(fallback)
}

// parseSlice parses a comma-separated string into a slice
func parseSlice(value string) []string {
	if value == "" {
		return []string{}
	}
	return splitAndTrim(value, ",")
}

// splitAndTrim splits a string by a separator and trims spaces from each element
func splitAndTrim(s, sep string) []string {
	parts := []string{}
	for _, part := range splitString(s, sep) {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

// splitString splits a string by a separator
func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, sep)
}
