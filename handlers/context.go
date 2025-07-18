package handlers

import (
	"log"
	"log/slog"
	"os"

	sc "github.com/gorilla/securecookie"
	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/database"
	"github.com/regcomp/gdpr/logging"
)

var STX *ServiceContext

type ServiceContext struct {
	AuthProvider  auth.IAuthProvider
	CookieKeys    *sc.SecureCookie
	SessionStore  auth.ISessionStore
	DatabaseStore database.IDatabaseStore
	RequestStore  IRequestStore
	NonceStore    auth.INonceStore

	RequestLogger *slog.Logger
	RequestTracer logging.IRequestTracer

	HostPath        string
	SessionDuration int
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
	requestStore := CreateRequestStore()
	nonceStore := auth.CreateNonceStore()

	return &ServiceContext{
		AuthProvider:  authProvider,
		RequestLogger: requestlogger,
		RequestTracer: requestTracer,
		CookieKeys:    cookieKeys,
		SessionStore:  sessionStore,
		DatabaseStore: databaseStore,
		RequestStore:  requestStore,
		NonceStore:    nonceStore,
		HostPath:      "localhost:8080",
	}
}

func LinkServiceContext(stx *ServiceContext) {
	STX = stx
}
