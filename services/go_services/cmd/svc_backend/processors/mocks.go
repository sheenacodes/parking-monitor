package processors

import (
	"time"
)

// MockDataStore is a mock implementation of the DataStore interface.
type MockDataStore struct {
	AddFieldToHashFunc func(hashKey string, fieldName string, fieldValue time.Time) error
	GetFieldAsTimeFunc func(hashKey, fieldName, layout string) (time.Time, error)
}

func (m *MockDataStore) AddFieldToHash(hashKey string, fieldName string, fieldValue time.Time) error {
	if m.AddFieldToHashFunc != nil {
		return m.AddFieldToHashFunc(hashKey, fieldName, fieldValue)
	}
	return nil
}

func (m *MockDataStore) GetFieldAsTime(hashKey string, fieldName string, layout string) (time.Time, error) {
	if m.GetFieldAsTimeFunc != nil {
		return m.GetFieldAsTimeFunc(hashKey, fieldName, layout)
	}
	return time.Time{}, nil
}

// MockSummaryPoster is a mock implementation of the SummaryPoster interface for testing.
type MockSummaryPoster struct {
	PostSummaryFunc func(data interface{}) error
}

func (m *MockSummaryPoster) PostSummary(data interface{}) error {
	if m.PostSummaryFunc != nil {
		return m.PostSummaryFunc(data)
	}
	return nil
}
