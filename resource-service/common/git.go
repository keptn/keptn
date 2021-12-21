package common

import (
	"errors"
	"github.com/go-git/go-git/v5"
	config2 "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"io"
	"io/ioutil"
	"strings"
)

const masterBranch = "master"

// GitCredentials contains git credentials info
type GitCredentials struct {
	User      string `json:"user,omitempty"`
	Token     string `json:"token,omitempty"`
	RemoteURI string `json:"remoteURI,omitempty"`
}

type GitContext struct {
	Project     string //TODO would it make sense to store branch ? no
	Credentials *GitCredentials
}

// IGit provides functions to interact with the git repository of a project
//go:generate moq -pkg common_mock -skip-ensure -out ./fake/git_mock.go . IGit
type IGit interface {
	ProjectExists(gitContext GitContext) bool
	CloneRepo(gitContext GitContext) (bool, error)
	StageAndCommitAll(gitContext GitContext, message string) error
	Push(gitContext GitContext) error
	Pull(gitContext GitContext) error
	CreateBranch(gitContext GitContext, branch string, sourceBranch string) error
	CheckoutBranch(gitContext GitContext, branch string) error
	GetFileRevision(gitContext GitContext, revision string, file string) ([]byte, error)
	GetDefaultBranch(gitContext GitContext) (string, error)
}

type Git struct{}

func (g Git) CloneRepo(gitContext GitContext) (bool, error) {
	projectPath := GetProjectConfigPath(gitContext.Project)
	if g.ProjectRepoExists(gitContext.Project) {
		// if project exist we do not clone again
		return true, nil
	}

	clone, err := git.PlainClone(projectPath, false,
		&git.CloneOptions{
			URL: gitContext.Credentials.RemoteURI,
			Auth: &http.BasicAuth{
				Username: gitContext.Credentials.User,
				Password: gitContext.Credentials.Token,
			},
			Depth: 1,
		},
	)

	if err != nil {
		if strings.Contains(err.Error(), "empty") {
			// TODO empty remote leads to an error
			init, err := git.PlainInit(projectPath, false)
			if err != nil {
				return false, err
			}
			clone = init
		} else {
			return false, err
		}
	}

	_, err = clone.Head()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (g Git) StageAndCommitAll(gitContext GitContext, message string) error {
	panic("implement me")
}

func (g Git) Push(gitContext GitContext) error {
	var err error
	if gitContext.Credentials == nil {
		return errors.New("Could not push, invalid credentials")
	}
	repo, _, err := getWorkTree(gitContext)
	if err != nil {
		return err
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: gitContext.Credentials.User,
			Password: gitContext.Credentials.Token,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (g Git) Pull(gitContext GitContext) error {
	panic("implement me")
}

func (g Git) CreateBranch(gitContext GitContext, branch string, sourceBranch string) error {
	// move head to sourceBranch
	g.CheckoutBranch(gitContext, sourceBranch)

	//create new branch
	return checkoutBranch(gitContext, &git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branch),
		Create: true,
	})
}

func (g Git) CheckoutBranch(gitContext GitContext, branch string) error {

	return checkoutBranch(gitContext, &git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branch),
	})
}

func checkoutBranch(gitContext GitContext, options *git.CheckoutOptions) error {
	_, w, err := getWorkTree(gitContext)
	if err != nil {
		return err
	}
	err = w.Checkout(options)
	return err
}

func (g Git) GetFileRevision(gitContext GitContext, revision string, file string) ([]byte, error) {
	path := GetProjectConfigPath(gitContext.Project)
	r, err := git.PlainOpen(path)
	if err != nil {
		return []byte{}, err
	}
	h, err := r.ResolveRevision(plumbing.Revision(revision))
	if h == nil {
		return []byte{}, errors.New("resolved nil hash " + revision)
	}

	obj, err := r.Object(plumbing.CommitObject, *h)

	if err != nil {
		return []byte{}, err
	}
	if obj == nil {
		return []byte{}, errors.New("not found")
	}
	blob, err := resolve(obj, file)

	if err != nil {
		return []byte{}, err
	}

	var re (io.Reader)
	re, err = blob.Reader()

	return ioutil.ReadAll(re)
}

func (g Git) GetDefaultBranch(gitContext GitContext) (string, error) {
	//checkoutBranch(gitContext, &git.CheckoutOptions{Branch: masterBranch})
	return "", nil
}

func (g *Git) ProjectExists(gitContext GitContext) bool {
	if g.ProjectRepoExists(gitContext.Project) {
		return true
	}
	// if not, try to clone
	_, err := g.CloneRepo(gitContext)
	if err != nil {
		return false
	}
	return true
}

func (g Git) ProjectRepoExists(project string) bool {
	_, err := git.PlainOpen(GetProjectConfigPath(project))
	if err == nil {
		return true
	}
	return false
}

func getWorkTree(gitContext GitContext) (*git.Repository, *git.Worktree, error) {
	projectConfigPath := GetProjectConfigPath(gitContext.Project)
	// check if we already have a repository
	repo, err := git.PlainOpen(projectConfigPath)
	if err != nil {
		return nil, nil, err
	}

	// check if remote matches with the credentials
	err = ensureRemoteMatchesCredentials(repo, gitContext)
	if err != nil {
		return nil, nil, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, nil, err
	}
	return repo, worktree, nil
}

func ensureRemoteMatchesCredentials(repo *git.Repository, gitContext GitContext) error {
	remote, err := repo.Remote("origin")
	if err != nil {
		return err
	}
	if remote.Config().URLs[0] != gitContext.Credentials.RemoteURI {
		err := repo.DeleteRemote("origin")
		if err != nil {
			return err
		}
		_, err = repo.CreateRemote(&config2.RemoteConfig{
			Name: "origin",
			URLs: []string{gitContext.Credentials.RemoteURI},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func resolve(obj object.Object, path string) (*object.Blob, error) {
	switch o := obj.(type) {
	case *object.Commit:
		t, err := o.Tree()
		if err != nil {
			return nil, err
		}
		return resolve(t, path)
	case *object.Tag:
		target, err := o.Object()
		if err != nil {
			return nil, err
		}
		return resolve(target, path)
	case *object.Tree:
		file, err := o.File(path)
		if err != nil {
			return nil, err
		}
		return &file.Blob, nil
	case *object.Blob:
		return o, nil
	default:
		return nil, object.ErrUnsupportedObject
	}
}
