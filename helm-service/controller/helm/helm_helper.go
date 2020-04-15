package helm

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	keptnevents "github.com/keptn/go-utils/pkg/lib"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"

	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func getInClusterConfig() bool {
	if os.Getenv("ENVIRONMENT") == "production" {
		return true
	}
	return false
}

// HelmExecutor provides util functions to execute helm commands
type HelmExecutor struct {
	logger keptnevents.LoggerInterface
}

// NewHelmExecutor creates a new HelmExecutor
func NewHelmExecutor(logger keptnevents.LoggerInterface) *HelmExecutor {
	return &HelmExecutor{logger: logger}
}

func (h *HelmExecutor) newActionConfig(config *rest.Config, namespace string) (*action.Configuration, error) {

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

func (h *HelmExecutor) newConfigFlags(config *rest.Config, namespace string) *genericclioptions.ConfigFlags {
	return &genericclioptions.ConfigFlags{
		Namespace:   &namespace,
		APIServer:   &config.Host,
		CAFile:      &config.CAFile,
		BearerToken: &config.BearerToken,
	}
}

func (h *HelmExecutor) getKubeRestConfig() (config *rest.Config, err error) {

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
func (h *HelmExecutor) GetManifest(releaseName, namespace string) (string, error) {

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
func (h *HelmExecutor) UpgradeChart(ch *chart.Chart, releaseName, namespace string, vals map[string]interface{}) error {

	if len(ch.Templates) > 0 {
		h.logger.Info(fmt.Sprintf("Start upgrading chart %s in namespace %s", releaseName, namespace))
		config, err := h.getKubeRestConfig()
		if err != nil {
			return err
		}
		cfg, err := h.newActionConfig(config, namespace)
		if err != nil {
			return err
		}

		iCli := action.NewUpgrade(cfg)
		iCli.Namespace = namespace
		iCli.Wait = true
		iCli.Install = true
		iCli.ResetValues = true
		u, err := iCli.Run(releaseName, ch, vals)
		if err != nil {
			return fmt.Errorf("Error when upgrading chart %s in namespace %s: %s",
				releaseName, namespace, err.Error())
		}
		h.logger.Debug(u.Manifest)

		if err := h.waitForDeploymentsOfHelmRelease(u.Manifest); err != nil {
			return err
		}
		h.logger.Info(fmt.Sprintf("Finished upgrading chart %s in namespace %s", releaseName, namespace))
	} else {
		h.logger.Debug("Upgrade not done as this is an empty chart")
	}
	return nil
}

func (h *HelmExecutor) waitForDeploymentsOfHelmRelease(helmManifest string) error {
	depls := GetDeployments(helmManifest)
	for _, depl := range depls {
		if err := keptnutils.WaitForDeploymentToBeRolledOut(getInClusterConfig(), depl.Name, depl.Namespace); err != nil {
			return fmt.Errorf("Error when waiting for deployment %s in namespace %s: %s", depl.Name, depl.Namespace, err.Error())
		}
	}
	return nil
}

// GetServices returns all services contained in the Helm manifest
func GetServices(helmManifest string) []*corev1.Service {

	services := []*corev1.Service{}
	dec := kyaml.NewYAMLToJSONDecoder(strings.NewReader(helmManifest))
	for {
		var svc corev1.Service
		err := dec.Decode(&svc)
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if keptnutils.IsService(&svc) {
			services = append(services, &svc)
		}
	}

	return services
}

// GetDeployments returns all deployments contained in the Helm manifest
func GetDeployments(helmManifest string) []*appsv1.Deployment {

	deployments := []*appsv1.Deployment{}
	dec := kyaml.NewYAMLToJSONDecoder(strings.NewReader(helmManifest))
	for {
		var dpl appsv1.Deployment
		err := dec.Decode(&dpl)
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if keptnutils.IsDeployment(&dpl) {
			deployments = append(deployments, &dpl)
		}
	}
	return deployments
}
