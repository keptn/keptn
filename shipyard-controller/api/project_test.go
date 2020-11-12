package api

import (
	"encoding/json"
	"errors"
	"github.com/go-test/deep"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const testShipyardContent = `apiVersion: spec.keptn.sh/0.2.0
kind: Shipyard
metadata:
  name: shipyard-sockshop
spec:
  stages:
  - name: dev
    sequences:
    - name: artifact-delivery
      tasks:
      - name: deployment
        properties:  
          strategy: direct
      - name: test
        properties:
          kind: functional
      - name: evaluation 
      - name: release 
  - name: hardening
    sequences:
    - name: artifact-delivery
      triggers:
      - dev.artifact-delivery.finished
      tasks:
      - name: deployment
        properties: 
          strategy: blue_green_service
      - name: test
        properties:  
          kind: performance
      - name: evaluation
      - name: release
        
  - name: production
    sequences:
    - name: artifact-delivery 
      triggers:
      - hardening.artifact-delivery.finished
      tasks:
      - name: deployment
        properties:
          strategy: blue_green
      - name: release
      
    - name: remediation
      tasks:
      - name: remediation
      - name: evaluation`

const testShipyardContentWithInvalidVersion = `apiVersion: spec.keptn.sh/0.1.0
kind: Shipyard
metadata:
  name: shipyard-sockshop
spec:
  stages:
  - name: dev
    sequences:
    - name: artifact-delivery
      tasks:
      - name: deployment
        properties:  
          strategy: direct
      - name: test
        properties:
          kind: functional
      - name: evaluation 
      - name: release`

const testShipyardContentWithInvalidStage = `apiVersion: spec.keptn.sh/0.1.0
kind: Shipyard
metadata:
  name: shipyard-sockshop
spec:
  stages:
  - sequences:
    - name: artifact-delivery
      tasks:
      - name: deployment
        properties:  
          strategy: direct
      - name: test
        properties:
          kind: functional
      - name: evaluation 
      - name: release`

const testBase64EncodedShipyardContent = `YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiBzaGlweWFyZC1zb2Nrc2hvcApzcGVjOgogIHN0YWdlczoKICAtIG5hbWU6IGRldgogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0YXNrczoKICAgICAgLSBuYW1lOiBkZXBsb3ltZW50CiAgICAgICAgcHJvcGVydGllczogIAogICAgICAgICAgc3RyYXRlZ3k6IGRpcmVjdAogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOgogICAgICAgICAga2luZDogZnVuY3Rpb25hbAogICAgICAtIG5hbWU6IGV2YWx1YXRpb24gCiAgICAgIC0gbmFtZTogcmVsZWFzZSAKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg==`
const testBase64EncodedShipyardContentWithInvalidVersion = `YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjEuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiBzaGlweWFyZC1zb2Nrc2hvcApzcGVjOgogIHN0YWdlczoKICAtIG5hbWU6IGRldgogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0YXNrczoKICAgICAgLSBuYW1lOiBkZXBsb3ltZW50CiAgICAgICAgcHJvcGVydGllczogIAogICAgICAgICAgc3RyYXRlZ3k6IGRpcmVjdAogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOgogICAgICAgICAga2luZDogZnVuY3Rpb25hbAogICAgICAtIG5hbWU6IGV2YWx1YXRpb24gCiAgICAgIC0gbmFtZTogcmVsZWFzZQ==`
const testBase64EncodedShipyardContentWithInvalidStage = `YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjEuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiBzaGlweWFyZC1zb2Nrc2hvcApzcGVjOgogIHN0YWdlczoKICAtIHNlcXVlbmNlczoKICAgIC0gbmFtZTogYXJ0aWZhY3QtZGVsaXZlcnkKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6ICAKICAgICAgICAgIHN0cmF0ZWd5OiBkaXJlY3QKICAgICAgLSBuYW1lOiB0ZXN0CiAgICAgICAgcHJvcGVydGllczoKICAgICAgICAgIGtpbmQ6IGZ1bmN0aW9uYWwKICAgICAgLSBuYW1lOiBldmFsdWF0aW9uIAogICAgICAtIG5hbWU6IHJlbGVhc2U=`

type mockSecretStore struct {
	create func(name string, content map[string][]byte) error
	delete func(name string) error
	get    func(name string) (map[string][]byte, error)
}

func (ms *mockSecretStore) CreateSecret(name string, content map[string][]byte) error {
	return ms.create(name, content)
}

func (ms *mockSecretStore) DeleteSecret(name string) error {
	return ms.delete(name)
}

func (ms *mockSecretStore) GetSecret(name string) (map[string][]byte, error) {
	return ms.get(name)
}

type mockConfigurationService struct {
	projects          []*keptnapimodels.Project
	receivedResources []string
	server            *httptest.Server
}

func (mcs *mockConfigurationService) get(path string) (interface{}, error) {
	if strings.Contains(path, "/service/") {
		for _, project := range mcs.projects {
			if strings.Contains(path, "/project/"+project.ProjectName) {
				for _, stage := range project.Stages {
					if strings.Contains(path, "stage/"+stage.StageName) {
						for _, service := range stage.Services {
							if strings.Contains(path, "/service/"+service.ServiceName) {
								return service, nil
							}
						}
					}
				}
			}
		}

	} else if strings.Contains(path, "/stage/") {
		for _, project := range mcs.projects {
			if strings.Contains(path, "/project/"+project.ProjectName) {
				for _, stage := range project.Stages {
					if strings.Contains(path, "stage/"+stage.StageName) {
						return stage, nil
					}
				}
			}
		}

	} else if strings.Contains(path, "/project/") {
		for _, project := range mcs.projects {
			if strings.Contains(path, "/project/"+project.ProjectName) {
				return project, nil
			}
		}
	}
	return nil, nil
}

func (mcs *mockConfigurationService) post(body interface{}, path string) (interface{}, error) {
	marshal, _ := json.Marshal(body)
	if strings.Contains(path, "/service") {
		service := &keptnapimodels.Service{}
		_ = json.Unmarshal(marshal, service)
		for _, project := range mcs.projects {
			if strings.Contains(path, "/project/"+project.ProjectName) {
				for _, stage := range project.Stages {
					if strings.Contains(path, "stage/"+stage.StageName) {
						stage.Services = append(stage.Services, service)
						return nil, nil
					}
				}
			}
		}

	} else if strings.Contains(path, "/stage") {
		stage := &keptnapimodels.Stage{}
		_ = json.Unmarshal(marshal, stage)
		for _, project := range mcs.projects {
			if strings.Contains(path, "/project/"+project.ProjectName) {
				project.Stages = append(project.Stages, stage)
			}
		}
	} else if strings.Contains(path, "/resource") {
		resources := &keptnapimodels.Resources{}
		_ = json.Unmarshal(marshal, resources)
		if len(resources.Resources) > 0 {
			mcs.receivedResources = append(mcs.receivedResources, *resources.Resources[0].ResourceURI)
		}
		return &keptnapimodels.Version{
			Version: "",
		}, nil
	} else if strings.Contains(path, "/project") {
		project := &keptnapimodels.Project{}
		_ = json.Unmarshal(marshal, project)
		mcs.projects = append(mcs.projects, project)
		return nil, nil
	}
	return nil, nil
}

func (mcs *mockConfigurationService) put(body interface{}, path string) (interface{}, error) {
	return nil, nil
}

func (mcs *mockConfigurationService) delete(path string) (interface{}, error) {
	if strings.Contains(path, "/service") {
		for _, project := range mcs.projects {
			if strings.Contains(path, "/project/"+project.ProjectName) {
				for _, stage := range project.Stages {
					newServices := []*keptnapimodels.Service{}
					for svcI, svc := range stage.Services {
						if !strings.Contains(path, "/service/"+svc.ServiceName) {
							newServices = append(newServices, stage.Services[svcI])
						}
					}
					stage.Services = newServices
				}
			}
		}
		return nil, nil
	} else if strings.Contains(path, "/project") {
		newProjects := []*keptnapimodels.Project{}

		for index, project := range mcs.projects {
			if !strings.Contains(path, "/project/"+project.ProjectName) {
				newProjects = append(newProjects, mcs.projects[index])
			}
		}
		mcs.projects = newProjects
		return nil, nil
	}
	return nil, nil
}

func newSimpleMockConfigurationService() *mockConfigurationService {
	mcs := &mockConfigurationService{
		projects:          []*keptnapimodels.Project{},
		receivedResources: []string{},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var itf interface{}
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(""))
		}
		json.Unmarshal(bytes, &itf)

		var response interface{}
		switch r.Method {
		case http.MethodGet:
			response, _ = mcs.get(r.URL.Path)
		case http.MethodPost:
			response, _ = mcs.post(itf, r.URL.Path)
		case http.MethodDelete:
			response, _ = mcs.delete(r.URL.Path)
		case http.MethodPut:
			response, _ = mcs.put(itf, r.URL.Path)
		}
		if response != nil {
			w.WriteHeader(http.StatusOK)
			marshal, _ := json.Marshal(response)
			_, _ = w.Write(marshal)
		} else {
			if r.Method == http.MethodGet {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}

	}))
	mcs.server = ts
	return mcs
}

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
					Shipyard: stringp(testShipyardFile),
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
					secretStore: &mockSecretStore{
						create: func(name string, content map[string][]byte) error {
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
					secretStore: &mockSecretStore{
						create: func(name string, content map[string][]byte) error {
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

func Test_projectManager_CreateProjectScenario1(t *testing.T) {
	mockEV := newMockEventbroker(t, func(meb *mockEventBroker, event *models.Event) {
		meb.receivedEvents = append(meb.receivedEvents, *event)
	}, func(meb *mockEventBroker) {

	})

	defer mockEV.server.Close()
	_ = os.Setenv("EVENTBROKER", mockEV.server.URL)

	mockCS := newSimpleMockConfigurationService()
	defer mockCS.server.Close()
	_ = os.Setenv("CONFIGURATION_SERVICE", mockCS.server.URL)

	csEndpoint, _ := keptncommon.GetServiceEndpoint("CONFIGURATION_SERVICE")

	pm := &projectManager{
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

	if diff := deep.Equal(expectedProjects, mockCS.projects); len(diff) > 0 {
		t.Errorf("project has not been created correctly")
		for _, d := range diff {
			t.Log(d)
		}
	}

	expectedResources := []string{"shipyard.yaml"}
	if diff := deep.Equal(expectedResources, mockCS.receivedResources); len(diff) > 0 {
		t.Errorf("project resources have not been created correctly")
		for _, d := range diff {
			t.Log(d)
		}
	}

	if shouldContainEvent(t, mockEV.receivedEvents, keptnv2.GetStartedEventType(keptnv2.ProjectCreateTaskName), "", nil) {
		t.Error("event broker did not receive project.create.started event")
	}

	if shouldContainEvent(t, mockEV.receivedEvents, keptnv2.GetFinishedEventType(keptnv2.ProjectCreateTaskName), "", nil) {
		t.Error("event broker did not receive project.create.started event")
	}

	// create the project again - should return error
	_, err := pm.createProject(createParams)
	if err != errProjectAlreadyExists {
		t.Errorf("expected errProjectAlreadyExists")
	}
}

func Test_projectManager_DeleteProject(t *testing.T) {
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
			secretStore: &mockSecretStore{
				create: func(name string, content map[string][]byte) error {
					return nil
				},
				delete: func(name string) error {
					return nil
				},
				get: func(name string) (map[string][]byte, error) {
					return map[string][]byte{}, nil
				},
			},
			logger: keptncommon.NewLogger("", "", "shipyard-controller"),
		},
		eventRepo: &mockEventRepo{
			deleteCollections: func(project string) error {
				return nil
			},
		},
		taskSequenceRepo: &mockTaskSequenceRepo{
			deleteTaskSequenceCollection: func(project string) error {
				return nil
			},
		},
	}

	_, _ = pm.deleteProject("my-project")

	// verify
	expectedProjects := []*keptnapimodels.Project{}

	if diff := deep.Equal(expectedProjects, mockCS.projects); len(diff) > 0 {
		t.Errorf("project has not been created correctly")
		for _, d := range diff {
			t.Log(d)
		}
	}

}
