package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/templates/pages"
)

func (stx *ServiceContext) GetDashboard(w http.ResponseWriter, r *http.Request) {
	dashboard := pages.Dashboard()
	dashboard.Render(r.Context(), w)
}
