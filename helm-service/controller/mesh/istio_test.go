package mesh

import (
	"testing"

	"github.com/kinbiko/jsonassert"
)

func TestGenerateHTTPGateway(t *testing.T) {

	data, err := GenerateHTTPGateway("sockshop-dev-gateway")
	if err != nil {
		t.Error(err)
	}

	ja := jsonassert.New(t)
	// find some sort of payload
	ja.Assertf(string(data), `
    {
		"apiVersion": "networking.istio.io/v1alpha3",
		"kind": "Gateway",
		"metadata": {
		  "name": "sockshop-dev-gateway",
		  "creationTimestamp": null
		},
		"spec": {
		  "selector": {
			"istio": "ingressgateway"
		  },
		  "servers": [
			{
			  "port": {
				"number": 80,
				"name": "http",
				"protocol": "HTTP"
			  },
			  "hosts": [
				"*"
			  ]
			}
		  ]
		}
	  }`)
}

func TestDestinationRule(t *testing.T) {

	data, err := GenerateDestinationRule("carts-primary", "carts-primary.sockshop-dev.svc.cluster.local")
	if err != nil {
		t.Error(err)
	}

	ja := jsonassert.New(t)
	// find some sort of payload
	ja.Assertf(string(data), `
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

	routeDestinations := []HTTPRouteDestination{HTTPRouteDestination{Host: "carts-primary.sockshop-dev.svc.cluster.local", Weight: 50},
		HTTPRouteDestination{Host: "carts-canary.sockshop-dev.svc.cluster.local", Weight: 50}}

	data, err := GenerateVirtualService("carts", []string{"sockshop-dev-gateway"},
		[]string{"carts.sockshop-dev.35.226.86.78.xip.io"}, routeDestinations)
	if err != nil {
		t.Error(err)
	}

	ja := jsonassert.New(t)
	// find some sort of payload
	ja.Assertf(string(data), `
    {
		"apiVersion": "networking.istio.io/v1alpha3",
		"kind": "VirtualService",
		"metadata": {
		  "name": "carts",
		  "creationTimestamp": null
		},
		"spec": {
		  "gateways": [
			"sockshop-dev-gateway"
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
