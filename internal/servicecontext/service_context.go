package servicecontext

import (
	"os"

	"github.com/regcomp/gdpr/internal/auth"
	"github.com/regcomp/gdpr/internal/caching"
	"github.com/regcomp/gdpr/internal/config"
	"github.com/regcomp/gdpr/internal/database"
	"github.com/regcomp/gdpr/internal/logging"
	"github.com/regcomp/gdpr/internal/secrets"
)

type ServiceContext struct {
	AuthProvider auth.IAuthProvider
	ConfigStore  config.IConfigStore

	SessionManager  *caching.SessionManager
	NonceManager    *caching.NonceManager
	CookieManager   *caching.CookieManager
	RequestManager  *caching.RequestManager
	DatabaseManager *database.DatabaseManager

	RequestLogger logging.ILogger
}

func CreateServiceContext(
	serviceCache caching.IServiceCache,
	configStore config.IConfigStore,
	secretManager secrets.ISecretStore,
) (*ServiceContext, error) {
	// other context setup goes here, like getting certs/keys
	authConfig, err := configStore.GetAuthProviderConfig()
	if err != nil {
		return nil, err
	}
	authSecrets, err := secretManager.GetAuthProviderSecrets()
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
	databaseManagerSecrets, err := secretManager.GetDatabaseManagerSecrets()
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

	cookieManager := caching.CreateCookieManager(serviceCache)
	sessionStore := caching.CreateSessionManager(serviceCache)
	nonceStash := caching.CreateNonceManager(serviceCache)
	requestStash := caching.CreateRequestManager(serviceCache)

	return &ServiceContext{
		AuthProvider:    authProvider,
		RequestLogger:   requestlogger,
		CookieManager:   cookieManager,
		SessionManager:  sessionStore,
		DatabaseManager: databaseManager,
		RequestManager:  requestStash,
		NonceManager:    nonceStash,
		ConfigStore:     configStore,
	}, nil
}
