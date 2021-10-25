package db

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/events_operations_mock.go . EventsDbOperations
type EventsDbOperations interface {
	UpdateShipyard(projectName string, shipyardContent string) error
}
