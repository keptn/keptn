package handler

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
)

func Test_validateServiceName(t *testing.T) {
	testcases := []struct {
		name          string
		projectName   string
		stageName     string
		serviceName   string
		expectedError bool
	}{
		{
			name:          "Valid Service name",
			projectName:   "project-1",
			stageName:     "testing",
			serviceName:   "my-service",
			expectedError: false,
		},
		{
			name:          "Invalid Service name",
			projectName:   "project-honk",
			stageName:     "production",
			serviceName:   "my-honk-service-invalid",
			expectedError: true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := validateServiceName(tc.projectName, tc.stageName, tc.serviceName)
			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_serviceManager_createService(t *testing.T) {
	projectName := "my-project"
	serviceName := "my-service"

	mockEV := fake.NewEventBroker(t, func(meb *fake.EventBroker, event *models.Event) {
		meb.ReceivedEvents = append(meb.ReceivedEvents, *event)
	}, func(meb *fake.EventBroker) {

	})

	defer mockEV.Server.Close()
	err := os.Setenv("EVENTBROKER", mockEV.Server.URL)
	require.NoError(t, err)

	mockCS := fake.NewSimpleMockConfigurationService()
	defer mockCS.Server.Close()
	err = os.Setenv("CONFIGURATION_SERVICE", mockCS.Server.URL)
	require.NoError(t, err)

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

	csEndpoint, err := keptncommon.GetServiceEndpoint("CONFIGURATION_SERVICE")
	require.NoError(t, err)

	sm := &serviceManager{
		apiBase: &apiBase{
			projectAPI:  keptnapi.NewProjectHandler(csEndpoint.String()),
			stagesAPI:   keptnapi.NewStageHandler(csEndpoint.String()),
			servicesAPI: keptnapi.NewServiceHandler(csEndpoint.String()),
			resourceAPI: keptnapi.NewResourceHandler(csEndpoint.String()),
			secretStore: &fake.SecretStore{
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

	err = sm.createService(projectName, params)
	require.NoError(t, err)

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

	require.Equal(t, expectedProjects, mockCS.Projects)
	require.False(t, fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetStartedEventType(keptnv2.ServiceCreateTaskName), "", nil))
	require.False(t, fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetFinishedEventType(keptnv2.ServiceCreateTaskName), "", nil))

	// create the service again - should return error
	err = sm.createService(projectName, params)
	require.Error(t, err)
	require.EqualError(t, err, errServiceAlreadyExists.Error())

	// should not create the service, name too long
	serviceName = "name-too-long-honk-1234567890"
	err = sm.createService(projectName, params)
	require.Error(t, err)
}

func Test_serviceManager_deleteService(t *testing.T) {
	projectName := "my-project"
	serviceName := "my-service"
	mockEV := fake.NewEventBroker(t, func(meb *fake.EventBroker, event *models.Event) {
		meb.ReceivedEvents = append(meb.ReceivedEvents, *event)
	}, func(meb *fake.EventBroker) {

	})

	defer mockEV.Server.Close()
	err := os.Setenv("EVENTBROKER", mockEV.Server.URL)
	require.NoError(t, err)

	mockCS := fake.NewSimpleMockConfigurationService()
	defer mockCS.Server.Close()
	err = os.Setenv("CONFIGURATION_SERVICE", mockCS.Server.URL)
	require.NoError(t, err)

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

	csEndpoint, err := keptncommon.GetServiceEndpoint("CONFIGURATION_SERVICE")
	require.NoError(t, err)

	sm := &serviceManager{
		apiBase: &apiBase{
			projectAPI:  keptnapi.NewProjectHandler(csEndpoint.String()),
			stagesAPI:   keptnapi.NewStageHandler(csEndpoint.String()),
			servicesAPI: keptnapi.NewServiceHandler(csEndpoint.String()),
			resourceAPI: keptnapi.NewResourceHandler(csEndpoint.String()),
			secretStore: &fake.SecretStore{
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

	err = sm.deleteService(projectName, serviceName)
	require.NoError(t, err)

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

	require.Equal(t, expectedProjects, mockCS.Projects)
	require.False(t, fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetStartedEventType(keptnv2.ServiceDeleteTaskName), "", nil))
	require.False(t, fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetFinishedEventType(keptnv2.ServiceDeleteTaskName), "", nil))
}
