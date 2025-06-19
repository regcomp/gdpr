package auth

import (
	"net/http"
	"net/url"
)

const (
	accessToken  = "foo"
	refreshToken = "bar"
)

type MockProvider struct{}

func createMockAuthProvider() *MockProvider {
	return &MockProvider{}
}

func (mp *MockProvider) GetProviderType() ProviderType {
	return MOCK
}

func (mp *MockProvider) AuthenticateUser(w http.ResponseWriter, r *http.Request, callback url.URL) {
	redirectURL, err := url.Parse(callback.String())
	if err != nil {
		// TODO:
	}

	params := url.Values{}
	params.Add("access", accessToken)
	params.Add("refresh", refreshToken)

	redirectURL.RawQuery = params.Encode()
	http.Redirect(w, r, redirectURL.String(), http.StatusTemporaryRedirect)
}

func (mp *MockProvider) IsValidAccessToken(token string) bool {
	return true
}

func (mp *MockProvider) GetNewAccessToken(refreshToken string) (string, error) {
	return accessToken, nil
}
