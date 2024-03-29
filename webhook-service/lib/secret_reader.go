package lib

import (
	"context"
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//go:generate moq  -pkg fake -out ./fake/secret_reader_mock.go . ISecretReader
type ISecretReader interface {
	ReadSecret(name, key string) (string, error)
}

type K8sSecretReater struct {
	k8sClient kubernetes.Interface
}

func NewK8sSecretReader(k8sClient kubernetes.Interface) *K8sSecretReater {
	return &K8sSecretReater{k8sClient: k8sClient}
}

func (sr *K8sSecretReater) ReadSecret(name, key string) (string, error) {
	secret, err := sr.k8sClient.CoreV1().Secrets(GetNamespaceFromEnvVar()).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	// only allow reading from secrets that are managed by Keptn's secret-service
	if secret.Labels["app.kubernetes.io/managed-by"] != "keptn-secret-service" {
		return "", errors.New("only secrets managed by Keptn's secret-service can be referenced")
	}
	return string(secret.Data[key]), nil
}
