package kube

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/go-version"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	kubeServerVersion  = regexp.MustCompile(`Server Version: version\.Info{Major:"(\d+)", Minor:"(\d+.*?)\+{0,1}"`)
	executeCommandFunc = keptnutils.ExecuteCommand
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/namespace_handler.go . IKeptnNamespaceHandler
type IKeptnNamespaceHandler interface {
	ExistsNamespace(useInClusterConfig bool, namespace string) (bool, error)
	CreateNamespace(useInClusterConfig bool, namespace string, namespaceMetadata ...metav1.ObjectMeta) error
	PatchKeptnManagedNamespace(useInClusterConfig bool, namespace string) error
}

type KubernetesUtilsKeptnNamespaceHandler struct{}

func NewKubernetesUtilsKeptnNamespaceHandler() *KubernetesUtilsKeptnNamespaceHandler {
	return &KubernetesUtilsKeptnNamespaceHandler{}
}

func (KubernetesUtilsKeptnNamespaceHandler) ExistsNamespace(useInClusterConfig bool, namespace string) (bool, error) {
	return keptnutils.ExistsNamespace(useInClusterConfig, namespace)
}

func (KubernetesUtilsKeptnNamespaceHandler) CreateNamespace(useInClusterConfig bool, namespace string, namespaceMetadata ...metav1.ObjectMeta) error {
	return keptnutils.CreateNamespace(useInClusterConfig, namespace, namespaceMetadata...)
}

func (KubernetesUtilsKeptnNamespaceHandler) PatchKeptnManagedNamespace(useInClusterConfig bool, namespace string) error {
	return keptnutils.PatchKeptnManagedNamespace(useInClusterConfig, namespace)
}

// IsKubectlAvailable checks whether kubectl is available
func IsKubectlAvailable() (bool, error) {

	_, err := executeCommandFunc("kubectl", []string{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func getKubeServerVersion() (string, error) {

	out, err := executeCommandFunc("kubectl", []string{"version"})
	if err != nil {
		return "", err
	}
	submatches := kubeServerVersion.FindStringSubmatch(out)
	if submatches == nil {
		return "", errors.New("Server Version not found: " + out)
	}
	return submatches[1] + "." + submatches[2], nil
}

// CheckKubeServerVersion checks the Kubernetes Server version against the given constraints and returns true if the version is higher than the given constraints
func CheckKubeServerVersion(constraints string) (bool, error) {

	serverVersion, err := getKubeServerVersion()
	if err != nil {
		return false, err
	}
	hVersion, err := version.NewVersion(serverVersion)
	if err != nil {
		return false, err
	}
	hConstraints, err := version.NewConstraint(constraints)
	if err != nil {
		return false, err
	}

	isNewerVersion, allConstraintsValid := checkIsNewVersion(hConstraints, hVersion)

	if allConstraintsValid {
		return false, nil
	}

	return isNewerVersion, fmt.Errorf("The Kubernetes Server Version '%s' doesn't satisfy constraints '%s'", serverVersion, constraints)
}

// checkIsNewVersion checks if the provided version is higher than all constraints and if all constraints are valid
func checkIsNewVersion(hConstraints version.Constraints, serverVersion *version.Version) (bool, bool) {
	isNewerVersion := false
	allConstraintsValid := true
	for _, constraint := range hConstraints {
		validConstraint := constraint.Check(serverVersion)
		if !validConstraint {
			allConstraintsValid = false
			constraintVersionString := getConstraintVersion(constraint.String())

			constraintVersion, err := version.NewVersion(constraintVersionString)
			if err != nil {
				return false, false
			}
			if serverVersion.GreaterThan(constraintVersion) {
				isNewerVersion = true
			} else {
				return false, false
			}
		}
	}
	return isNewerVersion, allConstraintsValid
}

// getConstraintVersion returns the version of a constraint without leading spaces, <, >, =
func getConstraintVersion(constraint string) string {
	for index, character := range constraint {
		if character != '<' && character != '>' && character != ' ' && character != '=' {
			return constraint[index:]
		}
	}
	return constraint
}

// CheckDeploymentManagedByHelm implements a naive check if the deployment with the given name in the given namespace
// was installed by helm by checking if the label "app.kubernetes.io/managed-by=Helm" is present on the deployment
func CheckDeploymentManagedByHelm(deploymentName, namespace string) (bool, error) {
	errstr := "Failed to check if deployment %s in namespace %s is managed by Helm: %s"

	type CmdResponse struct {
		Metadata struct {
			Labels map[string]string `json:"labels"`
		} `json:"metadata"`
	}

	out, err := executeCommandFunc("kubectl", []string{"get", "deployments", deploymentName, "-n", namespace, "-o", "json"})
	if err != nil {
		return false, fmt.Errorf(errstr, deploymentName, namespace, err.Error())
	}
	var response CmdResponse
	if err = json.Unmarshal([]byte(out), &response); err != nil {
		return false, fmt.Errorf(errstr, deploymentName, namespace, err.Error())
	}

	if value, keyExists := response.Metadata.Labels["app.kubernetes.io/managed-by"]; keyExists {
		if value == "Helm" {
			return true, nil
		}
	}

	return false, nil
}

// CheckDeploymentAvailable implements a check whether a deployment with the given name in the given namespace exists
func CheckDeploymentAvailable(deploymentName, namespace string) (bool, error) {

	type Metadata struct {
		Name string `json:"name"`
	}
	type Item struct {
		Metadata Metadata `json:"metadata"`
	}
	type CmdResponse struct {
		Items []Item `json:"items"`
	}

	errstr := "Failed to check if deployment %s is available in namespace %s: %s"
	out, err := executeCommandFunc("kubectl", []string{"get", "deployments", "-n", namespace, "-o", "json"})
	if err != nil {
		return false, fmt.Errorf(errstr, deploymentName, namespace, err.Error())
	}

	var response CmdResponse
	if err = json.Unmarshal([]byte(out), &response); err != nil {
		return false, fmt.Errorf(errstr, deploymentName, namespace, err.Error())
	}

	for _, item := range response.Items {
		if item.Metadata.Name == deploymentName {
			return true, nil
		}
	}
	return false, nil
}
