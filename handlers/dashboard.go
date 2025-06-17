package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/templates/pages"
)

func (stx *ServiceContext) GetDashboard(w http.ResponseWriter, r *http.Request) {
	// WARN: TEMPORARY!!

	ac, _ := r.Cookie("access-token")
	at := ac.Value
	rc, _ := r.Cookie("refresh-token")
	rt := rc.Value
	user := "Baz"

	pages.Dashboard(at, rt, user).Render(r.Context(), w)
	// -----

	// dashboard := pages.Dashboard()
	// dashboard.Render(r.Context(), w)
}
