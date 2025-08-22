package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/internal/views"
	"github.com/regcomp/gdpr/pkg/helpers"
	"github.com/regcomp/gdpr/pkg/logging"
)

func RecordsComponent(w http.ResponseWriter, r *http.Request) {
	logging.RT.UpdateRequestTrace(r, "RecordsComponent")
	err := views.WriteRecordsManagement(w, r.Context())
	if err != nil {
		helpers.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}
}
