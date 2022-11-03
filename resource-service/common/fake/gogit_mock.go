// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package common_mock

import (
	"github.com/go-git/go-git/v5"
	"github.com/keptn/keptn/resource-service/common_models"
	"sync"
)

// GogitMock is a mock implementation of common.Gogit.
//
// 	func TestSomethingThatUsesGogit(t *testing.T) {
//
// 		// make and configure a mocked common.Gogit
// 		mockedGogit := &GogitMock{
// 			FetchFunc: func(gitContext common_models.GitContext, repository *git.Repository, options *git.FetchOptions) error {
// 				panic("mock out the Fetch method")
// 			},
// 			PlainCloneFunc: func(gitContext common_models.GitContext, path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
// 				panic("mock out the PlainClone method")
// 			},
// 			PlainInitFunc: func(gitContext common_models.GitContext, path string, isBare bool) (*git.Repository, error) {
// 				panic("mock out the PlainInit method")
// 			},
// 			PlainOpenFunc: func(path string) (*git.Repository, error) {
// 				panic("mock out the PlainOpen method")
// 			},
// 			PullFunc: func(gitContext common_models.GitContext, worktree *git.Worktree, options *git.PullOptions) error {
// 				panic("mock out the Pull method")
// 			},
// 			PushFunc: func(gitContext common_models.GitContext, repository *git.Repository, options *git.PushOptions) error {
// 				panic("mock out the Push method")
// 			},
// 		}
//
// 		// use mockedGogit in code that requires common.Gogit
// 		// and then make assertions.
//
// 	}
type GogitMock struct {
	// FetchFunc mocks the Fetch method.
	FetchFunc func(gitContext common_models.GitContext, repository *git.Repository, options *git.FetchOptions) error

	// PlainCloneFunc mocks the PlainClone method.
	PlainCloneFunc func(gitContext common_models.GitContext, path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)

	// PlainInitFunc mocks the PlainInit method.
	PlainInitFunc func(gitContext common_models.GitContext, path string, isBare bool) (*git.Repository, error)

	// PlainOpenFunc mocks the PlainOpen method.
	PlainOpenFunc func(path string) (*git.Repository, error)

	// PullFunc mocks the Pull method.
	PullFunc func(gitContext common_models.GitContext, worktree *git.Worktree, options *git.PullOptions) error

	// PushFunc mocks the Push method.
	PushFunc func(gitContext common_models.GitContext, repository *git.Repository, options *git.PushOptions) error

	// calls tracks calls to the methods.
	calls struct {
		// Fetch holds details about calls to the Fetch method.
		Fetch []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
			// Repository is the repository argument value.
			Repository *git.Repository
			// Options is the options argument value.
			Options *git.FetchOptions
		}
		// PlainClone holds details about calls to the PlainClone method.
		PlainClone []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
			// Path is the path argument value.
			Path string
			// IsBare is the isBare argument value.
			IsBare bool
			// O is the o argument value.
			O *git.CloneOptions
		}
		// PlainInit holds details about calls to the PlainInit method.
		PlainInit []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
			// Path is the path argument value.
			Path string
			// IsBare is the isBare argument value.
			IsBare bool
		}
		// PlainOpen holds details about calls to the PlainOpen method.
		PlainOpen []struct {
			// Path is the path argument value.
			Path string
		}
		// Pull holds details about calls to the Pull method.
		Pull []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
			// Worktree is the worktree argument value.
			Worktree *git.Worktree
			// Options is the options argument value.
			Options *git.PullOptions
		}
		// Push holds details about calls to the Push method.
		Push []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
			// Repository is the repository argument value.
			Repository *git.Repository
			// Options is the options argument value.
			Options *git.PushOptions
		}
	}
	lockFetch      sync.RWMutex
	lockPlainClone sync.RWMutex
	lockPlainInit  sync.RWMutex
	lockPlainOpen  sync.RWMutex
	lockPull       sync.RWMutex
	lockPush       sync.RWMutex
}

// Fetch calls FetchFunc.
func (mock *GogitMock) Fetch(gitContext common_models.GitContext, repository *git.Repository, options *git.FetchOptions) error {
	if mock.FetchFunc == nil {
		panic("GogitMock.FetchFunc: method is nil but Gogit.Fetch was just called")
	}
	callInfo := struct {
		GitContext common_models.GitContext
		Repository *git.Repository
		Options    *git.FetchOptions
	}{
		GitContext: gitContext,
		Repository: repository,
		Options:    options,
	}
	mock.lockFetch.Lock()
	mock.calls.Fetch = append(mock.calls.Fetch, callInfo)
	mock.lockFetch.Unlock()
	return mock.FetchFunc(gitContext, repository, options)
}

// FetchCalls gets all the calls that were made to Fetch.
// Check the length with:
//     len(mockedGogit.FetchCalls())
func (mock *GogitMock) FetchCalls() []struct {
	GitContext common_models.GitContext
	Repository *git.Repository
	Options    *git.FetchOptions
} {
	var calls []struct {
		GitContext common_models.GitContext
		Repository *git.Repository
		Options    *git.FetchOptions
	}
	mock.lockFetch.RLock()
	calls = mock.calls.Fetch
	mock.lockFetch.RUnlock()
	return calls
}

// PlainClone calls PlainCloneFunc.
func (mock *GogitMock) PlainClone(gitContext common_models.GitContext, path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
	if mock.PlainCloneFunc == nil {
		panic("GogitMock.PlainCloneFunc: method is nil but Gogit.PlainClone was just called")
	}
	callInfo := struct {
		GitContext common_models.GitContext
		Path       string
		IsBare     bool
		O          *git.CloneOptions
	}{
		GitContext: gitContext,
		Path:       path,
		IsBare:     isBare,
		O:          o,
	}
	mock.lockPlainClone.Lock()
	mock.calls.PlainClone = append(mock.calls.PlainClone, callInfo)
	mock.lockPlainClone.Unlock()
	return mock.PlainCloneFunc(gitContext, path, isBare, o)
}

// PlainCloneCalls gets all the calls that were made to PlainClone.
// Check the length with:
//     len(mockedGogit.PlainCloneCalls())
func (mock *GogitMock) PlainCloneCalls() []struct {
	GitContext common_models.GitContext
	Path       string
	IsBare     bool
	O          *git.CloneOptions
} {
	var calls []struct {
		GitContext common_models.GitContext
		Path       string
		IsBare     bool
		O          *git.CloneOptions
	}
	mock.lockPlainClone.RLock()
	calls = mock.calls.PlainClone
	mock.lockPlainClone.RUnlock()
	return calls
}

// PlainInit calls PlainInitFunc.
func (mock *GogitMock) PlainInit(gitContext common_models.GitContext, path string, isBare bool) (*git.Repository, error) {
	if mock.PlainInitFunc == nil {
		panic("GogitMock.PlainInitFunc: method is nil but Gogit.PlainInit was just called")
	}
	callInfo := struct {
		GitContext common_models.GitContext
		Path       string
		IsBare     bool
	}{
		GitContext: gitContext,
		Path:       path,
		IsBare:     isBare,
	}
	mock.lockPlainInit.Lock()
	mock.calls.PlainInit = append(mock.calls.PlainInit, callInfo)
	mock.lockPlainInit.Unlock()
	return mock.PlainInitFunc(gitContext, path, isBare)
}

// PlainInitCalls gets all the calls that were made to PlainInit.
// Check the length with:
//     len(mockedGogit.PlainInitCalls())
func (mock *GogitMock) PlainInitCalls() []struct {
	GitContext common_models.GitContext
	Path       string
	IsBare     bool
} {
	var calls []struct {
		GitContext common_models.GitContext
		Path       string
		IsBare     bool
	}
	mock.lockPlainInit.RLock()
	calls = mock.calls.PlainInit
	mock.lockPlainInit.RUnlock()
	return calls
}

// PlainOpen calls PlainOpenFunc.
func (mock *GogitMock) PlainOpen(path string) (*git.Repository, error) {
	if mock.PlainOpenFunc == nil {
		panic("GogitMock.PlainOpenFunc: method is nil but Gogit.PlainOpen was just called")
	}
	callInfo := struct {
		Path string
	}{
		Path: path,
	}
	mock.lockPlainOpen.Lock()
	mock.calls.PlainOpen = append(mock.calls.PlainOpen, callInfo)
	mock.lockPlainOpen.Unlock()
	return mock.PlainOpenFunc(path)
}

// PlainOpenCalls gets all the calls that were made to PlainOpen.
// Check the length with:
//     len(mockedGogit.PlainOpenCalls())
func (mock *GogitMock) PlainOpenCalls() []struct {
	Path string
} {
	var calls []struct {
		Path string
	}
	mock.lockPlainOpen.RLock()
	calls = mock.calls.PlainOpen
	mock.lockPlainOpen.RUnlock()
	return calls
}

// Pull calls PullFunc.
func (mock *GogitMock) Pull(gitContext common_models.GitContext, worktree *git.Worktree, options *git.PullOptions) error {
	if mock.PullFunc == nil {
		panic("GogitMock.PullFunc: method is nil but Gogit.Pull was just called")
	}
	callInfo := struct {
		GitContext common_models.GitContext
		Worktree   *git.Worktree
		Options    *git.PullOptions
	}{
		GitContext: gitContext,
		Worktree:   worktree,
		Options:    options,
	}
	mock.lockPull.Lock()
	mock.calls.Pull = append(mock.calls.Pull, callInfo)
	mock.lockPull.Unlock()
	return mock.PullFunc(gitContext, worktree, options)
}

// PullCalls gets all the calls that were made to Pull.
// Check the length with:
//     len(mockedGogit.PullCalls())
func (mock *GogitMock) PullCalls() []struct {
	GitContext common_models.GitContext
	Worktree   *git.Worktree
	Options    *git.PullOptions
} {
	var calls []struct {
		GitContext common_models.GitContext
		Worktree   *git.Worktree
		Options    *git.PullOptions
	}
	mock.lockPull.RLock()
	calls = mock.calls.Pull
	mock.lockPull.RUnlock()
	return calls
}

// Push calls PushFunc.
func (mock *GogitMock) Push(gitContext common_models.GitContext, repository *git.Repository, options *git.PushOptions) error {
	if mock.PushFunc == nil {
		panic("GogitMock.PushFunc: method is nil but Gogit.Push was just called")
	}
	callInfo := struct {
		GitContext common_models.GitContext
		Repository *git.Repository
		Options    *git.PushOptions
	}{
		GitContext: gitContext,
		Repository: repository,
		Options:    options,
	}
	mock.lockPush.Lock()
	mock.calls.Push = append(mock.calls.Push, callInfo)
	mock.lockPush.Unlock()
	return mock.PushFunc(gitContext, repository, options)
}

// PushCalls gets all the calls that were made to Push.
// Check the length with:
//     len(mockedGogit.PushCalls())
func (mock *GogitMock) PushCalls() []struct {
	GitContext common_models.GitContext
	Repository *git.Repository
	Options    *git.PushOptions
} {
	var calls []struct {
		GitContext common_models.GitContext
		Repository *git.Repository
		Options    *git.PushOptions
	}
	mock.lockPush.RLock()
	calls = mock.calls.Push
	mock.lockPush.RUnlock()
	return calls
}
