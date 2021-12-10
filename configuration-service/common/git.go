package common

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	config2 "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/keptn/go-utils/pkg/common/retry"
	"github.com/keptn/keptn/configuration-service/common_models"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	utils "github.com/keptn/kubernetes-utils/pkg"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var namespace = os.Getenv("POD_NAMESPACE")

const masterBranch = "master"
const mainBranch = "main"

const gitKeptnUserEnvVar = "GIT_KEPTN_USER"
const gitKeptnEmailEnvVar = "GIT_KEPTN_EMAIL"

const gitKeptnUserDefault = "keptn"
const gitKeptnEmailDefault = "keptn@keptn.sh"

//go:generate moq -pkg common_mock -skip-ensure -out ./fake/command_executor_mock.go . CommandExecutor
type CommandExecutor interface {
	ExecuteCommand(command string, args []string, directory string) (string, error)
}

//go:generate moq -pkg common_mock -skip-ensure -out ./fake/credential_reader_mock.go . CredentialReader
type CredentialReader interface {
	GetCredentials(project string) (*common_models.GitCredentials, error)
}

type KeptnUtilsCommandExecutor struct{}

func (KeptnUtilsCommandExecutor) ExecuteCommand(command string, args []string, directory string) (string, error) {
	return utils.ExecuteCommandInDirectory(command, args, directory)
}

type K8sCredentialReader struct{}

func (K8sCredentialReader) GetCredentials(project string) (*common_models.GitCredentials, error) {
	clientSet, err := getK8sClient()
	if err != nil {
		return nil, err
	}

	secretName := fmt.Sprintf("git-credentials-%s", project)

	secret, err := clientSet.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		// secret must be available since upstream repo is mandatory now
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	// secret found -> unmarshal it
	var credentials common_models.GitCredentials
	err = json.Unmarshal(secret.Data["git-credentials"], &credentials)
	if err != nil {
		return nil, err
	}
	if credentials.User != "" && credentials.Token != "" && credentials.RemoteURI != "" {
		return &credentials, nil
	}
	return nil, nil
}

type GitClient struct {
	CredentialReader CredentialReader
}

func NewGitClient() *GitClient {
	return &GitClient{CredentialReader: &K8sCredentialReader{}}
}

func (g *GitClient) ProjectExists(project string) bool {
	if g.ProjectRepoExists(project) {
		return true
	}
	// if not, try to clone
	err := g.CloneRepo(project)
	if err != nil {
		return false
	}
	return true
}

func (g *GitClient) GetCommitIdFromPath(path string) (string, error) {
	r, err := git.PlainOpen(path)
	ref, err := r.Head()
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return "", err
	}
	return commit.Hash.String(), nil
}

func (g *GitClient) GetFileByPath(path string, revision string, file string) (string, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}

	h, err := r.ResolveRevision(plumbing.Revision(revision))
	if h == nil {
		return "", errors.New("resolved nil hash " + revision)
	}

	obj, err := r.Object(plumbing.CommitObject, *h)

	if err != nil {
		return "", err
	}
	if obj == nil {
		return "", errors.New("not found")
	}
	blob, err := resolve(obj, file)

	if err != nil {
		return "", err
	}

	var re (io.Reader)
	re, err = blob.Reader()

	if err != nil {
		return "", err
	}
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(read)
	data, err := ioutil.ReadAll(re)
	if err != nil {
		return "", err
	}
	newStr := base64.StdEncoding.EncodeToString(data)
	return newStr, nil
}

// resolve blob at given path from obj. obj can be a commit, tag, tree, or blob.
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

func (g *GitClient) ProjectRepoExists(project string) bool {
	_, err := git.PlainOpen(GetProjectConfigPath(project))
	if err == nil {
		return true
	}
	return false
}

func (g *GitClient) StageAndCommitAll(project, message string) (string, error) {
	credentials, err := g.CredentialReader.GetCredentials(project)
	if err != nil {
		return "", err
	}

	commitID, err := g.CommitChanges(project, credentials, message)
	if err != nil {
		return "", err
	}

	err = retry.Retry(func() error {
		err = g.PullUpstreamChanges(project, credentials)
		if err != nil {
			logger.WithError(err).Warn("could not pull")
			return err
		}

		err = g.PushUpstreamChanges(project, credentials)
		if err != nil {
			logger.WithError(err).Warn("could not push")
			return err
		}
		return  nil
	}, retry.NumberOfRetries(5), retry.DelayBetweenRetries(1*time.Second))

	if err != nil {
		return "",err
	}
	return commitID, nil
}

func (g *GitClient) PullUpstreamChanges(project string, credentials *common_models.GitCredentials) error {

	gitCLIExecutor := NewGit(&KeptnUtilsCommandExecutor{}, &K8sCredentialReader{})

	return gitCLIExecutor.PullUpstream(project)

	//var err error
	//if credentials == nil {
	//	credentials, err = g.CredentialReader.GetCredentials(project)
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//repo, worktree, err := g.getWorkTree(project, credentials)
	//if err != nil {
	//	return err
	//}
	//
	//head, err := repo.Head()
	//if err != nil {
	//	return err
	//}
	//err = worktree.Checkout(&git.CheckoutOptions{Branch: head.Name()})
	//if err != nil {
	//	return err
	//}
	//
	//err = worktree.Pull(&git.PullOptions{
	//	ReferenceName: head.Name(),
	//	RemoteName:    "origin",
	//	Auth: &http.BasicAuth{
	//		Username: credentials.User,
	//		Password: credentials.Token,
	//	},
	//	Force: true,
	//})
	//
	//if err != nil {
	//	if errors.Is(err, git.NoErrAlreadyUpToDate) {
	//		return nil
	//	}
	//	return err
	//}
	//
	//return nil
}

func (g *GitClient) CommitChanges(project string, credentials *common_models.GitCredentials, message string) (string, error) {
	var err error
	if credentials == nil {
		credentials, err = g.CredentialReader.GetCredentials(project)
		if err != nil {
			return "", err
		}
	}
	_, workTree, err := g.getWorkTree(project, credentials)
	if err != nil {
		return "", err
	}

	err = workTree.AddWithOptions(&git.AddOptions{All: true})
	if err != nil {
		return "", err
	}
	if message == "" {
		message = "commit changes"
	}
	hash, err := workTree.Commit(message, &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  getGitKeptnUser(),
			Email: getGitKeptnEmail(),
		},
	})
	if err != nil {
		return "",err
	}

	return hash.String(), nil
}

func (g *GitClient) PushUpstreamChanges(project string, credentials *common_models.GitCredentials) error {
	var err error
	if credentials == nil {
		credentials, err = g.CredentialReader.GetCredentials(project)
		if err != nil {
			return err
		}
	}
	repo, _, err := g.getWorkTree(project, credentials)
	if err != nil {
		return err
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: credentials.User,
			Password: credentials.Token,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (g *GitClient) CloneRepo(project string) error {

	projectConfigPath := GetProjectConfigPath(project)
	// check if we already have a repository
	_, err := git.PlainOpen(projectConfigPath)
	if err == nil {
		return nil
	}

	credentials, err := g.CredentialReader.GetCredentials(project)
	if err != nil {
		return err
	}

	clone, err := git.PlainClone(GetProjectConfigPath(project), false,
		&git.CloneOptions{
			URL: credentials.RemoteURI,
			Auth: &http.BasicAuth{
				Username: credentials.User,
				Password: credentials.Token,
			},
		},
	)

	if err != nil {
		if strings.Contains(err.Error(), "empty") {
			// TODO empty remote leads to an error
			init, err := git.PlainInit(projectConfigPath, false)
			if err != nil {
				return err
			}
			clone = init
		} else {
			return err
		}
	}

	err = ConfigureGitUser(project)
	if err != nil {
		return err
	}

	_, err = clone.Head()
	if err != nil {
		// empty repository, create project metadata
		newProjectMetadata := &ProjectMetadata{
			ProjectName:               project,
			CreationTimestamp:         time.Now().String(),
			IsUsingDirectoryStructure: true,
		}

		metadataString, err := yaml.Marshal(newProjectMetadata)

		err = WriteFile(GetProjectConfigPath(project)+"/metadata.yaml", metadataString)
		if err != nil {
			return obfuscateErrorMessage(fmt.Errorf("could not write metadata.yaml during creating project %s: %w", project, err), credentials)
		}

		_, err = clone.CreateRemote(&config2.RemoteConfig{
			Name: "origin",
			URLs: []string{credentials.RemoteURI},
		})
		if err != nil {
			return err
		}
		workTree, err := clone.Worktree()
		if err != nil {
			return err
		}

		_, err = workTree.Add("metadata.yaml")
		if err != nil {
			return err
		}
		_, err = workTree.Commit("added resource", &git.CommitOptions{
			All: true,
			Author: &object.Signature{
				Name:  getGitKeptnUser(),
				Email: getGitKeptnEmail(),
			},
		})
		if err != nil {
			return err
		}
		err = clone.Push(&git.PushOptions{
			RemoteName: "origin",
			Auth: &http.BasicAuth{
				Username: credentials.User,
				Password: credentials.Token,
			}})
		if err != nil {
			if errors.Is(err, git.NoErrAlreadyUpToDate) {
				return nil
			}
			return err
		}
		return nil
	}

	return nil
}

// MigrateProject checks whether a project has already been migrated to the directory-based structure for branches. if not, it will do the migration
func (g *GitClient) MigrateProject(project string) error {
	if !g.ProjectExists(project) {
		return errors.New("project does not exist")
	}
	credentials, err := g.CredentialReader.GetCredentials(project)
	if err != nil {
		return err
	}
	err = g.PullUpstreamChanges(project, credentials)

	metadata := &ProjectMetadata{}
	metadataFileContent, err := os.ReadFile(GetProjectConfigPath(project) + "/metadata.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(metadataFileContent, metadata)
	if err != nil {
		return err
	}

	// check metadata if already migrated - if yes, no need to do anything
	if metadata.IsUsingDirectoryStructure {
		return nil
	}

	tmpProjectPath := GetTmpProjectConfigPath(project)
	// create a new repo from the previous upstream - store it in a tmp directory
	tmpClone, err := git.PlainClone(tmpProjectPath, false,
		&git.CloneOptions{
			URL: credentials.RemoteURI,
			Auth: &http.BasicAuth{
				Username: credentials.User,
				Password: credentials.Token,
			},
		},
	)

	// check out branches of the remote and store the content in the master branch of the tmp repo
	oldRepo, oldRepoWorktree, err := g.getWorkTree(project, credentials)
	if err != nil {
		return err
	}

	branches, err := oldRepo.Branches()

	err = branches.ForEach(func(branch *plumbing.Reference) error {
		err = oldRepoWorktree.Checkout(&git.CheckoutOptions{Branch: branch.Name()})
		if err != nil {
			return err
		}
		err = ensureDirectoryExists(GetTmpProjectConfigPath(project) + "/keptn-stages")
		if err != nil {
			return err
		}

		err = ensureDirectoryExists(GetTmpProjectConfigPath(project) + "/keptn-stages/" + branch.Name().Short())
		if err != nil {
			return err
		}
		//newStageMetadata := &StageMetadata{
		//	StageName:         branch.Name().String(),
		//	CreationTimestamp: time.Now().String(),
		//}
		//
		//metadataString, err := yaml.Marshal(newStageMetadata)
		//if err != nil {
		//	return err
		//}
		files, err := ioutil.ReadDir(GetProjectConfigPath(project))
		if err != nil {
			return err
		}

		for _, file := range files {
			if file.Name() == ".git" {
				continue
			}
			err := os.Rename(GetProjectConfigPath(project)+"/"+file.Name(), GetTmpProjectConfigPath(project)+"/keptn-stages/"+branch.Name().Short()+"/"+file.Name())
			if err != nil {
				return err
			}
		}
		err = oldRepoWorktree.Reset(&git.ResetOptions{Mode: git.HardReset})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	tmpWorktree, err := tmpClone.Worktree()
	err = tmpWorktree.AddWithOptions(&git.AddOptions{All: true})
	if err != nil {
		return err
	}
	_, err = tmpWorktree.Commit("migrated project structure", &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  getGitKeptnUser(),
			Email: getGitKeptnEmail(),
		},
	})
	if err != nil {
		return err
	}

	err = tmpClone.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: credentials.User,
			Password: credentials.Token,
		},
	})

	if err != nil {
		return err
	}

	// TODO: remove old repository, move migrated from tmp to /data/config, update metadata.yaml of project with isUsingDirectoryStructure = true

	return nil
}

func ensureDirectoryExists(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path, os.ModePerm); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (g *GitClient) getWorkTree(project string, credentials *common_models.GitCredentials) (*git.Repository, *git.Worktree, error) {
	projectConfigPath := GetProjectConfigPath(project)
	// check if we already have a repository
	repo, err := git.PlainOpen(projectConfigPath)
	if err != nil {
		return nil, nil, err
	}

	// check if remote matches with the credentials
	err = g.ensureRemoteMatchesCredentials(repo, credentials)
	if err != nil {
		return nil, nil, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, nil, err
	}
	return repo, worktree, nil
}

func (g *GitClient) ensureRemoteMatchesCredentials(repo *git.Repository, credentials *common_models.GitCredentials) error {
	remote, err := repo.Remote("origin")
	if err != nil {
		return err
	}
	if remote.Config().URLs[0] != credentials.RemoteURI {
		err := repo.DeleteRemote("origin")
		if err != nil {
			return err
		}
		_, err = repo.CreateRemote(&config2.RemoteConfig{
			Name: "origin",
			URLs: []string{credentials.RemoteURI},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type Git struct {
	Executor         CommandExecutor
	CredentialReader CredentialReader
}

func NewGit(e CommandExecutor, c CredentialReader) Git {
	return Git{
		Executor:         e,
		CredentialReader: c,
	}
}

func (g *Git) CloneRepo(project string) (bool, error) {
	credentials, err := g.CredentialReader.GetCredentials(project)
	if err != nil {
		return false, err
	}
	uri := getRepoURI(credentials.RemoteURI, credentials.User, credentials.Token)

	if directoryExists(GetProjectConfigPath(project)) {
		return true, nil
	}
	msg, err := g.Executor.ExecuteCommand("git", []string{"clone", uri, project}, config.ConfigDir)
	const emptyRepoWarning = "warning: You appear to have cloned an empty repository."
	if strings.Contains(msg, emptyRepoWarning) {
		_, err := g.Executor.ExecuteCommand("git", []string{"init"}, GetProjectConfigPath(project))
		if err != nil {
			return false, obfuscateErrorMessage(fmt.Errorf("could not init repository for project %s: %w", project, err), credentials)
		}
		newProjectMetadata := &ProjectMetadata{
			ProjectName:       project,
			CreationTimestamp: time.Now().String(),
		}

		metadataString, err := yaml.Marshal(newProjectMetadata)

		err = WriteFile(GetProjectConfigPath(project)+"/metadata.yaml", metadataString)
		if err != nil {
			return false, obfuscateErrorMessage(fmt.Errorf("could not write metadata.yaml during creating project %s: %w", project, err), credentials)
		}

		//if err := UpdateOrCreateOrigin(project); err != nil {
		//	return false, obfuscateErrorMessage(fmt.Errorf("could not write metadata.yaml during creation of project %s: %w", project, err), credentials)
		//}

		err = g.StageAndCommitAll(project, "Added metadata.yaml", false)
		if err != nil {
			return false, obfuscateErrorMessage(fmt.Errorf("could not write metadata.yaml during creation of project %s: %w", project, err), credentials)
		}
		return false, obfuscateErrorMessage(err, credentials)
	} else if err != nil {
		return false, obfuscateErrorMessage(err, credentials)
	}
	return true, nil
}

func (g *Git) PullUpstream(project string) error {
	credentials, err := g.CredentialReader.GetCredentials(project)
	if err == nil && credentials != nil {
		repoURI := getRepoURI(credentials.RemoteURI, credentials.User, credentials.Token)
		err = g.pullUpstreamChanges(err, repoURI, GetProjectConfigPath(project), credentials)
		if err != nil {
			return obfuscateErrorMessage(err, credentials)
		}
	}
	return nil
}

// CreateStage creates a new stage
func (g *Git) CreateStage(project string, branch string, sourceBranch string) error {
	projectConfigPath := config.ConfigDir + "/" + project
	err := g.PullUpstream(project)
	if err != nil {
		return err
	}
	_, err = g.Executor.ExecuteCommand("git", []string{"checkout", "-b", branch}, projectConfigPath)
	if err != nil {
		return err
	}

	// if an upstream has been defined, push the new branch
	credentials, err := g.CredentialReader.GetCredentials(project)
	if err == nil && credentials != nil {
		repoURI := getRepoURI(credentials.RemoteURI, credentials.User, credentials.Token)
		_, err = utils.ExecuteCommandInDirectory("git", []string{"push", "--set-upstream", repoURI, branch}, projectConfigPath)
		if err != nil {
			return obfuscateErrorMessage(err, credentials)
		}
	}

	return nil
}

// UpdateOrCreateOrigin tries to update the remote origin.
// If no remote origin exists, it will add one
func (g *Git) UpdateOrCreateOrigin(project string) error {

	projectConfigPath := config.ConfigDir + "/" + project
	credentials, err := g.CredentialReader.GetCredentials(project)

	if err == nil && credentials != nil {
		repoURI := getRepoURI(credentials.RemoteURI, credentials.User, credentials.Token)

		// try to update existing remote origin
		_, err := g.Executor.ExecuteCommand("git", []string{"remote", "set-url", "origin", repoURI}, projectConfigPath)
		if err != nil {
			// create new remote origin in case updating was not possible
			_, err := g.Executor.ExecuteCommand("git", []string{"remote", "add", "origin", repoURI}, projectConfigPath)
			if err != nil {
				err2 := removeRemoteOrigin(project)
				if err2 != nil {
					return err2
				}

				return obfuscateErrorMessage(err, credentials)
			}
		}
		if err := setUpstreamsAndPush(project, credentials, repoURI); err != nil {
			err2 := removeRemoteOrigin(project)
			if err2 != nil {
				return err2
			}
			return fmt.Errorf("failed to set upstream: %v", err)
		}
	}
	return nil
}

func (g *Git) removeRemoteOrigin(project string) error {
	projectConfigPath := config.ConfigDir + "/" + project
	_, err := g.Executor.ExecuteCommand("git", []string{"remote", "remove", "origin"}, projectConfigPath)
	return err
}

func (g *Git) setUpstreamsAndPush(project string, credentials *common_models.GitCredentials, repoURI string) error {
	projectConfigPath := config.ConfigDir + "/" + project

	defaultBranch, err := g.GetDefaultBranch(project)
	if err != nil {
		return obfuscateErrorMessage(err, credentials)
	}

	err = g.pullUpstreamChanges(err, repoURI, projectConfigPath, credentials)
	if err != nil {
		// continue if the error indicated that no remote ref HEAD has been found (e.g. in an uninitialized repo)
		if !isNoRemoteHeadFoundError(err) {
			return obfuscateErrorMessage(err, credentials)
		}
	}
	_, err = g.Executor.ExecuteCommand("git", []string{"push", "--set-upstream", repoURI, defaultBranch}, projectConfigPath)
	if err != nil {
		return obfuscateErrorMessage(err, credentials)
	}

	return nil
}

func (g *Git) pullUpstreamChanges(err error, repoURI string, projectConfigPath string, credentials *common_models.GitCredentials) error {
	_, err = g.Executor.ExecuteCommand("git", []string{"pull", "-s", "recursive", "-X", "theirs", repoURI}, projectConfigPath)
	if err != nil {
		return obfuscateErrorMessage(err, credentials)
	}
	return err
}

// StageAndCommitAll stages all current changes and commits them to the current branch
func (g *Git) StageAndCommitAll(project string, message string, withPull bool) error {
	projectConfigPath := config.ConfigDir + "/" + project
	_, err := g.Executor.ExecuteCommand("git", []string{"add", "."}, projectConfigPath)
	if err != nil {
		return err
	}

	_, err = g.Executor.ExecuteCommand("git", []string{"commit", "-m", message}, projectConfigPath)
	// in this case, ignore errors since the only instance when this can occur at this stage is when there is nothing to commit (no delta)
	credentials, err := g.CredentialReader.GetCredentials(project)
	if err == nil && credentials != nil {
		// TODO: likely we'll need a retry loop here when multiple replicas can write
		repoURI := getRepoURI(credentials.RemoteURI, credentials.User, credentials.Token)
		if withPull {
			_, err = g.Executor.ExecuteCommand("git", []string{"pull", "-s", "recursive", "-X", "theirs", repoURI}, projectConfigPath)
			if err != nil {
				return obfuscateErrorMessage(err, credentials)
			}
		}
		msg, err := g.Executor.ExecuteCommand("git", []string{"push", repoURI}, projectConfigPath)
		fmt.Println(msg)
		if err != nil {
			return obfuscateErrorMessage(err, credentials)
		}
	}
	return nil
}

// GetCurrentVersion gets the latest version (i.e. commit hash) of the currently checked out branch
func (g *Git) GetCurrentVersion(project string) (string, error) {
	projectConfigPath := config.ConfigDir + "/" + project
	out, err := g.Executor.ExecuteCommand("git", []string{"rev-parse", "HEAD"}, projectConfigPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(out, "\n"), nil
}

// GetBranches returns a list of branches within the project
func (g *Git) GetStages(project string) ([]string, error) {
	stagesConfigPath := config.ConfigDir + "/" + project + "/" + StageDirectoryName
	var stages = []string{}
	if _, err := os.Stat(stagesConfigPath); err != nil {
		if os.IsNotExist(err) {
			// in that case we simply don't have any stages
			return stages, nil
		}
		return stages, err
	}

	err := filepath.Walk(stagesConfigPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.Contains(path, ".git") {
				return nil
			}

			// TODO: for now assume that all directories in StageDirectoryName are representing a stage
			if info.IsDir() {
				stages = append(stages, strings.Trim(info.Name(), "/"))
			}
			return nil
		})
	return stages, err
}

// GetDefaultBranch returns the name of the default branch of the repo
func (g *Git) GetDefaultBranch(project string) (string, error) {
	projectConfigPath := GetProjectConfigPath(project)

	credentials, err := g.CredentialReader.GetCredentials(project)
	if err != nil {
		return "", errors.New("could not determine default branch: " + err.Error())
	}
	if credentials != nil {
		retries := 2

		for i := 0; i < retries; i = i + 1 {
			out, err := g.Executor.ExecuteCommand("git", []string{"remote", "show", "origin"}, projectConfigPath)
			if err != nil {
				return "", err
			}
			lines := strings.Split(out, "\n")

			for _, line := range lines {
				if strings.Contains(line, "HEAD branch") {
					// if we get an ambiguous HEAD, we need to fall back to master/main
					if strings.Contains(line, "remote HEAD is ambiguous") {
						branches, err := g.GetStages(project)
						if err != nil {
							return "", obfuscateErrorMessage(err, credentials)
						}
						for _, branch := range branches {
							if branch == masterBranch || branch == mainBranch {
								return branch, nil
							}
						}
					}
					split := strings.Split(line, ":")
					if len(split) > 1 {
						defaultBranch := strings.TrimSpace(split[1])
						if defaultBranch != "(unknown)" && defaultBranch != "" {
							return defaultBranch, nil
						}
					}
				}
			}
			<-time.After(3 * time.Second)
		}
	}
	return masterBranch, nil
}

func (g *Git) Reset(project string) error {
	projectConfigPath := config.ConfigDir + "/" + project
	_, err := g.Executor.ExecuteCommand("git", []string{"reset", "--hard"}, projectConfigPath)
	if err != nil {
		return err
	}
	return nil
}

func (g *Git) ConfigureGitUser(project string) error {
	projectConfigPath := config.ConfigDir + "/" + project
	_, err := g.Executor.ExecuteCommand("git", []string{"config", "user.name", getGitKeptnUser()}, projectConfigPath)
	if err != nil {
		return fmt.Errorf("could not set git user.name: %w", err)
	}
	_, err = g.Executor.ExecuteCommand("git", []string{"config", "user.email", getGitKeptnEmail()}, projectConfigPath)
	if err != nil {
		return fmt.Errorf("could not set git user.email: %w", err)
	}
	return nil
}

// ==============================

// CloneRepo clones an upstream repository into a local folder "project" and returns
// whether the Git repo is already initialized.
func CloneRepo(project string) error {
	g := NewGitClient()
	return g.CloneRepo(project)
}

func isNoRemoteHeadFoundError(err error) bool {
	return strings.Contains(err.Error(), "Couldn't find remote ref HEAD")
}

func getRepoURI(uri string, user string, token string) string {
	if strings.Contains(user, "@") {
		// username contains an @, probably an e-mail; need to encode it
		// see https://stackoverflow.com/a/29356143
		user = url.QueryEscape(user)
	}
	token = url.QueryEscape(token)
	if strings.Contains(uri, user+"@") {
		uri = strings.Replace(uri, "://"+user+"@", "://"+user+":"+token+"@", 1)
	}

	if !strings.Contains(uri, user+":"+token+"@") {
		uri = strings.Replace(uri, "://", "://"+user+":"+token+"@", 1)
	}

	return uri
}

// PullUpstream pulls changes from the upstream
func PullUpstream(project string) error {
	g := NewGitClient()
	return g.PullUpstreamChanges(project, nil)
}

//GetFileByPath returns a file resource as a string taking the latest version of it
func GetFileByPath(project string, revision string, file string) (string, error) {
	g := NewGitClient()
	return g.GetFileByPath(project, revision, file)
}

// Reset resets the current branch to the latest commit
func Reset(project string) error {
	g := NewGit(&KeptnUtilsCommandExecutor{}, &K8sCredentialReader{})
	return g.Reset(project)
}

// UpdateOrCreateOrigin tries to update the remote origin.
// If no remote origin exists, it will add one
func UpdateOrCreateOrigin(project string) error {
	g := NewGitClient()
	return g.PullUpstreamChanges(project, nil)
}

func removeRemoteOrigin(project string) error {
	g := NewGit(&KeptnUtilsCommandExecutor{}, &K8sCredentialReader{})
	return g.removeRemoteOrigin(project)
}

func setUpstreamsAndPush(project string, credentials *common_models.GitCredentials, repoURI string) error {
	g := NewGit(&KeptnUtilsCommandExecutor{}, &K8sCredentialReader{})
	return g.setUpstreamsAndPush(project, credentials, repoURI)
}

// StageAndCommitAll stages all current changes and commits them to the current branch
func StageAndCommitAll(project string, message string) (string, error) {
	g := NewGitClient()
	return g.StageAndCommitAll(project, message)
}

func obfuscateErrorMessage(err error, credentials *common_models.GitCredentials) error {
	if err != nil && credentials != nil && credentials.Token != "" {
		errorMessage := strings.ReplaceAll(err.Error(), credentials.Token, "********")
		return errors.New(errorMessage)
	}
	return err
}

// GetCurrentVersion gets the latest version (i.e. commit hash) of the currently checked out branch
func GetCurrentVersion(project string) (string, error) {
	g := NewGit(&KeptnUtilsCommandExecutor{}, &K8sCredentialReader{})
	return g.GetCurrentVersion(project)
}

// GetBranches returns a list of branches within the project
func GetBranches(project string) ([]string, error) {
	g := NewGit(&KeptnUtilsCommandExecutor{}, &K8sCredentialReader{})
	return g.GetStages(project)
}

// GetDefaultBranch returns the name of the default branch of the repo
func GetDefaultBranch(project string) (string, error) {
	g := NewGit(&KeptnUtilsCommandExecutor{}, &K8sCredentialReader{})
	return g.GetDefaultBranch(project)
}

func MigrateProject(project string) error {
	g := NewGitClient()

	return g.MigrateProject(project)
}

// ProjectExists checks if a project exists
func ProjectExists(project string) bool {
	g := NewGitClient()

	return g.ProjectExists(project)
}

func directoryExists(path string) bool {
	// check if the project exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if !info.IsDir() {
		return false
	}
	return true
}

// StageExists checks if a stage in a given project exists
func StageExists(project string, stage string) bool {
	if !ProjectExists(project) {
		return false
	}
	stageDir := GetStageConfigPath(project, stage)
	if directoryExists(stageDir) {
		return true
	}
	// if not found, try to pull the upstream
	err := PullUpstream(project)
	if err != nil {
		return false
	}
	return directoryExists(stageDir)
}

// ServiceExists checks if a service exists in a given stage of a project
func ServiceExists(project string, stage string, service string) bool {
	if !StageExists(project, stage) {
		return false
	}
	serviceDir := GetServiceConfigPath(project, stage, service)
	if directoryExists(serviceDir) {
		return true
	}
	// if not found, try to pull the upstream
	err := PullUpstream(project)
	if err != nil {
		return false
	}
	return directoryExists(serviceDir)
}

// GetCredentials returns the git upstream credentials for a given project (stored as a secret), if available
func GetCredentials(project string) (*common_models.GitCredentials, error) {
	clientSet, err := getK8sClient()
	if err != nil {
		return nil, err
	}

	secretName := fmt.Sprintf("git-credentials-%s", project)

	secret, err := clientSet.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		// if no secret was found, we just assume the user doesn't want a git upstream repo for this project
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// secret found -> unmarshal it
	var credentials common_models.GitCredentials
	err = json.Unmarshal(secret.Data["git-credentials"], &credentials)
	if err != nil {
		return nil, err
	}
	if credentials.User != "" && credentials.Token != "" && credentials.RemoteURI != "" {
		return &credentials, nil
	}
	return nil, nil
}

func getK8sClient() (*kubernetes.Clientset, error) {
	var clientSet *kubernetes.Clientset
	var useInClusterConfig bool
	if os.Getenv("env") == "production" {
		useInClusterConfig = true
	} else {
		useInClusterConfig = false
	}
	clientSet, err := utils.GetClientset(useInClusterConfig)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}

// GetResourceMetadata godoc
func GetResourceMetadata(project string) *models.Version {
	result := &models.Version{}

	credentials, err := GetCredentials(project)

	if err == nil && credentials != nil {
		addRepoURIToMetadata(credentials, result)
	}
	addVersionToMetadata(project, result)

	return result
}

// ConfigureGitUser sets the properties user.name and user.email needed for interacting with git in the given project's git repository
func ConfigureGitUser(project string) error {
	g := NewGit(&KeptnUtilsCommandExecutor{}, &K8sCredentialReader{})
	return g.ConfigureGitUser(project)
}

func addRepoURIToMetadata(credentials *common_models.GitCredentials, metadata *models.Version) {
	// the git token should not be included in the repo URI in the first place, but let's make sure it's hidden in any case
	remoteURI := credentials.RemoteURI
	remoteURI = strings.Replace(remoteURI, credentials.Token, "********", -1)
	metadata.UpstreamURL = remoteURI
}

func addVersionToMetadata(project string, metadata *models.Version) {
	version, err := GetCurrentVersion(project)
	if err == nil {
		metadata.Version = version
	}
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
