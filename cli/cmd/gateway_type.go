package cmd

// gateway is an enum for the type of the gateway
type gateway int

const (
	// NodePort uses a port of the node for the Gateway
	NodePort gateway = iota
	// LoadBalancer uses a LoadBalancer for the Gateway
	LoadBalancer
)

func (i gateway) String() string {
	return [...]string{"NodePort", "LoadBalancer"}[i]
}

var gatewayToID = map[string]gateway{
	"NodePort":     NodePort,
	"LoadBalancer": LoadBalancer,
}
