package handlers

import (
	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/logging"
)

var STX *ServiceContext

type ServiceContext struct {
	Logger               logging.Logger
	AuthProvider         auth.Provider
	AccessTokenDuration  int
	RefreshTokenDuration int
}

func CreateServiceContext(getenv func(string) string) *ServiceContext {
	// other context setup goes here, like getting certs/keys
	logger := &logging.MockLogger{}
	authProvider, err := auth.GetProvider(getenv)
	if err != nil {
		// TODO: handle
	}

	return &ServiceContext{
		Logger:       logger,
		AuthProvider: authProvider,
	}
}

func LinkServiceContext(stx *ServiceContext) {
	STX = stx
}
