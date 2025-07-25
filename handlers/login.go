package handlers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/constants"
	"github.com/regcomp/gdpr/logging"
	"github.com/regcomp/gdpr/templates/pages"
)

// Will likely need some context from the config for what to display
func LoginPage(w http.ResponseWriter, r *http.Request) {
	logging.RT.UpdateRequestTrace(r, "LoginPage")
	page := pages.Login()
	page.Render(r.Context(), w)
}

func SubmitLoginCredentials(authProvider auth.IAuthProvider, config config.IConfigStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "SubmotLoginCredentials")
		urlString := config.GetServiceURL() + constants.PathAuthRenewToken
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
		// TODO: Vendor implementations go here
		case auth.MOCK:
			credentials.AccessToken = r.URL.Query().Get(constants.QueryParamAccessToken)
			credentials.RefreshToken = r.URL.Query().Get(constants.QueryParamRefreshToken)
		default:
			http.Error(w, "auth provider not implemented", http.StatusInternalServerError)
		}

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
		http.Redirect(w, r, constants.EndpointDashboard, http.StatusSeeOther)
	}
}

func RenewAccessToken(authProvider auth.IAuthProvider, cookieManager *auth.CookieManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "RenewAccessToken")

		refreshToken, err := cookieManager.GetRefreshToken(r)
		if err != nil {
			// TODO:
		}
		accessToken, err := authProvider.GetNewAccessToken(refreshToken, r)
		if err != nil {
			// TODO:
		}

		accessCookie, err := cookieManager.CreateAccessCookie(accessToken)
		if err != nil {
			// TODO:
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
