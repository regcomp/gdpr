package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/internal/views/templ/components"
	"github.com/regcomp/gdpr/pkg/logging"
)

func TestEndpoint(w http.ResponseWriter, r *http.Request) {
	logging.RT.UpdateRequestTrace(r, "TestEndpoint")
	component := components.TestComponent()
	component.Render(r.Context(), w)
}
