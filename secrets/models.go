package secrets

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
