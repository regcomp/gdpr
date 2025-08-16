package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/internal/auth"
	"github.com/regcomp/gdpr/internal/caching"
	"github.com/regcomp/gdpr/pkg/helpers"
	"github.com/regcomp/gdpr/pkg/logging"
)

func RenewAccessToken(ap auth.IAuthProvider, cm *caching.CookieManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "RenewAccessToken")

		refreshToken, err := cm.GetRefreshToken(r)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusUnauthorized)
			return
		}
		accessToken, err := ap.GetNewAccessToken(refreshToken, r)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		accessCookie, err := cm.CreateAccessCookie(accessToken)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, accessCookie)
		w.WriteHeader(http.StatusOK)
	}
}
