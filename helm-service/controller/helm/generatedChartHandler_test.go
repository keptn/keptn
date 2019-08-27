package helm

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/keptn/keptn/helm-service/controller/jsonutils"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
)

const projectName = "sockshop"
const serviceName = "carts"

type GeneratedResource struct {
	URI         string
	FileContent []string
}

var cartsCanaryIstioDestinationRuleGen = GeneratedResource{
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

var cartsIstioVirtualserviceGen = GeneratedResource{
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
      "sockshop-production-gateway"
    ],
    "hosts": [
      "carts.sockshop-production.mydomain.sh"
    ],
    "http": [
      {
        "route": [
          {
            "destination": {
              "host": "carts-canary.sockshop-production.svc.cluster.local"
            },
            "weight": 80
          },
          {
            "destination": {
              "host": "carts-primary.sockshop-production.svc.cluster.local"
            },
            "weight": 20
          }
        ]
      }
    ]
  }
}
`},
}

var cartsPrimaryIstioDestinationRuleGen = GeneratedResource{
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

var deploymentGen = GeneratedResource{
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
            "env": [
              {
                "name": "DT_TAGS",
                "value": "application={{ .Chart.Name }}"
              },
              {
                "name": "POD_NAME",
                "valueFrom": {
                  "fieldRef": {
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "DEPLOYMENT_NAME",
                "valueFrom": {
                  "fieldRef": {
                    "fieldPath": "metadata.labels['deployment']"
                  }
                }
              },
              {
                "name": "CONTAINER_IMAGE",
                "value": "{{ .Values.image }}"
              },
              {
                "name": "KEPTN_PROJECT",
                "value": "{{ .Chart.Name }}"
              },
              {
                "name": "KEPTN_STAGE",
                "valueFrom": {
                  "fieldRef": {
                    "fieldPath": "metadata.namespace"
                  }
                }
              },
              {
                "name": "KEPTN_SERVICE",
                "value": "carts"
              }
            ],
            "image": "{{ .Values.image }}",
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

var serviceGen = GeneratedResource{
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

var chartGen = GeneratedResource{
	URI: "Chart.yaml",
	FileContent: []string{`
{
  "apiVersion": "v1",
  "description": "A Helm chart for service carts (generated)",
  "name": "carts-generated",
  "version": "0.1.0"
}`},
}

var valuesGen = GeneratedResource{
	URI: "values.yaml",
	FileContent: []string{`
{
  "image": "docker.io/keptnexamples/carts:0.8.1"
}`},
}

func TestGenerateManagedChart(t *testing.T) {

	data := CreateHelmChartData(t)
	event := keptnevents.ServiceCreateEventData{Project: projectName, Service: serviceName, HelmChart: data}

	istioMesh := mesh.NewIstioMesh()
	gen, err := GenerateManagedChart(&event, "production", istioMesh, "mydomain.sh")
	assert.Nil(t, err, "Generating the managed Chart should not return any error")

	workingPath, err := ioutil.TempDir("", "helm-test")
	// defer os.RemoveAll(workingPath)
	fmt.Println(workingPath)
	keptnutils.Untar(workingPath, bytes.NewReader(gen))

	generatedResources := []GeneratedResource{cartsCanaryIstioDestinationRuleGen, cartsIstioVirtualserviceGen, cartsPrimaryIstioDestinationRuleGen,
		deploymentGen, serviceGen, chartGen, valuesGen}

	for _, resource := range generatedResources {

		f, err := os.Open(filepath.Join(workingPath, resource.URI))
		assert.Nil(t, err, "Reading generated data should not return any error")

		decoder := kyaml.NewDocumentDecoder(f)

		for i := 0; ; i++ {
			b1 := make([]byte, 4096)
			n1, err := decoder.Read(b1)
			if err == io.EOF {
				break
			}
			assert.Nil(t, err, "")

			jsonData, err := jsonutils.ToJSON(b1[:n1])
			if err != nil {
				t.Error(err)
			}

			ja := jsonassert.New(t)
			ja.Assertf(string(jsonData), resource.FileContent[i])
		}
	}
}
