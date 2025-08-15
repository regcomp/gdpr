package secrets

/*
Things that should live in the secrets store:
  - auth provider credentials
  - all database credentials
    - support for multiple databases
  - logger/audit credentials
    - where should logs/audits be sent

  - basically all credentials
*/

import (
	"fmt"

	"github.com/regcomp/gdpr/internal/config"
)

const (
	awsSecretsManager = "AWSSM"
	mockStoreType     = "MOCK"
)

type ISecretStore interface {
	GetAllSecrets()

	GetServiceCacheSecrets() (*ServiceCacheSecrets, error)
	GetAuthProviderSecrets() (*AuthProviderSecrets, error)
	GetDatabaseManagerSecrets() (*DatabaseManagerSecrets, error)
}

type SecretStoreSecrets struct {
	//
}

type ServiceCacheSecrets struct {
	//
}

type AuthProviderSecrets struct {
	//
}

type DatabaseManagerSecrets struct {
	RecordsSecrets  DatabaseSecrets
	ProviderSecrets map[string]DatabaseSecrets
}

// DatabaseSecrets put all possible fields in here
type DatabaseSecrets struct {
	URL string
	Key string
}

func CreateSecretStore(config *config.SecretStoreConfig) (ISecretStore, error) {
	// this function dispatches to a function that creates a wrapper around
	// a connection to a secret store that satisfies the interface
	switch config.StoreType {
	case mockStoreType:
		return CreateLocalSecretStore(), nil
	default:
		return nil, fmt.Errorf("unknown store type=%s", config.StoreType)
	}
}
