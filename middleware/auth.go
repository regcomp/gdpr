package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/constants"
	"github.com/regcomp/gdpr/logging"
)

func SkipIfAuthenticated(
	authProvider auth.IAuthProvider,
	cookieManager *auth.CookieManager,
	config config.IConfigStore,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.RT.UpdateRequestTrace(r, "SkipifAuthenticated")
			accessToken, err := cookieManager.GetAccessToken(r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			_, err = authProvider.ValidateAccessToken(accessToken)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			http.Redirect(w, r, config.GetServiceURL()+":"+config.GetDefaultPort()+constants.EndpointDashboard, http.StatusSeeOther)
		})
	}
}

func RequiresAuthentication(authProvider auth.IAuthProvider, cookieManager *auth.CookieManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.RT.UpdateRequestTrace(r, "RequiresAuthentication")
			accessToken, err := cookieManager.GetAccessToken(r)
			if err != nil {
				http.Error(w, "requires authentication", http.StatusUnauthorized)
				return
			}

			claims, err := authProvider.ValidateAccessToken(accessToken)
			if err != nil {
				w.Header().Add(constants.HeaderRenewAccessToken, constants.ValueTrueString)
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
