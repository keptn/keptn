package helmtest

import (
	"bytes"
	"encoding/json"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/keptn/keptn/helm-service/pkg/objectutils"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/chart/loader"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
)

const serviceContent = `--- 
apiVersion: v1
kind: Service
metadata: 
  name: carts
spec: 
  type: LoadBalancer
  ports: 
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8080
  selector: 
    app: carts
`

const deploymentContent = `--- 
apiVersion: apps/v1
kind: Deployment
metadata:
  name: carts
spec:
  replicas: {{ .Values.replicas }}
  strategy:
    rollingUpdate:
      maxUnavailable: 0
    type: RollingUpdate
  selector:
    matchLabels:
      app: carts
  template:
    metadata:
      labels:
        app: carts
    spec:
      containers:
      - name: carts
        image: "{{ .Values.image }}"
        imagePullPolicy: IfNotPresent
        ports:
        - name: http
          protocol: TCP
          containerPort: 8080
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 60
          periodSeconds: 10
          timeoutSeconds: 15
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 60
          periodSeconds: 10
          timeoutSeconds: 15
        resources:
          limits:
              cpu: 500m
              memory: 2048Mi
          requests:
              cpu: 250m
              memory: 1024Mi
`

const secretContent = `
apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: YWRtaW4=
  password: MWYyZDFlMmU2N2Rm
`

const valuesContent = `
image: docker.io/keptnexamples/carts:0.8.1
replicas: 1
`

const chartContent = `
apiVersion: v1
description: A Helm chart for service carts
name: carts
version: 0.1.0
`

func check(e error, t *testing.T) {
	if e != nil {
		t.Error(e)
	}
}

// CreateHelmChartData creates a new Helm chart tgz and returns its data
func CreateHelmChartData(t *testing.T) []byte {

	err := os.MkdirAll("carts/templates", 0777)
	check(err, t)
	err = ioutil.WriteFile("carts/Chart.yaml", []byte(chartContent), 0644)
	check(err, t)
	err = ioutil.WriteFile("carts/values.yaml", []byte(valuesContent), 0644)
	check(err, t)
	err = ioutil.WriteFile("carts/templates/deployment.yml", []byte(deploymentContent), 0644)
	check(err, t)
	err = ioutil.WriteFile("carts/templates/service.yaml", []byte(serviceContent), 0644)
	check(err, t)
	err = ioutil.WriteFile("carts/templates/secret.yaml", []byte(secretContent), 0644)
	check(err, t)

	ch, err := loader.Load("carts")
	if err != nil {
		check(err, t)
	}

	name, err := chartutil.Save(ch, ".")
	if err != nil {
		check(err, t)
	}
	defer os.RemoveAll(name)
	defer os.RemoveAll("carts")

	bytes, err := ioutil.ReadFile(name)
	check(err, t)
	return bytes
}

type GeneratedResource struct {
	URI         string
	FileContent []string
}

func Equals(actual *chart.Chart, valuesExpected GeneratedResource, templatesExpected []GeneratedResource, t *testing.T) {

	// Compare values
	jsonData, err := json.Marshal(actual.Values)
	if err != nil {
		t.Error(err)
	}

	ja := jsonassert.New(t)
	ja.Assertf(string(jsonData), valuesExpected.FileContent[0])

	for _, resource := range templatesExpected {

		reader := ioutil.NopCloser(bytes.NewReader(GetTemplateByName(actual, resource.URI).Data))
		decoder := kyaml.NewDocumentDecoder(reader)

		for i := 0; ; i++ {
			b1 := make([]byte, 4096)
			n1, err := decoder.Read(b1)
			if err == io.EOF {
				break
			}
			assert.Nil(t, err, "")

			jsonData, err := objectutils.ToJSON(b1[:n1])
			if err != nil {
				t.Error(err)
			}

			ja := jsonassert.New(t)
			ja.Assertf(string(jsonData), resource.FileContent[i])
		}
	}
}

func GetTemplateByName(chart *chart.Chart, templateName string) *chart.File {

	for _, template := range chart.Templates {
		if template.Name == templateName {
			return template
		}
	}
	return nil
}
