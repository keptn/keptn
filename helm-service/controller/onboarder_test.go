package controller

import (
	"encoding/base64"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/golang/mock/gomock"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/helm-service/mocks"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/chart"
	"testing"

	"github.com/keptn/keptn/helm-service/pkg/helm"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	"github.com/keptn/go-utils/pkg/api/models"
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

//mockColleciton holds all the mocks
type mocksCollection struct {
	mockedBaseHandler      *MockedHandler
	mockedMesh             *mocks.MockMesh
	mockedProjectHandler   *mocks.MockIProjectHandler
	mockedNamespaceManager *mocks.MockINamespaceManager
	mockedStagesHandler    *mocks.MockIStagesHandler
	mockedServiceHandler   *mocks.MockIServiceHandler
	mockedChartStorer      *mocks.MockIChartStorer
	mockedChartGenerator   *mocks.MockChartGenerator
	mockedChartPackager    *mocks.MockIChartPackager
}

// testOnboarderCreator is an oboarder which has only mocked dependencies
type testOnboarderCreator struct {
}

// newTestOnboarderCreator creates an instance of testOnboarderCreator which uses only mocks
func newTestOnboarderCreator(t *testing.T, mockedBaseHandlerOptions ...MockedHandlerOption) (*gomock.Controller, Onboarder, mocksCollection) {

	ctrl := gomock.NewController(t)
	mockedBaseHandler := NewMockedHandler(createKeptn(), "", mockedBaseHandlerOptions...)
	mockedMesh := mocks.NewMockMesh(ctrl)
	mockedProjectHandler := mocks.NewMockIProjectHandler(ctrl)
	mockedNamespaceManager := mocks.NewMockINamespaceManager(ctrl)
	mockedStagesHandler := mocks.NewMockIStagesHandler(ctrl)
	mockedServiceHandler := mocks.NewMockIServiceHandler(ctrl)
	mockedChartStorer := mocks.NewMockIChartStorer(ctrl)
	mockedChartGenerator := mocks.NewMockChartGenerator(ctrl)
	mockedChartPackager := mocks.NewMockIChartPackager(ctrl)

	onboarder := onboarder{
		Handler:          mockedBaseHandler,
		mesh:             mockedMesh,
		projectHandler:   mockedProjectHandler,
		namespaceManager: mockedNamespaceManager,
		stagesHandler:    mockedStagesHandler,
		serviceHandler:   mockedServiceHandler,
		chartStorer:      mockedChartStorer,
		chartGenerator:   mockedChartGenerator,
		chartPackager:    mockedChartPackager,
	}

	mocksCol := mocksCollection{
		mockedBaseHandler:      mockedBaseHandler,
		mockedMesh:             mockedMesh,
		mockedProjectHandler:   mockedProjectHandler,
		mockedNamespaceManager: mockedNamespaceManager,
		mockedStagesHandler:    mockedStagesHandler,
		mockedServiceHandler:   mockedServiceHandler,
		mockedChartStorer:      mockedChartStorer,
		mockedChartGenerator:   mockedChartGenerator,
		mockedChartPackager:    mockedChartPackager,
	}
	return ctrl, &onboarder, mocksCol

}

func TestCreateOnboarder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedMesh := mocks.NewMockMesh(ctrl)
	mockedProjectHandler := mocks.NewMockIProjectHandler(ctrl)
	mockedNamespaceManager := mocks.NewMockINamespaceManager(ctrl)
	mockedStagesHandler := mocks.NewMockIStagesHandler(ctrl)
	mockedServiceHandler := mocks.NewMockIServiceHandler(ctrl)
	mockedChartStorer := mocks.NewMockIChartStorer(ctrl)
	mockedChartGenerator := mocks.NewMockChartGenerator(ctrl)
	mockedChartPackager := mocks.NewMockIChartPackager(ctrl)

	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})
	onboarder := NewOnboarder(
		keptn,
		mockedMesh,
		mockedProjectHandler,
		mockedNamespaceManager,
		mockedStagesHandler,
		mockedServiceHandler,
		mockedChartStorer,
		mockedChartGenerator,
		mockedChartPackager,
		"")

	assert.NotNil(t, onboarder)

}

func TestHandleEvent_WhenPassingUnparsableEvent_ThenHandleErrorIsCalled(t *testing.T) {

	ctrl, instance, _ := newTestOnboarderCreator(t)
	defer ctrl.Finish()
	instance.HandleEvent(createUnparsableEvent())
}

func TestHandleEvent_WhenHelmChartMissing_ThenNothingHappens(t *testing.T) {
	ctrl, instance, moqs := newTestOnboarderCreator(t)
	defer ctrl.Finish()

	instance.HandleEvent(cloudevents.NewEvent())
	assert.Empty(t, moqs.mockedBaseHandler.sentCloudEvents)
	assert.Empty(t, moqs.mockedBaseHandler.handledErrorEvents)
}

func TestHandleEvent_WhenNoProjectExists_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl, instance, moqs := newTestOnboarderCreator(t)
	defer ctrl.Finish()

	moqs.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(nil, &models.Error{Message: stringp("")})

	instance.HandleEvent(createEvent(t, "EVENT_ID"))
	assert.Equal(t, 1, len(moqs.mockedBaseHandler.handledErrorEvents))
}

func TestHandleEvent_WhenNoStagesDefined_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl, instance, moqs := newTestOnboarderCreator(t)
	defer ctrl.Finish()

	moqs.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{}, nil)
	moqs.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return([]*models.Stage{}, nil)

	instance.HandleEvent(createEvent(t, "EVENT_ID"))
	assert.Equal(t, 1, len(moqs.mockedBaseHandler.handledErrorEvents))
}

func TestHandleEvent_WhenStagesCannotBeFetched_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl, instance, moqs := newTestOnboarderCreator(t)
	defer ctrl.Finish()

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}

	moqs.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	moqs.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(nil, errors.New("unable to fetch stages"))

	instance.HandleEvent(createEvent(t, "EVENT_ID"))
	assert.Equal(t, 1, len(moqs.mockedBaseHandler.handledErrorEvents))
}

func TestHandleEvent_WhenInitNamespacesFails_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl, instance, moqs := newTestOnboarderCreator(t)
	defer ctrl.Finish()

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}

	moqs.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	moqs.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(stages, nil)
	moqs.mockedNamespaceManager.EXPECT().InitNamespaces("my-project", []string{"dev"}).Return(errors.New("namespace initialization failed"))

	instance.HandleEvent(createEvent(t, "EVENT_ID"))
	assert.Equal(t, 1, len(moqs.mockedBaseHandler.handledErrorEvents))
}

func TestHandleEvent_WhenPassingInvalidServiceName_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl, instance, moqs := newTestOnboarderCreator(t)
	defer ctrl.Finish()

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}

	moqs.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)

	instance.HandleEvent(createEventWith(t, "EVENT_ID", keptnv2.EventData{Project: "my-project", Stage: "dev", Service: "!§$%&/"}))
	assert.Equal(t, 1, len(moqs.mockedBaseHandler.handledErrorEvents))
}

func TestHandleEvent_WhenUnableToStoreChart_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl, instance, moqs := newTestOnboarderCreator(t)
	defer ctrl.Finish()

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}

	moqs.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	moqs.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(stages, nil)
	moqs.mockedNamespaceManager.EXPECT().InitNamespaces(gomock.Any(), gomock.Any())
	moqs.mockedServiceHandler.EXPECT().GetService(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	moqs.mockedChartStorer.EXPECT().Store(gomock.Any()).Return("", errors.New("unable to store chart"))

	instance.HandleEvent(createEvent(t, "EVENT_ID"))
	assert.Equal(t, 1, len(moqs.mockedBaseHandler.handledErrorEvents))
}

func TestHandleEvent_WhenSendingFinishedEventFails_ThenHandleErrorisCalled(t *testing.T) {

	opt := func(o *MockedHandlerOptions) {
		o.SendEventBehavior = func(eventType string) bool {
			if eventType == "sh.keptn.event.service.create.finished" {
				return false
			}
			return true
		}
	}

	ctrl, instance, moqs := newTestOnboarderCreator(t, opt)
	defer ctrl.Finish()

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}

	moqs.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	moqs.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(stages, nil)
	moqs.mockedNamespaceManager.EXPECT().InitNamespaces("my-project", []string{"dev"})
	moqs.mockedServiceHandler.EXPECT().GetService("my-project", "dev", "carts").Return(nil, nil)
	moqs.mockedChartStorer.EXPECT().Store(gomock.Any()).Return("", nil)

	instance.HandleEvent(createEvent(t, "EVENT_ID"))
	assert.Equal(t, 1, len(moqs.mockedBaseHandler.handledErrorEvents))
}

func TestHandleEvent_OnboardsService(t *testing.T) {
	ctrl, instance, moqs := newTestOnboarderCreator(t)
	defer ctrl.Finish()

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}

	moqs.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	moqs.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(stages, nil)
	moqs.mockedNamespaceManager.EXPECT().InitNamespaces("my-project", []string{"dev"})
	moqs.mockedServiceHandler.EXPECT().GetService("my-project", "dev", "carts").Return(nil, nil)
	moqs.mockedChartStorer.EXPECT().Store(gomock.Any()).Return("", nil)

	instance.HandleEvent(createEvent(t, "EVENT_ID"))
	assert.Equal(t, 1, len(moqs.mockedBaseHandler.sentCloudEvents))
}

func TestOnboardGeneratedChart_withDirectStrategy(t *testing.T) {
	ctrl, instance, moqs := newTestOnboarderCreator(t)
	defer ctrl.Finish()

	testChart := createChart()

	expectedStoreChartOpts := keptnutils.StoreChartOptions{
		Project:   "myproject",
		Service:   "myservice",
		Stage:     "mydev",
		ChartName: "myservice-generated",
		HelmChart: []byte("chart_as_bytes"),
	}
	moqs.mockedChartStorer.EXPECT().Store(gomock.Eq(expectedStoreChartOpts))
	moqs.mockedChartGenerator.EXPECT().GenerateMeshChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&testChart, nil)
	moqs.mockedChartPackager.EXPECT().Package(&testChart).Return([]byte("chart_as_bytes"), nil)

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Direct)
	assert.Equal(t, &testChart, generatedChart)
	assert.Nil(t, err)
}

func TestOnboardGeneratedChart_withDirectStrategy_chartGenerationFails(t *testing.T) {
	ctrl, instance, moqs := newTestOnboarderCreator(t)
	defer ctrl.Finish()

	moqs.mockedChartGenerator.EXPECT().GenerateMeshChart(helmManifestResource, "myproject", "mydev", "myservice").Return(nil, errors.New("chart generation failed"))

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Direct)
	assert.NotNil(t, err)
	assert.Nil(t, generatedChart)
}

func TestOnboardGeneratedChart_withDirectStrategy_chartPackagingFails(t *testing.T) {
	ctrl, instance, mocksCol := newTestOnboarderCreator(t)
	defer ctrl.Finish()

	testChart := createChart()

	mocksCol.mockedChartGenerator.EXPECT().GenerateMeshChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&testChart, nil)
	mocksCol.mockedChartPackager.EXPECT().Package(&testChart).Return(nil, errors.New("chart packaging failed"))
	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Direct)

	assert.Nil(t, generatedChart)
	assert.NotNil(t, err)
}

func TestOnboardGeneratedChart_withDuplicateStrategy(t *testing.T) {
	ctrl, instance, moqs := newTestOnboarderCreator(t)
	defer ctrl.Finish()

	testChart := createChart()

	expectedStoreChartOpts := keptnutils.StoreChartOptions{
		Project:   "myproject",
		Service:   "myservice",
		Stage:     "mydev",
		ChartName: "myservice-generated",
		HelmChart: []byte("chart_as_bytes"),
	}
	moqs.mockedChartGenerator.EXPECT().GenerateDuplicateChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&testChart, nil)
	moqs.mockedNamespaceManager.EXPECT().InjectIstio("myproject", "mydev")
	moqs.mockedChartPackager.EXPECT().Package(&testChart).Return([]byte("chart_as_bytes"), nil)
	moqs.mockedChartStorer.EXPECT().Store(gomock.Eq(expectedStoreChartOpts))

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Duplicate)
	assert.Equal(t, &testChart, generatedChart)
	assert.Nil(t, err)
}

func TestOnboardGeneratedChart_injectingIstioConfigFails(t *testing.T) {
	ctrl, instance, moqs := newTestOnboarderCreator(t)
	defer ctrl.Finish()

	testChart := createChart()

	moqs.mockedChartGenerator.EXPECT().GenerateDuplicateChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&testChart, nil)
	moqs.mockedNamespaceManager.EXPECT().InjectIstio("myproject", "mydev").Return(errors.New("failed to inject istio"))

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Duplicate)
	assert.NotNil(t, err)
	assert.Nil(t, generatedChart)
}

func TestOnboardGeneratedChart_chartGenerationFails(t *testing.T) {
	ctrl, instance, moqs := newTestOnboarderCreator(t)
	defer ctrl.Finish()

	moqs.mockedChartGenerator.EXPECT().GenerateDuplicateChart(helmManifestResource, "myproject", "mydev", "myservice").Return(nil, errors.New("chart generation failed"))

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Duplicate)
	assert.NotNil(t, err)
	assert.Nil(t, generatedChart)
}

func TestOnboardGeneratedChart_chartStorageFails(t *testing.T) {
	ctrl, instance, moqs := newTestOnboarderCreator(t)
	defer ctrl.Finish()

	testChart := createChart()

	expectedStoreChartOpts := keptnutils.StoreChartOptions{
		Project:   "myproject",
		Service:   "myservice",
		Stage:     "mydev",
		ChartName: "myservice-generated",
		HelmChart: []byte("chart_as_bytes"),
	}
	moqs.mockedChartGenerator.EXPECT().GenerateDuplicateChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&testChart, nil)
	moqs.mockedNamespaceManager.EXPECT().InjectIstio("myproject", "mydev")
	moqs.mockedChartPackager.EXPECT().Package(&testChart).Return([]byte("chart_as_bytes"), nil)
	moqs.mockedChartStorer.EXPECT().Store(gomock.Eq(expectedStoreChartOpts)).Return("", errors.New("storing chart failed"))

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Duplicate)
	assert.NotNil(t, err)
	assert.Nil(t, generatedChart)
}

func TestCheckAndSetServiceName(t *testing.T) {

	o := onboarder{
		Handler: &HandlerBase{},
		mesh:    nil,
	}
	data := helm.CreateTestHelmChartData(t)

	testCases := []struct {
		name        string
		event       *keptnv2.ServiceCreateFinishedEventData
		error       error
		serviceName string
	}{
		{"Mismatch", &keptnv2.ServiceCreateFinishedEventData{EventData: keptnv2.EventData{Service: "carts-1"},
			Helm: keptnv2.Helm{Chart: base64.StdEncoding.EncodeToString(data)}},
			errors.New("Provided Keptn service name \"carts-1\" does not match Kubernetes service name \"carts\""), "carts-1"},
		{"Mismatch2", &keptnv2.ServiceCreateFinishedEventData{EventData: keptnv2.EventData{Service: "carts-1"},
			Helm: keptnv2.Helm{Chart: "%%%"}},
			errors.New("Error when decoding the Helm Chart: illegal base64 data at input byte 0"), "carts-1"},
		{"Match", &keptnv2.ServiceCreateFinishedEventData{EventData: keptnv2.EventData{Service: ""},
			Helm: keptnv2.Helm{Chart: base64.StdEncoding.EncodeToString(data)}},
			nil, "carts"},
		{"Set", &keptnv2.ServiceCreateFinishedEventData{EventData: keptnv2.EventData{Service: ""},
			Helm: keptnv2.Helm{Chart: base64.StdEncoding.EncodeToString(data)}},
			nil, "carts"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			res := o.checkAndSetServiceName(tt.event)
			if res == nil && res != tt.error {
				t.Errorf("got nil, want %s", tt.error.Error())
			} else if res != nil && tt.error != nil && res.Error() != tt.error.Error() {
				t.Errorf("got %s, want %s", res.Error(), tt.error.Error())
			} else if res != nil && tt.error == nil {
				t.Errorf("got %s, want nil", res.Error())
			}

			if tt.event.Service != tt.serviceName {
				t.Errorf("got %s, want %s", tt.event.Service, tt.serviceName)
			}
		})
	}
}

func stringp(s string) *string {
	return &s
}

func createEventWith(t *testing.T, id string, eventData keptnv2.EventData) event.Event {
	testEvent := cloudevents.NewEvent()
	testEvent.SetID(id)
	_ = testEvent.SetData(cloudevents.ApplicationJSON, keptnv2.ServiceCreateFinishedEventData{
		EventData: eventData,
		Helm: keptnv2.Helm{
			Chart: base64.StdEncoding.EncodeToString(helm.CreateTestHelmChartData(t)),
		},
	})
	return testEvent
}

func createEvent(t *testing.T, id string) event.Event {
	testEvent := cloudevents.NewEvent()
	testEvent.SetID(id)
	_ = testEvent.SetData(cloudevents.ApplicationJSON, keptnv2.ServiceCreateFinishedEventData{
		EventData: keptnv2.EventData{
			Project: "my-project",
			Stage:   "dev",
			Service: "carts",
		},
		Helm: keptnv2.Helm{
			Chart: base64.StdEncoding.EncodeToString(helm.CreateTestHelmChartData(t)),
		},
	})
	return testEvent
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

func createKeptn() *keptnv2.Keptn {
	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})
	return keptn
}
