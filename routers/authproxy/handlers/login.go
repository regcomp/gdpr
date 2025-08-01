package handlers

import (
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
