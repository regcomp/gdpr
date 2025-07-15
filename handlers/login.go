package handlers

import (
	"log"
	"net/http"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/templates/pages"
)

func (stx *ServiceContext) GetLogin(w http.ResponseWriter, r *http.Request) {
	stx.RequestTracer.UpdateRequestTrace(r, "GetLogin")
	page := pages.Login()
	page.Render(r.Context(), w)
}

func (stx *ServiceContext) PostLogin(w http.ResponseWriter, r *http.Request) {
	stx.RequestTracer.UpdateRequestTrace(r, "PostLogin")
	callbackURL := NewURL("https", stx.HostPath, AuthRouterPathPrefix+LoginCallbackPath)
	stx.AuthProvider.AuthenticateUser(w, r, callbackURL)
}

func (stx *ServiceContext) LoginCallback(w http.ResponseWriter, r *http.Request) {
	stx.RequestTracer.UpdateRequestTrace(r, "LoginCallback")
	credentials := auth.Credentials{}

	switch stx.AuthProvider.GetProviderType() {
	// TODO: Vendor implementations go here
	case auth.MOCK:
		credentials.AccessToken = r.URL.Query().Get("access")
		credentials.RefreshToken = r.URL.Query().Get("refresh")
	default:
		http.Error(w, "auth provider not implemented", http.StatusInternalServerError)
	}

	accessCookie, err := auth.CreateAccessCookie(credentials.AccessToken, stx.CookieKeys)
	if err != nil {
		log.Panic("could not create access cookie")
	}
	http.SetCookie(w, accessCookie)

	refreshCookie, err := auth.CreateRefreshCookie(credentials.RefreshToken, stx.CookieKeys)
	if err != nil {
		// TODO:
	}
	http.SetCookie(w, refreshCookie)

	sessionID := stx.SessionStore.CreateSession()
	sessionCookie, err := auth.CreateSessionCookie(sessionID, stx.CookieKeys)
	if err != nil {
		// TODO:
	}
	http.SetCookie(w, sessionCookie)

	csrfToken := auth.NewCSRFToken(sessionID, stx.HMACSecret)
	csrfCookie, err := auth.CreateCSRFCookie(csrfToken, stx.CookieKeys)
	if err != nil {
		// TODO:
	}
	http.SetCookie(w, csrfCookie)

	// NOTE: This redirect may want to instead reference where a user was when a refresh token expired.
	http.Redirect(w, r, DashboardPath, http.StatusSeeOther)
}

func (stx *ServiceContext) PostRefresh(w http.ResponseWriter, r *http.Request) {
	stx.RequestTracer.UpdateRequestTrace(r, "PostRefresh")

	refreshToken, err := auth.GetRefreshToken(r, stx.CookieKeys)
	if err != nil {
		// TODO:
	}
	accessToken, err := stx.AuthProvider.GetNewAccessToken(refreshToken, r)
	if err != nil {
		// TODO:
	}

	accessCookie, err := auth.CreateAccessCookie(accessToken, stx.CookieKeys)
	if err != nil {
		// TODO:
	}

	http.SetCookie(w, accessCookie)
	w.WriteHeader(http.StatusOK)
}

func (stx *ServiceContext) PostLogout(w http.ResponseWriter, r *http.Request) {
	stx.RequestTracer.UpdateRequestTrace(r, "PostLogout")
	auth.DestroyAllCookies(r)
	http.Redirect(w, r, LoginPath, http.StatusSeeOther)
}
