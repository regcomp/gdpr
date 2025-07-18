package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/templates/pages"
)

func (stx *ServiceContext) TestEndpoint(w http.ResponseWriter, r *http.Request) {
	stx.RequestTracer.UpdateRequestTrace(r, "TestEndpoint")
	page := pages.TestPage()
	page.Render(r.Context(), w)
}
