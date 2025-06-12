package handlers

import (
	"net/http"
)

func (stx *ServiceContext) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO:

		next.ServeHTTP(w, r)
	})
}

func (stx *ServiceContext) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO:

		next.ServeHTTP(w, r)
	})
}
