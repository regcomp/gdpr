package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/templates/pages"
)

func (stx *ServiceContext) HandleTest(w http.ResponseWriter, r *http.Request) {
	if stx.Testing == "" {
		RespondWithError(w, http.StatusInternalServerError, "No testing value initiated")
	}

	page := pages.Page(stx.Testing, "this is a test")
	page.Render(r.Context(), w)
}
