package handler

import (
	"errors"
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

const testBase64EncodedShipyardContent = `YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiBzaGlweWFyZC1zb2Nrc2hvcApzcGVjOgogIHN0YWdlczoKICAtIG5hbWU6IGRldgogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0YXNrczoKICAgICAgLSBuYW1lOiBkZXBsb3ltZW50CiAgICAgICAgcHJvcGVydGllczogIAogICAgICAgICAgc3RyYXRlZ3k6IGRpcmVjdAogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOgogICAgICAgICAga2luZDogZnVuY3Rpb25hbAogICAgICAtIG5hbWU6IGV2YWx1YXRpb24gCiAgICAgIC0gbmFtZTogcmVsZWFzZSAKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg==`
const testBase64EncodedShipyardContentWithInvalidVersion = `YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjEuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiBzaGlweWFyZC1zb2Nrc2hvcApzcGVjOgogIHN0YWdlczoKICAtIG5hbWU6IGRldgogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0YXNrczoKICAgICAgLSBuYW1lOiBkZXBsb3ltZW50CiAgICAgICAgcHJvcGVydGllczogIAogICAgICAgICAgc3RyYXRlZ3k6IGRpcmVjdAogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOgogICAgICAgICAga2luZDogZnVuY3Rpb25hbAogICAgICAtIG5hbWU6IGV2YWx1YXRpb24gCiAgICAgIC0gbmFtZTogcmVsZWFzZQ==`
const testBase64EncodedShipyardContentWithInvalidStage = `YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjEuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiBzaGlweWFyZC1zb2Nrc2hvcApzcGVjOgogIHN0YWdlczoKICAtIHNlcXVlbmNlczoKICAgIC0gbmFtZTogYXJ0aWZhY3QtZGVsaXZlcnkKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6ICAKICAgICAgICAgIHN0cmF0ZWd5OiBkaXJlY3QKICAgICAgLSBuYW1lOiB0ZXN0CiAgICAgICAgcHJvcGVydGllczoKICAgICAgICAgIGtpbmQ6IGZ1bmN0aW9uYWwKICAgICAgLSBuYW1lOiBldmFsdWF0aW9uIAogICAgICAtIG5hbWU6IHJlbGVhc2U=`

func Test_validateUpdateProjectParams(t *testing.T) {
	type args struct {
		createProjectParams *operations.CreateProjectParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should contain valid Keptn entity name",
			args: args{
				createProjectParams: &operations.CreateProjectParams{
					Name:     stringp("my-project"),
					Shipyard: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "should contain valid Keptn entity name",
			args: args{
				createProjectParams: &operations.CreateProjectParams{
					Name:     stringp("my@project"),
					Shipyard: nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateUpdateProjectParams(tt.args.createProjectParams); (err != nil) != tt.wantErr {
				t.Errorf("validateUpdateProjectParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateCreateProjectParams(t *testing.T) {
	type args struct {
		createProjectParams *operations.CreateProjectParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid entity name",
			args: args{
				createProjectParams: &operations.CreateProjectParams{
					Name:     stringp("my@project"),
					Shipyard: nil,
				},
			},
			wantErr: true,
		},
		{
			name: "no shipyard provided",
			args: args{
				createProjectParams: &operations.CreateProjectParams{
					Name:     stringp("my-project"),
					Shipyard: nil,
				},
			},
			wantErr: true,
		},
		{
			name: "no shipyard provided 2",
			args: args{
				createProjectParams: &operations.CreateProjectParams{
					Name:     stringp("my-project"),
					Shipyard: stringp(""),
				},
			},
			wantErr: true,
		},
		{
			name: "no base64 encoded shipyard",
			args: args{
				createProjectParams: &operations.CreateProjectParams{
					Name:     stringp("my-project"),
					Shipyard: stringp(fake.TestShipyardFile),
				},
			},
			wantErr: true,
		},
		{
			name: "base64 encoded shipyard with invalid version",
			args: args{
				createProjectParams: &operations.CreateProjectParams{
					Name:     stringp("my-project"),
					Shipyard: stringp(testBase64EncodedShipyardContentWithInvalidVersion),
				},
			},
			wantErr: true,
		},
		{
			name: "base64 encoded shipyard with invalid stage",
			args: args{
				createProjectParams: &operations.CreateProjectParams{
					Name:     stringp("my-project"),
					Shipyard: stringp(testBase64EncodedShipyardContentWithInvalidStage),
				},
			},
			wantErr: true,
		},
		{
			name: "base64 encoded shipyard with valid version",
			args: args{
				createProjectParams: &operations.CreateProjectParams{
					Name:     stringp("my-project"),
					Shipyard: stringp(testBase64EncodedShipyardContent),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateCreateProjectParams(tt.args.createProjectParams); (err != nil) != tt.wantErr {
				t.Errorf("validateCreateProjectParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getUpstreamRepoCredsSecretName(t *testing.T) {
	type args struct {
		projectName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get secret name",
			args: args{
				projectName: "test-project",
			},
			want: "git-credentials-test-project",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getUpstreamRepoCredsSecretName(tt.args.projectName); got != tt.want {
				t.Errorf("getUpstreamRepoCredsSecretName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_projectManager_createUpstreamRepoCredentials(t *testing.T) {
	type fields struct {
		apiBase *apiBase
	}
	type args struct {
		params *operations.CreateProjectParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create secret succeeds",
			fields: fields{
				apiBase: &apiBase{
					secretStore: &fake.SecretStore{
						UpdateFunc: func(name string, content map[string][]byte) error {
							return nil
						},
					},
					logger: keptncommon.NewLogger("", "", ""),
				},
			},
			args: args{
				params: &operations.CreateProjectParams{
					GitRemoteURL: "my-url",
					GitToken:     "my-token",
					GitUser:      "my-user",
					Name:         stringp("test-project"),
					Shipyard:     nil,
				},
			},
			wantErr: false,
		},
		{
			name: "create secret does not succeed",
			fields: fields{
				apiBase: &apiBase{
					secretStore: &fake.SecretStore{
						UpdateFunc: func(name string, content map[string][]byte) error {
							return errors.New("")
						},
					},
					logger: keptncommon.NewLogger("", "", ""),
				},
			},
			args: args{
				params: &operations.CreateProjectParams{
					GitRemoteURL: "my-url",
					GitToken:     "my-token",
					GitUser:      "my-user",
					Name:         stringp("test-project"),
					Shipyard:     nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &projectManager{
				apiBase: tt.fields.apiBase,
			}
			if err := pm.createUpstreamRepoCredentials(tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("createUpstreamRepoCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_projectManager_CreateProjectTwice(t *testing.T) {
	mockEV := fake.NewEventBroker(t, func(meb *fake.EventBroker, event *models.Event) {
		meb.ReceivedEvents = append(meb.ReceivedEvents, *event)
	}, func(meb *fake.EventBroker) {

	})

	defer mockEV.Server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.Server.URL)

	mockCS := fake.NewSimpleMockConfigurationService()
	defer mockCS.Server.Close()
	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.Server.URL)

	csEndpoint, _ := keptncommon.GetServiceEndpoint("CONFIGURATION_SERVICE")

	pm := &projectManager{
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

	projectName := "my-project"
	createParams := &operations.CreateProjectParams{
		GitRemoteURL: "",
		GitToken:     "",
		GitUser:      "",
		Name:         &projectName,
		Shipyard:     stringp(testBase64EncodedShipyardContent),
	}
	_, _ = pm.createProject(createParams)

	// verify
	expectedProjects := []*keptnapimodels.Project{
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

	if diff := deep.Equal(expectedProjects, mockCS.Projects); len(diff) > 0 {
		t.Errorf("project has not been created correctly")
		for _, d := range diff {
			t.Log(d)
		}
	}

	expectedResources := []string{"shipyard.yaml"}
	if diff := deep.Equal(expectedResources, mockCS.ReceivedResources); len(diff) > 0 {
		t.Errorf("project resources have not been created correctly")
		for _, d := range diff {
			t.Log(d)
		}
	}

	if fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetStartedEventType(keptnv2.ProjectCreateTaskName), "", nil) {
		t.Error("event broker did not receive project.create.started event")
	}

	if fake.ShouldContainEvent(t, mockEV.ReceivedEvents, keptnv2.GetFinishedEventType(keptnv2.ProjectCreateTaskName), "", nil) {
		t.Error("event broker did not receive project.create.started event")
	}

	// create the project again - should return error
	_, err := pm.createProject(createParams)
	if err != errProjectAlreadyExists {
		t.Errorf("expected errProjectAlreadyExists")
	}
}

func Test_projectManager_DeleteProject(t *testing.T) {
	mockEV := fake.NewEventBroker(t, func(meb *fake.EventBroker, event *models.Event) {
		meb.ReceivedEvents = append(meb.ReceivedEvents, *event)
	}, func(meb *fake.EventBroker) {

	})

	defer mockEV.Server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.Server.URL)

	mockCS := fake.NewSimpleMockConfigurationService()
	defer mockCS.Server.Close()
	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.Server.URL)

	mockCS.Projects = []*keptnapimodels.Project{
		{
			ProjectName: "my-project",
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

	pm := &projectManager{
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
				GetFunc: func(name string) (map[string][]byte, error) {
					return map[string][]byte{}, nil
				},
			},
			logger: keptncommon.NewLogger("", "", "shipyard-controller"),
		},
		eventRepo: &fake.EventRepository{
			DeleteEventCollectionsFunc: func(project string) error {
				return nil
			},
		},
		taskSequenceRepo: &fake.TaskSequenceRepository{
			DeleteTaskSequenceCollectionFunc: func(project string) error {
				return nil
			},
		},
	}

	_, _ = pm.deleteProject("my-project")

	// verify
	expectedProjects := []*keptnapimodels.Project{}

	if diff := deep.Equal(expectedProjects, mockCS.Projects); len(diff) > 0 {
		t.Errorf("project has not been created correctly")
		for _, d := range diff {
			t.Log(d)
		}
	}

}
