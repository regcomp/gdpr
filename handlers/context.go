package handlers

import (
	"log"
	"log/slog"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/database"
	"github.com/regcomp/gdpr/logging"
)

var STX *ServiceContext

type ServiceContext struct {
	AuthProvider  auth.IProvider
	RequestLogger *slog.Logger
	RequestTracer logging.ITracer
	CookieKeys    *securecookie.SecureCookie
	SessionStore  *auth.SessionStore
	DatabaseStore *database.DatabaseStore

	HostPath        string
	SessionDuration int

	HMACSecret []byte
}

func CreateServiceContext(getenv func(string) string) *ServiceContext {
	// other context setup goes here, like getting certs/keys
	authProvider, err := auth.GetProvider(getenv)
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	requestlogger := logging.NewRequestLogger(os.Stdout)
	requestTracer := logging.NewTracer(getenv)

	databaseStore, err := database.CreateDatabaseStore(getenv)
	if err != nil {
		// TODO:
		// NOTE: a databaseProvider can fail to initialize. This should halt the service from running
		// NOTE: Shouldn't need to establish database connections until necessary
	}

	cookieKeys := auth.CreateCookieKeys()
	sessionStore := auth.CreateSessionStore()

	hmacSecret := auth.GenerateHMACSecret()

	return &ServiceContext{
		AuthProvider:  authProvider,
		RequestLogger: requestlogger,
		RequestTracer: requestTracer,
		CookieKeys:    cookieKeys,
		SessionStore:  sessionStore,
		DatabaseStore: databaseStore,
		HostPath:      "localhost:8080",
		HMACSecret:    hmacSecret,
	}
}

func LinkServiceContext(stx *ServiceContext) {
	STX = stx
}
