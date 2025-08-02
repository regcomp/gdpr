package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/secrets"
)

type IRecordsDatabase interface {
	GetScheduledDeletions(int, time.Time) ([]RecordOfDeletionRequest, error)
}

type IDatabaseProvider interface {
	BatchDeleteData(data []uuid.UUID) []error
}

type IDatabaseManager interface {
	GetScheduledDeletionRecords(int, time.Time) ([]RecordOfDeletionRequest, PaginationInfo, error)
}

func CreateDatabaseStore(config *config.DatabaseStoreConfig, secrets *secrets.DatabaseStoreSecrets) (IDatabaseManager, error) {
	// TODO: switch on the config and initialize db connections
	return createLocalDatabaseManager(), nil
}
