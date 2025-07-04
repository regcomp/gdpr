package auth

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
)

// WARN: reusable module level keys. does not appear to have a use anymore
var (
	hashKey        []byte = nil
	blockKey       []byte = nil
	hasInitialized        = false
)

const (
	AccessTokenString  = "access-token"
	RefreshTokenString = "refresh-token"
	SessionIDString    = "session-id"
)

type cookieOption func(*http.Cookie)

func createSecrets() {
	if hasInitialized {
		return
	}
	hashKey = securecookie.GenerateRandomKey(64)
	blockKey = securecookie.GenerateRandomKey(32)
	hasInitialized = true
}

func CreateCookieKeys() *securecookie.SecureCookie {
	createSecrets()
	return securecookie.New(hashKey, blockKey)
}

func DecodeCookie(
	name string,
	sc *securecookie.SecureCookie,
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
	sc *securecookie.SecureCookie,
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
	}

	for _, option := range options {
		option(cookie)
	}

	return cookie, nil
}

func getTokenFromCookie(
	name string,
	r *http.Request,
	sc *securecookie.SecureCookie,
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

func CreateAccessCookie(accessToken string, sc *securecookie.SecureCookie) (*http.Cookie, error) {
	return createCookie(
		AccessTokenString,
		accessToken,
		sc,
		// TODO: Configure
	)
}

func GetAccessToken(r *http.Request, sc *securecookie.SecureCookie) (string, error) {
	return getTokenFromCookie(AccessTokenString, r, sc)
}

func CreateRefreshCookie(refreshToken string, sc *securecookie.SecureCookie) (*http.Cookie, error) {
	return createCookie(
		RefreshTokenString,
		refreshToken,
		sc,
		// TODO: Configure
	)
}

func GetRefreshToken(r *http.Request, sc *securecookie.SecureCookie) (string, error) {
	return getTokenFromCookie(RefreshTokenString, r, sc)
}

func CreateSessionCookie(sc *securecookie.SecureCookie) (*http.Cookie, error) {
	return createCookie(
		SessionIDString,
		generateSessionID(),
		sc,
		// TODO: Configure
	)
}

func GetSessionID(r *http.Request, sc *securecookie.SecureCookie) (string, error) {
	return getTokenFromCookie(SessionIDString, r, sc)
}
