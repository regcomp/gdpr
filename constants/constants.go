package constants

// Subrouter path prefixes
const (
	RouterClientPathPrefix = "/"
	RouterApiPathPrefix    = "/api"
	RouterAuthPathPrefix   = "/auth"
)

// endpoints
const (
	// base
	EndpointHealthz = "/healthz"
	EndpointTest1   = "/test1"
	EndpointTest2   = "/test2"

	// auth
	EndpointLogin                 = "/login"
	EndpointLogout                = "/logout"
	EndpointLoginCallback         = "/logincallback"
	EndpointRenewToken            = "/renewtoken"
	EndpointRegisterServiceWorker = "/registerserviceworker"

	// client
	EndpointDashboard = "/dashboard"
)

// service workers
const (
	WorkerAuthRetryPath  = "/static/sw/auth_retry.js"
	WorkerAuthRetryScope = "/"
)

// config keys
const (
	ConfigServiceURLKey       = "SERVICE_URL"
	ConfigDefaultPortKey      = "DEFAULT_PORT"
	ConfigSessionDurationKey  = "SESSION_DURATION"
	ConfigServiceCacheTypeKey = "SERVICE_CACHE_TYPE"
	ConfigSecretStoreTypeKey  = "SECRET_STORE_TYPE"
	ConfigAuthProvierTypeKey  = "AUTH_PROVIDER"
	ConfigRequestTracerOnKey  = "REQUEST_TRACER_ON"
)

var ConfigAttrs = []string{
	ConfigServiceURLKey,
	ConfigDefaultPortKey,
	ConfigSessionDurationKey,
	ConfigServiceCacheTypeKey,
	ConfigSecretStoreTypeKey,
	ConfigAuthProvierTypeKey,
	ConfigRequestTracerOnKey,
}

// Cookies
const (
	AccessCookieName  = "access-token"
	RefreshCookieName = "refresh-token"
	SessionCookieName = "session-id"
)

// form values
const FormValueNonce = "nonce"

// local files
const (
	LocalEnvPath    = ".env"
	LocalConfigPath = "config/default.config"
)

// query parameters
const (
	QueryParamRedirectURL  = "redirect-url"
	QueryParamRequestID    = "request-id"
	QueryParamAccessToken  = "access-token"
	QueryParamRefreshToken = "refresh-token"
	QueryParamWorkerPath   = "sw-path"
	QueryParamWorkerScope  = "sw-scope"
)

// request context keys
const (
	ContextKeyClaims     = "claims"
	ContextKetSessionID  = "session-id"
	ContextKeyNonceToken = "nonce-token"
)
