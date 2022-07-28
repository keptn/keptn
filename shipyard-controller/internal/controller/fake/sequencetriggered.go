// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fake

import (
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"sync"
)

// ISequenceTriggeredHookMock is a mock implementation of controller.ISequenceTriggeredHook.
//
// 	func TestSomethingThatUsesISequenceTriggeredHook(t *testing.T) {
//
// 		// make and configure a mocked controller.ISequenceTriggeredHook
// 		mockedISequenceTriggeredHook := &ISequenceTriggeredHookMock{
// 			OnSequenceTriggeredFunc: func(keptnContextExtendedCE apimodels.KeptnContextExtendedCE)  {
// 				panic("mock out the OnSequenceTriggered method")
// 			},
// 		}
//
// 		// use mockedISequenceTriggeredHook in code that requires controller.ISequenceTriggeredHook
// 		// and then make assertions.
//
// 	}
type ISequenceTriggeredHookMock struct {
	// OnSequenceTriggeredFunc mocks the OnSequenceTriggered method.
	OnSequenceTriggeredFunc func(keptnContextExtendedCE apimodels.KeptnContextExtendedCE)

	// calls tracks calls to the methods.
	calls struct {
		// OnSequenceTriggered holds details about calls to the OnSequenceTriggered method.
		OnSequenceTriggered []struct {
			// KeptnContextExtendedCE is the keptnContextExtendedCE argument value.
			KeptnContextExtendedCE apimodels.KeptnContextExtendedCE
		}
	}
	lockOnSequenceTriggered sync.RWMutex
}

// OnSequenceTriggered calls OnSequenceTriggeredFunc.
func (mock *ISequenceTriggeredHookMock) OnSequenceTriggered(keptnContextExtendedCE apimodels.KeptnContextExtendedCE) {
	if mock.OnSequenceTriggeredFunc == nil {
		panic("ISequenceTriggeredHookMock.OnSequenceTriggeredFunc: method is nil but ISequenceTriggeredHook.OnSequenceTriggered was just called")
	}
	callInfo := struct {
		KeptnContextExtendedCE apimodels.KeptnContextExtendedCE
	}{
		KeptnContextExtendedCE: keptnContextExtendedCE,
	}
	mock.lockOnSequenceTriggered.Lock()
	mock.calls.OnSequenceTriggered = append(mock.calls.OnSequenceTriggered, callInfo)
	mock.lockOnSequenceTriggered.Unlock()
	mock.OnSequenceTriggeredFunc(keptnContextExtendedCE)
}

// OnSequenceTriggeredCalls gets all the calls that were made to OnSequenceTriggered.
// Check the length with:
//     len(mockedISequenceTriggeredHook.OnSequenceTriggeredCalls())
func (mock *ISequenceTriggeredHookMock) OnSequenceTriggeredCalls() []struct {
	KeptnContextExtendedCE apimodels.KeptnContextExtendedCE
} {
	var calls []struct {
		KeptnContextExtendedCE apimodels.KeptnContextExtendedCE
	}
	mock.lockOnSequenceTriggered.RLock()
	calls = mock.calls.OnSequenceTriggered
	mock.lockOnSequenceTriggered.RUnlock()
	return calls
}
