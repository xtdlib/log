package main

import (
	"log/slog"

	"github.com/xtdlib/log"
)

func main() {
	log.Debug().Msg("debug message")
	log.Info().Int("a", 3).Int("b", 4).Msg("hello world james")
	log.Critical().Int("werwer", 3).Msg("werwerwer")
	log.Error().Int("werwer", 3).Msg("werwerwer")
	log.Notice().Int("werwer", 3).Msg("werwerwer")
	log.Critical().Int("werwer", 3).Msg("critical")

	// log.Print("hello world")

	slog.Info("hello world")
}
