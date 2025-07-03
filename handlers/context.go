package handlers

import (
	"log"
	"log/slog"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/logging"
)

var STX *ServiceContext

type ServiceContext struct {
	AuthProvider  auth.Provider
	RequestLogger *slog.Logger
	RequestTracer logging.Tracer
	CookieKeys    *securecookie.SecureCookie
	SessionStore  *auth.SessionStore

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

	cookieKeys := auth.CreateCookieKeys()
	sessionStore := auth.CreateSessionStore()

	return &ServiceContext{
		AuthProvider:  authProvider,
		RequestLogger: requestlogger,
		RequestTracer: requestTracer,
		CookieKeys:    cookieKeys,
		SessionStore:  sessionStore,
		HostPath:      "localhost:8080",
	}
}

func LinkServiceContext(stx *ServiceContext) {
	STX = stx
}
