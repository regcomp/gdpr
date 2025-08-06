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

var STX *ServiceContext

type ServiceContext struct {
	AuthProvider  auth.IAuthProvider
	SessionStore  auth.ISessionStore
	NonceStore    *auth.NonceStore
	CookieManager *auth.CookieManager

	ConfigStore config.IConfigStore

	DatabaseStore *database.DatabaseManager
	RequestStore  caching.IRequestStore

	RequestLogger logging.ILogger
}

func CreateServiceContext(
	serviceCache caching.IServiceCache,
	configStore config.IConfigStore,
	secretStore secrets.ISecretStore,
) (*ServiceContext, error) {
	// other context setup goes here, like getting certs/keys
	authProvider, err := auth.CreateAuthProvider(
		configStore.GetAuthProviderConfig(),
		secretStore.GetAuthProviderSecrets(),
	)
	if err != nil {
		return nil, err
	}

	requestlogger := logging.NewRequestLogger(os.Stdout)

	databaseManagerConfig, err := configStore.GetDatabaseManagerConfig()
	if err != nil {
		return nil, err
	}

	databaseManager, err := database.CreateDatabaseManager(
		databaseManagerConfig,
		secretStore.GetDatabaseStoreSecrets(),
	)
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
		DatabaseStore: databaseManager,
		RequestStore:  requestStore,
		NonceStore:    nonceStore,
		ConfigStore:   configStore,
	}, nil
}
