package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/regcomp/gdpr/internal/auth"
	"github.com/regcomp/gdpr/internal/caching"
	"github.com/regcomp/gdpr/internal/config"
	"github.com/regcomp/gdpr/pkg/helpers"
	"github.com/regcomp/gdpr/pkg/logging"
)

func SkipIfAuthenticated(ap auth.IAuthProvider, cm *caching.CookieManager, cs config.IConfigStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.RT.UpdateRequestTrace(r, "SkipifAuthenticated")
			accessToken, err := cm.GetAccessToken(r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			_, err = ap.ValidateAccessToken(accessToken)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			urlWithPort, err := cs.GetServiceURLWithPort()
			if err != nil {
				helpers.RespondWithError(w, err, http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r,
				urlWithPort+config.PathClientDashboard,
				http.StatusSeeOther,
			)
		})
	}
}

func RequiresAuthentication(ap auth.IAuthProvider, cm *caching.CookieManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.RT.UpdateRequestTrace(r, "RequiresAuthentication")
			accessToken, err := cm.GetAccessToken(r)
			if err != nil {
				helpers.RespondWithError(w, err, http.StatusUnauthorized)
				return
			}

			claims, err := ap.ValidateAccessToken(accessToken)
			if err != nil {
				w.Header().Add(config.HeaderRenewAccessToken, config.ValueTrue)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// TODO: Add the claims to the request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, config.ContextKeyClaims, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func HasActiveSession(sm *caching.SessionManager, cm *caching.CookieManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.RT.UpdateRequestTrace(r, "HasActiveSession")
			sessionID, err := cm.GetSessionID(r)
			if err != nil {
				// TODO:
				// No session cookie
				log.Panic("No session cookie")
			}
			_, err = sm.GetSession(sessionID)
			if err != nil {
				// TODO:
				// no registered session. old cookie?
				log.Panic("sessionID not found")
			}
			ctx := context.WithValue(r.Context(), config.ContextKeySessionId, sessionID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AddNonceToRequest(nm *caching.NonceManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.RT.UpdateRequestTrace(r, "AddNonceToRequest")
			nonce := nm.Generate()
			// adding to response header
			w.Header().Set(config.HeaderNonceToken, nonce)
			// adding to context so it can be passed to templates
			r = r.WithContext(context.WithValue(r.Context(), config.ContextKeyNonceToken, nonce))

			next.ServeHTTP(w, r)
		})
	}
}
