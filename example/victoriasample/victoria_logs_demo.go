package main

import (
	"context"
	"os"
	"time"

	"github.com/xtdlib/log"
)

func main() {
	// To enable Victoria Logs, set the VICTORIA_LOGS_ENDPOINT environment variable:
	// export VICTORIA_LOGS_ENDPOINT=http://oci-aca-001:9428/insert/elasticsearch/_bulk

	// Or set it programmatically for this demo
	if os.Getenv("VICTORIA_LOGS_ENDPOINT") == "" {
		os.Setenv("VICTORIA_LOGS_ENDPOINT", "http://oci-aca-001:9428/insert/elasticsearch/_bulk")
	}

	// The logging library will automatically detect the environment variable
	// and add Victoria Logs handler during initialization

	// Basic logging examples
	log.Info("Application started", "version", "1.0.0")

	// Log with structured data
	log.Info("User logged in",
		"user_id", 12345,
		"username", "john_doe",
		"ip_address", "192.168.1.100",
		"session_duration", time.Hour,
	)

	// Different log levels
	log.Trace("Detailed trace information", "debug_data", map[string]interface{}{
		"memory_usage": "45MB",
		"goroutines":   10,
	})

	log.Debug("Debug message", "component", "auth")
	log.Warn("Warning: High memory usage", "percentage", 85.5)
	log.Error("Failed to connect to database", "error", "connection timeout", "retry_count", 3)

	// Using groups
	logger := log.WithGroup("api")
	logger = logger.With("request_id", "abc-123")

	logger.Info("API request received",
		"method", "GET",
		"path", "/api/v1/users",
		"duration_ms", 125,
	)

	// Context logging
	ctx := context.WithValue(context.Background(), "trace_id", "xyz-789")
	log.InfoContext(ctx, "Processing payment",
		"amount", 99.99,
		"currency", "USD",
		"payment_method", "credit_card",
	)

	// Emergency level (highest priority)
	log.Emergency("System critical failure detected!",
		"component", "database",
		"action", "immediate_restart_required",
	)

	// Ensure all logs are sent before exiting
	log.Close()

	// To verify logs in Victoria Logs:
	// curl http://oci-aca-001:9428/select/logsql/query -d 'query=*'
}
