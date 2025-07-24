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
	"github.com/regcomp/gdpr/cache"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/routers"
	"github.com/regcomp/gdpr/secrets"
	servicecontext "github.com/regcomp/gdpr/service_context"
)

const (
	envPath    = ".env"
	configPath = "config/default.config"
)

func run(
	ctx context.Context,
	getenv func(string) string,
	inStream io.Reader,
	outStream io.Writer,
) error {
	// Loads files in parameter order
	if err := godotenv.Load(configPath, envPath); err != nil {
		log.Fatalf("error loading .env: %s", err.Error())
	}

	secretStoreType := getenv(config.SecretStoreTypeKey)
	secretStoreConfig := secrets.LoadConfig(secretStoreType)

	// establish connection to secrets store, type configured in env variable passed to docker container
	secretStore, err := secrets.CreateSecretStore(secretStoreConfig)
	if err != nil {
		return err
	}

	// establish connection to/instantiate cache
	// cache needs secret store to get missing information
	serviceCacheType := getenv(config.ServiceCacheTypeKey)
	serviceCache, err := cache.CreateServiceCache(secretStore, serviceCacheType)

	// service context needs cache to pull neccessary data from
	stx, err := servicecontext.CreateServiceContext(serviceCache, getenv)
	if err != nil {
		return err
	}

	router := routers.CreateRouter(stx)

	cert, err := tls.X509KeyPair(certs.ServerCertPEMBlock, certs.ServerKeyPEMBlock)
	if err != nil {
		// TODO:
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	config := config.LoadConfig(getenv)
	server := &http.Server{
		Addr:      ":" + config.Port,
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
