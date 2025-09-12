package main

import "github.com/xtdlib/log"

func main() {
	log.Fatalf("hello")

	log.Println("hello")
	log.Fatalf("")

	// err := fmt.Errorf("this is an error")
	// log.Fatal(err.Error())
	// slog.Error(err.Error())
}
