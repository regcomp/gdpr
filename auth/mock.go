package auth

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

const (
	sessionID    = "baz"
	refreshToken = "bar"
	accessToken  = "foo"
)

type MockProvider struct{}

func createMockAuthProvider() *MockProvider {
	return &MockProvider{}
}

func (mp *MockProvider) GetProviderType() ProviderType {
	return MOCK
}

func (mp *MockProvider) AuthenticateUser(w http.ResponseWriter, r *http.Request, callback url.URL) {
	payload := Credentials{
		SessionID:    sessionID,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Panicf("could not marshal credentials in mock provider: %s", err.Error())
	}

	r2, _ := http.NewRequest("GET", r.URL.String(), bytes.NewReader(data))
	r2.ContentLength = int64(len(data))

	r2.Header.Set("Content-Type", "application/json")

	http.Redirect(w, r2, callback.String(), http.StatusTemporaryRedirect)
}

func (mp *MockProvider) IsValidAccessToken(token string) bool {
	return true
}

func (mp *MockProvider) GetNewAccessToken(refreshToken string) (string, error) {
	return accessToken, nil
}
