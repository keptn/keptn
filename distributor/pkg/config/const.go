package config

type ConnectionType string

var ConnectionTypeToLocation = map[ConnectionType]string{ConnectionTypeNATS: "control-plane", ConnectionTypeHTTP: "remote-execution-plane"}

const mongoDbURL = "/mongodb-datastore"
const configServiceURL = "/configuration-service"
const controlPlaneURL = "/controlPlane"

var InClusterAPIProxyMappings = map[string]string{
	mongoDbURL:       "mongodb-datastore:8080",
	configServiceURL: "configuration-service:8080",
	controlPlaneURL:  "shipyard-controller:8080",
}

var ExternalAPIProxyMappings = map[string]string{
	mongoDbURL:       mongoDbURL,
	configServiceURL: configServiceURL,
	controlPlaneURL:  controlPlaneURL,
}

const (
	ConnectionTypeNATS ConnectionType = "nats"
	ConnectionTypeHTTP ConnectionType = "http"
)

const (
	DefaultShipyardControllerBaseURL = "http://shipyard-controller:8080"
	DefaultEventsEndpoint            = DefaultShipyardControllerBaseURL + "/v1/event/triggered"
	DefaultPollingInterval           = 10
	DefaultAPIProxyHTTPTimeout       = 30
)
