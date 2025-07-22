package servicecontext

import (
	"log"
	"os"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/database"
	"github.com/regcomp/gdpr/logging"
)

var STX *ServiceContext

type ServiceContext struct {
	AuthProvider auth.IAuthProvider
	CookieCrypt  auth.ICookieCrypt
	SessionStore auth.ISessionStore
	NonceStore   auth.INonceStore

	ConfigStore config.IConfigStore

	DatabaseStore database.IDatabaseStore
	RequestStore  database.IRequestStore

	RequestLogger logging.ILogger
	RequestTracer logging.IRequestTracer
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

	cookieCrypt := auth.CreateCookieCrypt()

	sessionStore := auth.CreateSessionStore()
	requestStore := database.CreateRequestStore()
	nonceStore := auth.CreateNonceStore()
	configStore := config.NewConfigStore()

	return &ServiceContext{
		AuthProvider:  authProvider,
		RequestLogger: requestlogger,
		RequestTracer: requestTracer,
		CookieCrypt:   cookieCrypt,
		SessionStore:  sessionStore,
		DatabaseStore: databaseStore,
		RequestStore:  requestStore,
		NonceStore:    nonceStore,
		ConfigStore:   configStore,
	}
}
