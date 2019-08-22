package mesh

import (
	"encoding/json"

	"github.com/keptn/keptn/helm-service/pkg/apis/networking/istio/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateHTTPGateway generates a new Istio Gateway
func GenerateHTTPGateway(name string) ([]byte, error) {

	selector := map[string]string{
		"istio": "ingressgateway",
	}

	port := v1alpha3.Port{Number: 80, Protocol: "HTTP", Name: "http"}
	server := v1alpha3.Server{Hosts: []string{"*"}, Port: &port}

	spec := v1alpha3.GatewaySpec{Selector: selector, Servers: []*v1alpha3.Server{&server}}

	gw := v1alpha3.Gateway{TypeMeta: metav1.TypeMeta{Kind: "Gateway", APIVersion: "networking.istio.io/v1alpha3"},
		ObjectMeta: metav1.ObjectMeta{Name: name}, Spec: spec}

	return json.Marshal(gw)
}

// GenerateDestinationRule generates a new Istio DestinationRule
func GenerateDestinationRule(name string, host string) ([]byte, error) {

	dr := v1alpha3.DestinationRule{TypeMeta: metav1.TypeMeta{Kind: "DestinationRule", APIVersion: "networking.istio.io/v1alpha3"},
		ObjectMeta: metav1.ObjectMeta{Name: name}, Spec: v1alpha3.DestinationRuleSpec{Host: host}}
	return json.Marshal(dr)
}

// GenerateVirtualService generates a new Istio VirtualService
func GenerateVirtualService(name string, gateways []string, hosts []string, httpRouteDestinations []HTTPRouteDestination) ([]byte, error) {

	destinations := []*v1alpha3.HTTPRouteDestination{}
	for _, httpRouteDst := range httpRouteDestinations {
		dst := v1alpha3.Destination{Host: httpRouteDst.Host}
		routeDst := v1alpha3.HTTPRouteDestination{Weight: httpRouteDst.Weight, Destination: &dst}
		destinations = append(destinations, &routeDst)
	}

	httpRoute := v1alpha3.HTTPRoute{Route: destinations}
	spec := v1alpha3.VirtualServiceSpec{Gateways: gateways, Hosts: hosts, Http: []*v1alpha3.HTTPRoute{&httpRoute}}

	vs := v1alpha3.VirtualService{TypeMeta: metav1.TypeMeta{Kind: "VirtualService", APIVersion: "networking.istio.io/v1alpha3"},
		ObjectMeta: metav1.ObjectMeta{Name: name}, Spec: spec}
	return json.Marshal(vs)
}
