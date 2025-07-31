package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/handlers"
	"github.com/regcomp/gdpr/middleware"
	servicecontext "github.com/regcomp/gdpr/service_context"
)

type SubRouter struct {
	Path   string
	Router *chi.Mux
}

func CreateRouter(stx *servicecontext.ServiceContext) *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		middleware.TraceRequests,
	)

	router.Get(config.EndpointHealthz, healthz)
	router.Get(
		config.EndpointRegisterServiceWorker,
		handlers.RegisterServiceWorker(stx.RequestStore),
	)

	mountRouters(router,
		CreateStaticRouter(),
		CreateServiceRouter(stx),
	)

	return router
}

func CreateStaticRouter() SubRouter {
	static := chi.NewRouter()

	static.Use(
		middleware.ScopeAuthRetryAccess(),
	)

	static.Handle("/*", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static/"))))

	return SubRouter{"/static", static}
}

func CreateServiceRouter(stx *servicecontext.ServiceContext) SubRouter {
	service := chi.NewRouter()

	service.Use(
		// middleware.RequestLogging(stx.RequestLogger),
		middleware.SetHSTSPolicy,
		// TODO: Content policies/CORS/ect... go here

		middleware.VerifyAuthRetryIsRunning(stx.RequestStore),
	)

	mountRouters(service,
		CreateAuthRouter(stx),
		CreateClientRouter(stx),
		CreateAPIRouter(stx),
	)
	return SubRouter{"/", service}
}

func CreateAuthRouter(stx *servicecontext.ServiceContext) SubRouter {
	auth := chi.NewRouter()

	auth.Use(
		middleware.SkipIfAuthenticated(stx.AuthProvider, stx.CookieManager, stx.ConfigStore),
	)

	auth.Get(config.EndpointLogin, handlers.LoginPage)
	auth.Post(config.EndpointLogin, handlers.SubmitLoginCredentials(stx.AuthProvider, stx.ConfigStore))

	// Apparently some providers will hit with either GET or POST
	loginCallback := handlers.LoginCallback(stx.AuthProvider, stx.CookieManager, stx.SessionStore)
	auth.Get(config.EndpointLoginCallback, loginCallback)
	auth.Post(config.EndpointLoginCallback, loginCallback)

	auth.Post(config.EndpointRenewToken, handlers.RenewAccessToken(stx.AuthProvider, stx.CookieManager))

	return SubRouter{Path: config.RouterAuthPathPrefix, Router: auth}
}

func CreateClientRouter(stx *servicecontext.ServiceContext) SubRouter {
	client := chi.NewRouter()

	client.Use(
		middleware.RequiresAuthentication(stx.AuthProvider, stx.CookieManager),
		middleware.HasActiveSession(stx.SessionStore, stx.CookieManager),
		middleware.AddNonceToRequest(stx.NonceStore),
	)

	client.Get(config.EndpointDashboard, handlers.DashboardPage(stx.CookieManager))
	client.Get(config.EndpointTest1, handlers.TestEndpoint1)
	client.Get(config.EndpointTest2, handlers.TestEndpoint2)

	// NOTE: Not sure were this should go
	client.Post(config.EndpointLogout, handlers.Logout(stx.CookieManager))

	return SubRouter{Path: config.RouterClientPathPrefix, Router: client}
}

func CreateAPIRouter(stx *servicecontext.ServiceContext) SubRouter {
	api := chi.NewRouter()

	api.Use(
		middleware.RequiresAuthentication(stx.AuthProvider, stx.CookieManager),
		middleware.HasActiveSession(stx.SessionStore, stx.CookieManager),
	)

	return SubRouter{Path: config.RouterApiPathPrefix, Router: api}
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
