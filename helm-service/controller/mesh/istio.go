package mesh

import (
	"errors"
	"fmt"
	"strings"

	"github.com/keptn/keptn/helm-service/pkg/apis/networking/istio/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// IstioMesh is a Istio implementation of interface Mesh
type IstioMesh struct {
}

// NewIstioMesh generates a new istio mesh
func NewIstioMesh() *IstioMesh {
	return &IstioMesh{}
}

// GenerateDestinationRule generates a new Istio DestinationRule
func (*IstioMesh) GenerateDestinationRule(name string, host string) ([]byte, error) {

	dr := v1alpha3.DestinationRule{TypeMeta: metav1.TypeMeta{Kind: "DestinationRule", APIVersion: "networking.istio.io/v1alpha3"},
		ObjectMeta: metav1.ObjectMeta{Name: name}, Spec: v1alpha3.DestinationRuleSpec{Host: host}}
	return yaml.Marshal(dr)
}

// GenerateVirtualService generates a new Istio VirtualService
func (*IstioMesh) GenerateVirtualService(name string, gateways []string, hosts []string, httpRouteDestinations []HTTPRouteDestination) ([]byte, error) {

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
	return yaml.Marshal(vs)
}

// UpdateWeights returns a virtual service with updated weights
func (*IstioMesh) UpdateWeights(virtualService []byte, canaryWeight int32) ([]byte, error) {

	vs := v1alpha3.VirtualService{}
	err := yaml.Unmarshal(virtualService, &vs)
	if err != nil {
		return nil, err
	}

	primaryWeight := int32(100 - canaryWeight)
	if primaryWeight < 0 {
		return nil, errors.New("Invalid canary weight")
	}

	if len(vs.Spec.Http) > 0 {
		fmt.Println("Updating weight for HTTP Traffic")
		for _, httpRoute := range vs.Spec.Http {
			for _, dst := range httpRoute.Route {

				if !strings.HasPrefix(dst.Destination.Host, vs.ObjectMeta.Name) {
					return nil, fmt.Errorf("Cannot update VirutalService because host has unexpected name %s", dst.Destination.Host)
				}
				if strings.HasPrefix(dst.Destination.Host, vs.ObjectMeta.Name+"-canary") {
					dst.Weight = canaryWeight
				}
				if strings.HasPrefix(dst.Destination.Host, vs.ObjectMeta.Name+"-primary") {
					dst.Weight = primaryWeight
				}
			}
		}
	}

	// check if Tcp traffic is routed with this VirtualService
	if len(vs.Spec.Tcp) > 0 {
		fmt.Println("Updating weight for TCP Traffic")
		for _, tcpRoute := range vs.Spec.Tcp {
			for _, dst := range tcpRoute.Route {
				if !strings.HasPrefix(dst.Destination.Host, vs.ObjectMeta.Name) {
					return nil, fmt.Errorf("Cannot update VirtualService because host has unexpected name %s", dst.Destination.Host)
				}
				if strings.HasPrefix(dst.Destination.Host, vs.ObjectMeta.Name+"-canary") {
					dst.Weight = canaryWeight
				}
				if strings.HasPrefix(dst.Destination.Host, vs.ObjectMeta.Name+"-primary") {
					dst.Weight = primaryWeight
				}
			}
		}
	}

	return yaml.Marshal(vs)
}

// GetDestinationRuleSuffix returns the file name suffix of destination rules
func (*IstioMesh) GetDestinationRuleSuffix() string {
	return "-istio-destinationrule.yaml"
}

// GetVirtualServiceSuffix returns the file name suffix of virtual services
func (*IstioMesh) GetVirtualServiceSuffix() string {
	return "-istio-virtualservice.yaml"
}
