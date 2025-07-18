package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/templates/pages"
)

func (stx *ServiceContext) RegisterAuthRetryWorker() http.HandlerFunc {
	return stx.RegisterServiceWorker(SWAuthRetryPath, SWAuthRetryScope, "RegisterAuthRetryWorker")
}

func (stx *ServiceContext) RegisterServiceWorker(swPath, swScope, functionTrace string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stx.RequestTracer.UpdateRequestTrace(r, functionTrace)
		redirectURL := r.URL.Query().Get("redirectURL")
		if redirectURL == "" {
			// TODO: fatal
		}
		requestID := r.URL.Query().Get("requestID")
		if requestID == "" {
			// TODO: fatal
		}
		cachedRequest, err := stx.RequestStore.RetrieveCachedRequest(requestID)
		if err != nil {
			// TODO: could not get cached request
		}

		pages.RegisterServiceWorker(swPath, swScope, cachedRequest).Render(r.Context(), w)
	})
}
