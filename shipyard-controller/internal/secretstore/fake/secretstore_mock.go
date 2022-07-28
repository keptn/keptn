// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fake

import (
	"github.com/keptn/keptn/shipyard-controller/internal/secretstore"
	"sync"
)

// Ensure, that SecretStoreMock does implement secretstore.SecretStore.
// If this is not the case, regenerate this file with moq.
var _ secretstore.SecretStore = &SecretStoreMock{}

// SecretStoreMock is a mock implementation of secretstore.SecretStore.
//
// 	func TestSomethingThatUsesSecretStore(t *testing.T) {
//
// 		// make and configure a mocked secretstore.SecretStore
// 		mockedSecretStore := &SecretStoreMock{
// 			CreateSecretFunc: func(name string, content map[string][]byte) error {
// 				panic("mock out the CreateSecret method")
// 			},
// 			DeleteSecretFunc: func(name string) error {
// 				panic("mock out the DeleteSecret method")
// 			},
// 			GetSecretFunc: func(name string) (map[string][]byte, error) {
// 				panic("mock out the GetSecret method")
// 			},
// 			UpdateSecretFunc: func(name string, content map[string][]byte) error {
// 				panic("mock out the UpdateSecret method")
// 			},
// 		}
//
// 		// use mockedSecretStore in code that requires secretstore.SecretStore
// 		// and then make assertions.
//
// 	}
type SecretStoreMock struct {
	// CreateSecretFunc mocks the CreateSecret method.
	CreateSecretFunc func(name string, content map[string][]byte) error

	// DeleteSecretFunc mocks the DeleteSecret method.
	DeleteSecretFunc func(name string) error

	// GetSecretFunc mocks the GetSecret method.
	GetSecretFunc func(name string) (map[string][]byte, error)

	// UpdateSecretFunc mocks the UpdateSecret method.
	UpdateSecretFunc func(name string, content map[string][]byte) error

	// calls tracks calls to the methods.
	calls struct {
		// CreateSecret holds details about calls to the CreateSecret method.
		CreateSecret []struct {
			// Name is the name argument value.
			Name string
			// Content is the content argument value.
			Content map[string][]byte
		}
		// DeleteSecret holds details about calls to the DeleteSecret method.
		DeleteSecret []struct {
			// Name is the name argument value.
			Name string
		}
		// GetSecret holds details about calls to the GetSecret method.
		GetSecret []struct {
			// Name is the name argument value.
			Name string
		}
		// UpdateSecret holds details about calls to the UpdateSecret method.
		UpdateSecret []struct {
			// Name is the name argument value.
			Name string
			// Content is the content argument value.
			Content map[string][]byte
		}
	}
	lockCreateSecret sync.RWMutex
	lockDeleteSecret sync.RWMutex
	lockGetSecret    sync.RWMutex
	lockUpdateSecret sync.RWMutex
}

// CreateSecret calls CreateSecretFunc.
func (mock *SecretStoreMock) CreateSecret(name string, content map[string][]byte) error {
	if mock.CreateSecretFunc == nil {
		panic("SecretStoreMock.CreateSecretFunc: method is nil but SecretStore.CreateSecret was just called")
	}
	callInfo := struct {
		Name    string
		Content map[string][]byte
	}{
		Name:    name,
		Content: content,
	}
	mock.lockCreateSecret.Lock()
	mock.calls.CreateSecret = append(mock.calls.CreateSecret, callInfo)
	mock.lockCreateSecret.Unlock()
	return mock.CreateSecretFunc(name, content)
}

// CreateSecretCalls gets all the calls that were made to CreateSecret.
// Check the length with:
//     len(mockedSecretStore.CreateSecretCalls())
func (mock *SecretStoreMock) CreateSecretCalls() []struct {
	Name    string
	Content map[string][]byte
} {
	var calls []struct {
		Name    string
		Content map[string][]byte
	}
	mock.lockCreateSecret.RLock()
	calls = mock.calls.CreateSecret
	mock.lockCreateSecret.RUnlock()
	return calls
}

// DeleteSecret calls DeleteSecretFunc.
func (mock *SecretStoreMock) DeleteSecret(name string) error {
	if mock.DeleteSecretFunc == nil {
		panic("SecretStoreMock.DeleteSecretFunc: method is nil but SecretStore.DeleteSecret was just called")
	}
	callInfo := struct {
		Name string
	}{
		Name: name,
	}
	mock.lockDeleteSecret.Lock()
	mock.calls.DeleteSecret = append(mock.calls.DeleteSecret, callInfo)
	mock.lockDeleteSecret.Unlock()
	return mock.DeleteSecretFunc(name)
}

// DeleteSecretCalls gets all the calls that were made to DeleteSecret.
// Check the length with:
//     len(mockedSecretStore.DeleteSecretCalls())
func (mock *SecretStoreMock) DeleteSecretCalls() []struct {
	Name string
} {
	var calls []struct {
		Name string
	}
	mock.lockDeleteSecret.RLock()
	calls = mock.calls.DeleteSecret
	mock.lockDeleteSecret.RUnlock()
	return calls
}

// GetSecret calls GetSecretFunc.
func (mock *SecretStoreMock) GetSecret(name string) (map[string][]byte, error) {
	if mock.GetSecretFunc == nil {
		panic("SecretStoreMock.GetSecretFunc: method is nil but SecretStore.GetSecret was just called")
	}
	callInfo := struct {
		Name string
	}{
		Name: name,
	}
	mock.lockGetSecret.Lock()
	mock.calls.GetSecret = append(mock.calls.GetSecret, callInfo)
	mock.lockGetSecret.Unlock()
	return mock.GetSecretFunc(name)
}

// GetSecretCalls gets all the calls that were made to GetSecret.
// Check the length with:
//     len(mockedSecretStore.GetSecretCalls())
func (mock *SecretStoreMock) GetSecretCalls() []struct {
	Name string
} {
	var calls []struct {
		Name string
	}
	mock.lockGetSecret.RLock()
	calls = mock.calls.GetSecret
	mock.lockGetSecret.RUnlock()
	return calls
}

// UpdateSecret calls UpdateSecretFunc.
func (mock *SecretStoreMock) UpdateSecret(name string, content map[string][]byte) error {
	if mock.UpdateSecretFunc == nil {
		panic("SecretStoreMock.UpdateSecretFunc: method is nil but SecretStore.UpdateSecret was just called")
	}
	callInfo := struct {
		Name    string
		Content map[string][]byte
	}{
		Name:    name,
		Content: content,
	}
	mock.lockUpdateSecret.Lock()
	mock.calls.UpdateSecret = append(mock.calls.UpdateSecret, callInfo)
	mock.lockUpdateSecret.Unlock()
	return mock.UpdateSecretFunc(name, content)
}

// UpdateSecretCalls gets all the calls that were made to UpdateSecret.
// Check the length with:
//     len(mockedSecretStore.UpdateSecretCalls())
func (mock *SecretStoreMock) UpdateSecretCalls() []struct {
	Name    string
	Content map[string][]byte
} {
	var calls []struct {
		Name    string
		Content map[string][]byte
	}
	mock.lockUpdateSecret.RLock()
	calls = mock.calls.UpdateSecret
	mock.lockUpdateSecret.RUnlock()
	return calls
}
