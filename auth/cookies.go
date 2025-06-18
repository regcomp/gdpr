package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var (
	hashKey  []byte = nil
	blockKey []byte = nil
)

const (
	AccessTokenString  = "access-token"
	RefreshTokenString = "refresh-token"
)

func createSecrets() {
	hashKey = securecookie.GenerateRandomKey(64)
	blockKey = securecookie.GenerateRandomKey(32)
}

func CreateCookieKeys() *securecookie.SecureCookie {
	if hashKey == nil {
		createSecrets()
	}
	return securecookie.New(hashKey, blockKey)
}

func CreateSessionStore() *sessions.CookieStore {
	if hashKey == nil {
		createSecrets()
	}
	return sessions.NewCookieStore(hashKey, blockKey)
}

func GenerateSessionID() string {
	return uuid.NewString()
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

func CreateAccessCookie(accessToken string, sc *securecookie.SecureCookie) *http.Cookie {
	value := map[string]string{
		AccessTokenString: accessToken,
	}
	encoded, err := sc.Encode(AccessTokenString, value)
	if err != nil {
		// TODO:
		log.Panicf("encoding access token cookie: %s", err.Error())
	}
	return &http.Cookie{
		Name:  AccessTokenString,
		Value: encoded,
		// TODO: Configure
	}
}

func GetAccessToken(r *http.Request, sc *securecookie.SecureCookie) (string, error) {
	accessCookie, err := r.Cookie(AccessTokenString)
	if err != nil {
		return "", fmt.Errorf("no access cookie")
	}
	accessValues, err := DecodeCookie(AccessTokenString, sc, accessCookie)
	if err != nil {
		log.Panicf("error decoding access cookie: %s", err.Error())
	}

	token, ok := accessValues[AccessTokenString]
	if !ok {
		return "", fmt.Errorf("blank access token")
	}

	return token, nil
}

func CreateRefreshCookie(refreshToken string, sc *securecookie.SecureCookie) *http.Cookie {
	value := map[string]string{
		RefreshTokenString: refreshToken,
	}
	encoded, err := sc.Encode(RefreshTokenString, value)
	if err != nil {
		// TODO:
		log.Panicf("encoding refresh token cookie: %s", err.Error())
	}
	return &http.Cookie{
		Name:  RefreshTokenString,
		Value: encoded,
		// TODO: Configure
	}
}

func GetRefreshToken(r *http.Request, sc *securecookie.SecureCookie) (string, error) {
	accessCookie, err := r.Cookie(RefreshTokenString)
	if err != nil {
		return "", fmt.Errorf("no access cookie")
	}
	accessValues, err := DecodeCookie(RefreshTokenString, sc, accessCookie)
	if err != nil {
		log.Panicf("error decoding refresh cookie: %s", err.Error())
	}

	token, ok := accessValues[RefreshTokenString]
	if !ok {
		return "", fmt.Errorf("blank access token")
	}

	return token, nil
}
