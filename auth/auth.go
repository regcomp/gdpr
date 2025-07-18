package auth

import (
	"encoding/json"
	"fmt"
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

// NOTE: The shape of this may change. This is the struct that Auth responses will be converted into for the
// service to manage auth
type Credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	SessionID    string `json:"session_id"`
}

type IProvider interface {
	GetProviderType() ProviderType
	AuthenticateUser(http.ResponseWriter, *http.Request, url.URL) // NOTE: This may require more fields
	ValidateAccessToken(string) (*CustomClaims, error)

	// NOTE: This likely doesnt need the request as information can be passed from
	// the access token
	GetNewAccessToken(string, *http.Request) (string, error)
}

func GetProvider(getenv func(string) string) (IProvider, error) {
	provider := getenv(providerConfigString)
	switch provider {
	case "MOCK":
		return createMockAuthProvider(), nil
	default:
		return nil, fmt.Errorf("unknown auth provider=%s", provider)
	}
}

func FillCredentialsFromRequestBody(r *http.Request, credentials *Credentials) {
	err := json.NewDecoder(r.Body).Decode(credentials)
	if err != nil {
		log.Panicf("could not decode request body: %s", err.Error())
	}
}
