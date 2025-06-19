package database

import "time"

type ScheduledDeletion struct {
	Key         string
	RequestedOn time.Time
}
