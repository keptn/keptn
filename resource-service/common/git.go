package common

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	nethttp "net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/keptn/keptn/resource-service/common_models"
	kerrors "github.com/keptn/keptn/resource-service/errors"
	logger "github.com/sirupsen/logrus"
	ssh2 "golang.org/x/crypto/ssh"
)

// IGit provides functions to interact with the git repository of a project
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
	if (gitContext == common_models.GitContext{}) || (*gitContext.Credentials == common_models.GitCredentials{}) {
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

	gitCloneOptions := &git.CloneOptions{
		URL: gitContext.Credentials.RemoteURI,
	}
	logger.Info("som v clonerepo a idem newpublickeys a privatekey je ")
	if gitContext.Credentials.PrivateKey != "" {
		publicKey, err := ssh.NewPublicKeys("git", []byte(gitContext.Credentials.PrivateKey), "<password>")
		if err != nil {
			logger.Info("publickey nevyslo, error je ", err)
			return false, err
		}
		publicKey.HostKeyCallback = ssh2.InsecureIgnoreHostKey()
		logger.Info("publickey vyslo")
		gitCloneOptions.Auth = publicKey
	} else if gitContext.Credentials.Token != "" {
		logger.Info("!!!!!som v token vetve")
		if gitContext.Credentials.GitProxyUrl != "" {
			logger.Info("!!!!!som v proxy vetve")
			customClient := &nethttp.Client{
				Transport: &nethttp.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					Proxy: nethttp.ProxyURL(&url.URL{
						Scheme: gitContext.Credentials.GitProxyScheme,
						//User:   url.UserPassword(gitContext.Credentials.GitProxyUser, gitContext.Credentials.GitProxyPassword),
						Host: gitContext.Credentials.GitProxyUrl,
					}),
				},

				// 15 second timeout
				Timeout: 15 * time.Second,

				// don't follow redirect
				CheckRedirect: func(req *nethttp.Request, via []*nethttp.Request) error {
					return nethttp.ErrUseLastResponse
				},
			}

			newclient := http.NewClient(customClient)
			if newclient == nil {
				logger.Info("!!!!!newclient zlyhal")
			}

			logger.Info("!!!!!idem instalovat protokol")
			client.InstallProtocol("https", newclient)
			logger.Info("!!!!!nainstaloval som protokol")

			gitCloneOptions.Auth = &http.BasicAuth{
				Username: gitContext.Credentials.User,
				Password: gitContext.Credentials.Token,
			}
		} else {
			logger.Info("!!!!!som v jednoduchej token vetve vedla proxy")
			gitCloneOptions.Auth = &http.BasicAuth{
				Username: gitContext.Credentials.User,
				Password: gitContext.Credentials.Token,
			}
		}
	}

	clone, err := g.git.PlainClone(projectPath, false, gitCloneOptions)

	if err != nil {
		logger.Info("plainclone nevyslo")
		if kerrors.ErrEmptyRemoteRepository.Is(err) {
			logger.Info("plainclone nevyslo222")
			clone, err = g.init(gitContext, projectPath)
			if err != nil {
				logger.Info("plainclone1 nevyslo, error je ", err)
				return false, fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "init", gitContext.Project, err)
			}
		} else {
			logger.Info("plainclone2 nevyslo, error je ", err)
			return false, fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "clone", gitContext.Project, err)
		}
	}

	logger.Info("plainclone vyslo2")

	err = configureGitUser(clone)
	if err != nil {
		logger.Info("plainclone nevyslo3333")
		return false, err
	}

	_, err = clone.Head()
	if err != nil {
		logger.Info("plainclone nevyslo333344444")
		return false, fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "clone", gitContext.Project, err)
	}
	logger.Info("plainclone vyslo az dokonca")
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
	if err != nil {
		return nil, err
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
				Name:  gitKeptnUserDefault,
				Email: gitKeptnEmailDefault,
				When:  time.Now(),
			},
		})
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
				Name:  gitKeptnUserDefault,
				Email: gitKeptnEmailDefault,
				When:  time.Now(),
			},
		})
	return id.String(), err
}

func (g Git) StageAndCommitAll(gitContext common_models.GitContext, message string) (string, error) {
	//os.Setenv("SSH_AUTH_SOCK", "/run/user/1000/keyring/ssh")
	logger.Info("som v stage and commit all")
	id, err := g.commitAll(gitContext, message)
	if err != nil {
		logger.Info("som v stage and commit all error1", err)
		return "", fmt.Errorf(kerrors.ErrMsgCouldNotCommit, gitContext.Project, err)
	}
	rollbackFunc := func() {
		err := g.resetHard(gitContext)
		if err != nil {
			logger.Info("som v stage and commit all error2", err)
			logger.WithError(err).Warn("could not reset")
		}
	}
	err = g.Pull(gitContext)
	if err != nil {
		logger.Info("som v stage and commit all error3333333 ", err) //tu to pada
		rollbackFunc()
		return "", err
	}

	err = g.Push(gitContext)
	if err != nil {
		logger.Info("som v stage and commit all error4", err)
		rollbackFunc()
		return "", err
	}

	id, updated, err := g.getCurrentRemoteRevision(gitContext)
	if err != nil {
		logger.Info("som v stage and commit all error5", err)
		return "", fmt.Errorf(kerrors.ErrMsgCouldNotCommit, gitContext.Project, err)
	}
	if !updated {
		return "", fmt.Errorf(kerrors.ErrMsgCouldNotCommit, gitContext.Project, kerrors.ErrForceNeeded)
	}

	logger.Info("vsetko ok v stageandcommit")

	return id, nil
}

func (g Git) Push(gitContext common_models.GitContext) error {
	var err error
	if gitContext.Credentials == nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "push", gitContext.Project, kerrors.ErrCredentialsNotFound)
	}
	repo, _, err := g.getWorkTree(gitContext)
	if err != nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "push", gitContext.Project, err)
	}

	gitPushOptions := &git.PushOptions{
		RemoteName: "origin",
	}
	if gitContext.Credentials.PrivateKey != "" {
		publicKey, err := ssh.NewPublicKeys("git", []byte(gitContext.Credentials.PrivateKey), "<password>")
		if err != nil {
			logger.Info("publickey v push nevyslo, error je ", err)
			return err
		}
		publicKey.HostKeyCallback = ssh2.InsecureIgnoreHostKey()
		logger.Info("publickey v push vyslo")
		gitPushOptions.Auth = publicKey
	} else if gitContext.Credentials.Token != "" {
		logger.Info("!!!!!som v token push vetve")
		if gitContext.Credentials.GitProxyUrl != "" {
			logger.Info("!!!!!som v proxy push vetve")
			customClient := &nethttp.Client{
				Transport: &nethttp.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					Proxy: nethttp.ProxyURL(&url.URL{
						Scheme: gitContext.Credentials.GitProxyScheme,
						//User:   url.UserPassword(gitContext.Credentials.GitProxyUser, gitContext.Credentials.GitProxyPassword),
						Host: gitContext.Credentials.GitProxyUrl,
					}),
				},

				// 15 second timeout
				Timeout: 15 * time.Second,

				// don't follow redirect
				CheckRedirect: func(req *nethttp.Request, via []*nethttp.Request) error {
					return nethttp.ErrUseLastResponse
				},
			}

			newclient := http.NewClient(customClient)
			if newclient == nil {
				logger.Info("!!!!!newclient push zlyhal")
			}

			logger.Info("!!!!!idem instalovat push protokol")
			client.InstallProtocol("https", newclient)
			logger.Info("!!!!!nainstaloval push som protokol")

			gitPushOptions.Auth = &http.BasicAuth{
				Username: gitContext.Credentials.User,
				Password: gitContext.Credentials.Token,
			}
		} else {
			logger.Info("!!!!!som v jednoduchej token push vetve vedla proxy")
			gitPushOptions.Auth = &http.BasicAuth{
				Username: gitContext.Credentials.User,
				Password: gitContext.Credentials.Token,
			}
		}
	}

	err = repo.Push(gitPushOptions)
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		if errors.Is(err, git.ErrForceNeeded) {
			return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "push", gitContext.Project, kerrors.ErrForceNeeded)
		}
		return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "push", gitContext.Project, err)
	}
	return nil
}

func (g *Git) Pull(gitContext common_models.GitContext) error {
	if g.ProjectExists(gitContext) {
		r, w, err := g.getWorkTree(gitContext)
		if err != nil {
			return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "pull", gitContext.Project, err)
		}

		head, err := r.Head()
		if err != nil {
			return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "pull", gitContext.Project, err)
		}

		gitPullOptions := &git.PullOptions{
			RemoteName:    "origin",
			Force:         true,
			ReferenceName: head.Name(),
		}
		if gitContext.Credentials.PrivateKey != "" {
			publicKey, err := ssh.NewPublicKeys("git", []byte(gitContext.Credentials.PrivateKey), "<password>")
			if err != nil {
				logger.Info("publickey v pull nevyslo, error je ", err)
				return err
			}
			publicKey.HostKeyCallback = ssh2.InsecureIgnoreHostKey()
			logger.Info("publickey v pull vyslo")
			gitPullOptions.Auth = publicKey
		} else if gitContext.Credentials.Token != "" {
			logger.Info("!!!!!som v token pull vetve")
			if gitContext.Credentials.GitProxyUrl != "" {
				logger.Info("!!!!!som v proxy pull vetve")
				customClient := &nethttp.Client{
					Transport: &nethttp.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
						Proxy: nethttp.ProxyURL(&url.URL{
							Scheme: gitContext.Credentials.GitProxyScheme,
							//User:   url.UserPassword(gitContext.Credentials.GitProxyUser, gitContext.Credentials.GitProxyPassword),
							Host: gitContext.Credentials.GitProxyUrl,
						}),
					},

					// 15 second timeout
					Timeout: 15 * time.Second,

					// don't follow redirect
					CheckRedirect: func(req *nethttp.Request, via []*nethttp.Request) error {
						return nethttp.ErrUseLastResponse
					},
				}

				newclient := http.NewClient(customClient)
				if newclient == nil {
					logger.Info("!!!!!newclient pull zlyhal")
				}

				logger.Info("!!!!!idem instalovat pull protokol")
				client.InstallProtocol("https", newclient)
				logger.Info("!!!!!nainstaloval pull som protokol")

				gitPullOptions.Auth = &http.BasicAuth{
					Username: gitContext.Credentials.User,
					Password: gitContext.Credentials.Token,
				}
			} else {
				logger.Info("!!!!!som v jednoduchej token pull vetve vedla proxy")
				gitPullOptions.Auth = &http.BasicAuth{
					Username: gitContext.Credentials.User,
					Password: gitContext.Credentials.Token,
				}
			}
		}

		err = w.Pull(gitPullOptions)
		if err != nil && errors.Is(err, plumbing.ErrReferenceNotFound) {
			// reference not there yet
			err = w.Pull(&git.PullOptions{RemoteName: "origin", Force: true, Auth: gitPullOptions.Auth})
		}
		if err != nil {
			if errors.Is(err, git.NoErrAlreadyUpToDate) || errors.Is(err, transport.ErrEmptyRemoteRepository) {
				return nil
			} else if errors.Is(err, git.ErrNonFastForwardUpdate) {
				return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "pull", gitContext.Project, kerrors.ErrNonFastForwardUpdate)
			}
			return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "pull", gitContext.Project, err)
		}
		return nil
	}
	return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "pull", gitContext.Project, kerrors.ErrProjectNotFound)
}
func (g *Git) GetCurrentRevision(gitContext common_models.GitContext) (string, error) {
	r, _, err := g.getWorkTree(gitContext)
	if err != nil {
		return "", fmt.Errorf(kerrors.ErrMsgCouldNotGetRevision, gitContext.Project, err)
	}
	ref, err := r.Head()
	if err != nil {
		return "", fmt.Errorf(kerrors.ErrMsgCouldNotGetRevision, gitContext.Project, err)
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

func (g *Git) CreateBranch(gitContext common_models.GitContext, branch string, sourceBranch string) error {
	// move head to sourceBranch
	logger.Info("idem vytvarat brach ")
	err := g.CheckoutBranch(gitContext, sourceBranch)
	if err != nil {
		logger.Info("idem vytvarat brach error1 ", err)
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
		logger.Info("idem vytvarat brach error2 ", err)
		return fmt.Errorf(kerrors.ErrMsgCouldNotCreate, branch, gitContext.Project, err)
	}

	// First try to check out branch
	err = w.Checkout(&git.CheckoutOptions{Create: false, Force: false, Branch: b})
	if err == nil {
		logger.Info("idem vytvarat brach error3 ", err)
		return fmt.Errorf(kerrors.ErrMsgCouldNotCreate, branch, gitContext.Project, kerrors.ErrBranchExists)
	}

	if err != nil {
		// got an error  - try to create it
		logger.Info("idem vytvarat brach error4 ", err)
		if err := w.Checkout(&git.CheckoutOptions{Create: true, Force: false, Branch: b}); err != nil {
			return fmt.Errorf(kerrors.ErrMsgCouldNotCreate, branch, gitContext.Project, err)
		}
	}

	err = r.CreateBranch(newBranch)
	if err != nil {
		logger.Info("idem vytvarat brach error5 ", err)
		return fmt.Errorf(kerrors.ErrMsgCouldNotCreate, branch, gitContext.Project, err)
	}

	logger.Info("vsetko vyslo v branchi, ideme prec")

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
	})
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return fmt.Errorf(kerrors.ErrMsgCouldNotCheckout, branch, kerrors.ErrReferenceNotFound)
		}
		return fmt.Errorf(kerrors.ErrMsgCouldNotCheckout, branch, err)
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

	gitFetchOptions := &git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{"+refs/*:refs/*"},
		// <src>:<dst>, + update the reference even if it isn’t a fast-forward.
		//// take all branch from remote and put them in the local repo as origin branches and as branches
		//RefSpecs: []config.RefSpec{"+refs/heads/*:refs/remotes/origin/*", "+refs/heads/*:refs/heads/*"},
		Force: true,
	}
	if gitContext.Credentials.PrivateKey != "" {
		publicKey, err := ssh.NewPublicKeys("git", []byte(gitContext.Credentials.PrivateKey), "<password>")
		if err != nil {
			logger.Info("publickey v fetch nevyslo, error je ", err)
			return err
		}
		publicKey.HostKeyCallback = ssh2.InsecureIgnoreHostKey()
		logger.Info("publickey v fetch vyslo")
		gitFetchOptions.Auth = publicKey
	} else if gitContext.Credentials.Token != "" {
		logger.Info("!!!!!som v token fetch vetve")
		if gitContext.Credentials.GitProxyUrl != "" {
			logger.Info("!!!!!som v proxy fetch vetve")
			customClient := &nethttp.Client{
				Transport: &nethttp.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					Proxy: nethttp.ProxyURL(&url.URL{
						Scheme: gitContext.Credentials.GitProxyScheme,
						//User:   url.UserPassword(gitContext.Credentials.GitProxyUser, gitContext.Credentials.GitProxyPassword),
						Host: gitContext.Credentials.GitProxyUrl,
					}),
				},

				// 15 second timeout
				Timeout: 45 * time.Second,

				// don't follow redirect
				CheckRedirect: func(req *nethttp.Request, via []*nethttp.Request) error {
					return nethttp.ErrUseLastResponse
				},
			}

			newclient := http.NewClient(customClient)
			if newclient == nil {
				logger.Info("!!!!!newclient fetch zlyhal")
			}

			logger.Info("!!!!!idem instalovat fetch protokol")
			client.InstallProtocol("https", newclient)
			logger.Info("!!!!!nainstaloval fetch som protokol")

			gitFetchOptions.Auth = &http.BasicAuth{
				Username: gitContext.Credentials.User,
				Password: gitContext.Credentials.Token,
			}
		} else {
			logger.Info("!!!!!som v jednoduchej token fetch vetve vedla proxy")
			gitFetchOptions.Auth = &http.BasicAuth{
				Username: gitContext.Credentials.User,
				Password: gitContext.Credentials.Token,
			}
		}
	}

	if err := r.Fetch(gitFetchOptions); err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
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
		def = "master"
	}
	return def, err
}

func (g *Git) ProjectExists(gitContext common_models.GitContext) bool {
	if g.ProjectRepoExists(gitContext.Project) {
		return true
	}
	clone, err := g.CloneRepo(gitContext)
	if err != nil {
		logger.Info("clonoval som repository a error je ", err)
	}
	return clone
}

func (g *Git) ProjectRepoExists(project string) bool {
	path := GetProjectConfigPath(project)
	//logger.Info("v projectrepoexists path je %s", path)
	_, err := os.Stat(path)
	if err == nil {
		// path exists
		logger.Info("idem do plainopen")
		_, err = g.git.PlainOpen(path)
		logger.Info("error z plaininfo je ", err)
		if err == nil {
			logger.Info("plainopen dobre")
			return true
		}
	}
	logger.Info("plainopen zle")
	return false
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

func (g Git) resetHard(gitContext common_models.GitContext) error {
	r, w, err := g.getWorkTree(gitContext)
	if err != nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotGitAction, "reset", gitContext.Project, err)
	}
	revision, err := r.ResolveRevision("HEAD~1")
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
			logger.Debugf("Could not resolve commit for path %s: %s ", path, err)
			return nil, err
		}
		return resolve(t, path)
	case *object.Tag:
		target, err := o.Object()
		if err != nil {
			logger.Debugf("Could not resolve tag for path %s: %s ", path, err)
			return nil, err
		}
		return resolve(target, path)
	case *object.Tree:
		file, err := o.File(path)
		if err != nil {
			logger.Debugf("Could not resolve file for path %s: %s ", path, err)
			return nil, err
		}
		return &file.Blob, nil
	case *object.Blob:
		return o, nil
	default:
		logger.Debug("Could not resolve unsupported object for path: ", path)
		return nil, object.ErrUnsupportedObject
	}
}
