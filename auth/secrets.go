package auth

type ISecretsStore interface {
	SecretsStore()
}

type MockSecretsStore struct {
	//
}

func (ss *MockSecretsStore) SecretsStore() {}
