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
	EndpointApiRecords = "/records"
	// auth
	EndpointAuthLogin = "/login"
	EndpointAuthLoginCallback = "/login-callback"
	EndpointAuthRenewToken = "/renew-token"
	// base
	EndpointBaseHealthz = "/healthz"
	EndpointBaseRegisterServiceWorker = "/register-service-worker"
	// client
	EndpointClientDashboard = "/dashboard"
	EndpointClientLogout = "/logout"
	EndpointClientRecords = "/records"
	EndpointClientTest = "/test"
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
	PathClientRecords = "/app/client/records"
	PathClientTest = "/app/client/test"
)
// service workers
const (
	WorkerAuthRetryPath = "/static/js/worker_auth_retry.js"
	WorkerAuthRetryScope = "/"
)
// config keys
const (
	ConfigAuthProviderTypeKey = "AUTH_PROVIDER_TYPE"
	ConfigDatabaseProviderTypesKey = "DATABASE_PROVIDER_TYPES"
	ConfigServiceUrlKey = "SERVICE_URL"
	ConfigSessionDurationKey = "SESSION_DURATION"
	ConfigSecretStoreTypeKey = "SECRET_STORE_TYPE"
	ConfigRecordsDatabaseTypeKey = "RECORDS_DATABASE_TYPE"
	ConfigDatabaseProviderNamesKey = "DATABASE_PROVIDER_NAMES"
	ConfigDatabaseProviderTableNamesKey = "DATABASE_PROVIDER_TABLE_NAMES"
	ConfigConfigStoreTypeKey = "CONFIG_STORE_TYPE"
	ConfigDefaultPortKey = "DEFAULT_PORT"
	ConfigServiceCacheTypeKey = "SERVICE_CACHE_TYPE"
)

var ConfigAttrs = []string{
	ConfigAuthProviderTypeKey,
	ConfigDatabaseProviderTypesKey,
	ConfigServiceUrlKey,
	ConfigSessionDurationKey,
	ConfigSecretStoreTypeKey,
	ConfigRecordsDatabaseTypeKey,
	ConfigDatabaseProviderNamesKey,
	ConfigDatabaseProviderTableNamesKey,
	ConfigConfigStoreTypeKey,
	ConfigDefaultPortKey,
	ConfigServiceCacheTypeKey,
}
// Cookies
const (
	CookieNameAccessToken = "access-token"
	CookieNameRefreshToken = "refresh-token"
	CookieNameSessionId = "session-id"
)

var CookieNames = []string{
	CookieNameAccessToken,
	CookieNameRefreshToken,
	CookieNameSessionId,
}
// form values
const FormValueNonce = "nonce"
// local files
const (
	LocalEnvPath = ".env"
	LocalDefaultConfigPath = "internal/config/default.config"
)
// query parameters
const (
	QueryParamRequestId = "request-id"
	QueryParamAccessToken = "access-token"
	QueryParamRefreshToken = "refresh-token"
	QueryParamWorkerPath = "worker-path"
	QueryParamWorkerScope = "worker-scope"
	QueryParamLimit = "limit"
	QueryParamAfter = "after"
	QueryParamRedirectUrl = "redirect-url"
)
// request context keys
const (
	ContextKeyClaims = "claims"
	ContextKeySessionId = "session-id"
	ContextKeyNonceToken = "nonce-token"
)
// headers
const (
	HeaderAuthRetryWorkerRunning = "Auth-Retry-Worker-Running"
	HeaderNonceToken = "Nonce-Token"
	HeaderRenewAccessToken = "Renew-Access-Token"
	HeaderServiceWorkerAllowed = "Service-Worker-Allowed"
)
// values
const (
	ValueEntryDelim = ";"
	ValueNameDelim = ":"
	ValueItemDelim = ","
	ValueTrue = "true"
	ValueLocalType = "LOCAL"
)

