package config

type RequestTracerConfig struct {
	TracerOn         bool
	DisplayResponses bool
}

type SecretStoreConfig struct {
	StoreType string
}

type ServiceCacheConfig struct {
	CacheType string
}

type AuthProviderConfig struct {
	ProviderType string
}

type DatabaseStoreConfig struct {
	//
}

type DatabaseProviderConfig struct {
	//
}
