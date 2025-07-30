# Victoria Logs Integration

This logging library now supports sending logs to Victoria Logs via HTTP API.

## Configuration

To enable Victoria Logs integration, simply set the `VICTORIA_LOGS_ENDPOINT` environment variable:

```bash
export VICTORIA_LOGS_ENDPOINT=http://oci-aca-001:9428/insert/elasticsearch/_bulk
```

When this environment variable is set, the library will automatically add a Victoria Logs handler during initialization with all features enabled by default:
- All log levels are accepted (TRACE through EMERGENCY)
- Source information is always included (file, line, function)
- Full structured logging support

## Features

- **Automatic detection**: Victoria Logs handler is added automatically when the environment variable is set
- **Elasticsearch bulk API format**: Uses the standard bulk API format for compatibility
- **Full structured logging**: All log attributes are preserved and sent to Victoria Logs
- **Source information**: Includes file, line, and function information when `AddSource: true`
- **Host information**: Automatically includes `host.name` field
- **Custom log levels**: Supports custom TRACE and EMERGENCY levels

## Log Format

Logs are sent in the following format:
```json
{
  "_msg": "log message",
  "_time": "2025-07-30T15:04:05.999999999Z",
  "level": "INFO",
  "host.name": "hostname",
  "source.file": "/path/to/file.go",
  "source.line": 42,
  "source.function": "main.main",
  "custom_field": "custom_value"
}
```

## Example Usage

```go
package main

import (
    "github.com/qwer-go/log"
)

func main() {
    // Logs will be sent to console, file, and Victoria Logs
    log.Info("Application started", "version", "1.0.0")
    
    // Structured logging with multiple attributes
    log.Info("User action", 
        "user_id", 12345,
        "action", "login",
        "ip", "192.168.1.100",
    )
    
    // Different log levels
    log.Trace("Detailed trace info")
    log.Debug("Debug message")
    log.Warn("Warning message")
    log.Error("Error occurred", "error", err)
    log.Emergency("Critical system failure!")
}
```

## Testing Victoria Logs

Run the example:
```bash
# Set the endpoint
export VICTORIA_LOGS_ENDPOINT=http://oci-aca-001:9428/insert/elasticsearch/_bulk

# Run the demo
go run example/victoria_logs_demo.go
```

Query logs from Victoria Logs:
```bash
# Get all logs
curl http://oci-aca-001:9428/select/logsql/query -d 'query=*'

# Filter by level
curl http://oci-aca-001:9428/select/logsql/query -d 'query=level:ERROR'

# Filter by host
curl http://oci-aca-001:9428/select/logsql/query -d 'query=host.name:yourhostname'
```

## Performance Considerations

- Logs are sent asynchronously to Victoria Logs using a buffered channel
- The handler uses a dedicated goroutine to process log entries without blocking
- Default buffer size is 1000 log entries
- If the buffer is full, new logs will be dropped (configurable behavior)
- Network errors are handled gracefully and won't crash your application

## Graceful Shutdown

The Victoria Logs handler runs asynchronously, so your application can exit normally. If you want to ensure the log file is properly closed:

```go
import "github.com/qwer-go/log"

func main() {
    // Your application code
    
    // Before exiting, close the log file
    defer log.Close()
}
```