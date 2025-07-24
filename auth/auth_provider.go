package auth

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

type IAuthProvider interface {
	GetProviderType() ProviderType
	AuthenticateUser(http.ResponseWriter, *http.Request, url.URL) // NOTE: This may require more fields
	ValidateAccessToken(string) (*CustomClaims, error)

	// NOTE: This likely doesnt need the request as information can be passed from
	// the access token
	GetNewAccessToken(string, *http.Request) (string, error)
}

func GetProvider(getenv func(string) string) (IAuthProvider, error) {
	provider := getenv(providerConfigString)
	switch provider {
	case "MOCK":
		return createMockAuthProvider(), nil
	default:
		return nil, fmt.Errorf("unknown auth provider=%s", provider)
	}
}

// func FillCredentialsFromRequestBody(r *http.Request, credentials *Credentials) {
// 	err := json.NewDecoder(r.Body).Decode(credentials)
// 	if err != nil {
// 		log.Panicf("could not decode request body: %s", err.Error())
// 	}
// }

type MockProvider struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

const (
	accessIssuer  = "mock-access"
	refreshIssuer = "mock-refresh"
)

func createMockAuthProvider() *MockProvider {
	mp := &MockProvider{}

	privateKey, err := generateRSAKeyPair()
	if err != nil {
		// TODO:
	}
	mp.privateKey = privateKey
	mp.publicKey = &privateKey.PublicKey

	return mp
}

func (mp *MockProvider) generateJWTWithClaims(audience, issuer string, tokenDuration time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   getUserID(),
		Audience:  []string{audience},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration * time.Second)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ID:        fmt.Sprintf("jwt-%d", time.Now().Unix()),
	}

	customClaims := CustomClaims{
		claims,
	}

	return generateJWTWithClaims(&customClaims, mp.privateKey)
}

func getUserID() string {
	return "user"
}

func (mp *MockProvider) GetProviderType() ProviderType {
	return MOCK
}

func (mp *MockProvider) AuthenticateUser(w http.ResponseWriter, r *http.Request, callback url.URL) {
	// NOTE: This assumes that there has been valid authentication and issues credentials

	redirectURL, err := url.Parse(callback.String())
	if err != nil {
		// TODO:
	}

	accessToken, err := mp.generateJWTWithClaims(r.URL.Hostname(), accessIssuer, 10)
	if err != nil {
		// TODO:
	}

	refreshToken, err := mp.generateJWTWithClaims(r.URL.Hostname(), refreshIssuer, 20)
	if err != nil {
		// TODO:
	}

	params := url.Values{}
	params.Add("access", accessToken)
	params.Add("refresh", refreshToken)

	redirectURL.RawQuery = params.Encode()
	http.Redirect(w, r, redirectURL.String(), http.StatusTemporaryRedirect)
}

func (mp *MockProvider) ValidateAccessToken(token string) (*CustomClaims, error) {
	return verifyJWTWithClaims(token, &CustomClaims{}, mp.publicKey)
}

func (mp *MockProvider) GetNewAccessToken(refreshToken string, r *http.Request) (string, error) {
	// TODO: Verify refresh token

	// -----

	accessToken, err := mp.generateJWTWithClaims(r.URL.Hostname(), accessIssuer, 10)
	if err != nil {
		// TODO:
	}
	return accessToken, nil
}
