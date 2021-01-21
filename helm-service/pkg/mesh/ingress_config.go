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

// GetLocalDeploymentURI returns URIs where a service is accessible from within the cluster
func GetLocalDeploymentURI(event keptnv2.EventData, port string) []string {
	return []string{"http://" + event.Service + "." + event.Project + "-" + event.Stage + ":" + port}
}

// GetPublicDeploymentURI returns URIs where a service is exposed
func GetPublicDeploymentURI(event keptnv2.EventData) []string {
	publicHostName := GetPublicDeploymentHostNameTemplate()

	publicHostName = strings.ReplaceAll(publicHostName, "${INGRESS_PROTOCOL}", GetIngressProtocol())
	publicHostName = strings.ReplaceAll(publicHostName, "${SERVICE}", event.Service)
	publicHostName = strings.ReplaceAll(publicHostName, "${PROJECT}", event.Project)
	publicHostName = strings.ReplaceAll(publicHostName, "${STAGE}", event.Stage)
	publicHostName = strings.ReplaceAll(publicHostName, "${INGRESS_HOSTNAME_SUFFIX}", GetIngressHostnameSuffix())
	publicHostName = strings.ReplaceAll(publicHostName, "${INGRESS_PORT}", GetIngressPort())

	return []string{publicHostName}
}

// GetPublicDeploymentHostNameTemplate returns the HostName of the service
func GetPublicDeploymentHostNameTemplate() string {
	hostNameTemplate := os.Getenv("HOSTNAME_TEMPLATE")
	if hostNameTemplate == "" {
		return "${INGRESS_PROTOCOL}://${SERVICE}.${PROJECT}-${STAGE}.${INGRESS_HOSTNAME_SUFFIX}:${INGRESS_PORT}"
	}
	return strings.ToUpper(hostNameTemplate)
}
