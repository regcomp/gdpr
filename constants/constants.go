package constants

import (
	"encoding/json"
	"log"
	"os"
)

// CONSTANTS SHARED BETWEEN JS AND GO ---------
const SharedConstantsFilePath = "./constants/shared.json"

var (
	// headers
	HeaderNonceToken             string
	HeaderRenewAccessToken       string
	HeaderServiceWorkerAllowed   string
	HeaderAuthRetryWorkerRunning string

	// paths
	PathAuthRenewToken string
)

func InitializeSharedConstants() {
	shared := struct { // NOTE: should mirror the json structure
		Headers map[string]string `json:"HEADERS"`
		Paths   map[string]string `json:"PATHS"`
	}{
		Headers: make(map[string]string),
		Paths:   make(map[string]string),
	}

	sharedBytes, err := os.ReadFile(SharedConstantsFilePath)
	if err != nil {
		log.Panicf("could not read file %s, err=%s", SharedConstantsFilePath, err.Error())
	}

	err = json.Unmarshal(sharedBytes, &shared)
	if err != nil {
		log.Panicf("invalid JSON: err=%s", err.Error())
	}

	// headers
	HeaderNonceToken = shared.Headers["NONCE_TOKEN"]
	HeaderRenewAccessToken = shared.Headers["RENEW_ACCESS_TOKEN"]
	HeaderServiceWorkerAllowed = shared.Headers["SERVICE_WORKER_ALLOWED"]
	HeaderAuthRetryWorkerRunning = shared.Headers["AUTH_RETRY_WORKER_RUNNING"]

	// paths
	PathAuthRenewToken = shared.Paths["AUTH_RENEW"]
}

// --------------------------------------------

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
	AuthRetryWorkerPath  = "/static/sw/auth_retry.js"
	AuthRetryWorkerScope = "/"
)

// config keys
const (
	ConfigServiceURLKey       = "SERVICE_URL"
	ConfigDefaultPortKey      = "DEFAULT_PORT"
	ConfigSessionDurationKey  = "SESSION_DURATION"
	ConfigServiceCacheTypeKey = "SERVICE_CACHE_TYPE"
	ConfigSecretStoreTypeKey  = "SECRET_STORE_TYPE"
	ConfigAuthProvierKey      = "AUTH_PROVIDER"
	ConfigDebugTraceRequests  = "DEBUG_TRACE_REQUESTS"
)

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
	QueryParamSWPath       = "sw-path"
	QueryParamSWScope      = "sw-scope"
)

// request context keys
const (
	ContextKeyClaims     = "claims"
	ContextKetSessionID  = "session-id"
	ContextKeyNonceToken = "nonce-token"
)
