package handlers

import (
	"log"
	"net/http"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/logging"
	"github.com/regcomp/gdpr/templates/pages"
)

func DashboardPage(cookieManager *auth.CookieManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "DashboardPage")

		// WARN: TEMPORARY!!

		accessToken, err := cookieManager.GetAccessToken(r)
		if err != nil {
			log.Panicf("%s", err.Error())
		}
		// refreshToken, err := auth.GetRefreshToken(r, stx.CookieKeys)
		// if err != nil {
		// 	log.Panicf("%s", err.Error())
		// }

		sessionID, err := cookieManager.GetSessionID(r)
		if err != nil {
			log.Panicf("%s", err.Error())
		}

		pages.Dashboard(accessToken, "---", sessionID).Render(r.Context(), w)
		// -----

		// dashboard := pages.Dashboard()
		// dashboard.Render(r.Context(), w)
	}
}
