package controller

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/google/uuid"
	"github.com/keptn/keptn/helm-service/controller/mesh"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnmodels "github.com/keptn/go-utils/pkg/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/stretchr/testify/assert"
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

func check(e error, t *testing.T) {
	if e != nil {
		t.Error(e)
	}
}

func getHelmChartData(t *testing.T) []byte {
	err := os.MkdirAll("tar/carts/templates", 0777)
	check(err, t)
	err = ioutil.WriteFile("tar/carts/values.yaml", []byte(values), 0644)
	check(err, t)
	err = ioutil.WriteFile("tar/carts/templates/deployment.yml", []byte(deployment), 0644)
	check(err, t)
	err = ioutil.WriteFile("tar/carts/templates/service.yaml", []byte(service), 0644)
	check(err, t)

	f, err := os.Create("carts.tgz")
	check(err, t)
	err = keptnutils.Tar("tar", f)
	check(err, t)

	data, err := ioutil.ReadFile("carts.tgz")
	check(err, t)
	return data
}

func cleanupHelmChart(t *testing.T) {
	err := os.Remove("carts.tgz")
	check(err, t)
	err = os.RemoveAll("tar")
	check(err, t)
}

const configBaseURL = "localhost:8080"
const projectName = "sockshop"
const serviceName = "carts"
const stage1 = "dev"
const stage2 = "prod"

func createTestProjet(t *testing.T) {

	prjHandler := keptnutils.NewProjectHandler(configBaseURL)
	prj := keptnmodels.Project{ProjectName: projectName}
	respErr, err := prjHandler.CreateProject(prj)
	check(err, t)
	assert.Nil(t, respErr, "Creating a project failed")
}

func createTestStages(t *testing.T) {

	stageHandler := keptnutils.NewStageHandler(configBaseURL)
	respErr, err := stageHandler.CreateStage(projectName, stage1)
	check(err, t)
	assert.Nil(t, respErr, "Creating a stage failed")
	respErr, err = stageHandler.CreateStage(projectName, stage2)
	check(err, t)
	assert.Nil(t, respErr, "Creating a stage failed")
}

func TestDoOnboard(t *testing.T) {

	createTestProjet(t)
	createTestStages(t)

	data := getHelmChartData(t)
	defer cleanupHelmChart(t)

	ce := cloudevents.New("0.2")
	dataBytes, err := json.Marshal(keptnevents.ServiceCreateEventData{Project: projectName, Service: serviceName, HelmChart: data})
	check(err, t)
	ce.Data = dataBytes

	id := uuid.New().String()
	err = DoOnboard(ce, mesh.NewIstioMesh(), keptnutils.NewLogger(id, "service.create", "helm-service"), id, configBaseURL)

	check(err, t)
}
