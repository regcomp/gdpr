package database

type IDatabaseProvider interface {
	GetScheduledDeletions() ([]ScheduledDeletion, error)
	BatchDeleteData(data []ScheduledDeletion) []error
}

type DatabaseStore struct {
	databases []IDatabaseProvider
}

func CreateDatabaseStore(getenv func(string) string) (*DatabaseStore, error) {
	return &DatabaseStore{
		databases: make([]IDatabaseProvider, 8),
	}, nil
}
