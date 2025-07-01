package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/logging"
)

func (stx *ServiceContext) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/favicon.ico" {
			logging.LogRequest(stx.RequestLogger, r)
		}
		next.ServeHTTP(w, r)
	})
}

func (stx *ServiceContext) VerifyServiceWorkerIsRunning(swPath, swScope, swHeader string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/static/sw/") {
				next.ServeHTTP(w, r)
				return
			}
			if r.Header.Get(swHeader) == "" {
				RegisterServiceWorker(swPath, swScope).ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// TODO: Figure out where this should go

type ContextKey string

const claimsContextKey ContextKey = "claims"

// -----

func (stx *ServiceContext) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := auth.GetAccessToken(r, stx.CookieKeys)
		if err != nil {
			auth.DestroyAllCookies(r)
			http.Error(w, "access token required", http.StatusUnauthorized)
			return
		}

		claims, err := stx.AuthProvider.ValidateAccessToken(accessToken)
		if err != nil {
			w.Header().Add("Refresh-Access-Token", "true")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// TODO: Add the claims to the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, claimsContextKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ScopeServiceWorkerAccess(swPath, accessPath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == swPath {
				w.Header().Add("Service-Worker-Allowed", accessPath)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (stx *ServiceContext) SetHSTSPolicy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		next.ServeHTTP(w, r)
	})
}
