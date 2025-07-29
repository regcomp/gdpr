package main

import (
	"context"
	"crypto/tls"
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
	certs "github.com/regcomp/gdpr/auth/local_certs"
	"github.com/regcomp/gdpr/caching"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/constants"
	"github.com/regcomp/gdpr/logging"
	"github.com/regcomp/gdpr/routers"
	"github.com/regcomp/gdpr/secrets"
	servicecontext "github.com/regcomp/gdpr/service_context"
)

func run(
	ctx context.Context,
	getenv func(string) string,
	outStream io.Writer,
) error {
	if err := godotenv.Load(constants.LocalDefaultconfigPath); err != nil {
		log.Fatalf("error loading .env: %s", err.Error())
	}

	// pull in and parse relevant env variables for external secrets store
	// should be passed in by docker env variables at runtime
	configStore := config.NewConfigStore(getenv, getenv)

	// establish connection to secrets store
	secretStore, err := secrets.CreateSecretStore(configStore.GetSecretStoreConfig())
	if err != nil {
		return err
	}

	// establish connection to/instantiate cache
	// cache needs secret store to get missing information
	serviceCache, err := caching.CreateServiceCache(
		configStore.GetServiceCacheConfig(),
		secretStore.GetServiceCacheSecrets(),
	)
	if err != nil {
		// TODO:
	}

	// service context needs cache to pull neccessary data from
	stx, err := servicecontext.CreateServiceContext(serviceCache, configStore, secretStore)
	if err != nil {
		return err
	}

	// WARN: FOR DEBUGGING
	logging.NewRequestTracer(&config.RequestTracerConfig{
		TracerOn:         true,
		DisplayResponses: true,
	})

	router := routers.CreateRouter(stx)

	cert, err := tls.X509KeyPair(certs.ServerCertPEMBlock, certs.ServerKeyPEMBlock)
	if err != nil {
		// TODO:
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	server := &http.Server{
		Addr:      ":" + stx.ConfigStore.GetDefaultPort(),
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Fprintf(outStream, "listening on %s\n\n", server.Addr)
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
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
