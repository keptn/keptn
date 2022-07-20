// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fake

import (
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"sync"
)

// ISequenceTaskFinishedHookMock is a mock implementation of sequencehooks.ISequenceTaskFinishedHook.
//
// 	func TestSomethingThatUsesISequenceTaskFinishedHook(t *testing.T) {
//
// 		// make and configure a mocked sequencehooks.ISequenceTaskFinishedHook
// 		mockedISequenceTaskFinishedHook := &ISequenceTaskFinishedHookMock{
// 			OnSequenceTaskFinishedFunc: func(event apimodels.KeptnContextExtendedCE)  {
// 				panic("mock out the OnSequenceTaskFinished method")
// 			},
// 		}
//
// 		// use mockedISequenceTaskFinishedHook in code that requires sequencehooks.ISequenceTaskFinishedHook
// 		// and then make assertions.
//
// 	}
type ISequenceTaskFinishedHookMock struct {
	// OnSequenceTaskFinishedFunc mocks the OnSequenceTaskFinished method.
	OnSequenceTaskFinishedFunc func(event apimodels.KeptnContextExtendedCE)

	// calls tracks calls to the methods.
	calls struct {
		// OnSequenceTaskFinished holds details about calls to the OnSequenceTaskFinished method.
		OnSequenceTaskFinished []struct {
			//models.KeptnContextExtendedCEis the event argument value.
			Event apimodels.KeptnContextExtendedCE
		}
	}
	lockOnSequenceTaskFinished sync.RWMutex
}

// OnSequenceTaskFinished calls OnSequenceTaskFinishedFunc.
func (mock *ISequenceTaskFinishedHookMock) OnSequenceTaskFinished(event apimodels.KeptnContextExtendedCE) {
	if mock.OnSequenceTaskFinishedFunc == nil {
		panic("ISequenceTaskFinishedHookMock.OnSequenceTaskFinishedFunc: method is nil but ISequenceTaskFinishedHook.OnSequenceTaskFinished was just called")
	}
	callInfo := struct {
		Event apimodels.KeptnContextExtendedCE
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
//     len(mockedISequenceTaskFinishedHook.OnSequenceTaskFinishedCalls())
func (mock *ISequenceTaskFinishedHookMock) OnSequenceTaskFinishedCalls() []struct {
	Event apimodels.KeptnContextExtendedCE
} {
	var calls []struct {
		Event apimodels.KeptnContextExtendedCE
	}
	mock.lockOnSequenceTaskFinished.RLock()
	calls = mock.calls.OnSequenceTaskFinished
	mock.lockOnSequenceTaskFinished.RUnlock()
	return calls
}