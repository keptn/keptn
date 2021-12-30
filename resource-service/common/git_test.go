package common

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/keptn/keptn/resource-service/common_models"
	"k8s.io/client-go/kubernetes"
	"reflect"
	"testing"
)

func TestGit_CloneRepo(t *testing.T) {

	tests := []struct {
		name       string
		gitContext common_models.GitContext
		want       bool
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Git{}
			got, err := g.CloneRepo(tt.gitContext)
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

				// add dummy file to check if branch exist
				/*
						f, err := w.Filesystem.Create("fo/fool.go")
						c.Assert(err, IsNil)
						f.Write([]byte(fmt.Sprintf("%s", "foo ciao")))
						f.Close()

						_,  err = w.Add("fo/fool.go")
						c.Assert(err, IsNil)

						_, err = w.Commit("added a file",
							&git.CommitOptions{
							All: true,
							Author: &object.Signature{
									Name:  "Test Create Branch",
									Email: "createBranch@gogit-test.com",
									When:  time.Now(),
								},
							})

						c.Assert(err, IsNil)
						//push to repo
					/*	err = r.Push(&git.PushOptions{
							//Force: true,
							Auth: &http.BasicAuth{
							Username: tt.gitContext.Credentials.User,
							Password: tt.gitContext.Credentials.Token,
						}})
						c.Assert(err, IsNil)*/

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
			if _, err := g.StageAndCommitAll(tt.args.gitContext, tt.args.message); (err != nil) != tt.wantErr {
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
