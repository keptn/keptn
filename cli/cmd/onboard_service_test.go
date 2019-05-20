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

const manifestContent = `---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: carts
  namespace: dev
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: carts
        version: v1
    spec:
      containers:
      - name: carts
        image: dynatracesockshop/carts:0.6.0
        env:
        - name: JAVA_OPTS
          value: -Xms128m -Xmx512m -XX:PermSize=128m -XX:MaxPermSize=128m -XX:+UseG1GC -Djava.security.egd=file:/dev/urandom
        - name: DT_TAGS
          value: "application=sockshop"
        - name: DT_CUSTOM_PROP
          value: "SERVICE_TYPE=BACKEND"
        resources:
          limits:
            cpu: 500m
            memory: 1024Mi
          requests:
            cpu: 400m
            memory: 768Mi
        ports:
        - containerPort: 8080
        volumeMounts:
        - mountPath: /tmp
          name: tmp-volume
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
      volumes:
        - name: tmp-volume
          emptyDir:
            medium: Memory
      nodeSelector:
        beta.kubernetes.io/os: linux
---
apiVersion: v1
kind: Service
metadata:
  name: carts
  labels:
    app: carts
  namespace: dev
spec:
  ports:
  - name: http
    port: 80
    targetPort: 8080
  selector:
    app: carts
  type: LoadBalancer`

const valuesContent = `replicaCount: 1
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

const valuesIncorrectContent = `replicaCount: 1
image:
  repository: null
  tag: null
  pullPolicy: IfNotPresent
service:
  name: carts.service
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

const deploymentContent = `apiVersion: extensions/v1beta1
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

const serviceContent = `apiVersion: v1
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

func init() {
	utils.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func TestOnboardServiceCmdUsingHelm(t *testing.T) {

	credentialmanager.MockCreds = true

	// Write temporary files
	const tmpValues = "valuesTest.tpl"
	const tmpDeployment = "deploymentTest.tpl"
	const tmpService = "serviceTest.tpl"

	ioutil.WriteFile(tmpValues, []byte(valuesContent), 0644)
	ioutil.WriteFile(tmpDeployment, []byte(deploymentContent), 0644)
	ioutil.WriteFile(tmpService, []byte(serviceContent), 0644)

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"onboard",
		"service",
		"--project=carts",
		fmt.Sprintf("--values=%s", tmpValues),
		fmt.Sprintf("--deployment=%s", tmpDeployment),
		fmt.Sprintf("--service=%s", tmpService),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	*valuesFilePath = ""
	*serviceFilePath = ""
	*deploymentFilePath = ""

	// Delete temporary files
	os.Remove(tmpValues)
	os.Remove(tmpDeployment)
	os.Remove(tmpService)

	if err != nil {
		t.Errorf("An error occured %v", err)
	}
}

func TestOnboardServiceCmdUsingHelmIncorrectName(t *testing.T) {

	credentialmanager.MockCreds = true

	// Write temporary files
	const tmpValues = "valuesTest.tpl"
	const tmpDeployment = "deploymentTest.tpl"
	const tmpService = "serviceTest.tpl"

	ioutil.WriteFile(tmpValues, []byte(valuesIncorrectContent), 0644)
	ioutil.WriteFile(tmpDeployment, []byte(deploymentContent), 0644)
	ioutil.WriteFile(tmpService, []byte(serviceContent), 0644)

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"onboard",
		"service",
		"--project=carts",
		fmt.Sprintf("--values=%s", tmpValues),
		fmt.Sprintf("--deployment=%s", tmpDeployment),
		fmt.Sprintf("--service=%s", tmpService),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	*valuesFilePath = ""
	*serviceFilePath = ""
	*deploymentFilePath = ""

	// Delete temporary files
	os.Remove(tmpValues)
	os.Remove(tmpDeployment)
	os.Remove(tmpService)

	if err != nil {
		if !utils.ErrorContains(err, "Service name as defined in the .yaml file includes invalid characters or is not well-formed.") {
			t.Errorf("An error occured: %v", err)
		}
	} else {
		t.Fail()
	}
}

// func TestOnboardServiceCmdUsingManifest(t *testing.T) {

// 	credentialmanager.MockCreds = true

// 	// Write temporary files
// 	const tmpManifest = "manifestTest.yml"

// 	ioutil.WriteFile(tmpManifest, []byte(manifestContent), 0644)

// 	buf := new(bytes.Buffer)
// 	rootCmd.SetOutput(buf)

// 	args := []string{
// 		"onboard",
// 		"service",
// 		"--project=carts",
// 		fmt.Sprintf("--manifest=%s", tmpManifest),
// 	}
// 	rootCmd.SetArgs(args)
// 	err := rootCmd.Execute()

// 	*manifestFilePath = ""

// 	// Delete temporary files
// 	os.Remove(tmpManifest)

// 	if err != nil {
// 		t.Errorf("An error occured %v", err)
// 	}
// }

// func TestOnboardServiceCmdUsingInvalidArguments(t *testing.T) {

// 	credentialmanager.MockCreds = true

// 	// Write temporary files
// 	const tmpManifest = "manifestTest.yml"
// 	const tmpValues = "valuesTest.tpl"

// 	ioutil.WriteFile(tmpManifest, []byte(manifestContent), 0644)
// 	ioutil.WriteFile(tmpValues, []byte(valuesContent), 0644)

// 	buf := new(bytes.Buffer)
// 	rootCmd.SetOutput(buf)

// 	args := []string{
// 		"onboard",
// 		"service",
// 		"--project=carts",
// 		fmt.Sprintf("--manifest=%s", tmpManifest),
// 		fmt.Sprintf("--values=%s", tmpValues),
// 	}
// 	rootCmd.SetArgs(args)
// 	err := rootCmd.Execute()

// 	*manifestFilePath = ""
// 	*valuesFilePath = ""

// 	// Delete temporary files
// 	os.Remove(tmpManifest)
// 	os.Remove(tmpValues)

// 	expectedError := errors.New("Error specifying a Helm description as well as a k8s manifest. Only use one option")

// 	if err.Error() != expectedError.Error() {
// 		t.Error("Actual error does not match expected error")
// 	}
// }
