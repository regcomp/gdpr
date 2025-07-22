package handlers

// Subrouter path prefixes
const (
	ClientRouterPathPrefix = "/"
	ApiRouterPathPrefix    = "/api"
	AuthRouterPathPrefix   = "/auth"
)

// endpoints
const (
	// base
	HealthzPath = "/healthz"
	Test1Path   = "/test1"
	Test2Path   = "/test2"

	// auth
	LoginPath                 = "/login"
	LogoutPath                = "/logout"
	LoginCallbackPath         = "/logincallback"
	RefreshPath               = "/refresh"
	RegisterServiceWorkerPath = "/registerserviceworker"

	// client
	DashboardPath = "/dashboard"
)

// service workers
const (
	SWAuthRetryPath  = "/static/sw/auth_retry.js"
	SWAuthRetryScope = "/"
)
