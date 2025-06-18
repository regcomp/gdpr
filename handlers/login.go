package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/templates/pages"
)

func (stx *ServiceContext) GetLogin(w http.ResponseWriter, r *http.Request) {
	page := pages.Login()
	page.Render(r.Context(), w)
}

func (stx *ServiceContext) PostLogin(w http.ResponseWriter, r *http.Request) {
	callbackURL := NewURL("http", stx.HostPath, LoginCallbackPath)
	stx.AuthProvider.AuthenticateUser(w, r, callbackURL)
}

func (stx *ServiceContext) LoginCallback(w http.ResponseWriter, r *http.Request) {
	credentials := auth.Credentials{}

	switch stx.AuthProvider.GetProviderType() {
	// TODO: Vendor implementations go here
	case auth.MOCK:
		credentials.AccessToken = r.URL.Query().Get("access")
		credentials.RefreshToken = r.URL.Query().Get("refresh")
	default:
		http.Error(w, "auth provider not implemented", http.StatusInternalServerError)
	}

	accessCookie, err := auth.CreateAccessCookie(credentials.AccessToken, stx.CookieKeys)
	if err != nil {
		// TODO:
	}
	http.SetCookie(w, accessCookie)

	refreshCookie, err := auth.CreateRefreshCookie(credentials.RefreshToken, stx.CookieKeys)
	if err != nil {
		// TODO:
	}
	http.SetCookie(w, refreshCookie)

	sessionCookie, err := auth.CreateSessionCookie(stx.CookieKeys)
	if err != nil {
		// TODO:
	}
	http.SetCookie(w, sessionCookie)

	// NOTE: This redirect may want to instead reference where a user was when a refresh token expired.
	http.Redirect(w, r, DashboardPath, http.StatusSeeOther)
}
