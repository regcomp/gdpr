package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/logging"
)

func Logout(cookieManager *auth.CookieManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "PostLogout")
		cookieManager.DestroyAllCookies(w, r)
		w.WriteHeader(http.StatusNoContent)
	}
}
