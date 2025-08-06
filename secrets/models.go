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

type DatabaseSecrets struct {
	// All possible information that any database host could need
}
