package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	sc "github.com/gorilla/securecookie"
	"github.com/regcomp/gdpr/caching"
	"github.com/regcomp/gdpr/constants"
)

var cookieNames = []string{
	constants.AccessCookieName,
	constants.RefreshCookieName,
	constants.SessionCookieName,
}

type cookieOption func(*http.Cookie)

type CookieManager struct {
	// cache caching.IServiceCache
	keys *sc.SecureCookie
}

type cookieKeys struct {
	Hash  []byte `json:"hash"`
	Block []byte `json:"block"`
}

func newCookieKeys() *cookieKeys {
	return &cookieKeys{
		Hash:  sc.GenerateRandomKey(64),
		Block: sc.GenerateRandomKey(32),
	}
}

func CreateCookieManager(serviceCache caching.IServiceCache) *CookieManager {
	cookieHashes := serviceCache.CookieHashesGet()
	var freshKeys *cookieKeys
	if cookieHashes == nil {
		freshKeys = newCookieKeys()
		cookieHashesBytes, err := json.Marshal(freshKeys)
		if err != nil {
			// TODO:
		}
		serviceCache.CookieHashesSet(cookieHashesBytes)
	} else {
		err := json.Unmarshal(cookieHashes, freshKeys)
		if err != nil {
			// TODO:
		}
	}

	keys := sc.New(freshKeys.Hash, freshKeys.Block)
	return &CookieManager{
		// cache: serviceCache,
		keys: keys,
	}
}

func (cm *CookieManager) CreateAccessCookie(accessToken string) (*http.Cookie, error) {
	return createCookie(
		constants.AccessCookieName,
		accessToken,
		cm.keys,
		// TODO: Configure
	)
}

func (cm *CookieManager) GetAccessToken(r *http.Request) (string, error) {
	return getTokenFromCookie(r, constants.AccessCookieName, cm.keys)
}

func (cm *CookieManager) CreateRefreshCookie(refreshToken string) (*http.Cookie, error) {
	return createCookie(
		constants.RefreshCookieName,
		refreshToken,
		cm.keys,
		// TODO: Configure
		func(c *http.Cookie) {
			c.Path = constants.RouterAuthPathPrefix + constants.EndpointRenewToken
		},
	)
}

func (cm *CookieManager) GetRefreshToken(r *http.Request) (string, error) {
	return getTokenFromCookie(r, constants.RefreshCookieName, cm.keys)
}

func (cm *CookieManager) CreateSessionCookie(sessionID string) (*http.Cookie, error) {
	return createCookie(
		constants.SessionCookieName,
		sessionID,
		cm.keys,
		// TODO: Configure
	)
}

func (cm *CookieManager) GetSessionID(r *http.Request) (string, error) {
	return getTokenFromCookie(r, constants.SessionCookieName, cm.keys)
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
