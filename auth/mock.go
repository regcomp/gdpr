package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	accessToken  = "foo"
	refreshToken = "bar"
	userID       = "baz"
)

type MockProvider struct{}

func createMockAuthProvider() *MockProvider {
	return &MockProvider{}
}

func (mp *MockProvider) GetProviderType() ProviderType {
	return MOCK
}

func (mp *MockProvider) AuthenticateUser(w http.ResponseWriter, r *http.Request, callback url.URL) {
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

	r.Body = io.NopCloser(bytes.NewBuffer(data))
	r.ContentLength = int64(len(data))

	r.Header.Set("Content-Type", "application/json")

	http.Redirect(w, r, callback.String(), http.StatusTemporaryRedirect)
}

func (mp *MockProvider) HasValidAuthentication(r *http.Request) bool {
	_, err := r.Cookie("access_token")
	return err == nil
}
