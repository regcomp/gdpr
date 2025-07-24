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

type StoreType string

const (
	AWSSecretsManager StoreType = "AWSSM"
	MockStoreType     StoreType = "MOCK"
)

type ISecretStore interface {
	SecretStore()
}

func CreateSecretStore(config ISecretStoreConfig) (ISecretStore, error) {
	// this function dispatches to a function that creates a wrapper around
	// a connection to a secret store that satisfies the interface
	switch config.StoreType() {
	default:
		return CreateMockSecretStore(), nil
	}
}

type MockSecretStore struct {
	//
}

func CreateMockSecretStore() *MockSecretStore {
	return &MockSecretStore{}
}

func (mss *MockSecretStore) SecretStore() {}

type ISecretStoreConfig interface {
	StoreType() StoreType
}

func LoadConfig(secretStoreType string) ISecretStoreConfig {
	// This function should grab all relevant information for whatever provider
	// specified in the secret store type from the environment. Data is passed
	// through env variables at runtime with docker. There shouldn't be any sensitive
	// data needed.
	switch secretStoreType {
	default:
		return NewMockSecretStoreConfig()
	}
}

type MockSecretStoreConfig struct {
	//
}

func NewMockSecretStoreConfig() *MockSecretStoreConfig {
	return &MockSecretStoreConfig{}
}

func (mssc *MockSecretStoreConfig) StoreType() StoreType { return MockStoreType }
