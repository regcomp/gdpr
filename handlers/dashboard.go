package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/templates/pages"
)

func (stx *ServiceContext) GetDashboard(w http.ResponseWriter, r *http.Request) {
	at, err := r.Cookie("access_token")
	if err != nil {
		// TODO:
	}

	rt, err := r.Cookie("access_token")
	if err != nil {
		// TODO:
	}

	dashboard := pages.Dashboard(at.Value, rt.Value)
	dashboard.Render(r.Context(), w)
}
