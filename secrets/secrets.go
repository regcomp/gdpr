package secrets

import (
	"fmt"

	"github.com/regcomp/gdpr/config"
)

/*

Things that should live in the secrets store:
  - auth provider credentials
  - all database credentials
    - support for multiple databases
  - logger/audit credentials
    - where should logs/audits be sent

  - basically all credentials

*/

const (
	awsSecretsManager = "AWSSM"
	mockStoreType     = "MOCK"
)

type ISecretStore interface {
	getAllSecrets()

	GetServiceCacheSecrets() (*ServiceCacheSecrets, error)
	GetAuthProviderSecrets() (*AuthProviderSecrets, error)
	GetDatabaseManagerSecrets() (*DatabaseManagerSecrets, error)
}

func CreateSecretStore(config *config.SecretStoreConfig) (ISecretStore, error) {
	// this function dispatches to a function that creates a wrapper around
	// a connection to a secret store that satisfies the interface
	switch config.StoreType {
	case mockStoreType:
		return createLocalSecretStore(), nil
	default:
		return nil, fmt.Errorf("unknown store type=%s", config.StoreType)
	}
}
