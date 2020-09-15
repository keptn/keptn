package mesh

import (
	"testing"

	"github.com/keptn/keptn/helm-service/pkg/objectutils"
	"github.com/kinbiko/jsonassert"
)

func TestDestinationRule(t *testing.T) {

	istioMesh := NewIstioMesh()
	data, err := istioMesh.GenerateDestinationRule("carts-primary", "carts-primary.sockshop-dev.svc.cluster.local")
	if err != nil {
		t.Error(err)
	}
	jsonData, err := objectutils.ToJSON(data)
	if err != nil {
		t.Error(err)
	}

	ja := jsonassert.New(t)
	// find some sort of payload
	ja.Assertf(string(jsonData), `
    {
		"apiVersion": "networking.istio.io/v1alpha3",
		"kind": "DestinationRule",
		"metadata": {
		  "name": "carts-primary",
		  "creationTimestamp": null
		},
		"spec": {
		  "host": "carts-primary.sockshop-dev.svc.cluster.local"
		}
	  }`)
}

func TestVirtualService(t *testing.T) {

	routeDestinations := []HTTPRouteDestination{{Host: "carts-primary.sockshop-dev.svc.cluster.local", Weight: 50},
		{Host: "carts-canary.sockshop-dev.svc.cluster.local", Weight: 50}}

	istioMesh := NewIstioMesh()
	data, err := istioMesh.GenerateVirtualService("carts", []string{"public-gateway.istio-system"}, []string{"carts.sockshop-dev.35.226.86.78.xip.io"}, routeDestinations)
	if err != nil {
		t.Error(err)
	}
	jsonData, err := objectutils.ToJSON(data)
	if err != nil {
		t.Error(err)
	}

	ja := jsonassert.New(t)
	// find some sort of payload
	ja.Assertf(string(jsonData), `
    {
		"apiVersion": "networking.istio.io/v1alpha3",
		"kind": "VirtualService",
		"metadata": {
		  "name": "carts",
		  "creationTimestamp": null
		},
		"spec": {
		  "gateways": [
			"public-gateway.istio-system"
		  ],
		  "hosts": [
			"carts.sockshop-dev.35.226.86.78.xip.io"
		  ],
		  "http": [
			{
			  "route": [
				{
				  "destination": {
					"host": "carts-primary.sockshop-dev.svc.cluster.local"
				  },
				  "weight": 50
				},
				{
				  "destination": {
					"host": "carts-canary.sockshop-dev.svc.cluster.local"
				  },
				  "weight": 50
				}
			  ]
			}
		  ]
		}
	  }`)
}
