package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/regcomp/gdpr/handlers"
)

const (
	swAuthRetryPath  = "/static/sw/auth_retry.js"
	swAuthRetryScope = "/"
)

type SubRouter struct {
	Path   string
	Router *chi.Mux
}

func CreateRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		handlers.STX.TraceRequests,
	)

	router.Get(handlers.HealthzPath, healthz)
	router.Get(handlers.Test, handlers.STX.TestEndpoint)

	mountRouters(router,
		CreateStaticRouter(),
		CreateServiceRouter(),
	)

	return router
}

func CreateStaticRouter() SubRouter {
	static := chi.NewRouter()

	static.Use(
		handlers.STX.ScopeServiceWorkerAccess(swAuthRetryPath, swAuthRetryScope),
	)

	static.Handle("/*", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static/"))))

	return SubRouter{"/static", static}
}

func CreateServiceRouter() SubRouter {
	service := chi.NewRouter()

	service.Use(
		// handlers.STX.Logging,
		handlers.STX.SetHSTSPolicy,
		// TODO: Content policies/CORS/ect... go here

		handlers.STX.VerifyServiceWorkerIsRunning(
			swAuthRetryPath,
			swAuthRetryScope,
			"SW-Auth-Retry-Running",
		),
	)

	mountRouters(service,
		CreateAuthRouter(),
		CreateClientRouter(),
		CreateAPIRouter(),
	)
	return SubRouter{"/", service}
}

func CreateAuthRouter() SubRouter {
	auth := chi.NewRouter()

	auth.Get(handlers.LoginPath, handlers.STX.GetLogin)
	auth.Post(handlers.LoginPath, handlers.STX.PostLogin)

	auth.Get(handlers.LoginCallbackPath, handlers.STX.LoginCallback)
	auth.Post(handlers.LoginCallbackPath, handlers.STX.LoginCallback)

	auth.Post(handlers.RefreshPath, handlers.STX.PostRefresh)
	auth.Post(handlers.LogoutPath, handlers.STX.PostLogout)

	return SubRouter{Path: handlers.AuthRouterPathPrefix, Router: auth}
}

func CreateClientRouter() SubRouter {
	client := chi.NewRouter()

	client.Use(
		handlers.STX.IsAuthenticated,
		handlers.STX.HasActiveSession,
	)

	client.Get(handlers.DashboardPath, handlers.STX.GetDashboard)

	return SubRouter{Path: handlers.ClientRouterPathPrefix, Router: client}
}

func CreateAPIRouter() SubRouter {
	api := chi.NewRouter()

	api.Use(
		handlers.STX.IsAuthenticated,
		handlers.STX.HasActiveSession,
	)

	return SubRouter{Path: handlers.ApiRouterPathPrefix, Router: api}
}

func mountRouters(main *chi.Mux, subrouters ...SubRouter) {
	if len(subrouters) < 1 {
		return
	}

	for _, subrouter := range subrouters {
		main.Mount(subrouter.Path, subrouter.Router)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
