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
	mockedBaseHandler      *MockHandler
	mockedMesh             *mocks.MockMesh
	mockedProjectHandler   *mocks.MockProjectOperator
	mockedNamespaceManager *mocks.MockINamespaceManager
	mockedStagesHandler    *mocks.MockIStagesHandler
	mockedServiceHandler   *mocks.MockIServiceHandler
	mockedChartStorer      *mocks.MockChartStorer
	mockedChartGenerator   *mocks.MockChartGenerator
	mockedChartPackager    *mocks.MockChartPackager
}

//testOnboarderCreator is an oboarder which has only mocked dependencies
type testOnboarderCreator struct {
}

//Create creates an instance of testOnboarderCreator which full of mocks
func (toc testOnboarderCreator) Create(t *testing.T) (*gomock.Controller, *Onboarder, *mocksCollection) {
	ctrl := gomock.NewController(t)
	mockedBaseHandler := NewMockHandler(ctrl)
	mockedMesh := mocks.NewMockMesh(ctrl)
	mockedProjectHandler := mocks.NewMockProjectOperator(ctrl)
	mockedNamespaceManager := mocks.NewMockINamespaceManager(ctrl)
	mockedStagesHandler := mocks.NewMockIStagesHandler(ctrl)
	mockedServiceHandler := mocks.NewMockIServiceHandler(ctrl)
	mockedChartStorer := mocks.NewMockChartStorer(ctrl)
	mockedChartGenerator := mocks.NewMockChartGenerator(ctrl)
	mockedChartPackager := mocks.NewMockChartPackager(ctrl)

	onboarder := Onboarder{
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
	return ctrl, &onboarder, &mocksCol

}

func TestCreateOnboarder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedMesh := mocks.NewMockMesh(ctrl)
	mockedProjectHandler := mocks.NewMockProjectOperator(ctrl)
	mockedNamespaceManager := mocks.NewMockINamespaceManager(ctrl)
	mockedStagesHandler := mocks.NewMockIStagesHandler(ctrl)
	mockedServiceHandler := mocks.NewMockIServiceHandler(ctrl)
	mockedChartStorer := mocks.NewMockChartStorer(ctrl)
	mockedChartGenerator := mocks.NewMockChartGenerator(ctrl)
	mockedChartPackager := mocks.NewMockChartPackager(ctrl)

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

	toc := testOnboarderCreator{}
	ctrl, instance, mocksCol := toc.Create(t)
	defer ctrl.Finish()

	mocksCol.mockedBaseHandler.EXPECT().handleError("EVENT_ID", gomock.Any(), "service.create", gomock.Any())

	instance.HandleEvent(createUnparsableEvent(), nilCloser)
}

func TestHandleEvent_WhenHelmChartMissing_ThenNothingHappens(t *testing.T) {
	toc := testOnboarderCreator{}
	ctrl, instance, _ := toc.Create(t)
	defer ctrl.Finish()

	instance.HandleEvent(cloudevents.NewEvent(), nilCloser)
}

func TestHandleEvent_WhenNoProjectExists_ThenHandleErrorIsCalled(t *testing.T) {
	toc := testOnboarderCreator{}
	ctrl, instance, mocksCol := toc.Create(t)
	defer ctrl.Finish()

	mocksCol.mockedBaseHandler.EXPECT().getKeptnHandler().Return(nil)
	mocksCol.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(nil, &models.Error{Message: stringp("")})
	mocksCol.mockedBaseHandler.EXPECT().handleError("EVENT_ID", gomock.Any(), "service.create", gomock.Any())

	instance.HandleEvent(createEvent(t, "EVENT_ID"), nilCloser)
}

func TestHandleEvent_WhenNoStagesDefined_ThenHandleErrorIsCalled(t *testing.T) {
	toc := testOnboarderCreator{}
	ctrl, instance, mocksCol := toc.Create(t)
	defer ctrl.Finish()

	mocksCol.mockedBaseHandler.EXPECT().getKeptnHandler().Return(nil)
	mocksCol.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{}, nil)
	mocksCol.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return([]*models.Stage{}, nil)
	mocksCol.mockedBaseHandler.EXPECT().handleError("EVENT_ID", gomock.Any(), "service.create", gomock.Any())

	instance.HandleEvent(createEvent(t, "EVENT_ID"), nilCloser)
}

func TestHandleEvent_WhenStagesCannotBeFetched_ThenHandleErrorIsCalled(t *testing.T) {
	toc := testOnboarderCreator{}
	ctrl, instance, mocksCol := toc.Create(t)
	defer ctrl.Finish()

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}

	mocksCol.mockedBaseHandler.EXPECT().getKeptnHandler().Return(nil)
	mocksCol.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	mocksCol.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(nil, errors.New("unable to fetch stages"))
	mocksCol.mockedBaseHandler.EXPECT().handleError("EVENT_ID", gomock.Any(), "service.create", gomock.Any())

	instance.HandleEvent(createEvent(t, "EVENT_ID"), nilCloser)
}

func TestHandleEvent_WhenInitNamespacesFails_ThenHandleErrorIsCalled(t *testing.T) {
	toc := testOnboarderCreator{}
	ctrl, instance, mocksCol := toc.Create(t)
	defer ctrl.Finish()

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}

	mocksCol.mockedBaseHandler.EXPECT().getKeptnHandler().Return(nil)
	mocksCol.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	mocksCol.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(stages, nil)
	mocksCol.mockedNamespaceManager.EXPECT().InitNamespaces("my-project", []string{"dev"}).Return(errors.New("namespace initialization failed"))
	mocksCol.mockedBaseHandler.EXPECT().handleError("EVENT_ID", gomock.Any(), "service.create", gomock.Any())

	instance.HandleEvent(createEvent(t, "EVENT_ID"), nilCloser)
}

func TestHandleEvent_WhenPassingInvalidServiceName_ThenHandleErrorIsCalled(t *testing.T) {
	toc := testOnboarderCreator{}
	ctrl, instance, mocksCol := toc.Create(t)
	defer ctrl.Finish()

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}

	mocksCol.mockedBaseHandler.EXPECT().getKeptnHandler().Return(nil)
	mocksCol.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	mocksCol.mockedBaseHandler.EXPECT().handleError("EVENT_ID", gomock.Any(), "service.create", gomock.Any())

	instance.HandleEvent(createEventWith(t, "EVENT_ID", keptnv2.EventData{Project: "my-project", Stage: "dev", Service: "!§$%&/"}), nilCloser)
}

func TestHandleEvent_WhenUnableToStoreChart_ThenHandleErrorIsCalled(t *testing.T) {
	toc := testOnboarderCreator{}
	ctrl, instance, mocksCol := toc.Create(t)
	defer ctrl.Finish()

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}

	mocksCol.mockedBaseHandler.EXPECT().getKeptnHandler().Return(nil)
	mocksCol.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	mocksCol.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(stages, nil)
	mocksCol.mockedNamespaceManager.EXPECT().InitNamespaces(gomock.Any(), gomock.Any())
	mocksCol.mockedServiceHandler.EXPECT().GetService(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	mocksCol.mockedChartStorer.EXPECT().StoreChart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New("unable to store chart"))
	mocksCol.mockedBaseHandler.EXPECT().getKeptnHandler().Return(createKeptn())
	mocksCol.mockedBaseHandler.EXPECT().handleError("EVENT_ID", gomock.Any(), "service.create", gomock.Any())

	instance.HandleEvent(createEvent(t, "EVENT_ID"), nilCloser)
}

func TestHandleEvent_WhenSendingFinishedEventFails_ThenHandleErrorisCalled(t *testing.T) {
	toc := testOnboarderCreator{}
	ctrl, instance, mocksCol := toc.Create(t)
	defer ctrl.Finish()

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}

	mocksCol.mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes()
	mocksCol.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	mocksCol.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(stages, nil)
	mocksCol.mockedNamespaceManager.EXPECT().InitNamespaces("my-project", []string{"dev"})
	mocksCol.mockedServiceHandler.EXPECT().GetService("my-project", "dev", "carts").Return(nil, nil)
	mocksCol.mockedChartStorer.EXPECT().StoreChart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil)
	mocksCol.mockedBaseHandler.EXPECT().sendEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to send finished event"))
	mocksCol.mockedBaseHandler.EXPECT().handleError("EVENT_ID", gomock.Any(), "service.create", gomock.Any())

	instance.HandleEvent(createEvent(t, "EVENT_ID"), nilCloser)
}

func TestHandleEvent_OnboardsService(t *testing.T) {
	toc := testOnboarderCreator{}
	ctrl, instance, mocksCol := toc.Create(t)
	defer ctrl.Finish()

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}

	mocksCol.mockedBaseHandler.EXPECT().getKeptnHandler().Return(nil)
	mocksCol.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	mocksCol.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(stages, nil)
	mocksCol.mockedNamespaceManager.EXPECT().InitNamespaces("my-project", []string{"dev"})
	mocksCol.mockedServiceHandler.EXPECT().GetService("my-project", "dev", "carts").Return(nil, nil)
	mocksCol.mockedChartStorer.EXPECT().StoreChart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil)
	mocksCol.mockedBaseHandler.EXPECT().sendEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	instance.HandleEvent(createEvent(t, "EVENT_ID"), nilCloser)
}

func TestOnboardGeneratedChart_withDirectStrategy(t *testing.T) {
	toc := testOnboarderCreator{}
	ctrl, instance, mocksCol := toc.Create(t)
	defer ctrl.Finish()

	testChart := createChart()

	expectedChartAsBytes := []byte("chart_as_bytes")
	mocksCol.mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(createKeptn())
	mocksCol.mockedChartStorer.EXPECT().StoreChart("myproject", "myservice", "mydev", "myservice-generated", gomock.Eq(expectedChartAsBytes))
	mocksCol.mockedChartGenerator.EXPECT().GenerateMeshChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&testChart, nil)
	mocksCol.mockedChartPackager.EXPECT().PackageChart(&testChart).Return(expectedChartAsBytes, nil)

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Direct)
	assert.Equal(t, &testChart, generatedChart)
	assert.Nil(t, err)
}

func TestOnboardGeneratedChart_withDirectStrategy_chartGenerationFails(t *testing.T) {
	toc := testOnboarderCreator{}
	ctrl, instance, mocksCol := toc.Create(t)
	defer ctrl.Finish()

	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})

	mocksCol.mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(keptn)
	mocksCol.mockedChartGenerator.EXPECT().GenerateMeshChart(helmManifestResource, "myproject", "mydev", "myservice").Return(nil, errors.New("chart generation failed"))

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Direct)
	assert.NotNil(t, err)
	assert.Nil(t, generatedChart)
}

func TestOnboardGeneratedChart_withDirectStrategy_chartPackagingFails(t *testing.T) {
	toc := testOnboarderCreator{}
	ctrl, instance, mocksCol := toc.Create(t)
	defer ctrl.Finish()

	testChart := createChart()

	mocksCol.mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(createKeptn())
	mocksCol.mockedChartGenerator.EXPECT().GenerateMeshChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&testChart, nil)
	mocksCol.mockedChartPackager.EXPECT().PackageChart(&testChart).Return(nil, errors.New("chart packaging failed"))
	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Direct)

	assert.Nil(t, generatedChart)
	assert.NotNil(t, err)
}

func TestOnboardGeneratedChart_withDuplicateStrategy(t *testing.T) {
	toc := testOnboarderCreator{}
	ctrl, instance, mocksCol := toc.Create(t)
	defer ctrl.Finish()

	testChart := createChart()

	expectedChartAsBytes := []byte("chart_as_bytes")
	mocksCol.mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(createKeptn())
	mocksCol.mockedChartGenerator.EXPECT().GenerateDuplicateChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&testChart, nil)
	mocksCol.mockedNamespaceManager.EXPECT().InjectIstio("myproject", "mydev")
	mocksCol.mockedChartPackager.EXPECT().PackageChart(&testChart).Return(expectedChartAsBytes, nil)
	mocksCol.mockedChartStorer.EXPECT().StoreChart("myproject", "myservice", "mydev", "myservice-generated", gomock.Eq(expectedChartAsBytes))

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Duplicate)
	assert.Equal(t, &testChart, generatedChart)
	assert.Nil(t, err)
}

func TestOnboardGeneratedChart_injectingIstioConfigFails(t *testing.T) {
	toc := testOnboarderCreator{}
	ctrl, instance, mocksCol := toc.Create(t)
	defer ctrl.Finish()

	testChart := createChart()

	mocksCol.mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(createKeptn())
	mocksCol.mockedChartGenerator.EXPECT().GenerateDuplicateChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&testChart, nil)
	mocksCol.mockedNamespaceManager.EXPECT().InjectIstio("myproject", "mydev").Return(errors.New("failed to inject istio"))

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Duplicate)
	assert.NotNil(t, err)
	assert.Nil(t, generatedChart)
}

func TestOnboardGeneratedChart_chartGenerationFails(t *testing.T) {
	toc := testOnboarderCreator{}
	ctrl, instance, mocksCol := toc.Create(t)
	defer ctrl.Finish()

	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})

	mocksCol.mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(keptn)
	mocksCol.mockedChartGenerator.EXPECT().GenerateDuplicateChart(helmManifestResource, "myproject", "mydev", "myservice").Return(nil, errors.New("chart generation failed"))

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Duplicate)
	assert.NotNil(t, err)
	assert.Nil(t, generatedChart)
}

func TestOnboardGeneratedChart_chartStorageFails(t *testing.T) {
	toc := testOnboarderCreator{}
	ctrl, instance, mocksCol := toc.Create(t)
	defer ctrl.Finish()

	testChart := createChart()

	expectedChartAsBytes := []byte("chart_as_bytes")
	mocksCol.mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(createKeptn())
	mocksCol.mockedChartGenerator.EXPECT().GenerateDuplicateChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&testChart, nil)
	mocksCol.mockedNamespaceManager.EXPECT().InjectIstio("myproject", "mydev")
	mocksCol.mockedChartPackager.EXPECT().PackageChart(&testChart).Return(expectedChartAsBytes, nil)
	mocksCol.mockedChartStorer.EXPECT().StoreChart("myproject", "myservice", "mydev", "myservice-generated", gomock.Eq(expectedChartAsBytes)).Return("", errors.New("storing chart failed"))

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}, keptnevents.Duplicate)
	assert.NotNil(t, err)
	assert.Nil(t, generatedChart)
}

func TestCheckAndSetServiceName(t *testing.T) {

	o := Onboarder{
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

func nilCloser(keptnHandler *keptnv2.Keptn) {
	//No-op
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
