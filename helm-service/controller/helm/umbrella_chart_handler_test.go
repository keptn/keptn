package helm

import (
	"testing"

	"github.com/keptn/keptn/helm-service/controller/mesh"
	"github.com/keptn/keptn/helm-service/pkg/objectutils"
	"github.com/kinbiko/jsonassert"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	"gotest.tools/assert"
)

func TestCreateRootChartResource(t *testing.T) {

	event := keptnevents.ServiceCreateEventData{Project: "sockshop", Service: "carts"}

	c := NewUmbrellaChartHandler(mesh.NewIstioMesh())
	resource, err := c.createRootChartResource(&event)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, *resource.ResourceURI, "Chart.yaml", "URI is wrong")

	jsonData, err := objectutils.ToJSON([]byte(resource.ResourceContent))
	if err != nil {
		t.Error(err)
	}
	ja := jsonassert.New(t)
	// find some sort of payload
	ja.Assertf(string(jsonData), `
    {
		"apiVersion": "v1",
		"description": "A Helm chart for project sockshop-umbrella",
		"name": "sockshop-umbrella",
		"version": "0.1.0"
	  }`)
}

func TestCreateRequirementsResource(t *testing.T) {

	c := NewUmbrellaChartHandler(mesh.NewIstioMesh())
	resource, err := c.createRequirementsResource()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, *resource.ResourceURI, "requirements.yaml", "URI is wrong")

	jsonData, err := objectutils.ToJSON([]byte(resource.ResourceContent))
	if err != nil {
		t.Error(err)
	}
	ja := jsonassert.New(t)
	// find some sort of payload
	ja.Assertf(string(jsonData), `
    {
		"dependencies": null
	}`)
}

func TestCreateValuesResource(t *testing.T) {

	c := NewUmbrellaChartHandler(mesh.NewIstioMesh())
	resource, err := c.createValuesResource()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, *resource.ResourceURI, "values.yaml", "URI is wrong")

	jsonData, err := objectutils.ToJSON([]byte(resource.ResourceContent))
	if err != nil {
		t.Error(err)
	}
	ja := jsonassert.New(t)
	// find some sort of payload
	ja.Assertf(string(jsonData), `
    {
	}`)
}

func TestCreateGatewayResource(t *testing.T) {

	c := NewUmbrellaChartHandler(mesh.NewIstioMesh())
	event := keptnevents.ServiceCreateEventData{Project: "sockshop", Service: "carts"}
	resource, err := c.createGatewayResource(&event, "dev")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, *resource.ResourceURI, "templates/istio-gateway.yaml", "URI is wrong")

	jsonData, err := objectutils.ToJSON([]byte(resource.ResourceContent))
	if err != nil {
		t.Error(err)
	}
	ja := jsonassert.New(t)
	// find some sort of payload
	ja.Assertf(string(jsonData), `
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
