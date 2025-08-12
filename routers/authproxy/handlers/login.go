package handlers

import (
	"net/http"
	"net/url"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/helpers"
	"github.com/regcomp/gdpr/logging"
	"github.com/regcomp/gdpr/templates/pages"
)

// LoginPage Will likely need some context from the config for what to display
func LoginPage(w http.ResponseWriter, r *http.Request) {
	logging.RT.UpdateRequestTrace(r, "LoginPage")
	page := pages.Login()
	page.Render(r.Context(), w)
}

func SubmitLoginCredentials(authProvider auth.IAuthProvider, configStore config.IConfigStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "SubmitLoginCredentials")
		urlWithPort, err := configStore.GetServiceURLWithPort()
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}
		urlString := urlWithPort + config.PathAuthLoginCallback
		callbackURL, err := url.Parse(urlString)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}
		authProvider.AuthenticateUser(w, r, callbackURL)
	})
}
