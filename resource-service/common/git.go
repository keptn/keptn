package common

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/keptn/keptn/resource-service/common_models"
	"io"
	"io/ioutil"
	"strings"
)

type Git struct {
	git Gogit
}

func NewGit(g Gogit) Git {
	return Git{
		git: g,
	}
}

func (g Git) CloneRepo(gitContext common_models.GitContext) (bool, error) {
	projectPath := GetProjectConfigPath(gitContext.Project)
	if g.ProjectRepoExists(gitContext.Project) {
		// if project exist we do not clone again
		return true, nil
	}

	clone, err := g.git.PlainClone(projectPath, false,
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
			init, err := g.git.PlainInit(projectPath, false)
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

func (g Git) StageAndCommitAll(gitContext common_models.GitContext, message string) error {
	panic("implement me")
}

func (g Git) Push(gitContext common_models.GitContext) error {
	var err error
	if gitContext.Credentials == nil {
		return errors.New("Could not push, invalid credentials")
	}
	repo, _, err := g.getWorkTree(gitContext)
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

func (g *Git) Pull(gitContext common_models.GitContext) error {
	panic("implement me")
}

func (g *Git) CreateBranch(gitContext common_models.GitContext, branch string, sourceBranch string) error {
	// move head to sourceBranch
	g.CheckoutBranch(gitContext, sourceBranch)

	//create new branch
	return g.checkoutBranch(gitContext, &git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branch),
		Create: true,
	})
}

func (g *Git) CheckoutBranch(gitContext common_models.GitContext, branch string) error {
	//  short path
	b := plumbing.NewBranchReferenceName(branch)

	//  complete reference path
	if strings.HasPrefix(branch, "refs") {
		b = plumbing.ReferenceName(branch)
	}

	return g.checkoutBranch(gitContext, &git.CheckoutOptions{
		Branch: b,
	})
}

func (g *Git) checkoutBranch(gitContext common_models.GitContext, options *git.CheckoutOptions) error {
	if g.ProjectExists(gitContext) {
		_, w, err := g.getWorkTree(gitContext)
		if err != nil {
			return err
		}
		return w.Checkout(options)
	}
	return errors.New("Could not find project")
}

func (g *Git) GetFileRevision(gitContext common_models.GitContext, revision string, file string) ([]byte, error) {
	path := GetProjectConfigPath(gitContext.Project)
	r, err := g.git.PlainOpen(path)
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

func (g *Git) GetDefaultBranch(gitContext common_models.GitContext) (string, error) {
	//checkoutBranch(gitContext, &git.CheckoutOptions{Branch: masterBranch})
	return "", nil
}

func (g *Git) ProjectExists(gitContext common_models.GitContext) bool {
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

func (g *Git) ProjectRepoExists(project string) bool {
	_, err := g.git.PlainOpen(GetProjectConfigPath(project))
	if err == nil {
		return true
	}
	return false
}

func (g *Git) getWorkTree(gitContext common_models.GitContext) (*git.Repository, *git.Worktree, error) {
	projectConfigPath := GetProjectConfigPath(gitContext.Project)
	// check if we already have a repository
	repo, err := g.git.PlainOpen(projectConfigPath)
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

func ensureRemoteMatchesCredentials(repo *git.Repository, gitContext common_models.GitContext) error {
	remote, err := repo.Remote("origin")
	if err != nil {
		return err
	}
	if remote.Config().URLs[0] != gitContext.Credentials.RemoteURI {
		err := repo.DeleteRemote("origin")
		if err != nil {
			return err
		}
		_, err = repo.CreateRemote(&config.RemoteConfig{
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
