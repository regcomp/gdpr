package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/caching"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/constants"
	"github.com/regcomp/gdpr/logging"
	"github.com/regcomp/gdpr/templates/pages"
)

func RegisterServiceWorker(requestStore caching.IRequestStore, configStore config.IConfigStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "RegisterServiceWorker")
		requestID := r.URL.Query().Get(constants.QueryParamRequestID)
		if requestID == "" {
			// TODO: fatal
		}
		cachedRequest, err := requestStore.RetrieveRequest(requestID)
		if err != nil {
			cachedRequest = &caching.CachedRequest{
				URL:    configStore.GetServiceURL(),
				Method: "GET",
			}
		}

		swPath := r.URL.Query().Get(constants.QueryParamSWPath)
		swScope := r.URL.Query().Get(constants.QueryParamSWScope)
		if swPath == "" || swScope == "" {
			http.Error(w, "missing service worker information", http.StatusBadRequest)
		}

		pages.RegisterServiceWorker(swPath, swScope, cachedRequest).Render(r.Context(), w)
	})
}
