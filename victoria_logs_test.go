package log

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestVictoriaLogsHandler(t *testing.T) {
	// Create a test server to mock Victoria Logs
	var receivedRequests []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receivedRequests = append(receivedRequests, string(body))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"ok"}`))
	}))
	defer server.Close()

	// Create Victoria Logs handler with test server endpoint
	handler := NewVictoriaLogsHandler(server.URL + "/insert/elasticsearch/_bulk")

	// Create logger with Victoria Logs handler
	logger := slog.New(handler)

	// Test basic logging
	logger.Info("test message", "key", "value")

	// Give time for async request
	time.Sleep(50 * time.Millisecond)

	// Verify request was sent
	if len(receivedRequests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(receivedRequests))
	}

	// Parse the bulk request
	lines := strings.Split(strings.TrimSpace(receivedRequests[0]), "\n")
	if len(lines) != 2 {
		t.Fatalf("Expected 2 lines in bulk request, got %d", len(lines))
	}

	// Verify create line
	if lines[0] != `{"create":{}}` {
		t.Errorf("Expected create line, got: %s", lines[0])
	}

	// Verify log entry
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(lines[1]), &logEntry); err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	// Check required fields
	if msg, ok := logEntry["_msg"].(string); !ok || msg != "test message" {
		t.Errorf("Expected _msg='test message', got: %v", logEntry["_msg"])
	}

	if level, ok := logEntry["level"].(string); !ok || level != "INFO" {
		t.Errorf("Expected level='INFO', got: %v", logEntry["level"])
	}

	if _, ok := logEntry["_time"].(string); !ok {
		t.Error("Missing _time field")
	}

	if _, ok := logEntry["host"].(string); !ok {
		t.Error("Missing host field")
	}

	// Check custom attribute
	if val, ok := logEntry["key"].(string); !ok || val != "value" {
		t.Errorf("Expected key='value', got: %v", logEntry["key"])
	}

	// Check source info (since AddSource is true)
	if _, ok := logEntry["source.file"].(string); !ok {
		t.Error("Missing source.file field")
	}
}

func TestVictoriaLogsHandlerWithGroups(t *testing.T) {
	var receivedRequests []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receivedRequests = append(receivedRequests, string(body))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	handler := NewVictoriaLogsHandler(server.URL + "/insert/elasticsearch/_bulk")
	logger := slog.New(handler)

	// Test with groups
	groupedLogger := logger.WithGroup("app").With("version", "1.0.0")
	groupedLogger.Info("grouped message", "status", "ok")

	time.Sleep(10 * time.Millisecond)

	if len(receivedRequests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(receivedRequests))
	}

	lines := strings.Split(strings.TrimSpace(receivedRequests[0]), "\n")
	var logEntry map[string]interface{}
	json.Unmarshal([]byte(lines[1]), &logEntry)

	// Check grouped attributes
	if val, ok := logEntry["app.version"].(string); !ok || val != "1.0.0" {
		t.Errorf("Expected app.version='1.0.0', got: %v", logEntry["app.version"])
	}

	if val, ok := logEntry["app.status"].(string); !ok || val != "ok" {
		t.Errorf("Expected app.status='ok', got: %v", logEntry["app.status"])
	}
}

func TestVictoriaLogsHandlerLevels(t *testing.T) {
	var receivedLevels []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		lines := strings.Split(strings.TrimSpace(string(body)), "\n")
		if len(lines) >= 2 {
			var logEntry map[string]interface{}
			json.Unmarshal([]byte(lines[1]), &logEntry)
			if level, ok := logEntry["level"].(string); ok {
				receivedLevels = append(receivedLevels, level)
			}
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	handler := NewVictoriaLogsHandler(server.URL + "/insert/elasticsearch/_bulk")
	logger := slog.New(handler)

	// Test all log levels
	logger.Log(context.Background(), LevelTrace, "trace msg")
	logger.Debug("debug msg")
	logger.Info("info msg")
	logger.Warn("warn msg")
	logger.Error("error msg")
	logger.Log(context.Background(), LevelEmergency, "emergency msg")

	time.Sleep(50 * time.Millisecond)

	expectedLevels := []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "EMERGENCY"}
	if len(receivedLevels) != len(expectedLevels) {
		t.Fatalf("Expected %d levels, got %d", len(expectedLevels), len(receivedLevels))
	}

	for i, expected := range expectedLevels {
		if receivedLevels[i] != expected {
			t.Errorf("Expected level[%d]='%s', got: '%s'", i, expected, receivedLevels[i])
		}
	}
}

func TestVictoriaLogsHandlerErrorHandling(t *testing.T) {
	// Test with invalid endpoint
	handler := NewVictoriaLogsHandler("http://invalid-endpoint:9999")
	logger := slog.New(handler)

	// This should not panic, just fail silently
	logger.Info("test message")

	// Test with server returning error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	}))
	defer server.Close()

	handler2 := NewVictoriaLogsHandler(server.URL)

	// This should not return an error since it's async now
	err := handler2.Handle(context.Background(), slog.NewRecord(time.Now(), slog.LevelInfo, "test", 0))
	if err != nil {
		t.Error("Expected no error for async handling")
	}
}

func TestVictoriaLogsHandlerWithBuffer(t *testing.T) {
	// Test to ensure handler works with multiHandler
	var buf bytes.Buffer
	consoleHandler := newConsoleHandler(&buf, nil)
	
	receivedCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	victoriaHandler := NewVictoriaLogsHandler(server.URL)
	
	multiHandler := &multiHandler{
		handlers: []slog.Handler{consoleHandler, victoriaHandler},
	}
	
	logger := slog.New(multiHandler)
	logger.Info("multi handler test", "handler", "both")

	time.Sleep(10 * time.Millisecond)

	// Check console output
	if !strings.Contains(buf.String(), "multi handler test") {
		t.Error("Expected message in console output")
	}

	// Check Victoria Logs received the log
	if receivedCount != 1 {
		t.Errorf("Expected 1 request to Victoria Logs, got %d", receivedCount)
	}
}