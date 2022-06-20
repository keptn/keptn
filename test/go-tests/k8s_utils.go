package go_tests

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"

	appsv1 "k8s.io/api/apps/v1"

	apierr "k8s.io/apimachinery/pkg/api/errors"

	typesv1 "k8s.io/api/core/v1"
)

func SetEnvVarsOfDeployment(deploymentName string, containerName string, envVars []v1.EnvVar) error {
	clientset, err := GetClientset(false)
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

	return WaitForDeploymentToBeRolledOut(false, deploymentName, GetKeptnNameSpaceFromEnv())
}

func GetImageOfDeploymentContainer(deploymentName, containerName string) (string, error) {
	clientset, err := GetClientset(false)
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
	clientset, err := GetClientset(false)
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

	return WaitForDeploymentToBeRolledOut(false, deploymentName, GetKeptnNameSpaceFromEnv())
}

type WaitForDeploymentOptions struct {
	WithImageName string
}

func WaitAndCheckDeployment(deploymentName, namespace string, timeout time.Duration, options WaitForDeploymentOptions) error {
	clientset, _ := GetClientset(false)
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
	client, _ := GetClientset(false)
	cm, err := client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return getDataByKeyFn(cm.Data), nil
}

func UpdateConfigMap(namespace string, configMapName string, replaceConfig func(cm *v1.ConfigMap)) error {
	client, _ := GetClientset(false)
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
	return WaitForDeploymentToBeRolledOut(false, deploymentName, namespace)
}

func WaitForPodOfDeployment(deploymentName string) error {
	return WaitForDeploymentToBeRolledOut(false, deploymentName, GetKeptnNameSpaceFromEnv())
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
	events, err := ExecuteKubeCommand(kubectlExecutable, []string{"get", "events", "--sort-by=’.lastTimestamp’", "-n=default", "-o=json"})
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
	api, err := GetKubeAPI(false)
	if err != nil {
		return false, err
	}

	service, err := api.Services(GetKeptnNameSpaceFromEnv()).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	if service.Spec.Selector["app.kubernetes.io/name"] == deploymentName {
		return true, nil
	}

	return false, nil
}

// GetClientset returns the kubernetes Clientset
func GetClientset(useInClusterConfig bool) (*kubernetes.Clientset, error) {

	var config *rest.Config
	var err error
	if useInClusterConfig {
		config, err = rest.InClusterConfig()
	} else {
		var kubeconfig string
		if os.Getenv("KUBECONFIG") != "" {
			kubeconfig = ExpandTilde(os.Getenv("KUBECONFIG"))
		} else {
			kubeconfig = filepath.Join(
				UserHomeDir(), ".kube", "config",
			)
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

// GetKubeAPI returns the CoreV1Interface
func GetKubeAPI(useInClusterConfig bool) (corev1.CoreV1Interface, error) {

	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1(), nil
}

// UserHomeDir returns the HOME directory by taking into account the operating system
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// ExpandTilde expands ~ to HOME
func ExpandTilde(fileName string) string {
	if fileName == "~" {
		return UserHomeDir()
	} else if strings.HasPrefix(fileName, "~/") {
		return filepath.Join(UserHomeDir(), fileName[2:])
	}
	return fileName
}

// WaitForDeploymentToBeRolledOut waits until the deployment is Available
func WaitForDeploymentToBeRolledOut(useInClusterConfig bool, deploymentName string, namespace string) error {
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return err
	}

	const maxWaitForDeploymentRetries = 90
	deployment, err := getDeployment(clientset, namespace, deploymentName)
	retries := 0
	for {

		var cond *appsv1.DeploymentCondition

		for i := range deployment.Status.Conditions {
			c := deployment.Status.Conditions[i]
			if c.Type == appsv1.DeploymentProgressing {
				cond = &c
				break
			}
		}

		if cond != nil && cond.Reason == "ProgressDeadlineExceeded" {
			return fmt.Errorf("Deployment %q exceeded its progress deadline", deployment.Name)
		}
		if !(deployment.Spec.Replicas != nil && deployment.Status.UpdatedReplicas < *deployment.Spec.Replicas ||
			deployment.Status.Replicas > deployment.Status.UpdatedReplicas ||
			deployment.Status.AvailableReplicas < deployment.Status.UpdatedReplicas) {
			return nil
		}

		time.Sleep(2 * time.Second)
		deployment, err = getDeployment(clientset, namespace, deploymentName)
		if err != nil {
			return err
		}
		retries = retries + 1
		if retries >= maxWaitForDeploymentRetries {
			return fmt.Errorf("Timed out waiting for deployment %q", deployment.Name)
		}
	}
}

func getDeployment(clientset *kubernetes.Clientset, namespace string, deploymentName string) (*appsv1.Deployment, error) {
	dep, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil &&
		strings.Contains(err.Error(), "the object has been modified; please apply your changes to the latest version and try again") {
		time.Sleep(10 * time.Second)
		return clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	}
	return dep, nil
}

// ExecuteCommand exectues the command using the args
func ExecuteKubeCommand(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("Error executing command %s %s: %s\n%s", command, strings.Join(args, " "), err.Error(), string(out))
	}
	return string(out), nil
}

// ExistsNamespace checks whether a namespace with the provided name exists
func ExistsNamespace(useInClusterConfig bool, namespace string) (bool, error) {
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return false, err
	}
	_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*apierr.StatusError); ok && statusErr.ErrStatus.Reason == metav1.StatusReasonNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetKeptnEndpointFromService returns the loadbalancer service IP from Keptn Installation
func GetKeptnEndpointFromService(useInClusterConfig bool, namespace string, serviceName string) (string, error) {
	var keptnService *typesv1.Service
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return "", err
	}
	keptnService, err = clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	switch keptnService.Spec.Type {
	case "LoadBalancer":
		if len(keptnService.Status.LoadBalancer.Ingress) > 0 {
			return keptnService.Status.LoadBalancer.Ingress[0].IP, nil
		}
		return "", fmt.Errorf("Loadbalancer IP isn't found")
	default:
		return "", fmt.Errorf("It doesn't support ClusterIP & NodePort type service for fetching endpoint automatically")
	}
}

// GetKeptnAPITokenFromSecret returns the `keptn-api-token` data secret from Keptn Installation
func GetKeptnAPITokenFromSecret(useInClusterConfig bool, namespace string, secretName string) (string, error) {
	var keptnSecret *typesv1.Secret
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return "", err
	}
	keptnSecret, err = clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	if apitoken, ok := keptnSecret.Data["keptn-api-token"]; ok {
		return string(apitoken), nil
	}
	return "", fmt.Errorf("data 'keptn-api-token' not found")
}

func ScaleDeployment(useInClusterConfig bool, deployment string, namespace string, replicas int32) error {
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return err
	}
	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := deploymentsClient.Get(context.TODO(), deployment, metav1.GetOptions{})
		if getErr != nil {
			return fmt.Errorf("Failed to get latest version of Deployment: %v", getErr)
		}

		result.Spec.Replicas = int32Ptr(replicas)
		_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	return retryErr
}

func int32Ptr(i int32) *int32 { return &i }

// RestartPodsWithSelector restarts the pods which are found in the provided namespace and selector
func RestartPodsWithSelector(useInClusterConfig bool, namespace string, selector string) error {
	clientset, err := GetKubeAPI(useInClusterConfig)
	if err != nil {
		return err
	}
	pods, err := clientset.Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return err
	}
	for _, pod := range pods.Items {
		if err := clientset.Pods(namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{}); err != nil {
			return err
		}
	}
	return nil
}
