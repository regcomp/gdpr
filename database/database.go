package database

type IDatabaseProvider interface {
	GetScheduledDeletions() ([]ScheduledDeletion, error)
	BatchDeleteData(data []ScheduledDeletion) []error
}

type IDatabaseStore interface {
	IDatabaseStore()
}

type DatabaseStore struct {
	databases []IDatabaseProvider
}

func CreateDatabaseStore(getenv func(string) string) (*DatabaseStore, error) {
	return &DatabaseStore{
		databases: make([]IDatabaseProvider, 8),
	}, nil
}

func (dbs *DatabaseStore) IDatabaseStore() {}
