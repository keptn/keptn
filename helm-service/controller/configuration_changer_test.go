package controller

import (
	"testing"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"github.com/keptn/keptn/helm-service/pkg/helmtest"
	"github.com/stretchr/testify/assert"
)

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
    "replicas" : 2,
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

func TestChangePrimaryDeployment(t *testing.T) {

	data := helmtest.CreateHelmChartData(t)

	h := helm.NewGeneratedChartHandler(mesh.NewIstioMesh(), helm.NewCanaryOnDeploymentGenerator(), "mydomain.sh")
	inputChart, err := keptnutils.LoadChart(data)
	if err != nil {
		t.Error(err)
	}
	gen, err := h.GenerateManagedChart(inputChart, "sockshop", "production")
	assert.Nil(t, err, "Generating the managed Chart should not return any error")

	generatedChart, err := keptnutils.LoadChart(gen)
	if err != nil {
		t.Error(err)
	}

	change := keptnevents.PropertyChange{PropertyPath: "spec.replicas", Value: 2}
	event := keptnevents.ConfigurationChangeEventData{DeploymentChanges: []keptnevents.PropertyChange{change}}
	changePrimaryDeployment(&event, generatedChart)

	// Compare templates
	generatedTemplateResources := []helmtest.GeneratedResource{cartsCanaryIstioDestinationRuleGen, cartsIstioVirtualserviceGen, cartsPrimaryIstioDestinationRuleGen,
		deploymentGen, serviceGen}
	helmtest.Equals(generatedChart, valuesGen, generatedTemplateResources, t)
}
