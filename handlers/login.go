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
	default:
		auth.FillCredentialsFromRequestBody(r, &credentials)
	}

	accessCookie := auth.CreateAccessCookie(credentials.AccessToken, stx.CookieKeys)
	http.SetCookie(w, accessCookie)

	refreshCookie := auth.CreateRefreshCookie(credentials.RefreshToken, stx.CookieKeys)
	http.SetCookie(w, refreshCookie)

	sessionCookie, err := stx.SessionStore.Get(r, "session-id")
	if err != nil {
		// TODO: 500
	}
	sessionCookie.Values["session-id"] = auth.GenerateSessionID()
	sessionCookie.Save(r, w)

	// NOTE: This redirect may want to instead reference where a user was when a refresh token expired.
	http.Redirect(w, r, DashboardPath, http.StatusSeeOther)
}
