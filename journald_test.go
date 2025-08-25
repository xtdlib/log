package log

import (
	"context"
	"log/slog"
	"os/exec"
	"testing"
	"time"
)

func TestJournaldHandler(t *testing.T) {
	// Skip if systemd-cat is not available
	if _, err := exec.LookPath("systemd-cat"); err != nil {
		t.Skip("systemd-cat not available, skipping journald tests")
	}

	handler := NewJournaldHandler()
	
	// Test basic logging
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	record.AddAttrs(slog.String("key", "value"))
	
	err := handler.Handle(context.Background(), record)
	if err != nil {
		t.Errorf("Failed to handle log record: %v", err)
	}
}

func TestJournaldHandlerLevels(t *testing.T) {
	// Skip if systemd-cat is not available
	if _, err := exec.LookPath("systemd-cat"); err != nil {
		t.Skip("systemd-cat not available, skipping journald tests")
	}

	handler := NewJournaldHandler()
	
	testCases := []struct {
		level slog.Level
		msg   string
	}{
		{LevelTrace, "trace message"},
		{LevelDebug, "debug message"},
		{LevelInfo, "info message"},
		{LevelWarn, "warning message"},
		{LevelError, "error message"},
		{LevelEmergency, "emergency message"},
	}
	
	for _, tc := range testCases {
		record := slog.NewRecord(time.Now(), tc.level, tc.msg, 0)
		err := handler.Handle(context.Background(), record)
		if err != nil {
			t.Errorf("Failed to handle %s: %v", tc.level, err)
		}
	}
}

func TestJournaldHandlerWithGroups(t *testing.T) {
	// Skip if systemd-cat is not available
	if _, err := exec.LookPath("systemd-cat"); err != nil {
		t.Skip("systemd-cat not available, skipping journald tests")
	}

	handler := NewJournaldHandler()
	
	// Test with groups
	groupHandler := handler.WithGroup("mygroup")
	attrHandler := groupHandler.WithAttrs([]slog.Attr{
		slog.String("app", "test"),
	})
	
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "grouped message", 0)
	record.AddAttrs(slog.String("field", "value"))
	
	err := attrHandler.Handle(context.Background(), record)
	if err != nil {
		t.Errorf("Failed to handle grouped log record: %v", err)
	}
}