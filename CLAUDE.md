# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go logging library (`github.com/xtdlib/log`) that extends Go's built-in `log/slog` package with enhanced features:

- **Colorful console output** with custom formatting
- **Triple output**: Console with colors + JSON file logging + Victoria Logs via HTTP API
- **Extended log levels**: TRACE (-8) and EMERGENCY (12) in addition to standard slog levels
- **Drop-in replacement** for `log/slog` with same API
- **Custom handlers**: `consoleHandler` for colored output, `multiHandler` for multiple destinations, `VictoriaLogsHandler` for remote logging
- **Victoria Logs integration**: Automatic HTTP-based logging to Victoria Logs (defaults to `oci-aca-001:9428`)

## Development Commands

### Testing
```bash
# Run all tests
go test .

# Run tests with verbose output
go test -v .

# Run specific test
go test -run TestBasicLogging

# Run Victoria Logs tests only
go test -v -run TestVictoriaLogs

# Run tests with coverage
go test -cover .
```

### Building
```bash
# Build the module (verify compilation)
go build .

# Run examples
go run example/simple/demo.go
go run example/levels.go  
go run example/level/emergency.go
go run example/victoria_logs_demo.go
```

### Module Management
```bash
# Tidy dependencies
go mod tidy

# Verify module
go mod verify
```

## Architecture

### Handler System

The library uses a **multi-handler architecture** where logs are simultaneously sent to multiple destinations:

1. **Console Handler** (`consoleHandler`): Provides ANSI-colored console output with structured formatting
2. **File Handler** (JSON): Writes structured JSON logs to timestamped files
3. **Victoria Logs Handler** (`VictoriaLogsHandler`): Asynchronously sends logs to Victoria Logs via HTTP API

### Core Components

- **`init.go`**: Initialization logic that automatically sets up all handlers based on environment
- **`log.go`**: Core implementation with custom handlers, ANSI colors, and logging functions
- **`victoria_logs.go`**: Optimized async HTTP handler with buffer pooling and pre-calculated constants
- **`multiHandler`**: Multiplexes logging to multiple handlers simultaneously

### Automatic Configuration

The library automatically configures itself during `init()`:

1. **Log Directory Selection**: 
   - Primary: `/var/log/{appname}/` (standard Linux location)
   - Fallback: `/tmp/{appname}/logs/` (when no `/var/log` access)

2. **Victoria Logs**: Enabled when `VICTORIA_LOGS_ENDPOINT` environment variable is set
   - Default endpoint: `http://oci-aca-001:9428/insert/elasticsearch/_bulk`
   - Uses Elasticsearch bulk API format
   - Async processing with 2000-entry buffer and connection pooling

3. **File Naming**: Timestamped files using format `YYYYMMDDTHHMMSS.log`

### Performance Optimizations

The Victoria Logs handler includes several performance optimizations:
- **Buffer pooling** (`sync.Pool`) for byte buffers
- **Pre-calculated constants** (create line, level names map)
- **Pre-sized maps** to reduce allocations
- **Async processing** with buffered channels
- **HTTP connection reuse** with optimized client

### Testing Architecture

- **In-memory buffers** for capturing and verifying log output
- **Mock HTTP servers** for testing Victoria Logs integration
- **Compatibility tests** ensuring drop-in replacement for `slog`
- **Async handling verification** with proper timing

## Usage Patterns

### Basic Usage (Drop-in slog replacement)
```go
import "github.com/xtdlib/log"

log.Info("message", "key", "value")  // Instead of slog.Info(...)
log.Emergency("critical failure", "component", "database")
```

### Victoria Logs Integration
Set environment variable (no code changes needed):
```bash
export VICTORIA_LOGS_ENDPOINT=http://oci-aca-001:9428/insert/elasticsearch/_bulk
```

### Graceful Shutdown
```go
defer log.Close()  // Ensures log file is closed properly
```

## Testing Notes

- Tests use mock HTTP servers to verify Victoria Logs integration
- Buffer timing of 50ms is used in async tests to ensure proper message delivery
- All handlers are tested independently and in combination via `multiHandler`
- Tests verify ANSI color output, JSON structure, and Victoria Logs bulk API format