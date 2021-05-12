package helm

import (
	"testing"

	"gotest.tools/assert"
)

const helmManifestResource = `
---
# Source: carts/templates/service.yaml
apiVersion: v1
kind: Service
metadata: 
  name: carts
spec: 
  type: ClusterIP
  ports: 
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8080
  selector: 
    app: carts
---
# Source: carts/templates/deployment.yaml
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
        image: "docker.io/keptnexamples/carts:0.10.1"
        imagePullPolicy: IfNotPresent
        ports:
        - name: http
          protocol: TCP
          containerPort: 8080
        env:
        - name: DT_CUSTOM_PROP
          value: "keptn_project=sockshop keptn_service=carts keptn_stage=dev keptn_deployment=direct"
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: "metadata.name"
        - name: DEPLOYMENT_NAME
          valueFrom:
            fieldRef:
              fieldPath: "metadata.labels['deployment']"
        - name: CONTAINER_IMAGE
          value: "docker.io/keptnexamples/carts:0.10.1"
        - name: KEPTN_PROJECT
          value: "carts"
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
              cpu: 1000m
              memory: 2048Mi
          requests:
              cpu: 500m
              memory: 1024Mi`

func TestGetServices(t *testing.T) {

	services := GetServices(helmManifestResource)
	assert.Equal(t, 1, len(services))
}

func TestGetDeployments(t *testing.T) {

	deployments := GetDeployments(helmManifestResource)
	assert.Equal(t, 1, len(deployments))
}

func TestGetReleaseName(t *testing.T) {
	var tests = []struct {
		project   string
		stage     string
		service   string
		generated bool
		out       string
	}{
		// legacy support: make sure helm chart release name looks the same way as before
		{"sockshop", "dev", "carts", false, "sockshop-dev-carts"},
		{"sockshop", "dev", "carts", true, "sockshop-dev-carts-generated"},
		{"sockshop", "dev", "carts-db", false, "sockshop-dev-carts-db"},
		{"sockshop", "dev", "carts-db", true, "sockshop-dev-carts-db-generated"},

		// now make project and stage name much longer, such that adding -generated would result in a name >= 53 chars
		{"sockshop-enhanced-version", "development", "carts", false, "sockshop-enhanced-version-development-carts"},
		{"sockshop-enhanced-version", "development", "carts", true, "development-carts-generated"},

		// In addition, use a very long service name, such that neither project nor stage name can be used
		{"sockshop-enhanced-version", "development", "my-carts-service-is-the-best", false, "development-my-carts-service-is-the-best"},
		{"sockshop-enhanced-version", "development", "my-carts-service-is-the-best", true, "development-my-carts-service-is-the-best-generated"},

		// finally, test the case where the service name itself is so big that it needs to be sliced to fit the -generated suffix
		// Note: this should really never be the case, but we will cover it anyway, just to be safe
		{"sockshop", "dev", "this-is-my-very-very-very-very-very-long-servicename", false, "this-is-my-very-very-very-very-very-long-servicename"},
		{"sockshop", "dev", "this-is-my-very-very-very-very-very-long-servicename", true, "this-is-my-very-very-very-very-very-long-s-generated"},
	}

	for _, tt := range tests {
		t.Run(tt.out, func(t *testing.T) {
			s := GetReleaseName(tt.project, tt.stage, tt.service, tt.generated)
			if s != tt.out {
				t.Errorf("got %q, want %q", s, tt.out)
			}
			// also, verify that s is <= 53 characters (just to be sure)
			if len(s) >= 53 {
				t.Errorf("len(%q) >= 53, but should be < 53", s)
			}
		})
	}
}
