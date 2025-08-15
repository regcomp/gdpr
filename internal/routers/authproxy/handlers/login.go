package handlers

import (
	"net/http"
	"net/url"

	"github.com/regcomp/gdpr/internal/auth"
	"github.com/regcomp/gdpr/internal/config"
	"github.com/regcomp/gdpr/internal/views"
	"github.com/regcomp/gdpr/pkg/helpers"
	"github.com/regcomp/gdpr/pkg/logging"
)

// LoginPage Will likely need some context from the config for what to display
func LoginPage(w http.ResponseWriter, r *http.Request) {
	logging.RT.UpdateRequestTrace(r, "LoginPage")
	err := views.ServeLogin(w, r.Context())
	if err != nil {
		helpers.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}
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
