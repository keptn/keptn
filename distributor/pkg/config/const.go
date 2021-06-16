package config

type ConnectionType string

var ConnectionTypeToLocation = map[ConnectionType]string{ConnectionTypeNATS: "control-plane", ConnectionTypeHTTP: "remote-execution-plane"}

const (
	ConnectionTypeNATS ConnectionType = "nats"
	ConnectionTypeHTTP ConnectionType = "http"
)

const (
	DefaultShipyardControllerBaseURL = "http://shipyard-controller:8080"
	DefaultEventsEndpoint            = DefaultShipyardControllerBaseURL + "/v1/event/triggered"
	DefaultPollingInterval           = 10
)
