package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/templates/components"
)

func (stx *ServiceContext) TestEndpoint1(w http.ResponseWriter, r *http.Request) {
	stx.RequestTracer.UpdateRequestTrace(r, "TestEndpoint1")
	component := components.TestComponent1()
	component.Render(r.Context(), w)
}

func (stx *ServiceContext) TestEndpoint2(w http.ResponseWriter, r *http.Request) {
	stx.RequestTracer.UpdateRequestTrace(r, "TestEndpoint2")
	component := components.TestComponent2()
	component.Render(r.Context(), w)
}
