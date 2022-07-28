// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package common_mock

import (
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/internal/configurationstore"
	"sync"
)

// Ensure, that ConfigurationStoreMock does implement configurationstore.ConfigurationStore.
// If this is not the case, regenerate this file with moq.
var _ configurationstore.ConfigurationStore = &ConfigurationStoreMock{}

// ConfigurationStoreMock is a mock implementation of configurationstore.ConfigurationStore.
//
// 	func TestSomethingThatUsesConfigurationStore(t *testing.T) {
//
// 		// make and configure a mocked configurationstore.ConfigurationStore
// 		mockedConfigurationStore := &ConfigurationStoreMock{
// 			CreateProjectFunc: func(project apimodels.Project) error {
// 				panic("mock out the CreateProject method")
// 			},
// 			CreateProjectShipyardFunc: func(projectName string, resources []*apimodels.Resource) error {
// 				panic("mock out the CreateProjectShipyard method")
// 			},
// 			CreateServiceFunc: func(projectName string, stageName string, serviceName string) error {
// 				panic("mock out the CreateService method")
// 			},
// 			CreateStageFunc: func(projectName string, stage string) error {
// 				panic("mock out the CreateStage method")
// 			},
// 			DeleteProjectFunc: func(projectName string) error {
// 				panic("mock out the DeleteProject method")
// 			},
// 			DeleteServiceFunc: func(projectName string, stageName string, serviceName string) error {
// 				panic("mock out the DeleteService method")
// 			},
// 			GetProjectResourceFunc: func(projectName string, resourceURI string) (*apimodels.Resource, error) {
// 				panic("mock out the GetProjectResource method")
// 			},
// 			GetStageResourceFunc: func(projectName string, stageName string, resourceURI string) (*apimodels.Resource, error) {
// 				panic("mock out the GetStageResource method")
// 			},
// 			UpdateProjectFunc: func(project apimodels.Project) error {
// 				panic("mock out the UpdateProject method")
// 			},
// 			UpdateProjectResourceFunc: func(projectName string, resource *apimodels.Resource) error {
// 				panic("mock out the UpdateProjectResource method")
// 			},
// 		}
//
// 		// use mockedConfigurationStore in code that requires configurationstore.ConfigurationStore
// 		// and then make assertions.
//
// 	}
type ConfigurationStoreMock struct {
	// CreateProjectFunc mocks the CreateProject method.
	CreateProjectFunc func(project apimodels.Project) error

	// CreateProjectShipyardFunc mocks the CreateProjectShipyard method.
	CreateProjectShipyardFunc func(projectName string, resources []*apimodels.Resource) error

	// CreateServiceFunc mocks the CreateService method.
	CreateServiceFunc func(projectName string, stageName string, serviceName string) error

	// CreateStageFunc mocks the CreateStage method.
	CreateStageFunc func(projectName string, stage string) error

	// DeleteProjectFunc mocks the DeleteProject method.
	DeleteProjectFunc func(projectName string) error

	// DeleteServiceFunc mocks the DeleteService method.
	DeleteServiceFunc func(projectName string, stageName string, serviceName string) error

	// GetProjectResourceFunc mocks the GetProjectResource method.
	GetProjectResourceFunc func(projectName string, resourceURI string) (*apimodels.Resource, error)

	// GetStageResourceFunc mocks the GetStageResource method.
	GetStageResourceFunc func(projectName string, stageName string, resourceURI string) (*apimodels.Resource, error)

	// UpdateProjectFunc mocks the UpdateProject method.
	UpdateProjectFunc func(project apimodels.Project) error

	// UpdateProjectResourceFunc mocks the UpdateProjectResource method.
	UpdateProjectResourceFunc func(projectName string, resource *apimodels.Resource) error

	// calls tracks calls to the methods.
	calls struct {
		// CreateProject holds details about calls to the CreateProject method.
		CreateProject []struct {
			// Project is the project argument value.
			Project apimodels.Project
		}
		// CreateProjectShipyard holds details about calls to the CreateProjectShipyard method.
		CreateProjectShipyard []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
			// Resources is the resources argument value.
			Resources []*apimodels.Resource
		}
		// CreateService holds details about calls to the CreateService method.
		CreateService []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
			// StageName is the stageName argument value.
			StageName string
			// ServiceName is the serviceName argument value.
			ServiceName string
		}
		// CreateStage holds details about calls to the CreateStage method.
		CreateStage []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
			// Stage is the stage argument value.
			Stage string
		}
		// DeleteProject holds details about calls to the DeleteProject method.
		DeleteProject []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
		}
		// DeleteService holds details about calls to the DeleteService method.
		DeleteService []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
			// StageName is the stageName argument value.
			StageName string
			// ServiceName is the serviceName argument value.
			ServiceName string
		}
		// GetProjectResource holds details about calls to the GetProjectResource method.
		GetProjectResource []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
			// ResourceURI is the resourceURI argument value.
			ResourceURI string
		}
		// GetStageResource holds details about calls to the GetStageResource method.
		GetStageResource []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
			// StageName is the stageName argument value.
			StageName string
			// ResourceURI is the resourceURI argument value.
			ResourceURI string
		}
		// UpdateProject holds details about calls to the UpdateProject method.
		UpdateProject []struct {
			// Project is the project argument value.
			Project apimodels.Project
		}
		// UpdateProjectResource holds details about calls to the UpdateProjectResource method.
		UpdateProjectResource []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
			// Resource is the resource argument value.
			Resource *apimodels.Resource
		}
	}
	lockCreateProject         sync.RWMutex
	lockCreateProjectShipyard sync.RWMutex
	lockCreateService         sync.RWMutex
	lockCreateStage           sync.RWMutex
	lockDeleteProject         sync.RWMutex
	lockDeleteService         sync.RWMutex
	lockGetProjectResource    sync.RWMutex
	lockGetStageResource      sync.RWMutex
	lockUpdateProject         sync.RWMutex
	lockUpdateProjectResource sync.RWMutex
}

// CreateProject calls CreateProjectFunc.
func (mock *ConfigurationStoreMock) CreateProject(project apimodels.Project) error {
	if mock.CreateProjectFunc == nil {
		panic("ConfigurationStoreMock.CreateProjectFunc: method is nil but ConfigurationStore.CreateProject was just called")
	}
	callInfo := struct {
		Project apimodels.Project
	}{
		Project: project,
	}
	mock.lockCreateProject.Lock()
	mock.calls.CreateProject = append(mock.calls.CreateProject, callInfo)
	mock.lockCreateProject.Unlock()
	return mock.CreateProjectFunc(project)
}

// CreateProjectCalls gets all the calls that were made to CreateProject.
// Check the length with:
//     len(mockedConfigurationStore.CreateProjectCalls())
func (mock *ConfigurationStoreMock) CreateProjectCalls() []struct {
	Project apimodels.Project
} {
	var calls []struct {
		Project apimodels.Project
	}
	mock.lockCreateProject.RLock()
	calls = mock.calls.CreateProject
	mock.lockCreateProject.RUnlock()
	return calls
}

// CreateProjectShipyard calls CreateProjectShipyardFunc.
func (mock *ConfigurationStoreMock) CreateProjectShipyard(projectName string, resources []*apimodels.Resource) error {
	if mock.CreateProjectShipyardFunc == nil {
		panic("ConfigurationStoreMock.CreateProjectShipyardFunc: method is nil but ConfigurationStore.CreateProjectShipyard was just called")
	}
	callInfo := struct {
		ProjectName string
		Resources   []*apimodels.Resource
	}{
		ProjectName: projectName,
		Resources:   resources,
	}
	mock.lockCreateProjectShipyard.Lock()
	mock.calls.CreateProjectShipyard = append(mock.calls.CreateProjectShipyard, callInfo)
	mock.lockCreateProjectShipyard.Unlock()
	return mock.CreateProjectShipyardFunc(projectName, resources)
}

// CreateProjectShipyardCalls gets all the calls that were made to CreateProjectShipyard.
// Check the length with:
//     len(mockedConfigurationStore.CreateProjectShipyardCalls())
func (mock *ConfigurationStoreMock) CreateProjectShipyardCalls() []struct {
	ProjectName string
	Resources   []*apimodels.Resource
} {
	var calls []struct {
		ProjectName string
		Resources   []*apimodels.Resource
	}
	mock.lockCreateProjectShipyard.RLock()
	calls = mock.calls.CreateProjectShipyard
	mock.lockCreateProjectShipyard.RUnlock()
	return calls
}

// CreateService calls CreateServiceFunc.
func (mock *ConfigurationStoreMock) CreateService(projectName string, stageName string, serviceName string) error {
	if mock.CreateServiceFunc == nil {
		panic("ConfigurationStoreMock.CreateServiceFunc: method is nil but ConfigurationStore.CreateService was just called")
	}
	callInfo := struct {
		ProjectName string
		StageName   string
		ServiceName string
	}{
		ProjectName: projectName,
		StageName:   stageName,
		ServiceName: serviceName,
	}
	mock.lockCreateService.Lock()
	mock.calls.CreateService = append(mock.calls.CreateService, callInfo)
	mock.lockCreateService.Unlock()
	return mock.CreateServiceFunc(projectName, stageName, serviceName)
}

// CreateServiceCalls gets all the calls that were made to CreateService.
// Check the length with:
//     len(mockedConfigurationStore.CreateServiceCalls())
func (mock *ConfigurationStoreMock) CreateServiceCalls() []struct {
	ProjectName string
	StageName   string
	ServiceName string
} {
	var calls []struct {
		ProjectName string
		StageName   string
		ServiceName string
	}
	mock.lockCreateService.RLock()
	calls = mock.calls.CreateService
	mock.lockCreateService.RUnlock()
	return calls
}

// CreateStage calls CreateStageFunc.
func (mock *ConfigurationStoreMock) CreateStage(projectName string, stage string) error {
	if mock.CreateStageFunc == nil {
		panic("ConfigurationStoreMock.CreateStageFunc: method is nil but ConfigurationStore.CreateStage was just called")
	}
	callInfo := struct {
		ProjectName string
		Stage       string
	}{
		ProjectName: projectName,
		Stage:       stage,
	}
	mock.lockCreateStage.Lock()
	mock.calls.CreateStage = append(mock.calls.CreateStage, callInfo)
	mock.lockCreateStage.Unlock()
	return mock.CreateStageFunc(projectName, stage)
}

// CreateStageCalls gets all the calls that were made to CreateStage.
// Check the length with:
//     len(mockedConfigurationStore.CreateStageCalls())
func (mock *ConfigurationStoreMock) CreateStageCalls() []struct {
	ProjectName string
	Stage       string
} {
	var calls []struct {
		ProjectName string
		Stage       string
	}
	mock.lockCreateStage.RLock()
	calls = mock.calls.CreateStage
	mock.lockCreateStage.RUnlock()
	return calls
}

// DeleteProject calls DeleteProjectFunc.
func (mock *ConfigurationStoreMock) DeleteProject(projectName string) error {
	if mock.DeleteProjectFunc == nil {
		panic("ConfigurationStoreMock.DeleteProjectFunc: method is nil but ConfigurationStore.DeleteProject was just called")
	}
	callInfo := struct {
		ProjectName string
	}{
		ProjectName: projectName,
	}
	mock.lockDeleteProject.Lock()
	mock.calls.DeleteProject = append(mock.calls.DeleteProject, callInfo)
	mock.lockDeleteProject.Unlock()
	return mock.DeleteProjectFunc(projectName)
}

// DeleteProjectCalls gets all the calls that were made to DeleteProject.
// Check the length with:
//     len(mockedConfigurationStore.DeleteProjectCalls())
func (mock *ConfigurationStoreMock) DeleteProjectCalls() []struct {
	ProjectName string
} {
	var calls []struct {
		ProjectName string
	}
	mock.lockDeleteProject.RLock()
	calls = mock.calls.DeleteProject
	mock.lockDeleteProject.RUnlock()
	return calls
}

// DeleteService calls DeleteServiceFunc.
func (mock *ConfigurationStoreMock) DeleteService(projectName string, stageName string, serviceName string) error {
	if mock.DeleteServiceFunc == nil {
		panic("ConfigurationStoreMock.DeleteServiceFunc: method is nil but ConfigurationStore.DeleteService was just called")
	}
	callInfo := struct {
		ProjectName string
		StageName   string
		ServiceName string
	}{
		ProjectName: projectName,
		StageName:   stageName,
		ServiceName: serviceName,
	}
	mock.lockDeleteService.Lock()
	mock.calls.DeleteService = append(mock.calls.DeleteService, callInfo)
	mock.lockDeleteService.Unlock()
	return mock.DeleteServiceFunc(projectName, stageName, serviceName)
}

// DeleteServiceCalls gets all the calls that were made to DeleteService.
// Check the length with:
//     len(mockedConfigurationStore.DeleteServiceCalls())
func (mock *ConfigurationStoreMock) DeleteServiceCalls() []struct {
	ProjectName string
	StageName   string
	ServiceName string
} {
	var calls []struct {
		ProjectName string
		StageName   string
		ServiceName string
	}
	mock.lockDeleteService.RLock()
	calls = mock.calls.DeleteService
	mock.lockDeleteService.RUnlock()
	return calls
}

// GetProjectResource calls GetProjectResourceFunc.
func (mock *ConfigurationStoreMock) GetProjectResource(projectName string, resourceURI string) (*apimodels.Resource, error) {
	if mock.GetProjectResourceFunc == nil {
		panic("ConfigurationStoreMock.GetProjectResourceFunc: method is nil but ConfigurationStore.GetProjectResource was just called")
	}
	callInfo := struct {
		ProjectName string
		ResourceURI string
	}{
		ProjectName: projectName,
		ResourceURI: resourceURI,
	}
	mock.lockGetProjectResource.Lock()
	mock.calls.GetProjectResource = append(mock.calls.GetProjectResource, callInfo)
	mock.lockGetProjectResource.Unlock()
	return mock.GetProjectResourceFunc(projectName, resourceURI)
}

// GetProjectResourceCalls gets all the calls that were made to GetProjectResource.
// Check the length with:
//     len(mockedConfigurationStore.GetProjectResourceCalls())
func (mock *ConfigurationStoreMock) GetProjectResourceCalls() []struct {
	ProjectName string
	ResourceURI string
} {
	var calls []struct {
		ProjectName string
		ResourceURI string
	}
	mock.lockGetProjectResource.RLock()
	calls = mock.calls.GetProjectResource
	mock.lockGetProjectResource.RUnlock()
	return calls
}

// GetStageResource calls GetStageResourceFunc.
func (mock *ConfigurationStoreMock) GetStageResource(projectName string, stageName string, resourceURI string) (*apimodels.Resource, error) {
	if mock.GetStageResourceFunc == nil {
		panic("ConfigurationStoreMock.GetStageResourceFunc: method is nil but ConfigurationStore.GetStageResource was just called")
	}
	callInfo := struct {
		ProjectName string
		StageName   string
		ResourceURI string
	}{
		ProjectName: projectName,
		StageName:   stageName,
		ResourceURI: resourceURI,
	}
	mock.lockGetStageResource.Lock()
	mock.calls.GetStageResource = append(mock.calls.GetStageResource, callInfo)
	mock.lockGetStageResource.Unlock()
	return mock.GetStageResourceFunc(projectName, stageName, resourceURI)
}

// GetStageResourceCalls gets all the calls that were made to GetStageResource.
// Check the length with:
//     len(mockedConfigurationStore.GetStageResourceCalls())
func (mock *ConfigurationStoreMock) GetStageResourceCalls() []struct {
	ProjectName string
	StageName   string
	ResourceURI string
} {
	var calls []struct {
		ProjectName string
		StageName   string
		ResourceURI string
	}
	mock.lockGetStageResource.RLock()
	calls = mock.calls.GetStageResource
	mock.lockGetStageResource.RUnlock()
	return calls
}

// UpdateProject calls UpdateProjectFunc.
func (mock *ConfigurationStoreMock) UpdateProject(project apimodels.Project) error {
	if mock.UpdateProjectFunc == nil {
		panic("ConfigurationStoreMock.UpdateProjectFunc: method is nil but ConfigurationStore.UpdateProject was just called")
	}
	callInfo := struct {
		Project apimodels.Project
	}{
		Project: project,
	}
	mock.lockUpdateProject.Lock()
	mock.calls.UpdateProject = append(mock.calls.UpdateProject, callInfo)
	mock.lockUpdateProject.Unlock()
	return mock.UpdateProjectFunc(project)
}

// UpdateProjectCalls gets all the calls that were made to UpdateProject.
// Check the length with:
//     len(mockedConfigurationStore.UpdateProjectCalls())
func (mock *ConfigurationStoreMock) UpdateProjectCalls() []struct {
	Project apimodels.Project
} {
	var calls []struct {
		Project apimodels.Project
	}
	mock.lockUpdateProject.RLock()
	calls = mock.calls.UpdateProject
	mock.lockUpdateProject.RUnlock()
	return calls
}

// UpdateProjectResource calls UpdateProjectResourceFunc.
func (mock *ConfigurationStoreMock) UpdateProjectResource(projectName string, resource *apimodels.Resource) error {
	if mock.UpdateProjectResourceFunc == nil {
		panic("ConfigurationStoreMock.UpdateProjectResourceFunc: method is nil but ConfigurationStore.UpdateProjectResource was just called")
	}
	callInfo := struct {
		ProjectName string
		Resource    *apimodels.Resource
	}{
		ProjectName: projectName,
		Resource:    resource,
	}
	mock.lockUpdateProjectResource.Lock()
	mock.calls.UpdateProjectResource = append(mock.calls.UpdateProjectResource, callInfo)
	mock.lockUpdateProjectResource.Unlock()
	return mock.UpdateProjectResourceFunc(projectName, resource)
}

// UpdateProjectResourceCalls gets all the calls that were made to UpdateProjectResource.
// Check the length with:
//     len(mockedConfigurationStore.UpdateProjectResourceCalls())
func (mock *ConfigurationStoreMock) UpdateProjectResourceCalls() []struct {
	ProjectName string
	Resource    *apimodels.Resource
} {
	var calls []struct {
		ProjectName string
		Resource    *apimodels.Resource
	}
	mock.lockUpdateProjectResource.RLock()
	calls = mock.calls.UpdateProjectResource
	mock.lockUpdateProjectResource.RUnlock()
	return calls
}
