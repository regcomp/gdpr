package handlers

import (
	"fmt"
	"net/http"

	"github.com/regcomp/gdpr/internal/auth"
	"github.com/regcomp/gdpr/internal/caching"
	"github.com/regcomp/gdpr/internal/config"
	"github.com/regcomp/gdpr/pkg/helpers"
	"github.com/regcomp/gdpr/pkg/logging"
)

func LoginCallback(ap auth.IAuthProvider, cm *caching.CookieManager, sm *caching.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "LoginCallback")

		credentials := auth.Credentials{}

		switch ap.GetProviderType() {
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

		accessCookie, err := cm.CreateAccessCookie(credentials.AccessToken)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, accessCookie)

		refreshCookie, err := cm.CreateRefreshCookie(credentials.RefreshToken)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, refreshCookie)

		sessionID, err := sm.CreateSession()
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}
		sessionCookie, err := cm.CreateSessionCookie(sessionID)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, sessionCookie)

		// NOTE: This redirect may want to instead reference where a user was when a refresh token expired.
		http.Redirect(w, r, config.PathClientDashboard, http.StatusSeeOther)
	}
}
