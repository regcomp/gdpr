package auth

type ISecretsStore interface {
	SecretsStore()
}

type SecretsStore struct {
	//
}

func (ss *SecretsStore) SecretsStore() {}
