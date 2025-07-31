package log

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var httpClient *http.Client

var appName = filepath.Base(os.Args[0])

func init() {
	if appName == "" {
		appName = "unknown"
	}

	httpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			// ForceAttemptHTTP2:     true,
			MaxIdleConnsPerHost: runtime.GOMAXPROCS(0) + 1,
		},
	}

	// Create colorful console handler
	consoleHandler := newConsoleHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     LevelDebug,
	})

	handlers := []slog.Handler{consoleHandler}

	// Try standard log directories in order of preference
	var logDir string
	var logPath string

	// 1. Try /var/log/{appname}/ (standard Linux location)
	logDir = filepath.Join("/var/log", appName)
	logFileName := fmt.Sprintf("%s.log", appStartTime.Format("20060102T150405"))
	logPath = filepath.Join(logDir, logFileName)

	// Test if we can write to /var/log
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// 2. Fallback to temp directory
		logDir = filepath.Join(os.TempDir(), appName, "logs")
		logPath = filepath.Join(logDir, logFileName)
	}

	if err := os.MkdirAll(logDir, 0755); err == nil {
		// Create log file with timestamp

		if file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
			logFile = file
			// ReplaceAttr for JSON to handle custom level names
			replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.LevelKey {
					level := a.Value.Any().(slog.Level)
					switch level {
					case LevelTrace:
						return slog.String(slog.LevelKey, "TRACE")
					case LevelEmergency:
						return slog.String(slog.LevelKey, "EMER")
					}
				}
				return a
			}
			fileHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
				AddSource:   true,
				Level:       LevelTrace,
				ReplaceAttr: replaceAttr,
			})
			handlers = append(handlers, fileHandler)
		} else {
			defer slog.Error(fmt.Sprintf("log: %v", err))
		}
	}

	// Add Victoria Logs handler if configured
	victoriaHandler := NewVictoriaLogsHandler(os.Getenv("VICTORIA_LOGS_ENDPOINT"))
	handlers = append(handlers, victoriaHandler)

	// Create multi-handler
	handler := &multiHandler{handlers: handlers}
	defaultLogger = slog.New(handler)

	// Set log as the default slog logger
	slog.SetDefault(defaultLogger)
	// slog.Info(fmt.Sprintf("log: log to %v", logPath))
}

func SetDefault(l *slog.Logger) {
	defaultLogger = l
	slog.SetDefault(l)
}

func Default() *slog.Logger {
	return defaultLogger
}
