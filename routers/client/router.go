package client

import (
	"github.com/go-chi/chi/v5"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/middleware"
	"github.com/regcomp/gdpr/routers/client/handlers"
	servicecontext "github.com/regcomp/gdpr/service_context"
)

func CreateClientRouter(stx *servicecontext.ServiceContext) *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(
		middleware.AddNonceToRequest(stx.NonceManager),
	)

	mux.Get(config.EndpointDashboard, handlers.DashboardPage(stx.CookieManager))
	mux.Get(config.EndpointTest1, handlers.TestEndpoint1)
	mux.Get(config.EndpointTest2, handlers.TestEndpoint2)

	mux.Post(config.EndpointLogout, handlers.Logout(stx.CookieManager))

	return mux
}
