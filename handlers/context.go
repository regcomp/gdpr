package handlers

import (
	"log/slog"
	"os"

	"github.com/regcomp/gdpr/auth"
	"github.com/regcomp/gdpr/logging"
	"github.com/regcomp/gdpr/sessions"
)

var STX *ServiceContext

type ServiceContext struct {
	AuthProvider   auth.Provider
	SessionManager *sessions.SessionManager

	RequestLogger *slog.Logger

	HostPath             string
	AccessTokenDuration  int
	RefreshTokenDuration int
}

func CreateServiceContext(getenv func(string) string) *ServiceContext {
	// other context setup goes here, like getting certs/keys
	authProvider, err := auth.GetProvider(getenv)
	if err != nil {
		// TODO:
	}

	requestlogger := logging.NewRequestLogger(os.Stdout)

	return &ServiceContext{
		AuthProvider:  authProvider,
		RequestLogger: requestlogger,
		HostPath:      "localhost:8080",
	}
}

func LinkServiceContext(stx *ServiceContext) {
	STX = stx
}
