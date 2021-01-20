package controller

import (
	"encoding/base64"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/golang/mock/gomock"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/helm-service/mocks"
	"github.com/stretchr/testify/assert"
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
	mockedOnBoarder        *mocks.MockOnboarder
}

// testOnboarderCreator is an oboarder which has only mocked dependencies
type testOnboarderCreator struct {
}

// newTestOnboardHandlerCreator creates an instance of testOnboarderCreator which uses only mocks
func newTestOnboardHandlerCreator(t *testing.T, mockedBaseHandlerOptions ...MockedHandlerOption) (*gomock.Controller, OnboardHandler, mocksCollection) {

	ctrl := gomock.NewController(t)
	mockedBaseHandler := NewMockedHandler(createKeptn(), "", mockedBaseHandlerOptions...)
	mockedMesh := mocks.NewMockMesh(ctrl)
	mockedProjectHandler := mocks.NewMockIProjectHandler(ctrl)
	mockedNamespaceManager := mocks.NewMockINamespaceManager(ctrl)
	mockedStagesHandler := mocks.NewMockIStagesHandler(ctrl)
	mockedServiceHandler := mocks.NewMockIServiceHandler(ctrl)
	mockedOnBoarder := mocks.NewMockOnboarder(ctrl)

	onboardHandler := OnboardHandler{
		Handler:          mockedBaseHandler,
		projectHandler:   mockedProjectHandler,
		namespaceManager: mockedNamespaceManager,
		stagesHandler:    mockedStagesHandler,
		onboarder:        mockedOnBoarder,
	}

	mocksCol := mocksCollection{
		mockedBaseHandler:      mockedBaseHandler,
		mockedMesh:             mockedMesh,
		mockedProjectHandler:   mockedProjectHandler,
		mockedNamespaceManager: mockedNamespaceManager,
		mockedStagesHandler:    mockedStagesHandler,
		mockedServiceHandler:   mockedServiceHandler,
		mockedOnBoarder:        mockedOnBoarder,
	}
	return ctrl, onboardHandler, mocksCol

}

func TestCreateOnboarderHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedProjectHandler := mocks.NewMockIProjectHandler(ctrl)
	mockedNamespaceManager := mocks.NewMockINamespaceManager(ctrl)
	mockedStagesHandler := mocks.NewMockIStagesHandler(ctrl)
	mockedOnboarder := mocks.NewMockOnboarder(ctrl)

	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})
	handler := NewOnboardHandler(
		keptn,
		mockedProjectHandler,
		mockedNamespaceManager,
		mockedStagesHandler,
		mockedOnboarder,
		"")

	assert.NotNil(t, handler)

}

func TestHandleEvent_WhenPassingUnparsableEvent_ThenHandleErrorIsCalled(t *testing.T) {

	ctrl, instance, _ := newTestOnboardHandlerCreator(t)
	defer ctrl.Finish()
	instance.HandleEvent(createUnparsableEvent())
}

func TestHandleEvent_WhenHelmChartMissing_ThenNothingHappens(t *testing.T) {
	ctrl, instance, moqs := newTestOnboardHandlerCreator(t)
	defer ctrl.Finish()

	instance.HandleEvent(cloudevents.NewEvent())
	assert.Empty(t, moqs.mockedBaseHandler.sentCloudEvents)
	assert.Empty(t, moqs.mockedBaseHandler.handledErrorEvents)
}

func TestHandleEvent_WhenNoProjectExists_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl, instance, moqs := newTestOnboardHandlerCreator(t)
	defer ctrl.Finish()

	moqs.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(nil, &models.Error{Message: stringp("")})

	instance.HandleEvent(createEvent(t, "EVENT_ID"))
	assert.Equal(t, 1, len(moqs.mockedBaseHandler.handledErrorEvents))
}

func TestHandleEvent_WhenNoStagesDefined_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl, instance, moqs := newTestOnboardHandlerCreator(t)
	defer ctrl.Finish()

	moqs.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{}, nil)
	moqs.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return([]*models.Stage{}, nil)

	instance.HandleEvent(createEvent(t, "EVENT_ID"))
	assert.Equal(t, 1, len(moqs.mockedBaseHandler.handledErrorEvents))
}

func TestHandleEvent_WhenStagesCannotBeFetched_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl, instance, moqs := newTestOnboardHandlerCreator(t)
	defer ctrl.Finish()

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}

	moqs.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	moqs.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(nil, errors.New("unable to fetch stages"))

	instance.HandleEvent(createEvent(t, "EVENT_ID"))
	assert.Equal(t, 1, len(moqs.mockedBaseHandler.handledErrorEvents))
}

func TestHandleEvent_WhenInitNamespacesFails_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl, instance, moqs := newTestOnboardHandlerCreator(t)
	defer ctrl.Finish()

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}

	moqs.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	moqs.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(stages, nil)
	moqs.mockedNamespaceManager.EXPECT().InitNamespaces("my-project", []string{"dev"}).Return(errors.New("namespace initialization failed"))

	instance.HandleEvent(createEvent(t, "EVENT_ID"))
	assert.Equal(t, 1, len(moqs.mockedBaseHandler.handledErrorEvents))
}

func TestHandleEvent_WhenPassingInvalidServiceName_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl, instance, moqs := newTestOnboardHandlerCreator(t)
	defer ctrl.Finish()

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}

	moqs.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)

	instance.HandleEvent(createEventWith(t, "EVENT_ID", keptnv2.EventData{Project: "my-project", Stage: "dev", Service: "otherthancarts"}))
	assert.Equal(t, 1, len(moqs.mockedBaseHandler.handledErrorEvents))
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
