package helm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/keptn/keptn/helm-service/pkg/namespacemanager"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"

	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"

	"helm.sh/helm/v3/pkg/release"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"

	"github.com/keptn/go-utils/pkg/common/fileutils"
	"github.com/keptn/go-utils/pkg/common/kubeutils"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func getInClusterConfig() bool {
	if os.Getenv("ENVIRONMENT") == "production" {
		return true
	}
	return false
}

// HelmV3Executor provides util functions to execute helm commands
type HelmV3Executor struct {
	logger           keptncommon.LoggerInterface
	namespaceManager namespacemanager.INamespaceManager
}

// NewHelmV3Executor creates a new HelmV3Executor
func NewHelmV3Executor(logger keptncommon.LoggerInterface, nsManager namespacemanager.INamespaceManager) *HelmV3Executor {
	return &HelmV3Executor{
		logger:           logger,
		namespaceManager: nsManager,
	}
}

func (h *HelmV3Executor) newActionConfig(config *rest.Config, namespace string) (*action.Configuration, error) {

	logFunc := func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}

	restClientGetter := h.newConfigFlags(config, namespace)
	kubeClient := &kube.Client{
		Factory: cmdutil.NewFactory(restClientGetter),
		Log:     logFunc,
	}
	client, err := kubeClient.Factory.KubernetesClientSet()
	if err != nil {
		return nil, err
	}

	s := driver.NewSecrets(client.CoreV1().Secrets(namespace))
	s.Log = logFunc

	return &action.Configuration{
		RESTClientGetter: restClientGetter,
		Releases:         storage.Init(s),
		KubeClient:       kubeClient,
		Log:              logFunc,
	}, nil
}

func (h *HelmV3Executor) newConfigFlags(config *rest.Config, namespace string) *genericclioptions.ConfigFlags {
	return &genericclioptions.ConfigFlags{
		Namespace:   &namespace,
		APIServer:   &config.Host,
		CAFile:      &config.CAFile,
		BearerToken: &config.BearerToken,
	}
}

func (h *HelmV3Executor) getKubeRestConfig() (config *rest.Config, err error) {

	if getInClusterConfig() {
		config, err = rest.InClusterConfig()
	} else {
		kubeconfig := filepath.Join(
			fileutils.UserHomeDir(), ".kube", "config",
		)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return
}

// GetManifest returns the manifest for the provided release
func (h *HelmV3Executor) GetManifest(releaseName, namespace string) (string, error) {

	config, err := h.getKubeRestConfig()
	if err != nil {
		return "", err
	}
	cfg, err := h.newActionConfig(config, namespace)
	if err != nil {
		return "", err
	}
	getAction := action.NewGet(cfg)

	release, err := getAction.Run(releaseName)
	if err != nil {
		return "", fmt.Errorf("Error when querying the manifest of chart %s in namespace %s: %s",
			releaseName, namespace, err.Error())
	}
	return release.Manifest, nil
}

// UpgradeChart upgrades the provided chart and waits for all deployments
func (h *HelmV3Executor) UpgradeChart(ch *chart.Chart, releaseName, namespace string, vals map[string]interface{}) error {

	if len(ch.Templates) > 0 {
		h.logger.Info(fmt.Sprintf("Creating namespace %s if not present", namespace))
		if err := h.namespaceManager.CreateNamespaceIfNotExists(namespace); err != nil {
			return err
		}

		h.logger.Info(fmt.Sprintf("Start upgrading release %s in namespace %s", releaseName, namespace))
		config, err := h.getKubeRestConfig()
		if err != nil {
			return err
		}
		cfg, err := h.newActionConfig(config, namespace)
		if err != nil {
			return err
		}

		histClient := action.NewHistory(cfg)
		var release *release.Release

		if _, err = histClient.Run(releaseName); err == driver.ErrReleaseNotFound {
			iCli := action.NewInstall(cfg)
			iCli.Namespace = namespace
			iCli.ReleaseName = releaseName
			iCli.Wait = true
			iCli.Timeout = time.Minute * 3
			iCli.Atomic = true
			release, err = iCli.Run(ch, vals)
		} else {
			iCli := action.NewUpgrade(cfg)
			iCli.Namespace = namespace
			iCli.Wait = true
			iCli.ResetValues = true
			iCli.Timeout = time.Minute * 3
			iCli.Atomic = true
			release, err = iCli.Run(releaseName, ch, vals)
		}
		if err != nil {
			return fmt.Errorf("Error when installing/upgrading chart %s in namespace %s: %s",
				releaseName, namespace, err.Error())
		}
		if release != nil {
			h.logger.Debug(release.Manifest)
			if err := h.waitForDeploymentsOfHelmRelease(release.Manifest); err != nil {
				return err
			}
		} else {
			h.logger.Debug("Release is nil")
		}
		h.logger.Info(fmt.Sprintf("Finished upgrading chart %s in namespace %s", releaseName, namespace))
	} else {
		h.logger.Debug("Upgrade not done as this is an empty chart")
	}
	return nil
}

func (h *HelmV3Executor) waitForDeploymentsOfHelmRelease(helmManifest string) error {
	depls := GetDeployments(helmManifest)
	for _, depl := range depls {
		if err := waitForDeploymentToBeRolledOut(getInClusterConfig(), depl.Name, depl.Namespace); err != nil {
			return fmt.Errorf("Error when waiting for deployment %s in namespace %s: %s", depl.Name, depl.Namespace, err.Error())
		}
	}
	return nil
}

func waitForDeploymentToBeRolledOut(useInClusterConfig bool, deploymentName string, namespace string) error {
	clientset, err := kubeutils.GetClientSet(useInClusterConfig)
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

// UninstallRelease uninstalls the specified release in the namespace
func (h *HelmV3Executor) UninstallRelease(releaseName, namespace string) error {

	h.logger.Debug(fmt.Sprintf("Start uninstalling Helm release %s in namespace %s", releaseName, namespace))
	config, err := h.getKubeRestConfig()
	if err != nil {
		return err
	}
	cfg, err := h.newActionConfig(config, namespace)
	if err != nil {
		return err
	}

	histClient := action.NewHistory(cfg)
	if _, err = histClient.Run(releaseName); err == driver.ErrReleaseNotFound {
		h.logger.Info(fmt.Sprintf("No Helm release with name %s was found in namespace %s", releaseName, namespace))
		return nil
	}

	iCli := action.NewUninstall(cfg)
	if _, err := iCli.Run(releaseName); err != nil {
		return fmt.Errorf("Error when uninstalling release %s in namespace %s: %s",
			releaseName, namespace, err.Error())
	}
	h.logger.Debug(fmt.Sprintf("Successfully uninstall Helm release %s in namespace %s", releaseName, namespace))
	return nil
}
