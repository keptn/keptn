// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package db_mock

import (
	"github.com/keptn/keptn/shipyard-controller/models"
	"sync"
)

// ServicesDbOperationsMock is a mock implementation of db.ServicesDbOperations.
//
// 	func TestSomethingThatUsesServicesDbOperations(t *testing.T) {
//
// 		// make and configure a mocked db.ServicesDbOperations
// 		mockedServicesDbOperations := &ServicesDbOperationsMock{
// 			CreateServiceFunc: func(project string, stage string, service string) error {
// 				panic("mock out the CreateService method")
// 			},
// 			DeleteServiceFunc: func(project string, stage string, service string) error {
// 				panic("mock out the DeleteService method")
// 			},
// 			GetProjectFunc: func(projectName string) (*models.ExpandedProject, error) {
// 				panic("mock out the GetProject method")
// 			},
// 		}
//
// 		// use mockedServicesDbOperations in code that requires db.ServicesDbOperations
// 		// and then make assertions.
//
// 	}
type ServicesDbOperationsMock struct {
	// CreateServiceFunc mocks the CreateService method.
	CreateServiceFunc func(project string, stage string, service string) error

	// DeleteServiceFunc mocks the DeleteService method.
	DeleteServiceFunc func(project string, stage string, service string) error

	// GetProjectFunc mocks the GetProject method.
	GetProjectFunc func(projectName string) (*models.ExpandedProject, error)

	// calls tracks calls to the methods.
	calls struct {
		// CreateService holds details about calls to the CreateService method.
		CreateService []struct {
			// Project is the project argument value.
			Project string
			// Stage is the stage argument value.
			Stage string
			// Service is the service argument value.
			Service string
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
		// GetProject holds details about calls to the GetProject method.
		GetProject []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
		}
	}
	lockCreateService sync.RWMutex
	lockDeleteService sync.RWMutex
	lockGetProject    sync.RWMutex
}

// CreateService calls CreateServiceFunc.
func (mock *ServicesDbOperationsMock) CreateService(project string, stage string, service string) error {
	if mock.CreateServiceFunc == nil {
		panic("ServicesDbOperationsMock.CreateServiceFunc: method is nil but ServicesDbOperations.CreateService was just called")
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
//     len(mockedServicesDbOperations.CreateServiceCalls())
func (mock *ServicesDbOperationsMock) CreateServiceCalls() []struct {
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

// DeleteService calls DeleteServiceFunc.
func (mock *ServicesDbOperationsMock) DeleteService(project string, stage string, service string) error {
	if mock.DeleteServiceFunc == nil {
		panic("ServicesDbOperationsMock.DeleteServiceFunc: method is nil but ServicesDbOperations.DeleteService was just called")
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
//     len(mockedServicesDbOperations.DeleteServiceCalls())
func (mock *ServicesDbOperationsMock) DeleteServiceCalls() []struct {
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

// GetProject calls GetProjectFunc.
func (mock *ServicesDbOperationsMock) GetProject(projectName string) (*models.ExpandedProject, error) {
	if mock.GetProjectFunc == nil {
		panic("ServicesDbOperationsMock.GetProjectFunc: method is nil but ServicesDbOperations.GetProject was just called")
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
//     len(mockedServicesDbOperations.GetProjectCalls())
func (mock *ServicesDbOperationsMock) GetProjectCalls() []struct {
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
