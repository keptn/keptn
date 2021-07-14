package config

type ConnectionType string

var ConnectionTypeToLocation = map[ConnectionType]string{ConnectionTypeNATS: "control-plane", ConnectionTypeHTTP: "remote-execution-plane"}

const mongoDbUrl =       "/mongodb-datastore"
const configServiceUrl = "/configuration-service"
const controlPlaneUrl =  "/controlPlane"

var InClusterAPIProxyMappings = map[string]string{
	mongoDbUrl:       "mongodb-datastore:8080",
	configServiceUrl: "configuration-service:8080",
	controlPlaneUrl:  "shipyard-controller:8080",
}

var ExternalAPIProxyMappings = map[string]string{
	mongoDbUrl:       mongoDbUrl,
	configServiceUrl: configServiceUrl,
	controlPlaneUrl:  controlPlaneUrl,
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
