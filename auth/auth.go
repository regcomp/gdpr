package auth

import (
	"net/http"
	"net/url"
)

const providerConfigString = "AUTH_PROVIDER"

type ProviderType int

const (
	AUTH0 ProviderType = iota
	OKTA
	//
	MOCK
)

// The shape of this may change. This is the struct that Auth responses will be converted into for the
// service to manage auth
type Credentials struct {
	UserId       string
	RefreshToken string
	AccessToken  string
}

type Provider interface {
	GetProviderType() ProviderType
	AuthenticateUser(http.ResponseWriter, *http.Request, url.URL) // NOTE: This may require more fields
	HasValidAuthentication(*http.Request) bool
}

func GetProvider(getenv func(string) string) (Provider, error) {
	provider := getenv(providerConfigString)
	switch provider {
	default:
		return createMockAuthProvider(), nil
	}
}
