package helm

import (
	"strings"

	"helm.sh/helm/v3/pkg/chart"
)

// HelmMockExecutor mocks Helm operations
type HelmMockExecutor struct {
}

// NewHelmMockExecutor creates a new HelmMockExecutor
func NewHelmMockExecutor() *HelmMockExecutor {
	return &HelmMockExecutor{}
}

const userService = `--- 
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

const userDeployment = `--- 
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
        image: "docker.io/keptnexamples/carts:0.8.1"
        imagePullPolicy: IfNotPresent        
`

// GeneratedCanaryDestinationRule is a DestinationRule manifest for tests
const GeneratedCanaryDestinationRule = `---
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  creationTimestamp: null
  name: carts-canary
spec:
  host: carts-canary.sockshop-NAMESPACE_PLACEHOLDER.svc.cluster.local
`

// GeneratedPrimaryDestinationRule is a DestinationRule manifest for tests
const GeneratedPrimaryDestinationRule = `---
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  creationTimestamp: null
  name: carts-primary
spec:
  host: carts-primary.sockshop-NAMESPACE_PLACEHOLDER.svc.cluster.local
`

// GeneratedCanaryService is a Service manifest for tests
const GeneratedCanaryService = `---
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  name: carts-canary
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: carts
  type: ClusterIP
status:
  loadBalancer: {}
`

// GeneratedPrimaryService is a Service manifest for tests
const GeneratedPrimaryService = `---
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  name: carts-primary
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: carts-primary
  type: ClusterIP
status:
  loadBalancer: {}
`

// GeneratedVirtualService is a VirtualService manifest for tests
const GeneratedVirtualService = `---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  creationTimestamp: null
  name: carts
spec:
  gateways:
  - public-gateway.istio-system
  - mesh
  hosts:
  - carts.sockshop-NAMESPACE_PLACEHOLDER.demo.keptn.sh
  - carts
  http:
  - route:
    - destination:
        host: carts-canary.sockshop-NAMESPACE_PLACEHOLDER.svc.cluster.local
    - destination:
        host: carts-primary.sockshop-NAMESPACE_PLACEHOLDER.svc.cluster.local
      weight: 100
`

// GeneratedPrimaryDeployment is a Deployment manifest for tests
const GeneratedPrimaryDeployment = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  name: carts-primary
spec:
  replicas: 1
  selector:
    matchLabels:
      app: carts-primary
  strategy:
    rollingUpdate:
      maxUnavailable: 0
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: carts-primary
    spec:
      containers:
      - image: docker.io/keptnexamples/carts:0.8.1
        imagePullPolicy: IfNotPresent
        name: carts
        resources: {}
status: {}    
`

// GetManifest returns test/sample manifests
func (h *HelmMockExecutor) GetManifest(releaseName, namespace string) (string, error) {

	if strings.HasSuffix(releaseName, "-generated") {
		genManifests := GeneratedPrimaryDeployment + GeneratedCanaryService + GeneratedPrimaryService + GeneratedCanaryDestinationRule +
			GeneratedPrimaryDestinationRule + GeneratedVirtualService
		return strings.ReplaceAll(genManifests, "NAMESPACE_PLACEHOLDER", namespace), nil
	}
	return userDeployment + userService, nil
}

// UpgradeChart does not execute any action
func (h *HelmMockExecutor) UpgradeChart(ch *chart.Chart, releaseName, namespace string, vals map[string]interface{}) error {
	return nil
}

func (h *HelmMockExecutor) UninstallRelease(releaseName, namespace string) error {
	return nil
}
