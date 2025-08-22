package database

import (
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/regcomp/gdpr/internal/config"
)

type LocalRecordsDatabase struct {
	records        []RecordOfDeletionRequest
	deletionsQueue []uuid.UUID
	// also audit trail of actions taken by users
}

func makeDummyData(n int) []RecordOfDeletionRequest {
	records := make([]RecordOfDeletionRequest, 0, 16)
	for i := range n {
		records = append(records, RecordOfDeletionRequest{
			ID:           uuid.New(),
			CustomerID:   uuid.New(),
			CustomerName: strconv.Itoa(i),
			CreatedAt:    time.Unix(int64(i), 0),
			UpdatedAt:    time.Unix(int64(i), 0),
			RequestedAt:  time.Unix(int64(i), 0),
		})
	}
	return records
}

func CreateLocalRecordsDatabase(cfg config.RecordsDatabaseConfig) *LocalRecordsDatabase {
	return &LocalRecordsDatabase{
		records:        makeDummyData(64),
		deletionsQueue: make([]uuid.UUID, 0, 16),
	}
}

// GetRecordsOfDeletionRequest will return an array of length 1 greater than the limit if there is enough data.
// The extra element is the next cursor start and signals that there is more data for pagination.
func (lrd *LocalRecordsDatabase) GetRecordsOfDeletionRequest(limit int, start time.Time) ([]RecordOfDeletionRequest, error) {
	slices.SortFunc(lrd.records, func(left, right RecordOfDeletionRequest) int {
		if left.CreatedAt.Compare(right.CreatedAt) < 0 {
			return -1
		} else if left.CreatedAt.Compare(right.CreatedAt) > 0 {
			return 1
		} else {
			return strings.Compare(left.CustomerName, right.CustomerName)
		}
	})
	startIdx := getIdxFromCreatedAt(lrd.records, start)
	if startIdx < 0 {
		return nil, ErrNoMoreData
	}

	endIdx := getEndIdxFromLimit(lrd.records, startIdx, limit)
	result := slices.Clone(lrd.records[startIdx:endIdx])

	log.Printf("-------\ngetting records, cursor=%s, startIdx=%d, endIdx=%d\n-------\n",
		start.String(), startIdx, endIdx,
	)
	return result, nil
}

func (lrd *LocalRecordsDatabase) AddToDataDeletionsQueue(ids []uuid.UUID) error {
	lrd.deletionsQueue = append(lrd.deletionsQueue, ids...)
	return nil
}

func (lrd *LocalRecordsDatabase) GetDataDeletionsQueue() []uuid.UUID {
	return lrd.deletionsQueue
}

// getIdxFromCreatedAt will return the next closest match in the case that the exact time
// isn't found
func getIdxFromCreatedAt(records []RecordOfDeletionRequest, start time.Time) int {
	for idx, record := range records {
		if record.CreatedAt.Compare(start) < 0 {
			continue
		}
		return idx
	}

	return -1
}

func getEndIdxFromLimit(records []RecordOfDeletionRequest, startIdx, limit int) int {
	availableRecords := len(records) - startIdx
	if availableRecords > limit+1 {
		return startIdx + limit + 1
	}
	return len(records)
}

type localTable struct {
	entries []map[string]string
}

type LocalDatabaseProvider struct {
	// All the tables
	tables map[string]localTable

	// This should be set from the configuration of the provider
	registeredTablesNames []string
}

func CreateLocalDatabaseProvider(providerConfig config.ProviderDatabaseConfig) *LocalDatabaseProvider {
	tables := make(map[string]localTable, 8)
	for _, name := range providerConfig.TableNames {
		table := localTable{entries: make([]map[string]string, 32)}
		tables[name] = table
	}

	return &LocalDatabaseProvider{
		tables:                tables,
		registeredTablesNames: providerConfig.TableNames,
	}
}

func (ldp *LocalDatabaseProvider) GetRegisteredTableNames() []string {
	return ldp.registeredTablesNames
}

func (ldp *LocalDatabaseProvider) DeleteDataFromRegisteredTables(names []string, ids []uuid.UUID) []error {
	return nil
}
