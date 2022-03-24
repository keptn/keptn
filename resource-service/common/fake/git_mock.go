// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package common_mock

import (
	"github.com/keptn/keptn/resource-service/common_models"
	"sync"
)

// IGitMock is a mock implementation of common.IGit.
//
// 	func TestSomethingThatUsesIGit(t *testing.T) {
//
// 		// make and configure a mocked common.IGit
// 		mockedIGit := &IGitMock{
// 			CheckoutBranchFunc: func(gitContext common_models.GitContext, branch string) error {
// 				panic("mock out the CheckoutBranch method")
// 			},
// 			CloneRepoFunc: func(gitContext common_models.GitContext) (bool, error) {
// 				panic("mock out the CloneRepo method")
// 			},
// 			CreateBranchFunc: func(gitContext common_models.GitContext, branch string, sourceBranch string) error {
// 				panic("mock out the CreateBranch method")
// 			},
// 			GetCurrentRevisionFunc: func(gitContext common_models.GitContext) (string, error) {
// 				panic("mock out the GetCurrentRevision method")
// 			},
// 			GetDefaultBranchFunc: func(gitContext common_models.GitContext) (string, error) {
// 				panic("mock out the GetDefaultBranch method")
// 			},
// 			GetFileRevisionFunc: func(gitContext common_models.GitContext, revision string, file string) ([]byte, error) {
// 				panic("mock out the GetFileRevision method")
// 			},
// 			MigrateProjectFunc: func(gitContext common_models.GitContext, newMetadatacontent []byte) error {
// 				panic("mock out the MigrateProject method")
// 			},
// 			ProjectExistsFunc: func(gitContext common_models.GitContext) bool {
// 				panic("mock out the ProjectExists method")
// 			},
// 			ProjectRepoExistsFunc: func(projectName string) bool {
// 				panic("mock out the ProjectRepoExists method")
// 			},
// 			PullFunc: func(gitContext common_models.GitContext) error {
// 				panic("mock out the Pull method")
// 			},
// 			PushFunc: func(gitContext common_models.GitContext) error {
// 				panic("mock out the Push method")
// 			},
// 			ResetHardFunc: func(gitContext common_models.GitContext) error {
// 				panic("mock out the ResetHard method")
// 			},
// 			StageAndCommitAllFunc: func(gitContext common_models.GitContext, message string) (string, error) {
// 				panic("mock out the StageAndCommitAll method")
// 			},
// 		}
//
// 		// use mockedIGit in code that requires common.IGit
// 		// and then make assertions.
//
// 	}
type IGitMock struct {
	// CheckoutBranchFunc mocks the CheckoutBranch method.
	CheckoutBranchFunc func(gitContext common_models.GitContext, branch string) error

	// CloneRepoFunc mocks the CloneRepo method.
	CloneRepoFunc func(gitContext common_models.GitContext) (bool, error)

	// CreateBranchFunc mocks the CreateBranch method.
	CreateBranchFunc func(gitContext common_models.GitContext, branch string, sourceBranch string) error

	// GetCurrentRevisionFunc mocks the GetCurrentRevision method.
	GetCurrentRevisionFunc func(gitContext common_models.GitContext) (string, error)

	// GetDefaultBranchFunc mocks the GetDefaultBranch method.
	GetDefaultBranchFunc func(gitContext common_models.GitContext) (string, error)

	// GetFileRevisionFunc mocks the GetFileRevision method.
	GetFileRevisionFunc func(gitContext common_models.GitContext, revision string, file string) ([]byte, error)

	// MigrateProjectFunc mocks the MigrateProject method.
	MigrateProjectFunc func(gitContext common_models.GitContext, newMetadatacontent []byte) error

	// ProjectExistsFunc mocks the ProjectExists method.
	ProjectExistsFunc func(gitContext common_models.GitContext) bool

	// ProjectRepoExistsFunc mocks the ProjectRepoExists method.
	ProjectRepoExistsFunc func(projectName string) bool

	// PullFunc mocks the Pull method.
	PullFunc func(gitContext common_models.GitContext) error

	// PushFunc mocks the Push method.
	PushFunc func(gitContext common_models.GitContext) error

	// ResetHardFunc mocks the ResetHard method.
	ResetHardFunc func(gitContext common_models.GitContext) error

	// StageAndCommitAllFunc mocks the StageAndCommitAll method.
	StageAndCommitAllFunc func(gitContext common_models.GitContext, message string) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// CheckoutBranch holds details about calls to the CheckoutBranch method.
		CheckoutBranch []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
			// Branch is the branch argument value.
			Branch string
		}
		// CloneRepo holds details about calls to the CloneRepo method.
		CloneRepo []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
		}
		// CreateBranch holds details about calls to the CreateBranch method.
		CreateBranch []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
			// Branch is the branch argument value.
			Branch string
			// SourceBranch is the sourceBranch argument value.
			SourceBranch string
		}
		// GetCurrentRevision holds details about calls to the GetCurrentRevision method.
		GetCurrentRevision []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
		}
		// GetDefaultBranch holds details about calls to the GetDefaultBranch method.
		GetDefaultBranch []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
		}
		// GetFileRevision holds details about calls to the GetFileRevision method.
		GetFileRevision []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
			// Revision is the revision argument value.
			Revision string
			// File is the file argument value.
			File string
		}
		// MigrateProject holds details about calls to the MigrateProject method.
		MigrateProject []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
			// NewMetadatacontent is the newMetadatacontent argument value.
			NewMetadatacontent []byte
		}
		// ProjectExists holds details about calls to the ProjectExists method.
		ProjectExists []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
		}
		// ProjectRepoExists holds details about calls to the ProjectRepoExists method.
		ProjectRepoExists []struct {
			// ProjectName is the projectName argument value.
			ProjectName string
		}
		// Pull holds details about calls to the Pull method.
		Pull []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
		}
		// Push holds details about calls to the Push method.
		Push []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
		}
		// ResetHard holds details about calls to the ResetHard method.
		ResetHard []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
		}
		// StageAndCommitAll holds details about calls to the StageAndCommitAll method.
		StageAndCommitAll []struct {
			// GitContext is the gitContext argument value.
			GitContext common_models.GitContext
			// Message is the message argument value.
			Message string
		}
	}
	lockCheckoutBranch     sync.RWMutex
	lockCloneRepo          sync.RWMutex
	lockCreateBranch       sync.RWMutex
	lockGetCurrentRevision sync.RWMutex
	lockGetDefaultBranch   sync.RWMutex
	lockGetFileRevision    sync.RWMutex
	lockMigrateProject     sync.RWMutex
	lockProjectExists      sync.RWMutex
	lockProjectRepoExists  sync.RWMutex
	lockPull               sync.RWMutex
	lockPush               sync.RWMutex
	lockResetHard          sync.RWMutex
	lockStageAndCommitAll  sync.RWMutex
}

// CheckoutBranch calls CheckoutBranchFunc.
func (mock *IGitMock) CheckoutBranch(gitContext common_models.GitContext, branch string) error {
	if mock.CheckoutBranchFunc == nil {
		panic("IGitMock.CheckoutBranchFunc: method is nil but IGit.CheckoutBranch was just called")
	}
	callInfo := struct {
		GitContext common_models.GitContext
		Branch     string
	}{
		GitContext: gitContext,
		Branch:     branch,
	}
	mock.lockCheckoutBranch.Lock()
	mock.calls.CheckoutBranch = append(mock.calls.CheckoutBranch, callInfo)
	mock.lockCheckoutBranch.Unlock()
	return mock.CheckoutBranchFunc(gitContext, branch)
}

// CheckoutBranchCalls gets all the calls that were made to CheckoutBranch.
// Check the length with:
//     len(mockedIGit.CheckoutBranchCalls())
func (mock *IGitMock) CheckoutBranchCalls() []struct {
	GitContext common_models.GitContext
	Branch     string
} {
	var calls []struct {
		GitContext common_models.GitContext
		Branch     string
	}
	mock.lockCheckoutBranch.RLock()
	calls = mock.calls.CheckoutBranch
	mock.lockCheckoutBranch.RUnlock()
	return calls
}

// CloneRepo calls CloneRepoFunc.
func (mock *IGitMock) CloneRepo(gitContext common_models.GitContext) (bool, error) {
	if mock.CloneRepoFunc == nil {
		panic("IGitMock.CloneRepoFunc: method is nil but IGit.CloneRepo was just called")
	}
	callInfo := struct {
		GitContext common_models.GitContext
	}{
		GitContext: gitContext,
	}
	mock.lockCloneRepo.Lock()
	mock.calls.CloneRepo = append(mock.calls.CloneRepo, callInfo)
	mock.lockCloneRepo.Unlock()
	return mock.CloneRepoFunc(gitContext)
}

// CloneRepoCalls gets all the calls that were made to CloneRepo.
// Check the length with:
//     len(mockedIGit.CloneRepoCalls())
func (mock *IGitMock) CloneRepoCalls() []struct {
	GitContext common_models.GitContext
} {
	var calls []struct {
		GitContext common_models.GitContext
	}
	mock.lockCloneRepo.RLock()
	calls = mock.calls.CloneRepo
	mock.lockCloneRepo.RUnlock()
	return calls
}

// CreateBranch calls CreateBranchFunc.
func (mock *IGitMock) CreateBranch(gitContext common_models.GitContext, branch string, sourceBranch string) error {
	if mock.CreateBranchFunc == nil {
		panic("IGitMock.CreateBranchFunc: method is nil but IGit.CreateBranch was just called")
	}
	callInfo := struct {
		GitContext   common_models.GitContext
		Branch       string
		SourceBranch string
	}{
		GitContext:   gitContext,
		Branch:       branch,
		SourceBranch: sourceBranch,
	}
	mock.lockCreateBranch.Lock()
	mock.calls.CreateBranch = append(mock.calls.CreateBranch, callInfo)
	mock.lockCreateBranch.Unlock()
	return mock.CreateBranchFunc(gitContext, branch, sourceBranch)
}

// CreateBranchCalls gets all the calls that were made to CreateBranch.
// Check the length with:
//     len(mockedIGit.CreateBranchCalls())
func (mock *IGitMock) CreateBranchCalls() []struct {
	GitContext   common_models.GitContext
	Branch       string
	SourceBranch string
} {
	var calls []struct {
		GitContext   common_models.GitContext
		Branch       string
		SourceBranch string
	}
	mock.lockCreateBranch.RLock()
	calls = mock.calls.CreateBranch
	mock.lockCreateBranch.RUnlock()
	return calls
}

// GetCurrentRevision calls GetCurrentRevisionFunc.
func (mock *IGitMock) GetCurrentRevision(gitContext common_models.GitContext) (string, error) {
	if mock.GetCurrentRevisionFunc == nil {
		panic("IGitMock.GetCurrentRevisionFunc: method is nil but IGit.GetCurrentRevision was just called")
	}
	callInfo := struct {
		GitContext common_models.GitContext
	}{
		GitContext: gitContext,
	}
	mock.lockGetCurrentRevision.Lock()
	mock.calls.GetCurrentRevision = append(mock.calls.GetCurrentRevision, callInfo)
	mock.lockGetCurrentRevision.Unlock()
	return mock.GetCurrentRevisionFunc(gitContext)
}

// GetCurrentRevisionCalls gets all the calls that were made to GetCurrentRevision.
// Check the length with:
//     len(mockedIGit.GetCurrentRevisionCalls())
func (mock *IGitMock) GetCurrentRevisionCalls() []struct {
	GitContext common_models.GitContext
} {
	var calls []struct {
		GitContext common_models.GitContext
	}
	mock.lockGetCurrentRevision.RLock()
	calls = mock.calls.GetCurrentRevision
	mock.lockGetCurrentRevision.RUnlock()
	return calls
}

// GetDefaultBranch calls GetDefaultBranchFunc.
func (mock *IGitMock) GetDefaultBranch(gitContext common_models.GitContext) (string, error) {
	if mock.GetDefaultBranchFunc == nil {
		panic("IGitMock.GetDefaultBranchFunc: method is nil but IGit.GetDefaultBranch was just called")
	}
	callInfo := struct {
		GitContext common_models.GitContext
	}{
		GitContext: gitContext,
	}
	mock.lockGetDefaultBranch.Lock()
	mock.calls.GetDefaultBranch = append(mock.calls.GetDefaultBranch, callInfo)
	mock.lockGetDefaultBranch.Unlock()
	return mock.GetDefaultBranchFunc(gitContext)
}

// GetDefaultBranchCalls gets all the calls that were made to GetDefaultBranch.
// Check the length with:
//     len(mockedIGit.GetDefaultBranchCalls())
func (mock *IGitMock) GetDefaultBranchCalls() []struct {
	GitContext common_models.GitContext
} {
	var calls []struct {
		GitContext common_models.GitContext
	}
	mock.lockGetDefaultBranch.RLock()
	calls = mock.calls.GetDefaultBranch
	mock.lockGetDefaultBranch.RUnlock()
	return calls
}

// GetFileRevision calls GetFileRevisionFunc.
func (mock *IGitMock) GetFileRevision(gitContext common_models.GitContext, revision string, file string) ([]byte, error) {
	if mock.GetFileRevisionFunc == nil {
		panic("IGitMock.GetFileRevisionFunc: method is nil but IGit.GetFileRevision was just called")
	}
	callInfo := struct {
		GitContext common_models.GitContext
		Revision   string
		File       string
	}{
		GitContext: gitContext,
		Revision:   revision,
		File:       file,
	}
	mock.lockGetFileRevision.Lock()
	mock.calls.GetFileRevision = append(mock.calls.GetFileRevision, callInfo)
	mock.lockGetFileRevision.Unlock()
	return mock.GetFileRevisionFunc(gitContext, revision, file)
}

// GetFileRevisionCalls gets all the calls that were made to GetFileRevision.
// Check the length with:
//     len(mockedIGit.GetFileRevisionCalls())
func (mock *IGitMock) GetFileRevisionCalls() []struct {
	GitContext common_models.GitContext
	Revision   string
	File       string
} {
	var calls []struct {
		GitContext common_models.GitContext
		Revision   string
		File       string
	}
	mock.lockGetFileRevision.RLock()
	calls = mock.calls.GetFileRevision
	mock.lockGetFileRevision.RUnlock()
	return calls
}

// MigrateProject calls MigrateProjectFunc.
func (mock *IGitMock) MigrateProject(gitContext common_models.GitContext, newMetadatacontent []byte) error {
	if mock.MigrateProjectFunc == nil {
		panic("IGitMock.MigrateProjectFunc: method is nil but IGit.MigrateProject was just called")
	}
	callInfo := struct {
		GitContext         common_models.GitContext
		NewMetadatacontent []byte
	}{
		GitContext:         gitContext,
		NewMetadatacontent: newMetadatacontent,
	}
	mock.lockMigrateProject.Lock()
	mock.calls.MigrateProject = append(mock.calls.MigrateProject, callInfo)
	mock.lockMigrateProject.Unlock()
	return mock.MigrateProjectFunc(gitContext, newMetadatacontent)
}

// MigrateProjectCalls gets all the calls that were made to MigrateProject.
// Check the length with:
//     len(mockedIGit.MigrateProjectCalls())
func (mock *IGitMock) MigrateProjectCalls() []struct {
	GitContext         common_models.GitContext
	NewMetadatacontent []byte
} {
	var calls []struct {
		GitContext         common_models.GitContext
		NewMetadatacontent []byte
	}
	mock.lockMigrateProject.RLock()
	calls = mock.calls.MigrateProject
	mock.lockMigrateProject.RUnlock()
	return calls
}

// ProjectExists calls ProjectExistsFunc.
func (mock *IGitMock) ProjectExists(gitContext common_models.GitContext) bool {
	if mock.ProjectExistsFunc == nil {
		panic("IGitMock.ProjectExistsFunc: method is nil but IGit.ProjectExists was just called")
	}
	callInfo := struct {
		GitContext common_models.GitContext
	}{
		GitContext: gitContext,
	}
	mock.lockProjectExists.Lock()
	mock.calls.ProjectExists = append(mock.calls.ProjectExists, callInfo)
	mock.lockProjectExists.Unlock()
	return mock.ProjectExistsFunc(gitContext)
}

// ProjectExistsCalls gets all the calls that were made to ProjectExists.
// Check the length with:
//     len(mockedIGit.ProjectExistsCalls())
func (mock *IGitMock) ProjectExistsCalls() []struct {
	GitContext common_models.GitContext
} {
	var calls []struct {
		GitContext common_models.GitContext
	}
	mock.lockProjectExists.RLock()
	calls = mock.calls.ProjectExists
	mock.lockProjectExists.RUnlock()
	return calls
}

// ProjectRepoExists calls ProjectRepoExistsFunc.
func (mock *IGitMock) ProjectRepoExists(projectName string) bool {
	if mock.ProjectRepoExistsFunc == nil {
		panic("IGitMock.ProjectRepoExistsFunc: method is nil but IGit.ProjectRepoExists was just called")
	}
	callInfo := struct {
		ProjectName string
	}{
		ProjectName: projectName,
	}
	mock.lockProjectRepoExists.Lock()
	mock.calls.ProjectRepoExists = append(mock.calls.ProjectRepoExists, callInfo)
	mock.lockProjectRepoExists.Unlock()
	return mock.ProjectRepoExistsFunc(projectName)
}

// ProjectRepoExistsCalls gets all the calls that were made to ProjectRepoExists.
// Check the length with:
//     len(mockedIGit.ProjectRepoExistsCalls())
func (mock *IGitMock) ProjectRepoExistsCalls() []struct {
	ProjectName string
} {
	var calls []struct {
		ProjectName string
	}
	mock.lockProjectRepoExists.RLock()
	calls = mock.calls.ProjectRepoExists
	mock.lockProjectRepoExists.RUnlock()
	return calls
}

// Pull calls PullFunc.
func (mock *IGitMock) Pull(gitContext common_models.GitContext) error {
	if mock.PullFunc == nil {
		panic("IGitMock.PullFunc: method is nil but IGit.Pull was just called")
	}
	callInfo := struct {
		GitContext common_models.GitContext
	}{
		GitContext: gitContext,
	}
	mock.lockPull.Lock()
	mock.calls.Pull = append(mock.calls.Pull, callInfo)
	mock.lockPull.Unlock()
	return mock.PullFunc(gitContext)
}

// PullCalls gets all the calls that were made to Pull.
// Check the length with:
//     len(mockedIGit.PullCalls())
func (mock *IGitMock) PullCalls() []struct {
	GitContext common_models.GitContext
} {
	var calls []struct {
		GitContext common_models.GitContext
	}
	mock.lockPull.RLock()
	calls = mock.calls.Pull
	mock.lockPull.RUnlock()
	return calls
}

// Push calls PushFunc.
func (mock *IGitMock) Push(gitContext common_models.GitContext) error {
	if mock.PushFunc == nil {
		panic("IGitMock.PushFunc: method is nil but IGit.Push was just called")
	}
	callInfo := struct {
		GitContext common_models.GitContext
	}{
		GitContext: gitContext,
	}
	mock.lockPush.Lock()
	mock.calls.Push = append(mock.calls.Push, callInfo)
	mock.lockPush.Unlock()
	return mock.PushFunc(gitContext)
}

// PushCalls gets all the calls that were made to Push.
// Check the length with:
//     len(mockedIGit.PushCalls())
func (mock *IGitMock) PushCalls() []struct {
	GitContext common_models.GitContext
} {
	var calls []struct {
		GitContext common_models.GitContext
	}
	mock.lockPush.RLock()
	calls = mock.calls.Push
	mock.lockPush.RUnlock()
	return calls
}

// ResetHard calls ResetHardFunc.
func (mock *IGitMock) ResetHard(gitContext common_models.GitContext, revision string) error {
	if mock.ResetHardFunc == nil {
		panic("IGitMock.ResetHardFunc: method is nil but IGit.ResetHard was just called")
	}
	callInfo := struct {
		GitContext common_models.GitContext
	}{
		GitContext: gitContext,
	}
	mock.lockResetHard.Lock()
	mock.calls.ResetHard = append(mock.calls.ResetHard, callInfo)
	mock.lockResetHard.Unlock()
	return mock.ResetHardFunc(gitContext)
}

// ResetHardCalls gets all the calls that were made to ResetHard.
// Check the length with:
//     len(mockedIGit.ResetHardCalls())
func (mock *IGitMock) ResetHardCalls() []struct {
	GitContext common_models.GitContext
} {
	var calls []struct {
		GitContext common_models.GitContext
	}
	mock.lockResetHard.RLock()
	calls = mock.calls.ResetHard
	mock.lockResetHard.RUnlock()
	return calls
}

// StageAndCommitAll calls StageAndCommitAllFunc.
func (mock *IGitMock) StageAndCommitAll(gitContext common_models.GitContext, message string) (string, error) {
	if mock.StageAndCommitAllFunc == nil {
		panic("IGitMock.StageAndCommitAllFunc: method is nil but IGit.StageAndCommitAll was just called")
	}
	callInfo := struct {
		GitContext common_models.GitContext
		Message    string
	}{
		GitContext: gitContext,
		Message:    message,
	}
	mock.lockStageAndCommitAll.Lock()
	mock.calls.StageAndCommitAll = append(mock.calls.StageAndCommitAll, callInfo)
	mock.lockStageAndCommitAll.Unlock()
	return mock.StageAndCommitAllFunc(gitContext, message)
}

// StageAndCommitAllCalls gets all the calls that were made to StageAndCommitAll.
// Check the length with:
//     len(mockedIGit.StageAndCommitAllCalls())
func (mock *IGitMock) StageAndCommitAllCalls() []struct {
	GitContext common_models.GitContext
	Message    string
} {
	var calls []struct {
		GitContext common_models.GitContext
		Message    string
	}
	mock.lockStageAndCommitAll.RLock()
	calls = mock.calls.StageAndCommitAll
	mock.lockStageAndCommitAll.RUnlock()
	return calls
}
