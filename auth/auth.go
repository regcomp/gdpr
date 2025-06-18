package auth

import (
	"encoding/json"
	"log"
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
	SessionID    string `json:"session_id"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type Provider interface {
	GetProviderType() ProviderType
	AuthenticateUser(http.ResponseWriter, *http.Request, url.URL) // NOTE: This may require more fields
	IsValidAccessToken(string) bool
	GetNewAccessToken(string) (string, error)
}

func GetProvider(getenv func(string) string) (Provider, error) {
	provider := getenv(providerConfigString)
	switch provider {
	default:
		return createMockAuthProvider(), nil
	}
}

func FillCredentialsFromRequestBody(r *http.Request, credentials *Credentials) {
	err := json.NewDecoder(r.Body).Decode(credentials)
	if err != nil {
		log.Panicf("could not decode request body: %s", err.Error())
	}
}
