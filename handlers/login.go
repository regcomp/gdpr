package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/templates/pages"
)

func (stx *ServiceContext) GetLogin(w http.ResponseWriter, r *http.Request) {
	stx.Logger.Info("GetLogin", nil)
	page := pages.Login()
	page.Render(r.Context(), w)
}

func (stx *ServiceContext) PostLogin(w http.ResponseWriter, r *http.Request) {
	stx.Logger.Info("PostLogin", nil)
	// Authenticate user
	credentials, err := stx.AuthProvider.AuthenticateUser(r)
	if err != nil {
		// TODO: handle
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    credentials.AccessToken,
		HttpOnly: true,                    // Stops XSS
		Secure:   true,                    // HTTPS only
		SameSite: http.SameSiteStrictMode, // Partially mitigates CSRF
		MaxAge:   stx.AccessTokenDuration,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    credentials.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   stx.RefreshTokenDuration,
	})

	http.Redirect(w, r, DashboardPath, http.StatusSeeOther)
}
