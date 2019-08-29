package mesh

// Mesh abstracts the underlying mesh router
type Mesh interface {
	GenerateHTTPGateway(name string) ([]byte, error)
	GenerateDestinationRule(name string, host string) ([]byte, error)
	GenerateVirtualService(name string, gateways []string, hosts []string, httpRouteDestinations []HTTPRouteDestination) ([]byte, error)
	UpdateWeights(virtualService []byte, canaryWeight int32) ([]byte, error)
	GetDestinationRuleSuffix() string
	GetVirtualServiceSuffix() string
}

// HTTPRouteDestination helper struct for route destinations in a VirtualService
type HTTPRouteDestination struct {
	Host   string
	Weight int32
}
