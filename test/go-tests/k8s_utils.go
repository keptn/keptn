package go_tests

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/keptn/go-utils/pkg/common/kubeutils"
	keptn2 "github.com/keptn/go-utils/pkg/lib"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

func SetEnvVarsOfDeployment(deploymentName string, containerName string, envVars []v1.EnvVar) error {
	clientset, err := kubeutils.GetClientset(false)
	if err != nil {
		return err
	}
	depl, err := clientset.AppsV1().Deployments(GetKeptnNameSpaceFromEnv()).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	for index, container := range depl.Spec.Template.Spec.Containers {
		if "distributor" == container.Name {
			for _, e := range envVars {
				replaced := false
				for ii, existingEnvVar := range depl.Spec.Template.Spec.Containers[index].Env {
					// if we find an already existing env war with the same name, we need to replace it
					if existingEnvVar.Name == e.Name {
						depl.Spec.Template.Spec.Containers[index].Env[ii] = e
						replaced = true
						break
					}
				}
				// if we did not replace an env var, we need to append it
				if !replaced {
					depl.Spec.Template.Spec.Containers[index].Env = append(depl.Spec.Template.Spec.Containers[index].Env, e)
					replaced = false
				}
			}
		}
	}

	_, err = clientset.AppsV1().Deployments(GetKeptnNameSpaceFromEnv()).Update(context.TODO(), depl, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return kubeutils.WaitForDeploymentToBeRolledOut(false, deploymentName, GetKeptnNameSpaceFromEnv())
}

func GetImageOfDeploymentContainer(deploymentName, containerName string) (string, error) {
	clientset, err := kubeutils.GetClientset(false)
	if err != nil {
		return "", err
	}
	depl, err := clientset.AppsV1().Deployments(GetKeptnNameSpaceFromEnv()).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	for _, container := range depl.Spec.Template.Spec.Containers {
		if containerName == container.Name {
			return container.Image, nil
		}
	}
	return "", fmt.Errorf("container %s not found in deployment %s", containerName, deploymentName)
}

func SetImageOfDeploymentContainer(deploymentName, containerName, image string) error {
	clientset, err := kubeutils.GetClientset(false)
	if err != nil {
		return err
	}

	depl, err := clientset.AppsV1().Deployments(GetKeptnNameSpaceFromEnv()).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	for index, container := range depl.Spec.Template.Spec.Containers {
		if containerName == container.Name {
			depl.Spec.Template.Spec.Containers[index].Image = image
			depl.Spec.Template.Spec.Containers[index].ImagePullPolicy = "Always"
		}
	}
	_, err = clientset.AppsV1().Deployments(GetKeptnNameSpaceFromEnv()).Update(context.TODO(), depl, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return kubeutils.WaitForDeploymentToBeRolledOut(false, deploymentName, GetKeptnNameSpaceFromEnv())
}

type WaitForDeploymentOptions struct {
	WithImageName string
}

func WaitAndCheckDeployment(deploymentName, namespace string, timeout time.Duration, options WaitForDeploymentOptions) error {
	clientset, _ := kubeutils.GetClientset(false)
	return wait.PollImmediate(time.Second*3, timeout, checkDeployment(clientset, deploymentName, namespace, options))
}

func checkDeployment(client *kubernetes.Clientset, deploymentName, namespace string, options WaitForDeploymentOptions) wait.ConditionFunc {
	return func() (bool, error) {
		deployment, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		// check replicas count
		if deployment.Status.ReadyReplicas != *(deployment.Spec.Replicas) {
			return false, nil
		}

		// check for image name
		if options.WithImageName != "" {
			containerNameOfDeployment := deployment.Spec.Template.Spec.Containers[0].Image
			if !(containerNameOfDeployment == options.WithImageName) {
				return false, nil
			}
		}
		return true, nil
	}
}

func WaitForURL(url string, timeout time.Duration) error {
	return wait.PollImmediate(time.Second*3, timeout, checkURL(url))
}

func checkURL(url string) wait.ConditionFunc {
	return func() (bool, error) {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		_, err := http.Get(url)
		if err != nil {
			return false, nil
		}
		return true, nil
	}
}

func GetFromConfigMap(namespace string, configMapName string, getDataByKeyFn func(data map[string]string) string) (string, error) {
	client, _ := kubeutils.GetClientset(false)
	cm, err := client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return getDataByKeyFn(cm.Data), nil
}

func UpdateConfigMap(namespace string, configMapName string, replaceConfig func(cm *v1.ConfigMap)) error {
	client, _ := kubeutils.GetClientset(false)
	cm, err := client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	replaceConfig(cm)
	_, err = client.CoreV1().ConfigMaps(namespace).Update(context.TODO(), cm, metav1.UpdateOptions{})
	return err
}

func PutConfigMapDataVal(namespace string, configMapName string, key string, val string) error {
	return UpdateConfigMap(namespace, configMapName, func(cm *v1.ConfigMap) {
		cm.Data[key] = val
	})
}

// WaitForDeploymentInNamespace
// deprecated, use WaitAndCheckDeployment
func WaitForDeploymentInNamespace(deploymentName, namespace string) error {
	return kubeutils.WaitForDeploymentToBeRolledOut(false, deploymentName, namespace)
}

func WaitForPodOfDeployment(deploymentName string) error {
	return kubeutils.WaitForDeploymentToBeRolledOut(false, deploymentName, GetKeptnNameSpaceFromEnv())
}

type K8SEvent struct {
	Reason        string    `json:"reason"`
	Type          string    `json:"type"`
	Message       string    `json:"message"`
	LastTimestamp time.Time `json:"lastTimestamp"`
}

type K8SEventArray struct {
	Items []K8SEvent `json:"items"`
}

func GetOOMEvents() (K8SEventArray, error) {
	events, err := keptn2.ExecuteCommand(kubectlExecutable, []string{"get", "events", "--sort-by=’.lastTimestamp’", "-n=default", "-o=json"})
	TimeInterval := time.Now().Add(-1 * time.Hour)
	if err != nil {
		return K8SEventArray{}, err
	}

	eventArray := K8SEventArray{}
	err = json.Unmarshal([]byte(events), &eventArray)

	if err != nil {
		return K8SEventArray{}, err
	}

	oomEvents := K8SEventArray{
		Items: []K8SEvent{},
	}
	for _, event := range eventArray.Items {
		if event.LastTimestamp.Before(TimeInterval) {
			break
		}
		if strings.Contains(event.Reason, "OOM") {
			oomEvents.Items = append(oomEvents.Items, event)
		}
	}
	return oomEvents, err
}

func CompareServiceNameWithDeploymentName(serviceName string, deploymentName string) (bool, error) {
	api, err := kubeutils.GetKubeAPI(false)
	if err != nil {
		return false, err
	}

	service, err := api.Services(GetKeptnNameSpaceFromEnv()).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	if service.Spec.Selector["app.kubernetes.io/component"] == deploymentName {
		return true, nil
	}

	return false, nil
}
