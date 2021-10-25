// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package db_mock

import (
	"github.com/keptn/keptn/shipyard-controller/models"
	"sync"
)

// ProjectMVRepoMock is a mock implementation of db.ProjectMVRepo.
//
// 	func TestSomethingThatUsesProjectMVRepo(t *testing.T) {
//
// 		// make and configure a mocked db.ProjectMVRepo
// 		mockedProjectMVRepo := &ProjectMVRepoMock{
// 			CloseOpenRemediationsFunc: func(project string, stage string, service string, keptnContext string) error {
// 				panic("mock out the CloseOpenRemediations method")
// 			},
// 			CreateProjectFunc: func(prj *models.ExpandedProject) error {
// 				panic("mock out the CreateProject method")
// 			},
// 			CreateRemediationFunc: func(project string, stage string, service string, remediation *models.Remediation) error {
// 				panic("mock out the CreateRemediation method")
// 			},
// 			CreateServiceFunc: func(project string, stage string, service string) error {
// 				panic("mock out the CreateService method")
// 			},
// 			CreateStageFunc: func(project string, stage string) error {
// 				panic("mock out the CreateStage method")
// 			},
// 			DeleteProjectFunc: func(projectName string) error {
// 				panic("mock out the DeleteProject method")
// 			},
// 			DeleteServiceFunc: func(project string, stage string, service string) error {
// 				panic("mock out the DeleteService method")
// 			},
// 			DeleteStageFunc: func(project string, stage string) error {
// 				panic("mock out the DeleteStage method")
// 			},
// 			DeleteUpstreamInfoFunc: func(projectName string) error {
// 				panic("mock out the DeleteUpstreamInfo method")
// 			},
// 			GetProjectFunc: func(projectName string) (*models.ExpandedProject, error) {
// 				panic("mock out the GetProject method")
// 			},
// 			GetProjectsFunc: func() ([]*models.ExpandedProject, error) {
// 				panic("mock out the GetProjects method")
// 			},
// 			GetServiceFunc: func(projectName string, stageName string, serviceName string) (*models.ExpandedService, error) {
// 				panic("mock out the GetService method")
// 			},
// 			OnSequenceTaskFinishedFunc: func(event models.Event)  {
// 				panic("mock out the OnSequenceTaskFinished method")
// 			},
// 			OnSequenceTaskStartedFunc: func(event models.Event)  {
// 				panic("mock out the OnSequenceTaskStarted method")
// 			},
// 			UpdateEventOfServiceFunc: func(e models.Event) error {
// 				panic("mock out the UpdateEventOfService method")
// 			},
// 			UpdateProjectFunc: func(prj *models.ExpandedProject) error {
// 				panic("mock out the UpdateProject method")
// 			},
// 			UpdateShipyardFunc: func(projectName string, shipyardContent string) error {
// 				panic("mock out the UpdateShipyard method")
// 			},
// 			UpdateUpstreamInfoFunc: func(projectName string, uri string, user string) error {
// 				panic("mock out the UpdateUpstreamInfo method")
// 			},
// 			UpdatedShipyardFunc: func(projectName string, shipyard string) error {
// 				panic("mock out the UpdatedShipyard method")
// 			},
// 		}
//
// 		// use mockedProjectMVRepo in code that requires db.ProjectMVRepo
// 		// and then make assertions.
//
// 	}
type ProjectMVRepoMock struct {
	// CloseOpenRemediationsFunc mocks the CloseOpenRemediations method.
	CloseOpenRemediationsFunc func(project string, stage string, service string, keptnContext string) error

	// CreateProjectFunc mocks the CreateProject method.
	CreateProjectFunc func(prj *models.ExpandedProject) error

	// CreateRemediationFunc mocks the CreateRemediation method.
	CreateRemediationFunc func(project string, stage string, service string, remediation *models.Remediation) error

	// CreateServiceFunc mocks the CreateService method.
	CreateServiceFunc func(project string, stage string, service string) error

	// CreateStageFunc mocks the CreateStage method.
	CreateStageFunc func(project string, stage string) error

	// DeleteProjectFunc mocks the DeleteProject method.
	DeleteProjectFunc func(projectName string) error

	// DeleteServiceFunc mocks the DeleteService method.
	DeleteServiceFunc func(project string, stage string, service string) error

	// DeleteStageFunc mocks the DeleteStage method.
	DeleteStageFunc func(project string, stage string) error

	// DeleteUpstreamInfoFunc mocks the DeleteUpstreamInfo method.
	DeleteUpstreamInfoFunc func(projectName string) error

	// GetProjectFunc mocks the GetProject method.
	GetProjectFunc func(projectName string) (*models.ExpandedProject, error)

	// GetProjectsFunc mocks the GetProjects method.
	GetProjectsFunc func() ([]*models.ExpandedProject, error)

	// GetServiceFunc mocks the GetService method.
	GetServiceFunc func(projectName string, stageName string, serviceName string) (*models.ExpandedService, error)

	// OnSequenceTaskFinishedFunc mocks the OnSequenceTaskFinished method.
	OnSequenceTaskFinishedFunc func(event models.Event)

	// OnSequenceTaskStartedFunc mocks the OnSequenceTaskStarted method.
	OnSequenceTaskStartedFunc func(event models.Event)

	// UpdateEventOfServiceFunc mocks the UpdateEventOfService method.
	UpdateEventOfServiceFunc func(e models.Event) error

	// UpdateProjectFunc mocks the UpdateProject method.
	UpdateProjectFunc func(prj *models.ExpandedProject) error

	// UpdateShipyardFunc mocks the UpdateShipyard method.
	UpdateShipyardFunc func(projectName string, shipyardContent string) error

	// UpdateUpstreamInfoFunc mocks the UpdateUpstreamInfo method.
	UpdateUpstreamInfoFunc func(projectName string, uri string, user string) error

	// UpdatedShipyardFunc mocks the UpdatedShipyard method.
	UpdatedShipyardFunc func(projectName string, shipyard string) error

	// calls tracks calls to the methods.
	calls struct {
		// CloseOpenRemediations holds details about calls to the CloseOpenRemediations method.
		CloseOpenRemediations []struct {
			// Project is the project argument value.
			Project string
			// Stage is the stage argument value.
			Stage string
			// Service is the service argument value.
			Service string
			// KeptnContext is the keptnContext argument value.
			KeptnContext string
		}
		// CreateProject holds details about calls to the CreateProject method.
		CreateProject []struct {
			// Prj is the prj argument value.
			Prj *models.ExpandedProject
		}
		// CreateRemediation holds details about calls to the CreateRemediation method.
		CreateRemediation []struct {
			// Project is the project argument value.
			Project string
			// Stage is the stage argument value.
			Stage string
			// Service is the service argument value.
			Service string
			// Remediation is the remediation argument value.
			Remediation *models.Remediation
		}
		// CreateService holds details about calls to the CreateService method.
		CreateService []struct {
			// Project is the project argument value.
			Project string
			// Stage is the stage argument value.
			Stage string
			// Service is the service argument value.
			Service string
		}
		// CreateStage holds details about calls to the CreateStage method.
		CreateStage []struct {
			// Project is the project argument value.
			Project string
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
			// Project is the project argument value.
			Project string
			// Stage is the stage argument value.
			Stage string
			// Service is the service argument value.
			Service string
		}
		// DeleteStage holds details about calls to the DeleteStage method.
		DeleteStage []struct {
			// Project is the project argument value.
			Project string
			// Stage is the stage argument value.
			Stage string
		}
		// DeleteUpstreamInfo holds details about calls to the DeleteUpstreamInfo method.
		DeleteUpstreamInfo []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
		}
		// GetProject holds details about calls to the GetProject method.
		GetProject []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
		}
		// GetProjects holds details about calls to the GetProjects method.
		GetProjects []struct {
		}
		// GetService holds details about calls to the GetService method.
		GetService []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
			// StageName is the stageName argument value.
			StageName string
			// ServiceName is the serviceName argument value.
			ServiceName string
		}
		// OnSequenceTaskFinished holds details about calls to the OnSequenceTaskFinished method.
		OnSequenceTaskFinished []struct {
			// Event is the event argument value.
			Event models.Event
		}
		// OnSequenceTaskStarted holds details about calls to the OnSequenceTaskStarted method.
		OnSequenceTaskStarted []struct {
			// Event is the event argument value.
			Event models.Event
		}
		// UpdateEventOfService holds details about calls to the UpdateEventOfService method.
		UpdateEventOfService []struct {
			// E is the e argument value.
			E models.Event
		}
		// UpdateProject holds details about calls to the UpdateProject method.
		UpdateProject []struct {
			// Prj is the prj argument value.
			Prj *models.ExpandedProject
		}
		// UpdateShipyard holds details about calls to the UpdateShipyard method.
		UpdateShipyard []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
			// ShipyardContent is the shipyardContent argument value.
			ShipyardContent string
		}
		// UpdateUpstreamInfo holds details about calls to the UpdateUpstreamInfo method.
		UpdateUpstreamInfo []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
			// URI is the uri argument value.
			URI string
			// User is the user argument value.
			User string
		}
		// UpdatedShipyard holds details about calls to the UpdatedShipyard method.
		UpdatedShipyard []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
			// Shipyard is the shipyard argument value.
			Shipyard string
		}
	}
	lockCloseOpenRemediations  sync.RWMutex
	lockCreateProject          sync.RWMutex
	lockCreateRemediation      sync.RWMutex
	lockCreateService          sync.RWMutex
	lockCreateStage            sync.RWMutex
	lockDeleteProject          sync.RWMutex
	lockDeleteService          sync.RWMutex
	lockDeleteStage            sync.RWMutex
	lockDeleteUpstreamInfo     sync.RWMutex
	lockGetProject             sync.RWMutex
	lockGetProjects            sync.RWMutex
	lockGetService             sync.RWMutex
	lockOnSequenceTaskFinished sync.RWMutex
	lockOnSequenceTaskStarted  sync.RWMutex
	lockUpdateEventOfService   sync.RWMutex
	lockUpdateProject          sync.RWMutex
	lockUpdateShipyard         sync.RWMutex
	lockUpdateUpstreamInfo     sync.RWMutex
	lockUpdatedShipyard        sync.RWMutex
}

// CloseOpenRemediations calls CloseOpenRemediationsFunc.
func (mock *ProjectMVRepoMock) CloseOpenRemediations(project string, stage string, service string, keptnContext string) error {
	if mock.CloseOpenRemediationsFunc == nil {
		panic("ProjectMVRepoMock.CloseOpenRemediationsFunc: method is nil but ProjectMVRepo.CloseOpenRemediations was just called")
	}
	callInfo := struct {
		Project      string
		Stage        string
		Service      string
		KeptnContext string
	}{
		Project:      project,
		Stage:        stage,
		Service:      service,
		KeptnContext: keptnContext,
	}
	mock.lockCloseOpenRemediations.Lock()
	mock.calls.CloseOpenRemediations = append(mock.calls.CloseOpenRemediations, callInfo)
	mock.lockCloseOpenRemediations.Unlock()
	return mock.CloseOpenRemediationsFunc(project, stage, service, keptnContext)
}

// CloseOpenRemediationsCalls gets all the calls that were made to CloseOpenRemediations.
// Check the length with:
//     len(mockedProjectMVRepo.CloseOpenRemediationsCalls())
func (mock *ProjectMVRepoMock) CloseOpenRemediationsCalls() []struct {
	Project      string
	Stage        string
	Service      string
	KeptnContext string
} {
	var calls []struct {
		Project      string
		Stage        string
		Service      string
		KeptnContext string
	}
	mock.lockCloseOpenRemediations.RLock()
	calls = mock.calls.CloseOpenRemediations
	mock.lockCloseOpenRemediations.RUnlock()
	return calls
}

// CreateProject calls CreateProjectFunc.
func (mock *ProjectMVRepoMock) CreateProject(prj *models.ExpandedProject) error {
	if mock.CreateProjectFunc == nil {
		panic("ProjectMVRepoMock.CreateProjectFunc: method is nil but ProjectMVRepo.CreateProject was just called")
	}
	callInfo := struct {
		Prj *models.ExpandedProject
	}{
		Prj: prj,
	}
	mock.lockCreateProject.Lock()
	mock.calls.CreateProject = append(mock.calls.CreateProject, callInfo)
	mock.lockCreateProject.Unlock()
	return mock.CreateProjectFunc(prj)
}

// CreateProjectCalls gets all the calls that were made to CreateProject.
// Check the length with:
//     len(mockedProjectMVRepo.CreateProjectCalls())
func (mock *ProjectMVRepoMock) CreateProjectCalls() []struct {
	Prj *models.ExpandedProject
} {
	var calls []struct {
		Prj *models.ExpandedProject
	}
	mock.lockCreateProject.RLock()
	calls = mock.calls.CreateProject
	mock.lockCreateProject.RUnlock()
	return calls
}

// CreateRemediation calls CreateRemediationFunc.
func (mock *ProjectMVRepoMock) CreateRemediation(project string, stage string, service string, remediation *models.Remediation) error {
	if mock.CreateRemediationFunc == nil {
		panic("ProjectMVRepoMock.CreateRemediationFunc: method is nil but ProjectMVRepo.CreateRemediation was just called")
	}
	callInfo := struct {
		Project     string
		Stage       string
		Service     string
		Remediation *models.Remediation
	}{
		Project:     project,
		Stage:       stage,
		Service:     service,
		Remediation: remediation,
	}
	mock.lockCreateRemediation.Lock()
	mock.calls.CreateRemediation = append(mock.calls.CreateRemediation, callInfo)
	mock.lockCreateRemediation.Unlock()
	return mock.CreateRemediationFunc(project, stage, service, remediation)
}

// CreateRemediationCalls gets all the calls that were made to CreateRemediation.
// Check the length with:
//     len(mockedProjectMVRepo.CreateRemediationCalls())
func (mock *ProjectMVRepoMock) CreateRemediationCalls() []struct {
	Project     string
	Stage       string
	Service     string
	Remediation *models.Remediation
} {
	var calls []struct {
		Project     string
		Stage       string
		Service     string
		Remediation *models.Remediation
	}
	mock.lockCreateRemediation.RLock()
	calls = mock.calls.CreateRemediation
	mock.lockCreateRemediation.RUnlock()
	return calls
}

// CreateService calls CreateServiceFunc.
func (mock *ProjectMVRepoMock) CreateService(project string, stage string, service string) error {
	if mock.CreateServiceFunc == nil {
		panic("ProjectMVRepoMock.CreateServiceFunc: method is nil but ProjectMVRepo.CreateService was just called")
	}
	callInfo := struct {
		Project string
		Stage   string
		Service string
	}{
		Project: project,
		Stage:   stage,
		Service: service,
	}
	mock.lockCreateService.Lock()
	mock.calls.CreateService = append(mock.calls.CreateService, callInfo)
	mock.lockCreateService.Unlock()
	return mock.CreateServiceFunc(project, stage, service)
}

// CreateServiceCalls gets all the calls that were made to CreateService.
// Check the length with:
//     len(mockedProjectMVRepo.CreateServiceCalls())
func (mock *ProjectMVRepoMock) CreateServiceCalls() []struct {
	Project string
	Stage   string
	Service string
} {
	var calls []struct {
		Project string
		Stage   string
		Service string
	}
	mock.lockCreateService.RLock()
	calls = mock.calls.CreateService
	mock.lockCreateService.RUnlock()
	return calls
}

// CreateStage calls CreateStageFunc.
func (mock *ProjectMVRepoMock) CreateStage(project string, stage string) error {
	if mock.CreateStageFunc == nil {
		panic("ProjectMVRepoMock.CreateStageFunc: method is nil but ProjectMVRepo.CreateStage was just called")
	}
	callInfo := struct {
		Project string
		Stage   string
	}{
		Project: project,
		Stage:   stage,
	}
	mock.lockCreateStage.Lock()
	mock.calls.CreateStage = append(mock.calls.CreateStage, callInfo)
	mock.lockCreateStage.Unlock()
	return mock.CreateStageFunc(project, stage)
}

// CreateStageCalls gets all the calls that were made to CreateStage.
// Check the length with:
//     len(mockedProjectMVRepo.CreateStageCalls())
func (mock *ProjectMVRepoMock) CreateStageCalls() []struct {
	Project string
	Stage   string
} {
	var calls []struct {
		Project string
		Stage   string
	}
	mock.lockCreateStage.RLock()
	calls = mock.calls.CreateStage
	mock.lockCreateStage.RUnlock()
	return calls
}

// DeleteProject calls DeleteProjectFunc.
func (mock *ProjectMVRepoMock) DeleteProject(projectName string) error {
	if mock.DeleteProjectFunc == nil {
		panic("ProjectMVRepoMock.DeleteProjectFunc: method is nil but ProjectMVRepo.DeleteProject was just called")
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
//     len(mockedProjectMVRepo.DeleteProjectCalls())
func (mock *ProjectMVRepoMock) DeleteProjectCalls() []struct {
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
func (mock *ProjectMVRepoMock) DeleteService(project string, stage string, service string) error {
	if mock.DeleteServiceFunc == nil {
		panic("ProjectMVRepoMock.DeleteServiceFunc: method is nil but ProjectMVRepo.DeleteService was just called")
	}
	callInfo := struct {
		Project string
		Stage   string
		Service string
	}{
		Project: project,
		Stage:   stage,
		Service: service,
	}
	mock.lockDeleteService.Lock()
	mock.calls.DeleteService = append(mock.calls.DeleteService, callInfo)
	mock.lockDeleteService.Unlock()
	return mock.DeleteServiceFunc(project, stage, service)
}

// DeleteServiceCalls gets all the calls that were made to DeleteService.
// Check the length with:
//     len(mockedProjectMVRepo.DeleteServiceCalls())
func (mock *ProjectMVRepoMock) DeleteServiceCalls() []struct {
	Project string
	Stage   string
	Service string
} {
	var calls []struct {
		Project string
		Stage   string
		Service string
	}
	mock.lockDeleteService.RLock()
	calls = mock.calls.DeleteService
	mock.lockDeleteService.RUnlock()
	return calls
}

// DeleteStage calls DeleteStageFunc.
func (mock *ProjectMVRepoMock) DeleteStage(project string, stage string) error {
	if mock.DeleteStageFunc == nil {
		panic("ProjectMVRepoMock.DeleteStageFunc: method is nil but ProjectMVRepo.DeleteStage was just called")
	}
	callInfo := struct {
		Project string
		Stage   string
	}{
		Project: project,
		Stage:   stage,
	}
	mock.lockDeleteStage.Lock()
	mock.calls.DeleteStage = append(mock.calls.DeleteStage, callInfo)
	mock.lockDeleteStage.Unlock()
	return mock.DeleteStageFunc(project, stage)
}

// DeleteStageCalls gets all the calls that were made to DeleteStage.
// Check the length with:
//     len(mockedProjectMVRepo.DeleteStageCalls())
func (mock *ProjectMVRepoMock) DeleteStageCalls() []struct {
	Project string
	Stage   string
} {
	var calls []struct {
		Project string
		Stage   string
	}
	mock.lockDeleteStage.RLock()
	calls = mock.calls.DeleteStage
	mock.lockDeleteStage.RUnlock()
	return calls
}

// DeleteUpstreamInfo calls DeleteUpstreamInfoFunc.
func (mock *ProjectMVRepoMock) DeleteUpstreamInfo(projectName string) error {
	if mock.DeleteUpstreamInfoFunc == nil {
		panic("ProjectMVRepoMock.DeleteUpstreamInfoFunc: method is nil but ProjectMVRepo.DeleteUpstreamInfo was just called")
	}
	callInfo := struct {
		ProjectName string
	}{
		ProjectName: projectName,
	}
	mock.lockDeleteUpstreamInfo.Lock()
	mock.calls.DeleteUpstreamInfo = append(mock.calls.DeleteUpstreamInfo, callInfo)
	mock.lockDeleteUpstreamInfo.Unlock()
	return mock.DeleteUpstreamInfoFunc(projectName)
}

// DeleteUpstreamInfoCalls gets all the calls that were made to DeleteUpstreamInfo.
// Check the length with:
//     len(mockedProjectMVRepo.DeleteUpstreamInfoCalls())
func (mock *ProjectMVRepoMock) DeleteUpstreamInfoCalls() []struct {
	ProjectName string
} {
	var calls []struct {
		ProjectName string
	}
	mock.lockDeleteUpstreamInfo.RLock()
	calls = mock.calls.DeleteUpstreamInfo
	mock.lockDeleteUpstreamInfo.RUnlock()
	return calls
}

// GetProject calls GetProjectFunc.
func (mock *ProjectMVRepoMock) GetProject(projectName string) (*models.ExpandedProject, error) {
	if mock.GetProjectFunc == nil {
		panic("ProjectMVRepoMock.GetProjectFunc: method is nil but ProjectMVRepo.GetProject was just called")
	}
	callInfo := struct {
		ProjectName string
	}{
		ProjectName: projectName,
	}
	mock.lockGetProject.Lock()
	mock.calls.GetProject = append(mock.calls.GetProject, callInfo)
	mock.lockGetProject.Unlock()
	return mock.GetProjectFunc(projectName)
}

// GetProjectCalls gets all the calls that were made to GetProject.
// Check the length with:
//     len(mockedProjectMVRepo.GetProjectCalls())
func (mock *ProjectMVRepoMock) GetProjectCalls() []struct {
	ProjectName string
} {
	var calls []struct {
		ProjectName string
	}
	mock.lockGetProject.RLock()
	calls = mock.calls.GetProject
	mock.lockGetProject.RUnlock()
	return calls
}

// GetProjects calls GetProjectsFunc.
func (mock *ProjectMVRepoMock) GetProjects() ([]*models.ExpandedProject, error) {
	if mock.GetProjectsFunc == nil {
		panic("ProjectMVRepoMock.GetProjectsFunc: method is nil but ProjectMVRepo.GetProjects was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetProjects.Lock()
	mock.calls.GetProjects = append(mock.calls.GetProjects, callInfo)
	mock.lockGetProjects.Unlock()
	return mock.GetProjectsFunc()
}

// GetProjectsCalls gets all the calls that were made to GetProjects.
// Check the length with:
//     len(mockedProjectMVRepo.GetProjectsCalls())
func (mock *ProjectMVRepoMock) GetProjectsCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetProjects.RLock()
	calls = mock.calls.GetProjects
	mock.lockGetProjects.RUnlock()
	return calls
}

// GetService calls GetServiceFunc.
func (mock *ProjectMVRepoMock) GetService(projectName string, stageName string, serviceName string) (*models.ExpandedService, error) {
	if mock.GetServiceFunc == nil {
		panic("ProjectMVRepoMock.GetServiceFunc: method is nil but ProjectMVRepo.GetService was just called")
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
	mock.lockGetService.Lock()
	mock.calls.GetService = append(mock.calls.GetService, callInfo)
	mock.lockGetService.Unlock()
	return mock.GetServiceFunc(projectName, stageName, serviceName)
}

// GetServiceCalls gets all the calls that were made to GetService.
// Check the length with:
//     len(mockedProjectMVRepo.GetServiceCalls())
func (mock *ProjectMVRepoMock) GetServiceCalls() []struct {
	ProjectName string
	StageName   string
	ServiceName string
} {
	var calls []struct {
		ProjectName string
		StageName   string
		ServiceName string
	}
	mock.lockGetService.RLock()
	calls = mock.calls.GetService
	mock.lockGetService.RUnlock()
	return calls
}

// OnSequenceTaskFinished calls OnSequenceTaskFinishedFunc.
func (mock *ProjectMVRepoMock) OnSequenceTaskFinished(event models.Event) {
	if mock.OnSequenceTaskFinishedFunc == nil {
		panic("ProjectMVRepoMock.OnSequenceTaskFinishedFunc: method is nil but ProjectMVRepo.OnSequenceTaskFinished was just called")
	}
	callInfo := struct {
		Event models.Event
	}{
		Event: event,
	}
	mock.lockOnSequenceTaskFinished.Lock()
	mock.calls.OnSequenceTaskFinished = append(mock.calls.OnSequenceTaskFinished, callInfo)
	mock.lockOnSequenceTaskFinished.Unlock()
	mock.OnSequenceTaskFinishedFunc(event)
}

// OnSequenceTaskFinishedCalls gets all the calls that were made to OnSequenceTaskFinished.
// Check the length with:
//     len(mockedProjectMVRepo.OnSequenceTaskFinishedCalls())
func (mock *ProjectMVRepoMock) OnSequenceTaskFinishedCalls() []struct {
	Event models.Event
} {
	var calls []struct {
		Event models.Event
	}
	mock.lockOnSequenceTaskFinished.RLock()
	calls = mock.calls.OnSequenceTaskFinished
	mock.lockOnSequenceTaskFinished.RUnlock()
	return calls
}

// OnSequenceTaskStarted calls OnSequenceTaskStartedFunc.
func (mock *ProjectMVRepoMock) OnSequenceTaskStarted(event models.Event) {
	if mock.OnSequenceTaskStartedFunc == nil {
		panic("ProjectMVRepoMock.OnSequenceTaskStartedFunc: method is nil but ProjectMVRepo.OnSequenceTaskStarted was just called")
	}
	callInfo := struct {
		Event models.Event
	}{
		Event: event,
	}
	mock.lockOnSequenceTaskStarted.Lock()
	mock.calls.OnSequenceTaskStarted = append(mock.calls.OnSequenceTaskStarted, callInfo)
	mock.lockOnSequenceTaskStarted.Unlock()
	mock.OnSequenceTaskStartedFunc(event)
}

// OnSequenceTaskStartedCalls gets all the calls that were made to OnSequenceTaskStarted.
// Check the length with:
//     len(mockedProjectMVRepo.OnSequenceTaskStartedCalls())
func (mock *ProjectMVRepoMock) OnSequenceTaskStartedCalls() []struct {
	Event models.Event
} {
	var calls []struct {
		Event models.Event
	}
	mock.lockOnSequenceTaskStarted.RLock()
	calls = mock.calls.OnSequenceTaskStarted
	mock.lockOnSequenceTaskStarted.RUnlock()
	return calls
}

// UpdateEventOfService calls UpdateEventOfServiceFunc.
func (mock *ProjectMVRepoMock) UpdateEventOfService(e models.Event) error {
	if mock.UpdateEventOfServiceFunc == nil {
		panic("ProjectMVRepoMock.UpdateEventOfServiceFunc: method is nil but ProjectMVRepo.UpdateEventOfService was just called")
	}
	callInfo := struct {
		E models.Event
	}{
		E: e,
	}
	mock.lockUpdateEventOfService.Lock()
	mock.calls.UpdateEventOfService = append(mock.calls.UpdateEventOfService, callInfo)
	mock.lockUpdateEventOfService.Unlock()
	return mock.UpdateEventOfServiceFunc(e)
}

// UpdateEventOfServiceCalls gets all the calls that were made to UpdateEventOfService.
// Check the length with:
//     len(mockedProjectMVRepo.UpdateEventOfServiceCalls())
func (mock *ProjectMVRepoMock) UpdateEventOfServiceCalls() []struct {
	E models.Event
} {
	var calls []struct {
		E models.Event
	}
	mock.lockUpdateEventOfService.RLock()
	calls = mock.calls.UpdateEventOfService
	mock.lockUpdateEventOfService.RUnlock()
	return calls
}

// UpdateProject calls UpdateProjectFunc.
func (mock *ProjectMVRepoMock) UpdateProject(prj *models.ExpandedProject) error {
	if mock.UpdateProjectFunc == nil {
		panic("ProjectMVRepoMock.UpdateProjectFunc: method is nil but ProjectMVRepo.UpdateProject was just called")
	}
	callInfo := struct {
		Prj *models.ExpandedProject
	}{
		Prj: prj,
	}
	mock.lockUpdateProject.Lock()
	mock.calls.UpdateProject = append(mock.calls.UpdateProject, callInfo)
	mock.lockUpdateProject.Unlock()
	return mock.UpdateProjectFunc(prj)
}

// UpdateProjectCalls gets all the calls that were made to UpdateProject.
// Check the length with:
//     len(mockedProjectMVRepo.UpdateProjectCalls())
func (mock *ProjectMVRepoMock) UpdateProjectCalls() []struct {
	Prj *models.ExpandedProject
} {
	var calls []struct {
		Prj *models.ExpandedProject
	}
	mock.lockUpdateProject.RLock()
	calls = mock.calls.UpdateProject
	mock.lockUpdateProject.RUnlock()
	return calls
}

// UpdateShipyard calls UpdateShipyardFunc.
func (mock *ProjectMVRepoMock) UpdateShipyard(projectName string, shipyardContent string) error {
	if mock.UpdateShipyardFunc == nil {
		panic("ProjectMVRepoMock.UpdateShipyardFunc: method is nil but ProjectMVRepo.UpdateShipyard was just called")
	}
	callInfo := struct {
		ProjectName     string
		ShipyardContent string
	}{
		ProjectName:     projectName,
		ShipyardContent: shipyardContent,
	}
	mock.lockUpdateShipyard.Lock()
	mock.calls.UpdateShipyard = append(mock.calls.UpdateShipyard, callInfo)
	mock.lockUpdateShipyard.Unlock()
	return mock.UpdateShipyardFunc(projectName, shipyardContent)
}

// UpdateShipyardCalls gets all the calls that were made to UpdateShipyard.
// Check the length with:
//     len(mockedProjectMVRepo.UpdateShipyardCalls())
func (mock *ProjectMVRepoMock) UpdateShipyardCalls() []struct {
	ProjectName     string
	ShipyardContent string
} {
	var calls []struct {
		ProjectName     string
		ShipyardContent string
	}
	mock.lockUpdateShipyard.RLock()
	calls = mock.calls.UpdateShipyard
	mock.lockUpdateShipyard.RUnlock()
	return calls
}

// UpdateUpstreamInfo calls UpdateUpstreamInfoFunc.
func (mock *ProjectMVRepoMock) UpdateUpstreamInfo(projectName string, uri string, user string) error {
	if mock.UpdateUpstreamInfoFunc == nil {
		panic("ProjectMVRepoMock.UpdateUpstreamInfoFunc: method is nil but ProjectMVRepo.UpdateUpstreamInfo was just called")
	}
	callInfo := struct {
		ProjectName string
		URI         string
		User        string
	}{
		ProjectName: projectName,
		URI:         uri,
		User:        user,
	}
	mock.lockUpdateUpstreamInfo.Lock()
	mock.calls.UpdateUpstreamInfo = append(mock.calls.UpdateUpstreamInfo, callInfo)
	mock.lockUpdateUpstreamInfo.Unlock()
	return mock.UpdateUpstreamInfoFunc(projectName, uri, user)
}

// UpdateUpstreamInfoCalls gets all the calls that were made to UpdateUpstreamInfo.
// Check the length with:
//     len(mockedProjectMVRepo.UpdateUpstreamInfoCalls())
func (mock *ProjectMVRepoMock) UpdateUpstreamInfoCalls() []struct {
	ProjectName string
	URI         string
	User        string
} {
	var calls []struct {
		ProjectName string
		URI         string
		User        string
	}
	mock.lockUpdateUpstreamInfo.RLock()
	calls = mock.calls.UpdateUpstreamInfo
	mock.lockUpdateUpstreamInfo.RUnlock()
	return calls
}

// UpdatedShipyard calls UpdatedShipyardFunc.
func (mock *ProjectMVRepoMock) UpdatedShipyard(projectName string, shipyard string) error {
	if mock.UpdatedShipyardFunc == nil {
		panic("ProjectMVRepoMock.UpdatedShipyardFunc: method is nil but ProjectMVRepo.UpdatedShipyard was just called")
	}
	callInfo := struct {
		ProjectName string
		Shipyard    string
	}{
		ProjectName: projectName,
		Shipyard:    shipyard,
	}
	mock.lockUpdatedShipyard.Lock()
	mock.calls.UpdatedShipyard = append(mock.calls.UpdatedShipyard, callInfo)
	mock.lockUpdatedShipyard.Unlock()
	return mock.UpdatedShipyardFunc(projectName, shipyard)
}

// UpdatedShipyardCalls gets all the calls that were made to UpdatedShipyard.
// Check the length with:
//     len(mockedProjectMVRepo.UpdatedShipyardCalls())
func (mock *ProjectMVRepoMock) UpdatedShipyardCalls() []struct {
	ProjectName string
	Shipyard    string
} {
	var calls []struct {
		ProjectName string
		Shipyard    string
	}
	mock.lockUpdatedShipyard.RLock()
	calls = mock.calls.UpdatedShipyard
	mock.lockUpdatedShipyard.RUnlock()
	return calls
}
