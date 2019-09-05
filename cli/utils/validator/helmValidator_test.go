package validator

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/utils"
	"gotest.tools/assert"
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
const invalidValues = `
carts:
  image: docker.io/keptnexamples/carts:0.8.1
`

func TestValidateTemplateRequirements1(t *testing.T) {

	const testFile = "test.yaml"
	// Write test file
	d := []byte(service)

	if err := ioutil.WriteFile(testFile, d, 0644); err != nil {
		t.Error(err)
	}

	testedFiles := []string{testFile}
	res, err := validateTemplateRequirements(testedFiles)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res, false, "Deployment is missing")

	if err := os.Remove(testFile); err != nil {
		t.Error(err)
	}
}

func TestValidateTemplateRequirements2(t *testing.T) {

	const testFile = "test.yaml"
	// Write test file
	d := []byte(deployment)

	if err := ioutil.WriteFile(testFile, d, 0644); err != nil {
		t.Error(err)
	}

	testedFiles := []string{testFile}
	res, err := validateTemplateRequirements(testedFiles)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res, false, "Service is missing")

	if err := os.Remove(testFile); err != nil {
		t.Error(err)
	}
}

func TestValidateTemplateRequirements3(t *testing.T) {

	const testFile = "test.yaml"
	// Write test file
	d := []byte(service + deployment)

	if err := ioutil.WriteFile(testFile, d, 0644); err != nil {
		t.Error(err)
	}

	testedFiles := []string{testFile}
	res, err := validateTemplateRequirements(testedFiles)
	if err != nil {
		t.Error(err)
	}
	assert.Assert(t, res, "Template should be valid.")

	if err := os.Remove(testFile); err != nil {
		t.Error(err)
	}
}

func TestValidateValuesRequirements(t *testing.T) {

	const testFile = "test.yaml"
	// Write test file
	d := []byte(invalidValues)

	if err := ioutil.WriteFile(testFile, d, 0644); err != nil {
		t.Error(err)
	}

	res, err := validateValues(testFile)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res, false, "Values should be invalid.")

	if err := os.Remove(testFile); err != nil {
		t.Error(err)
	}
}

func TestValidateHelmChart1(t *testing.T) {

	err := os.MkdirAll("tar/carts/templates", 0777)
	check(err)
	err = ioutil.WriteFile("tar/carts/values.yaml", []byte(values), 0644)
	check(err)
	err = ioutil.WriteFile("tar/carts/templates/deployment.yml", []byte(deployment), 0644)
	check(err)
	err = ioutil.WriteFile("tar/carts/templates/service.yaml", []byte(service), 0644)
	check(err)

	f, err := os.Create("carts.tgz")
	check(err)
	err = utils.Tar("tar", f)
	check(err)

	dat, err := ioutil.ReadFile("carts.tgz")
	if err != nil {
		t.Error(err)
	}

	res, err := ValidateHelmChart(dat)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res, true, "Helm chart should be valid")

	if err := os.Remove("carts.tgz"); err != nil {
		t.Error(err)
	}
	if err := os.RemoveAll("tar"); err != nil {
		t.Error(err)
	}
}

func TestValidateHelmChart2(t *testing.T) {

	err := os.MkdirAll("tar/carts/templates", 0777)
	check(err)
	err = ioutil.WriteFile("tar/carts/values.yaml", []byte(values), 0644)
	check(err)
	err = ioutil.WriteFile("tar/carts/templates/deployment.yml", []byte(deployment), 0644)
	check(err)

	f, err := os.Create("carts.tgz")
	check(err)
	err = utils.Tar("tar", f)
	check(err)

	dat, err := ioutil.ReadFile("carts.tgz")
	if err != nil {
		t.Error(err)
	}

	res, err := ValidateHelmChart(dat)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res, false, "Helm chart should be valid")

	if err := os.Remove("carts.tgz"); err != nil {
		t.Error(err)
	}
	if err := os.RemoveAll("tar"); err != nil {
		t.Error(err)
	}
}

func TestValidateHelmChart3(t *testing.T) {

	err := os.MkdirAll("tar/carts/templates", 0777)
	check(err)
	err = ioutil.WriteFile("tar/carts/values.yaml", []byte(values+deployment), 0644)
	check(err)
	err = ioutil.WriteFile("tar/carts/templates/deployment.yml", []byte(deployment), 0644)
	check(err)

	f, err := os.Create("carts.tgz")
	check(err)
	err = utils.Tar("tar", f)
	check(err)

	dat, err := ioutil.ReadFile("carts.tgz")
	if err != nil {
		t.Error(err)
	}

	res, err := ValidateHelmChart(dat)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res, false, "Helm chart should be valid")

	if err := os.Remove("carts.tgz"); err != nil {
		t.Error(err)
	}
	if err := os.RemoveAll("tar"); err != nil {
		t.Error(err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
