// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fake

import (
	"net/http"
	"sync"
)

// KeptnEndpointProviderMock is a mock implementation of execute.KeptnEndpointProvider.
//
// 	func TestSomethingThatUsesKeptnEndpointProvider(t *testing.T) {
//
// 		// make and configure a mocked execute.KeptnEndpointProvider
// 		mockedKeptnEndpointProvider := &KeptnEndpointProviderMock{
// 			GetControlPlaneEndpointFunc: func() string {
// 				panic("mock out the GetControlPlaneEndpoint method")
// 			},
// 			GetSecretsServiceEndpointFunc: func() string {
// 				panic("mock out the GetSecretsServiceEndpoint method")
// 			},
// 		}
//
// 		// use mockedKeptnEndpointProvider in code that requires execute.KeptnEndpointProvider
// 		// and then make assertions.
//
// 	}
type KeptnEndpointProviderMock struct {
	// GetControlPlaneEndpointFunc mocks the GetControlPlaneEndpoint method.
	GetControlPlaneEndpointFunc func() string

	// GetSecretsServiceEndpointFunc mocks the GetSecretsServiceEndpoint method.
	GetSecretsServiceEndpointFunc func() string

	// calls tracks calls to the methods.
	calls struct {
		// GetControlPlaneEndpoint holds details about calls to the GetControlPlaneEndpoint method.
		GetControlPlaneEndpoint []struct {
		}
		// GetSecretsServiceEndpoint holds details about calls to the GetSecretsServiceEndpoint method.
		GetSecretsServiceEndpoint []struct {
		}
	}
	lockGetControlPlaneEndpoint   sync.RWMutex
	lockGetSecretsServiceEndpoint sync.RWMutex
}

// GetControlPlaneEndpoint calls GetControlPlaneEndpointFunc.
func (mock *KeptnEndpointProviderMock) GetControlPlaneEndpoint() string {
	if mock.GetControlPlaneEndpointFunc == nil {
		panic("KeptnEndpointProviderMock.GetControlPlaneEndpointFunc: method is nil but KeptnEndpointProvider.GetControlPlaneEndpoint was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetControlPlaneEndpoint.Lock()
	mock.calls.GetControlPlaneEndpoint = append(mock.calls.GetControlPlaneEndpoint, callInfo)
	mock.lockGetControlPlaneEndpoint.Unlock()
	return mock.GetControlPlaneEndpointFunc()
}

// GetControlPlaneEndpointCalls gets all the calls that were made to GetControlPlaneEndpoint.
// Check the length with:
//     len(mockedKeptnEndpointProvider.GetControlPlaneEndpointCalls())
func (mock *KeptnEndpointProviderMock) GetControlPlaneEndpointCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetControlPlaneEndpoint.RLock()
	calls = mock.calls.GetControlPlaneEndpoint
	mock.lockGetControlPlaneEndpoint.RUnlock()
	return calls
}

// GetSecretsServiceEndpoint calls GetSecretsServiceEndpointFunc.
func (mock *KeptnEndpointProviderMock) GetSecretsServiceEndpoint() string {
	if mock.GetSecretsServiceEndpointFunc == nil {
		panic("KeptnEndpointProviderMock.GetSecretsServiceEndpointFunc: method is nil but KeptnEndpointProvider.GetSecretsServiceEndpoint was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetSecretsServiceEndpoint.Lock()
	mock.calls.GetSecretsServiceEndpoint = append(mock.calls.GetSecretsServiceEndpoint, callInfo)
	mock.lockGetSecretsServiceEndpoint.Unlock()
	return mock.GetSecretsServiceEndpointFunc()
}

// GetSecretsServiceEndpointCalls gets all the calls that were made to GetSecretsServiceEndpoint.
// Check the length with:
//     len(mockedKeptnEndpointProvider.GetSecretsServiceEndpointCalls())
func (mock *KeptnEndpointProviderMock) GetSecretsServiceEndpointCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetSecretsServiceEndpoint.RLock()
	calls = mock.calls.GetSecretsServiceEndpoint
	mock.lockGetSecretsServiceEndpoint.RUnlock()
	return calls
}

// MockHTTPDoer is a mock implementation of execute.httpdoer.
//
// 	func TestSomethingThatUseshttpdoer(t *testing.T) {
//
// 		// make and configure a mocked execute.httpdoer
// 		mockedhttpdoer := &MockHTTPDoer{
// 			DoFunc: func(r *http.Request) (*http.Response, error) {
// 				panic("mock out the Do method")
// 			},
// 		}
//
// 		// use mockedhttpdoer in code that requires execute.httpdoer
// 		// and then make assertions.
//
// 	}
type MockHTTPDoer struct {
	// DoFunc mocks the Do method.
	DoFunc func(r *http.Request) (*http.Response, error)

	// calls tracks calls to the methods.
	calls struct {
		// Do holds details about calls to the Do method.
		Do []struct {
			// R is the r argument value.
			R *http.Request
		}
	}
	lockDo sync.RWMutex
}

// Do calls DoFunc.
func (mock *MockHTTPDoer) Do(r *http.Request) (*http.Response, error) {
	if mock.DoFunc == nil {
		panic("MockHTTPDoer.DoFunc: method is nil but httpdoer.Do was just called")
	}
	callInfo := struct {
		R *http.Request
	}{
		R: r,
	}
	mock.lockDo.Lock()
	mock.calls.Do = append(mock.calls.Do, callInfo)
	mock.lockDo.Unlock()
	return mock.DoFunc(r)
}

// DoCalls gets all the calls that were made to Do.
// Check the length with:
//     len(mockedhttpdoer.DoCalls())
func (mock *MockHTTPDoer) DoCalls() []struct {
	R *http.Request
} {
	var calls []struct {
		R *http.Request
	}
	mock.lockDo.RLock()
	calls = mock.calls.Do
	mock.lockDo.RUnlock()
	return calls
}
