package auth

import (
	"fmt"
	"net/http"

	sc "github.com/gorilla/securecookie"
	"github.com/regcomp/gdpr/cache"
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

type CookieManager struct {
	serviceCache cache.IServiceCache
	keys         *sc.SecureCookie
}

func CreateCookieManager(serviceCache cache.IServiceCache) *CookieManager {
	// TODO: This data needs to be pulled from the cache or generated and added to it
	hashKey := sc.GenerateRandomKey(64)
	blockKey := sc.GenerateRandomKey(32)
	keys := sc.New(hashKey, blockKey)
	return &CookieManager{
		serviceCache: serviceCache,
		keys:         keys,
	}
}

func (cm *CookieManager) CreateAccessCookie(accessToken string) (*http.Cookie, error) {
	return createCookie(
		AccessCookieName,
		accessToken,
		cm.keys,
		// TODO: Configure
	)
}

func (cm *CookieManager) GetAccessToken(r *http.Request) (string, error) {
	return getTokenFromCookie(r, AccessCookieName, cm.keys)
}

func (cm *CookieManager) CreateRefreshCookie(refreshToken string) (*http.Cookie, error) {
	return createCookie(
		RefreshCookieName,
		refreshToken,
		cm.keys,
		// TODO: Configure
		func(c *http.Cookie) {
			c.Path = "/auth/refresh/"
		},
	)
}

func (cm *CookieManager) GetRefreshToken(r *http.Request) (string, error) {
	return getTokenFromCookie(r, RefreshCookieName, cm.keys)
}

func (cm *CookieManager) CreateSessionCookie(sessionID string) (*http.Cookie, error) {
	return createCookie(
		SessionCookieName,
		sessionID,
		cm.keys,
		// TODO: Configure
	)
}

func (cm *CookieManager) GetSessionID(r *http.Request) (string, error) {
	return getTokenFromCookie(r, SessionCookieName, cm.keys)
}

func (cm *CookieManager) DestroyAllCookies(w http.ResponseWriter, r *http.Request) {
	for _, cookieName := range cookieNames {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			continue
		}
		destroyCookie(w, cookie)
	}
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

func getTokenFromCookie(
	r *http.Request,
	name string,
	sc *sc.SecureCookie,
) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("could not find %s cookie=%s", name, err.Error())
	}
	decodedValues, err := decodeCookie(name, sc, cookie)
	if err != nil {
		return "", fmt.Errorf("could not decode %s cookie=%s", name, err.Error())
	}

	token, ok := decodedValues[name]
	if !ok {
		return "", fmt.Errorf("could not find %s token=%s", name, err.Error())
	}

	return token, nil
}

func decodeCookie(
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

func destroyCookie(w http.ResponseWriter, cookie *http.Cookie) {
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}
