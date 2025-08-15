package secrets

import (
	"sync"
)

type LocalSecretStore struct {
	mu  sync.RWMutex
	dat map[string][]byte
}

func CreateLocalSecretStore() *LocalSecretStore {
	return &LocalSecretStore{
		mu:  sync.RWMutex{},
		dat: make(map[string][]byte, 8),
	}
}

func (ss *LocalSecretStore) GetAllSecrets() {}
func (ss *LocalSecretStore) GetServiceCacheSecrets() (*ServiceCacheSecrets, error) {
	return &ServiceCacheSecrets{}, nil
}

func (ss *LocalSecretStore) GetAuthProviderSecrets() (*AuthProviderSecrets, error) {
	return &AuthProviderSecrets{}, nil
}

func (ss *LocalSecretStore) GetDatabaseManagerSecrets() (*DatabaseManagerSecrets, error) {
	providerSecrets := make(map[string]DatabaseSecrets)
	// WARN: This is coupled to the config. hard coded for ease of debugging
	databases := []string{"one", "two", "three"}
	for _, database := range databases {
		providerSecrets[database] = DatabaseSecrets{}
	}
	return &DatabaseManagerSecrets{
		ProviderSecrets: providerSecrets,
	}, nil
}
