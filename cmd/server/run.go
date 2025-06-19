package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/handlers"
	"github.com/regcomp/gdpr/routers"
)

const (
	envPath    = ".env"
	configPath = "config/default.config"
)

func run(
	ctx context.Context,
	// args []string, // TODO: add later as further configuration
	getenv func(string) string,
	inStream io.Reader,
	outStream io.Writer,
) error {
	// Loads files in parameter order
	if err := godotenv.Load(configPath, envPath); err != nil {
		log.Fatalf("error loading .env: %s", err.Error())
	}

	// TODO: Load args

	// TODO: if --config-path=... exists, load that into env last

	config := config.LoadConfig(getenv)

	stx := handlers.CreateServiceContext(getenv)
	handlers.LinkServiceContext(stx)

	router := routers.CreateRouter(
		routers.CreateApiRouter(),
		routers.CreateClientRouter(),
	)

	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Fprintf(outStream, "listening on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(outStream, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(outStream, "error shutting down http server: %v\n", err)
		}
	}()

	wg.Wait()
	return nil
}
