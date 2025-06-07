package handlers

import (
	"net/http"
)

func (stx *ServiceContext) HandleTest(w http.ResponseWriter, r *http.Request) {
	if stx.Testing == "" {
		RespondWithError(w, http.StatusInternalServerError, "No testing value initiated")
	}

	payload := struct {
		Result string `json:"result"`
	}{
		Result: stx.Testing,
	}

	RespondWithJSON(w, http.StatusOK, payload)
}
