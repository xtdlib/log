package main

import (
	"os"

	"github.com/phuslu/log"
)

func main() {
	logger := log.Logger{
		// TimeFormat: "01-02 15:04:05",
		// TimeFormat: time.DateTime,
		// TimeFormat: time.RFC3339Nano,
		// TimeFormat: phuslog.TimeFormatUnixMs,
		// Writer:     phuslog.IOWriter{Writer: _defaultOutput},

		Writer: &log.ConsoleWriter{
			Writer: os.Stdout,
			ColorOutput: true,
			QuoteString: true,
			EndWithMessage:  true,
		},
		Level: log.DebugLevel,
		// Caller: 2,
	}
	logger.Info().Int("a", 3).Int("b", 4).Msg("hello world james")
}
