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
	ConfigConfigStoreTypeKey = "CONFIG_STORE_TYPE"
	ConfigServiceUrlKey = "SERVICE_URL"
	ConfigDefaultPortKey = "DEFAULT_PORT"
	ConfigSecretStoreTypeKey = "SECRET_STORE_TYPE"
	ConfigAuthProviderTypeKey = "AUTH_PROVIDER_TYPE"
	ConfigRecordsDatabaseTypeKey = "RECORDS_DATABASE_TYPE"
	ConfigDatabaseProviderTypesKey = "DATABASE_PROVIDER_TYPES"
	ConfigSessionDurationKey = "SESSION_DURATION"
	ConfigServiceCacheTypeKey = "SERVICE_CACHE_TYPE"
	ConfigDatabaseProviderNamesKey = "DATABASE_PROVIDER_NAMES"
	ConfigDatabaseProviderTableNamesKey = "DATABASE_PROVIDER_TABLE_NAMES"
)

var ConfigAttrs = []string{
	ConfigConfigStoreTypeKey,
	ConfigServiceUrlKey,
	ConfigDefaultPortKey,
	ConfigSecretStoreTypeKey,
	ConfigAuthProviderTypeKey,
	ConfigRecordsDatabaseTypeKey,
	ConfigDatabaseProviderTypesKey,
	ConfigSessionDurationKey,
	ConfigServiceCacheTypeKey,
	ConfigDatabaseProviderNamesKey,
	ConfigDatabaseProviderTableNamesKey,
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

