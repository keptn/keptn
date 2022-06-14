package common

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/keptn/keptn/resource-service/common_models"
	errors2 "github.com/keptn/keptn/resource-service/errors"
	logger "github.com/sirupsen/logrus"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//go:generate moq -pkg common_mock -skip-ensure -out ./fake/credential_reader_mock.go . CredentialReader
type CredentialReader interface {
	GetCredentials(project string) (*common_models.GitCredentials, error)
}

type K8sCredentialReader struct {
	k8sClient kubernetes.Interface
}

func NewK8sCredentialReader(k8sClient kubernetes.Interface) *K8sCredentialReader {
	return &K8sCredentialReader{k8sClient: k8sClient}
}

func (kr K8sCredentialReader) GetCredentials(project string) (*common_models.GitCredentials, error) {
	secretName := fmt.Sprintf("git-credentials-%s", project)

	secret, err := kr.k8sClient.CoreV1().Secrets(GetKeptnNamespace()).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		logger.Debug("Could not retrieve credentials named: ", secretName)
		return nil, errors2.ErrCredentialsNotFound
	}
	if err != nil {
		logger.Debug("Could not retrieve credentials named: ", secretName)
		return nil, err
	}

	// secret found -> unmarshal it
	credentials := &common_models.GitCredentials{}
	if err := json.Unmarshal(secret.Data["git-credentials"], credentials); err != nil {
		return nil, errors2.ErrMalformedCredentials
	}
	if err := credentials.Validate(); err != nil {
		logger.Error("Issue with credentials : ", err)
		return nil, err
	}
	return credentials, nil
}

func GetKeptnNamespace() string {
	return os.Getenv("POD_NAMESPACE")
}
