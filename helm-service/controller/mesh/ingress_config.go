package mesh

import (
	"os"
	"strings"
)

// GetIngressHostnameSuffix returns the ingress hostname suffix
func GetIngressHostnameSuffix() string {
	if os.Getenv("INGRESS_HOSTNAME_SUFFIX") != "" {
		return os.Getenv("INGRESS_HOSTNAME_SUFFIX")
	}
	return "svc.cluster.local"
}

// GetIngressProtocol returns the ingress protocol
func GetIngressProtocol() string {
	if os.Getenv("INGRESS_PROTOCOL") != "" {
		return strings.ToLower(os.Getenv("INGRESS_PROTOCOL"))
	}
	return "http"
}

// GetIngressPort returns the ingress port
func GetIngressPort() string {
	if os.Getenv("INGRESS_PORT") != "" {
		return os.Getenv("INGRESS_PORT")
	}
	return "80"
}

// GetIngressGateway returns the ingress gateway
func GetIngressGateway() string {
	if os.Getenv("INGRESS_GATEWAY") != "" {
		return os.Getenv("INGRESS_GATEWAY")
	}
	return "public-gateway.istio-system"
}
