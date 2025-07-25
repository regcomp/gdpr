//go:generate go run ../../constants/generators/shared.go

package main

import (
	"context"
	"os"
)

func main() {
	ctx := context.Background()
	err := run(ctx, os.Getenv, os.Stdin, os.Stderr)
	if err != nil {
		// TODO: log the server crash
	}
	// println("tada")
}
