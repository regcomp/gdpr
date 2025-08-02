package database

import (
	"time"

	"github.com/google/uuid"
)

type RecordOfDeletionRequest struct {
	ID          uuid.UUID
	CustomerID  uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	RequestedOn time.Time
}

type PaginationInfo struct {
	HasMore    bool
	NextCursor time.Time
	Total      int
}

type CustomerData struct {
	ID uuid.UUID
}
