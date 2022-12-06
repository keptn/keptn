package common

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport"

	envconfig "github.com/keptn/keptn/resource-service/config"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/keptn/keptn/resource-service/common_models"
	kerrors "github.com/keptn/keptn/resource-service/errors"
	logger "github.com/sirupsen/logrus"
)

const gitHeadFilePath = "/.git/HEAD"

// IGit provides functions to interact with the git repository of a project
//
//go:generate moq -pkg common_mock -skip-ensure -out ./fake/git_mock.go . IGit
type IGit interface {
	ProjectExists(gitContext common_models.GitContext) bool
	ProjectRepoExists(projectName string) bool
	CloneRepo(gitContext common_models.GitContext) (bool, error)
	StageAndCommitAll(gitContext common_models.GitContext, message string) (string, error)
	Push(gitContext common_models.GitContext) error
	Pull(gitContext common_models.GitContext) error
	CreateBranch(gitContext common_models.GitContext, branch string, sourceBranch string) error
	CheckoutBranch(gitContext common_models.GitContext, branch string) error
	GetFileRevision(gitContext common_models.GitContext, revision string, file string) ([]byte, error)
	GetCurrentRevision(gitContext common_models.GitContext) (string, error)
	GetDefaultBranch(gitContext common_models.GitContext) (string, error)
	MigrateProject(gitContext common_models.GitContext, newMetadatacontent []byte) error
	ResetHard(gitContext common_models.GitContext, revision string) error
	MoveToNewUpstream(currentContext common_models.GitContext, newContext common_models.GitContext) error
	CheckUpstreamConnection(gitContext common_models.GitContext) error
}

type Git struct {
	git Gogit
}

func NewGit(git Gogit) *Git {
	return &Git{git: git}
}

func configureGitUser(repository *git.Repository) error {

	c, err := repository.Config()
	c.User.Name = getGitKeptnUser()
	c.User.Email = getGitKeptnEmail()
	if err != nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotSetUser, err)
	}
	repository.SetConfig(c)
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
	if (gitContext.Credentials == nil) || (*gitContext.Credentials == common_models.GitCredentials{}) {
		return false, fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "clone", "project", kerrors.ErrInvalidGitContext)
	}

	projectPath := GetProjectConfigPath(gitContext.Project)
	if g.ProjectRepoExists(gitContext.Project) {
		// if project exist we do not clone again
		return true, nil
	}
	err := ensureDirectoryExists(projectPath)
	if err != nil {
		return false, fmt.Errorf(kerrors.ErrMsgCouldNotCreatePath, projectPath, err)
	}
	clone, err := g.git.PlainClone(gitContext, projectPath, false,
		&git.CloneOptions{
			URL:             gitContext.Credentials.RemoteURL,
			Auth:            gitContext.AuthMethod.GoGitAuth,
			InsecureSkipTLS: retrieveInsecureSkipTLS(gitContext.Credentials),
		},
	)

	if err != nil {
		if kerrors.ErrEmptyRemoteRepository.Is(err) {
			clone, err = g.init(gitContext, projectPath)
			if err != nil {
				return false, fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "init", gitContext.Project, mapError(err))
			}
		} else {
			return false, fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "clone", gitContext.Project, mapError(err))
		}
	}

	err = configureGitUser(clone)
	if err != nil {
		return false, err
	}

	head, err := clone.Head()
	if err != nil {
		return false, fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "clone", gitContext.Project, mapError(err))
	}

	if err = g.fetch(gitContext, clone); err != nil {
		return false, fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "fetch", gitContext.Project, mapError(err))
	}

	if err := g.storeDefaultBranchConfig(gitContext, err, clone, head); err != nil {
		return false, err
	}

	return true, nil
}

func (g Git) storeDefaultBranchConfig(gitContext common_models.GitContext, err error, clone *git.Repository, head *plumbing.Reference) error {
	cfg, err := clone.Config()
	if err != nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "get config", gitContext.Project, mapError(err))
	}

	cfg.Init.DefaultBranch = head.Name().String()

	err = clone.SetConfig(cfg)
	if err != nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "set config", gitContext.Project, mapError(err))
	}
	return nil
}

func (g Git) rewriteDefaultBranch(path string, env envconfig.EnvConfig) error {
	defaultBranch := env.RetrieveDefaultBranchFromEnv()
	if defaultBranch != common_models.GitInitDefaultBranchName {
		logger.Infof("Setting default branch to %s", defaultBranch)
		input, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		output := bytes.Replace(input, []byte(common_models.GitInitDefaultBranchName), []byte(defaultBranch), -1)

		if err = ioutil.WriteFile(path, output, 0700); err != nil {
			return err
		}
	}
	return nil
}

func (g Git) init(gitContext common_models.GitContext, projectPath string) (*git.Repository, error) {
	init, err := g.git.PlainInit(projectPath, false)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			init, err = g.git.PlainOpen(projectPath)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if err := g.rewriteDefaultBranch(projectPath+gitHeadFilePath, envconfig.Global); err != nil {
		return nil, err
	}

	if _, err := init.Remote("origin"); errors.Is(err, git.ErrRemoteNotFound) {
		_, err = init.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{gitContext.Credentials.RemoteURL},
		})
		if err != nil {
			return nil, err
		}
	}

	f, err := os.Create(projectPath + "/metadata.yaml")
	if err != nil {
		return nil, err
	}
	_, err = f.Write([]byte{})
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}

	os.MkdirAll(projectPath+"/.git", 0700)
	w, err := init.Worktree()
	if err != nil {
		return nil, err
	}

	w.Add(projectPath + "/metadata.yaml")
	_, err = w.Commit("init git empty repo",
		&git.CommitOptions{
			All: true,
			Author: &object.Signature{
				Name:  getGitKeptnUser(),
				Email: getGitKeptnEmail(),
				When:  time.Now(),
			},
		})
	if err != nil {
		return nil, err
	}

	err = g.Push(gitContext)
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

	err = w.AddWithOptions(&git.AddOptions{All: true})
	if err != nil {
		return "", err
	}
	id, err := w.Commit(message,
		&git.CommitOptions{
			All: true,
			Author: &object.Signature{
				Name:  getGitKeptnUser(),
				Email: getGitKeptnEmail(),
				When:  time.Now(),
			},
		})
	return id.String(), err
}

func (g Git) StageAndCommitAll(gitContext common_models.GitContext, message string) (string, error) {
	id, err := g.commitAll(gitContext, message)
	if err != nil {
		if err = g.ResetHard(gitContext, "HEAD~0"); err != nil {
			logger.Warnf("Could not reset after commitAll: %v", err)
		} else {
			logger.Warn("Untracked changes were removed")
		}
		return "", fmt.Errorf(kerrors.ErrMsgCouldNotCommit, gitContext.Project, mapError(err))
	}
	rollbackFunc := func() {
		err := g.ResetHard(gitContext, "HEAD~1")
		if err != nil {
			logger.Warnf("Could not reset: %v", err)
		} else {
			logger.Warn("Committed changes were removed")
		}
	}
	err = g.Pull(gitContext)
	if err != nil {
		rollbackFunc()
		return "", err
	}

	err = g.Push(gitContext)
	if err != nil {
		rollbackFunc()
		return "", err
	}

	id, updated, err := g.getCurrentRemoteRevision(gitContext)
	if err != nil {
		return "", fmt.Errorf(kerrors.ErrMsgCouldNotCommit, gitContext.Project, mapError(err))
	}
	if !updated {
		return "", fmt.Errorf(kerrors.ErrMsgCouldNotCommit, gitContext.Project, kerrors.ErrForceNeeded)
	}

	return id, nil
}

func (g Git) Push(gitContext common_models.GitContext) error {
	var err error
	if gitContext.Credentials == nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "push", gitContext.Project, kerrors.ErrCredentialsNotFound)
	}
	repo, _, err := g.getWorkTree(gitContext)
	if err != nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "push", gitContext.Project, mapError(err))
	}
	err = repo.Push(&git.PushOptions{
		RemoteName:      "origin",
		Auth:            gitContext.AuthMethod.GoGitAuth,
		InsecureSkipTLS: retrieveInsecureSkipTLS(gitContext.Credentials),
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "push", gitContext.Project, mapError(err))
	}
	return nil
}

func (g *Git) Pull(gitContext common_models.GitContext) error {
	if !g.ProjectExists(gitContext) {
		return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "pull", gitContext.Project, kerrors.ErrProjectNotFound)
	}

	r, w, err := g.getWorkTree(gitContext)
	if err != nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "pull", gitContext.Project, mapError(err))
	}

	head, err := r.Head()
	if err != nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "pull", gitContext.Project, mapError(err))
	}
	err = w.Pull(&git.PullOptions{
		RemoteName:      "origin",
		Force:           true,
		ReferenceName:   head.Name(),
		Auth:            gitContext.AuthMethod.GoGitAuth,
		InsecureSkipTLS: retrieveInsecureSkipTLS(gitContext.Credentials),
	})
	if err != nil && errors.Is(err, plumbing.ErrReferenceNotFound) {
		// reference not there yet
		err = w.Pull(&git.PullOptions{RemoteName: "origin", Force: true, Auth: gitContext.AuthMethod.GoGitAuth, InsecureSkipTLS: retrieveInsecureSkipTLS(gitContext.Credentials)})
	}
	if err != nil {
		// do not return an error if we are alread< up to date or if the repository is empty
		if errors.Is(err, git.NoErrAlreadyUpToDate) || errors.Is(err, transport.ErrEmptyRemoteRepository) {
			return nil
		}
		return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "pull", gitContext.Project, mapError(err))
	}

	return nil
}

// mapError translates errors that are specific to the go-git library to errors that are understood by the other resource-service components
func mapError(err error) error {
	if errors.Is(err, git.ErrNonFastForwardUpdate) {
		return kerrors.ErrNonFastForwardUpdate
	}
	if errors.Is(err, transport.ErrAuthenticationRequired) {
		return kerrors.ErrAuthenticationRequired
	}
	if errors.Is(err, transport.ErrAuthorizationFailed) {
		return kerrors.ErrAuthorizationFailed
	}
	if errors.Is(err, git.ErrForceNeeded) {
		return kerrors.ErrForceNeeded
	}
	return err
}

func (g *Git) GetCurrentRevision(gitContext common_models.GitContext) (string, error) {
	r, _, err := g.getWorkTree(gitContext)
	if err != nil {
		return "", fmt.Errorf(kerrors.ErrMsgCouldNotGetRevision, gitContext.Project, mapError(err))
	}
	ref, err := r.Head()
	if err != nil {
		return "", fmt.Errorf(kerrors.ErrMsgCouldNotGetRevision, gitContext.Project, mapError(err))
	}
	hash := ref.Hash()
	return hash.String(), nil
}

// returns what is the current commit id of remote and if the remote is up-to-date with the local branch
func (g *Git) getCurrentRemoteRevision(gitContext common_models.GitContext) (string, bool, error) {
	repo, _, err := g.getWorkTree(gitContext)
	if err != nil {
		return "", false, fmt.Errorf(kerrors.ErrMsgCouldNotGetRevision, gitContext.Project, err)
	}

	headRef, err := repo.Head()
	if err != nil {
		return "", false, fmt.Errorf(kerrors.ErrMsgCouldNotGetRevision, gitContext.Project, err)
	}

	// get hash
	branch := headRef.Name().Short()
	revision := plumbing.Revision("origin/" + branch)
	revHash, err := repo.ResolveRevision(revision)

	if err != nil {
		return "", false, fmt.Errorf(kerrors.ErrMsgCouldNotGetRevision, gitContext.Project, err)
	}

	// ... retrieving the commit objects
	revCommit, err := repo.CommitObject(*revHash)
	if err != nil {
		return "", false, fmt.Errorf(kerrors.ErrMsgCouldNotGetRevision, gitContext.Project, err)
	}

	headCommit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return "", false, fmt.Errorf(kerrors.ErrMsgCouldNotGetRevision, gitContext.Project, err)
	}

	//check if latest repo commit is in remote
	isAncestor, err := headCommit.IsAncestor(revCommit)

	if err != nil {
		return "", false, fmt.Errorf(kerrors.ErrMsgCouldNotGetRevision, gitContext.Project, err)
	}
	return revHash.String(), isAncestor, nil
}

func retrieveInsecureSkipTLS(credentials *common_models.GitCredentials) bool {
	if credentials != nil && credentials.HttpsAuth != nil {
		return credentials.HttpsAuth.InsecureSkipTLS
	}
	return false
}

func (g *Git) CreateBranch(gitContext common_models.GitContext, branch string, sourceBranch string) error {
	// move head to sourceBranch
	err := g.CheckoutBranch(gitContext, sourceBranch)
	if err != nil {
		return err
	}
	b := plumbing.NewBranchReferenceName(branch)
	newBranch := &config.Branch{
		Name:   branch,
		Remote: "origin",
		Merge:  b,
	}
	r, w, err := g.getWorkTree(gitContext)
	if err != nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotCreate, branch, gitContext.Project, mapError(err))
	}

	// First try to check out branch
	err = w.Checkout(&git.CheckoutOptions{Create: false, Force: false, Branch: b})
	if err == nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotCreate, branch, gitContext.Project, kerrors.ErrBranchExists)
	}

	if err != nil {
		// got an error  - try to create it
		if err := w.Checkout(&git.CheckoutOptions{Create: true, Force: false, Branch: b}); err != nil {
			return fmt.Errorf(kerrors.ErrMsgCouldNotCreate, branch, gitContext.Project, mapError(err))
		}
	}

	err = r.CreateBranch(newBranch)
	if err != nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotCreate, branch, gitContext.Project, mapError(err))
	}

	return nil
}

func (g *Git) CheckoutBranch(gitContext common_models.GitContext, branch string) error {
	//  short path
	b := plumbing.NewBranchReferenceName(branch)

	//  complete reference path
	if strings.HasPrefix(branch, "refs") {
		b = plumbing.ReferenceName(branch)
	}

	err := g.checkoutBranch(gitContext, &git.CheckoutOptions{
		Branch: b,
		Force:  true,
	})
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return fmt.Errorf(kerrors.ErrMsgCouldNotCheckout, branch, kerrors.ErrReferenceNotFound)
		}
		return fmt.Errorf(kerrors.ErrMsgCouldNotCheckout, branch, mapError(err))
	}
	return nil
}

func (g *Git) checkoutBranch(gitContext common_models.GitContext, options *git.CheckoutOptions) error {
	if g.ProjectExists(gitContext) {
		r, w, err := g.getWorkTree(gitContext)
		if err != nil {
			return err
		}
		if err = g.fetch(gitContext, r); err != nil {
			return err
		}

		return w.Checkout(options)
	}
	return kerrors.ErrProjectNotFound
}

func (g *Git) fetch(gitContext common_models.GitContext, r *git.Repository) error {
	if err := r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{"+refs/*:refs/*"},
		// <src>:<dst>, + update the reference even if it isn’t a fast-forward.
		//// take all branch from remote and put them in the local repo as origin branches and as branches
		//RefSpecs: []config.RefSpec{"+refs/heads/*:refs/remotes/origin/*", "+refs/heads/*:refs/heads/*"},
		Force:           true,
		Auth:            gitContext.AuthMethod.GoGitAuth,
		InsecureSkipTLS: retrieveInsecureSkipTLS(gitContext.Credentials),
	}); err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}
	return nil
}

func (g *Git) GetFileRevision(gitContext common_models.GitContext, revision string, file string) ([]byte, error) {
	path := GetProjectConfigPath(gitContext.Project)
	r, err := g.git.PlainOpen(path)
	if err != nil {
		logger.Debugf("Could not open project %s: %s", file, err)
		return []byte{},
			fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "open", gitContext.Project, err)
	}
	h, err := r.ResolveRevision(plumbing.Revision(revision))
	if err != nil {
		logger.Debugf("Could not resolve revision for %s: %s", revision, err)
		return []byte{},
			fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "retrieve revision in ", gitContext.Project, err)
	}
	if h == nil {
		return []byte{},
			fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "open", gitContext.Project, kerrors.ErrResolvedNilHash)
	}

	obj, err := r.Object(plumbing.CommitObject, *h)

	if err != nil {
		return []byte{},
			fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "retrieve revision in ", gitContext.Project, err)
	}
	if obj == nil {
		return []byte{}, fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "retrieve revision in ", gitContext.Project, kerrors.ErrResolveRevision)
	}
	blob, err := resolve(obj, file)

	if err != nil {
		if errors.Is(err, object.ErrFileNotFound) {
			return []byte{},
				fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "retrieve revision in ", gitContext.Project, kerrors.ErrResourceNotFound)
		}
		return []byte{},
			fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "retrieve revision in ", gitContext.Project, err)
	}

	var re (io.Reader)
	re, err = blob.Reader()

	if err != nil {
		return []byte{},
			fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "retrieve revision in ", gitContext.Project, err)
	}

	return ioutil.ReadAll(re)
}

func (g *Git) GetDefaultBranch(gitContext common_models.GitContext) (string, error) {
	r, _, err := g.getWorkTree(gitContext)
	if err != nil {
		return "", fmt.Errorf(kerrors.ErrMsgCouldNotGetDefBranch, gitContext.Project, err)
	}
	repoConfig, err := r.Config()
	if err != nil {
		return "", fmt.Errorf(kerrors.ErrMsgCouldNotGetDefBranch, gitContext.Project, err)
	}
	def := repoConfig.Init.DefaultBranch
	if def == "" {
		head, err := r.Head()
		if err != nil {
			return "", fmt.Errorf(kerrors.ErrMsgCouldNotGetDefBranch, gitContext.Project, err)
		}
		return string(head.Name()), nil
	}
	return def, err
}

func (g *Git) ProjectExists(gitContext common_models.GitContext) bool {
	if g.ProjectRepoExists(gitContext.Project) {
		return true
	}
	clone, err := g.CloneRepo(gitContext)
	if err != nil {
		logger.Errorf("Could not check for project availability: %v", err)
	}
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

func (g *Git) MoveToNewUpstream(currentContext common_models.GitContext, newContext common_models.GitContext) error {
	if err := g.Pull(currentContext); err != nil {
		return err
	}

	currentRepo, currentRepoWorktree, err := g.getWorkTree(currentContext)
	if err != nil {
		return err
	}

	tmpOrigin := "tmp-origin"

	// first, make sure the temporary remote is gone before creating it
	if err := currentRepo.DeleteRemote(tmpOrigin); err != nil && err != git.ErrRemoteNotFound {
		return mapError(err)
	}

	if _, err := currentRepo.CreateRemote(&config.RemoteConfig{
		Name: tmpOrigin,
		URLs: []string{newContext.Credentials.RemoteURL},
	}); err != nil {
		return mapError(err)
	}

	if err := g.fetch(currentContext, currentRepo); err != nil {
		return err
	}
	branches, err := currentRepo.Branches()
	err = branches.ForEach(func(branch *plumbing.Reference) error {
		err := currentRepoWorktree.Checkout(&git.CheckoutOptions{Branch: branch.Name()})
		if err != nil {
			return err
		}

		err = currentRepo.Push(&git.PushOptions{
			RemoteName:      tmpOrigin,
			Auth:            newContext.AuthMethod.GoGitAuth,
			Force:           false,
			InsecureSkipTLS: retrieveInsecureSkipTLS(newContext.Credentials),
		})
		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "push", newContext.Project, mapError(err))
		}
		return nil
	})
	if err != nil {
		return mapError(err)
	}

	if err := ensureRemoteMatchesCredentials(currentRepo, newContext); err != nil {
		return mapError(err)
	}

	if err := currentRepo.DeleteRemote(tmpOrigin); err != nil {
		return mapError(err)
	}

	return nil
}

func (g *Git) CheckUpstreamConnection(gitContext common_models.GitContext) error {
	if err := g.Pull(gitContext); err != nil {
		return mapError(err)
	}
	return nil
}

func (g *Git) MigrateProject(gitContext common_models.GitContext, newMetadataContent []byte) error {
	if err := g.Pull(gitContext); err != nil {
		return err
	}

	tmpGitContext := gitContext
	tmpGitContext.Project = "_keptn-tmp_" + gitContext.Project

	tmpProjectPath := GetProjectConfigPath(tmpGitContext.Project)
	projectPath := GetProjectConfigPath(gitContext.Project)

	defaultBranch, err := g.GetDefaultBranch(gitContext)
	if err != nil {
		return err
	}

	if _, err := g.CloneRepo(tmpGitContext); err != nil {
		return err
	}

	// check out branches of the tmp remote and store the content in the master branch of the repo
	oldRepo, oldRepoWorktree, err := g.getWorkTree(tmpGitContext)
	if err != nil {
		return err
	}

	if err := g.fetch(tmpGitContext, oldRepo); err != nil {
		return err
	}
	branches, err := oldRepo.Branches()
	err = branches.ForEach(func(branch *plumbing.Reference) error {
		if branch.Name().Short() != defaultBranch {
			return g.migrateBranch(branch, oldRepoWorktree, projectPath, tmpProjectPath)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(GetProjectMetadataFilePath(gitContext.Project), newMetadataContent, os.ModePerm); err != nil {
		return err
	}

	_, err = g.StageAndCommitAll(gitContext, "migrated project structure")
	if err != nil {
		return err
	}

	if err := os.RemoveAll(tmpProjectPath); err != nil {
		return err
	}

	return nil
}

func (g *Git) migrateBranch(branch *plumbing.Reference, oldRepoWorktree *git.Worktree, projectPath string, tmpProjectPath string) error {
	err := oldRepoWorktree.Checkout(&git.CheckoutOptions{Branch: branch.Name()})
	if err != nil {
		return err
	}

	err = ensureDirectoryExists(projectPath + "/" + StageDirectoryName)
	if err != nil {
		return err
	}

	err = ensureDirectoryExists(projectPath + "/" + StageDirectoryName + "/" + branch.Name().Short())
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(tmpProjectPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.Name() == ".git" {
			continue
		}
		err := os.Rename(tmpProjectPath+"/"+file.Name(), projectPath+"/"+StageDirectoryName+"/"+branch.Name().Short()+"/"+file.Name())
		if err != nil {
			return err
		}
	}
	err = oldRepoWorktree.Reset(&git.ResetOptions{Mode: git.HardReset})
	if err != nil {
		return err
	}

	return nil
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

func (g Git) ResetHard(gitContext common_models.GitContext, rev string) error {
	r, w, err := g.getWorkTree(gitContext)
	if err != nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "reset", gitContext.Project, err)
	}
	revision, err := r.ResolveRevision(plumbing.Revision(rev))
	if err != nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "reset", gitContext.Project, err)
	}
	return w.Reset(&git.ResetOptions{
		Commit: *revision,
		Mode:   git.HardReset,
	})
}

func ensureRemoteMatchesCredentials(repo *git.Repository, gitContext common_models.GitContext) error {
	remote, err := repo.Remote("origin")
	if err != nil {
		return err
	}
	if remote.Config().URLs[0] != gitContext.Credentials.RemoteURL {
		err := repo.DeleteRemote("origin")
		if err != nil {
			return err
		}
		_, err = repo.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{gitContext.Credentials.RemoteURL},
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
			logger.Debugf("Could not resolve commit for path %s: %v", path, err)
			return nil, err
		}
		return resolve(t, path)
	case *object.Tag:
		target, err := o.Object()
		if err != nil {
			logger.Debugf("Could not resolve tag for path %s: %v", path, err)
			return nil, err
		}
		return resolve(target, path)
	case *object.Tree:
		file, err := o.File(path)
		if err != nil {
			logger.Debugf("Could not resolve file for path %s: %v", path, err)
			return nil, err
		}
		return &file.Blob, nil
	case *object.Blob:
		return o, nil
	default:
		logger.Debugf("Could not resolve unsupported object for path: %s", path)
		return nil, object.ErrUnsupportedObject
	}
}
