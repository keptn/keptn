package helm

import (
	"strings"
	"testing"

	"sigs.k8s.io/yaml"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"github.com/keptn/keptn/helm-service/pkg/apis/networking/istio/v1alpha3"
	"github.com/keptn/keptn/helm-service/pkg/helmtest"
	"github.com/stretchr/testify/assert"
)

const projectName = "sockshop"
const serviceName = "carts"

var cartsCanaryIstioDestinationRuleGen = helmtest.GeneratedResource{
	URI: "templates/carts-canary-istio-destinationrule.yaml",
	FileContent: []string{`
{
  "apiVersion": "networking.istio.io/v1alpha3",
  "kind": "DestinationRule",
  "metadata": {
    "creationTimestamp": null,
    "name": "carts-canary"
  },
  "spec": {
    "host": "carts-canary.sockshop-production.svc.cluster.local"
  }
}
`},
}

var cartsIstioVirtualserviceGen = helmtest.GeneratedResource{
	URI: "templates/carts-istio-virtualservice.yaml",
	FileContent: []string{`
{
  "apiVersion": "networking.istio.io/v1alpha3",
  "kind": "VirtualService",
  "metadata": {
    "creationTimestamp": null,
    "name": "carts"
  },
  "spec": {
    "gateways": [
      "sockshop-production-gateway.sockshop-production",
      "mesh"
    ],
    "hosts": [
      "carts.sockshop-production.mydomain.sh",
      "carts",
      "carts.sockshop-production"
    ],
    "http": [
      {
        "route": [
          {
            "destination": {
              "host": "carts-canary.sockshop-production.svc.cluster.local"
            }
          },
          {
            "destination": {
              "host": "carts-primary.sockshop-production.svc.cluster.local"
            },
            "weight": 100
          }
        ]
      }
    ]
  }
}
`},
}

var cartsPrimaryIstioDestinationRuleGen = helmtest.GeneratedResource{
	URI: "templates/carts-primary-istio-destinationrule.yaml",
	FileContent: []string{`
{
  "apiVersion": "networking.istio.io/v1alpha3",
  "kind": "DestinationRule",
  "metadata": {
    "creationTimestamp": null,
    "name": "carts-primary"
  },
  "spec": {
    "host": "carts-primary.sockshop-production.svc.cluster.local"
  }
}
`},
}

var deploymentGen = helmtest.GeneratedResource{
	URI: "templates/deployment.yml",
	FileContent: []string{`
{
  "apiVersion": "apps/v1",
  "kind": "Deployment",
  "metadata": {
    "creationTimestamp": null,
    "name": "carts-primary"
  },
  "spec": {
    "replicas" : 1,
    "selector": {
      "matchLabels": {
        "app": "carts-primary"
      }
    },
    "strategy": {
      "rollingUpdate": {
        "maxUnavailable": 0
      },
      "type": "RollingUpdate"
    },
    "template": {
      "metadata": {
        "creationTimestamp": null,
        "labels": {
          "app": "carts-primary"
        }
      },
      "spec": {
        "containers": [
          {            
            "image": "docker.io/keptnexamples/carts:0.8.1",
            "imagePullPolicy": "IfNotPresent",
            "livenessProbe": {
              "httpGet": {
                "path": "/health",
                "port": 8080
              },
              "initialDelaySeconds": 60,
              "periodSeconds": 10,
              "timeoutSeconds": 15
            },
            "name": "carts",
            "ports": [
              {
                "containerPort": 8080,
                "name": "http",
                "protocol": "TCP"
              }
            ],
            "readinessProbe": {
              "httpGet": {
                "path": "/health",
                "port": 8080
              },
              "initialDelaySeconds": 60,
              "periodSeconds": 10,
              "timeoutSeconds": 15
            },
            "resources": {
              "limits": {
                "cpu": "500m",
                "memory": "2Gi"
              },
              "requests": {
                "cpu": "250m",
                "memory": "1Gi"
              }
            }
          }
        ]
      }
    }
  },
  "status": {
  }
}`},
}

var serviceGen = helmtest.GeneratedResource{
	URI: "templates/service.yaml",
	FileContent: []string{`
{
  "apiVersion": "v1",
  "kind": "Service",
  "metadata": {
    "creationTimestamp": null,
    "name": "carts-canary"
  },
  "spec": {
    "ports": [
      {
        "name": "http",
        "port": 80,
        "protocol": "TCP",
        "targetPort": 8080
      }
    ],
    "selector": {
      "app": "carts"
    },
    "type": "LoadBalancer"
  },
  "status": {
    "loadBalancer": {
    }
  }
}`, `
{
  "apiVersion": "v1",
  "kind": "Service",
  "metadata": {
    "creationTimestamp": null,
    "name": "carts-primary"
  },
  "spec": {
    "ports": [
      {
        "name": "http",
        "port": 80,
        "protocol": "TCP",
        "targetPort": 8080
      }
    ],
    "selector": {
      "app": "carts-primary"
    },
    "type": "LoadBalancer"
  },
  "status": {
    "loadBalancer": {
    }
  }
}`}}

var valuesGen = helmtest.GeneratedResource{
	URI: "values.yaml",
	FileContent: []string{`
{
  "image": "docker.io/keptnexamples/carts:0.8.1",
  "replicas": 1
}`},
}

func TestGenerateManagedChart(t *testing.T) {

	data := helmtest.CreateHelmChartData(t)

	h := NewGeneratedChartHandler(mesh.NewIstioMesh(), NewCanaryOnDeploymentGenerator(), "mydomain.sh")
	inputChart, err := keptnutils.LoadChart(data)
	if err != nil {
		t.Error(err)
	}
	gen, err := h.GenerateManagedChart(inputChart, projectName, "production")
	assert.Nil(t, err, "Generating the managed Chart should not return any error")

	ch, err := keptnutils.LoadChart(gen)

	// Compare templates
	generatedTemplateResources := []helmtest.GeneratedResource{cartsCanaryIstioDestinationRuleGen, cartsIstioVirtualserviceGen, cartsPrimaryIstioDestinationRuleGen,
		deploymentGen, serviceGen}
	helmtest.Equals(ch, valuesGen, generatedTemplateResources, t)
}

func TestUpdateCanaryWeight(t *testing.T) {
	data := helmtest.CreateHelmChartData(t)

	h := NewGeneratedChartHandler(mesh.NewIstioMesh(), NewCanaryOnDeploymentGenerator(), "mydomain.sh")
	chart, err := keptnutils.LoadChart(data)
	if err != nil {
		t.Error(err)
	}
	genChartData, err := h.GenerateManagedChart(chart, "sockshop", "production")
	if err != nil {
		t.Error(err)
	}
	genChart, err := keptnutils.LoadChart(genChartData)
	if err != nil {
		t.Error(err)
	}

	h.UpdateCanaryWeight(genChart, 48)
	assert.Nil(t, err, "Generating the managed Chart should not return any error")

	template := helmtest.GetTemplateByName(genChart, "templates/carts-istio-virtualservice.yaml")
	assert.NotNil(t, template, "Template must not be null")
	vs := v1alpha3.VirtualService{}
	err = yaml.Unmarshal(template.Data, &vs)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, len(vs.Spec.Http))
	assert.Equal(t, 2, len(vs.Spec.Http[0].Route))
	assert.Equal(t, int32(48), vs.Spec.Http[0].Route[0].Weight)
	assert.True(t, strings.Contains(vs.Spec.Http[0].Route[0].Destination.Host, "canary"))
	assert.Equal(t, int32(52), vs.Spec.Http[0].Route[1].Weight)
	assert.True(t, strings.Contains(vs.Spec.Http[0].Route[1].Destination.Host, "primary"))
}
