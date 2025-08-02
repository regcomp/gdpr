package database

import "time"

type LocalDatabaseManager struct {
	records   IRecordsDatabase
	providers []IDatabaseProvider
}

func createLocalDatabaseManager() *LocalDatabaseManager {
	return &LocalDatabaseManager{
		// TODO: TODO TODO TODO TODO
		records:   nil,
		providers: nil,
	}
}

func (mp *LocalDatabaseManager) GetScheduledDeletionRecords(limit int, start time.Time) ([]RecordOfDeletionRequest, PaginationInfo, error) {
	//
	return nil, PaginationInfo{}, nil
}

func (mp *LocalDatabaseManager) BatchDeleteData(data []RecordOfDeletionRequest) []error {
	return nil
}
