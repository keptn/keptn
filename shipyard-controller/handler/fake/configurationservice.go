package fake

import (
	"encoding/json"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

type ConfigurationService struct {
	Projects          []*keptnapimodels.Project
	ReceivedResources []string
	Server            *httptest.Server
}

func NewConfigurationService(shipyardContent string) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(shipyardContent))
	}))
	return ts
}

func (mcs *ConfigurationService) get(path string) (interface{}, error) {
	if strings.Contains(path, "/service/") {
		for _, project := range mcs.Projects {
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
		for _, project := range mcs.Projects {
			if strings.Contains(path, "/project/"+project.ProjectName) {
				for _, stage := range project.Stages {
					if strings.Contains(path, "stage/"+stage.StageName) {
						return stage, nil
					}
				}
			}
		}

	} else if strings.Contains(path, "/project/") {
		for _, project := range mcs.Projects {
			if strings.Contains(path, "/project/"+project.ProjectName) {
				return project, nil
			}
		}
	}
	return nil, nil
}

func (mcs *ConfigurationService) post(body interface{}, path string) (interface{}, error) {
	marshal, _ := json.Marshal(body)
	if strings.Contains(path, "/service") {
		service := &keptnapimodels.Service{}
		_ = json.Unmarshal(marshal, service)
		for _, project := range mcs.Projects {
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
		for _, project := range mcs.Projects {
			if strings.Contains(path, "/project/"+project.ProjectName) {
				project.Stages = append(project.Stages, stage)
			}
		}
	} else if strings.Contains(path, "/resource") {
		resources := &keptnapimodels.Resources{}
		_ = json.Unmarshal(marshal, resources)
		if len(resources.Resources) > 0 {
			mcs.ReceivedResources = append(mcs.ReceivedResources, *resources.Resources[0].ResourceURI)
		}
		return &keptnapimodels.Version{
			Version: "",
		}, nil
	} else if strings.Contains(path, "/project") {
		project := &keptnapimodels.Project{}
		_ = json.Unmarshal(marshal, project)
		mcs.Projects = append(mcs.Projects, project)
		return nil, nil
	}
	return nil, nil
}

func (mcs *ConfigurationService) put(body interface{}, path string) (interface{}, error) {
	return nil, nil
}

func (mcs *ConfigurationService) delete(path string) (interface{}, error) {
	if strings.Contains(path, "/service") {
		for _, project := range mcs.Projects {
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

		for index, project := range mcs.Projects {
			if !strings.Contains(path, "/project/"+project.ProjectName) {
				newProjects = append(newProjects, mcs.Projects[index])
			}
		}
		mcs.Projects = newProjects
		return nil, nil
	}
	return nil, nil
}

func NewSimpleMockConfigurationService() *ConfigurationService {
	mcs := &ConfigurationService{
		Projects:          []*keptnapimodels.Project{},
		ReceivedResources: []string{},
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
	mcs.Server = ts
	return mcs
}
