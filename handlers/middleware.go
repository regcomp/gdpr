package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
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
			stx.RequestTracer.UpdateActiveTrace("VerifyServiceWorkerIsRunning")
			if strings.HasPrefix(r.URL.Path, "/static/sw") {
				next.ServeHTTP(w, r)
				return
			}
			if r.Header.Get(swHeader) == "" {
				stx.RegisterServiceWorker(swPath, swScope).ServeHTTP(w, r)
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
		stx.RequestTracer.UpdateActiveTrace("IsAuthenticated")
		accessToken, err := auth.GetAccessToken(r, stx.CookieKeys)
		if err != nil {
			auth.DestroyAllCookies(r)
			http.Error(w, "access token required", http.StatusUnauthorized)
			return
		}

		claims, err := stx.AuthProvider.ValidateAccessToken(accessToken)
		if err != nil {
			w.Header().Add("Refresh-Access-Token", "true")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// TODO: Add the claims to the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, claimsContextKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (stx *ServiceContext) ScopeServiceWorkerAccess(swPath, accessPath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			stx.RequestTracer.UpdateActiveTrace("ScopeServiceWorkerAccess")
			if r.URL.Path == swPath {
				w.Header().Add("Service-Worker-Allowed", accessPath)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (stx *ServiceContext) SetHSTSPolicy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stx.RequestTracer.UpdateActiveTrace("SetHSTSPolicy")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		next.ServeHTTP(w, r)
	})
}

func (stx *ServiceContext) TraceRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cw := &logging.CustomWriter{ResponseWriter: w, Code: http.StatusOK}
		stx.RequestTracer.NewRequestTrace(r)
		next.ServeHTTP(cw, r)
		stx.RequestTracer.DumpActiveTrace()
		fmt.Printf("[RESPONSE HEADER %d]\n", cw.Code)
		cw.Header().Write(os.Stdout)
		fmt.Printf("[BODY]\n%s\n", cw.Body.String())
		fmt.Println("")
	})
}
