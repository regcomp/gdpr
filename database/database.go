package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/helpers"
	"github.com/regcomp/gdpr/secrets"
)

var (
	errEntryNotFound error = errors.New("entry not found")
	errNoMoreData    error = errors.New("no more data")
)

type IRecordsDatabase interface {
	GetRecordsOfDeletionRequest(limit int, start time.Time) ([]RecordOfDeletionRequest, error)
	AddToDataDeletionsQueue([]uuid.UUID) error
	GetDataDeletionsQueue() []uuid.UUID
}

func CreateRecordsDatabase(cfg config.RecordsDatabaseConfig, secrets secrets.DatabaseSecrets) (IRecordsDatabase, error) {
	switch cfg.ProviderType {
	case config.ValueLocalType:
		return createLocalRecordsDatabase(cfg), nil
	default:
		return nil, fmt.Errorf("unknown database type=%s", cfg.ProviderType)
	}
}

type IDatabaseProvider interface {
	GetRegisteredTableNames() []string
	DeleteDataFromRegisteredTables([]string, []uuid.UUID) []error
}

func CreateDatabaseProvider(cfg config.ProviderDatabaseConfig, secrets secrets.DatabaseSecrets) (IDatabaseProvider, error) {
	switch cfg.ProviderType {
	case config.ValueLocalType:
		return createLocalDatabaseProvider(cfg), nil
	default:
		return nil, fmt.Errorf("unknown database type=%s", cfg.ProviderType)
	}
}

type DatabaseManager struct {
	records   IRecordsDatabase
	providers []IDatabaseProvider
}

func CreateDatabaseManager(cfg *config.DatabaseManagerConfig, secrets *secrets.DatabaseManagerSecrets) (*DatabaseManager, error) {
	records, err := CreateRecordsDatabase(cfg.RecordsConfig, secrets.RecordsSecrets)
	if err != nil {
		return nil, err
	}

	providers := make([]IDatabaseProvider, 0, 8)
	for _, providerConfig := range cfg.ProviderConfigs {
		providerSecrets, ok := secrets.ProviderSecrets[providerConfig.ProviderName]
		if !ok {
			return nil, fmt.Errorf("no secrets found for database provider=%s", providerConfig.ProviderName)
		}
		provider, err := CreateDatabaseProvider(providerConfig, providerSecrets)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}

	return &DatabaseManager{
		records:   records,
		providers: providers,
	}, nil
}

func (dbm *DatabaseManager) GetDeletionRequestsAndPaginationInfo(limit int, start time.Time) ([]RecordOfDeletionRequest, PaginationInfo, error) {
	deletionRequests, err := dbm.records.GetRecordsOfDeletionRequest(limit, start)
	if err != nil {
		return nil, PaginationInfo{}, err
	}

	var hasMore bool
	var nextCursor time.Time
	if len(deletionRequests) < limit+1 {
		// too few items returned. no more data to paginate
		hasMore = false
		nextCursor = time.Now()
	} else {
		hasMore = true
		nextCursor = deletionRequests[len(deletionRequests)-1].CreatedAt
	}

	paginationInfo := PaginationInfo{
		HasMore:    hasMore,
		NextCursor: nextCursor,
	}

	return deletionRequests[:len(deletionRequests)-helpers.Btoi(hasMore)],
		paginationInfo,
		nil
}

func (dbm *DatabaseManager) AddToDataDeletionQueue(ids []uuid.UUID) error {
	err := dbm.records.AddToDataDeletionsQueue(ids)
	if err != nil {
		return err
	}
	return nil
}

func (dbm *DatabaseManager) RunDataDeletionsQueue() []error {
	deletions := dbm.records.GetDataDeletionsQueue()
	if len(deletions) == 0 {
		return nil
	}

	errs := make([]error, 0, 16)
	for _, provider := range dbm.providers {
		names := provider.GetRegisteredTableNames()
		errs = append(errs, provider.DeleteDataFromRegisteredTables(names, deletions)...)
	}

	// TODO: update deletion queue

	return errs
}
