package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/helpers"
	"github.com/regcomp/gdpr/logging"
)

func RenewAccessToken(authProvider auth.IAuthProvider, cookieManager *auth.CookieManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "RenewAccessToken")

		refreshToken, err := cookieManager.GetRefreshToken(r)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}
		accessToken, err := authProvider.GetNewAccessToken(refreshToken, r)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		accessCookie, err := cookieManager.CreateAccessCookie(accessToken)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, accessCookie)
		w.WriteHeader(http.StatusOK)
	}
}
