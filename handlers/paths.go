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
	Test        = "/test"

	// auth
	LoginPath         = "/login"
	LogoutPath        = "/logout"
	LoginCallbackPath = "/logincallback"
	RefreshPath       = "/refresh"

	// client
	DashboardPath = "/dashboard"
)
