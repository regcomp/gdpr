// Auto-generated from /config/config.json - DO NOT EDIT

package config 

// Subrouter path prefixes
const (
	RouterApiPathPrefix = "/api"
	RouterAppPathPrefix = "/app"
	RouterAuthPathPrefix = "/auth"
	RouterBasePathPrefix = "/"
	RouterClientPathPrefix = "/client"
	RouterServicePathPrefix = "/"
	RouterStaticPathPrefix = "/static"
)

// endpoints
const (
	// api
	EndpointRecords = "/records"
	// auth
	EndpointLogin = "/login"
	EndpointLoginCallback = "/login-callback"
	EndpointRenewToken = "/renew-token"
	// base
	EndpointHealthz = "/healthz"
	EndpointRegisterServiceWorker = "/register-service-worker"
	// client
	EndpointDashboard = "/dashboard"
	EndpointLogout = "/logout"
	EndpointTest1 = "/test1"
	EndpointTest2 = "/test2"
)

// full paths
const (
	// api
	PathApiRecords = "/app/api/records"
	// auth
	PathAuthLogin = "/auth/login"
	PathAuthLoginCallback = "/auth/login-callback"
	PathAuthRenewToken = "/auth/renew-token"
	// base
	PathBaseHealthz = "/healthz"
	PathBaseRegisterServiceWorker = "/register-service-worker"
	// client
	PathClientDashboard = "/app/client/dashboard"
	PathClientLogout = "/app/client/logout"
	PathClientTest1 = "/app/client/test1"
	PathClientTest2 = "/app/client/test2"
)
// service workers
const (
	WorkerAuthRetryPath = "/static/sw/auth_retry.js"
	WorkerAuthRetryScope = "/"
)
// config keys
const (
	ConfigConfigStoreTypeKey = "CONFIG_STORE_TYPE"
	ConfigDefaultPortKey = "DEFAULT_PORT"
	ConfigSessionDurationKey = "SESSION_DURATION"
	ConfigSecretStoreTypeKey = "SECRET_STORE_TYPE"
	ConfigDatabaseProviderTableNamesKey = "DATABASE_PROVIDER_TABLE_NAMES"
	ConfigServiceUrlKey = "SERVICE_URL"
	ConfigServiceCacheTypeKey = "SERVICE_CACHE_TYPE"
	ConfigAuthProviderTypeKey = "AUTH_PROVIDER_TYPE"
	ConfigRecordsDatabaseTypeKey = "RECORDS_DATABASE_TYPE"
	ConfigDatabaseProviderNamesKey = "DATABASE_PROVIDER_NAMES"
	ConfigDatabaseProviderTypesKey = "DATABASE_PROVIDER_TYPES"
)

var ConfigAttrs = []string{
	ConfigConfigStoreTypeKey,
	ConfigDefaultPortKey,
	ConfigSessionDurationKey,
	ConfigSecretStoreTypeKey,
	ConfigDatabaseProviderTableNamesKey,
	ConfigServiceUrlKey,
	ConfigServiceCacheTypeKey,
	ConfigAuthProviderTypeKey,
	ConfigRecordsDatabaseTypeKey,
	ConfigDatabaseProviderNamesKey,
	ConfigDatabaseProviderTypesKey,
}
// Cookies
const (
	CookieNameAccessToken = "access-token"
	CookieNameRefreshToken = "refresh-token"
	CookieNameSessionId = "session-id"
)
// form values
const FormValueNonce = "nonce"
// local files
const (
	LocalDefaultConfigPath = "config/default.config"
	LocalEnvPath = ".env"
)
// query parameters
const (
	QueryParamWorkerScope = "worker-scope"
	QueryParamLimit = "limit"
	QueryParamAfter = "after"
	QueryParamRedirectUrl = "redirect-url"
	QueryParamRequestId = "request-id"
	QueryParamAccessToken = "access-token"
	QueryParamRefreshToken = "refresh-token"
	QueryParamWorkerPath = "worker-path"
)
// request context keys
const (
	ContextKeyNonceToken = "nonce-token"
	ContextKeyClaims = "claims"
	ContextKeySessionId = "session-id"
)
// headers
const (
	HeaderNonceToken = "Nonce-Token"
	HeaderRenewAccessToken = "Renew-Access-Token"
	HeaderServiceWorkerAllowed = "Service-Worker-Allowed"
	HeaderAuthRetryWorkerRunning = "Auth-Retry-Worker-Running"
)
// values
const (
	ValueNameDelim = ":"
	ValueItemDelim = ","
	ValueTrue = "true"
	ValueLocalType = "LOCAL"
	ValueEntryDelim = ";"
)

