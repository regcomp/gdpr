package database

type DatabaseProvider interface {
	GetScheduledDeletions() ([]ScheduledDeletion, error)
	BatchDeleteData(data []ScheduledDeletion) []error
}
