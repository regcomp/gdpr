package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/templates/pages"
)

func (stx *ServiceContext) TestEndpoint(w http.ResponseWriter, r *http.Request) {
	page := pages.TestPage()
	page.Render(r.Context(), w)
}
