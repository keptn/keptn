package db

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/events_operations_mock.go . EventsDbOperations
type EventsDbOperations interface {
	UpdateEventOfService(event interface{}, eventType string, keptnContext string, eventID string, triggeredID string) error
	UpdateShipyard(projectName string, shipyardContent string) error
}
