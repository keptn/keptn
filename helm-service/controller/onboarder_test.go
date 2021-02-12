package controller

import (
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/golang/mock/gomock"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/mocks"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/chart"
	"testing"
)

//
//import (
//	"encoding/base64"
//	"errors"
//	cloudevents "github.com/cloudevents/sdk-go/v2"
//	"github.com/cloudevents/sdk-go/v2/event"
//	"github.com/golang/mock/gomock"
//	keptnevents "github.com/keptn/go-utils/pkg/lib"
//	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
//	"github.com/keptn/keptn/helm-service/mocks"
//	keptnutils "github.com/keptn/kubernetes-utils/pkg"
//	"github.com/stretchr/testify/assert"
//	"helm.sh/helm/v3/pkg/chart"
//	"testing"
//
//	"github.com/keptn/keptn/helm-service/pkg/helm"
//
//	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
//
//	"github.com/keptn/go-utils/pkg/api/models"
//)
//
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

func TestOnboardGeneratedChart_withDirectStrategy(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockNamespaceManager := mocks.NewMockINamespaceManager(ctrl)
	mockChartStorer := mocks.NewMockIChartStorer(ctrl)
	mockChartGenerator := mocks.NewMockChartGenerator(ctrl)
	mockChartPackager := mocks.NewMockIChartPackager(ctrl)
	mockBaseHandler := NewMockedHandler(createKeptn(), "")
	testChart := createChart()

	expectedStoreChartOpts := keptnutils.StoreChartOptions{
		Project:   "myproject",
		Service:   "myservice",
		Stage:     "mydev",
		ChartName: "myservice-generated",
		HelmChart: []byte("chart_as_bytes"),
	}
	mockChartStorer.EXPECT().Store(gomock.Eq(expectedStoreChartOpts))
	mockChartGenerator.EXPECT().GenerateMeshChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&testChart, nil)
	mockChartPackager.EXPECT().Package(&testChart).Return([]byte("chart_as_bytes"), nil)

	instance := NewOnboarder(mockBaseHandler, mockNamespaceManager, mockChartStorer, mockChartGenerator, mockChartPackager)
	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Direct)
	assert.Equal(t, &testChart, generatedChart)
	assert.Nil(t, err)
}

func TestOnboardGeneratedChart_withUserManagedStrategy(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockNamespaceManager := mocks.NewMockINamespaceManager(ctrl)
	mockChartStorer := mocks.NewMockIChartStorer(ctrl)
	mockChartGenerator := mocks.NewMockChartGenerator(ctrl)
	mockChartPackager := mocks.NewMockIChartPackager(ctrl)
	mockBaseHandler := NewMockedHandler(createKeptn(), "")

	mockChartGenerator.EXPECT().GenerateDuplicateChart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockNamespaceManager.EXPECT().InjectIstio(gomock.Any(), gomock.Any()).Times(0)
	mockChartPackager.EXPECT().Package(gomock.Any()).Times(0)
	mockChartStorer.EXPECT().Store(gomock.Any()).Times(0)
	instance := NewOnboarder(mockBaseHandler, mockNamespaceManager, mockChartStorer, mockChartGenerator, mockChartPackager)
	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.UserManaged)
	assert.Equal(t, &chart.Chart{}, generatedChart)
	assert.Nil(t, err)
}

func TestOnboardGeneratedChart_withDirectStrategy_chartGenerationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockNamespaceManager := mocks.NewMockINamespaceManager(ctrl)
	mockChartStorer := mocks.NewMockIChartStorer(ctrl)
	mockChartGenerator := mocks.NewMockChartGenerator(ctrl)
	mockChartPackager := mocks.NewMockIChartPackager(ctrl)
	mockBaseHandler := NewMockedHandler(createKeptn(), "")
	instance := NewOnboarder(mockBaseHandler, mockNamespaceManager, mockChartStorer, mockChartGenerator, mockChartPackager)

	mockChartGenerator.EXPECT().GenerateMeshChart(helmManifestResource, "myproject", "mydev", "myservice").Return(nil, errors.New("chart generation failed"))

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Direct)
	assert.NotNil(t, err)
	assert.Nil(t, generatedChart)

}

func TestOnboardGeneratedChart_withDirectStrategy_chartPackagingFails(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockNamespaceManager := mocks.NewMockINamespaceManager(ctrl)
	mockChartStorer := mocks.NewMockIChartStorer(ctrl)
	mockChartGenerator := mocks.NewMockChartGenerator(ctrl)
	mockChartPackager := mocks.NewMockIChartPackager(ctrl)
	mockBaseHandler := NewMockedHandler(createKeptn(), "")
	testChart := createChart()
	instance := NewOnboarder(mockBaseHandler, mockNamespaceManager, mockChartStorer, mockChartGenerator, mockChartPackager)

	mockChartGenerator.EXPECT().GenerateMeshChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&testChart, nil)
	mockChartPackager.EXPECT().Package(&testChart).Return(nil, errors.New("chart packaging failed"))
	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Direct)

	assert.Nil(t, generatedChart)
	assert.NotNil(t, err)

}

func TestOnboardGeneratedChart_withDuplicateStrategy(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockNamespaceManager := mocks.NewMockINamespaceManager(ctrl)
	mockChartStorer := mocks.NewMockIChartStorer(ctrl)
	mockChartGenerator := mocks.NewMockChartGenerator(ctrl)
	mockChartPackager := mocks.NewMockIChartPackager(ctrl)
	mockBaseHandler := NewMockedHandler(createKeptn(), "")
	testChart := createChart()
	instance := NewOnboarder(mockBaseHandler, mockNamespaceManager, mockChartStorer, mockChartGenerator, mockChartPackager)

	expectedStoreChartOpts := keptnutils.StoreChartOptions{
		Project:   "myproject",
		Service:   "myservice",
		Stage:     "mydev",
		ChartName: "myservice-generated",
		HelmChart: []byte("chart_as_bytes"),
	}

	mockChartGenerator.EXPECT().GenerateDuplicateChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&testChart, nil)
	mockNamespaceManager.EXPECT().InjectIstio("myproject", "mydev")
	mockChartPackager.EXPECT().Package(&testChart).Return([]byte("chart_as_bytes"), nil)
	mockChartStorer.EXPECT().Store(gomock.Eq(expectedStoreChartOpts))

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Duplicate)
	assert.Equal(t, &testChart, generatedChart)
	assert.Nil(t, err)
}

//
func TestOnboardGeneratedChart_injectingIstioConfigFails(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockNamespaceManager := mocks.NewMockINamespaceManager(ctrl)
	mockChartStorer := mocks.NewMockIChartStorer(ctrl)
	mockChartGenerator := mocks.NewMockChartGenerator(ctrl)
	mockChartPackager := mocks.NewMockIChartPackager(ctrl)
	mockBaseHandler := NewMockedHandler(createKeptn(), "")
	testChart := createChart()
	instance := NewOnboarder(mockBaseHandler, mockNamespaceManager, mockChartStorer, mockChartGenerator, mockChartPackager)

	mockChartGenerator.EXPECT().GenerateDuplicateChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&testChart, nil)
	mockNamespaceManager.EXPECT().InjectIstio("myproject", "mydev").Return(errors.New("failed to inject istio"))

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Duplicate)
	assert.NotNil(t, err)
	assert.Nil(t, generatedChart)
}

//
func TestOnboardGeneratedChart_chartGenerationFails(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockNamespaceManager := mocks.NewMockINamespaceManager(ctrl)
	mockChartStorer := mocks.NewMockIChartStorer(ctrl)
	mockChartGenerator := mocks.NewMockChartGenerator(ctrl)
	mockChartPackager := mocks.NewMockIChartPackager(ctrl)
	mockBaseHandler := NewMockedHandler(createKeptn(), "")
	instance := NewOnboarder(mockBaseHandler, mockNamespaceManager, mockChartStorer, mockChartGenerator, mockChartPackager)

	mockChartGenerator.EXPECT().GenerateDuplicateChart(helmManifestResource, "myproject", "mydev", "myservice").Return(nil, errors.New("chart generation failed"))

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Duplicate)
	assert.NotNil(t, err)
	assert.Nil(t, generatedChart)

}

//
func TestOnboardGeneratedChart_chartStorageFails(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockNamespaceManager := mocks.NewMockINamespaceManager(ctrl)
	mockChartStorer := mocks.NewMockIChartStorer(ctrl)
	mockChartGenerator := mocks.NewMockChartGenerator(ctrl)
	mockChartPackager := mocks.NewMockIChartPackager(ctrl)
	mockBaseHandler := NewMockedHandler(createKeptn(), "")
	testChart := createChart()
	instance := NewOnboarder(mockBaseHandler, mockNamespaceManager, mockChartStorer, mockChartGenerator, mockChartPackager)

	expectedStoreChartOpts := keptnutils.StoreChartOptions{
		Project:   "myproject",
		Service:   "myservice",
		Stage:     "mydev",
		ChartName: "myservice-generated",
		HelmChart: []byte("chart_as_bytes"),
	}
	mockChartGenerator.EXPECT().GenerateDuplicateChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&testChart, nil)
	mockNamespaceManager.EXPECT().InjectIstio("myproject", "mydev")
	mockChartPackager.EXPECT().Package(&testChart).Return([]byte("chart_as_bytes"), nil)
	mockChartStorer.EXPECT().Store(gomock.Eq(expectedStoreChartOpts)).Return("", errors.New("storing chart failed"))

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Duplicate)
	assert.NotNil(t, err)
	assert.Nil(t, generatedChart)
}

func createUnparsableEvent() event.Event {
	testEvent := cloudevents.NewEvent()
	_ = testEvent.SetData(cloudevents.ApplicationJSON, "WEIRD_JSON_CONTENT")
	testEvent.SetID("EVENT_ID")
	return testEvent
}

func createChart() chart.Chart {
	return chart.Chart{
		Metadata: &chart.Metadata{
			APIVersion: "v2",
			Name:       "myservice-generated",
			Keywords:   []string{"deployment_strategy=" + keptnevents.Direct.String()},
			Version:    "0.1.0",
		},
	}
}

func createKeptnBaseHandlerMock() Handler {
	return &MockedHandler{}
}

//
func createKeptn() *keptnv2.Keptn {
	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})
	return keptn
}
