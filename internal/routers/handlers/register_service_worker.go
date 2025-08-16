package handlers

import (
	"fmt"
	"net/http"

	"github.com/regcomp/gdpr/internal/caching"
	"github.com/regcomp/gdpr/internal/config"
	"github.com/regcomp/gdpr/internal/views"
	"github.com/regcomp/gdpr/pkg/helpers"
	"github.com/regcomp/gdpr/pkg/logging"
)

func RegisterServiceWorker(requestStore *caching.RequestManager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "RegisterServiceWorker")
		requestID := r.URL.Query().Get(config.QueryParamRequestId)
		if requestID == "" {
			err := fmt.Errorf("missing request id parameter")
			helpers.RespondWithError(w, err, http.StatusBadRequest)
			return
		}
		cachedRequest, err := requestStore.RetrieveRequest(requestID)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		swPath := r.URL.Query().Get(config.QueryParamWorkerPath)
		swScope := r.URL.Query().Get(config.QueryParamWorkerScope)
		if swPath == "" || swScope == "" {
			err := fmt.Errorf("missing service worker information")
			helpers.RespondWithError(w, err, http.StatusBadRequest)
			return
		}

		err = views.ServeRegisterServiceWorker(w, r.Context(), cachedRequest, swPath, swScope)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
		}
	})
}
