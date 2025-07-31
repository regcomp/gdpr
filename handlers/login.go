package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/logging"
	"github.com/regcomp/gdpr/templates/pages"
)

// Will likely need some context from the config for what to display
func LoginPage(w http.ResponseWriter, r *http.Request) {
	logging.RT.UpdateRequestTrace(r, "LoginPage")
	page := pages.Login()
	page.Render(r.Context(), w)
}

func SubmitLoginCredentials(authProvider auth.IAuthProvider, configStore config.IConfigStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "SubmitLoginCredentials")
		urlString := configStore.GetServiceURLWithPort() + config.PathAuthLoginCallback
		callbackURL, err := url.Parse(urlString)
		if err != nil {
			// TODO:
		}
		authProvider.AuthenticateUser(w, r, callbackURL)
	})
}

func LoginCallback(
	authProvider auth.IAuthProvider,
	cookieManager *auth.CookieManager,
	sessionStore auth.ISessionStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "LoginCallback")
		credentials := auth.Credentials{}

		switch authProvider.GetProviderType() {
		// NOTE: Vendor implementations go here
		case auth.MockProviderType:
			credentials.AccessToken = r.URL.Query().Get(config.QueryParamAccessToken)
			credentials.RefreshToken = r.URL.Query().Get(config.QueryParamRefreshToken)
		default:
			http.Error(w, "auth provider not implemented", http.StatusInternalServerError)
		}

		// TODO: VALIDATE THE JWTS RECIEVED

		accessCookie, err := cookieManager.CreateAccessCookie(credentials.AccessToken)
		if err != nil {
			log.Panic("could not create access cookie")
		}
		http.SetCookie(w, accessCookie)

		refreshCookie, err := cookieManager.CreateRefreshCookie(credentials.RefreshToken)
		if err != nil {
			// TODO:
		}
		http.SetCookie(w, refreshCookie)

		sessionID := sessionStore.CreateSession()
		sessionCookie, err := cookieManager.CreateSessionCookie(sessionID)
		if err != nil {
			// TODO:
		}
		http.SetCookie(w, sessionCookie)

		// NOTE: This redirect may want to instead reference where a user was when a refresh token expired.
		http.Redirect(w, r, config.EndpointDashboard, http.StatusSeeOther)
	}
}

func RenewAccessToken(authProvider auth.IAuthProvider, cookieManager *auth.CookieManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "RenewAccessToken")

		refreshToken, err := cookieManager.GetRefreshToken(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not get refresh token, err=%s", err.Error()), http.StatusInternalServerError)
		}
		accessToken, err := authProvider.GetNewAccessToken(refreshToken, r)
		if err != nil {
			http.Error(w, "could not renew access token", http.StatusInternalServerError)
		}

		accessCookie, err := cookieManager.CreateAccessCookie(accessToken)
		if err != nil {
			http.Error(w, "could not create access cookie", http.StatusInternalServerError)
		}

		http.SetCookie(w, accessCookie)
		w.WriteHeader(http.StatusOK)
	}
}

func Logout(cookieManager *auth.CookieManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "PostLogout")
		cookieManager.DestroyAllCookies(w, r)
		w.WriteHeader(http.StatusOK)
	}
}
