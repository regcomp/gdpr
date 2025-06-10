package auth

import "net/http"

const authProviderConfigString = "AUTH_PROVIDER"

type Credentials struct {
	UserId       string
	RefreshToken string
	AccessToken  string
}

type Provider interface {
	AuthenticateUser(*http.Request) (Credentials, error)
}

func GetProvider(getenv func(string) string) (Provider, error) {
	provider := getenv(authProviderConfigString)
	switch provider {
	default:
		return createMockProvider(), nil
	}
}
