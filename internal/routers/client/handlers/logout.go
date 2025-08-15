package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/internal/caching"
	"github.com/regcomp/gdpr/pkg/logging"
)

func Logout(cm *caching.CookieManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "PostLogout")
		cm.DestroyAllCookies(w, r)
		w.WriteHeader(http.StatusNoContent)
	}
}
