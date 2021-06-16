package config

type ConnectionType string

var ConnectionTypeToLocation = map[ConnectionType]string{ConnectionTypeNATS: "control-plane", ConnectionTypeHTTP: "remote-execution-plane"}

var InClusterAPIProxyMappings = map[string]string{
	"/mongodb-datastore":     "mongodb-datastore:8080",
	"/configuration-service": "configuration-service:8080",
	"/controlPlane":          "shipyard-controller:8080",
}

var ExternalAPIProxyMappings = map[string]string{
	"/mongodb-datastore":     "/mongodb-datastore",
	"/configuration-service": "/configuration-service",
	"/controlPlane":          "/controlPlane",
}

const (
	ConnectionTypeNATS ConnectionType = "nats"
	ConnectionTypeHTTP ConnectionType = "http"
)

const (
	DefaultShipyardControllerBaseURL = "http://shipyard-controller:8080"
	DefaultEventsEndpoint            = DefaultShipyardControllerBaseURL + "/v1/event/triggered"
	DefaultPollingInterval           = 10
)
