package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/regcomp/gdpr/internal/config"
	"github.com/regcomp/gdpr/internal/secrets"
	"github.com/regcomp/gdpr/pkg/helpers"
)

var (
	ErrEntryNotFound error = errors.New("entry not found")
	ErrNoMoreData    error = errors.New("no more data")
)

type IRecordsDatabase interface {
	GetRecordsOfDeletionRequest(limit int, start time.Time) ([]RecordOfDeletionRequest, error)
	AddToDataDeletionsQueue([]uuid.UUID) error
	GetDataDeletionsQueue() []uuid.UUID
}

type IDatabaseProvider interface {
	GetRegisteredTableNames() []string
	DeleteDataFromRegisteredTables([]string, []uuid.UUID) []error
}

type RecordOfDeletionRequest struct {
	ID            uuid.UUID `json:"id"`
	CustomerID    uuid.UUID `json:"customerId"`
	CustomerName  string    `json:"name"`
	CustomerEmail string    `json:"email"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	RequestedAt   time.Time `json:"requestedAt"`
}

type PaginationInfo struct {
	HasMore    bool
	NextCursor time.Time
	Total      int
}

type CustomerData struct {
	ID uuid.UUID
}

func CreateRecordsDatabase(cfg config.RecordsDatabaseConfig, secrets secrets.DatabaseSecrets) (IRecordsDatabase, error) {
	switch cfg.ProviderType {
	case config.ValueLocalType:
		return CreateLocalRecordsDatabase(cfg), nil
	default:
		return nil, fmt.Errorf("unknown database type=%s", cfg.ProviderType)
	}
}

func CreateDatabaseProvider(cfg config.ProviderDatabaseConfig, secrets secrets.DatabaseSecrets) (IDatabaseProvider, error) {
	switch cfg.ProviderType {
	case config.ValueLocalType:
		return CreateLocalDatabaseProvider(cfg), nil
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
