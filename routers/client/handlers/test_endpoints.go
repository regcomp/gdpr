package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/logging"
	"github.com/regcomp/gdpr/templates/components"
)

func TestEndpoint1(w http.ResponseWriter, r *http.Request) {
	logging.RT.UpdateRequestTrace(r, "TestEndpoint1")
	component := components.DeletionRequests()
	component.Render(r.Context(), w)
}

func TestEndpoint2(w http.ResponseWriter, r *http.Request) {
	logging.RT.UpdateRequestTrace(r, "TestEndpoint2")
	component := components.TestComponent2()
	component.Render(r.Context(), w)
}
