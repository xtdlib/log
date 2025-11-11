package log

import (
	"log/slog"
	"os"
	"path/filepath"
)

var appName = filepath.Base(os.Args[0])

func init() {
	if appName == "" {
		appName = "unknown"
	}

	var handlers []slog.Handler

	// Check if we're running in a TTY (console) or not (like systemd service/kubernetes)
	// if isInteractive() {
	if true {
		// Running in console - use colorful console handler
		consoleHandler := newConsoleHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     LevelDebug,
		})
		handlers = []slog.Handler{consoleHandler}
	} else {
		// Running as service/kubernetes - use JSON handler to stdout
		jsonHandler := newJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     LevelDebug,
		})
		handlers = []slog.Handler{jsonHandler}
	}

	// Create multi-handler
	handler := &multiHandler{handlers: handlers}
	defaultLogger = slog.New(handler)

	// Set log as the default slog logger
	slog.SetDefault(defaultLogger)
}

func SetDefault(l *slog.Logger) {
	defaultLogger = l
	slog.SetDefault(l)
}

func Default() *slog.Logger {
	return defaultLogger
}

// isInteractive checks if we're running in an interactive terminal (TTY)
func isInteractive() bool {
	// Check if stdout is a terminal
	if fileInfo, err := os.Stdout.Stat(); err == nil {
		return (fileInfo.Mode() & os.ModeCharDevice) != 0
	}
	return false
}
