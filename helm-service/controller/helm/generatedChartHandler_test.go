package helm

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"k8s.io/helm/pkg/proto/hapi/chart"

	"github.com/keptn/keptn/helm-service/controller/jsonutils"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/helm/pkg/chartutil"
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
      "sockshop-production-gateway",
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

var valuesGen = GeneratedResource{
	URI: "values.yaml",
	FileContent: []string{`
{
  "image": "docker.io/keptnexamples/carts:0.8.1"
}`},
}

func TestGenerateManagedChart(t *testing.T) {

	data := CreateHelmChartData(t)

	h := NewGeneratedChartHandler(mesh.NewIstioMesh(), NewCanaryOnDeploymentGenerator(), "mydomain.sh")
	inputChart, err := LoadChart(data)
	if err != nil {
		t.Error(err)
	}
	gen, err := h.GenerateManagedChart(inputChart, projectName, "production")
	assert.Nil(t, err, "Generating the managed Chart should not return any error")

	workingPath, err := ioutil.TempDir("", "helm-test")
	defer os.RemoveAll(workingPath)
	packagedChartFilePath := filepath.Join(workingPath, serviceName)
	err = ioutil.WriteFile(packagedChartFilePath, gen, 0644)
	if err != nil {
		t.Error(err)
	}

	ch, err := chartutil.Load(packagedChartFilePath)

	// Compare values
	yReader := kyaml.NewYAMLReader(bufio.NewReader(bytes.NewReader([]byte(ch.Values.Raw))))
	yamlData, err := yReader.Read()
	if err != nil {
		t.Error(err)
	}
	jsonData, err := jsonutils.ToJSON(yamlData)
	if err != nil {
		t.Error(err)
	}

	ja := jsonassert.New(t)
	ja.Assertf(string(jsonData), valuesGen.FileContent[0])

	// Compare templates
	generatedTemplateResources := []GeneratedResource{cartsCanaryIstioDestinationRuleGen, cartsIstioVirtualserviceGen, cartsPrimaryIstioDestinationRuleGen,
		deploymentGen, serviceGen}

	for _, resource := range generatedTemplateResources {

		reader := ioutil.NopCloser(bytes.NewReader(getTemplateByName(ch, resource.URI).Data))
		decoder := kyaml.NewDocumentDecoder(reader)

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

func getTemplateByName(chart *chart.Chart, templateName string) *chart.Template {

	for _, template := range chart.Templates {
		if template.Name == templateName {
			return template
		}
	}
	return nil
}
