package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/regcomp/gdpr/handlers"
)

type SubRouter struct {
	Path   string
	Router *chi.Mux
}

func CreateRouter(subRouters ...SubRouter) *chi.Mux {
	router := chi.NewRouter()

	router.Use(handlers.STX.Logging)

	router.Handle("/static/*", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static/"))))

	router.Get(handlers.HealthzPath, healthz)

	router.Get(handlers.LoginPath, handlers.STX.GetLogin)
	router.Post(handlers.LoginPath, handlers.STX.PostLogin)
	// router.Route(handlers.LoginCallbackPath, func(r chi.Router) {
	// 	r.Get("/", handlers.STX.LoginCallback)
	// 	r.Post("/", handlers.STX.LoginCallback)
	// })
	router.Post(handlers.LoginCallbackPath, handlers.STX.LoginCallback)

	router.Get(handlers.Test, handlers.STX.TestEndpoint)

	mountRouters(router, subRouters...)

	return router
}

func CreateClientRouter() SubRouter {
	client := chi.NewRouter()

	client.Use(handlers.STX.IsAuthenticated)

	client.Get(handlers.DashboardPath, handlers.STX.GetDashboard)

	return SubRouter{Path: handlers.ClientRouterPathPrefix, Router: client}
}

func CreateApiRouter() SubRouter {
	api := chi.NewRouter()

	api.Use(handlers.STX.IsAuthenticated)

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
