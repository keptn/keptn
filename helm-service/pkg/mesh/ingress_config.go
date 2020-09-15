package mesh

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
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
	if os.Getenv("ISTIO_GATEWAY") != "" {
		return os.Getenv("ISTIO_GATEWAY")
	}
	return "public-gateway.istio-system"
}

func GetLocalDeploymentURI(event keptnv2.EventData) []string {
	return []string{"http://" + event.Service + "." + event.Project + "-" + event.Stage}
}

func GetPublicDeploymentURI(event keptnv2.EventData) []string {
	return []string{GetIngressProtocol() + "://" + event.GetService() + "." + event.GetProject() + "-" + event.GetStage() + "." + GetIngressHostnameSuffix() + ":" + GetIngressPort()}
}
