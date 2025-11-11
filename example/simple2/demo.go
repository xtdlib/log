package main

import (
	"log/slog"

	"github.com/xtdlib/log"
)

func main() {
	log.Println("hello")
	slog.Error("hello")
	log.Fatalf("hello")
	log.Fatalf("")

	// err := fmt.Errorf("this is an error")
	// log.Fatal(err.Error())
	// slog.Error(err.Error())
}
