package auth

import "net/http"

type MockProvider struct{}

func createMockProvider() *MockProvider {
	return &MockProvider{}
}

func (mp *MockProvider) AuthenticateUser(r *http.Request) (Credentials, error) {
	return Credentials{
		UserId:       "user-id",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}, nil
}
