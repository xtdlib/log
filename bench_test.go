package log

import (
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	phuslog "github.com/phuslu/log"
)

func BenchmarkXtdlog(b *testing.B) {
	f, err := os.Create("/home/rok/src/github.com/xtdlib/log/example/1")
	if err != nil {
		b.Fatal(err)
	}
	SetWriter(f)
	for b.Loop() {
		Info().Int("a", 3).Int("b", 4).Msg("hello world james")
	}
}

func BenchmarkSlog(b *testing.B) {
	f, err := os.Create("/home/rok/src/github.com/xtdlib/log/example/2")
	if err != nil {
		b.Fatal(err)
	}
	slogger := slog.New(slog.NewJSONHandler(f, nil))
	for b.Loop() {
		slogger.Info("hello world james", "a", 3, "b", 4)
	}
}

func BenchmarkPatchedSlog(b *testing.B) {
	f, err := os.Create("/home/rok/src/github.com/xtdlib/log/example/3")
	if err != nil {
		b.Fatal(err)
	}
	_default = phuslog.Logger{
		// TimeFormat: "01-02 15:04:05",
		// TimeFormat: time.DateTime,
		TimeFormat: time.RFC3339Nano,
		// TimeFormat: phuslog.TimeFormatUnixMs,
		Writer: phuslog.IOWriter{Writer: f},
		Level:  phuslog.DebugLevel,
		// Caller: 2,
	}
	slog.SetDefault(_default.Slog())
	for b.Loop() {
		slog.Info("hello world james", "a", 3, "b", 4)
	}
}

func TestMain(m *testing.M) {
	_defaultOutput = io.Discard
	os.Exit(m.Run())
}

func BenchmarkNoCaller(b *testing.B) {
	for b.Loop() {
		Info().Int("a", 3).Int("b", 4).Msg("hello world james")
	}
}

func BenchmarkCaller(b *testing.B) {
	for b.Loop() {
		Notice().Int("a", 3).Int("b", 4).Msg("hello world james")
	}
}
