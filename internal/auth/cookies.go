package auth

import (
	"fmt"
	"net/http"

	sc "github.com/gorilla/securecookie"
)

type cookieOption func(*http.Cookie)

func CreateCookie(
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

func GetTokenFromCookie(
	r *http.Request,
	name string,
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

func DestroyCookie(w http.ResponseWriter, cookie *http.Cookie) {
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}
