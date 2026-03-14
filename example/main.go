package main

import (
	"log/slog"

	"github.com/xtdlib/log"
)

func main() {
	log.Info().Int("a", 3).Int("b", 4).Msg("hello world james")
	log.Emerg().Int("werwer", 3).Msg("werwerwer")
	log.Error().Int("werwer", 3).Msg("werwerwer")
	log.Notice().Int("werwer", 3).Msg("werwerwer")

	log.Print("hello world")

	slog.Info("hello world")
}
