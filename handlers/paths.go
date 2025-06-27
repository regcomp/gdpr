package handlers

// Subrouter path prefixes
const (
	ClientRouterPathPrefix = "/"
	ApiRouterPathPrefix    = "/api"
)

// endpoints
const (
	HealthzPath       = "/healthz"
	LoginPath         = "/login"
	LoginCallbackPath = "/logincallback"
	RefreshPath       = "/auth/refresh"
	DashboardPath     = "/dashboard"
	Test              = "/test"
)
