package common

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/keptn/keptn/configuration-service/config"
	utils "github.com/keptn/kubernetes-utils/pkg"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GitCredentials contains git ccredentials info
type GitCredentials struct {
	User      string `json:"user,omitempty"`
	Token     string `json:"token,omitempty"`
	RemoteURI string `json:"remoteURI,omitempty"`
}

// CloneRepo clones an upstream repository into a local folder "project"
func CloneRepo(project string, user string, token string, uri string) error {
	uri = getRepoURI(uri, user, token)

	_, err := utils.ExecuteCommandInDirectory("git", []string{"clone", uri, project}, config.ConfigDir)

	return err
}

func getRepoURI(uri string, user string, token string) string {
	if strings.Contains(user, "@") {
		// username contains an @, probably an e-mail; need to encode it
		// see https://stackoverflow.com/a/29356143
		user = url.QueryEscape(user)
	}

	if strings.Contains(uri, user+"@") {
		uri = strings.Replace(uri, "https://"+user+"@", "https://"+user+":"+token+"@", 1)
	}

	if !strings.Contains(uri, user+":"+token+"@") {
		uri = strings.Replace(uri, "https://", "https://"+user+":"+token+"@", 1)
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
			return err
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
			return errors.New("Could not push to upstream")
		}
	}

	return nil
}

// StageAndCommitAll stages all current changes and commits them to the current branch
func StageAndCommitAll(project string, message string) error {
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
		_, err = utils.ExecuteCommandInDirectory("git", []string{"pull", "-s", "recursive", "-X", "theirs", repoURI}, projectConfigPath)
		if err != nil {
			return err
		}
		_, err = utils.ExecuteCommandInDirectory("git", []string{"push", repoURI}, projectConfigPath)
		if err != nil {
			return err
		}
	}
	return nil
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

// StoreGitCredentials stores the specified git credentials as a secret in the cluster
func StoreGitCredentials(project string, user string, token string, remoteURI string) error {

	clientSet, err := getK8sClient()
	if err != nil {
		return err
	}

	credentials := &GitCredentials{
		User:      user,
		Token:     token,
		RemoteURI: remoteURI,
	}

	credsEncoded, err := json.Marshal(credentials)
	if err != nil {
		return err
	}
	secret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "git-credentials-" + project,
			Namespace: "keptn",
		},
		Data: map[string][]byte{
			"git-credentials": credsEncoded,
		},
		Type: "Opaque",
	}
	_, err = clientSet.CoreV1().Secrets("keptn").Create(secret)
	if err != nil {
		return err
	}
	return nil
}

// GetCredentials returns the credentials for a given project, if available
func GetCredentials(project string) (*GitCredentials, error) {
	clientSet, err := getK8sClient()
	if err != nil {
		return nil, err
	}

	secret, err := clientSet.CoreV1().Secrets("keptn").Get("git-credentials-"+project, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var credentials GitCredentials
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

// DeleteCredentials deletes the credentials of a given project
func DeleteCredentials(project string) error {
	clientSet, err := getK8sClient()
	if err != nil {
		return err
	}

	err = clientSet.CoreV1().Secrets("keptn").Delete("git-credentials-"+project, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

// GetBranches returns a list of branches within the project
func GetBranches(project string) ([]string, error) {
	projectConfigPath := config.ConfigDir + "/" + project
	out, err := utils.ExecuteCommandInDirectory("git", []string{"for-each-ref", `--format=%(refname:short)`, "refs/heads/*"}, projectConfigPath)
	if err != nil {
		return nil, err
	}
	branches := strings.Split(out, "\n")

	return branches, nil
}
