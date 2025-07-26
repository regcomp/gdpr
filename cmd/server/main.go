//go:generate go run ../../constants/generators/gen_shared.go

package main

import (
	"context"
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
		// TODO: log the server crash
	}
	// println("tada")
}
