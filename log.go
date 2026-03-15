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

// slog-rs
// https://github.com/slog-rs/slog/blob/1adf6422ca472ce29b1e48c99142eca2f3193d39/src/lib.rs#L2199
//     ["OFF", "CRIT", "ERRO", "WARN", "INFO", "DEBG", "TRCE"];

// ● Fetch(https://raw.githubusercontent.com/golang/glog/master/internal/logsink/logsink.go)
//   ⎿  Received 11.9KB (200 OK)
//
// ● logsink.go line 213-214에서 확인됩니다:
//
//   const severityChar = "IWEF"
//   buf.WriteByte(severityChar[m.Severity])
//
//   로그 헤더를 조립할 때 I, W, E, F 한 글자를 prefix로 씁니다. 실제 glog 출력이 이렇게 생겼습니다:
//
//   I0314 10:22:05.123456 12345 main.go:42] hello world
//   W0314 10:22:05.123457 12345 main.go:43] something bad

func init() {
	phuslog.TimeKey = "ts"
	phuslog.CallerKey = "src"
	phuslog.CallerFuncKey = "func"
	phuslog.MessageKey = "msg"
	phuslog.LevelString = [8]string{
		phuslog.TraceLevel: "TRAC",
		phuslog.DebugLevel: "DEBG",
		phuslog.InfoLevel:  "INFO",
		phuslog.WarnLevel:  "NOTI",
		phuslog.ErrorLevel: "ERRO",
		phuslog.FatalLevel: "FATAL",
		phuslog.PanicLevel: "PANIC",
	}

	_default = phuslog.Logger{
		// TimeFormat: "01-02 15:04:05",
		// TimeFormat: time.DateTime,
		// TimeFormat: time.RFC3339Nano,
		TimeFormat: phuslog.TimeFormatUnixMs,
		// Writer:     phuslog.IOWriter{Writer: _defaultOutput},
		Writer: &phuslog.ConsoleWriter{
			Formatter: phuslog.LogfmtFormatter{"ts"}.Formatter,
			Writer:    io.MultiWriter(os.Stdout, os.Stderr),
		},

		// Writer: &phuslog.ConsoleWriter{
		// 	Writer:         os.Stdout,
		// 	ColorOutput:    true,
		// 	QuoteString:    true,
		// 	EndWithMessage: true,
		// },
		Level: phuslog.DebugLevel,
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
	return _default.Log().Str("level", "TRAC")
}

func Debug() (e *phuslog.Entry) {
	return _default.Log().Str("level", "DEBG")
}

func Info() (e *phuslog.Entry) {
	return _default.Log().Str("level", "INFO")
}

func Notice() (e *phuslog.Entry) {
	return _default.Log().Str("level", "NOTI")
}

// ["OFF", "CRIT", "ERRO", "WARN", "INFO", "DEBG", "TRCE"];
func Error() (e *phuslog.Entry) {
	return _default.Log().Str("level", "ERRO").Caller(2)
}

func Critical() (e *phuslog.Entry) {
	return _default.Log().Str("level", "FATL").Caller(2)
}

func Print(args ...any) {
	_default.Log().Str("level", "INFO").Msgs(args...)
}
