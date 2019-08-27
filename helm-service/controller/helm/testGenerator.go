package helm

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

const service = `--- 
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

const deployment = `--- 
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

const values = `
image: docker.io/keptnexamples/carts:0.8.1
`

const chart = `
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
	err := os.MkdirAll("tar/carts/templates", 0777)
	check(err, t)
	err = ioutil.WriteFile("tar/carts/Chart.yaml", []byte(chart), 0644)
	check(err, t)
	err = ioutil.WriteFile("tar/carts/values.yaml", []byte(values), 0644)
	check(err, t)
	err = ioutil.WriteFile("tar/carts/templates/deployment.yml", []byte(deployment), 0644)
	check(err, t)
	err = ioutil.WriteFile("tar/carts/templates/service.yaml", []byte(service), 0644)
	check(err, t)

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err = keptnutils.Tar("tar", writer)
	check(err, t)

	err = os.RemoveAll("tar")
	check(err, t)
	writer.Flush()
	return b.Bytes()
}
