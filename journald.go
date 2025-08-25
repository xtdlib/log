package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// JournaldHandler sends logs to systemd journal using systemd-cat in JSON format
type JournaldHandler struct {
	attrs []slog.Attr
	group string
}

// NewJournaldHandler creates a new handler that sends JSON logs to systemd journal
func NewJournaldHandler() *JournaldHandler {
	return &JournaldHandler{}
}

func (h *JournaldHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *JournaldHandler) Handle(ctx context.Context, r slog.Record) error {
	// Map slog levels to systemd priority levels
	priority := h.getPriority(r.Level)
	
	// Build JSON log entry
	entry := make(map[string]interface{})
	
	// Standard fields
	entry["message"] = r.Message
	entry["timestamp"] = r.Time.Format(time.RFC3339Nano)
	entry["level"] = h.getLevelName(r.Level)
	entry["app"] = appName
	
	// Add hostname
	if hostname, err := os.Hostname(); err == nil {
		entry["host"] = hostname
	}
	
	// Add source information if available
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		entry["source"] = map[string]interface{}{
			"file":     f.File,
			"line":     f.Line,
			"function": f.Function,
		}
	}
	
	// Add prepended attributes
	for _, a := range h.attrs {
		h.addAttrToMap(entry, a, h.group)
	}
	
	// Add record attributes
	r.Attrs(func(a slog.Attr) bool {
		h.addAttrToMap(entry, a, h.group)
		return true
	})
	
	// Convert to JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}
	
	// Use systemd-cat to send to journal with a prefix to distinguish from other logs
	// Using "xtdlog-" prefix for easy identification in journald
	identifier := fmt.Sprintf("xtdlog-%s", appName)
	cmd := exec.Command("systemd-cat", "-t", identifier, "-p", priority)
	cmd.Stdin = bytes.NewReader(jsonData)
	
	return cmd.Run()
}

func (h *JournaldHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &JournaldHandler{
		attrs: append(append([]slog.Attr{}, h.attrs...), attrs...),
		group: h.group,
	}
}

func (h *JournaldHandler) WithGroup(name string) slog.Handler {
	var newGroup string
	if h.group != "" {
		newGroup = h.group + "." + name
	} else {
		newGroup = name
	}
	return &JournaldHandler{
		attrs: append([]slog.Attr{}, h.attrs...),
		group: newGroup,
	}
}

func (h *JournaldHandler) getPriority(level slog.Level) string {
	// Map slog levels to systemd priority levels (0-7)
	// 0=emerg, 1=alert, 2=crit, 3=err, 4=warning, 5=notice, 6=info, 7=debug
	switch {
	case level >= LevelEmergency:
		return "0" // emergency
	case level >= LevelError:
		return "3" // error
	case level >= LevelWarn:
		return "4" // warning  
	case level >= LevelInfo:
		return "6" // info
	case level >= LevelDebug:
		return "7" // debug
	default:
		return "7" // trace and below -> debug
	}
}

func (h *JournaldHandler) getLevelName(level slog.Level) string {
	switch level {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelEmergency:
		return "EMERGENCY"
	default:
		return level.String()
	}
}

func (h *JournaldHandler) addAttrToMap(m map[string]interface{}, a slog.Attr, group string) {
	if a.Equal(slog.Attr{}) {
		return
	}
	
	key := a.Key
	if group != "" {
		key = group + "." + key
	}
	
	switch a.Value.Kind() {
	case slog.KindGroup:
		// Handle nested groups
		groupMap := make(map[string]interface{})
		for _, attr := range a.Value.Group() {
			h.addAttrToMap(groupMap, attr, "")
		}
		m[key] = groupMap
	default:
		m[key] = a.Value.Any()
	}
}