package servicecontext

import (
	"os"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/caching"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/database"
	"github.com/regcomp/gdpr/logging"
)

var STX *ServiceContext

type ServiceContext struct {
	AuthProvider  auth.IAuthProvider
	SessionStore  auth.ISessionStore
	NonceStore    *auth.NonceStore
	CookieManager *auth.CookieManager

	ConfigStore config.IConfigStore

	DatabaseStore database.IDatabaseStore
	RequestStore  caching.IRequestStore

	RequestLogger logging.ILogger
}

func CreateServiceContext(serviceCache caching.IServiceCache, configStore config.IConfigStore) (*ServiceContext, error) {
	// other context setup goes here, like getting certs/keys
	authProvider, err := auth.GetProvider(configStore.GetAuthProvider())
	if err != nil {
		return nil, err
	}

	requestlogger := logging.NewRequestLogger(os.Stdout)
	logging.NewTracer(configStore.GetTracerLevel())

	databaseStore, err := database.CreateDatabaseStore(getConfig)
	if err != nil {
		// TODO:
		// NOTE: a databaseProvider can fail to initialize. This should halt the service from running
		// NOTE: Shouldn't need to establish database connections until necessary
	}

	cookieManager := auth.CreateCookieManager(serviceCache)
	sessionStore := auth.CreateSessionStore(serviceCache)
	nonceStore := auth.CreateNonceStore(serviceCache)
	requestStore := caching.CreateRequestStore(serviceCache)

	return &ServiceContext{
		AuthProvider:  authProvider,
		RequestLogger: requestlogger,
		CookieManager: cookieManager,
		SessionStore:  sessionStore,
		DatabaseStore: databaseStore,
		RequestStore:  requestStore,
		NonceStore:    nonceStore,
		ConfigStore:   configStore,
	}, nil
}
