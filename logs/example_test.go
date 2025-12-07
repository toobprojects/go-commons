package logs_test

import (
	"log/slog"
	"os"

	"github.com/toobprojects/go-commons/logs"
)

// ExampleInit_json demonstrates using JSON output for log aggregators like Splunk/ELK
func ExampleInit_json() {
	// Initialize with JSON output
	logs.Init(logs.Config{
		JSON:  true,
		Level: slog.LevelInfo,
		Out:   os.Stdout,
	})

	// Log some messages - these will be in JSON format
	logs.Info("Application started", "version", "1.0.0", "env", "production")
	logs.Info("Processing request", "user_id", 12345, "action", "login")
	logs.Warn("Cache miss", "key", "user:12345", "ttl", 300)
	logs.Error("Failed to connect", "service", "database", "retry", 3)
}

// ExampleInit_text demonstrates using human-readable text output for development
func ExampleInit_text() {
	// Initialize with text output
	logs.Init(logs.Config{
		JSON:  false,
		Level: slog.LevelDebug,
		Out:   os.Stdout,
		Color: false,
	})

	logs.Debug("Debug message")
	logs.Info("Info message", "key", "value")
}
