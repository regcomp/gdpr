package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/regcomp/gdpr/internal/config"
	"github.com/regcomp/gdpr/internal/middleware"
	"github.com/regcomp/gdpr/internal/routers/api"
	authproxy "github.com/regcomp/gdpr/internal/routers/authproxy"
	"github.com/regcomp/gdpr/internal/routers/client"
	"github.com/regcomp/gdpr/internal/routers/handlers"
	"github.com/regcomp/gdpr/internal/servicecontext"
	"github.com/regcomp/gdpr/pkg/logging"
)

type SubRouter struct {
	MountPath string
	Router    *chi.Mux
}

func CreateRouter(stx *servicecontext.ServiceContext) *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		logging.TraceRequests,
	)

	router.Get(config.EndpointBaseHealthz, handlers.Healthz)
	router.Get(
		config.EndpointBaseRegisterServiceWorker,
		handlers.RegisterServiceWorker(stx.RequestManager),
	)
	router.Get(
		"/favicon.ico",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "/static/favicon.ico")
		}),
	)

	mountRouters(router,
		CreateStaticRouter(),
		CreateAuthProxyRouter(stx),
		CreateAppRouter(stx),
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

	return SubRouter{config.RouterStaticPathPrefix, static}
}

func CreateServiceRouter(stx *servicecontext.ServiceContext) SubRouter {
	service := chi.NewRouter()

	service.Use(
		// middleware.RequestLogging(stx.RequestLogger),
		middleware.SetHSTSPolicy,
		// TODO: Content policies/CORS/ect... go here

		middleware.VerifyAuthRetryIsRunning(stx.RequestManager),
	)

	mountRouters(service,
		CreateAuthProxyRouter(stx),
		CreateAppRouter(stx),
	)
	return SubRouter{config.RouterServicePathPrefix, service}
}

func CreateAppRouter(stx *servicecontext.ServiceContext) SubRouter {
	requiresValidAuth := chi.NewRouter()

	requiresValidAuth.Use(
		middleware.RequiresAuthentication(stx.AuthProvider, stx.CookieManager),
		middleware.HasActiveSession(stx.SessionManager, stx.CookieManager),
		middleware.SetHSTSPolicy,
		// TODO: Content policies/CORS/ect... go here

		middleware.VerifyAuthRetryIsRunning(stx.RequestManager),
	)

	mountRouters(requiresValidAuth,
		CreateClientRouter(stx),
		CreateAPIRouter(stx),
	)

	return SubRouter{MountPath: config.RouterAppPathPrefix, Router: requiresValidAuth}
}

func mountRouters(main *chi.Mux, subrouters ...SubRouter) {
	if len(subrouters) < 1 {
		return
	}

	for _, subrouter := range subrouters {
		main.Mount(subrouter.MountPath, subrouter.Router)
	}
}

func CreateAuthProxyRouter(stx *servicecontext.ServiceContext) SubRouter {
	return SubRouter{
		MountPath: config.RouterAuthPathPrefix,
		Router:    authproxy.CreateAuthProxyRouter(stx),
	}
}

func CreateClientRouter(stx *servicecontext.ServiceContext) SubRouter {
	return SubRouter{
		MountPath: config.RouterClientPathPrefix,
		Router:    client.CreateClientRouter(stx),
	}
}

func CreateAPIRouter(stx *servicecontext.ServiceContext) SubRouter {
	return SubRouter{
		MountPath: config.RouterApiPathPrefix,
		Router:    api.CreateApiRouter(stx),
	}
}
