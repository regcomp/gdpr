package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/auth"
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
		accessToken, err := auth.GetAccessToken(r, stx.CookieKeys)
		if err != nil {
			// TODO: No access token cookie, kill all cookies, kick to login
		}

		// WARN: This whole section is suspect. Not all requests should include the refresh
		// token in cookies. If the access token is invalid, a response should be issued
		// that triggers a request to the endpoint which refresh tokens are sent to
		if !stx.AuthProvider.IsValidAccessToken(accessToken) {
			refreshToken, err := auth.GetRefreshToken(r, stx.CookieKeys)
			if err != nil {
				// TODO: A RESPONSE TO THE CLIENT TO BE SENT THAT SHOULD TRIGGER
				// A REQUEST TO THE YET TO BE DEFINED REFRESH ENDPOINT
			}
			accessToken, err := stx.AuthProvider.GetNewAccessToken(refreshToken)
			if err != nil {
				// TODO: Unable to refresh, kill cookies, kick to login
			}
			accessCookie := auth.CreateAccessCookie(accessToken, stx.CookieKeys)
			http.SetCookie(w, accessCookie)
		}

		next.ServeHTTP(w, r)
	})
}
