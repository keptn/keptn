package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	utils "github.com/keptn/kubernetes-utils/pkg"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/url"
	"os"
	"strings"
	"time"
)

// GitCredentials contains git ccredentials info
type GitCredentials struct {
	User      string `json:"user,omitempty"`
	Token     string `json:"token,omitempty"`
	RemoteURI string `json:"remoteURI,omitempty"`
}

var namespace = os.Getenv("POD_NAMESPACE")

// CloneRepo clones an upstream repository into a local folder "project" and returns
// whether the Git repo is already initialized.
func CloneRepo(project string, credentials GitCredentials) (bool, error) {
	uri := getRepoURI(credentials.RemoteURI, credentials.User, credentials.Token)

	msg, err := utils.ExecuteCommandInDirectory("git", []string{"clone", uri, project}, config.ConfigDir)
	const emptyRepoWarning = "warning: You appear to have cloned an empty repository."
	if strings.Contains(msg, emptyRepoWarning) {
		return false, obfuscateErrorMessage(err, &credentials)
	} else if err != nil {
		return false, obfuscateErrorMessage(err, &credentials)
	}
	return true, nil
}

func getRepoURI(uri string, user string, token string) string {
	if strings.Contains(user, "@") {
		// username contains an @, probably an e-mail; need to encode it
		// see https://stackoverflow.com/a/29356143
		user = url.QueryEscape(user)
	}

	if strings.Contains(uri, user+"@") {
		uri = strings.Replace(uri, "://"+user+"@", "://"+user+":"+token+"@", 1)
	}

	if !strings.Contains(uri, user+":"+token+"@") {
		uri = strings.Replace(uri, "://", "://"+user+":"+token+"@", 1)
	}

	return uri
}

// CheckoutBranch checks out the given branch
func CheckoutBranch(project string, branch string, disableUpstreamSync bool) error {
	projectConfigPath := config.ConfigDir + "/" + project
	_, err := utils.ExecuteCommandInDirectory("git", []string{"checkout", branch}, projectConfigPath)
	if err != nil {
		return err
	}
	if disableUpstreamSync {
		return nil
	}
	credentials, err := GetCredentials(project)
	if err == nil && credentials != nil {
		repoURI := getRepoURI(credentials.RemoteURI, credentials.User, credentials.Token)
		_, err = utils.ExecuteCommandInDirectory("git", []string{"pull", "-s", "recursive", "-X", "theirs", repoURI}, projectConfigPath)
		if err != nil {
			return obfuscateErrorMessage(err, credentials)
		}
	}
	return nil
}

// CreateBranch creates a new branch
func CreateBranch(project string, branch string, sourceBranch string) error {
	projectConfigPath := config.ConfigDir + "/" + project
	err := CheckoutBranch(project, sourceBranch, false)
	if err != nil {
		return err
	}
	_, err = utils.ExecuteCommandInDirectory("git", []string{"checkout", "-b", branch}, projectConfigPath)
	if err != nil {
		return err
	}

	// if an upstream has been defined, push the new branch
	credentials, err := GetCredentials(project)
	if err == nil && credentials != nil {
		repoURI := getRepoURI(credentials.RemoteURI, credentials.User, credentials.Token)
		_, err = utils.ExecuteCommandInDirectory("git", []string{"push", "--set-upstream", repoURI, branch}, projectConfigPath)
		if err != nil {
			return obfuscateErrorMessage(err, credentials)
		}
	}

	return nil
}

// AddOrigin adds a remote Git repository
func AddOrigin(project string) error {
	projectConfigPath := config.ConfigDir + "/" + project

	// if an upstream has been defined, add the origin and push
	credentials, err := GetCredentials(project)
	if err == nil && credentials != nil {
		repoURI := getRepoURI(credentials.RemoteURI, credentials.User, credentials.Token)
		_, err = utils.ExecuteCommandInDirectory("git", []string{"remote", "add", "origin", repoURI}, projectConfigPath)
		if err != nil {
			removeRemoteOrigin(project)
			return obfuscateErrorMessage(err, credentials)
		}

		if err := setUpstreamsAndPush(project, credentials, repoURI); err != nil {
			removeRemoteOrigin(project)
			return fmt.Errorf("failed to set upstream: %v\nKeptn requires an uninitialized repo", err)
		}
	}
	return err
}

// UpdateOrCreateOrigin tries to update the remote origin.
// If no remote origin exists, it will add one
func UpdateOrCreateOrigin(project string) error {

	projectConfigPath := config.ConfigDir + "/" + project
	credentials, err := GetCredentials(project)

	if err == nil && credentials != nil {
		repoURI := getRepoURI(credentials.RemoteURI, credentials.User, credentials.Token)

		// try to update existing remote origin
		_, err := utils.ExecuteCommandInDirectory("git", []string{"remote", "set-url", "origin", repoURI}, projectConfigPath)
		if err != nil {
			// create new remote origin in case updating was not possible
			_, err := utils.ExecuteCommandInDirectory("git", []string{"remote", "add", "origin", repoURI}, projectConfigPath)
			if err != nil {
				removeRemoteOrigin(project)
				return obfuscateErrorMessage(err, credentials)
			}
		}
		if err := setUpstreamsAndPush(project, credentials, repoURI); err != nil {
			removeRemoteOrigin(project)
			return fmt.Errorf("failed to set upstream: %v\nKeptn requires an uninitialized repo", err)
		}
	}
	return nil
}

func removeRemoteOrigin(project string) error {
	projectConfigPath := config.ConfigDir + "/" + project
	_, err := utils.ExecuteCommandInDirectory("git", []string{"remote", "remove", "origin"}, projectConfigPath)
	return err
}

func setUpstreamsAndPush(project string, credentials *GitCredentials, repoURI string) error {
	projectConfigPath := config.ConfigDir + "/" + project
	branches, err := GetBranches(project)
	if err != nil {
		return obfuscateErrorMessage(err, credentials)
	}

	defaultBranch, err := GetDefaultBranch(project)
	if err != nil {
		return obfuscateErrorMessage(err, credentials)
	}

	// first, make sure to push the master/main branch first
	for _, branch := range branches {
		if branch == defaultBranch {
			err := CheckoutBranch(project, branch, true)
			if err != nil {
				return obfuscateErrorMessage(err, credentials)
			}
			_, err = utils.ExecuteCommandInDirectory("git", []string{"push", "--set-upstream", repoURI, branch}, projectConfigPath)
			if err != nil {
				return obfuscateErrorMessage(err, credentials)
			}

			break
		}
	}

	for _, branch := range branches {
		if branch == defaultBranch {
			continue
		}
		err := CheckoutBranch(project, branch, true)
		if err != nil {
			return obfuscateErrorMessage(err, credentials)
		}
		_, err = utils.ExecuteCommandInDirectory("git", []string{"push", "--set-upstream", repoURI, branch}, projectConfigPath)
		if err != nil {
			return obfuscateErrorMessage(err, credentials)
		}
	}
	return nil
}

// StageAndCommitAll stages all current changes and commits them to the current branch
func StageAndCommitAll(project string, message string, withPull bool) error {
	projectConfigPath := config.ConfigDir + "/" + project
	_, err := utils.ExecuteCommandInDirectory("git", []string{"add", "."}, projectConfigPath)
	if err != nil {
		return err
	}

	_, err = utils.ExecuteCommandInDirectory("git", []string{"commit", "-m", message}, projectConfigPath)
	// in this case, ignore errors since the only instance when this can occur at this stage is when there is nothin to commit (no delta)
	credentials, err := GetCredentials(project)
	if err == nil && credentials != nil {
		repoURI := getRepoURI(credentials.RemoteURI, credentials.User, credentials.Token)
		if withPull {
			_, err = utils.ExecuteCommandInDirectory("git", []string{"pull", "-s", "recursive", "-X", "theirs", repoURI}, projectConfigPath)
			if err != nil {
				return obfuscateErrorMessage(err, credentials)
			}
		}
		_, err = utils.ExecuteCommandInDirectory("git", []string{"push", repoURI}, projectConfigPath)
		if err != nil {
			return obfuscateErrorMessage(err, credentials)
		}
	}
	return nil
}

func obfuscateErrorMessage(err error, credentials *GitCredentials) error {
	if err != nil && credentials != nil && credentials.Token != "" {
		errorMessage := strings.ReplaceAll(err.Error(), credentials.Token, "********")
		return errors.New(errorMessage)
	}
	return err
}

// GetCurrentVersion gets the latest version (i.e. commit hash) of the currently checked out branch
func GetCurrentVersion(project string) (string, error) {
	projectConfigPath := config.ConfigDir + "/" + project
	out, err := utils.ExecuteCommandInDirectory("git", []string{"rev-parse", "HEAD"}, projectConfigPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(out, "\n"), nil
}

// ProjectExists checks if a project exists
func ProjectExists(project string) bool {
	projectConfigPath := config.ConfigDir + "/" + project
	// check if the project exists
	_, err := os.Stat(projectConfigPath)
	// create file if not exists
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// StageExists checks if a stage in a given project exists
func StageExists(project string, stage string, disableUpstreamSync bool) bool {
	if !ProjectExists(project) {
		return false
	}
	// try to checkout the branch containing the stage config
	err := CheckoutBranch(project, stage, disableUpstreamSync)
	if err != nil {
		return false
	}
	return true
}

// ServiceExists checks if a service exists in a given stage of a project
func ServiceExists(project string, stage string, service string, disableUpstreamSync bool) bool {
	if !ProjectExists(project) {
		return false
	}
	// try to checkout the branch containing the stage config
	err := CheckoutBranch(project, stage, disableUpstreamSync)
	if err != nil {
		return false
	}
	serviceConfigPath := config.ConfigDir + "/" + project + "/" + service
	_, err = os.Stat(serviceConfigPath)
	// create file if not exists
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// GetCredentials returns the credentials for a given project, if available
func GetCredentials(project string) (*GitCredentials, error) {
	clientSet, err := getK8sClient()
	if err != nil {
		return nil, err
	}

	secret, err := clientSet.CoreV1().Secrets(namespace).Get("git-credentials-"+project, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var credentials GitCredentials
	err = json.Unmarshal(secret.Data["git-credentials"], &credentials)
	if err != nil {
		return nil, err
	}
	if credentials.User != "" && credentials.Token != "" && credentials.RemoteURI != "" {
		mv := GetProjectsMaterializedView()
		// try to update the materialized view, but continue if it does not work for some reason
		_ = mv.UpdateUpstreamInfo(project, credentials.RemoteURI, credentials.User)
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

// GetBranches returns a list of branches within the project
func GetBranches(project string) ([]string, error) {
	projectConfigPath := config.ConfigDir + "/" + project
	out, err := utils.ExecuteCommandInDirectory("git", []string{"for-each-ref", `--format=%(refname:short)`, "refs/heads/*"}, projectConfigPath)
	if err != nil {
		return nil, err
	}
	branches := strings.Split(strings.TrimSpace(out), "\n")

	return branches, nil
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

// GetDefaultBranch returns the name of the default branch of the repo
func GetDefaultBranch(project string) (string, error) {
	projectConfigPath := config.ConfigDir + "/" + project

	credentials, err := GetCredentials(project)
	if err != nil {
		return "", errors.New("could not determine default branch: " + err.Error())
	}
	if credentials != nil {
		retries := 5

		for i := 0; i < retries; i = i + 1 {
			out, err := utils.ExecuteCommandInDirectory("git", []string{"remote", "show", "origin"}, projectConfigPath)
			if err != nil {
				return "", err
			}
			lines := strings.Split(out, "\n")

			for _, line := range lines {
				if strings.Contains(line, "HEAD branch") {
					split := strings.Split(line, ":")
					if len(split) > 1 {
						defaultBranch := strings.TrimSpace(split[1])
						if defaultBranch != "(unknown)" {
							return defaultBranch, nil
						}
					}
				}
			}
			<-time.After(3 * time.Second)
		}
	}
	return "master", nil
}

func addRepoURIToMetadata(credentials *GitCredentials, metadata *models.Version) {
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
