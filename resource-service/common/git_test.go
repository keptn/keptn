package common

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	common_mock "github.com/keptn/keptn/resource-service/common/fake"
	"k8s.io/client-go/kubernetes"
	"reflect"
	"testing"
)

func TestGit_CheckoutBranch(t *testing.T) {

	tests := []struct {
		name       string
		gitContext GitContext
		branch     string
		wantErr    bool
	}{
		{name: "checkout master branch",
			gitContext: git.context,
			branch:     "master",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedGogit := &common_mock.GogitMock{
				PlainCloneFunc: func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {

				},
				PlainInitFunc: func(path string, isBare bool) (*git.Repository, error) {
					panic("mock out the PlainInit method")
				},
				PlainOpenFunc: func(path string) (*git.Repository, error) {
					panic("mock out the PlainOpen method")
				},
			}
			g := Git{
				git: mockedGogit,
			}
			if err := g.CheckoutBranch(tt.gitContext, tt.branch); (err != nil) != tt.wantErr {
				t.Errorf("CheckoutBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGit_CloneRepo(t *testing.T) {
	type args struct {
		gitContext GitContext
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Git{}
			got, err := g.CloneRepo(tt.args.gitContext)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloneRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CloneRepo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGit_CreateBranch(t *testing.T) {
	type args struct {
		gitContext   GitContext
		branch       string
		sourceBranch string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Git{}
			if err := g.CreateBranch(tt.args.gitContext, tt.args.branch, tt.args.sourceBranch); (err != nil) != tt.wantErr {
				t.Errorf("CreateBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGit_GetDefaultBranch(t *testing.T) {
	type args struct {
		gitContext GitContext
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Git{}
			got, err := g.GetDefaultBranch(tt.args.gitContext)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDefaultBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDefaultBranch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGit_GetFileRevision(t *testing.T) {
	type args struct {
		gitContext GitContext
		revision   string
		file       string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Git{}
			got, err := g.GetFileRevision(tt.args.gitContext, tt.args.revision, tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileRevision() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFileRevision() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGit_ProjectExists(t *testing.T) {
	type args struct {
		gitContext GitContext
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Git{}
			if got := g.ProjectExists(tt.args.gitContext); got != tt.want {
				t.Errorf("ProjectExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGit_ProjectRepoExists(t *testing.T) {
	type args struct {
		project string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Git{}
			if got := g.ProjectRepoExists(tt.args.project); got != tt.want {
				t.Errorf("ProjectRepoExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGit_Pull(t *testing.T) {
	type args struct {
		gitContext GitContext
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Git{}
			if err := g.Pull(tt.args.gitContext); (err != nil) != tt.wantErr {
				t.Errorf("Pull() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGit_Push(t *testing.T) {
	type args struct {
		gitContext GitContext
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Git{}
			if err := g.Push(tt.args.gitContext); (err != nil) != tt.wantErr {
				t.Errorf("Push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGit_StageAndCommitAll(t *testing.T) {
	type args struct {
		gitContext GitContext
		message    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Git{}
			if err := g.StageAndCommitAll(tt.args.gitContext, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("StageAndCommitAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestK8sCredentialReader_GetCredentials(t *testing.T) {
	type args struct {
		project string
	}
	tests := []struct {
		name    string
		args    args
		want    *GitCredentials
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k8 := K8sCredentialReader{}
			got, err := k8.GetCredentials(tt.args.project)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCredentials() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkoutBranch(t *testing.T) {
	type args struct {
		gitContext GitContext
		options    *git.CheckoutOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkoutBranch(tt.args.gitContext, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("checkoutBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ensureRemoteMatchesCredentials(t *testing.T) {
	type args struct {
		repo       *git.Repository
		gitContext GitContext
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ensureRemoteMatchesCredentials(tt.args.repo, tt.args.gitContext); (err != nil) != tt.wantErr {
				t.Errorf("ensureRemoteMatchesCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getK8sClient(t *testing.T) {
	tests := []struct {
		name    string
		want    *kubernetes.Clientset
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getK8sClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("getK8sClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getK8sClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getWorkTree(t *testing.T) {
	type args struct {
		gitContext GitContext
	}
	tests := []struct {
		name    string
		args    args
		want    *git.Repository
		want1   *git.Worktree
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getWorkTree(tt.args.gitContext)
			if (err != nil) != tt.wantErr {
				t.Errorf("getWorkTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getWorkTree() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getWorkTree() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_resolve(t *testing.T) {
	type args struct {
		obj  object.Object
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *object.Blob
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolve(tt.args.obj, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolve() got = %v, want %v", got, tt.want)
			}
		})
	}
}
