// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fake

import (
	"io"
	"net/http"
	"sync"
)

// KeptnEndpointProviderMock is a mock implementation of execute.KeptnEndpointProvider.
//
// 	func TestSomethingThatUsesKeptnEndpointProvider(t *testing.T) {
//
// 		// make and configure a mocked execute.KeptnEndpointProvider
// 		mockedKeptnEndpointProvider := &KeptnEndpointProviderMock{
// 			GetConfigurationServiceEndpointFunc: func() string {
// 				panic("mock out the GetConfigurationServiceEndpoint method")
// 			},
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
	// GetConfigurationServiceEndpointFunc mocks the GetConfigurationServiceEndpoint method.
	GetConfigurationServiceEndpointFunc func() string

	// GetControlPlaneEndpointFunc mocks the GetControlPlaneEndpoint method.
	GetControlPlaneEndpointFunc func() string

	// GetSecretsServiceEndpointFunc mocks the GetSecretsServiceEndpoint method.
	GetSecretsServiceEndpointFunc func() string

	// calls tracks calls to the methods.
	calls struct {
		// GetConfigurationServiceEndpoint holds details about calls to the GetConfigurationServiceEndpoint method.
		GetConfigurationServiceEndpoint []struct {
		}
		// GetControlPlaneEndpoint holds details about calls to the GetControlPlaneEndpoint method.
		GetControlPlaneEndpoint []struct {
		}
		// GetSecretsServiceEndpoint holds details about calls to the GetSecretsServiceEndpoint method.
		GetSecretsServiceEndpoint []struct {
		}
	}
	lockGetConfigurationServiceEndpoint sync.RWMutex
	lockGetControlPlaneEndpoint         sync.RWMutex
	lockGetSecretsServiceEndpoint       sync.RWMutex
}

// GetConfigurationServiceEndpoint calls GetConfigurationServiceEndpointFunc.
func (mock *KeptnEndpointProviderMock) GetConfigurationServiceEndpoint() string {
	if mock.GetConfigurationServiceEndpointFunc == nil {
		panic("KeptnEndpointProviderMock.GetConfigurationServiceEndpointFunc: method is nil but KeptnEndpointProvider.GetConfigurationServiceEndpoint was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetConfigurationServiceEndpoint.Lock()
	mock.calls.GetConfigurationServiceEndpoint = append(mock.calls.GetConfigurationServiceEndpoint, callInfo)
	mock.lockGetConfigurationServiceEndpoint.Unlock()
	return mock.GetConfigurationServiceEndpointFunc()
}

// GetConfigurationServiceEndpointCalls gets all the calls that were made to GetConfigurationServiceEndpoint.
// Check the length with:
//     len(mockedKeptnEndpointProvider.GetConfigurationServiceEndpointCalls())
func (mock *KeptnEndpointProviderMock) GetConfigurationServiceEndpointCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetConfigurationServiceEndpoint.RLock()
	calls = mock.calls.GetConfigurationServiceEndpoint
	mock.lockGetConfigurationServiceEndpoint.RUnlock()
	return calls
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

// MockResourcePusher is a mock implementation of execute.resourcePusher.
//
// 	func TestSomethingThatUsesresourcePusher(t *testing.T) {
//
// 		// make and configure a mocked execute.resourcePusher
// 		mockedresourcePusher := &MockResourcePusher{
// 			PushToServiceFunc: func(project string, stage string, service string, content io.ReadCloser, resourceURI string) (any, error) {
// 				panic("mock out the PushToService method")
// 			},
// 			PushToStageFunc: func(project string, stage string, content io.ReadCloser, resourceURI string) (any, error) {
// 				panic("mock out the PushToStage method")
// 			},
// 		}
//
// 		// use mockedresourcePusher in code that requires execute.resourcePusher
// 		// and then make assertions.
//
// 	}
type MockResourcePusher struct {
	// PushToServiceFunc mocks the PushToService method.
	PushToServiceFunc func(project string, stage string, service string, content io.ReadCloser, resourceURI string) (any, error)

	// PushToStageFunc mocks the PushToStage method.
	PushToStageFunc func(project string, stage string, content io.ReadCloser, resourceURI string) (any, error)

	// calls tracks calls to the methods.
	calls struct {
		// PushToService holds details about calls to the PushToService method.
		PushToService []struct {
			// Project is the project argument value.
			Project string
			// Stage is the stage argument value.
			Stage string
			// Service is the service argument value.
			Service string
			// Content is the content argument value.
			Content io.ReadCloser
			// ResourceURI is the resourceURI argument value.
			ResourceURI string
		}
		// PushToStage holds details about calls to the PushToStage method.
		PushToStage []struct {
			// Project is the project argument value.
			Project string
			// Stage is the stage argument value.
			Stage string
			// Content is the content argument value.
			Content io.ReadCloser
			// ResourceURI is the resourceURI argument value.
			ResourceURI string
		}
	}
	lockPushToService sync.RWMutex
	lockPushToStage   sync.RWMutex
}

// PushToService calls PushToServiceFunc.
func (mock *MockResourcePusher) PushToService(project string, stage string, service string, content io.ReadCloser, resourceURI string) (any, error) {
	if mock.PushToServiceFunc == nil {
		panic("MockResourcePusher.PushToServiceFunc: method is nil but resourcePusher.PushToService was just called")
	}
	callInfo := struct {
		Project     string
		Stage       string
		Service     string
		Content     io.ReadCloser
		ResourceURI string
	}{
		Project:     project,
		Stage:       stage,
		Service:     service,
		Content:     content,
		ResourceURI: resourceURI,
	}
	mock.lockPushToService.Lock()
	mock.calls.PushToService = append(mock.calls.PushToService, callInfo)
	mock.lockPushToService.Unlock()
	return mock.PushToServiceFunc(project, stage, service, content, resourceURI)
}

// PushToServiceCalls gets all the calls that were made to PushToService.
// Check the length with:
//     len(mockedresourcePusher.PushToServiceCalls())
func (mock *MockResourcePusher) PushToServiceCalls() []struct {
	Project     string
	Stage       string
	Service     string
	Content     io.ReadCloser
	ResourceURI string
} {
	var calls []struct {
		Project     string
		Stage       string
		Service     string
		Content     io.ReadCloser
		ResourceURI string
	}
	mock.lockPushToService.RLock()
	calls = mock.calls.PushToService
	mock.lockPushToService.RUnlock()
	return calls
}

// PushToStage calls PushToStageFunc.
func (mock *MockResourcePusher) PushToStage(project string, stage string, content io.ReadCloser, resourceURI string) (any, error) {
	if mock.PushToStageFunc == nil {
		panic("MockResourcePusher.PushToStageFunc: method is nil but resourcePusher.PushToStage was just called")
	}
	callInfo := struct {
		Project     string
		Stage       string
		Content     io.ReadCloser
		ResourceURI string
	}{
		Project:     project,
		Stage:       stage,
		Content:     content,
		ResourceURI: resourceURI,
	}
	mock.lockPushToStage.Lock()
	mock.calls.PushToStage = append(mock.calls.PushToStage, callInfo)
	mock.lockPushToStage.Unlock()
	return mock.PushToStageFunc(project, stage, content, resourceURI)
}

// PushToStageCalls gets all the calls that were made to PushToStage.
// Check the length with:
//     len(mockedresourcePusher.PushToStageCalls())
func (mock *MockResourcePusher) PushToStageCalls() []struct {
	Project     string
	Stage       string
	Content     io.ReadCloser
	ResourceURI string
} {
	var calls []struct {
		Project     string
		Stage       string
		Content     io.ReadCloser
		ResourceURI string
	}
	mock.lockPushToStage.RLock()
	calls = mock.calls.PushToStage
	mock.lockPushToStage.RUnlock()
	return calls
}
