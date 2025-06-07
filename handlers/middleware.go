package handlers

import (
	"net/http"
)

func (stx *ServiceContext) HasAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: integrate auth checks

		next.ServeHTTP(w, r)
	})
}
