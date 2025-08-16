package main

import (
	"context"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	err := run(
		ctx,
		os.Getenv,
		os.Stderr,
	)
	if err != nil {
		log.Printf("Server Crash: %s", err.Error())
		// TODO: log the server crash
	}
	// println("tada")
}
