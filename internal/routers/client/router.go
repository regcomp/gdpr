package client

import (
	"github.com/go-chi/chi/v5"
	"github.com/regcomp/gdpr/internal/config"
	"github.com/regcomp/gdpr/internal/middleware"
	"github.com/regcomp/gdpr/internal/routers/client/handlers"
	"github.com/regcomp/gdpr/internal/servicecontext"
)

func CreateClientRouter(stx *servicecontext.ServiceContext) *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(
		middleware.AddNonceToRequest(stx.NonceManager),
	)

	mux.Get(config.EndpointClientDashboard, handlers.DashboardPage(stx.CookieManager))
	mux.Get(config.EndpointClientRecords, handlers.RecordsComponent)
	mux.Get(config.EndpointClientTest, handlers.TestEndpoint)

	mux.Post(config.EndpointClientLogout, handlers.Logout(stx.CookieManager))

	return mux
}
