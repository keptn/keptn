package validator

import (
	"io/ioutil"
	"os"
	"testing"

	keptnutils "github.com/keptn/kubernetes-utils/pkg"

	"github.com/stretchr/testify/assert"
)

const defaultService = `--- 
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

const defaultDeployment = `--- 
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
`

const defaultValues = `
image: docker.io/keptnexamples/carts:0.8.1
replicas: 1
`

const defaultChart = `
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

func TestCheckValues(t *testing.T) {

	// wrong name of image property
	const image = `
dockerImage: docker.io/keptnexamples/carts:0.8.1
`

	invalidValues := []string{image}

	for _, invalidValue := range invalidValues {
		err := os.MkdirAll("carts/templates", 0777)
		check(err, t)
		err = ioutil.WriteFile("carts/Chart.yaml", []byte(defaultChart), 0644)
		check(err, t)
		err = ioutil.WriteFile("carts/values.yaml", []byte(invalidValue), 0644)
		check(err, t)
		err = ioutil.WriteFile("carts/templates/deployment.yml", []byte(defaultDeployment), 0644)
		check(err, t)
		err = ioutil.WriteFile("carts/templates/service.yaml", []byte(defaultService), 0644)
		check(err, t)

		ch, err := keptnutils.LoadChartFromPath("carts")
		check(err, t)

		res, err := ValidateHelmChart(ch, "carts")
		check(err, t)
		assert.False(t, res)
		os.RemoveAll("carts")
	}
}

func TestCheckService(t *testing.T) {

	const invalidService1 = `--- 
apiVersion: v1
kind: Service2
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

	const invalidService2 = `--- 
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
`
	const invalidService3 = `--- 
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
    app2: carts
`

	const invalidService4 = `---`

	invalidServices := []string{invalidService1, invalidService2, invalidService3, invalidService4}

	for _, invalidService := range invalidServices {
		err := os.MkdirAll("carts/templates", 0777)
		check(err, t)
		err = ioutil.WriteFile("carts/Chart.yaml", []byte(defaultChart), 0644)
		check(err, t)
		err = ioutil.WriteFile("carts/values.yaml", []byte(defaultValues), 0644)
		check(err, t)
		err = ioutil.WriteFile("carts/templates/deployment.yml", []byte(defaultDeployment), 0644)
		check(err, t)
		err = ioutil.WriteFile("carts/templates/service.yaml", []byte(invalidService), 0644)
		check(err, t)

		ch, err := keptnutils.LoadChartFromPath("carts")
		check(err, t)

		res, err := ValidateHelmChart(ch, "carts")
		check(err, t)
		assert.False(t, res)
		os.RemoveAll("carts")
	}
}

func TestCheckDeployment(t *testing.T) {

	const invalidDeployment1 = `--- 
apiVersion: apps/v1
kind: Deployment2
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
`

	const invalidDeployment2 = `--- 
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
  template:
    metadata:
      labels:
        app: carts
    spec:
      containers:
      - name: carts
        image: "{{ .Values.image }}"
        imagePullPolicy: IfNotPresent        
`

	const invalidDeployment3 = `--- 
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
      app2: carts
  template:
    metadata:
      labels:
        app: carts
    spec:
      containers:
      - name: carts
        image: "{{ .Values.image }}"
        imagePullPolicy: IfNotPresent        
`

	const invalidDeployment4 = `--- 
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
    spec:
      containers:
      - name: carts
        image: "{{ .Values.image }}"
        imagePullPolicy: IfNotPresent        
`

	const invalidDeployment5 = `--- 
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
        app2: carts
    spec:
      containers:
      - name: carts
        image: "{{ .Values.image }}"
        imagePullPolicy: IfNotPresent        
`

	const invalidDeployment6 = `---`

	invalidDeployments := []string{invalidDeployment1, invalidDeployment2, invalidDeployment3,
		invalidDeployment4, invalidDeployment5, invalidDeployment6}

	for _, invalidDeployment := range invalidDeployments {
		err := os.MkdirAll("carts/templates", 0777)
		check(err, t)
		err = ioutil.WriteFile("carts/Chart.yaml", []byte(defaultChart), 0644)
		check(err, t)
		err = ioutil.WriteFile("carts/values.yaml", []byte(defaultValues), 0644)
		check(err, t)
		err = ioutil.WriteFile("carts/templates/deployment.yml", []byte(invalidDeployment), 0644)
		check(err, t)
		err = ioutil.WriteFile("carts/templates/service.yaml", []byte(defaultService), 0644)
		check(err, t)

		ch, err := keptnutils.LoadChartFromPath("carts")
		check(err, t)

		res, err := ValidateHelmChart(ch, "carts")
		check(err, t)
		assert.False(t, res)
		os.RemoveAll("carts")
	}
}

func TestTemplateFileNames(t *testing.T) {
	err := os.MkdirAll("carts/templates", 0777)
	check(err, t)

	err = ioutil.WriteFile("carts/Chart.yaml", []byte(defaultChart), 0644)
	check(err, t)
	err = ioutil.WriteFile("carts/values.yaml", []byte(defaultValues), 0644)
	check(err, t)
	err = ioutil.WriteFile("carts/templates/deployment.yml", []byte(defaultDeployment), 0644)
	check(err, t)
	err = ioutil.WriteFile("carts/templates/service.yaml", []byte(defaultService), 0644)
	check(err, t)

	empty := make([]byte, 0)
	err = ioutil.WriteFile("carts/templates/test-istio-destinationrule.yaml", empty, 0644)
	check(err, t)

	ch, err := keptnutils.LoadChartFromPath("carts")
	check(err, t)

	res, err := ValidateHelmChart(ch, "carts")
	check(err, t)
	assert.False(t, res)
	os.RemoveAll("carts")
}
