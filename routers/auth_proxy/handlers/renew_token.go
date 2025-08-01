package handlers

import (
	"fmt"
	"net/http"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/logging"
)

func RenewAccessToken(authProvider auth.IAuthProvider, cookieManager *auth.CookieManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "RenewAccessToken")

		refreshToken, err := cookieManager.GetRefreshToken(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not get refresh token, err=%s", err.Error()), http.StatusInternalServerError)
		}
		accessToken, err := authProvider.GetNewAccessToken(refreshToken, r)
		if err != nil {
			http.Error(w, "could not renew access token", http.StatusInternalServerError)
		}

		accessCookie, err := cookieManager.CreateAccessCookie(accessToken)
		if err != nil {
			http.Error(w, "could not create access cookie", http.StatusInternalServerError)
		}

		http.SetCookie(w, accessCookie)
		w.WriteHeader(http.StatusOK)
	}
}
