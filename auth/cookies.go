package auth

import (
	"fmt"
	"net/http"

	sc "github.com/gorilla/securecookie"
)

// Cookies
const (
	AccessCookieName  = "access-token"
	RefreshCookieName = "refresh-token"
	SessionCookieName = "session-id"
)

var cookieNames = []string{
	AccessCookieName,
	RefreshCookieName,
	SessionCookieName,
}

type cookieOption func(*http.Cookie)

// NOTE: this is too implementation specific
type ICookieCrypt interface {
	Keys() *sc.SecureCookie
}

type CookieCrypt struct {
	keys *sc.SecureCookie
}

func CreateCookieCrypt() *CookieCrypt {
	hashKey := sc.GenerateRandomKey(64)
	blockKey := sc.GenerateRandomKey(32)
	return &CookieCrypt{keys: sc.New(hashKey, blockKey)}
}

func (cc *CookieCrypt) Keys() *sc.SecureCookie {
	return cc.keys
}

func DecodeCookie(
	name string,
	sc *sc.SecureCookie,
	encryptedCookie *http.Cookie,
) (map[string]string, error) {
	value := make(map[string]string)
	if err := sc.Decode(name, encryptedCookie.Value, &value); err != nil {
		return nil, err
	}
	return value, nil
}

func createCookie(
	name, value string,
	sc *sc.SecureCookie,
	options ...cookieOption,
) (*http.Cookie, error) {
	values := map[string]string{
		name: value,
	}

	encoded, err := sc.Encode(name, values)
	if err != nil {
		return nil, fmt.Errorf("could not encode cookie values: %s, %s=%s", name, value, err.Error())
	}

	cookie := &http.Cookie{
		Name:  name,
		Value: encoded,

		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	for _, option := range options {
		option(cookie)
	}

	return cookie, nil
}

func DestroyAllCookies(w http.ResponseWriter, r *http.Request) {
	for _, cookieName := range cookieNames {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			continue
		}
		destroyCookie(w, cookie)
	}
}

func destroyCookie(w http.ResponseWriter, cookie *http.Cookie) {
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}

func getTokenFromCookie(
	name string,
	r *http.Request,
	sc *sc.SecureCookie,
) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("could not find %s cookie=%s", name, err.Error())
	}
	decodedValues, err := DecodeCookie(name, sc, cookie)
	if err != nil {
		return "", fmt.Errorf("could not decode %s cookie=%s", name, err.Error())
	}

	token, ok := decodedValues[name]
	if !ok {
		return "", fmt.Errorf("could not find %s token=%s", name, err.Error())
	}

	return token, nil
}

func CreateAccessCookie(accessToken string, sc *sc.SecureCookie) (*http.Cookie, error) {
	return createCookie(
		AccessCookieName,
		accessToken,
		sc,
		// TODO: Configure
	)
}

func GetAccessToken(r *http.Request, sc *sc.SecureCookie) (string, error) {
	return getTokenFromCookie(AccessCookieName, r, sc)
}

func CreateRefreshCookie(refreshToken string, sc *sc.SecureCookie) (*http.Cookie, error) {
	return createCookie(
		RefreshCookieName,
		refreshToken,
		sc,
		// TODO: Configure
		func(c *http.Cookie) {
			c.Path = "/auth/refresh/"
		},
	)
}

func GetRefreshToken(r *http.Request, sc *sc.SecureCookie) (string, error) {
	return getTokenFromCookie(RefreshCookieName, r, sc)
}

func CreateSessionCookie(sessionID string, sc *sc.SecureCookie) (*http.Cookie, error) {
	return createCookie(
		SessionCookieName,
		sessionID,
		sc,
		// TODO: Configure
	)
}

func GetSessionID(r *http.Request, sc *sc.SecureCookie) (string, error) {
	return getTokenFromCookie(SessionCookieName, r, sc)
}
