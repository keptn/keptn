// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package event_handler_mock

import (
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"sync"
)

// ServiceHandlerMock is a mock implementation of event_handler.ServiceHandler.
//
// 	func TestSomethingThatUsesServiceHandler(t *testing.T) {
//
// 		// make and configure a mocked event_handler.ServiceHandler
// 		mockedServiceHandler := &ServiceHandlerMock{
// 			GetServiceFunc: func(project string, stage string, service string) (*apimodels.Service, error) {
// 				panic("mock out the GetService method")
// 			},
// 		}
//
// 		// use mockedServiceHandler in code that requires event_handler.ServiceHandler
// 		// and then make assertions.
//
// 	}
type ServiceHandlerMock struct {
	// GetServiceFunc mocks the GetService method.
	GetServiceFunc func(project string, stage string, service string) (*apimodels.Service, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetService holds details about calls to the GetService method.
		GetService []struct {
			// Project is the project argument value.
			Project string
			// Stage is the stage argument value.
			Stage string
			// Service is the service argument value.
			Service string
		}
	}
	lockGetService sync.RWMutex
}

// GetService calls GetServiceFunc.
func (mock *ServiceHandlerMock) GetService(project string, stage string, service string) (*apimodels.Service, error) {
	if mock.GetServiceFunc == nil {
		panic("ServiceHandlerMock.GetServiceFunc: method is nil but ServiceHandler.GetService was just called")
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
	mock.lockGetService.Lock()
	mock.calls.GetService = append(mock.calls.GetService, callInfo)
	mock.lockGetService.Unlock()
	return mock.GetServiceFunc(project, stage, service)
}

// GetServiceCalls gets all the calls that were made to GetService.
// Check the length with:
//     len(mockedServiceHandler.GetServiceCalls())
func (mock *ServiceHandlerMock) GetServiceCalls() []struct {
	Project string
	Stage   string
	Service string
} {
	var calls []struct {
		Project string
		Stage   string
		Service string
	}
	mock.lockGetService.RLock()
	calls = mock.calls.GetService
	mock.lockGetService.RUnlock()
	return calls
}
