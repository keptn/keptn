package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/helm/pkg/proto/hapi/chart"
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
  replicas: 1
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
        image: mongo 
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

func TestIncreaseReplicaCount(t *testing.T) {

	const expectedDeploymentContent = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  name: carts
spec:
  replicas: 3
  selector:
    matchLabels:
      app: carts
  strategy:
    rollingUpdate:
      maxUnavailable: 0
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: carts
    spec:
      containers:
      - image: mongo
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 60
          periodSeconds: 10
          timeoutSeconds: 15
        name: carts
        ports:
        - containerPort: 8080
          name: http
          protocol: TCP
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
            memory: 2Gi
          requests:
            cpu: 250m
            memory: 1Gi
status: {}
`

	meta := &chart.Metadata{
		Name: "test.chart",
	}
	const templateFileName = "templates/deployment.yaml"
	inputChart := chart.Chart{Metadata: meta}
	template := &chart.Template{Name: templateFileName, Data: []byte(deploymentContent)}
	inputChart.Templates = append(inputChart.Templates, template)

	scaler := NewScaler()
	changedTemplates, err := scaler.increaseReplicaCount(&inputChart, 2)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(changedTemplates))
	assert.Equal(t, templateFileName, changedTemplates[0].Name)
	assert.Equal(t, expectedDeploymentContent, string(changedTemplates[0].Data))
}

func TestIncreaseReplicaCountWithMultipleDocuments(t *testing.T) {

	const expectedDeploymentContent = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  name: carts
spec:
  replicas: 3
  selector:
    matchLabels:
      app: carts
  strategy:
    rollingUpdate:
      maxUnavailable: 0
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: carts
    spec:
      containers:
      - image: mongo
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 60
          periodSeconds: 10
          timeoutSeconds: 15
        name: carts
        ports:
        - containerPort: 8080
          name: http
          protocol: TCP
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
            memory: 2Gi
          requests:
            cpu: 250m
            memory: 1Gi
status: {}`

	const expectedServiceContent = `---
apiVersion: v1
kind: Service
metadata:
  name: carts
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: carts
  type: LoadBalancer
`

	meta := &chart.Metadata{
		Name: "test.chart",
	}
	const templateFileName = "templates/deployment.yaml"
	inputChart := chart.Chart{Metadata: meta}
	template := &chart.Template{Name: templateFileName, Data: []byte(deploymentContent + "\n" + serviceContent)}
	inputChart.Templates = append(inputChart.Templates, template)
	emptyTemplate := &chart.Template{Name: "test.yaml", Data: []byte{}}
	inputChart.Templates = append(inputChart.Templates, emptyTemplate)

	scaler := NewScaler()
	changedTemplates, err := scaler.increaseReplicaCount(&inputChart, 2)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(changedTemplates))
	assert.Equal(t, templateFileName, changedTemplates[0].Name)
	assert.Equal(t, expectedDeploymentContent+"\n"+expectedServiceContent, string(changedTemplates[0].Data))
}
