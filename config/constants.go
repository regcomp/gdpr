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
	ConfigServiceUrlKey = "SERVICE_URL"
	ConfigDefaultPortKey = "DEFAULT_PORT"
	ConfigSessionDurationKey = "SESSION_DURATION"
	ConfigServiceCacheTypeKey = "SERVICE_CACHE_TYPE"
	ConfigSecretStoreTypeKey = "SECRET_STORE_TYPE"
	ConfigAuthProviderTypeKey = "AUTH_PROVIDER_TYPE"
)

var ConfigAttrs = []string{
	ConfigConfigStoreTypeKey,
	ConfigServiceUrlKey,
	ConfigDefaultPortKey,
	ConfigSessionDurationKey,
	ConfigServiceCacheTypeKey,
	ConfigSecretStoreTypeKey,
	ConfigAuthProviderTypeKey,
}
// Cookies
const (
	CookieNameSessionId = "session-id"
	CookieNameAccessToken = "access-token"
	CookieNameRefreshToken = "refresh-token"
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
	QueryParamRedirectUrl = "redirect-url"
	QueryParamRequestId = "request-id"
	QueryParamAccessToken = "access-token"
	QueryParamRefreshToken = "refresh-token"
	QueryParamWorkerPath = "worker-path"
	QueryParamWorkerScope = "worker-scope"
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
	ValueTrue = "true"
)

