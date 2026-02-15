package main

import (
	stdlog "log"
	"log/slog"
	"os"

	"github.com/xtdlib/log/example/simple/pkg"
)

func XX() (int, error) {
	return 0, nil
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))
	slog.Info("wer")
	stdlog.SetFlags(stdlog.Lshortfile | stdlog.LstdFlags)
	stdlog.Println(XX())
	pkg.Foo()
}
