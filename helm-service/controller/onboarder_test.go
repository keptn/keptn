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
	configmodels "github.com/keptn/go-utils/pkg/api/models"
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

//MOCKS
var mockedBaseHandler *MockHandler
var mockedMesh *mocks.MockMesh
var mockedProjectHandler *mocks.MockProjectOperator
var mockedNamespaceManager *mocks.MockINamespaceManager
var mockedStagesHandler *mocks.MockIStagesHandler
var mockedServiceHandler *mocks.MockIServiceHandler
var mockedChartStorer *mocks.MockChartStorer
var mockedChartGenerator *mocks.MockChartGenerator
var mockedChartPackager *mocks.MockChartPackager

func createMocks(t *testing.T) *gomock.Controller {

	ctrl := gomock.NewController(t)
	mockedBaseHandler = NewMockHandler(ctrl)
	mockedMesh = mocks.NewMockMesh(ctrl)
	mockedProjectHandler = mocks.NewMockProjectOperator(ctrl)
	mockedNamespaceManager = mocks.NewMockINamespaceManager(ctrl)
	mockedStagesHandler = mocks.NewMockIStagesHandler(ctrl)
	mockedServiceHandler = mocks.NewMockIServiceHandler(ctrl)
	mockedChartStorer = mocks.NewMockChartStorer(ctrl)
	mockedChartGenerator = mocks.NewMockChartGenerator(ctrl)
	mockedChartPackager = mocks.NewMockChartPackager(ctrl)
	return ctrl

}

func TestCreateOnboarder(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})
	if onboarder := NewOnboarder(
		keptn,
		mockedMesh,
		mockedProjectHandler,
		mockedNamespaceManager,
		mockedStagesHandler,
		mockedServiceHandler,
		mockedChartStorer,
		mockedChartGenerator,
		mockedChartPackager,
		""); onboarder == nil {

		t.Error("onboarder instance is nil")
	}
}

func TestHandleEvent_WhenPassingUnparsableEvent_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	mockedBaseHandler.EXPECT().handleError("EVENT_ID", gomock.Any(), "service.create", gomock.Any())

	instance := Onboarder{
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

	instance.HandleEvent(createUnparsableEvent(), nilCloser)
}

func TestHandleEvent_WhenHelmChartMissing_ThenNothingHappens(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	instance := Onboarder{
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

	event := cloudevents.NewEvent()
	instance.HandleEvent(event, nilCloser)
}

func TestHandleEvent_WhenNoProjectExists_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	mockedBaseHandler.EXPECT().getKeptnHandler().Return(nil)
	mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(nil, &models.Error{Message: stringp("")})
	mockedBaseHandler.EXPECT().handleError("EVENT_ID", gomock.Any(), "service.create", gomock.Any())

	instance := Onboarder{
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

	instance.HandleEvent(createEvent(t, "EVENT_ID"), nilCloser)
}

func TestHandleEvent_WhenNoStagesDefined_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	mockedBaseHandler.EXPECT().getKeptnHandler().Return(nil)
	mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{}, nil)
	mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return([]*models.Stage{}, nil)
	mockedBaseHandler.EXPECT().handleError("EVENT_ID", gomock.Any(), "service.create", gomock.Any())

	instance := Onboarder{
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

	instance.HandleEvent(createEvent(t, "EVENT_ID"), nilCloser)
}

func TestHandleEvent_WhenInitNamespacesFails_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	stages := []*models.Stage{&models.Stage{Services: []*models.Service{&models.Service{}}, StageName: "dev"}}

	mockedBaseHandler.EXPECT().getKeptnHandler().Return(nil)
	mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(stages, nil)
	mockedNamespaceManager.EXPECT().InitNamespaces("my-project", []string{"dev"}).Return(errors.New("Namespace initialization failed :("))
	mockedBaseHandler.EXPECT().handleError("EVENT_ID", gomock.Any(), "service.create", gomock.Any())

	instance := Onboarder{
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

	instance.HandleEvent(createEvent(t, "EVENT_ID"), nilCloser)
}

func TestHandleEvent_WhenPassingInvalidServiceName_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	stages := []*models.Stage{&models.Stage{Services: []*models.Service{&models.Service{}}, StageName: "dev"}}

	mockedBaseHandler.EXPECT().getKeptnHandler().Return(nil)
	mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	mockedBaseHandler.EXPECT().handleError("EVENT_ID", gomock.Any(), "service.create", gomock.Any())

	instance := Onboarder{
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

	instance.HandleEvent(createEventWith(t, "EVENT_ID", keptnv2.EventData{Project: "my-project", Stage: "dev", Service: "!§$%&/"}), nilCloser)
}

func TestHandleEvent_WhenUnableToStoreChart_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	stages := []*models.Stage{&models.Stage{Services: []*models.Service{&models.Service{}}, StageName: "dev"}}

	mockedBaseHandler.EXPECT().getKeptnHandler().Return(nil)
	mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(stages, nil)
	mockedNamespaceManager.EXPECT().InitNamespaces(gomock.Any(), gomock.Any())
	mockedServiceHandler.EXPECT().GetService(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	mockedChartStorer.EXPECT().StoreChart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New("Unable to store chart :("))
	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})
	mockedBaseHandler.EXPECT().getKeptnHandler().Return(keptn)
	mockedBaseHandler.EXPECT().handleError("EVENT_ID", gomock.Any(), "service.create", gomock.Any())

	instance := Onboarder{
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

	instance.HandleEvent(createEvent(t, "EVENT_ID"), nilCloser)
}

func TestHandleEvent_WhenSendingFinishedEventFails_ThenHandleErrorisCalled(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	stages := []*models.Stage{&models.Stage{Services: []*models.Service{&models.Service{}}, StageName: "dev"}}

	mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes()
	mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(stages, nil)
	mockedNamespaceManager.EXPECT().InitNamespaces("my-project", []string{"dev"})
	mockedServiceHandler.EXPECT().GetService("my-project", "dev", "carts").Return(nil, nil)
	mockedChartStorer.EXPECT().StoreChart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil)
	mockedBaseHandler.EXPECT().sendEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Failed to send finished event :("))
	mockedBaseHandler.EXPECT().handleError("EVENT_ID", gomock.Any(), "service.create", gomock.Any())

	instance := Onboarder{
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

	instance.HandleEvent(createEvent(t, "EVENT_ID"), nilCloser)

}

func TestHandleEvent_OnboardsService(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	stages := []*models.Stage{&models.Stage{Services: []*models.Service{&models.Service{}}, StageName: "dev"}}

	mockedBaseHandler.EXPECT().getKeptnHandler().Return(nil)
	mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(stages, nil)
	mockedNamespaceManager.EXPECT().InitNamespaces("my-project", []string{"dev"})
	mockedServiceHandler.EXPECT().GetService("my-project", "dev", "carts").Return(nil, nil)
	mockedChartStorer.EXPECT().StoreChart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil)
	mockedBaseHandler.EXPECT().sendEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	instance := Onboarder{
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

	instance.HandleEvent(createEvent(t, "EVENT_ID"), nilCloser)

}

func TestOnboardGeneratedChart_withDirectStrategy(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})

	chart := chart.Chart{
		Metadata: &chart.Metadata{
			APIVersion: "v2",
			Name:       "myservice-generated",
			Keywords:   []string{"deployment_strategy=" + keptnevents.Direct.String()},
			Version:    "0.1.0",
		},
	}

	expectedChartAsBytes := []byte("chart_as_bytes")
	mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(keptn)
	mockedChartStorer.EXPECT().StoreChart("myproject", "myservice", "mydev", "myservice-generated", gomock.Eq(expectedChartAsBytes))
	mockedChartGenerator.EXPECT().GenerateMeshChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&chart, nil)
	mockedChartPackager.EXPECT().PackageChart(&chart).Return(expectedChartAsBytes, nil)

	instance := Onboarder{
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

	eventData := keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, eventData, keptnevents.Direct)
	assert.Equal(t, &chart, generatedChart)
	assert.Nil(t, err)
}

func TestOnboardGeneratedChart_withDirectStrategy_chartGenerationFails(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})

	mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(keptn)
	mockedChartGenerator.EXPECT().GenerateMeshChart(helmManifestResource, "myproject", "mydev", "myservice").Return(nil, errors.New("Chart generation failed"))

	instance := Onboarder{
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

	eventData := keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, eventData, keptnevents.Direct)

	assert.NotNil(t, err)
	assert.Nil(t, generatedChart)
}

func TestOnboardGeneratedChart_withDirectStrategy_chartPackagingFails(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})

	chart := chart.Chart{
		Metadata: &chart.Metadata{
			APIVersion: "v2",
			Name:       "myservice-generated",
			Keywords:   []string{"deployment_strategy=" + keptnevents.Direct.String()},
			Version:    "0.1.0",
		},
	}

	mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(keptn)
	mockedChartGenerator.EXPECT().GenerateMeshChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&chart, nil)
	mockedChartPackager.EXPECT().PackageChart(&chart).Return(nil, errors.New("Chart packaging failed :("))

	instance := Onboarder{
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

	eventData := keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, eventData, keptnevents.Direct)

	assert.Nil(t, generatedChart)
	assert.NotNil(t, err)
}

func TestOnboardGeneratedChart_withDuplicateStrategy(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})

	chart := chart.Chart{
		Metadata: &chart.Metadata{
			APIVersion: "v2",
			Name:       "myservice-generated",
			Keywords:   []string{"deployment_strategy=" + keptnevents.Duplicate.String()},
			Version:    "0.1.0",
		},
	}

	expectedChartAsBytes := []byte("chart_as_bytes")
	mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(keptn)
	mockedChartGenerator.EXPECT().GenerateDuplicateChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&chart, nil)
	mockedNamespaceManager.EXPECT().InjectIstio("myproject", "mydev")
	mockedChartPackager.EXPECT().PackageChart(&chart).Return(expectedChartAsBytes, nil)
	mockedChartStorer.EXPECT().StoreChart("myproject", "myservice", "mydev", "myservice-generated", gomock.Eq(expectedChartAsBytes))

	instance := Onboarder{
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

	eventData := keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, eventData, keptnevents.Duplicate)
	assert.Equal(t, &chart, generatedChart)
	assert.Nil(t, err)
}

func TestOnboardGeneratedChart_injectingIstioConfigFails(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})

	chart := chart.Chart{
		Metadata: &chart.Metadata{
			APIVersion: "v2",
			Name:       "myservice-generated",
			Keywords:   []string{"deployment_strategy=" + keptnevents.Duplicate.String()},
			Version:    "0.1.0",
		},
	}

	mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(keptn)
	mockedChartGenerator.EXPECT().GenerateDuplicateChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&chart, nil)
	mockedNamespaceManager.EXPECT().InjectIstio("myproject", "mydev").Return(errors.New("Failed to inject istio :("))

	instance := Onboarder{
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

	eventData := keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, eventData, keptnevents.Duplicate)
	assert.NotNil(t, err)
	assert.Nil(t, generatedChart)
}

func TestOnboardGeneratedChart_chartGenerationFails(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})

	mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(keptn)
	mockedChartGenerator.EXPECT().GenerateDuplicateChart(helmManifestResource, "myproject", "mydev", "myservice").Return(nil, errors.New("Chart generation failed"))

	instance := Onboarder{
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

	eventData := keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, eventData, keptnevents.Duplicate)

	assert.NotNil(t, err)
	assert.Nil(t, generatedChart)
}

func TestOnboardGeneratedChart_chartStorageFails(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})

	chart := chart.Chart{
		Metadata: &chart.Metadata{
			APIVersion: "v2",
			Name:       "myservice-generated",
			Keywords:   []string{"deployment_strategy=" + keptnevents.Duplicate.String()},
			Version:    "0.1.0",
		},
	}

	expectedChartAsBytes := []byte("chart_as_bytes")
	mockedBaseHandler.EXPECT().getKeptnHandler().AnyTimes().Return(keptn)
	mockedChartGenerator.EXPECT().GenerateDuplicateChart(helmManifestResource, "myproject", "mydev", "myservice").Return(&chart, nil)
	mockedNamespaceManager.EXPECT().InjectIstio("myproject", "mydev")
	mockedChartPackager.EXPECT().PackageChart(&chart).Return(expectedChartAsBytes, nil)
	mockedChartStorer.EXPECT().StoreChart("myproject", "myservice", "mydev", "myservice-generated", gomock.Eq(expectedChartAsBytes)).Return("", errors.New("Storing chart failed :("))

	instance := Onboarder{
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

	eventData := keptnv2.EventData{Project: "myproject", Stage: "mydev", Service: "myservice"}

	generatedChart, err := instance.OnboardGeneratedChart(helmManifestResource, eventData, keptnevents.Duplicate)
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
		{"Match", &keptnv2.ServiceCreateFinishedEventData{EventData: keptnv2.EventData{Service: "carts"},
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

func check(e *configmodels.Error, t *testing.T) {
	if e != nil {
		t.Error(e.Message)
	}
}

func createEventWith(t *testing.T, id string, eventData keptnv2.EventData) event.Event {
	event := cloudevents.NewEvent()
	event.SetID(id)
	event.SetData(cloudevents.ApplicationJSON, keptnv2.ServiceCreateFinishedEventData{
		EventData: eventData,
		Helm: keptnv2.Helm{
			Chart: base64.StdEncoding.EncodeToString(helm.CreateTestHelmChartData(t)),
		},
	})
	return event
}

func createEvent(t *testing.T, id string) event.Event {
	event := cloudevents.NewEvent()
	event.SetID(id)
	event.SetData(cloudevents.ApplicationJSON, keptnv2.ServiceCreateFinishedEventData{
		EventData: keptnv2.EventData{
			Project: "my-project",
			Stage:   "dev",
			Service: "carts",
		},
		Helm: keptnv2.Helm{
			Chart: base64.StdEncoding.EncodeToString(helm.CreateTestHelmChartData(t)),
		},
	})
	return event
}

func createUnparsableEvent() event.Event {
	event := cloudevents.NewEvent()
	event.SetData(cloudevents.ApplicationJSON, "WEIRD_JSON_CONTENT")
	event.SetID("EVENT_ID")
	return event
}

func nilCloser(keptnHandler *keptnv2.Keptn) {
	//No-op
}
