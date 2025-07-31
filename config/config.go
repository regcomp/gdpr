package config

import (
	"log"
)

type IConfigStore interface {
	GetServiceURL() string
	GetServiceURLWithPort() string
	GetDefaultPort() string
	GetSessionDuration() string

	GetSecretStoreConfig() *SecretStoreConfig
	GetServiceCacheConfig() *ServiceCacheConfig
	GetAuthProviderConfig() *AuthProviderConfig
	GetDatabaseStoreConfig() *DatabaseStoreConfig
}

type LocalConfigStore struct {
	// mu    sync.RWMutex // May need in the future
	attrs map[string]string
}

func NewConfigStore(getenv func(string) string, getters ...func(string) string) IConfigStore {
	storeType := getenv(ConfigConfigStoreTypeKey)
	switch storeType {
	case "LOCAL":
		return newLocalConfigStore(getters...)
	default:
		log.Fatalf("unknown config store type: %s", storeType)
		return nil
	}
}

func newLocalConfigStore(getters ...func(string) string) *LocalConfigStore {
	store := &LocalConfigStore{
		attrs: make(map[string]string),
	}
	store.initializeLocalStore(getters...)
	return store
}

func (cs *LocalConfigStore) initializeLocalStore(getters ...func(string) string) {
	for _, getter := range getters {
		for _, attr := range ConfigAttrs {
			cs.attrs[attr] = getter(attr)
		}
	}
}

func (cs *LocalConfigStore) GetServiceURL() string {
	return cs.attrs[ConfigServiceUrlKey]
}

func (cs *LocalConfigStore) GetDefaultPort() string {
	return cs.attrs[ConfigDefaultPortKey]
}

func (cs *LocalConfigStore) GetServiceURLWithPort() string {
	return cs.attrs[ConfigServiceUrlKey] + ":" + cs.GetDefaultPort()
}

func (cs *LocalConfigStore) GetSessionDuration() string {
	return cs.attrs[ConfigSessionDurationKey]
}

func (cs *LocalConfigStore) GetSecretStoreConfig() *SecretStoreConfig {
	return &SecretStoreConfig{
		StoreType: cs.attrs[ConfigSecretStoreTypeKey],
	}
}

func (cs *LocalConfigStore) GetServiceCacheConfig() *ServiceCacheConfig {
	return &ServiceCacheConfig{
		CacheType: cs.attrs[ConfigServiceCacheTypeKey],
	}
}

func (cs *LocalConfigStore) GetAuthProviderConfig() *AuthProviderConfig {
	return &AuthProviderConfig{
		ProviderType: cs.attrs[ConfigAuthProviderTypeKey],
	}
}

func (cs *LocalConfigStore) GetDatabaseStoreConfig() *DatabaseStoreConfig {
	return &DatabaseStoreConfig{
		//
	}
}
