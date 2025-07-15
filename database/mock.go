package database

type MockProvider struct{}

func (mp *MockProvider) GetScheduledDeletions() ([]ScheduledDeletion, error) {
	//
	return nil, nil
}

func (mp *MockProvider) BatchDeleteData(data []ScheduledDeletion) []error {
	return nil
}
