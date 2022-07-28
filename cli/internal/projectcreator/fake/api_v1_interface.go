// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fake

import (
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"sync"
)

// APIV1InterfaceMock is a mock implementation of projectcreator.apiV1Interface.
//
// 	func TestSomethingThatUsesapiV1Interface(t *testing.T) {
//
// 		// make and configure a mocked projectcreator.apiV1Interface
// 		mockedapiV1Interface := &APIV1InterfaceMock{
// 			CreateProjectFunc: func(project apimodels.CreateProject) (string, *apimodels.Error) {
// 				panic("mock out the CreateProject method")
// 			},
// 			CreateServiceFunc: func(project string, service apimodels.CreateService) (string, *apimodels.Error) {
// 				panic("mock out the CreateService method")
// 			},
// 			DeleteProjectFunc: func(project apimodels.Project) (*apimodels.DeleteProjectResponse, *apimodels.Error) {
// 				panic("mock out the DeleteProject method")
// 			},
// 			DeleteServiceFunc: func(project string, service string) (*apimodels.DeleteServiceResponse, *apimodels.Error) {
// 				panic("mock out the DeleteService method")
// 			},
// 			GetMetadataFunc: func() (*apimodels.Metadata, *apimodels.Error) {
// 				panic("mock out the GetMetadata method")
// 			},
// 			SendEventFunc: func(event apimodels.KeptnContextExtendedCE) (*apimodels.EventContext, *apimodels.Error) {
// 				panic("mock out the SendEvent method")
// 			},
// 			TriggerEvaluationFunc: func(project string, stage string, service string, evaluation apimodels.Evaluation) (*apimodels.EventContext, *apimodels.Error) {
// 				panic("mock out the TriggerEvaluation method")
// 			},
// 			UpdateProjectFunc: func(project apimodels.CreateProject) (string, *apimodels.Error) {
// 				panic("mock out the UpdateProject method")
// 			},
// 		}
//
// 		// use mockedapiV1Interface in code that requires projectcreator.apiV1Interface
// 		// and then make assertions.
//
// 	}
type APIV1InterfaceMock struct {
	// CreateProjectFunc mocks the CreateProject method.
	CreateProjectFunc func(project apimodels.CreateProject) (string, *apimodels.Error)

	// CreateServiceFunc mocks the CreateService method.
	CreateServiceFunc func(project string, service apimodels.CreateService) (string, *apimodels.Error)

	// DeleteProjectFunc mocks the DeleteProject method.
	DeleteProjectFunc func(project apimodels.Project) (*apimodels.DeleteProjectResponse, *apimodels.Error)

	// DeleteServiceFunc mocks the DeleteService method.
	DeleteServiceFunc func(project string, service string) (*apimodels.DeleteServiceResponse, *apimodels.Error)

	// GetMetadataFunc mocks the GetMetadata method.
	GetMetadataFunc func() (*apimodels.Metadata, *apimodels.Error)

	// SendEventFunc mocks the SendEvent method.
	SendEventFunc func(event apimodels.KeptnContextExtendedCE) (*apimodels.EventContext, *apimodels.Error)

	// TriggerEvaluationFunc mocks the TriggerEvaluation method.
	TriggerEvaluationFunc func(project string, stage string, service string, evaluation apimodels.Evaluation) (*apimodels.EventContext, *apimodels.Error)

	// UpdateProjectFunc mocks the UpdateProject method.
	UpdateProjectFunc func(project apimodels.CreateProject) (string, *apimodels.Error)

	// calls tracks calls to the methods.
	calls struct {
		// CreateProject holds details about calls to the CreateProject method.
		CreateProject []struct {
			// Project is the project argument value.
			Project apimodels.CreateProject
		}
		// CreateService holds details about calls to the CreateService method.
		CreateService []struct {
			// Project is the project argument value.
			Project string
			// Service is the service argument value.
			Service apimodels.CreateService
		}
		// DeleteProject holds details about calls to the DeleteProject method.
		DeleteProject []struct {
			// Project is the project argument value.
			Project apimodels.Project
		}
		// DeleteService holds details about calls to the DeleteService method.
		DeleteService []struct {
			// Project is the project argument value.
			Project string
			// Service is the service argument value.
			Service string
		}
		// GetMetadata holds details about calls to the GetMetadata method.
		GetMetadata []struct {
		}
		// SendEvent holds details about calls to the SendEvent method.
		SendEvent []struct {
			// Event is the event argument value.
			Event apimodels.KeptnContextExtendedCE
		}
		// TriggerEvaluation holds details about calls to the TriggerEvaluation method.
		TriggerEvaluation []struct {
			// Project is the project argument value.
			Project string
			// Stage is the stage argument value.
			Stage string
			// Service is the service argument value.
			Service string
			// Evaluation is the evaluation argument value.
			Evaluation apimodels.Evaluation
		}
		// UpdateProject holds details about calls to the UpdateProject method.
		UpdateProject []struct {
			// Project is the project argument value.
			Project apimodels.CreateProject
		}
	}
	lockCreateProject     sync.RWMutex
	lockCreateService     sync.RWMutex
	lockDeleteProject     sync.RWMutex
	lockDeleteService     sync.RWMutex
	lockGetMetadata       sync.RWMutex
	lockSendEvent         sync.RWMutex
	lockTriggerEvaluation sync.RWMutex
	lockUpdateProject     sync.RWMutex
}

// CreateProject calls CreateProjectFunc.
func (mock *APIV1InterfaceMock) CreateProject(project apimodels.CreateProject) (string, *apimodels.Error) {
	if mock.CreateProjectFunc == nil {
		panic("APIV1InterfaceMock.CreateProjectFunc: method is nil but apiV1Interface.CreateProject was just called")
	}
	callInfo := struct {
		Project apimodels.CreateProject
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
//     len(mockedapiV1Interface.CreateProjectCalls())
func (mock *APIV1InterfaceMock) CreateProjectCalls() []struct {
	Project apimodels.CreateProject
} {
	var calls []struct {
		Project apimodels.CreateProject
	}
	mock.lockCreateProject.RLock()
	calls = mock.calls.CreateProject
	mock.lockCreateProject.RUnlock()
	return calls
}

// CreateService calls CreateServiceFunc.
func (mock *APIV1InterfaceMock) CreateService(project string, service apimodels.CreateService) (string, *apimodels.Error) {
	if mock.CreateServiceFunc == nil {
		panic("APIV1InterfaceMock.CreateServiceFunc: method is nil but apiV1Interface.CreateService was just called")
	}
	callInfo := struct {
		Project string
		Service apimodels.CreateService
	}{
		Project: project,
		Service: service,
	}
	mock.lockCreateService.Lock()
	mock.calls.CreateService = append(mock.calls.CreateService, callInfo)
	mock.lockCreateService.Unlock()
	return mock.CreateServiceFunc(project, service)
}

// CreateServiceCalls gets all the calls that were made to CreateService.
// Check the length with:
//     len(mockedapiV1Interface.CreateServiceCalls())
func (mock *APIV1InterfaceMock) CreateServiceCalls() []struct {
	Project string
	Service apimodels.CreateService
} {
	var calls []struct {
		Project string
		Service apimodels.CreateService
	}
	mock.lockCreateService.RLock()
	calls = mock.calls.CreateService
	mock.lockCreateService.RUnlock()
	return calls
}

// DeleteProject calls DeleteProjectFunc.
func (mock *APIV1InterfaceMock) DeleteProject(project apimodels.Project) (*apimodels.DeleteProjectResponse, *apimodels.Error) {
	if mock.DeleteProjectFunc == nil {
		panic("APIV1InterfaceMock.DeleteProjectFunc: method is nil but apiV1Interface.DeleteProject was just called")
	}
	callInfo := struct {
		Project apimodels.Project
	}{
		Project: project,
	}
	mock.lockDeleteProject.Lock()
	mock.calls.DeleteProject = append(mock.calls.DeleteProject, callInfo)
	mock.lockDeleteProject.Unlock()
	return mock.DeleteProjectFunc(project)
}

// DeleteProjectCalls gets all the calls that were made to DeleteProject.
// Check the length with:
//     len(mockedapiV1Interface.DeleteProjectCalls())
func (mock *APIV1InterfaceMock) DeleteProjectCalls() []struct {
	Project apimodels.Project
} {
	var calls []struct {
		Project apimodels.Project
	}
	mock.lockDeleteProject.RLock()
	calls = mock.calls.DeleteProject
	mock.lockDeleteProject.RUnlock()
	return calls
}

// DeleteService calls DeleteServiceFunc.
func (mock *APIV1InterfaceMock) DeleteService(project string, service string) (*apimodels.DeleteServiceResponse, *apimodels.Error) {
	if mock.DeleteServiceFunc == nil {
		panic("APIV1InterfaceMock.DeleteServiceFunc: method is nil but apiV1Interface.DeleteService was just called")
	}
	callInfo := struct {
		Project string
		Service string
	}{
		Project: project,
		Service: service,
	}
	mock.lockDeleteService.Lock()
	mock.calls.DeleteService = append(mock.calls.DeleteService, callInfo)
	mock.lockDeleteService.Unlock()
	return mock.DeleteServiceFunc(project, service)
}

// DeleteServiceCalls gets all the calls that were made to DeleteService.
// Check the length with:
//     len(mockedapiV1Interface.DeleteServiceCalls())
func (mock *APIV1InterfaceMock) DeleteServiceCalls() []struct {
	Project string
	Service string
} {
	var calls []struct {
		Project string
		Service string
	}
	mock.lockDeleteService.RLock()
	calls = mock.calls.DeleteService
	mock.lockDeleteService.RUnlock()
	return calls
}

// GetMetadata calls GetMetadataFunc.
func (mock *APIV1InterfaceMock) GetMetadata() (*apimodels.Metadata, *apimodels.Error) {
	if mock.GetMetadataFunc == nil {
		panic("APIV1InterfaceMock.GetMetadataFunc: method is nil but apiV1Interface.GetMetadata was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetMetadata.Lock()
	mock.calls.GetMetadata = append(mock.calls.GetMetadata, callInfo)
	mock.lockGetMetadata.Unlock()
	return mock.GetMetadataFunc()
}

// GetMetadataCalls gets all the calls that were made to GetMetadata.
// Check the length with:
//     len(mockedapiV1Interface.GetMetadataCalls())
func (mock *APIV1InterfaceMock) GetMetadataCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetMetadata.RLock()
	calls = mock.calls.GetMetadata
	mock.lockGetMetadata.RUnlock()
	return calls
}

// SendEvent calls SendEventFunc.
func (mock *APIV1InterfaceMock) SendEvent(event apimodels.KeptnContextExtendedCE) (*apimodels.EventContext, *apimodels.Error) {
	if mock.SendEventFunc == nil {
		panic("APIV1InterfaceMock.SendEventFunc: method is nil but apiV1Interface.SendEvent was just called")
	}
	callInfo := struct {
		Event apimodels.KeptnContextExtendedCE
	}{
		Event: event,
	}
	mock.lockSendEvent.Lock()
	mock.calls.SendEvent = append(mock.calls.SendEvent, callInfo)
	mock.lockSendEvent.Unlock()
	return mock.SendEventFunc(event)
}

// SendEventCalls gets all the calls that were made to SendEvent.
// Check the length with:
//     len(mockedapiV1Interface.SendEventCalls())
func (mock *APIV1InterfaceMock) SendEventCalls() []struct {
	Event apimodels.KeptnContextExtendedCE
} {
	var calls []struct {
		Event apimodels.KeptnContextExtendedCE
	}
	mock.lockSendEvent.RLock()
	calls = mock.calls.SendEvent
	mock.lockSendEvent.RUnlock()
	return calls
}

// TriggerEvaluation calls TriggerEvaluationFunc.
func (mock *APIV1InterfaceMock) TriggerEvaluation(project string, stage string, service string, evaluation apimodels.Evaluation) (*apimodels.EventContext, *apimodels.Error) {
	if mock.TriggerEvaluationFunc == nil {
		panic("APIV1InterfaceMock.TriggerEvaluationFunc: method is nil but apiV1Interface.TriggerEvaluation was just called")
	}
	callInfo := struct {
		Project    string
		Stage      string
		Service    string
		Evaluation apimodels.Evaluation
	}{
		Project:    project,
		Stage:      stage,
		Service:    service,
		Evaluation: evaluation,
	}
	mock.lockTriggerEvaluation.Lock()
	mock.calls.TriggerEvaluation = append(mock.calls.TriggerEvaluation, callInfo)
	mock.lockTriggerEvaluation.Unlock()
	return mock.TriggerEvaluationFunc(project, stage, service, evaluation)
}

// TriggerEvaluationCalls gets all the calls that were made to TriggerEvaluation.
// Check the length with:
//     len(mockedapiV1Interface.TriggerEvaluationCalls())
func (mock *APIV1InterfaceMock) TriggerEvaluationCalls() []struct {
	Project    string
	Stage      string
	Service    string
	Evaluation apimodels.Evaluation
} {
	var calls []struct {
		Project    string
		Stage      string
		Service    string
		Evaluation apimodels.Evaluation
	}
	mock.lockTriggerEvaluation.RLock()
	calls = mock.calls.TriggerEvaluation
	mock.lockTriggerEvaluation.RUnlock()
	return calls
}

// UpdateProject calls UpdateProjectFunc.
func (mock *APIV1InterfaceMock) UpdateProject(project apimodels.CreateProject) (string, *apimodels.Error) {
	if mock.UpdateProjectFunc == nil {
		panic("APIV1InterfaceMock.UpdateProjectFunc: method is nil but apiV1Interface.UpdateProject was just called")
	}
	callInfo := struct {
		Project apimodels.CreateProject
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
//     len(mockedapiV1Interface.UpdateProjectCalls())
func (mock *APIV1InterfaceMock) UpdateProjectCalls() []struct {
	Project apimodels.CreateProject
} {
	var calls []struct {
		Project apimodels.CreateProject
	}
	mock.lockUpdateProject.RLock()
	calls = mock.calls.UpdateProject
	mock.lockUpdateProject.RUnlock()
	return calls
}
