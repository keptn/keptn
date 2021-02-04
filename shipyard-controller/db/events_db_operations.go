package db

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/events_db_operations_moq.go . EventsDbOperations
type EventsDbOperations interface {
	UpdateEventOfService(event interface{}, eventType string, keptnContext string, eventID string, triggeredID string) error
}
