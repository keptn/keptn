// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fake

import (
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"sync"
)

// ISubSequenceFinishedHookMock is a mock implementation of sequencehooks.ISubSequenceFinishedHook.
//
// 	func TestSomethingThatUsesISubSequenceFinishedHook(t *testing.T) {
//
// 		// make and configure a mocked sequencehooks.ISubSequenceFinishedHook
// 		mockedISubSequenceFinishedHook := &ISubSequenceFinishedHookMock{
// 			OnSubSequenceFinishedFunc: func(event apimodels.KeptnContextExtendedCE)  {
// 				panic("mock out the OnSubSequenceFinished method")
// 			},
// 		}
//
// 		// use mockedISubSequenceFinishedHook in code that requires sequencehooks.ISubSequenceFinishedHook
// 		// and then make assertions.
//
// 	}
type ISubSequenceFinishedHookMock struct {
	// OnSubSequenceFinishedFunc mocks the OnSubSequenceFinished method.
	OnSubSequenceFinishedFunc func(event apimodels.KeptnContextExtendedCE)

	// calls tracks calls to the methods.
	calls struct {
		// OnSubSequenceFinished holds details about calls to the OnSubSequenceFinished method.
		OnSubSequenceFinished []struct {
			//models.KeptnContextExtendedCEis the event argument value.
			Event apimodels.KeptnContextExtendedCE
		}
	}
	lockOnSubSequenceFinished sync.RWMutex
}

// OnSubSequenceFinished calls OnSubSequenceFinishedFunc.
func (mock *ISubSequenceFinishedHookMock) OnSubSequenceFinished(event apimodels.KeptnContextExtendedCE) {
	if mock.OnSubSequenceFinishedFunc == nil {
		panic("ISubSequenceFinishedHookMock.OnSubSequenceFinishedFunc: method is nil but ISubSequenceFinishedHook.OnSubSequenceFinished was just called")
	}
	callInfo := struct {
		Event apimodels.KeptnContextExtendedCE
	}{
		Event: event,
	}
	mock.lockOnSubSequenceFinished.Lock()
	mock.calls.OnSubSequenceFinished = append(mock.calls.OnSubSequenceFinished, callInfo)
	mock.lockOnSubSequenceFinished.Unlock()
	mock.OnSubSequenceFinishedFunc(event)
}

// OnSubSequenceFinishedCalls gets all the calls that were made to OnSubSequenceFinished.
// Check the length with:
//     len(mockedISubSequenceFinishedHook.OnSubSequenceFinishedCalls())
func (mock *ISubSequenceFinishedHookMock) OnSubSequenceFinishedCalls() []struct {
	Event apimodels.KeptnContextExtendedCE
} {
	var calls []struct {
		Event apimodels.KeptnContextExtendedCE
	}
	mock.lockOnSubSequenceFinished.RLock()
	calls = mock.calls.OnSubSequenceFinished
	mock.lockOnSubSequenceFinished.RUnlock()
	return calls
}
