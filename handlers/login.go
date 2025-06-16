package handlers

import (
	"encoding/json"
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
	// TODO: implementations go here
	default:
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			// TODO:
		}
	}
	accessSession, err := stx.SessionManager.RefreshStore.Get(r, "access-token")
	if err != nil {
		// TODO:
	}
	accessSession.Values["access-token"] = credentials.AccessToken
	err = accessSession.Save(r, w)
	if err != nil {
		// TODO:
	}

	refreshSession, err := stx.SessionManager.RefreshStore.Get(r, "refresh-token")
	if err != nil {
		// TODO:
	}
	refreshSession.Values["access-token"] = credentials.RefreshToken
	err = refreshSession.Save(r, w)
	if err != nil {
		// TODO:
	}

	// NOTE: This redirect may want to instead reference where a user was when a refresh token expired.
	http.Redirect(w, r, DashboardPath, http.StatusMovedPermanently)
}
