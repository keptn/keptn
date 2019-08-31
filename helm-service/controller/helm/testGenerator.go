package helm

import (
	"io/ioutil"
	"os"
	"testing"

	"k8s.io/helm/pkg/chartutil"
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
        env:
        - name: DT_TAGS
          value: "application={{ .Chart.Name }}"
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: "metadata.name"
        - name: DEPLOYMENT_NAME
          valueFrom:
            fieldRef:
              fieldPath: "metadata.labels['deployment']"
        - name: CONTAINER_IMAGE
          value: "{{ .Values.image }}"
        - name: KEPTN_PROJECT
          value: "{{ .Chart.Name }}"
        - name: KEPTN_STAGE
          valueFrom:
            fieldRef:
              fieldPath: "metadata.namespace"
        - name: KEPTN_SERVICE
          value: "carts"
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

	ch, err := chartutil.LoadDir("carts")
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
