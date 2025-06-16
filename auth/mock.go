package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	accessToken  = "access-token"
	refreshToken = "refresh-token"
	userID       = "user-id"
)

type MockProvider struct {
	Proxy *http.Server
}

func createMockAuthProvider() *MockProvider {
	return &MockProvider{}
}

func (mp *MockProvider) GetProviderType() ProviderType {
	return MOCK
}

func (mp *MockProvider) AuthenticateUser(w http.ResponseWriter, r *http.Request, callback url.URL) {
	// spin up a mock auth proxy

	// send temporaryredirect to the auth proxy

	payload := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		UserID       string `json:"user_id"`
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       userID,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		// TODO:
	}

	body := bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, callback.String(), body)
	if err != nil {
		// TODO:
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		// TODO:
	}
}

func (mp *MockProvider) HasValidAuthentication(r *http.Request) bool {
	_, err := r.Cookie("access_token")
	return err == nil
}
