package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/caching"
	"github.com/regcomp/gdpr/constants"
	"github.com/regcomp/gdpr/logging"
)

func RequestLogging(requestLogger logging.ILogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/favicon.ico" {
				requestLogger.Log(r)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func VerifyAuthRetryIsRunning(requestStore caching.IRequestStore) func(http.Handler) http.Handler {
	return verifyServiceWorkerIsRunning(
		constants.AuthRetryWorkerPath,
		constants.AuthRetryWorkerScope,
		constants.HeaderAuthRetryWorkerRunning,
		requestStore,
	)
}

func verifyServiceWorkerIsRunning(
	swPath, swScope, swHeader string,
	requestStore caching.IRequestStore,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.RT.UpdateRequestTrace(r, "VerifyServiceWorkerIsRunning")
			if strings.HasPrefix(r.URL.Path, "/static/sw") {
				next.ServeHTTP(w, r)
				return
			}
			if r.Header.Get(swHeader) != "true" {
				// cache request
				requestID, err := requestStore.StoreRequest(r)
				if err != nil {
					// TODO:
				}

				// construct registration url
				registrationURL := fmt.Sprintf("%s?%s=%s&%s=%s&%s=%s",
					constants.EndpointRegisterServiceWorker,
					constants.QueryParamRequestID, requestID,
					constants.QueryParamSWPath, swPath,
					constants.QueryParamSWScope, swScope,
				)

				// redirect to the url
				http.Redirect(w, r, registrationURL, http.StatusTemporaryRedirect)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func IsAuthenticated(authProvider auth.IAuthProvider, cookieManager *auth.CookieManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.RT.UpdateRequestTrace(r, "IsAuthenticated")
			accessToken, err := cookieManager.GetAccessToken(r)
			if err != nil {
				http.Error(w, "requires authentication", http.StatusUnauthorized)
				return
			}

			claims, err := authProvider.ValidateAccessToken(accessToken)
			if err != nil {
				w.Header().Add(constants.HeaderRenewAccessToken, "true")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// TODO: Add the claims to the request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, constants.ContextKeyClaims, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func HasActiveSession(sessionStore auth.ISessionStore, cookieManager *auth.CookieManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.RT.UpdateRequestTrace(r, "HasActiveSession")
			sessionID, err := cookieManager.GetSessionID(r)
			if err != nil {
				// TODO:
				// No session cookie
				log.Panic("No session cookie")
			}
			_, err = sessionStore.GetSession(sessionID)
			if err != nil {
				// TODO:
				// no registered session. old cookie?
				log.Panic("sessionID not found")
			}
			ctx := context.WithValue(r.Context(), constants.ContextKetSessionID, sessionID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AddNonceToRequest(nonceStore *auth.NonceStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.RT.UpdateRequestTrace(r, "AddNonceToRequest")
			nonce := nonceStore.Generate()
			// adding to response header
			w.Header().Set(constants.HeaderNonceToken, nonce)
			// adding to context so it can be passed to templates
			r = r.WithContext(context.WithValue(r.Context(), constants.ContextKeyNonceToken, nonce))

			next.ServeHTTP(w, r)
		})
	}
}

func ScopeAuthRetryAccess() func(http.Handler) http.Handler {
	return ScopeServiceWorkerAccess(constants.AuthRetryWorkerPath, constants.AuthRetryWorkerScope)
}

func ScopeServiceWorkerAccess(swPath, accessPath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.RT.UpdateRequestTrace(r, "ScopeServiceWorkerAccess")
			if r.URL.Path == swPath {
				w.Header().Add(constants.HeaderServiceWorkerAllowed, accessPath)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func SetHSTSPolicy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "SetHSTSPolicy")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		next.ServeHTTP(w, r)
	})
}

func TraceRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cw := &logging.CustomWriter{ResponseWriter: w, Code: http.StatusOK}
		logging.RT.NewRequestTrace(r)
		next.ServeHTTP(cw, r)
		logging.RT.DumpRequestTrace(r)
		fmt.Printf("[RESPONSE HEADER %d]\n", cw.Code)
		cw.Header().Write(os.Stdout)
		fmt.Printf("[BODY]\n%s\n", cw.Body.String())
		fmt.Println("")
	})
}
