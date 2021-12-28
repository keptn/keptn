package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/keptn/keptn/resource-service/common_models"
	utils "github.com/keptn/kubernetes-utils/pkg"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
)

var ErrCouldNotReadCredentials = errors.New("could not get git credentials from client")
var ErrDecodeCredentialsError = errors.New("could not decode credentials")
var ErrNoCredentialsFound = errors.New("no credentials found")

//go:generate moq -pkg common_mock -skip-ensure -out ./fake/credential_reader_mock.go . CredentialReader
type CredentialReader interface {
	GetCredentials(project string) (*common_models.GitCredentials, error)
}

type K8sCredentialReader struct{}

func (K8sCredentialReader) GetCredentials(project string) (*common_models.GitCredentials, error) {
	clientSet, err := getK8sClient()
	if err != nil {
		return nil, ErrCouldNotReadCredentials
	}

	secretName := fmt.Sprintf("git-credentials-%s", project)

	secret, err := clientSet.CoreV1().Secrets(GetKeptnNamespace()).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		return nil, ErrNoCredentialsFound
	}
	if err != nil {
		return nil, ErrCouldNotReadCredentials
	}

	// secret found -> unmarshal it
	var credentials common_models.GitCredentials
	err = json.Unmarshal(secret.Data["git-credentials"], &credentials)
	if err != nil {
		return nil, ErrDecodeCredentialsError
	}
	if credentials.User != "" && credentials.Token != "" && credentials.RemoteURI != "" {
		return &credentials, nil
	}
	return nil, nil
}

func GetKeptnNamespace() string {
	return os.Getenv("POD_NAMESPACE")
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
