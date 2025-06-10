package handlers

import (
	"net/http"
)

func (stx *ServiceContext) TestEndpoint(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Test string `json:"test"`
	}{
		Test: "endpoint",
	}
	RespondWithJSON(w, http.StatusOK, payload)
}
