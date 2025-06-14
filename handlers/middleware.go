package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/logging"
)

func (stx *ServiceContext) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.LogRequest(stx.RequestLogger, r)
		next.ServeHTTP(w, r)
	})
}

func (stx *ServiceContext) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !stx.AuthProvider.HasValidAuthentication(r) {
			http.Redirect(w, r, LoginPath, http.StatusSeeOther)
		}

		next.ServeHTTP(w, r)
	})
}

func (stx *ServiceContext) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO:

		next.ServeHTTP(w, r)
	})
}
