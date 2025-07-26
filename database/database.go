package database

import (
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/secrets"
)

type IDatabaseProvider interface {
	GetScheduledDeletions() ([]ScheduledDeletion, error)
	BatchDeleteData(data []ScheduledDeletion) []error
}

type DatabaseStore struct {
	databases []IDatabaseProvider
}

func CreateDatabaseStore(config *config.DatabaseStoreConfig, secrets *secrets.DatabaseStoreSecrets) (*DatabaseStore, error) {
	return &DatabaseStore{
		databases: make([]IDatabaseProvider, 8),
	}, nil
}
