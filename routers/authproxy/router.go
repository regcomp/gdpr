package authproxy

import (
	"github.com/go-chi/chi/v5"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/middleware"
	"github.com/regcomp/gdpr/routers/authproxy/handlers"
	servicecontext "github.com/regcomp/gdpr/service_context"
)

func CreateAuthProxyRouter(stx *servicecontext.ServiceContext) *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(
		middleware.SkipIfAuthenticated(stx.AuthProvider, stx.CookieManager, stx.ConfigStore),
	)

	mux.Get(config.EndpointLogin, handlers.LoginPage)
	mux.Post(config.EndpointLogin, handlers.SubmitLoginCredentials(stx.AuthProvider, stx.ConfigStore))

	// Apparently some providers will hit with either GET or POST
	loginCallback := handlers.LoginCallback(stx.AuthProvider, stx.CookieManager, stx.SessionStore)
	mux.Get(config.EndpointLoginCallback, loginCallback)
	mux.Post(config.EndpointLoginCallback, loginCallback)

	mux.Post(config.EndpointRenewToken, handlers.RenewAccessToken(stx.AuthProvider, stx.CookieManager))

	return mux
}
