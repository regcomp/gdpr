package config

import (
	"github.com/regcomp/gdpr/constants"
)

type IConfigStore interface {
	GetServiceURL() string
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

func NewLocalConfigStore(getters ...func(string) string) *LocalConfigStore {
	store := &LocalConfigStore{
		attrs: make(map[string]string),
	}
	store.initializeStore(getters...)
	return store
}

func (cs *LocalConfigStore) initializeStore(getters ...func(string) string) {
	for _, getter := range getters {
		for _, attr := range constants.ConfigAttrs {
			cs.attrs[attr] = getter(attr)
		}
	}
}

func (cs *LocalConfigStore) GetServiceURL() string {
	return cs.attrs[constants.ConfigServiceURLKey]
}

func (cs *LocalConfigStore) GetDefaultPort() string {
	return cs.attrs[constants.ConfigDefaultPortKey]
}

func (cs *LocalConfigStore) GetSessionDuration() string {
	return cs.attrs[constants.ConfigSessionDurationKey]
}

func (cs *LocalConfigStore) GetTracerConfig() *RequestTracerConfig {
	return &RequestTracerConfig{
		TracerOn: cs.attrs[constants.ConfigRequestTracerOnKey],
	}
}

func (cs *LocalConfigStore) GetSecretStoreConfig() *SecretStoreConfig {
	return &SecretStoreConfig{
		StoreType: cs.attrs[constants.ConfigSecretStoreTypeKey],
	}
}

func (cs *LocalConfigStore) GetServiceCacheConfig() *ServiceCacheConfig {
	return &ServiceCacheConfig{
		CacheType: cs.attrs[constants.ConfigServiceCacheTypeKey],
	}
}

func (cs *LocalConfigStore) GetAuthProviderConfig() *AuthProviderConfig {
	return &AuthProviderConfig{
		ProviderType: cs.attrs[constants.ConfigAuthProvierTypeKey],
	}
}

func (cs *LocalConfigStore) GetDatabaseStoreConfig() *DatabaseStoreConfig {
	return &DatabaseStoreConfig{
		//
	}
}
