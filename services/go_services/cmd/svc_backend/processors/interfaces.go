package processors

import "time"

// DataStore defines the interface for adding and retrieving key value pairs
type DataStore interface {
	AddFieldToHash(hashKey string, fieldName string, fieldValue time.Time) error
	GetFieldAsTime(hashKey string, fieldName string, layout string) (time.Time, error)
}

// SummaryPoster defines the interface for posting summaries.
type SummaryPoster interface {
	PostSummary(data interface{}) error
}
