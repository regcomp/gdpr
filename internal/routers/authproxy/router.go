package authproxy

import (
	"github.com/go-chi/chi/v5"

	"github.com/regcomp/gdpr/internal/config"
	"github.com/regcomp/gdpr/internal/middleware"
	"github.com/regcomp/gdpr/internal/routers/authproxy/handlers"
	"github.com/regcomp/gdpr/internal/servicecontext"
)

func CreateAuthProxyRouter(stx *servicecontext.ServiceContext) *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(
		middleware.SkipIfAuthenticated(stx.AuthProvider, stx.CookieManager, stx.ConfigStore),
	)

	mux.Get(config.EndpointAuthLogin, handlers.LoginPage)
	mux.Post(config.EndpointAuthLogin, handlers.SubmitLoginCredentials(stx.AuthProvider, stx.ConfigStore))

	// Apparently some providers will hit with either GET or POST
	loginCallback := handlers.LoginCallback(stx.AuthProvider, stx.CookieManager, stx.SessionManager)
	mux.Get(config.EndpointAuthLoginCallback, loginCallback)
	mux.Post(config.EndpointAuthLoginCallback, loginCallback)

	mux.Post(config.EndpointAuthRenewToken, handlers.RenewAccessToken(stx.AuthProvider, stx.CookieManager))

	return mux
}
