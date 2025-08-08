package servicecontext

import (
	"os"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/caching"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/database"
	"github.com/regcomp/gdpr/logging"
	"github.com/regcomp/gdpr/secrets"
)

type ServiceContext struct {
	AuthProvider  auth.IAuthProvider
	SessionStore  *auth.SessionStore
	NonceManager  *auth.NonceManager
	CookieManager *auth.CookieManager

	ConfigStore config.IConfigStore

	DatabaseManager *database.DatabaseManager
	RequestStash    *caching.RequestStash

	RequestLogger logging.ILogger
}

func CreateServiceContext(
	serviceCache caching.IServiceCache,
	configStore config.IConfigStore,
	secretStore secrets.ISecretStore,
) (*ServiceContext, error) {
	// other context setup goes here, like getting certs/keys
	authConfig, err := configStore.GetAuthProviderConfig()
	if err != nil {
		return nil, err
	}
	authSecrets, err := secretStore.GetAuthProviderSecrets()
	if err != nil {
		return nil, err
	}

	authProvider, err := auth.CreateAuthProvider(authConfig, authSecrets)
	if err != nil {
		return nil, err
	}

	requestlogger := logging.NewRequestLogger(os.Stdout)

	databaseManagerConfig, err := configStore.GetDatabaseManagerConfig()
	if err != nil {
		return nil, err
	}
	databaseManagerSecrets, err := secretStore.GetDatabaseManagerSecrets()
	if err != nil {
		return nil, err
	}
	databaseManager, err := database.CreateDatabaseManager(databaseManagerConfig, databaseManagerSecrets)
	if err != nil {
		return nil, err

		// TODO:
		// NOTE: a databaseProvider can fail to initialize. This should halt the service from running
		// NOTE: Shouldn't need to establish database connections until necessary
	}

	cookieManager := auth.CreateCookieManager(serviceCache)
	sessionStore := auth.CreateSessionStore(serviceCache)
	nonceStash := auth.CreateNonceStash(serviceCache)
	requestStash := caching.CreateRequestStash(serviceCache)

	return &ServiceContext{
		AuthProvider:    authProvider,
		RequestLogger:   requestlogger,
		CookieManager:   cookieManager,
		SessionStore:    sessionStore,
		DatabaseManager: databaseManager,
		RequestStash:    requestStash,
		NonceManager:    nonceStash,
		ConfigStore:     configStore,
	}, nil
}
