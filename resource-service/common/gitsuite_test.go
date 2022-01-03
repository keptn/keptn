package common

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	common_mock "github.com/keptn/keptn/resource-service/common/fake"
	"github.com/keptn/keptn/resource-service/common_models"
	config2 "github.com/keptn/keptn/resource-service/config"
	. "gopkg.in/check.v1"
	"os"
	"reflect"
	"testing"
	"time"
)

func Test(t *testing.T) { TestingT(t) }

type BaseSuite struct {
	//Suite      fixtures.Suite
	Repository *git.Repository
	url        string
}

var _ = Suite(&BaseSuite{})

func (s *BaseSuite) SetUpSuite(c *C) {
	// init fixture repo
	//s.Suite.SetUpSuite(c)
	s.buildBasicRepository(c)
}

func (s *BaseSuite) TearDownSuite(c *C) {
	//s.Suite.TearDownSuite(c)
	//err := os.RemoveAll("./debug")
	//c.Assert(err, IsNil)
}

func (s *BaseSuite) SetUpTest(c *C) {
	s.SetUpSuite(c)
}

func (s *BaseSuite) buildBasicRepository(c *C) {
	err := os.RemoveAll("./debug")
	c.Assert(err, IsNil)
	//url := fixtures.ByURL("https://github.com/git-fixtures/basic.git").One().DotGit().Root()
	s.url = config2.ConfigDir + "/remote"

	// make a local remote
	_, err = git.PlainClone(s.url, true, &git.CloneOptions{URL: "https://github.com/git-fixtures/basic.git"})
	c.Assert(err, IsNil)

	// make local git repo
	//fs, err := memfs.New().Chroot(config2.ConfigDir + "/sockshop")
	//s.Repository, err = git.Clone(memory.NewStorage(), fs, &git.CloneOptions{URL: s.url})
	s.Repository, err = git.PlainClone(config2.ConfigDir+"/sockshop", false, &git.CloneOptions{URL: s.url})
	c.Assert(err, IsNil)
}

func (s *BaseSuite) TestGit_GetDefaultBranch(c *C) {

	tests := []struct {
		name       string
		gitContext common_models.GitContext
		want       string
		wantErr    bool
	}{
		{
			name:       "simple master",
			gitContext: s.NewGitContext(),
			want:       "master",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		g := Git{GogitReal{}}
		conf, err := s.Repository.Config()
		c.Assert(err, IsNil)
		conf.Init.DefaultBranch = tt.want
		s.Repository.SetConfig(conf)
		got, err := g.GetDefaultBranch(tt.gitContext)
		if (err != nil) != tt.wantErr {
			c.Errorf("GetDefaultBranch() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if got != tt.want {
			c.Errorf("GetDefaultBranch() got = %v, want %v", got, tt.want)
		}

	}
}

func (s *BaseSuite) TestGit_Pull(c *C) {

	tests := []struct {
		name       string
		gitContext common_models.GitContext
		expected   string
		wantErr    bool
		err        error
	}{
		{
			name:       "retrieve already uptodate sockshop",
			gitContext: s.NewGitContext(),
			wantErr:    false,
			expected: "[core]\n" + "\tbare = false\n" +
				"[remote \"origin\"]\n" +
				"\turl = ./debug/config/remote\n" +
				"\tfetch = +refs/heads/*:refs/remotes/origin/*\n" +
				"[branch \"master\"]\n" +
				"\tremote = origin\n" +
				"\tmerge = refs/heads/master\n",
		},
		{
			name: "retrieve from unexisting project",
			gitContext: common_models.GitContext{
				Project: "mine",
				Credentials: &common_models.GitCredentials{
					User:      "ssss",
					Token:     "bjh",
					RemoteURI: s.url},
			},
			wantErr: false,
			expected: "[core]\n" +
				"\tbare = false\n" +
				"[remote \"origin\"]\n" +
				"\turl = ./debug/config/remote\n" +
				"\tfetch = +refs/heads/*:refs/remotes/origin/*\n" +
				"[branch \"master\"]\n" +
				"\tremote = origin\n" +
				"\tmerge = refs/heads/master\n" +
				"[user]\n" +
				"\tname = keptn\n" +
				"\temail = keptn@keptn.sh\n",
		},
		{
			name: "retrieve from unexisting url",
			gitContext: common_models.GitContext{
				Project: "mine",
				Credentials: &common_models.GitCredentials{
					User:      "ssss",
					Token:     "bjh",
					RemoteURI: "jibberish"},
			},

			wantErr: true,
		},
	}

	for _, tt := range tests {
		c.Logf("Test %s", tt.name)
		g := Git{GogitReal{}}
		if err := g.Pull(tt.gitContext); (err != nil) != tt.wantErr {
			c.Errorf("Pull() error = %v, wantErr %v", err, tt.wantErr)
		}
		if !tt.wantErr {
			b, err := os.ReadFile(GetProjectConfigPath(tt.gitContext.Project + "/.git/config"))
			c.Assert(err, IsNil)
			c.Assert(string(b), Equals, tt.expected)
		}

	}
}

func (s *BaseSuite) Test_resolve(c *C) {

	tests := []struct {
		name    string
		obj     object.Object
		path    string
		want    *object.Blob
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {

		got, err := resolve(tt.obj, tt.path)
		if (err != nil) != tt.wantErr {
			c.Errorf("resolve() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			c.Errorf("resolve() got = %v, want %v", got, tt.want)
		}
	}
}

func (s *BaseSuite) TestGit_CloneRepo(c *C) {

	tests := []struct {
		name       string
		gitContext common_models.GitContext
		git        Gogit
		want       bool
		wantErr    bool
	}{
		{
			name: "clone sockshop from remote",
			git: &common_mock.GogitMock{
				PlainOpenFunc: func(path string) (*git.Repository, error) {
					return s.Repository, nil
				},
			},
			gitContext: s.NewGitContext(),
			wantErr:    false,
			want:       true,
		},
		{
			name: "clone existing sockshop",
			git: &common_mock.GogitMock{
				PlainOpenFunc: func(path string) (*git.Repository, error) {
					return s.Repository, nil
				},
			},
			gitContext: s.NewGitContext(),
			wantErr:    false,
			want:       true,
		},
		{
			name:       "empty context",
			gitContext: common_models.GitContext{},
			git:        GogitReal{},
			wantErr:    true,
			want:       false,
		},
		{ // TODO: do we worry here if url is not valid or while saving it?
			// go git seems to try to parse this wrong url
			name: "wrong url context",
			gitContext: common_models.GitContext{
				Project: "sockshop",
				Credentials: &common_models.GitCredentials{
					User:      "Me",
					Token:     "blabla",
					RemoteURI: "http//wrongurl"},
			},
			git: &common_mock.GogitMock{
				PlainCloneFunc: func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
					return nil, errors.New("auth error")
				},
				PlainInitFunc: func(path string, isBare bool) (*git.Repository, error) {
					return nil, errors.New("not exists")
				},
				PlainOpenFunc: func(path string) (*git.Repository, error) {
					return nil, errors.New("not exists")
				},
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "Wrong credential",
			gitContext: common_models.GitContext{
				Project: "so",
				Credentials: &common_models.GitCredentials{
					User:      "ssss",
					Token:     "bjh",
					RemoteURI: "https://github.com/git-fixtures/basic.git"},
			},
			git: &common_mock.GogitMock{
				PlainCloneFunc: func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
					return nil, errors.New("auth error")
				},
				PlainInitFunc: func(path string, isBare bool) (*git.Repository, error) {
					return nil, errors.New("not exists")
				},
				PlainOpenFunc: func(path string) (*git.Repository, error) {
					return nil, errors.New("not exists")
				},
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		c.Log("Test ", tt.name)
		g := Git{tt.git}
		got, err := g.CloneRepo(tt.gitContext)
		if (err != nil) != tt.wantErr {
			c.Errorf("CloneRepo() error = %v, wantErr %v", err, tt.wantErr)

		}
		if got != tt.want {
			c.Errorf("CloneRepo() got = %v, want %v", got, tt.want)
		}

	}
}

func (s *BaseSuite) TestGit_CreateBranch(c *C) {

	tests := []struct {
		name         string
		gitContext   common_models.GitContext
		branch       string
		sourceBranch string
		wantErr      bool
		error        error
	}{
		{
			name:         "simple branch from master",
			gitContext:   s.NewGitContext(),
			branch:       "dev",
			sourceBranch: "master",
			wantErr:      false,
			error:        nil,
		},
		{
			name:         "add existing",
			gitContext:   s.NewGitContext(),
			branch:       "dev",
			sourceBranch: "master",
			wantErr:      true,
			error:        errors.New("branch already exists"),
		},
		{
			name:         "illegal add to non existing branch",
			gitContext:   s.NewGitContext(),
			branch:       "dev",
			sourceBranch: "refs/heads/branch",
			wantErr:      true,
			error:        errors.New("reference not found"),
		},
	}
	r := s.Repository
	g := Git{
		s.NewTestGit(),
	}

	expected := []byte("[core]\n\tbare = false\n[remote \"origin\"]\n\turl = " +
		"./debug/config/remote\n\tfetch = +refs/heads/*:refs/remotes/origin/*\n[branch \"master\"]\n" +
		"\tremote = origin\n\tmerge = refs/heads/master\n[branch \"dev\"]\n" +
		"\tremote = origin\n\tmerge = refs/heads/dev\n")

	for _, tt := range tests {
		c.Logf("Test: %s", tt.name)

		err := g.CreateBranch(tt.gitContext, tt.branch, tt.sourceBranch)

		if (err != nil) && tt.wantErr {
			c.Assert(err.Error(), Equals, tt.error.Error())
			continue
		}
		if err != nil {
			c.Errorf("CreateBranch() error = %v, wantErr %v", err, tt.wantErr)
		}

		// check git config files
		cfg, err := r.Config()
		c.Assert(err, IsNil)
		marshaled, err := cfg.Marshal()
		c.Assert(err, IsNil)
		c.Assert(string(expected), Equals, string(marshaled))
	}

}

func (s *BaseSuite) TestGit_CheckoutBranch(c *C) {

	tests := []struct {
		name       string
		gitContext common_models.GitContext
		branch     string
		wantErr    bool
	}{
		{
			name:       "checkout master branch full ref",
			gitContext: s.NewGitContext(),
			branch:     "refs/heads/master",
		},
		{
			name:       "checkout master branch",
			gitContext: s.NewGitContext(),
			branch:     "master",
		},
		{
			name:       "checkout not existing branch",
			gitContext: s.NewGitContext(),
			branch:     "refs/heads/dev",
			wantErr:    true,
		},
	}
	g := Git{s.NewTestGit()}
	for _, tt := range tests {
		c.Log("Test: ", tt.name)
		if err := g.CheckoutBranch(tt.gitContext, tt.branch); (err != nil) != tt.wantErr {
			c.Errorf("CheckoutBranch() error = %v, wantErr %v", err, tt.wantErr)
		}

	}
}

func (s *BaseSuite) TestGit_GetFileRevision(c *C) {

	tests := []struct {
		name       string
		gitContext common_models.GitContext
		file       string
		content    string
		wantErr    bool
	}{
		{
			name:       "get from commitID",
			gitContext: s.NewGitContext(),
			file:       "foo/example.go",
			content:    "ciao",
			wantErr:    false,
		},
	}
	for _, tt := range tests {

		g := Git{s.NewTestGit()}
		id := s.commitAndPush(tt.file, tt.content, c)
		got, err := g.GetFileRevision(tt.gitContext, id.String(), tt.file)
		if (err != nil) != tt.wantErr {
			c.Errorf("GetFileRevision() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		b := []byte(fmt.Sprintf("%s", tt.content))
		if !reflect.DeepEqual(got, b) {
			c.Errorf("GetFileRevision() got = %v, want %v", got, b)
		}

	}
}

func (s *BaseSuite) TestGit_ProjectRepoExists(c *C) {

	tests := []struct {
		name    string
		project string
		want    bool
	}{
		{
			name:    "project exists",
			project: "sockshop",
			want:    true,
		},
		{
			name:    "project does not exists",
			project: "whatever",
			want:    false,
		},
	}
	for _, tt := range tests {
		if tt.want {
			os.Mkdir(GetProjectConfigPath(tt.project), os.ModePerm)
			git.PlainInit(GetProjectConfigPath(tt.project), false)
		}
		g := Git{GogitReal{}}
		if got := g.ProjectRepoExists(tt.project); got != tt.want {
			c.Errorf("ProjectRepoExists() = %v, want %v", got, tt.want)
		}

	}
}

func Test_getGitKeptnUser(t *testing.T) {
	tests := []struct {
		name        string
		envVarValue string
		want        string
	}{
		{
			name:        "default value",
			envVarValue: "",
			want:        gitKeptnUserDefault,
		},
		{
			name:        "env var value",
			envVarValue: "my-user",
			want:        "my-user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv(gitKeptnUserEnvVar, tt.envVarValue)
			if got := getGitKeptnUser(); got != tt.want {
				t.Errorf("getGitKeptnUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *BaseSuite) Test_getGitKeptnEmail(c *C) {
	tests := []struct {
		name        string
		envVarValue string
		want        string
	}{
		{
			name:        "default value",
			envVarValue: "",
			want:        gitKeptnEmailDefault,
		},
		{
			name:        "env var value",
			envVarValue: "my-user@keptn.sh",
			want:        "my-user@keptn.sh",
		},
	}
	for _, tt := range tests {
		_ = os.Setenv(gitKeptnEmailEnvVar, tt.envVarValue)
		if got := getGitKeptnEmail(); got != tt.want {
			c.Errorf("getGitKeptnEmail() = %v, want %v", got, tt.want)
		}
	}

}

func (s *BaseSuite) NewGitContext() common_models.GitContext {
	return common_models.GitContext{
		Project: "sockshop",
		Credentials: &common_models.GitCredentials{
			User:      "Me",
			Token:     "blabla",
			RemoteURI: s.url,
		},
	}
}

func (s *BaseSuite) NewTestGit() *common_mock.GogitMock {

	return &common_mock.GogitMock{
		PlainCloneFunc: func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return s.Repository, nil
		},
		PlainInitFunc: nil,
		PlainOpenFunc: func(path string) (*git.Repository, error) {
			return s.Repository, nil //git.PlainOpen(path)
		},
	}
}

func (s *BaseSuite) commitAndPush(file string, content string, c *C) plumbing.Hash {
	r := s.Repository
	w, err := r.Worktree()
	f, err := w.Filesystem.Create(file)
	c.Assert(err, IsNil)
	f.Write([]byte(fmt.Sprintf("%s", content)))
	f.Close()

	_, err = w.Add(file)
	c.Assert(err, IsNil)

	id, err := w.Commit("added a file",
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
	err = r.Push(&git.PushOptions{
		//Force: true,
		Auth: &http.BasicAuth{
			Username: "whatever",
			Password: "whatever",
		}})
	c.Assert(err, IsNil)

	return id
}
