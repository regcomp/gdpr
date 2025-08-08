package handlers

import (
	"fmt"
	"net/http"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/helpers"
	"github.com/regcomp/gdpr/logging"
)

func LoginCallback(
	authProvider auth.IAuthProvider,
	cookieManager *auth.CookieManager,
	sessionStore *auth.SessionStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "LoginCallback")
		credentials := auth.Credentials{}

		switch authProvider.GetProviderType() {
		// NOTE: Vendor implementations go here
		case auth.MockProviderType:
			credentials.AccessToken = r.URL.Query().Get(config.QueryParamAccessToken)
			credentials.RefreshToken = r.URL.Query().Get(config.QueryParamRefreshToken)
		default:
			err := fmt.Errorf("unknown auth provider")
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		// TODO: VALIDATE THE JWTS RECIEVED

		accessCookie, err := cookieManager.CreateAccessCookie(credentials.AccessToken)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, accessCookie)

		refreshCookie, err := cookieManager.CreateRefreshCookie(credentials.RefreshToken)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, refreshCookie)

		sessionID, err := sessionStore.CreateSession()
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}
		sessionCookie, err := cookieManager.CreateSessionCookie(sessionID)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, sessionCookie)

		// NOTE: This redirect may want to instead reference where a user was when a refresh token expired.
		http.Redirect(w, r, config.EndpointDashboard, http.StatusSeeOther)
	}
}
