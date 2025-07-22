package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

func (stx *ServiceContext) VerifyAuthRetryIsRunning() func(http.Handler) http.Handler {
	return stx.VerifyServiceWorkerIsRunning(SWAuthRetryPath, SWAuthRetryScope, "SW-Auth-Retry-Running")
}

func (stx *ServiceContext) VerifyServiceWorkerIsRunning(swPath, swScope, swHeader string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			stx.RequestTracer.UpdateRequestTrace(r, "VerifyServiceWorkerIsRunning")
			if strings.HasPrefix(r.URL.Path, "/static/sw") {
				next.ServeHTTP(w, r)
				return
			}
			if r.Header.Get(swHeader) != "true" {
				// cache request
				requestID, err := stx.RequestStore.StoreCachedRequest(r)
				if err != nil {
					// TODO:
				}

				// construct callback url
				bootstrapURL := fmt.Sprintf("%s?redirect=%s&requestID=%s",
					RegisterServiceWorkerPath,
					url.QueryEscape(r.URL.String()),
					requestID,
				)

				// redirect to the url
				http.Redirect(w, r, bootstrapURL, http.StatusTemporaryRedirect)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// TODO: Figure out where this should go

type ContextKey string

const (
	claimsContextKey   ContextKey = "claims"
	sessionIDContexKey ContextKey = "session-id"
)

// -----

func (stx *ServiceContext) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stx.RequestTracer.UpdateRequestTrace(r, "IsAuthenticated")
		accessToken, err := auth.GetAccessToken(r, stx.CookieKeys)
		if err != nil {
			stx.RequestTracer.UpdateRequestTrace(r, "DestroyAllCookies")
			http.Error(w, "requires authentication", http.StatusUnauthorized)
			return
		}

		stx.RequestTracer.UpdateRequestTrace(r, "ValidateAccessToken")
		claims, err := stx.AuthProvider.ValidateAccessToken(accessToken)
		if err != nil {
			stx.RequestTracer.UpdateRequestTrace(r, "Invalid Access Token")
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

func (stx *ServiceContext) HasActiveSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stx.RequestTracer.UpdateRequestTrace(r, "HasActiveSession")
		sessionID, err := auth.GetSessionID(r, stx.CookieKeys)
		if err != nil {
			// TODO:
			// No session cookie
			log.Panic("No session cookie")
		}
		_, err = stx.SessionStore.GetSession(sessionID)
		if err != nil {
			// TODO:
			// no registered session. old cookie?
			log.Panic("sessionID not found")
		}
		ctx := context.WithValue(r.Context(), sessionIDContexKey, sessionID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (stx *ServiceContext) AddNonceToRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonce := stx.NonceStore.Generate()
		// adding to response header
		w.Header().Set("XSRF-Nonce", nonce)
		// adding to context so it can be passed to templates
		r = r.WithContext(context.WithValue(r.Context(), "nonce", nonce))

		next.ServeHTTP(w, r)
	})
}

func (stx *ServiceContext) ScopeAuthRetryAccess() func(http.Handler) http.Handler {
	return stx.ScopeServiceWorkerAccess(SWAuthRetryPath, SWAuthRetryScope)
}

func (stx *ServiceContext) ScopeServiceWorkerAccess(swPath, accessPath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			stx.RequestTracer.UpdateRequestTrace(r, "ScopeServiceWorkerAccess")
			if r.URL.Path == swPath {
				w.Header().Add("Service-Worker-Allowed", accessPath)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (stx *ServiceContext) SetHSTSPolicy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stx.RequestTracer.UpdateRequestTrace(r, "SetHSTSPolicy")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		next.ServeHTTP(w, r)
	})
}

func (stx *ServiceContext) TraceRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cw := &logging.CustomWriter{ResponseWriter: w, Code: http.StatusOK}
		stx.RequestTracer.NewRequestTrace(r)
		next.ServeHTTP(cw, r)
		stx.RequestTracer.DumpRequestTrace(r)
		fmt.Printf("[RESPONSE HEADER %d]\n", cw.Code)
		cw.Header().Write(os.Stdout)
		fmt.Printf("[BODY]\n%s\n", cw.Body.String())
		fmt.Println("")
	})
}
