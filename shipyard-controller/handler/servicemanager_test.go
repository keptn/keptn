package handler

import (
	"github.com/go-test/deep"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"os"
	"testing"
)

func Test_serviceManager_createService(t *testing.T) {
	projectName := "my-project"
	serviceName := "my-service"

	mockEV := fake.NewMockEventbroker(t, func(meb *fake.MockEventBroker, event *models.Event) {
		meb.ReceivedEvents = append(meb.ReceivedEvents, *event)
	}, func(meb *fake.MockEventBroker) {

	})

	defer mockEV.Server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.Server.URL)

	mockCS := fake.NewSimpleMockConfigurationService()
	defer mockCS.Server.Close()
	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.Server.URL)

	mockCS.Projects = []*keptnapimodels.Project{
		{
			ProjectName: projectName,
			Stages: []*keptnapimodels.Stage{
				{
					StageName: "dev",
				},
				{
					StageName: "hardening",
				},
				{
					StageName: "production",
				},
			},
		},
	}

	csEndpoint, _ := keptncommon.GetServiceEndpoint("CONFIGURATION_SERVICE")

	sm := &serviceManager{
		apiBase: &apiBase{
			projectAPI:  keptnapi.NewProjectHandler(csEndpoint.String()),
			stagesAPI:   keptnapi.NewStageHandler(csEndpoint.String()),
			servicesAPI: keptnapi.NewServiceHandler(csEndpoint.String()),
			resourceAPI: keptnapi.NewResourceHandler(csEndpoint.String()),
			secretStore: &fake.MockSecretStore{
				CreateFunc: func(name string, content map[string][]byte) error {
					return nil
				},
				DeleteFunc: func(name string) error {
					return nil
				},
			},
			logger: keptncommon.NewLogger("", "", "shipyard-controller"),
		},
	}

	params := &operations.CreateServiceParams{
		ServiceName: &serviceName,
		HelmChart:   "",
	}

	if err := sm.createService(projectName, params); err != nil {
		t.Error("received error: " + err.Error())
	}

	expectedProjects := []*keptnapimodels.Project{
		{
			ProjectName: projectName,
			Stages: []*keptnapimodels.Stage{
				{
					StageName: "dev",
					Services: []*keptnapimodels.Service{
						{
							ServiceName: serviceName,
						},
					},
				},
				{
					StageName: "hardening",
					Services: []*keptnapimodels.Service{
						{
							ServiceName: serviceName,
						},
					},
				},
				{
					StageName: "production",
					Services: []*keptnapimodels.Service{
						{
							ServiceName: serviceName,
						},
					},
				},
			},
		},
	}

	if diff := deep.Equal(expectedProjects, mockCS.Projects); len(diff) > 0 {
		t.Errorf("project has not been created correctly")
		for _, d := range diff {
			t.Log(d)
		}
	}

	if fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetStartedEventType(keptnv2.ServiceCreateTaskName), "", nil) {
		t.Error("event broker did not receive service.create.started event")
	}

	if fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetFinishedEventType(keptnv2.ServiceCreateTaskName), "", nil) {
		t.Error("event broker did not receive service.create.started event")
	}

	// create the service again - should return error
	err := sm.createService(projectName, params)
	if err != errServiceAlreadyExists {
		t.Errorf("expected errProjectAlreadyExists")
	}
}

func Test_serviceManager_deleteService(t *testing.T) {
	projectName := "my-project"
	serviceName := "my-service"
	mockEV := fake.NewMockEventbroker(t, func(meb *fake.MockEventBroker, event *models.Event) {
		meb.ReceivedEvents = append(meb.ReceivedEvents, *event)
	}, func(meb *fake.MockEventBroker) {

	})

	defer mockEV.Server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.Server.URL)

	mockCS := fake.NewSimpleMockConfigurationService()
	defer mockCS.Server.Close()
	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.Server.URL)

	mockCS.Projects = []*keptnapimodels.Project{
		{
			ProjectName: projectName,
			Stages: []*keptnapimodels.Stage{
				{
					StageName: "dev",
					Services: []*keptnapimodels.Service{
						{
							ServiceName: serviceName,
						},
					},
				},
				{
					StageName: "hardening",
					Services: []*keptnapimodels.Service{
						{
							ServiceName: serviceName,
						},
					},
				},
				{
					StageName: "production",
					Services: []*keptnapimodels.Service{
						{
							ServiceName: serviceName,
						},
					},
				},
			},
		},
	}

	csEndpoint, _ := keptncommon.GetServiceEndpoint("CONFIGURATION_SERVICE")

	sm := &serviceManager{
		apiBase: &apiBase{
			projectAPI:  keptnapi.NewProjectHandler(csEndpoint.String()),
			stagesAPI:   keptnapi.NewStageHandler(csEndpoint.String()),
			servicesAPI: keptnapi.NewServiceHandler(csEndpoint.String()),
			resourceAPI: keptnapi.NewResourceHandler(csEndpoint.String()),
			secretStore: &fake.MockSecretStore{
				CreateFunc: func(name string, content map[string][]byte) error {
					return nil
				},
				DeleteFunc: func(name string) error {
					return nil
				},
			},
			logger: keptncommon.NewLogger("", "", "shipyard-controller"),
		},
	}

	if err := sm.deleteService(projectName, serviceName); err != nil {
		t.Error("received error: " + err.Error())
	}

	expectedProjects := []*keptnapimodels.Project{
		{
			ProjectName: "my-project",
			Stages: []*keptnapimodels.Stage{
				{
					StageName: "dev",
					Services:  []*keptnapimodels.Service{},
				},
				{
					StageName: "hardening",
					Services:  []*keptnapimodels.Service{},
				},
				{
					StageName: "production",
					Services:  []*keptnapimodels.Service{},
				},
			},
		},
	}

	if diff := deep.Equal(expectedProjects, mockCS.Projects); len(diff) > 0 {
		t.Errorf("project has not been created correctly")
		for _, d := range diff {
			t.Log(d)
		}
	}

	if fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetStartedEventType(keptnv2.ServiceDeleteTaskName), "", nil) {
		t.Error("event broker did not receive service.delete.started event")
	}

	if fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetFinishedEventType(keptnv2.ServiceDeleteTaskName), "", nil) {
		t.Error("event broker did not receive service.delete.finished event")
	}

}
