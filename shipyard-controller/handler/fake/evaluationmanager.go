// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fake

import (
	"github.com/keptn/keptn/shipyard-controller/models"
	"sync"
)

// IEvaluationManagerMock is a mock implementation of handler.IEvaluationManager.
//
// 	func TestSomethingThatUsesIEvaluationManager(t *testing.T) {
//
// 		// make and configure a mocked handler.IEvaluationManager
// 		mockedIEvaluationManager := &IEvaluationManagerMock{
// 			CreateEvaluationFunc: func(project string, stage string, service string, params *operations.CreateEvaluationParams) (*operations.CreateEvaluationResponse, *models.Error) {
// 				panic("mock out the CreateEvaluation method")
// 			},
// 		}
//
// 		// use mockedIEvaluationManager in code that requires handler.IEvaluationManager
// 		// and then make assertions.
//
// 	}
type IEvaluationManagerMock struct {
	// CreateEvaluationFunc mocks the CreateEvaluation method.
	CreateEvaluationFunc func(project string, stage string, service string, params *models.CreateEvaluationParams) (*models.CreateEvaluationResponse, *models.Error)

	// calls tracks calls to the methods.
	calls struct {
		// CreateEvaluation holds details about calls to the CreateEvaluation method.
		CreateEvaluation []struct {
			// Project is the project argument value.
			Project string
			// Stage is the stage argument value.
			Stage string
			// Service is the service argument value.
			Service string
			// Params is the params argument value.
			Params *models.CreateEvaluationParams
		}
	}
	lockCreateEvaluation sync.RWMutex
}

// CreateEvaluation calls CreateEvaluationFunc.
func (mock *IEvaluationManagerMock) CreateEvaluation(project string, stage string, service string, params *models.CreateEvaluationParams) (*models.CreateEvaluationResponse, *models.Error) {
	if mock.CreateEvaluationFunc == nil {
		panic("IEvaluationManagerMock.CreateEvaluationFunc: method is nil but IEvaluationManager.CreateEvaluation was just called")
	}
	callInfo := struct {
		Project string
		Stage   string
		Service string
		Params  *models.CreateEvaluationParams
	}{
		Project: project,
		Stage:   stage,
		Service: service,
		Params:  params,
	}
	mock.lockCreateEvaluation.Lock()
	mock.calls.CreateEvaluation = append(mock.calls.CreateEvaluation, callInfo)
	mock.lockCreateEvaluation.Unlock()
	return mock.CreateEvaluationFunc(project, stage, service, params)
}

// CreateEvaluationCalls gets all the calls that were made to CreateEvaluation.
// Check the length with:
//     len(mockedIEvaluationManager.CreateEvaluationCalls())
func (mock *IEvaluationManagerMock) CreateEvaluationCalls() []struct {
	Project string
	Stage   string
	Service string
	Params  *models.CreateEvaluationParams
} {
	var calls []struct {
		Project string
		Stage   string
		Service string
		Params  *models.CreateEvaluationParams
	}
	mock.lockCreateEvaluation.RLock()
	calls = mock.calls.CreateEvaluation
	mock.lockCreateEvaluation.RUnlock()
	return calls
}
