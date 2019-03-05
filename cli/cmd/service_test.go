package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
)

func init() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
}

func TestOnboardServiceCmd(t *testing.T) {

	credentialmanager.MockCreds = true

	// Write temporary files
	const tmpValues = "valuesTest.tpl"
	valuesContent := `replicaCount: 1
image:
  repository: null
  tag: null
  pullPolicy: IfNotPresent
service:
  name: carts
  type: LoadBalancer
  externalPort: 8080
  internalPort: 8080
resources:
  limits:
    cpu: 100m
    memory: 128Mi
  requests:
    cpu: 100m
    memory: 128Mi`

	const tmpDeployment = "deploymentTest.tpl"
	deploymentContent := `apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: {{ .Chart.Name }}-SERVICE_PLACEHOLDER_DEC
  labels:
    app: {{ .Chart.Name }}-selector-SERVICE_PLACEHOLDER_DEC
    chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
    release: "{{ .Release.Name }}"
    heritage: "{{ .Release.Service }}"
spec:
  replicas: {{ .Values.SERVICE_PLACEHOLDER_C.replicaCount }}
  template:
    metadata:
      labels:
        app: {{ .Chart.Name }}-selector-SERVICE_PLACEHOLDER_DEC
        deployment: SERVICE_PLACEHOLDER_DEC
    spec:
      containers:
      - name: {{ .Chart.Name }}
        image: "{{ .Values.SERVICE_PLACEHOLDER_C.image.repository }}:{{ .Values.SERVICE_PLACEHOLDER_C.image.tag }}"
        imagePullPolicy: {{ .Values.SERVICE_PLACEHOLDER_C.image.pullPolicy }}
        ports:
        - name: internalport
          containerPort: {{ .Values.SERVICE_PLACEHOLDER_C.service.internalPort }}
        resources: {{ toYaml .Values.resources | indent 12 }}`

	const tmpService = "serviceTest.tpl"
	serviceContent := `apiVersion: v1
kind: Service
metadata:
  name: {{ .Chart.Name }}-SERVICE_PLACEHOLDER_DEC
  labels:
    chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
spec:
  type: {{ .Values.SERVICE_PLACEHOLDER_C.service.type }}
  ports:
  - port: {{ .Values.SERVICE_PLACEHOLDER_C.service.externalPort }}
    targetPort: {{ .Values.SERVICE_PLACEHOLDER_C.service.internalPort }}
    protocol: TCP
    name: {{ .Values.SERVICE_PLACEHOLDER_C.service.name }}
  selector:
    app: {{ .Chart.Name }}-selector-SERVICE_PLACEHOLDER_DEC`

	ioutil.WriteFile(tmpValues, []byte(valuesContent), 0644)
	ioutil.WriteFile(tmpDeployment, []byte(deploymentContent), 0644)
	ioutil.WriteFile(tmpService, []byte(serviceContent), 0644)

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"onboard",
		"service",
		"--project=carts",
		fmt.Sprintf("--deployment=%s", tmpDeployment),
		fmt.Sprintf("--values=%s", tmpValues),
		// fmt.Sprintf("--service=%s", tmpService),
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	// Delete temporary shipyard.yml file
	os.Remove(tmpValues)
	os.Remove(tmpDeployment)
	os.Remove(tmpService)

	if err != nil {
		t.Errorf("An error occured %v", err)
	}
}
