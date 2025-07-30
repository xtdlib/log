package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)


// VictoriaLogsHandler sends logs to Victoria Logs via HTTP API
type VictoriaLogsHandler struct {
	endpoint string
	client   *http.Client
	attrs    []slog.Attr
	group    string
	logChan  chan []byte
}

var (
	victoriaLogsHandler *VictoriaLogsHandler
	// Buffer pool for reusing byte buffers
	bufferPool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
	// Pre-calculated constants
	createLineBytes = []byte(`{"create":{}}` + "\n")
	// Level name lookup map for better performance
	levelNames = map[slog.Level]string{
		LevelTrace:     "TRACE",
		LevelDebug:     "DEBUG",
		LevelInfo:      "INFO",
		LevelWarn:      "WARN",
		LevelError:     "ERROR",
		LevelEmergency: "EMERGENCY",
	}
	hostname, _ = os.Hostname()
)

// NewVictoriaLogsHandler creates a new handler that sends logs to Victoria Logs
func NewVictoriaLogsHandler(endpoint string) *VictoriaLogsHandler {
	if endpoint == "" {
		endpoint = "http://oci-aca-001:9428/insert/elasticsearch/_bulk"
	}

	h := &VictoriaLogsHandler{
		endpoint: endpoint,
		client:   httpClient,
		logChan:  make(chan []byte, 2000), // Buffer up to 2000 log entries
	}

	// Start the async worker
	go h.worker()

	// Store handler reference for cleanup
	victoriaLogsHandler = h

	return h
}

// worker processes log entries asynchronously
func (h *VictoriaLogsHandler) worker() {
	// Pre-create header for reuse
	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	
	for data := range h.logChan {
		req, err := http.NewRequest("POST", h.endpoint, bytes.NewReader(data))
		if err != nil {
			continue
		}

		req.Header = header

		resp, err := h.client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()
	}
}

func (h *VictoriaLogsHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Always enabled for all levels
	return true
}


func (h *VictoriaLogsHandler) Handle(ctx context.Context, r slog.Record) error {
	// Pre-size map with expected number of fields (7-10 typically)
	entry := make(map[string]interface{}, 10)

	// Standard fields
	entry["_msg"] = r.Message
	entry["_time"] = r.Time.Format(time.RFC3339Nano)
	entry["level"] = getLevelName(r.Level)
	entry["host"] = hostname
	entry["app"] = appName

	// Always add source information
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		entry["source.file"] = f.File
		entry["source.line"] = f.Line
		entry["source.function"] = f.Function
	}

	// Add prepended attributes
	for _, a := range h.attrs {
		addAttrToMap(entry, a, h.group)
	}

	// Add record attributes
	r.Attrs(func(a slog.Attr) bool {
		addAttrToMap(entry, a, h.group)
		return true
	})

	// Get buffer from pool
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	// Write create line (pre-calculated bytes)
	buf.Write(createLineBytes)

	// Write log entry
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(entry); err != nil {
		return fmt.Errorf("failed to encode log entry: %w", err)
	}

	// Copy bytes before returning buffer to pool
	data := make([]byte, buf.Len())
	copy(data, buf.Bytes())

	// Send to channel for async processing
	select {
	case h.logChan <- data:
		// Successfully queued
	default:
		// Channel is full, drop the log
	}

	return nil
}

func (h *VictoriaLogsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Share the same channel and worker for derived handlers
	return &VictoriaLogsHandler{
		endpoint: h.endpoint,
		client:   h.client,
		attrs:    append(append([]slog.Attr{}, h.attrs...), attrs...),
		group:    h.group,
		logChan:  h.logChan,
	}
}

func (h *VictoriaLogsHandler) WithGroup(name string) slog.Handler {
	var newGroup string
	if h.group != "" {
		newGroup = h.group + "." + name
	} else {
		newGroup = name
	}
	// Share the same channel and worker for derived handlers
	return &VictoriaLogsHandler{
		endpoint: h.endpoint,
		client:   h.client,
		attrs:    append([]slog.Attr{}, h.attrs...),
		group:    newGroup,
		logChan:  h.logChan,
	}
}

func getLevelName(level slog.Level) string {
	if name, ok := levelNames[level]; ok {
		return name
	}
	return level.String()
}

func addAttrToMap(m map[string]interface{}, a slog.Attr, group string) {
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
			addAttrToMap(groupMap, attr, "")
		}
		m[key] = groupMap
	default:
		m[key] = a.Value.Any()
	}
}
