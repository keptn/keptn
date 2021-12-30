package handler

import (
	"github.com/keptn/keptn/resource-service/common"
	common_mock "github.com/keptn/keptn/resource-service/common/fake"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type testResourceManagerFields struct {
	git              common.IGit
	credentialReader common.CredentialReader
	fileSystem       common.IFileSystem
}

func TestResourceManager_CreateResources(t *testing.T) {

}

type fakeFileInfo struct {
	name  string
	isDir bool
}

func newFakeFileInfo(name string, isDir bool) *fakeFileInfo {
	return &fakeFileInfo{name: name, isDir: isDir}
}

func (f fakeFileInfo) Name() string {
	return f.name
}

func (fakeFileInfo) Size() int64 {
	return 100
}

func (fakeFileInfo) Mode() fs.FileMode {
	return os.ModePerm
}

func (fakeFileInfo) ModTime() time.Time {
	return time.Time{}
}

func (f fakeFileInfo) IsDir() bool {
	return f.isDir
}

func (fakeFileInfo) Sys() interface{} {
	return nil
}

func getTestResourceManagerFields() testResourceManagerFields {
	return testResourceManagerFields{
		git: &common_mock.IGitMock{
			CheckoutBranchFunc: func(gitContext common.GitContext, branch string) error {
				return nil
			},
			CloneRepoFunc: func(gitContext common.GitContext) (bool, error) {
				return true, nil
			},
			CreateBranchFunc: func(gitContext common.GitContext, branch string, sourceBranch string) error {
				return nil
			},
			GetCurrentRevisionFunc: func(gitContext common.GitContext) (string, error) {
				return "my-revision", nil
			},
			GetDefaultBranchFunc: func(gitContext common.GitContext) (string, error) {
				return "main", nil
			},
			GetFileRevisionFunc: func(gitContext common.GitContext, path string, revision string, file string) ([]byte, error) {
				return []byte("my-content"), nil
			},
			ProjectExistsFunc: func(gitContext common.GitContext) bool {
				return true
			},
			ProjectRepoExistsFunc: func(projectName string) bool {
				return true
			},
			PullFunc: func(gitContext common.GitContext) error {
				return nil
			},
			PushFunc: func(gitContext common.GitContext) error {
				return nil
			},
			StageAndCommitAllFunc: func(gitContext common.GitContext, message string) (string, error) {
				return "my-revision", nil
			},
		},
		credentialReader: &common_mock.CredentialReaderMock{
			GetCredentialsFunc: func(project string) (*common.GitCredentials, error) {
				return &common.GitCredentials{
					User:      "user",
					Token:     "token",
					RemoteURI: "remote-url",
				}, nil
			},
		},
		fileSystem: &common_mock.IFileSystemMock{
			DeleteFileFunc: func(path string) error {
				return nil
			},
			FileExistsFunc: func(path string) bool {
				return true
			},
			MakeDirFunc: func(path string) error {
				return nil
			},
			ReadFileFunc: func(filename string) ([]byte, error) {
				return []byte("file-content"), nil
			},
			WalkPathFunc: func(path string, walkFunc filepath.WalkFunc) error {

				_ = walkFunc(path+"/file1", newFakeFileInfo("file1", false), nil)
				_ = walkFunc(path+"/file2", newFakeFileInfo("file2", false), nil)
				_ = walkFunc(path+"/file3", newFakeFileInfo("file2", false), nil)

				return nil
			},
			WriteBase64EncodedFileFunc: func(path string, content string) error {
				return nil
			},
			WriteFileFunc: func(path string, content []byte) error {
				return nil
			},
		},
	}
}
