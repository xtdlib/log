package main

import (
	"github.com/xtdlib/log"
)

func main() {
	// The library automatically configures journald if available
	
	log.Trace("This is a trace message", "level", "trace")
	log.Debug("This is a debug message", "level", "debug")
	log.Info("This is an info message", "level", "info")
	log.Warn("This is a warning message", "level", "warn")
	log.Error("This is an error message", "level", "error")
	log.Emergency("This is an emergency message", "level", "emergency")
	
	// Test with groups and attributes
	logger := log.With("component", "demo")
	logger.Info("Message with component", "key1", "value1", "key2", 42)
	
	// Test with nested group
	groupLogger := logger.WithGroup("database")
	groupLogger.Error("Database connection failed", "host", "localhost", "port", 5432)
	
	log.Info("Demo completed - check journalctl -t <appname> to see the logs")
}