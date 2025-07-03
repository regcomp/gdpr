package handlers

import (
	"log"
	"net/http"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/templates/pages"
)

func (stx *ServiceContext) GetDashboard(w http.ResponseWriter, r *http.Request) {
	stx.RequestTracer.UpdateActiveTrace("GetDashboard")

	// WARN: TEMPORARY!!

	accessToken, err := auth.GetAccessToken(r, stx.CookieKeys)
	if err != nil {
		log.Panicf("%s", err.Error())
	}
	refreshToken, err := auth.GetRefreshToken(r, stx.CookieKeys)
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	sessionID, err := auth.GetSessionID(r, stx.CookieKeys)
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	pages.Dashboard(accessToken, refreshToken, sessionID).Render(r.Context(), w)
	// -----

	// dashboard := pages.Dashboard()
	// dashboard.Render(r.Context(), w)
}
