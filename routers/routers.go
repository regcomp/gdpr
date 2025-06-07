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

	// TODO: Add service-wide middleware. eg, auth, logging, ect
	router.Use(handlers.STX.HasAuth)

	router.Handle("/static/*", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static/"))))

	router.Get("/healthz", healthz)

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
	router := chi.NewRouter()

	router.Get("/context", handlers.STX.HandleTest)

	return SubRouter{Path: opsPathPrefix, Router: router}
}

func CreateApiRouter() SubRouter {
	router := chi.NewRouter()

	return SubRouter{Path: apiPathPrefix, Router: router}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
