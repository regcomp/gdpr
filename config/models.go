package config

type RequestTracerConfig struct {
	TracerOn         bool
	DisplayResponses bool
}

type SecretStoreConfig struct {
	StoreType string
}

// NOTE: This may need configurations for all of the kinds of data that
// will be stored
type ServiceCacheConfig struct {
	CacheType string
}

type AuthProviderConfig struct {
	ProviderType string
}

type DatabaseManagerConfig struct {
	RecordsConfig   RecordsDatabaseConfig
	ProviderConfigs []ProviderDatabaseConfig
}

type ProviderDatabaseConfig struct {
	ProviderName string
	ProviderType string
	TableNames   []string
}

type RecordsDatabaseConfig struct {
	ProviderType string
}
