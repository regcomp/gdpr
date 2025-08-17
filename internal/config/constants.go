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
	WorkerAuthRetryPath = "/static/js/sw/auth_retry.sw.js"
	WorkerAuthRetryScope = "/"
)
// config keys
const (
	ConfigDatabaseProviderTypesKey = "DATABASE_PROVIDER_TYPES"
	ConfigConfigStoreTypeKey = "CONFIG_STORE_TYPE"
	ConfigServiceCacheTypeKey = "SERVICE_CACHE_TYPE"
	ConfigRecordsDatabaseTypeKey = "RECORDS_DATABASE_TYPE"
	ConfigDatabaseProviderNamesKey = "DATABASE_PROVIDER_NAMES"
	ConfigDatabaseProviderTableNamesKey = "DATABASE_PROVIDER_TABLE_NAMES"
	ConfigServiceUrlKey = "SERVICE_URL"
	ConfigDefaultPortKey = "DEFAULT_PORT"
	ConfigSessionDurationKey = "SESSION_DURATION"
	ConfigSecretStoreTypeKey = "SECRET_STORE_TYPE"
	ConfigAuthProviderTypeKey = "AUTH_PROVIDER_TYPE"
)

var ConfigAttrs = []string{
	ConfigDatabaseProviderTypesKey,
	ConfigConfigStoreTypeKey,
	ConfigServiceCacheTypeKey,
	ConfigRecordsDatabaseTypeKey,
	ConfigDatabaseProviderNamesKey,
	ConfigDatabaseProviderTableNamesKey,
	ConfigServiceUrlKey,
	ConfigDefaultPortKey,
	ConfigSessionDurationKey,
	ConfigSecretStoreTypeKey,
	ConfigAuthProviderTypeKey,
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
	QueryParamAfter = "after"
	QueryParamRedirectUrl = "redirect-url"
	QueryParamRequestId = "request-id"
	QueryParamAccessToken = "access-token"
	QueryParamRefreshToken = "refresh-token"
	QueryParamWorkerPath = "worker-path"
	QueryParamWorkerScope = "worker-scope"
	QueryParamLimit = "limit"
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
	ValueTrue = "true"
	ValueLocalType = "LOCAL"
	ValueEntryDelim = ";"
	ValueNameDelim = ":"
	ValueItemDelim = ","
)

