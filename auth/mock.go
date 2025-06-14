package auth

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

const (
	accessToken  = "access-token"
	refreshToken = "refresh-token"
	userID       = "user-id"
)

type MockProvider struct {
	Client *http.Client
}

func createMockAuthProvider() *MockProvider {
	return &MockProvider{
		Client: &http.Client{},
	}
}

func (mp *MockProvider) GetProviderType() ProviderType {
	return MOCK
}

func (mp *MockProvider) AuthenticateUser(w http.ResponseWriter, r *http.Request, callback url.URL) {
	// TODO: respond
	log.Println("AuthenticateUser hit...")
	// -----

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
	log.Printf("callback url: %s", callback.String())

	body := bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, url.QueryEscape(callback.String()), body)
	if err != nil {
		// TODO:
	}
	_, err = mp.Client.Do(req)
	if err != nil {
		log.Panic("DefaultClient request failed\n")
	}
}

func (mp *MockProvider) HasValidAuthentication(r *http.Request) bool {
	_, err := r.Cookie("access_token")
	return err == nil
}
