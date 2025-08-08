package config

import (
	"fmt"
)

type IConfigStore interface {
	GetServiceURL() (string, error)
	GetServiceURLWithPort() (string, error)
	GetDefaultPort() (string, error)
	GetSessionDuration() (string, error)

	GetSecretStoreConfig() (*SecretStoreConfig, error)
	GetServiceCacheConfig() (*ServiceCacheConfig, error)
	GetAuthProviderConfig() (*AuthProviderConfig, error)
	GetDatabaseManagerConfig() (*DatabaseManagerConfig, error)
}

func NewConfigStore(getenv func(string) string, getters ...func(string) string) (IConfigStore, error) {
	storeType := getenv(ConfigConfigStoreTypeKey)
	switch storeType {
	case ValueLocalType:
		return newLocalConfigStore(getters...)
	default:
		return nil, fmt.Errorf("unknown config store type=%s", storeType)
	}
}
