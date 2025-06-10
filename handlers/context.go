package handlers

import "github.com/regcomp/gdpr/auth"

var STX *ServiceContext

type ServiceContext struct {
	AuthProvider         auth.Provider
	AccessTokenDuration  int
	RefreshTokenDuration int
}

func CreateServiceContext(getenv func(string) string) *ServiceContext {
	// other context setup goes here, like getting certs/keys
	authProvider, err := auth.GetProvider(getenv)
	if err != nil {
		// TODO: handle
	}

	return &ServiceContext{
		AuthProvider: authProvider,
	}
}

func LinkServiceContext(stx *ServiceContext) {
	STX = stx
}
