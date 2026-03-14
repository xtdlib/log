package log

import (
	"io"
	"log/slog"
	"os"

	stdlog "log"

	phuslog "github.com/phuslu/log"
)

// "emerg" (0), "alert" (1), "crit" (2), "err" (3), "warning" (4), "notice" (5), "info" (6), "debug" (7).

var _default phuslog.Logger
var _defaultOutput io.Writer = os.Stdout

func init() {
	phuslog.TimeKey = "ts"
	phuslog.CallerKey = "src"
	phuslog.CallerFuncKey = "func"
	phuslog.MessageKey = "msg"
	phuslog.LevelString = [8]string{
		phuslog.TraceLevel: "TRACE",
		phuslog.DebugLevel: "DEBUG",
		phuslog.InfoLevel:  "INFO",
		phuslog.WarnLevel:  "NOTI",
		phuslog.ErrorLevel: "ERROR",
		phuslog.FatalLevel: "FATAL",
		phuslog.PanicLevel: "PANIC",
	}

	_default = phuslog.Logger{
		// TimeFormat: "01-02 15:04:05",
		// TimeFormat: time.DateTime,
		// TimeFormat: time.RFC3339Nano,
		TimeFormat: phuslog.TimeFormatUnixMs,
		Writer:     phuslog.IOWriter{Writer: _defaultOutput},
		Level:      phuslog.DebugLevel,
		// Caller: 2,
	}

	slog.SetDefault(slog.New(_default.Slog().Handler()))
}

func SetWriter(w io.Writer) {
	_default.Writer = phuslog.IOWriter{Writer: w}
}

func WithCaller(n int) {
	_default.Caller = n
}

var Println = stdlog.Println
var Printf = _default.Printf

func Trace() (e *phuslog.Entry) {
	return _default.WithLevel(1)
}

func Debug() (e *phuslog.Entry) {
	return _default.WithLevel(2)
}

func Info() (e *phuslog.Entry) {
	return _default.WithLevel(3)
}

func Notice() (e *phuslog.Entry) {
	return _default.WithLevel(4).Caller(2)
}

func Error() (e *phuslog.Entry) {
	return _default.WithLevel(5).Caller(2)
}

func Emerg() (e *phuslog.Entry) {
	return _default.Log().Str("level", "EMERG").Caller(2)
}

func Print(args ...any) {
	_default.Log().Str("level", "INFO").Msgs(args...)
}
