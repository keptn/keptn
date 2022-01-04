package common

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/keptn/go-utils/pkg/common/retry"
	"github.com/keptn/keptn/resource-service/common_models"
	utils "github.com/keptn/kubernetes-utils/pkg"
	logger "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Git struct {
	git Gogit
}

func configureGitUser(repository *git.Repository) error {

	config, err := repository.Config()
	config.User.Name = getGitKeptnUser()
	config.User.Email = getGitKeptnEmail()
	if err != nil {
		return fmt.Errorf("could not set git user: %w", err)
	}
	repository.SetConfig(config)
	return nil

}

func getGitKeptnUser() string {
	if keptnUser := os.Getenv(gitKeptnUserEnvVar); keptnUser != "" {
		return keptnUser
	}
	return gitKeptnUserDefault
}

func getGitKeptnEmail() string {
	if keptnEmail := os.Getenv(gitKeptnEmailEnvVar); keptnEmail != "" {
		return keptnEmail
	}
	return gitKeptnEmailDefault
}

func (g Git) CloneRepo(gitContext common_models.GitContext) (bool, error) {
	if (gitContext == common_models.GitContext{}) || (*gitContext.Credentials == common_models.GitCredentials{}) {
		return false, errors.New("Could not clone repo: " + InvalidContextErrorMsg)
	}

	projectPath := GetProjectConfigPath(gitContext.Project)
	if g.ProjectRepoExists(gitContext.Project) {
		// if project exist we do not clone again
		return true, nil
	}
	err := ensureDirectoryExists(projectPath)
	if err != nil {
		return false, err
	}
	clone, err := g.git.PlainClone(projectPath, false,
		&git.CloneOptions{
			URL: gitContext.Credentials.RemoteURI,
			Auth: &http.BasicAuth{
				Username: gitContext.Credentials.User,
				Password: gitContext.Credentials.Token,
			},
		},
	)

	if err != nil {
		if strings.Contains(err.Error(), "empty") {
			// TODO empty remote leads to an error
			clone, err = g.init(gitContext, projectPath)
			if err != nil {
				return false, err
			}
		} else {
			return false, err
		}
	}

	err = configureGitUser(clone)
	if err != nil {
		return false, err
	}

	_, err = clone.Head()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (g Git) init(gitContext common_models.GitContext, projectPath string) (*git.Repository, error) {
	init, err := g.git.PlainInit(projectPath, false)
	if err != nil {
		return nil, err
	}

	_, err = init.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{gitContext.Credentials.RemoteURI},
	})
	os.MkdirAll(projectPath+"/.git", 0700)
	w, err := init.Worktree()
	if err != nil {
		return nil, err
	}
	w.Add("/.git")
	_, err = g.commitAll(gitContext, "init repo")
	if err != nil {
		return nil, err
	}
	return init, nil
}

func (g Git) commitAll(gitContext common_models.GitContext, message string) (string, error) {
	_, w, err := g.getWorkTree(gitContext)
	if err != nil {
		return "", err
	}
	if message == "" {
		message = "commit changes"
	}
	id, err := w.Commit(message,
		&git.CommitOptions{
			All: true,
			Author: &object.Signature{
				Name:  gitKeptnUserDefault,
				Email: gitKeptnEmailDefault,
				When:  time.Now(),
			},
		})
	return id.String(), err
}

func (g Git) StageAndCommitAll(gitContext common_models.GitContext, message string) (string, error) {

	id, err := g.commitAll(gitContext, message)
	if err != nil {
		return "", err
	}
	err = retry.Retry(func() error {
		err = g.Pull(gitContext)
		if err != nil {
			logger.WithError(err).Warn("could not pull")
			return err
		}

		err = g.Push(gitContext)
		if err != nil {
			logger.WithError(err).Warn("could not push")
			return err
		}
		return nil
	}, retry.NumberOfRetries(5), retry.DelayBetweenRetries(1*time.Second))
	if err != nil {
		//TODO : test me
		id, err = utils.ExecuteCommandInDirectory(
			"git", []string{"push", "--set-upstream",
				gitContext.Credentials.RemoteURI},
			GetProjectConfigPath(gitContext.Project),
		)
		return id, err
	}
	return id, err
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
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}
	return nil
}

func (g *Git) Pull(gitContext common_models.GitContext) error {
	if g.ProjectExists(gitContext) {
		_, w, err := g.getWorkTree(gitContext)
		if err != nil {
			return err
		}
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			return nil
		}
		return err
	}
	return errors.New(ProjectDoesNotExistErrorMsg)
}
func (g *Git) GetCurrentRevision(gitContext common_models.GitContext) (string, error) {
	r, _, err := g.getWorkTree(gitContext)
	if err != nil {
		return "", err
	}
	ref, err := r.Head()
	hash := ref.Hash()
	return hash.String(), err
}

func (g *Git) CreateBranch(gitContext common_models.GitContext, branch string, sourceBranch string) error {
	// move head to sourceBranch
	err := g.CheckoutBranch(gitContext, sourceBranch)
	if err != nil {
		return err //errors.New("Could not create branch, source branch does not exist!")
	}
	b := plumbing.NewBranchReferenceName(branch)
	newBranch := &config.Branch{
		Name:   branch,
		Remote: "origin",
		Merge:  b,
	}
	r, _, err := g.getWorkTree(gitContext)
	if err != nil {
		return err
	}
	err = r.CreateBranch(newBranch)
	return err
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
	r, _, err := g.getWorkTree(gitContext)
	config, err := r.Config()
	if err != nil {
		return "", err
	}
	def := config.Init.DefaultBranch
	if def == "" {
		def = "master"
	}
	return def, err
}

func (g *Git) ProjectExists(gitContext common_models.GitContext) bool {
	if g.ProjectRepoExists(gitContext.Project) {
		return true
	}
	// if not, try to clone
	//_, err := os.Stat(path)
	clone, _ := g.CloneRepo(gitContext)
	return clone
}

func (g *Git) ProjectRepoExists(project string) bool {
	path := GetProjectConfigPath(project)
	_, err := os.Stat(path)
	if err == nil {
		// path exists
		_, err = g.git.PlainOpen(path)
		if err == nil {
			return true
		}
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

func ensureDirectoryExists(path string) error {
	if _, err := os.Stat(path); err != nil {

		if err := os.MkdirAll(path, 0700); err != nil {
			return err
		}
	}
	return nil
}
