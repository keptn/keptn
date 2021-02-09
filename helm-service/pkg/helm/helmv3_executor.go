package helm

import (
	"fmt"
	"github.com/keptn/keptn/helm-service/pkg/namespacemanager"
	"os"
	"path/filepath"
	"time"

	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"

	"helm.sh/helm/v3/pkg/release"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"

	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	keptnutils "github.com/keptn/kubernetes-utils/pkg"
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
			keptnutils.UserHomeDir(), ".kube", "config",
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
		return "", fmt.Errorf("Error when quering the manifest of chart %s in namespace %s: %s",
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
			release, err = iCli.Run(ch, vals)
		} else {
			iCli := action.NewUpgrade(cfg)
			iCli.Namespace = namespace
			iCli.Wait = true
			iCli.ResetValues = true
			iCli.Timeout = time.Minute * 3
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
		if err := keptnutils.WaitForDeploymentToBeRolledOut(getInClusterConfig(), depl.Name, depl.Namespace); err != nil {
			return fmt.Errorf("Error when waiting for deployment %s in namespace %s: %s", depl.Name, depl.Namespace, err.Error())
		}
	}
	return nil
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
