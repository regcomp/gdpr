// Auto-generated from /config/config.json - DO NOT EDIT

package config 

// Subrouter path prefixes
const (
	RouterApiPathPrefix = "/api"
	RouterAuthPathPrefix = "/auth"
	RouterBasePathPrefix = "/"
	RouterClientPathPrefix = "/"
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
	PathApiRecords = "/api/records"
	// auth
	PathAuthLogin = "/auth/login"
	PathAuthLoginCallback = "/auth/login-callback"
	PathAuthRenewToken = "/auth/renew-token"
	// base
	PathBaseHealthz = "/healthz"
	PathBaseRegisterServiceWorker = "/register-service-worker"
	// client
	PathClientDashboard = "/dashboard"
	PathClientLogout = "/logout"
	PathClientTest1 = "/test1"
	PathClientTest2 = "/test2"
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
	ConfigSecretStoreTypeKey = "SECRET_STORE_TYPE"
	ConfigServiceUrlKey = "SERVICE_URL"
	ConfigSessionDurationKey = "SESSION_DURATION"
	ConfigServiceCacheTypeKey = "SERVICE_CACHE_TYPE"
	ConfigAuthProviderTypeKey = "AUTH_PROVIDER_TYPE"
	ConfigRecordsDatabaseTypeKey = "RECORDS_DATABASE_TYPE"
	ConfigDatabaseProviderNamesKey = "DATABASE_PROVIDER_NAMES"
	ConfigDatabaseProviderTypesKey = "DATABASE_PROVIDER_TYPES"
	ConfigDatabaseProviderTableNamesKey = "DATABASE_PROVIDER_TABLE_NAMES"
)

var ConfigAttrs = []string{
	ConfigConfigStoreTypeKey,
	ConfigDefaultPortKey,
	ConfigSecretStoreTypeKey,
	ConfigServiceUrlKey,
	ConfigSessionDurationKey,
	ConfigServiceCacheTypeKey,
	ConfigAuthProviderTypeKey,
	ConfigRecordsDatabaseTypeKey,
	ConfigDatabaseProviderNamesKey,
	ConfigDatabaseProviderTypesKey,
	ConfigDatabaseProviderTableNamesKey,
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
	LocalEnvPath = ".env"
	LocalDefaultConfigPath = "config/default.config"
)
// query parameters
const (
	QueryParamRefreshToken = "refresh-token"
	QueryParamWorkerPath = "worker-path"
	QueryParamWorkerScope = "worker-scope"
	QueryParamLimit = "limit"
	QueryParamAfter = "after"
	QueryParamRedirectUrl = "redirect-url"
	QueryParamRequestId = "request-id"
	QueryParamAccessToken = "access-token"
)
// request context keys
const (
	ContextKeyClaims = "claims"
	ContextKeySessionId = "session-id"
	ContextKeyNonceToken = "nonce-token"
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
	ValueEntryDelim = ";"
	ValueNameDelim = ":"
	ValueItemDelim = ","
	ValueTrue = "true"
	ValueLocalType = "LOCAL"
)

