package cmd

// apiServiceType is an enum for the service type of the api
type apiServiceType int

const (
	// ClusterIP exposes the api on a cluster-internal IP
	ClusterIP apiServiceType = iota
	// LoadBalancer exposes the service externally using a cloud provider's load balancer
	LoadBalancer
	// NodePort exposes the service on each node's IP at a static port
	NodePort
)

func (i apiServiceType) String() string {
	return [...]string{"ClusterIP", "LoadBalancer", "NodePort"}[i]
}

var apiServiceTypeToID = map[string]apiServiceType{
	"ClusterIP":    ClusterIP,
	"LoadBalancer": LoadBalancer,
	"NodePort":     NodePort,
}
