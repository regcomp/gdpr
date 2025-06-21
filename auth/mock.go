package auth

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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

	log.Println("hit authenticate user")
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
