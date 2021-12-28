package common

import (
	"github.com/go-git/go-billy/v5/memfs"
	fixtures "github.com/go-git/go-git-fixtures/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	common_mock "github.com/keptn/keptn/resource-service/common/fake"
	"github.com/keptn/keptn/resource-service/common_models"
	"k8s.io/client-go/kubernetes"
	"reflect"
	"testing"
)

func NewTestGit() *common_mock.GogitMock {
	fs, _ := memfs.New().Chroot(".debug/config")
	mem := memory.NewStorage()
	url := fixtures.ByURL("https://github.com/git-fixtures/basic.git").One().DotGit().Root()
	return &common_mock.GogitMock{
		PlainCloneFunc: func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return git.Clone(mem, fs, &git.CloneOptions{URL: url})
		},
		PlainOpenFunc: func(path string) (*git.Repository, error) {
			return git.Open(mem, fs)
		},
	}

}

func TestGit_CheckoutBranch(t *testing.T) {

	tests := []struct {
		name       string
		gitContext common_models.GitContext
		branch     string
		wantErr    bool
	}{
		{
			name: "checkout master branch full ref",
			gitContext: common_models.GitContext{
				Project: "go",
				Credentials: &common_models.GitCredentials{
					User:      "Me",
					Token:     "blabla",
					RemoteURI: "https://github.com/git-fixtures/basic.git"},
			},
			branch: "refs/heads/master",
		},
		{
			name: "checkout master branch",
			gitContext: common_models.GitContext{
				Project: "go",
				Credentials: &common_models.GitCredentials{
					User:      "Me",
					Token:     "blabla",
					RemoteURI: "https://github.com/git-fixtures/basic.git"},
			},
			branch: "master",
		},
		{
			name: "checkout not existing branch",
			gitContext: common_models.GitContext{
				Project: "go",
				Credentials: &common_models.GitCredentials{
					User:      "Me",
					Token:     "blabla",
					RemoteURI: "https://github.com/git-fixtures/basic.git"},
			},
			branch:  "refs/heads/dev",
			wantErr: true,
		},
		{
			name: "checkout existing origin branch",
			gitContext: common_models.GitContext{
				Project: "go",
				Credentials: &common_models.GitCredentials{
					User:      "Me",
					Token:     "blabla",
					RemoteURI: "https://github.com/git-fixtures/basic.git"},
			},
			branch: "refs/remotes/origin/branch",
		},
	}
	g := Git{
		git: NewTestGit(),
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if err := g.CheckoutBranch(tt.gitContext, tt.branch); (err != nil) != tt.wantErr {
				t.Errorf("CheckoutBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGit_CloneRepo(t *testing.T) {
	type args struct {
		gitContext common_models.GitContext
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
		gitContext   common_models.GitContext
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
		gitContext common_models.GitContext
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

	tests := []struct {
		name       string
		gitContext common_models.GitContext
		revision   string
		file       string
		want       []byte
		wantErr    bool
	}{
		{
			name: "get from commitID",
			gitContext: common_models.GitContext{
				Project: "go",
				Credentials: &common_models.GitCredentials{
					User:      "Me",
					Token:     "blabla",
					RemoteURI: "https://github.com/git-fixtures/basic.git"},
			},
			file:     "example.go",
			revision: "918c48b83bd081e863dbe1b80f8998f058cd8294",
			want:     []byte{},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Git{NewTestGit()}
			got, err := g.GetFileRevision(tt.gitContext, tt.revision, tt.file)
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
		gitContext common_models.GitContext
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
		gitContext common_models.GitContext
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
		gitContext common_models.GitContext
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
		gitContext common_models.GitContext
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
		want    *common_models.GitCredentials
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

func Test_ensureRemoteMatchesCredentials(t *testing.T) {
	type args struct {
		repo       *git.Repository
		gitContext common_models.GitContext
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
