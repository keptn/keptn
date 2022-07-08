// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fake

import (
	"github.com/keptn/keptn/api/importer/model"
	"io"
	"net/http"
	"sync"
)

// MockRequestFactory is a mock implementation of execute.requestFactory.
//
// 	func TestSomethingThatUsesrequestFactory(t *testing.T) {
//
// 		// make and configure a mocked execute.requestFactory
// 		mockedrequestFactory := &MockRequestFactory{
// 			CreateRequestFunc: func(tCtx model.TaskContext, host string, body io.Reader) (*http.Request, error) {
// 				panic("mock out the CreateRequest method")
// 			},
// 		}
//
// 		// use mockedrequestFactory in code that requires execute.requestFactory
// 		// and then make assertions.
//
// 	}
type MockRequestFactory struct {
	// CreateRequestFunc mocks the CreateRequest method.
	CreateRequestFunc func(tCtx model.TaskContext, host string, body io.Reader) (*http.Request, error)

	// calls tracks calls to the methods.
	calls struct {
		// CreateRequest holds details about calls to the CreateRequest method.
		CreateRequest []struct {
			// TCtx is the tCtx argument value.
			TCtx model.TaskContext
			// Host is the host argument value.
			Host string
			// Body is the body argument value.
			Body io.Reader
		}
	}
	lockCreateRequest sync.RWMutex
}

// CreateRequest calls CreateRequestFunc.
func (mock *MockRequestFactory) CreateRequest(tCtx model.TaskContext, host string, body io.Reader) (*http.Request, error) {
	if mock.CreateRequestFunc == nil {
		panic("MockRequestFactory.CreateRequestFunc: method is nil but requestFactory.CreateRequest was just called")
	}
	callInfo := struct {
		TCtx model.TaskContext
		Host string
		Body io.Reader
	}{
		TCtx: tCtx,
		Host: host,
		Body: body,
	}
	mock.lockCreateRequest.Lock()
	mock.calls.CreateRequest = append(mock.calls.CreateRequest, callInfo)
	mock.lockCreateRequest.Unlock()
	return mock.CreateRequestFunc(tCtx, host, body)
}

// CreateRequestCalls gets all the calls that were made to CreateRequest.
// Check the length with:
//     len(mockedrequestFactory.CreateRequestCalls())
func (mock *MockRequestFactory) CreateRequestCalls() []struct {
	TCtx model.TaskContext
	Host string
	Body io.Reader
} {
	var calls []struct {
		TCtx model.TaskContext
		Host string
		Body io.Reader
	}
	mock.lockCreateRequest.RLock()
	calls = mock.calls.CreateRequest
	mock.lockCreateRequest.RUnlock()
	return calls
}