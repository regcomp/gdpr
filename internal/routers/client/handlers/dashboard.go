package handlers

import (
	"log"
	"net/http"

	"github.com/regcomp/gdpr/internal/caching"
	"github.com/regcomp/gdpr/internal/views"
	"github.com/regcomp/gdpr/pkg/helpers"
	"github.com/regcomp/gdpr/pkg/logging"
)

func DashboardPage(cm *caching.CookieManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "DashboardPage")

		// WARN: TEMPORARY!!

		accessToken, err := cm.GetAccessToken(r)
		if err != nil {
			log.Panicf("%s", err.Error())
		}

		sessionID, err := cm.GetSessionID(r)
		if err != nil {
			log.Panicf("%s", err.Error())
		}

		err = views.ServeDashboard(w, r.Context(), accessToken, "---", sessionID)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}
	}
}
