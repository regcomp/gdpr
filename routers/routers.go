package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/regcomp/gdpr/handlers"
)

const (
	opsPathPrefix = "/"
	apiPathPrefix = "/api"
)

type SubRouter struct {
	Path   string
	Router *chi.Mux
}

func CreateRouter(subRouters ...SubRouter) *chi.Mux {
	router := chi.NewRouter()

	// TODO: Add service-wide middleware.

	router.Handle("/static/*", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static/"))))

	router.Get("/healthz", healthz)

	router.Get("/login", handlers.STX.GetLogin)
	router.Post("/login", handlers.STX.PostLogin)

	mountRouters(router, subRouters...)

	return router
}

func mountRouters(main *chi.Mux, subrouters ...SubRouter) {
	if len(subrouters) < 1 {
		return
	}

	for _, subrouter := range subrouters {
		main.Mount(subrouter.Path, subrouter.Router)
	}
}

func CreateOpsRouter() SubRouter {
	ops := chi.NewRouter()

	ops.Use(handlers.STX.HasAuth)

	ops.Get("/dashboard", handlers.STX.GetDashboard)

	return SubRouter{Path: opsPathPrefix, Router: ops}
}

func CreateApiRouter() SubRouter {
	api := chi.NewRouter()

	api.Use(handlers.STX.HasAuth)

	return SubRouter{Path: apiPathPrefix, Router: api}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
