package api

import (
	"github.com/go-test/deep"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"os"
	"testing"
)

func Test_serviceManager_createService(t *testing.T) {
	projectName := "my-project"
	serviceName := "my-service"

	mockEV := newMockEventbroker(t, func(meb *mockEventBroker, event *models.Event) {
		meb.receivedEvents = append(meb.receivedEvents, *event)
	}, func(meb *mockEventBroker) {

	})

	defer mockEV.server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.server.URL)

	mockCS := newSimpleMockConfigurationService()
	defer mockCS.server.Close()
	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.server.URL)

	mockCS.projects = []*keptnapimodels.Project{
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
			secretStore: &mockSecretStore{
				create: func(name string, content map[string][]byte) error {
					return nil
				},
				delete: func(name string) error {
					return nil
				},
			},
			logger: keptncommon.NewLogger("", "", "shipyard-controller"),
		},
	}

	params := &operations.CreateServiceParams{
		Name: &serviceName,
		Helm: keptnv2.Helm{},
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

	if diff := deep.Equal(expectedProjects, mockCS.projects); len(diff) > 0 {
		t.Errorf("project has not been created correctly")
		for _, d := range diff {
			t.Log(d)
		}
	}

	if shouldContainEvent(t, mockEV.receivedEvents, keptnv2.GetStartedEventType(keptnv2.ServiceCreateTaskName), "", nil) {
		t.Error("event broker did not receive service.create.started event")
	}

	if shouldContainEvent(t, mockEV.receivedEvents, keptnv2.GetFinishedEventType(keptnv2.ServiceCreateTaskName), "", nil) {
		t.Error("event broker did not receive service.create.started event")
	}

	// create the service again - should return error
	err := sm.createService(projectName, params)
	if err != errServiceAlreadyExists {
		t.Errorf("expected errProjectAlreadyExists")
	}
}
